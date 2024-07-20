package service

import (
	"context"
	"log/slog"

	"github.com/thetherington/jobber-common/error-handling/svc"
	"github.com/thetherington/jobber-common/models/gig"
	"github.com/thetherington/jobber-common/models/search"
)

func (g *GigService) SearchGig(ctx context.Context, req search.SearchRequest) (*gig.ResponseSearchGigs, error) {
	hits, gigs, err := g.search.SearchGigs(ctx, req.SearchQuery, req.PaginateProps, req.DeliveryTime, req.Min, req.Max)
	if err != nil {
		slog.With("error", err).Error("failed to search gigs")
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	return &gig.ResponseSearchGigs{Message: "Search gig results", Total: int32(hits), Gigs: gigs}, nil
}

func (g *GigService) SearchGigCategory(ctx context.Context, username string) (*gig.ResponseSearchGigs, error) {
	category, err := g.cache.GetUserSelectedGigCategory(ctx, username)
	if err != nil {
		slog.With("error", err).Error("failure getting user selected category from redis")
	}

	hits, gigs, err := g.search.SearchGigsByCategory(ctx, category)
	if err != nil {
		slog.With("error", err).Error("failed to search gigs by category ")
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	return &gig.ResponseSearchGigs{Message: "Search gigs category results", Total: int32(hits), Gigs: gigs}, nil
}

func (g *GigService) SearchGigTop(ctx context.Context, username string) (*gig.ResponseSearchGigs, error) {
	category, err := g.cache.GetUserSelectedGigCategory(ctx, username)
	if err != nil {
		slog.With("error", err).Error("failure getting user selected category from redis")
	}

	hits, gigs, err := g.search.SearchTopRatedGigsbyCategory(ctx, category)
	if err != nil {
		slog.With("error", err).Error("failed to search top gigs by category ")
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	return &gig.ResponseSearchGigs{Message: "Search top gigs category results", Total: int32(hits), Gigs: gigs}, nil
}

func (g *GigService) SearchGigSimilar(ctx context.Context, id string) (*gig.ResponseSearchGigs, error) {
	hits, gigs, err := g.search.SearchSimiliarGigs(ctx, id)
	if err != nil {
		slog.With("error", err).Error("failed to search similiar gigs ")
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	return &gig.ResponseSearchGigs{Message: "More gigs like this results", Total: int32(hits), Gigs: gigs}, nil
}
