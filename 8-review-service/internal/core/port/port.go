package port

import (
	"context"

	"github.com/thetherington/jobber-common/models/review"
)

type ReviewRepository interface {
	CreateReview(ctx context.Context, review *review.ReviewDocument) (*review.ReviewDocument, error)
	GetReviewsByGigId(ctx context.Context, id string) ([]*review.ReviewDocument, error)
	GetReviewsBySellerId(ctx context.Context, id string) ([]*review.ReviewDocument, error)
}

type ReviewService interface {
	CreateReview(ctx context.Context, review *review.ReviewDocument) (*review.ReviewResponse, error)
	GetReviewsByGigId(ctx context.Context, id string) (*review.ReviewsResponse, error)
	GetReviewsBySellerId(ctx context.Context, id string) (*review.ReviewsResponse, error)
}

type ReviewProducer interface {
	PublishFanoutMessage(exchangeName string, data []byte) error
}
