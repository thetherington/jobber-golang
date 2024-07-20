package grpc

import (
	"context"

	"github.com/thetherington/jobber-common/models/gig"
	"github.com/thetherington/jobber-common/models/search"
	pb "github.com/thetherington/jobber-common/protogen/go/auth"
	pbGig "github.com/thetherington/jobber-common/protogen/go/gig"
	pbReview "github.com/thetherington/jobber-common/protogen/go/review"
	"github.com/thetherington/jobber-common/utils"
)

func createProtobufGig(gig *gig.SellerGig) *pbGig.GigMessage {
	g := pbGig.GigMessage{
		ES_ID:            &gig.ES_Id,
		ID:               &gig.ID,
		SellerId:         &gig.SellerId,
		Title:            &gig.Title,
		Username:         &gig.Username,
		ProfilePicture:   &gig.ProfilePicture,
		Email:            &gig.Email,
		Description:      gig.Description,
		Active:           gig.Active,
		Categories:       gig.Categories,
		SubCategories:    gig.SubCategories,
		Tags:             gig.Tags,
		RatingsCount:     gig.RatingsCount,
		RatingSum:        gig.RatingSum,
		ExpectedDelivery: gig.ExpectedDelivery,
		BasicTitle:       gig.BasicTitle,
		BasicDescription: gig.BasicDescription,
		Price:            gig.Price,
		CoverImage:       gig.CoverImage,
		CreatedAt:        utils.ToDateTime(gig.CreatedAt),
		SortId:           gig.SortId,
	}

	if gig.RatingCategories != nil {
		rc := gig.RatingCategories

		g.RatingCategories = &pbReview.RatingCategories{
			One:   &pbReview.RatingCategoryItem{Value: rc.One.Value, Count: rc.One.Count},
			Two:   &pbReview.RatingCategoryItem{Value: rc.Two.Value, Count: rc.Two.Count},
			Three: &pbReview.RatingCategoryItem{Value: rc.Three.Value, Count: rc.Three.Count},
			Four:  &pbReview.RatingCategoryItem{Value: rc.Four.Value, Count: rc.Four.Count},
			Five:  &pbReview.RatingCategoryItem{Value: rc.Five.Value, Count: rc.Five.Count},
		}
	}

	return &g
}

func (a *GrpcAdapter) GetGigById(ctx context.Context, req *pb.GetGigRequest) (*pb.GigResponse, error) {
	resp, err := a.searchService.GetGigByID(ctx, req.GetId())
	if err != nil {
		return nil, serviceError(err)
	}

	return &pb.GigResponse{
		Gig: createProtobufGig(resp),
	}, nil
}

func (a *GrpcAdapter) SearchGig(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
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

	resp, err := a.searchService.SearchGigs(ctx, searchRequest)
	if err != nil {
		return nil, serviceError(err)
	}

	hits := make([]*pbGig.GigMessage, 0)

	for _, h := range resp.Hits {
		hits = append(hits, createProtobufGig(h))
	}

	return &pb.SearchResponse{
		Total: resp.Total,
		Hits:  hits,
	}, nil
}
