package grpc

import (
	"context"

	"github.com/thetherington/jobber-common/error-handling/grpcerror"
	"github.com/thetherington/jobber-common/models/gig"
	pb "github.com/thetherington/jobber-common/protogen/go/gig"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (g *GrpcAdapter) CreateGig(ctx context.Context, req *pb.GigMessage) (*pb.GigResponse, error) {
	resp, err := g.gigService.CreateGig(ctx, gig.CreateSellerGig(req))
	if err != nil {
		return nil, serviceError(err)
	}

	return &pb.GigResponse{Message: resp.Message, Gig: gig.CreateGigMessage(resp.Gig)}, nil
}

func (g *GrpcAdapter) UpdateGig(ctx context.Context, req *pb.GigRequestUpdate) (*pb.GigResponse, error) {
	resp, err := g.gigService.UpdateGig(ctx, req.Id, gig.CreateSellerGig(req.Gig))
	if err != nil {
		return nil, serviceError(err)
	}

	return &pb.GigResponse{Message: resp.Message, Gig: gig.CreateGigMessage(resp.Gig)}, nil
}

func (g *GrpcAdapter) GetGigById(ctx context.Context, req *pb.GigRequestById) (*pb.GigResponse, error) {
	resp, err := g.gigService.GetGigById(ctx, req.Id)
	if err != nil {
		return nil, serviceError(err)
	}

	return &pb.GigResponse{Message: resp.Message, Gig: gig.CreateGigMessage(resp.Gig)}, nil
}

func (g *GrpcAdapter) DeleteGig(ctx context.Context, req *pb.GigDeleteRequest) (*pb.GigMessageResponse, error) {
	msg, err := g.gigService.DeleteGig(ctx, req.GigId, req.SellerId)
	if err != nil {
		return nil, serviceError(err)
	}

	return &pb.GigMessageResponse{Message: msg}, nil
}

func (g *GrpcAdapter) GetSellerGigs(ctx context.Context, req *pb.GigRequestById) (*pb.GigsResponse, error) {
	resp, err := g.gigService.GetSellerGigs(ctx, req.Id)
	if err != nil {
		return nil, serviceError(err)
	}

	var gigs []*pb.GigMessage
	for _, g := range resp.Gigs {
		gigs = append(gigs, gig.CreateGigMessage(g))
	}

	return &pb.GigsResponse{Message: resp.Message, Gigs: gigs}, nil
}

func (g *GrpcAdapter) GetSellerPausedGigs(ctx context.Context, req *pb.GigRequestById) (*pb.GigsResponse, error) {
	resp, err := g.gigService.GetSellerPausedGigs(ctx, req.Id)
	if err != nil {
		return nil, serviceError(err)
	}

	var gigs []*pb.GigMessage
	for _, g := range resp.Gigs {
		gigs = append(gigs, gig.CreateGigMessage(g))
	}

	return &pb.GigsResponse{Message: resp.Message, Gigs: gigs}, nil
}

func (g *GrpcAdapter) UpdateActiveGig(ctx context.Context, req *pb.GigUpdateActive) (*pb.GigResponse, error) {
	resp, err := g.gigService.UpdateActiveGig(ctx, req.GigId, req.Active)
	if err != nil {
		return nil, serviceError(err)
	}

	return &pb.GigResponse{Message: resp.Message, Gig: gig.CreateGigMessage(resp.Gig)}, nil
}

func (g *GrpcAdapter) SeedGigs(ctx context.Context, req *pb.GigSeedRequest) (*pb.GigMessageResponse, error) {
	msg, err := g.gigService.SeedGigs(ctx, req.Count)
	if err != nil {
		return nil, serviceError(err)
	}

	return &pb.GigMessageResponse{Message: msg}, nil
}
