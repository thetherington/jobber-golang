package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/thetherington/jobber-auth/internal/adapters/storage/elasticsearch"
	"github.com/thetherington/jobber-auth/internal/core/port"
	"github.com/thetherington/jobber-common/error-handling/svc"
	"github.com/thetherington/jobber-common/models/gig"
	"github.com/thetherington/jobber-common/models/search"
)

/**
 * SearchService implements
 */
type SearchService struct {
	client port.SearchClient
}

// NewSearchService creates a new search service instance
func NewSearchService(client port.SearchClient) *SearchService {
	return &SearchService{
		client,
	}
}

func (s *SearchService) GetGigByID(ctx context.Context, id string) (*gig.SellerGig, error) {
	if id == "" {
		slog.Debug("GetGigByID invalid id")
		return nil, svc.NewError(svc.ErrInvalidData, fmt.Errorf("gig id invalid"))
	}

	gig, err := s.client.GetDocumentById(ctx, id)
	if err != nil {
		if errors.Is(err, elasticsearch.ErrGigNotFound) {
			return nil, svc.NewError(svc.ErrNotFound, fmt.Errorf("gig id (%s) not found", id))
		}

		slog.With("error", err).Error("GetDocumentById failed", "gigId", id)
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	return gig, nil
}

func (s *SearchService) SearchGigs(ctx context.Context, req search.SearchRequest) (*search.SearchResponse, error) {
	hits, gigs, err := s.client.SearchGigs(ctx, req.SearchQuery, req.PaginateProps, req.DeliveryTime, req.Min, req.Max)
	if err != nil {
		slog.With("error", err).Error("failed to search gigs")
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	return &search.SearchResponse{
		Total: hits,
		Hits:  gigs,
	}, nil
}
