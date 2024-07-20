package grpc

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/thetherington/jobber-common/models/chat"
	pb "github.com/thetherington/jobber-common/protogen/go/chat"
	"google.golang.org/protobuf/types/known/emptypb"
)

type sub struct {
	stream   pb.ChatService_SubscribeServer // stream is the server side of the RPC stream
	finished chan<- bool                    // finished is used to signal closure of a client subscribing goroutine
}

// Subscribe handles a subscribe request from a client
func (g *GrpcAdapter) Subscribe(req *pb.RequestWithParam, stream pb.ChatService_SubscribeServer) error {
	// Handle subscribe request
	slog.Info("Received subscribe request from gateway", "id", req.Param)

	fin := make(chan bool)

	// Save the subscriber stream according to the given client Id
	g.subscribers.Store(req.Param, sub{stream: stream, finished: fin})

	ctx := stream.Context()

	// Keep this scope alive because once this scope exists - the stream is closed
	for {
		select {
		case <-fin:
			slog.Info("Closing stream for gateway client", "id", req.Param)
			return nil

		case <-ctx.Done():
			slog.Warn("Gateway client has disconnected", "id", req.Param)
			g.subscribers.Delete(req.Param)

			return nil
		}
	}
}

// Unsubscribe handles a unsubscribe request from a client
func (g *GrpcAdapter) Unsubscribe(ctx context.Context, req *pb.RequestWithParam) (*emptypb.Empty, error) {
	v, ok := g.subscribers.Load(req.Param)
	if !ok {
		return nil, serviceError(fmt.Errorf("failed to load subscriber key: %s", req.Param))
	}

	sub, ok := v.(sub)
	if !ok {
		return nil, serviceError(fmt.Errorf("failed to cast subscriber value: %T", v))
	}

	select {
	case sub.finished <- true:
		slog.Info("Unsubscribed client", "id", req.Param)
	default:
		// Default case is to avoid blocking in case client has already unsubscribed
	}

	g.subscribers.Delete(req.Param)

	return &emptypb.Empty{}, nil
}

func (g *GrpcAdapter) PushMessage(cmd string, message *chat.MessageDocument) {
	var unsubscribe []string

	// Iterate over all subscribers and send data to each client
	g.subscribers.Range(func(key, value any) bool {
		id, ok := key.(string)
		if !ok {
			slog.Error("failed to cast subscriber key", "key", key)
			return false
		}

		sub, ok := value.(sub)
		if !ok {
			slog.Error("failed to cast subscriber value", "value", value)
		}

		// Send data over gRPC stream to client
		err := sub.stream.Send(&pb.MessageResponse{
			Message:     cmd,
			MessageData: chat.CreateProtoMessageDocument(message),
		})
		if err != nil {
			slog.With("error", err).Error("failed to send data to client", "id", id)

			select {
			case sub.finished <- true:
				slog.Info("Unsubscribe client", "id", id)
			default:
				// Default case is to avoid blocking in case client has already unsubscribed
			}

			unsubscribe = append(unsubscribe, id)
		}

		return true
	})

	// Unsubscribe erroneous client streams
	for _, id := range unsubscribe {
		g.subscribers.Delete(id)
	}
}
