package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/thetherington/jobber-common/error-handling/svc"
	"github.com/thetherington/jobber-common/models/chat"
	"github.com/thetherington/jobber-common/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (c *ChatService) MarkMessageAsRead(ctx context.Context, messageId string) (*chat.MarkSingleMessageResponse, error) {
	// convert messageId string to ObjectId
	objectId, err := primitive.ObjectIDFromHex(messageId)
	if err != nil {
		return nil, svc.NewError(svc.ErrBadRequest, fmt.Errorf("id provided is invalid"))
	}

	var message *chat.MessageDocument

	// update message by id and get updated message from mongo
	err = messageDB.FindOneAndUpdate(ctx, bson.M{"_id": objectId}, bson.D{{
		Key:   "$set",
		Value: bson.M{"isRead": true},
	}}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&message)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, svc.NewError(svc.ErrNotFound, fmt.Errorf("message does not exist"))
		}

		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	return &chat.MarkSingleMessageResponse{Message: "Message marked as read", SingleMessage: message}, nil
}

func (c *ChatService) MarkManyMessagesAsRead(ctx context.Context, senderUsername, receiverUsername, messageId string) (*chat.MarkSingleMessageResponse, error) {
	filter := bson.D{
		{Key: "$and", Value: bson.A{
			bson.M{"senderUsername": senderUsername},
			bson.M{"receiverUsername": receiverUsername},
			bson.M{"isRead": false},
		}},
	}

	// update all messages with filter to be marked as read
	if _, err := messageDB.UpdateMany(ctx, filter, bson.D{{
		Key:   "$set",
		Value: bson.M{"isRead": true},
	}}); err != nil {
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	// convert messageId string to ObjectId
	objectId, err := primitive.ObjectIDFromHex(messageId)
	if err != nil {
		return nil, svc.NewError(svc.ErrBadRequest, fmt.Errorf("id provided is invalid"))
	}

	var message *chat.MessageDocument

	if err := messageDB.FindOne(ctx, bson.M{"_id": objectId}).Decode(&message); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, svc.NewError(svc.ErrNotFound, fmt.Errorf("message does not exist"))
		}

		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	return &chat.MarkSingleMessageResponse{Message: "Messages marked as read", SingleMessage: message}, nil
}

func (c *ChatService) UpdateOffer(ctx context.Context, messageId, action string) (*chat.MessageResponse, error) {
	// convert messageId string to ObjectId
	objectId, err := primitive.ObjectIDFromHex(messageId)
	if err != nil {
		return nil, svc.NewError(svc.ErrBadRequest, fmt.Errorf("id provided is invalid"))
	}

	var message *chat.MessageDocument

	// update message by id offer field of accepted or cancelled
	err = messageDB.FindOneAndUpdate(ctx, bson.M{"_id": objectId}, bson.D{{
		Key:   "$set",
		Value: bson.D{{Key: fmt.Sprintf("offer.%s", utils.LowerCase(action)), Value: true}}},
	}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&message)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, svc.NewError(svc.ErrNotFound, fmt.Errorf("message does not exist"))
		}

		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	return &chat.MessageResponse{Message: "Message offer updated", MessageData: message}, nil
}
