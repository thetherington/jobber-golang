package grpc

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/thetherington/jobber-common/models/order"
	pb "github.com/thetherington/jobber-common/protogen/go/order"
	"github.com/thetherington/jobber-common/utils"
	"google.golang.org/protobuf/types/known/emptypb"
)

type subNotify struct {
	stream   pb.NotificationService_SubscribeNotifyServer // stream is the server side of the RPC stream
	finished chan<- bool                                  // finished is used to signal closure of a client subscribing goroutine
}

// Subscribe handles a subscribe request from a client
func (g *GrpcAdapter) SubscribeNotify(req *pb.RequestWithParam, stream pb.NotificationService_SubscribeNotifyServer) error {
	// Handle subscribe request
	slog.Info("Received subscribe request from gateway for Notificatons", "id", req.Param)

	fin := make(chan bool)

	// Save the subscriber stream according to the given client Id
	g.notifySubscribers.Store(req.Param, subNotify{stream: stream, finished: fin})

	ctx := stream.Context()

	// Keep this scope alive because once this scope exists - the stream is closed
	for {
		select {
		case <-fin:
			slog.Info("Closing notification stream for gateway client", "id", req.Param)
			return nil

		case <-ctx.Done():
			slog.Warn("Gateway client has disconnected from notification stream", "id", req.Param)
			g.notifySubscribers.Delete(req.Param)

			return nil
		}
	}
}

// Unsubscribe handles a unsubscribe request from a client
func (g *GrpcAdapter) UnsubscribeNotify(ctx context.Context, req *pb.RequestWithParam) (*emptypb.Empty, error) {
	v, ok := g.notifySubscribers.Load(req.Param)
	if !ok {
		return nil, serviceError(fmt.Errorf("failed to load subscriber key: %s", req.Param))
	}

	sub, ok := v.(subNotify)
	if !ok {
		return nil, serviceError(fmt.Errorf("failed to cast subscriber value: %T", v))
	}

	select {
	case sub.finished <- true:
		slog.Info("Unsubscribed client", "id", req.Param, "stream", "notifications")
	default:
		// Default case is to avoid blocking in case client has already unsubscribed
	}

	g.notifySubscribers.Delete(req.Param)

	return &emptypb.Empty{}, nil
}

func (g *GrpcAdapter) PushMessage(cmd string, notifcation *order.Notification) {
	var unsubscribe []string

	// Iterate over all subscribers and send data to each client
	g.notifySubscribers.Range(func(key, value any) bool {
		id, ok := key.(string)
		if !ok {
			slog.Error("failed to cast subscriber key", "key", key)
			return false
		}

		sub, ok := value.(subNotify)
		if !ok {
			slog.Error("failed to cast subscriber value", "value", value)
		}

		// Send data over gRPC stream to client
		err := sub.stream.Send(&pb.NotificationMessage{
			Id:               notifcation.Id,
			UserTo:           notifcation.UserTo,
			SenderUsername:   notifcation.SenderUsername,
			SenderPicture:    notifcation.SenderPicture,
			ReceiverUsername: notifcation.ReceiverUsername,
			ReceiverPicture:  notifcation.ReceiverPicture,
			IsRead:           notifcation.IsRead,
			Message:          notifcation.Message,
			OrderId:          notifcation.OrderId,
			Cmd:              cmd,
			CreatedAt:        utils.ToDateTime(notifcation.CreatedAt),
		})
		if err != nil {
			slog.With("error", err).Error("failed to send data to client", "id", id, "stream", "notifications")

			select {
			case sub.finished <- true:
				slog.Info("Unsubscribe client", "id", id, "stream", "notifications")
			default:
				// Default case is to avoid blocking in case client has already unsubscribed
			}

			unsubscribe = append(unsubscribe, id)
		}

		return true
	})

	// Unsubscribe erroneous client streams
	for _, id := range unsubscribe {
		g.notifySubscribers.Delete(id)
	}
}
