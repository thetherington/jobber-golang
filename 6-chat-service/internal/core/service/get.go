package service

import (
	"context"
	"log/slog"

	"github.com/thetherington/jobber-common/error-handling/svc"
	"github.com/thetherington/jobber-common/models/chat"
	"github.com/thetherington/jobber-common/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (c *ChatService) GetMessages(ctx context.Context, sender, receiver string) (*chat.MessagesResponse, error) {
	match := bson.D{
		{Key: "$match", Value: bson.D{{
			Key: "$or",
			Value: bson.A{
				bson.D{{Key: "$and", Value: bson.A{
					bson.M{"senderUsername": utils.FirstLetterUpperCase(sender)},
					bson.M{"receiverUsername": utils.FirstLetterUpperCase(receiver)},
				}}},
				bson.D{{Key: "$and", Value: bson.A{
					bson.M{"senderUsername": utils.FirstLetterUpperCase(receiver)},
					bson.M{"receiverUsername": utils.FirstLetterUpperCase(sender)},
				}}},
			},
		}}},
	}

	sort := bson.D{
		{Key: "$sort", Value: bson.M{"createdAt": 1}},
	}

	cursor, err := messageDB.Aggregate(ctx, mongo.Pipeline{match, sort})
	if err != nil {
		slog.With("error", err).Error("failed to find messages")
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	var messages []*chat.MessageDocument

	if err = cursor.All(context.TODO(), &messages); err != nil {
		return nil, err
	}

	return &chat.MessagesResponse{Message: "Chat messages", Messages: messages}, nil
}

func (c *ChatService) GetUserMessages(ctx context.Context, conversationId string) (*chat.MessagesResponse, error) {
	match := bson.D{{Key: "$match", Value: bson.M{"conversationId": conversationId}}}

	sort := bson.D{{Key: "$sort", Value: bson.M{"createdAt": 1}}}

	cursor, err := messageDB.Aggregate(ctx, mongo.Pipeline{match, sort})
	if err != nil {
		slog.With("error", err).Error("failed to find messages")
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	var messages []*chat.MessageDocument

	if err = cursor.All(context.TODO(), &messages); err != nil {
		return nil, err
	}

	return &chat.MessagesResponse{Message: "Chat messages", Messages: messages}, nil
}

func (c *ChatService) GetUserConversationList(ctx context.Context, username string) (*chat.ConversationListResponse, error) {
	match := bson.D{
		{Key: "$match", Value: bson.D{{
			Key: "$or",
			Value: bson.A{
				bson.M{"senderUsername": utils.FirstLetterUpperCase(username)},
				bson.M{"receiverUsername": utils.FirstLetterUpperCase(username)},
			},
		}}},
	}

	group := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$conversationId"},
			{Key: "result", Value: bson.D{
				{Key: "$top", Value: bson.D{
					{Key: "output", Value: "$$ROOT"},
					{Key: "sortBy", Value: bson.M{"createdAt": -1}},
				}},
			}},
		}},
	}

	project := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "_id", Value: "$result._id"},
			{Key: "conversationId", Value: "$result.conversationId"},
			{Key: "sellerId", Value: "$result.sellerId"},
			{Key: "buyerId", Value: "$result.buyerId"},
			{Key: "receiverUsername", Value: "$result.receiverUsername"},
			{Key: "receiverPicture", Value: "$result.receiverPicture"},
			{Key: "senderUsername", Value: "$result.senderUsername"},
			{Key: "senderPicture", Value: "$result.senderPicture"},
			{Key: "body", Value: "$result.body"},
			{Key: "file", Value: "$result.file"},
			{Key: "gigId", Value: "$result.gigId"},
			{Key: "isRead", Value: "$result.isRead"},
			{Key: "hasOffer", Value: "$result.hasOffer"},
			{Key: "createdAt", Value: "$result.createdAt"},
		}},
	}

	cursor, err := messageDB.Aggregate(ctx, mongo.Pipeline{match, group, project})
	if err != nil {
		slog.With("error", err).Error("failed to find messages")
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	var messages []*chat.MessageDocument

	if err = cursor.All(context.TODO(), &messages); err != nil {
		return nil, err
	}

	return &chat.ConversationListResponse{Message: "Chat messages", Conversations: messages}, nil
}

func (c *ChatService) GetConversation(ctx context.Context, sender, receiver string) (*chat.ConversationsResponse, error) {
	match := bson.D{
		{Key: "$match", Value: bson.D{{
			Key: "$or",
			Value: bson.A{
				bson.D{{Key: "$and", Value: bson.A{
					bson.M{"senderUsername": utils.FirstLetterUpperCase(sender)},
					bson.M{"receiverUsername": utils.FirstLetterUpperCase(receiver)},
				}}},
				bson.D{{Key: "$and", Value: bson.A{
					bson.M{"senderUsername": utils.FirstLetterUpperCase(receiver)},
					bson.M{"receiverUsername": utils.FirstLetterUpperCase(sender)},
				}}},
			},
		}}},
	}

	cursor, err := conversationDB.Aggregate(ctx, mongo.Pipeline{match})
	if err != nil {
		slog.With("error", err).Error("failed to find conversation")
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	var conversations []*chat.ConversationDocument

	if err = cursor.All(context.TODO(), &conversations); err != nil {
		return nil, err
	}

	return &chat.ConversationsResponse{Message: "Chat conversation", Conversations: conversations}, nil
}
