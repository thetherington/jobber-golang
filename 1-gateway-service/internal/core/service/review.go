package service

import (
	"context"

	"github.com/thetherington/jobber-common/error-handling/svc"
	"github.com/thetherington/jobber-common/models/review"
	pb "github.com/thetherington/jobber-common/protogen/go/review"
	"github.com/thetherington/jobber-gateway/internal/core/port"
)

/**
 * ReviewService implements
 */
type ReviewService struct {
	client port.ReviewRPCClient
}

// NewReviewService creates a new review service instance
func NewReviewService(rpc port.ReviewRPCClient) *ReviewService {
	return &ReviewService{
		rpc,
	}
}

func (rs *ReviewService) CreateReview(ctx context.Context, payload *review.ReviewDocument) (*review.ReviewResponse, error) {
	resp, err := rs.client.CreateReview(ctx, payload.MarshalToProto())
	if err != nil {
		return nil, svc.GrpcErrorResolve(err, "CreateReview")
	}

	return &review.ReviewResponse{Message: resp.Message, Review: review.UnmarshalToDocument(resp.Review)}, nil
}

func (rs *ReviewService) GetReviewsByGigId(ctx context.Context, id string) (*review.ReviewsResponse, error) {
	resp, err := rs.client.GetReviewsByGigId(ctx, &pb.RequestById{Id: id})
	if err != nil {
		return nil, svc.GrpcErrorResolve(err, "GetReviewsByGigId")
	}

	reviews := make([]*review.ReviewDocument, 0)

	for _, r := range resp.Reviews {
		reviews = append(reviews, review.UnmarshalToDocument(r))
	}

	return &review.ReviewsResponse{Message: resp.Message, Reviews: reviews}, nil
}

func (rs *ReviewService) GetReviewsBySellerId(ctx context.Context, id string) (*review.ReviewsResponse, error) {
	resp, err := rs.client.GetReviewsBySellerId(ctx, &pb.RequestById{Id: id})
	if err != nil {
		return nil, svc.GrpcErrorResolve(err, "GetReviewsBySellerId")
	}

	reviews := make([]*review.ReviewDocument, 0)

	for _, r := range resp.Reviews {
		reviews = append(reviews, review.UnmarshalToDocument(r))
	}

	return &review.ReviewsResponse{Message: resp.Message, Reviews: reviews}, nil
}
