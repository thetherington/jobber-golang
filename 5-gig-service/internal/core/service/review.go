package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/thetherington/jobber-common/error-handling/svc"
	"github.com/thetherington/jobber-common/models/gig"
	"github.com/thetherington/jobber-common/models/review"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (g *GigService) UpdateGigReview(data *review.ReviewMessageDetails) error {
	objectId, err := primitive.ObjectIDFromHex(data.GigId)
	if err != nil {
		return svc.NewError(svc.ErrBadRequest, fmt.Errorf("id provided is invalid"))
	}

	ratingTypes := map[string]string{
		"1": "one",
		"2": "two",
		"3": "three",
		"4": "four",
		"5": "five",
	}
	ratingKey := ratingTypes[strconv.Itoa(int(data.Rating))]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var updatedGig gig.SellerGig

	// update mongodb
	err = gigDB.FindOneAndUpdate(ctx, bson.M{"_id": objectId}, bson.D{
		{Key: "$inc",
			Value: bson.D{
				{Key: "ratingsCount", Value: 1},
				{Key: "ratingSum", Value: data.Rating},
				{Key: fmt.Sprintf("ratingCategories.%s.value", ratingKey), Value: data.Rating},
				{Key: fmt.Sprintf("ratingCategories.%s.count", ratingKey), Value: 1},
			},
		},
	}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedGig)
	if err != nil {
		return err
	}

	// update elasticsearch
	if _, err = g.search.UpdateGig(ctx, data.GigId, &updatedGig); err != nil {
		return err
	}

	return err
}
