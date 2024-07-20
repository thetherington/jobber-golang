package service

import (
	"context"
	"log/slog"

	"github.com/thetherington/jobber-common/error-handling/svc"
	"github.com/thetherington/jobber-common/models/gig"
	"github.com/thetherington/jobber-common/models/search"
	pb "github.com/thetherington/jobber-common/protogen/go/gig"
	"github.com/thetherington/jobber-gateway/internal/core/port"
)

/**
 * GigService implements
 */
type GigService struct {
	client port.GigRPCClient
}

// NewGigService creates a new gig service instance
func NewGigService(rpc port.GigRPCClient) *GigService {
	return &GigService{
		rpc,
	}
}

func (g *GigService) CreateGig(ctx context.Context, req *gig.SellerGig) (*gig.ResponseGig, error) {
	resp, err := g.client.CreateGig(ctx, gig.CreateGigMessage(req))
	if err != nil {
		slog.With("error", err).Debug("CreateGig error")
		return nil, svc.GrpcErrorResolve(err, "CreateGig")
	}

	return &gig.ResponseGig{Message: resp.Message, Gig: gig.CreateSellerGig(resp.Gig)}, nil
}

func (g *GigService) UpdateGig(ctx context.Context, id string, req *gig.SellerGig) (*gig.ResponseGig, error) {
	resp, err := g.client.UpdateGig(ctx, &pb.GigRequestUpdate{Id: id, Gig: gig.CreateGigMessage(req)})
	if err != nil {
		slog.With("error", err).Debug("UpdateGig error")
		return nil, svc.GrpcErrorResolve(err, "UpdateGig")
	}

	return &gig.ResponseGig{Message: resp.Message, Gig: gig.CreateSellerGig(resp.Gig)}, nil
}

func (g *GigService) DeleteGig(ctx context.Context, gigId string, sellerId string) (string, error) {
	resp, err := g.client.DeleteGig(ctx, &pb.GigDeleteRequest{GigId: gigId, SellerId: sellerId})
	if err != nil {
		slog.With("error", err).Debug("DeleteGig error")
		return "", svc.GrpcErrorResolve(err, "DeleteGig")
	}

	return resp.Message, nil
}

func (g *GigService) UpdateActiveGig(ctx context.Context, id string, active bool) (*gig.ResponseGig, error) {
	resp, err := g.client.UpdateActiveGig(ctx, &pb.GigUpdateActive{GigId: id, Active: active})
	if err != nil {
		slog.With("error", err).Debug("UpdateActiveGig error")
		return nil, svc.GrpcErrorResolve(err, "UpdateActiveGig")
	}

	return &gig.ResponseGig{Message: resp.Message, Gig: gig.CreateSellerGig(resp.Gig)}, nil
}

func (g *GigService) SeedGigs(ctx context.Context, count int32) (string, error) {
	resp, err := g.client.SeedGigs(ctx, &pb.GigSeedRequest{Count: count})
	if err != nil {
		slog.With("error", err).Debug("SeedGigs error")
		return "", svc.GrpcErrorResolve(err, "SeedGigs")
	}

	return resp.Message, nil
}

func (g *GigService) GetGigById(ctx context.Context, id string) (*gig.ResponseGig, error) {
	resp, err := g.client.GetGigById(ctx, &pb.GigRequestById{Id: id})
	if err != nil {
		slog.With("error", err).Debug("GetGigById error")
		return nil, svc.GrpcErrorResolve(err, "GetGigById")
	}

	return &gig.ResponseGig{Message: resp.Message, Gig: gig.CreateSellerGig(resp.Gig)}, nil
}

func (g *GigService) GetSellerGigs(ctx context.Context, id string) (*gig.ResponseGigs, error) {
	resp, err := g.client.GetSellerGigs(ctx, &pb.GigRequestById{Id: id})
	if err != nil {
		slog.With("error", err).Debug("GetSellerGigs error")
		return nil, svc.GrpcErrorResolve(err, "GetSellerGigs")
	}

	gigs := make([]*gig.SellerGig, 0)
	for _, g := range resp.Gigs {
		gigs = append(gigs, gig.CreateSellerGig(g))
	}

	return &gig.ResponseGigs{Message: resp.Message, Gigs: gigs}, nil
}

