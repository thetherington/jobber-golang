package grpc

import (
	"context"
	"log/slog"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/thetherington/jobber-common/models/order"
	pb "github.com/thetherington/jobber-common/protogen/go/order"
	"github.com/thetherington/jobber-common/utils"
	"github.com/thetherington/jobber-gateway/internal/core/port"
	"go.elastic.co/apm/module/apmgrpc/v2"
	"google.golang.org/grpc"
)

type OrderServiceClient struct {
	idNotify string
	idOrder  string
	ws       port.NotificationSocketManager

	pb.OrderServiceClient
	pb.NotificationServiceClient
}

func NewOrderAdapter(address string, session *scs.SessionManager, ws port.NotificationSocketManager) *OrderServiceClient {
	ctxForwardInterceptor := CreateContextForwarderInterceptor(session) // create context interceptor forwarder
	tokenForwardInterceptor := CreateTokenForwarderInterceptor("order") // create gateway token and gRPC header

	streamTokenForwardInterceptor := CreateStreamTokenForwarderInterceptor("order") // create gateway token and gRPC header

	conn, err := NewClient(
		address,
		grpc.WithChainUnaryInterceptor(ctxForwardInterceptor, tokenForwardInterceptor, apmgrpc.NewUnaryClientInterceptor()),
		grpc.WithChainStreamInterceptor(streamTokenForwardInterceptor, apmgrpc.NewStreamClientInterceptor()),
	)
	if err != nil {
		slog.With("error", err).Error("Failed to create Order gRPC Client")
	}

	order := pb.NewOrderServiceClient(conn) // order crud rpc calls
	notification := pb.NewNotificationServiceClient(conn)

	id := uuid.NewString()
	idOrder := uuid.NewString()

	return &OrderServiceClient{
		id,
		idOrder,
		ws,
		order,
		notification,
	}
}

func (ch *OrderServiceClient) SubscribeNotifyStream() {
	var err error

	// stream is the client side of the RPC stream
	var stream pb.NotificationService_SubscribeNotifyClient

	for {
		// if stream is not ready try to subscribe
		if stream == nil {
			ch.idNotify = uuid.NewString()

			slog.Info("Subscribing to Notification Service Message Stream", "id", ch.idNotify)

			if stream, err = ch.SubscribeNotify(context.Background(), &pb.RequestWithParam{Param: ch.idNotify}); err != nil {
				slog.With("error", err).Error("failed to subscribe")

				// Retry on failure
				time.Sleep(time.Second * 5)
				continue
			}
		}

		response, err := stream.Recv()
		if err != nil {
			slog.With("error", err).Error("failed to receive message")
			// Clearing the stream will force the client to resubscribe on next iteration
			stream = nil

			// Retry on failure
			time.Sleep(time.Second * 5)
			continue
		}

		// send message to websocket connections
		ch.ws.DispatchNotification(response.Cmd, &order.Notification{
			Id:               response.Id,
			UserTo:           response.UserTo,
			SenderUsername:   response.SenderUsername,
			SenderPicture:    response.SenderPicture,
			ReceiverUsername: response.ReceiverUsername,
			ReceiverPicture:  response.ReceiverPicture,
			IsRead:           response.IsRead,
			Message:          response.Message,
			OrderId:          response.OrderId,
			CreatedAt:        utils.ToTime(response.CreatedAt),
		})
	}
}

func (ch *OrderServiceClient) DisconnectStream() {
	slog.Info("Unsubscribing to Notification Service Message Stream")

	if _, err := ch.UnsubscribeNotify(context.Background(), &pb.RequestWithParam{Param: ch.idNotify}); err != nil {
		slog.With("error", err).Error("unsubscribe failed")
	}
}

func (ch *OrderServiceClient) SubscribeOrderStream() {
	var err error

	// stream is the client side of the RPC stream
	var stream pb.OrderService_SubscribeOrderClient

	for {
		// if stream is not ready try to subscribe
		if stream == nil {
			ch.idOrder = uuid.NewString()

			slog.Info("Subscribing to Order Service Message Stream", "id", ch.idOrder)

			if stream, err = ch.SubscribeOrder(context.Background(), &pb.RequestById{Id: ch.idOrder}); err != nil {
				slog.With("error", err).Error("failed to subscribe to order stream")

				// Retry on failure
				time.Sleep(time.Second * 5)
				continue
			}
		}

		response, err := stream.Recv()
		if err != nil {
			slog.With("error", err).Error("failed to receive message")
			// Clearing the stream will force the client to resubscribe on next iteration
			stream = nil

			// Retry on failure
			time.Sleep(time.Second * 5)
			continue
		}

		// send message to websocket connections
		ch.ws.DispatchOrder("order_update", order.CreateOrderDocument(response))
	}
}

func (ch *OrderServiceClient) DisconnectOrderStream() {
	slog.Info("Unsubscribing to Order Service Order Stream")

	if _, err := ch.UnsubscribeOrder(context.Background(), &pb.RequestById{Id: ch.idOrder}); err != nil {
		slog.With("error", err).Error("unsubscribe failed")
	}
}
