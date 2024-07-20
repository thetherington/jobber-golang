package port

import (
	"context"

	"github.com/thetherington/jobber-common/models/review"
	pb "github.com/thetherington/jobber-common/protogen/go/review"
	"google.golang.org/grpc"
)

type ReviewService interface {
	CreateReview(ctx context.Context, review *review.ReviewDocument) (*review.ReviewResponse, error)
	GetReviewsByGigId(ctx context.Context, id string) (*review.ReviewsResponse, error)
	GetReviewsBySellerId(ctx context.Context, id string) (*review.ReviewsResponse, error)
}

type ReviewRPCClient interface {
	CreateReview(ctx context.Context, in *pb.ReviewDocument, opts ...grpc.CallOption) (*pb.ReviewResponse, error)
	GetReviewsBySellerId(ctx context.Context, in *pb.RequestById, opts ...grpc.CallOption) (*pb.ReviewsResponse, error)
	GetReviewsByGigId(ctx context.Context, in *pb.RequestById, opts ...grpc.CallOption) (*pb.ReviewsResponse, error)
}
