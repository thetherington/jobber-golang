package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/thetherington/jobber-common/error-handling/svc"
	"github.com/thetherington/jobber-common/middleware"
	"github.com/thetherington/jobber-common/models/gig"
	"github.com/thetherington/jobber-common/models/review"
	pborder "github.com/thetherington/jobber-common/protogen/go/order"
	"github.com/thetherington/jobber-gig/internal/core/port"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/proto"
)

var (
	validate *validator.Validate
	gigDB    *mongo.Collection
)

/**
 * GigService implements
 */
type GigService struct {
	search port.SearchClient
	queue  port.GigProducer
	cache  port.CacheRepository
	image  port.ImageUploader
}

// NewGigService creates a new gig service instance
func NewGigService(db *mongo.Database, search port.SearchClient, queue port.GigProducer, cache port.CacheRepository, image port.ImageUploader) *GigService {
	validate = validator.New(validator.WithRequiredStructEnabled())

	gigDB = db.Collection("Gig")

	gigDB.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "sellerId", Value: 1}}},
	})

	return &GigService{
		search: search,
		queue:  queue,
		cache:  cache,
		image:  image,
	}
}

func (g *GigService) SetQueue(queue port.GigProducer) {
	g.queue = queue
}

func (g *GigService) CreateGig(ctx context.Context, newGig *gig.SellerGig) (*gig.ResponseGig, error) {
	// Validate Gig payload
	if err := newGig.Validate(validate); err != nil {
		slog.With("error", err).Debug("Gig Create Validation Failed")
		return nil, svc.NewError(svc.ErrInvalidData, err)
	}

	// Get username from the user cookie session passed down into the context.
	username := ctx.Value(middleware.CtxUsernameKey)
	if username == nil {
		slog.Debug("Username in context is nil")
		return nil, svc.NewError(svc.ErrUnAuthorized, fmt.Errorf("session invalid. please login and try again"))
	}
	// Get email from the user cookie session passed down into the context.
	email := ctx.Value(middleware.CtxEmailKey)
	if email == nil {
		slog.Debug("Username in context is nil")
		return nil, svc.NewError(svc.ErrUnAuthorized, fmt.Errorf("session invalid. please login and try again"))
	}

	// get gigs count from elastic
	count, err := g.search.GetGigsCount(ctx)
	if err != nil {
		slog.With("error", err).Error("failed to get count")
	}

	// upload gig image to cloudinary if the image isn't not a url
	if err := validate.Var(newGig.CoverImage, "required,http_url"); err != nil {
		imagePublicId := uuid.New().String()

		link, err := g.image.UploadImage(ctx, newGig.CoverImage, imagePublicId, true, true)
		if err != nil {
			return nil, svc.NewError(svc.ErrInternalFailure, err)
		}

		newGig.CoverImage = link
	}

	// update gig properties
	newGig.Username = username.(string)
	newGig.Email = email.(string)
	newGig.SortId = count + 1
	newGig.Active = true

	newGig.RatingCategories = &review.RatingCategories{
		One:   review.RatingCategoryItem{Value: 0, Count: 0},
		Two:   review.RatingCategoryItem{Value: 0, Count: 0},
		Three: review.RatingCategoryItem{Value: 0, Count: 0},
		Four:  review.RatingCategoryItem{Value: 0, Count: 0},
		Five:  review.RatingCategoryItem{Value: 0, Count: 0},
	}

	// save gig into mongodb, get _id and use it in id
	result, err := gigDB.InsertOne(ctx, newGig)
	if err != nil {
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	if objId, ok := result.InsertedID.(primitive.ObjectID); ok {
		newGig.ID = objId.Hex()
	}

	// save gig into elasticsearch
	if _, err = g.search.InsertGig(ctx, newGig); err != nil {
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	// increment the seller gig count via the Users microservice from rabbitmq
	pbmsg := &pborder.SellerGigUpdate{
		Action:     pborder.Action_UpdateGigCount,
		SellerId:   newGig.SellerId,
		OrderProps: &pborder.OrderProps{GigCount: proto.Int32(1)},
	}

	if data, err := proto.Marshal(pbmsg); err == nil {
		if err := g.queue.PublishDirectMessage("jobber-seller-update", "user-seller", data); err != nil {
			slog.With("error", err).Error("Signup: Failed to send message to jobber-seller-update")
		}
	}

	return &gig.ResponseGig{Message: "Gig created successfully", Gig: newGig}, nil
}

func (g *GigService) UpdateGig(ctx context.Context, id string, req *gig.SellerGig) (*gig.ResponseGig, error) {
	// Validate Gig payload
	validateGig := &gig.UpdateSellerGig{
		Title:            req.Title,
		Description:      req.Description,
		Categories:       req.Categories,
		SubCategories:    req.SubCategories,
		Tags:             req.Tags,
		ExpectedDelivery: req.ExpectedDelivery,
		BasicTitle:       req.BasicTitle,
		BasicDescription: req.BasicDescription,
		Price:            req.Price,
	}

	if err := validateGig.Validate(validate); err != nil {
		slog.With("error", err).Debug("Gig Update Validation Failed")
		return nil, svc.NewError(svc.ErrInvalidData, err)
	}

	// upload gig image to cloudinary if the image isn't not a url
	if err := validate.Var(req.CoverImage, "required,http_url"); err != nil {
		imagePublicId := uuid.New().String()

		link, err := g.image.UploadImage(ctx, req.CoverImage, imagePublicId, true, true)
		if err != nil {
			return nil, svc.NewError(svc.ErrInternalFailure, err)
		}

		req.CoverImage = link
	}

	var updatedGig gig.SellerGig

	// convert id string to ObjectId
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, svc.NewError(svc.ErrBadRequest, fmt.Errorf("id provided is invalid"))
	}

	// update gig by id and get updated gig from mongo
	err = gigDB.FindOneAndUpdate(ctx, bson.M{"_id": objectId}, bson.D{{
		Key: "$set",
		Value: bson.D{
			{Key: "title", Value: req.Title},
			{Key: "description", Value: req.Description},
			{Key: "categories", Value: req.Categories},
			{Key: "subCategories", Value: req.SubCategories},
			{Key: "tags", Value: req.Tags},
			{Key: "price", Value: req.Price},
			{Key: "coverImage", Value: req.CoverImage},
			{Key: "expectedDelivery", Value: req.ExpectedDelivery},
			{Key: "basicTitle", Value: req.BasicTitle},
			{Key: "basicDescription", Value: req.BasicDescription},
		},
	}}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedGig)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, svc.NewError(svc.ErrNotFound, fmt.Errorf("gig does not exist"))
		}

		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	// update elasticsearch
	if _, err = g.search.UpdateGig(ctx, id, &updatedGig); err != nil {
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	updatedGig.ID = id

	return &gig.ResponseGig{Message: "Gig updated successfully", Gig: &updatedGig}, nil
}

func (g *GigService) DeleteGig(ctx context.Context, gigId string, sellerId string) (string, error) {
	// convert id string to ObjectId
	objectId, err := primitive.ObjectIDFromHex(gigId)
	if err != nil {
		return "", svc.NewError(svc.ErrBadRequest, fmt.Errorf("id provided is invalid"))
	}

	// delete gig from elastic first, this is more forgiving
	if err := g.search.DeleteGig(ctx, gigId); err != nil {
		return "", svc.NewError(svc.ErrInternalFailure, fmt.Errorf("failed to delete gig from index: %w", err))
	}

	// delete gig from mongodb
	result, err := gigDB.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		return "", svc.NewError(svc.ErrInternalFailure, err)
	}
	if result.DeletedCount == 0 {
		return "", svc.NewError(svc.ErrNotFound, fmt.Errorf("gig does not exist"))
	}

	// increment the seller gig count via the Users microservice from rabbitmq
	pbmsg := &pborder.SellerGigUpdate{
		Action:     pborder.Action_UpdateGigCount,
		SellerId:   sellerId,
		OrderProps: &pborder.OrderProps{GigCount: proto.Int32(-1)},
	}

	if data, err := proto.Marshal(pbmsg); err == nil {
		if err := g.queue.PublishDirectMessage("jobber-seller-update", "user-seller", data); err != nil {
			slog.With("error", err).Error("DeleteGig: Failed to send message to jobber-seller-update")
		}
	}

	return "Gig deleted successfully", nil
}

func (g *GigService) UpdateActiveGig(ctx context.Context, id string, active bool) (*gig.ResponseGig, error) {
	// convert id string to ObjectId
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, svc.NewError(svc.ErrBadRequest, fmt.Errorf("id provided is invalid"))
	}

	var updatedGig gig.SellerGig

	err = gigDB.FindOneAndUpdate(ctx, bson.M{"_id": objectId}, bson.D{{
		Key:   "$set",
		Value: bson.D{{Key: "active", Value: active}},
	}}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedGig)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, svc.NewError(svc.ErrNotFound, fmt.Errorf("gig does not exist"))
		}

		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	// update elasticsearch
	if _, err = g.search.UpdateGig(ctx, id, &updatedGig); err != nil {
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	msg := "Gig is now active"
	if !active {
		msg = "Gig is inactive"
	}

	return &gig.ResponseGig{Message: msg, Gig: &updatedGig}, nil
}
