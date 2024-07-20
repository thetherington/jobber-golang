package grpc

import (
	"context"
	"log/slog"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/thetherington/jobber-common/models/chat"
	pb "github.com/thetherington/jobber-common/protogen/go/chat"
	"github.com/thetherington/jobber-gateway/internal/core/port"
	"go.elastic.co/apm/module/apmgrpc/v2"
	"google.golang.org/grpc"
)

type ChatServiceClient struct {
	id string
	ws port.ChatSocketManager
	pb.ChatServiceClient
}

func NewChatAdapter(address string, session *scs.SessionManager, ws port.ChatSocketManager) *ChatServiceClient {
	ctxForwardInterceptor := CreateContextForwarderInterceptor(session)   // create context interceptor forwarder
	tokenForwardInterceptor := CreateTokenForwarderInterceptor("message") // create gateway token and gRPC header

	streamTokenForwardInterceptor := CreateStreamTokenForwarderInterceptor("message") // create gateway token and gRPC header

	conn, err := NewClient(
		address,
		grpc.WithChainUnaryInterceptor(ctxForwardInterceptor, tokenForwardInterceptor, apmgrpc.NewUnaryClientInterceptor()),
		grpc.WithChainStreamInterceptor(streamTokenForwardInterceptor, apmgrpc.NewStreamClientInterceptor()),
	)
	if err != nil {
		slog.With("error", err).Error("Failed to create Chat gRPC Client")
	}

	chat := pb.NewChatServiceClient(conn) // chat crud rpc calls
	id := uuid.NewString()

	return &ChatServiceClient{
		id,
		ws,
		chat,
	}
}

func (ch *ChatServiceClient) SubscribeStream() {
	var err error

	// stream is the client side of the RPC stream
	var stream pb.ChatService_SubscribeClient

	for {
		// if stream is not ready try to subscribe
		if stream == nil {
			ch.id = uuid.NewString()

			slog.Info("Subscribing to Chat Service Message Stream", "id", ch.id)

			if stream, err = ch.Subscribe(context.Background(), &pb.RequestWithParam{Param: ch.id}); err != nil {
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
		ch.ws.DispatchMessage(response.Message, chat.CreateMessageDocument(response.MessageData))
	}
}

func (ch *ChatServiceClient) DisconnectStream() {
	slog.Info("Unsubscribing to Chat Service Message Stream")

	if _, err := ch.Unsubscribe(context.Background(), &pb.RequestWithParam{Param: ch.id}); err != nil {
		slog.With("error", err).Error("unsubscribe failed")
	}
}
