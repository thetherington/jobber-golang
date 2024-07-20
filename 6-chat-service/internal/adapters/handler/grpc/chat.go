package grpc

import (
	"context"

	"github.com/thetherington/jobber-common/error-handling/grpcerror"
	"github.com/thetherington/jobber-common/models/chat"
	"github.com/thetherington/jobber-common/models/event"
	pb "github.com/thetherington/jobber-common/protogen/go/chat"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// parses service error and returns a equivelient gRPC error
func serviceError(err error) error {
	// try to cast the error to a grpcerror lookup
	if apiError, ok := grpcerror.FromError(err); ok {
		s := status.New(apiError.Status, apiError.Message)
		return s.Err()
	}

	// generic response
	s := status.New(codes.Internal, err.Error())
	return s.Err()
}

func (g *GrpcAdapter) CreateMessage(ctx context.Context, req *pb.MessageDocument) (*pb.MessageResponse, error) {
	resp, err := g.chatService.AddMessage(ctx, chat.CreateMessageDocument(req))
	if err != nil {
		return nil, serviceError(err)
	}

	// push message onto streaming response for front end web socket
	go g.PushMessage(event.EventMessageReceived, resp.MessageData)

	return &pb.MessageResponse{
		Message:        resp.Message,
		ConversationId: &resp.ConversationId,
		MessageData:    chat.CreateProtoMessageDocument(resp.MessageData),
	}, nil
}

func (g *GrpcAdapter) GetConversations(ctx context.Context, req *pb.RequestMsgConversations) (*pb.ConversationsResponse, error) {
	resp, err := g.chatService.GetConversation(ctx, req.SenderUsername, req.ReceiverUsername)
	if err != nil {
		return nil, serviceError(err)
	}

	conversations := make([]*pb.Conversation, 0)

	for _, c := range resp.Conversations {
		conversations = append(conversations, &pb.Conversation{
			ConversationId:   c.CoversationId,
			SenderUsername:   c.SenderUsername,
			ReceiverUsername: c.ReceiverUsername,
		})
	}

	return &pb.ConversationsResponse{Message: resp.Message, Conversations: conversations}, nil
}

func (g *GrpcAdapter) GetMessages(ctx context.Context, req *pb.RequestMsgConversations) (*pb.MessagesResponse, error) {
	resp, err := g.chatService.GetMessages(ctx, req.SenderUsername, req.ReceiverUsername)
	if err != nil {
		return nil, serviceError(err)
	}

	messages := make([]*pb.MessageDocument, 0)

	for _, m := range resp.Messages {
		messages = append(messages, chat.CreateProtoMessageDocument(m))
	}

	return &pb.MessagesResponse{Message: resp.Message, Messages: messages}, nil
}

func (g *GrpcAdapter) GetConversationList(ctx context.Context, req *pb.RequestWithParam) (*pb.MessagesResponse, error) {
	resp, err := g.chatService.GetUserConversationList(ctx, req.Param)
	if err != nil {
		return nil, serviceError(err)
	}

	messages := make([]*pb.MessageDocument, 0)

	for _, m := range resp.Conversations {
		messages = append(messages, chat.CreateProtoMessageDocument(m))
	}

	return &pb.MessagesResponse{Message: resp.Message, Messages: messages}, nil
}

func (g *GrpcAdapter) GetUserMessages(ctx context.Context, req *pb.RequestWithParam) (*pb.MessagesResponse, error) {
	resp, err := g.chatService.GetUserMessages(ctx, req.Param)
	if err != nil {
		return nil, serviceError(err)
	}

	messages := make([]*pb.MessageDocument, 0)

	for _, m := range resp.Messages {
		messages = append(messages, chat.CreateProtoMessageDocument(m))
	}

	return &pb.MessagesResponse{Message: resp.Message, Messages: messages}, nil
}

func (g *GrpcAdapter) MarkMultipleMessages(ctx context.Context, req *pb.RequestMsgConversations) (*pb.ResponseMessage, error) {
	resp, err := g.chatService.MarkManyMessagesAsRead(
		ctx,
		req.SenderUsername,
		req.ReceiverUsername,
		*req.MessageId,
	)
	if err != nil {
		return nil, serviceError(err)
	}

	// send message back to frontend  websocket via grpc streaming response.
	go g.PushMessage(event.EventMessageUpdated, resp.SingleMessage)

	return &pb.ResponseMessage{Message: resp.Message}, nil
}

func (g *GrpcAdapter) MarkSingleMessage(ctx context.Context, req *pb.RequestWithParam) (*pb.MessageResponse, error) {
	resp, err := g.chatService.MarkMessageAsRead(ctx, req.Param)
	if err != nil {
		return nil, serviceError(err)
	}

	// send message back to frontend  websocket via grpc streaming response.
	go g.PushMessage(event.EventMessageUpdated, resp.SingleMessage)

	return &pb.MessageResponse{Message: resp.Message, MessageData: chat.CreateProtoMessageDocument(resp.SingleMessage)}, nil
}

func (g *GrpcAdapter) UpdateOffer(ctx context.Context, req *pb.UpdateOfferRequest) (*pb.MessageResponse, error) {
	resp, err := g.chatService.UpdateOffer(ctx, req.MessageId, req.Action.String())
	if err != nil {
		return nil, serviceError(err)
	}

	return &pb.MessageResponse{Message: resp.Message, MessageData: chat.CreateProtoMessageDocument(resp.MessageData)}, nil
}
