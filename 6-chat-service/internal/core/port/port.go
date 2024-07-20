package port

import (
	"context"

	"github.com/thetherington/jobber-common/models/chat"
)

type ImageUploader interface {
	UploadImage(ctx context.Context, file string, public_id string, overwrite bool, invalidate bool) (string, error)
	UploadVideo(ctx context.Context, file string, public_id string, overwrite bool, invalidate bool) (string, error)
}

type ChatProducer interface {
	PublishDirectMessage(exchangeName string, routingKey string, data []byte) error
}

type ChatService interface {
	CreateConversation(ctx context.Context, conversationId string, sender string, receiver string) error
	AddMessage(ctx context.Context, msg *chat.MessageDocument) (*chat.MessageResponse, error)
	GetConversation(ctx context.Context, sender string, receiver string) (*chat.ConversationsResponse, error)
	GetMessages(ctx context.Context, sender string, receiver string) (*chat.MessagesResponse, error)
	GetUserConversationList(ctx context.Context, username string) (*chat.ConversationListResponse, error)
	GetUserMessages(ctx context.Context, conversationId string) (*chat.MessagesResponse, error)
	MarkManyMessagesAsRead(ctx context.Context, senderUsername string, receiverUsername string, messageId string) (*chat.MarkSingleMessageResponse, error)
	MarkMessageAsRead(ctx context.Context, messageId string) (*chat.MarkSingleMessageResponse, error)
	UpdateOffer(ctx context.Context, messageId string, action string) (*chat.MessageResponse, error)
}
