package port

import (
	"context"

	"github.com/thetherington/jobber-common/models/chat"
	pb "github.com/thetherington/jobber-common/protogen/go/chat"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ChatRPCClient interface {
	CreateMessage(ctx context.Context, in *pb.MessageDocument, opts ...grpc.CallOption) (*pb.MessageResponse, error)
	GetConversations(ctx context.Context, in *pb.RequestMsgConversations, opts ...grpc.CallOption) (*pb.ConversationsResponse, error)
	GetMessages(ctx context.Context, in *pb.RequestMsgConversations, opts ...grpc.CallOption) (*pb.MessagesResponse, error)
	GetConversationList(ctx context.Context, in *pb.RequestWithParam, opts ...grpc.CallOption) (*pb.MessagesResponse, error)
	GetUserMessages(ctx context.Context, in *pb.RequestWithParam, opts ...grpc.CallOption) (*pb.MessagesResponse, error)
	MarkMultipleMessages(ctx context.Context, in *pb.RequestMsgConversations, opts ...grpc.CallOption) (*pb.ResponseMessage, error)
	MarkSingleMessage(ctx context.Context, in *pb.RequestWithParam, opts ...grpc.CallOption) (*pb.MessageResponse, error)
	UpdateOffer(ctx context.Context, in *pb.UpdateOfferRequest, opts ...grpc.CallOption) (*pb.MessageResponse, error)
	Subscribe(ctx context.Context, in *pb.RequestWithParam, opts ...grpc.CallOption) (pb.ChatService_SubscribeClient, error)
	Unsubscribe(ctx context.Context, in *pb.RequestWithParam, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type ChatService interface {
	AddMessage(ctx context.Context, msg *chat.MessageDocument) (*chat.MessageResponse, error)
	GetConversation(ctx context.Context, sender string, receiver string) (*chat.ConversationsResponse, error)
	GetMessages(ctx context.Context, sender string, receiver string) (*chat.MessagesResponse, error)
	GetUserConversationList(ctx context.Context, username string) (*chat.ConversationListResponse, error)
	GetUserMessages(ctx context.Context, conversationId string) (*chat.MessagesResponse, error)
	MarkManyMessagesAsRead(ctx context.Context, senderUsername string, receiverUsername string, messageId string) (string, error)
	MarkMessageAsRead(ctx context.Context, messageId string) (*chat.MarkSingleMessageResponse, error)
	UpdateOffer(ctx context.Context, messageId string, action string) (*chat.MarkSingleMessageResponse, error)
}

type ChatSocketManager interface {
	DispatchMessage(cmd string, payload *chat.MessageDocument)
}
