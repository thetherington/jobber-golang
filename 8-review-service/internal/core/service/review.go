package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/thetherington/jobber-common/error-handling/svc"
	"github.com/thetherington/jobber-common/models/review"
	pb "github.com/thetherington/jobber-common/protogen/go/review"
	"github.com/thetherington/jobber-common/utils"
	"github.com/thetherington/jobber-review/internal/adapters/storage/postgres"
	"github.com/thetherington/jobber-review/internal/core/port"
	"google.golang.org/protobuf/proto"
)

var (
	validate *validator.Validate
)

/**
 * ReviewService implements
 */
type ReviewService struct {
	repo  port.ReviewRepository
	queue port.ReviewProducer
}

// NewReviewService creates a new review service instance
func NewReviewService(repo port.ReviewRepository, queue port.ReviewProducer) *ReviewService {
	validate = validator.New(validator.WithRequiredStructEnabled())

	return &ReviewService{
		repo,
		queue,
	}
}

func (rs *ReviewService) CreateReview(ctx context.Context, payload *review.ReviewDocument) (*review.ReviewResponse, error) {
	// Validate Review payload
	if err := payload.Validate(validate); err != nil {
		slog.With("error", err).Debug("Review Create Validation Failed")
		return nil, svc.NewError(svc.ErrInvalidData, err)
	}

	resp, err := rs.repo.CreateReview(ctx, payload)
	if err != nil {
		if errors.Is(err, postgres.ErrConflictingData) {
			return nil, svc.NewError(svc.ErrBadRequest, err)
		}

		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	pbMsg := &pb.ReviewMessageDetails{
		GigId:      resp.GigId,
		ReviewerId: resp.ReviewerId,
		SellerId:   resp.SellerId,
		Review:     resp.Review,
		Rating:     resp.Rating,
		OrderId:    resp.OrderId,
		Action:     pb.ReviewType_BuyerReview,
		CreatedAt:  utils.ToDateTime(&resp.CreatedAt),
	}

	if resp.ReviewType == "seller-review" {
		pbMsg.Action = pb.ReviewType_SellerReview
	}

	// send the review information to the order and user microservice via a fanout exchange.
	if data, err := proto.Marshal(pbMsg); err == nil {
		if err := rs.queue.PublishFanoutMessage("jobber-review", data); err != nil {
			slog.With("error", err).Error("CreateReview: Failed to send offer to jobber-review")
		}
	}

	return &review.ReviewResponse{Message: "Review message created", Review: resp}, nil
}

func (rs *ReviewService) GetReviewsByGigId(ctx context.Context, id string) (*review.ReviewsResponse, error) {
	reviews, err := rs.repo.GetReviewsByGigId(ctx, id)
	if err != nil {
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	return &review.ReviewsResponse{Message: "Review messages by gig id", Reviews: reviews}, nil
}

func (rs *ReviewService) GetReviewsBySellerId(ctx context.Context, id string) (*review.ReviewsResponse, error) {
	reviews, err := rs.repo.GetReviewsBySellerId(ctx, id)
	if err != nil {
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	return &review.ReviewsResponse{Message: "Review messages by gig id", Reviews: reviews}, nil
}
