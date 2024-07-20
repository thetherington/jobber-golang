package grpc

import (
	"context"

	"github.com/thetherington/jobber-common/error-handling/grpcerror"
	"github.com/thetherington/jobber-common/models/review"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/thetherington/jobber-common/protogen/go/review"
)

// parses service error and returns a equivelient gRPC error
func serviceError(err error) error {
	// try to cast the error to a grpcerror lookup
	if apiError, ok := grpcerror.FromError(err); ok {
		s := status.New(apiError.Status, apiError.Message)
		return s.Err()
	}

	// generic response
	s := status.New(codes.Internal, err.Error())
	return s.Err()
}

func (g *GrpcAdapter) CreateReview(ctx context.Context, req *pb.ReviewDocument) (*pb.ReviewResponse, error) {
	resp, err := g.reviewService.CreateReview(ctx, review.UnmarshalToDocument(req))
	if err != nil {
		return nil, serviceError(err)
	}

	return &pb.ReviewResponse{Message: resp.Message, Review: resp.Review.MarshalToProto()}, nil
}

func (g *GrpcAdapter) GetReviewsBySellerId(ctx context.Context, req *pb.RequestById) (*pb.ReviewsResponse, error) {
	resp, err := g.reviewService.GetReviewsBySellerId(ctx, req.Id)
	if err != nil {
		return nil, serviceError(err)
	}

	reviews := make([]*pb.ReviewDocument, 0)

	for _, r := range resp.Reviews {
		reviews = append(reviews, r.MarshalToProto())
	}

	return &pb.ReviewsResponse{Message: resp.Message, Reviews: reviews}, nil
}

func (g *GrpcAdapter) GetReviewsByGigId(ctx context.Context, req *pb.RequestById) (*pb.ReviewsResponse, error) {
	resp, err := g.reviewService.GetReviewsByGigId(ctx, req.Id)
	if err != nil {
		return nil, serviceError(err)
	}

	reviews := make([]*pb.ReviewDocument, 0)

	for _, r := range resp.Reviews {
		reviews = append(reviews, r.MarshalToProto())
	}

	return &pb.ReviewsResponse{Message: resp.Message, Reviews: reviews}, nil
}
