package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/thetherington/jobber-common/error-handling/svc"
	"github.com/thetherington/jobber-common/models/gig"
	"github.com/thetherington/jobber-gig/internal/adapters/storage/elasticsearch"
)

func (g *GigService) GetGigById(ctx context.Context, id string) (*gig.ResponseGig, error) {
	result, err := g.search.GetDocumentById(ctx, id)
	if err != nil {
		if errors.Is(err, elasticsearch.ErrGigNotFound) {
			return nil, svc.NewError(svc.ErrNotFound, fmt.Errorf("gig id (%s) not found", id))
		}

		slog.With("error", err).Error("GetDocumentById failed", "gigId", id)
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	return &gig.ResponseGig{Message: "Gig by id", Gig: result}, nil
}

func (g *GigService) GetSellerGigs(ctx context.Context, id string) (*gig.ResponseGigs, error) {
	gigs, err := g.search.GigsSearchBySellerId(ctx, id, true)
	if err != nil {
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	return &gig.ResponseGigs{Message: "Seller Gigs", Gigs: gigs}, nil
}

func (g *GigService) GetSellerPausedGigs(ctx context.Context, id string) (*gig.ResponseGigs, error) {
	gigs, err := g.search.GigsSearchBySellerId(ctx, id, false)
	if err != nil {
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	return &gig.ResponseGigs{Message: "Seller Paused Gigs", Gigs: gigs}, nil
}
