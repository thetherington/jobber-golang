package service

import (
	"context"
	"log/slog"

	"github.com/thetherington/jobber-common/error-handling/svc"
	"github.com/thetherington/jobber-common/models/chat"
	pb "github.com/thetherington/jobber-common/protogen/go/chat"
	"github.com/thetherington/jobber-gateway/internal/core/port"
)

/**
 * ChatService implements
 */
type ChatService struct {
	client port.ChatRPCClient
}

// NewChatService creates a new chat service instance
func NewChatService(rpc port.ChatRPCClient) *ChatService {
	return &ChatService{
		rpc,
	}
}

func (c *ChatService) AddMessage(ctx context.Context, msg *chat.MessageDocument) (*chat.MessageResponse, error) {
	req := chat.CreateProtoMessageDocument(msg)

	resp, err := c.client.CreateMessage(ctx, req)
	if err != nil {
		slog.With("error", err).Debug("AddMessage error")
		return nil, svc.GrpcErrorResolve(err, "AddMessage")
	}

	return &chat.MessageResponse{
		Message:        resp.Message,
		ConversationId: *resp.ConversationId,
		MessageData:    chat.CreateMessageDocument(resp.MessageData),
	}, nil
}

func (c *ChatService) GetMessages(ctx context.Context, sender, receiver string) (*chat.MessagesResponse, error) {
	resp, err := c.client.GetMessages(ctx, &pb.RequestMsgConversations{
		SenderUsername:   sender,
		ReceiverUsername: receiver,
	})
	if err != nil {
		slog.With("error", err).Debug("GetMessages error")
		return nil, svc.GrpcErrorResolve(err, "GetMessages")
	}

	messages := make([]*chat.MessageDocument, 0)

	for _, m := range resp.Messages {
		messages = append(messages, chat.CreateMessageDocument(m))
	}

	return &chat.MessagesResponse{Message: resp.Message, Messages: messages}, nil
}

func (c *ChatService) GetUserMessages(ctx context.Context, conversationId string) (*chat.MessagesResponse, error) {
	resp, err := c.client.GetUserMessages(ctx, &pb.RequestWithParam{
		Param: conversationId,
	})
	if err != nil {
		slog.With("error", err).Debug("GetUserMessages error")
		return nil, svc.GrpcErrorResolve(err, "GetUserMessages")
	}

	messages := make([]*chat.MessageDocument, 0)

	for _, m := range resp.Messages {
		messages = append(messages, chat.CreateMessageDocument(m))
	}

	return &chat.MessagesResponse{Message: resp.Message, Messages: messages}, nil
}

func (c *ChatService) GetUserConversationList(ctx context.Context, username string) (*chat.ConversationListResponse, error) {
	resp, err := c.client.GetConversationList(ctx, &pb.RequestWithParam{
		Param: username,
	})
	if err != nil {
		slog.With("error", err).Debug("GetUserConversationList error")
		return nil, svc.GrpcErrorResolve(err, "GetUserConversationList")
	}

	messages := make([]*chat.MessageDocument, 0)

	for _, m := range resp.Messages {
		messages = append(messages, chat.CreateMessageDocument(m))
	}

	return &chat.ConversationListResponse{Message: resp.Message, Conversations: messages}, nil
}

func (c *ChatService) GetConversation(ctx context.Context, sender, receiver string) (*chat.ConversationsResponse, error) {
	resp, err := c.client.GetConversations(ctx, &pb.RequestMsgConversations{
		SenderUsername:   sender,
		ReceiverUsername: receiver,
	})
	if err != nil {
		slog.With("error", err).Debug("GetConversation error")
		return nil, svc.GrpcErrorResolve(err, "GetConversation")
	}

	conversations := make([]*chat.ConversationDocument, 0)

	for _, c := range resp.Conversations {
		conversations = append(conversations, &chat.ConversationDocument{
			CoversationId:    c.ConversationId,
			SenderUsername:   c.SenderUsername,
			ReceiverUsername: c.ReceiverUsername,
		})
	}

	return &chat.ConversationsResponse{Message: resp.Message, Conversations: conversations}, nil
}

func (c *ChatService) MarkMessageAsRead(ctx context.Context, messageId string) (*chat.MarkSingleMessageResponse, error) {
	resp, err := c.client.MarkSingleMessage(ctx, &pb.RequestWithParam{Param: messageId})
	if err != nil {
		slog.With("error", err).Debug("MarkMessageAsRead error")
		return nil, svc.GrpcErrorResolve(err, "MarkMessageAsRead")
	}

	return &chat.MarkSingleMessageResponse{
		Message:       resp.Message,
		SingleMessage: chat.CreateMessageDocument(resp.MessageData),
	}, nil
}

func (c *ChatService) MarkManyMessagesAsRead(ctx context.Context, senderUsername, receiverUsername, messageId string) (string, error) {
	resp, err := c.client.MarkMultipleMessages(ctx, &pb.RequestMsgConversations{
		MessageId:        &messageId,
		SenderUsername:   senderUsername,
		ReceiverUsername: receiverUsername,
	})
	if err != nil {
		slog.With("error", err).Debug("MarkManyMessagesAsRead error")
		return "", svc.GrpcErrorResolve(err, "MarkManyMessagesAsRead")
	}

	return resp.Message, nil
}

func (c *ChatService) UpdateOffer(ctx context.Context, messageId, action string) (*chat.MarkSingleMessageResponse, error) {
	a := pb.OfferAction_Accepted

	if action == "cancelled" {
		a = pb.OfferAction_Cancelled
	}

	resp, err := c.client.UpdateOffer(ctx, &pb.UpdateOfferRequest{MessageId: messageId, Action: a})
	if err != nil {
		slog.With("error", err).Debug("UpdateOffer error")
		return nil, svc.GrpcErrorResolve(err, "UpdateOffer")
	}

	return &chat.MarkSingleMessageResponse{
		Message:       resp.Message,
		SingleMessage: chat.CreateMessageDocument(resp.MessageData),
	}, nil
}
