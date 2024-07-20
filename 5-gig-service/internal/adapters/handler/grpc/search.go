package grpc

import (
	"context"

	"github.com/thetherington/jobber-common/models/gig"
	"github.com/thetherington/jobber-common/models/search"
	pb "github.com/thetherington/jobber-common/protogen/go/gig"
)

func (g *GrpcAdapter) SearchGig(ctx context.Context, req *pb.GigSearchRequest) (*pb.SearchResponse, error) {
	var searchRequest search.SearchRequest

	searchRequest.SearchQuery = req.GetSearchQuery()

	if req.PaginateProps != nil {
		searchRequest.PaginateProps = &search.PaginateProps{
			From: req.PaginateProps.GetFrom(),
			Size: int(req.PaginateProps.GetSize()),
			Type: req.PaginateProps.GetType(),
		}
	}

	searchRequest.DeliveryTime = req.DeliveryTime
	searchRequest.Min = req.Min
	searchRequest.Max = req.Max

	resp, err := g.gigService.SearchGig(ctx, searchRequest)
	if err != nil {
		return nil, serviceError(err)
	}

	var gigs []*pb.GigMessage
	for _, g := range resp.Gigs {
		gigs = append(gigs, gig.CreateGigMessage(g))
	}

	return &pb.SearchResponse{Message: resp.Message, Total: resp.Total, Gigs: gigs}, nil
}

func (g *GrpcAdapter) SearchGigCategory(ctx context.Context, req *pb.SearchGigByValue) (*pb.SearchResponse, error) {
	resp, err := g.gigService.SearchGigCategory(ctx, req.Value)
	if err != nil {
		return nil, serviceError(err)
	}

	var gigs []*pb.GigMessage
	for _, g := range resp.Gigs {
		gigs = append(gigs, gig.CreateGigMessage(g))
	}

	return &pb.SearchResponse{Message: resp.Message, Total: resp.Total, Gigs: gigs}, nil
}

func (g *GrpcAdapter) SearchGigTop(ctx context.Context, req *pb.SearchGigByValue) (*pb.SearchResponse, error) {
	resp, err := g.gigService.SearchGigTop(ctx, req.Value)
	if err != nil {
		return nil, serviceError(err)
	}

	var gigs []*pb.GigMessage
	for _, g := range resp.Gigs {
		gigs = append(gigs, gig.CreateGigMessage(g))
	}

	return &pb.SearchResponse{Message: resp.Message, Total: resp.Total, Gigs: gigs}, nil
}

func (g *GrpcAdapter) SearchGigSimilar(ctx context.Context, req *pb.SearchGigByValue) (*pb.SearchResponse, error) {
	resp, err := g.gigService.SearchGigSimilar(ctx, req.Value)
	if err != nil {
		return nil, serviceError(err)
	}

	var gigs []*pb.GigMessage
	for _, g := range resp.Gigs {
		gigs = append(gigs, gig.CreateGigMessage(g))
	}

	return &pb.SearchResponse{Message: resp.Message, Total: resp.Total, Gigs: gigs}, nil
}
