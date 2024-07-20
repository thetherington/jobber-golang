package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/thetherington/jobber-common/error-handling/svc"
	"github.com/thetherington/jobber-common/models/users"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var buyerDB *mongo.Collection

/**
 * BuyerService implements
 */
type BuyerService struct {
}

// NewBuyerService creates a new buyer service instance
func NewBuyerService(db *mongo.Database) *BuyerService {
	buyerDB = db.Collection("Buyer")

	buyerDB.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "username", Value: 1}}},
		{Keys: bson.D{{Key: "email", Value: 1}}},
	})

	return &BuyerService{}
}

func (bs *BuyerService) CreateBuyer(b *users.Buyer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := buyerDB.InsertOne(ctx, b)
	if err != nil {
		slog.With("error", err).Error("Error creating buyer")
	}

	return nil
}

func (bs *BuyerService) GetBuyerByEmail(ctx context.Context, email string) (*users.Buyer, error) {
	var buyer *users.Buyer

	err := buyerDB.FindOne(ctx, bson.D{{Key: "email", Value: email}}).Decode(&buyer)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, svc.NewError(svc.ErrNotFound, fmt.Errorf("buyer email does not exist"))
		}

		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	return buyer, nil
}

func (bs *BuyerService) GetBuyerByUsername(ctx context.Context, username string) (*users.Buyer, error) {
	var b *users.Buyer

	err := buyerDB.FindOne(ctx, bson.D{{Key: "username", Value: username}}).Decode(&b)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, svc.NewError(svc.ErrNotFound, fmt.Errorf("buyer username does not exist"))
		}

		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	return b, nil
}

func (bs *BuyerService) GetRandomBuyers(ctx context.Context, count int) ([]*users.Buyer, error) {
	agg := bson.D{{Key: "$sample", Value: bson.D{{Key: "size", Value: count}}}}

	cursor, err := buyerDB.Aggregate(ctx, mongo.Pipeline{agg})
	if err != nil {
		return nil, err
	}

	var buyers []*users.Buyer

	if err = cursor.All(context.TODO(), &buyers); err != nil {
		return nil, err
	}

	return buyers, nil
}

func (bs *BuyerService) UpdateBuyerIsSeller(ctx context.Context, email string) error {
	_, err := buyerDB.UpdateOne(
		ctx,
		bson.D{{Key: "email", Value: email}}, // filter
		bson.D{{Key: "$set", Value: bson.D{{Key: "isSeller", Value: true}}}}, // update
	)

	if err != nil {
		slog.With("error", err).Error("failed to update isSeller field to true", "email", email)

		if errors.Is(err, mongo.ErrNoDocuments) {
			return svc.NewError(svc.ErrNotFound, fmt.Errorf("failed to lookup buyer to make sller"))
		}

		return svc.NewError(svc.ErrInternalFailure, fmt.Errorf("internal failure to make buyer a seller"))
	}

	return nil
}

func (bs *BuyerService) UpdateBuyerPurchasedGigs(buyerId string, pruchasedGigId string, action string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	method := "$pull"

	if action == "PURCHASED_GIG" {
		method = "$push"
	}

	id, _ := primitive.ObjectIDFromHex(buyerId)

	_, err := buyerDB.UpdateOne(ctx,
		bson.D{{Key: "_id", Value: id}},
		bson.D{{Key: method, Value: bson.D{
			{Key: "purchasedGigs", Value: pruchasedGigId},
		}}},
	)
	if err != nil {
		slog.With("error", err).Error("failed to update buyer purchasedGigs", "id", buyerId)
	}

	return nil
}
