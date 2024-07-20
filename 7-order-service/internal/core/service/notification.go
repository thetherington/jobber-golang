package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/thetherington/jobber-common/error-handling/svc"
	"github.com/thetherington/jobber-common/models/order"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	notificationDB *mongo.Collection
)

/**
 * NotificatonService implements
 */
type NotificatonService struct{}

// NewNotificatonService creates a new notification service instance
func NewNotificatonService(db *mongo.Database) *NotificatonService {
	notificationDB = db.Collection("OrderNotification")

	notificationDB.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "userTo", Value: 1}}},
	})

	return &NotificatonService{}
}

func (ns *NotificatonService) CreateNotification(data *order.Notification) (*order.Notification, error) {
	result, err := notificationDB.InsertOne(context.Background(), data)
	if err != nil {
		return nil, err
	}

	if objId, ok := result.InsertedID.(primitive.ObjectID); ok {
		data.Id = objId.Hex()
	}

	return data, nil
}

func (ns *NotificatonService) GetNotificationsById(ctx context.Context, userToId string) ([]*order.Notification, error) {
	match := bson.D{{Key: "$match", Value: bson.M{"userTo": userToId}}}

	cursor, err := notificationDB.Aggregate(ctx, mongo.Pipeline{match})
	if err != nil {
		slog.With("error", err).Error("failed to find notifications")
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	var notifications []*order.Notification

	if err = cursor.All(context.TODO(), &notifications); err != nil {
		return nil, err
	}

	return notifications, nil
}

func (ns *NotificatonService) MarkNotificationAsRead(ctx context.Context, notificationId string) (*order.Notification, error) {
	var notification *order.Notification

	// convert id string to ObjectId
	objectId, err := primitive.ObjectIDFromHex(notificationId)
	if err != nil {
		return nil, svc.NewError(svc.ErrBadRequest, fmt.Errorf("id provided is invalid"))
	}

	err = notificationDB.FindOneAndUpdate(ctx, bson.M{"_id": objectId},
		bson.D{{Key: "$set", Value: bson.D{{Key: "isRead", Value: true}}}},
		options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&notification)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, svc.NewError(svc.ErrNotFound, fmt.Errorf("order does not exist"))
		}

		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	return notification, nil
}

func (ns *NotificatonService) CreateReturnNotification(o *order.OrderDocument, userToId string, message string) (*order.Notification, error) {
	return ns.CreateNotification(&order.Notification{
		UserTo:           userToId,
		SenderUsername:   o.SellerUsername,
		SenderPicture:    o.SellerImage,
		ReceiverUsername: o.BuyerUsername,
		ReceiverPicture:  o.BuyerImage,
		IsRead:           false,
		Message:          message,
		OrderId:          o.OrderId,
		CreatedAt:        &time.Time{},
	})
}
