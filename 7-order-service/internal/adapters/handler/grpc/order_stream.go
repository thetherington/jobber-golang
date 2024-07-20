package grpc

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/thetherington/jobber-common/models/order"
	pb "github.com/thetherington/jobber-common/protogen/go/order"
	"google.golang.org/protobuf/types/known/emptypb"
)

type subOrder struct {
	stream   pb.OrderService_SubscribeOrderServer // stream is the server side of the RPC stream
	finished chan<- bool                          // finished is used to signal closure of a client subscribing goroutine
}

// Subscribe handles a subscribe request from a client
func (g *GrpcAdapter) SubscribeOrder(req *pb.RequestById, stream pb.OrderService_SubscribeOrderServer) error {
	// Handle subscribe request
	slog.Info("Received subscribe request from gateway for Orders", "id", req.Id)

	fin := make(chan bool)

	// Save the subscriber stream according to the given client Id
	g.orderSubscribers.Store(req.Id, subOrder{stream: stream, finished: fin})

	ctx := stream.Context()

	// Keep this scope alive because once this scope exists - the stream is closed
	for {
		select {
		case <-fin:
			slog.Info("Closing order stream for gateway client", "id", req.Id)
			return nil

		case <-ctx.Done():
			slog.Warn("Gateway client has disconnected from order stream", "id", req.Id)
			g.orderSubscribers.Delete(req.Id)

			return nil
		}
	}
}

// Unsubscribe handles a unsubscribe request from a client
func (g *GrpcAdapter) UnsubscribeOrder(ctx context.Context, req *pb.RequestById) (*emptypb.Empty, error) {
	v, ok := g.orderSubscribers.Load(req.Id)
	if !ok {
		return nil, serviceError(fmt.Errorf("failed to load subscriber key: %s", req.Id))
	}

	sub, ok := v.(subOrder)
	if !ok {
		return nil, serviceError(fmt.Errorf("failed to cast subscriber value: %T", v))
	}

	select {
	case sub.finished <- true:
		slog.Info("Unsubscribed client", "id", req.Id, "stream", "orders")
	default:
		// Default case is to avoid blocking in case client has already unsubscribed
	}

	g.orderSubscribers.Delete(req.GetId())

	return &emptypb.Empty{}, nil
}

func (g *GrpcAdapter) PushOrder(cmd string, order *order.OrderDocument) {
	var unsubscribe []string

	// Iterate over all subscribers and send data to each client
	g.orderSubscribers.Range(func(key, value any) bool {
		id, ok := key.(string)
		if !ok {
			slog.Error("failed to cast subscriber key", "key", key, "stream", "orders")
			return false
		}

		sub, ok := value.(subOrder)
		if !ok {
			slog.Error("failed to cast subscriber value", "value", value, "stream", "orders")
		}

		// Send data over gRPC stream to client
		err := sub.stream.Send(order.MarshalToProto())
		if err != nil {
			slog.With("error", err).Error("failed to send data to client", "id", id, "stream", "orders")

			select {
			case sub.finished <- true:
				slog.Info("Unsubscribe client", "id", id, "stream", "orders")
			default:
				// Default case is to avoid blocking in case client has already unsubscribed
			}

			unsubscribe = append(unsubscribe, id)
		}

		return true
	})

	// Unsubscribe erroneous client streams
	for _, id := range unsubscribe {
		g.orderSubscribers.Delete(id)
	}
}
