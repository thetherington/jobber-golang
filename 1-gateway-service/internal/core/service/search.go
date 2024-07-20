package service

import (
	"context"
	"log/slog"

	"github.com/thetherington/jobber-common/error-handling/svc"
	"github.com/thetherington/jobber-common/models/gig"
	"github.com/thetherington/jobber-common/models/review"
	"github.com/thetherington/jobber-common/models/search"
	pb "github.com/thetherington/jobber-common/protogen/go/auth"
	pbGig "github.com/thetherington/jobber-common/protogen/go/gig"
	"github.com/thetherington/jobber-common/utils"
	"github.com/thetherington/jobber-gateway/internal/core/port"
)

/**
 * SearchService implements
 */
type SearchService struct {
	client port.SearchRPCClient
}

// NewSearchServicee creates a new search service instance
func NewSearchService(rpc port.SearchRPCClient) *SearchService {
	return &SearchService{
		rpc,
	}
}

func createModelGig(gigpb *pbGig.GigMessage) *gig.SellerGig {
	g := gig.SellerGig{
		ES_Id:            *gigpb.ES_ID,
		ID:               *gigpb.ID,
		SellerId:         *gigpb.SellerId,
		Title:            *gigpb.Title,
		Username:         *gigpb.Username,
		ProfilePicture:   *gigpb.ProfilePicture,
		Email:            *gigpb.Email,
		Description:      gigpb.Description,
		Active:           gigpb.Active,
		Categories:       gigpb.Categories,
		SubCategories:    gigpb.SubCategories,
		Tags:             gigpb.Tags,
		RatingsCount:     gigpb.RatingsCount,
		RatingSum:        gigpb.RatingSum,
		ExpectedDelivery: gigpb.ExpectedDelivery,
		BasicTitle:       gigpb.BasicTitle,
		BasicDescription: gigpb.BasicDescription,
		Price:            gigpb.Price,
		CoverImage:       gigpb.CoverImage,
		CreatedAt:        utils.ToTime(gigpb.GetCreatedAt()),
		SortId:           gigpb.SortId,
	}

	if gigpb.RatingCategories != nil {
		rc := gigpb.RatingCategories

		g.RatingCategories = &review.RatingCategories{
			One:   review.RatingCategoryItem{Value: rc.One.Value, Count: rc.One.Count},
			Two:   review.RatingCategoryItem{Value: rc.Two.Value, Count: rc.Two.Count},
			Three: review.RatingCategoryItem{Value: rc.Three.Value, Count: rc.Three.Count},
			Four:  review.RatingCategoryItem{Value: rc.Four.Value, Count: rc.Four.Count},
			Five:  review.RatingCategoryItem{Value: rc.Five.Value, Count: rc.Five.Count},
		}
	}

	return &g
}

func (s *SearchService) GetGigByID(ctx context.Context, id string) (*gig.SellerGig, error) {
	resp, err := s.client.GetGigById(ctx, &pb.GetGigRequest{
		Id: id,
	})
	if err != nil {
		slog.With("error", err).Debug("get gig by id error")
		return nil, svc.GrpcErrorResolve(err, "getGigByID")
	}

	return createModelGig(resp.Gig), nil
}

func (s *SearchService) SearchGigs(ctx context.Context, req search.SearchRequest) (*search.SearchResponse, error) {
	protoRequest := pb.SearchRequest{
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

	resp, err := s.client.SearchGig(ctx, &protoRequest)
	if err != nil {
		slog.With("error", err).Debug("search gigs failed")
		return nil, svc.GrpcErrorResolve(err, "searchgigs")
	}

	hits := make([]*gig.SellerGig, 0)

	for _, hit := range resp.Hits {
		hits = append(hits, createModelGig(hit))
	}

	return &search.SearchResponse{
		Total: resp.Total,
		Hits:  hits,
	}, nil
}
