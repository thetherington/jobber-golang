package port

import (
	"context"

	"github.com/thetherington/jobber-common/models/gig"
	"github.com/thetherington/jobber-common/models/review"
	"github.com/thetherington/jobber-common/models/search"
)

type GigService interface {
	CreateGig(ctx context.Context, req *gig.SellerGig) (*gig.ResponseGig, error)
	DeleteGig(ctx context.Context, gigId string, sellerId string) (string, error)
	GetGigById(ctx context.Context, id string) (*gig.ResponseGig, error)
	GetSellerGigs(ctx context.Context, id string) (*gig.ResponseGigs, error)
	GetSellerPausedGigs(ctx context.Context, id string) (*gig.ResponseGigs, error)
	SearchGig(ctx context.Context, req search.SearchRequest) (*gig.ResponseSearchGigs, error)
	SearchGigCategory(ctx context.Context, username string) (*gig.ResponseSearchGigs, error)
	SearchGigSimilar(ctx context.Context, id string) (*gig.ResponseSearchGigs, error)
	SearchGigTop(ctx context.Context, username string) (*gig.ResponseSearchGigs, error)
	SeedGigs(ctx context.Context, count int32) (string, error)
	SeedData(ctx context.Context, sellers []any) error
	UpdateActiveGig(ctx context.Context, id string, active bool) (*gig.ResponseGig, error)
	UpdateGig(ctx context.Context, id string, req *gig.SellerGig) (*gig.ResponseGig, error)
	UpdateGigReview(data *review.ReviewMessageDetails) error
}

type GigProducer interface {
	PublishDirectMessage(exchangeName string, routingKey string, data []byte) error
}

type ImageUploader interface {
	UploadImage(ctx context.Context, file string, public_id string, overwrite bool, invalidate bool) (string, error)
	UploadVideo(ctx context.Context, file string, public_id string, overwrite bool, invalidate bool) (string, error)
}

type SearchClient interface {
	GetDocumentById(ctx context.Context, id string) (*gig.SellerGig, error)
	GetGigsCount(ctx context.Context) (int32, error)
	InsertGig(ctx context.Context, newGig *gig.SellerGig) (string, error)
	UpdateGig(ctx context.Context, id string, updateGig *gig.SellerGig) (string, error)
	DeleteGig(ctx context.Context, id string) error
	GigsSearchBySellerId(ctx context.Context, id string, active bool) ([]*gig.SellerGig, error)
	SearchGigs(ctx context.Context, searchQuery string, paginate *search.PaginateProps, deliveryTime *string, min *float64, max *float64) (int64, []*gig.SellerGig, error)
	SearchGigsByCategory(ctx context.Context, category string) (int64, []*gig.SellerGig, error)
	SearchSimiliarGigs(ctx context.Context, id string) (int64, []*gig.SellerGig, error)
	SearchTopRatedGigsbyCategory(ctx context.Context, category string) (int64, []*gig.SellerGig, error)
}

type CacheRepository interface {
	GetUserSelectedGigCategory(ctx context.Context, username string) (string, error)
}
