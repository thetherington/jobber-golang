package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/thetherington/jobber-common/error-handling/svc"
	"github.com/thetherington/jobber-common/models/order"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (o *OrderService) GetOrderById(ctx context.Context, id string) (*order.OrderResponse, error) {
	var orderDoc *order.OrderDocument

	err := orderDB.FindOne(ctx, bson.M{"orderId": id}).Decode(&orderDoc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, svc.NewError(svc.ErrNotFound, fmt.Errorf("order id does not exist"))
		}

		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	return &order.OrderResponse{Message: "Order by order id", Order: orderDoc}, nil
}

func (o *OrderService) GetSellerOrders(ctx context.Context, id string) (*order.OrdersResponse, error) {
	match := bson.D{{Key: "$match", Value: bson.M{"sellerId": id}}}

	cursor, err := orderDB.Aggregate(ctx, mongo.Pipeline{match})
	if err != nil {
		slog.With("error", err).Error("failed to find orders")
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	var orderDocs []*order.OrderDocument

	if err = cursor.All(context.TODO(), &orderDocs); err != nil {
		return nil, err
	}

	return &order.OrdersResponse{Message: "Seller orders", Orders: orderDocs}, nil
}

func (o *OrderService) GetBuyerOrders(ctx context.Context, id string) (*order.OrdersResponse, error) {
	match := bson.D{{Key: "$match", Value: bson.M{"buyerId": id}}}

	cursor, err := orderDB.Aggregate(ctx, mongo.Pipeline{match})
	if err != nil {
		slog.With("error", err).Error("failed to find orders")
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	var orderDocs []*order.OrderDocument

	if err = cursor.All(context.TODO(), &orderDocs); err != nil {
		return nil, err
	}

	return &order.OrdersResponse{Message: "Buyer orders", Orders: orderDocs}, nil
}