func (g *GigService) GetSellerPausedGigs(ctx context.Context, id string) (*gig.ResponseGigs, error) {
	resp, err := g.client.GetSellerPausedGigs(ctx, &pb.GigRequestById{Id: id})
	if err != nil {
		slog.With("error", err).Debug("GetSellerPausedGigs error")
		return nil, svc.GrpcErrorResolve(err, "GetSellerPausedGigs")
	}

	gigs := make([]*gig.SellerGig, 0)
	for _, g := range resp.Gigs {
		gigs = append(gigs, gig.CreateSellerGig(g))
	}

	return &gig.ResponseGigs{Message: resp.Message, Gigs: gigs}, nil
}

func (g *GigService) SearchGig(ctx context.Context, req search.SearchRequest) (*gig.ResponseSearchGigs, error) {
	protoRequest := pb.GigSearchRequest{
		SearchQuery: req.SearchQuery,
	}

	if req.PaginateProps != nil {
		protoRequest.PaginateProps = &pb.PaginateProps{
			From: req.PaginateProps.From,
			Size: int32(req.PaginateProps.Size),
			Type: req.PaginateProps.Type,
		}
	}

	protoRequest.DeliveryTime = req.DeliveryTime
	protoRequest.Min = req.Min
	protoRequest.Max = req.Max

	resp, err := g.client.SearchGig(ctx, &protoRequest)
	if err != nil {
		slog.With("error", err).Debug("SearchGig failed")
		return nil, svc.GrpcErrorResolve(err, "SearchGig")
	}

	gigs := make([]*gig.SellerGig, 0)
	for _, g := range resp.Gigs {
		gigs = append(gigs, gig.CreateSellerGig(g))
	}

	return &gig.ResponseSearchGigs{Message: resp.Message, Total: resp.Total, Gigs: gigs}, nil
}

func (g *GigService) SearchGigCategory(ctx context.Context, username string) (*gig.ResponseSearchGigs, error) {
	resp, err := g.client.SearchGigCategory(ctx, &pb.SearchGigByValue{Value: username})
	if err != nil {
		slog.With("error", err).Debug("SearchGigCategory failed")
		return nil, svc.GrpcErrorResolve(err, "SearchGigCategory")
	}

	gigs := make([]*gig.SellerGig, 0)
	for _, g := range resp.Gigs {
		gigs = append(gigs, gig.CreateSellerGig(g))
	}

	return &gig.ResponseSearchGigs{Message: resp.Message, Total: resp.Total, Gigs: gigs}, nil
}

func (g *GigService) SearchGigTop(ctx context.Context, username string) (*gig.ResponseSearchGigs, error) {
	resp, err := g.client.SearchGigTop(ctx, &pb.SearchGigByValue{Value: username})
	if err != nil {
		slog.With("error", err).Debug("SearchGigTop failed")
		return nil, svc.GrpcErrorResolve(err, "SearchGigTop")
	}

	gigs := make([]*gig.SellerGig, 0)
	for _, g := range resp.Gigs {
		gigs = append(gigs, gig.CreateSellerGig(g))
	}

	return &gig.ResponseSearchGigs{Message: resp.Message, Total: resp.Total, Gigs: gigs}, nil
}

func (g *GigService) SearchGigSimilar(ctx context.Context, id string) (*gig.ResponseSearchGigs, error) {
	resp, err := g.client.SearchGigSimilar(ctx, &pb.SearchGigByValue{Value: id})
	if err != nil {
		slog.With("error", err).Debug("SearchGigSimilar failed")
		return nil, svc.GrpcErrorResolve(err, "SearchGigSimilar")
	}

	gigs := make([]*gig.SellerGig, 0)
	for _, g := range resp.Gigs {
		gigs = append(gigs, gig.CreateSellerGig(g))
	}

	return &gig.ResponseSearchGigs{Message: resp.Message, Total: resp.Total, Gigs: gigs}, nil
}
