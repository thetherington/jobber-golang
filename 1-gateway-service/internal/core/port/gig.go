package port

import (
	"context"

	"github.com/thetherington/jobber-common/models/gig"
	"github.com/thetherington/jobber-common/models/search"
	pb "github.com/thetherington/jobber-common/protogen/go/gig"
	"google.golang.org/grpc"
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
	UpdateActiveGig(ctx context.Context, id string, active bool) (*gig.ResponseGig, error)
	UpdateGig(ctx context.Context, id string, req *gig.SellerGig) (*gig.ResponseGig, error)
}

type GigRPCClient interface {
	CreateGig(ctx context.Context, in *pb.GigMessage, opts ...grpc.CallOption) (*pb.GigResponse, error)
	DeleteGig(ctx context.Context, in *pb.GigDeleteRequest, opts ...grpc.CallOption) (*pb.GigMessageResponse, error)
	GetGigById(ctx context.Context, in *pb.GigRequestById, opts ...grpc.CallOption) (*pb.GigResponse, error)
	GetSellerGigs(ctx context.Context, in *pb.GigRequestById, opts ...grpc.CallOption) (*pb.GigsResponse, error)
	GetSellerPausedGigs(ctx context.Context, in *pb.GigRequestById, opts ...grpc.CallOption) (*pb.GigsResponse, error)
	SearchGig(ctx context.Context, in *pb.GigSearchRequest, opts ...grpc.CallOption) (*pb.SearchResponse, error)
	SearchGigCategory(ctx context.Context, in *pb.SearchGigByValue, opts ...grpc.CallOption) (*pb.SearchResponse, error)
	SearchGigSimilar(ctx context.Context, in *pb.SearchGigByValue, opts ...grpc.CallOption) (*pb.SearchResponse, error)
	SearchGigTop(ctx context.Context, in *pb.SearchGigByValue, opts ...grpc.CallOption) (*pb.SearchResponse, error)
	SeedGigs(ctx context.Context, in *pb.GigSeedRequest, opts ...grpc.CallOption) (*pb.GigMessageResponse, error)
	UpdateActiveGig(ctx context.Context, in *pb.GigUpdateActive, opts ...grpc.CallOption) (*pb.GigResponse, error)
	UpdateGig(ctx context.Context, in *pb.GigRequestUpdate, opts ...grpc.CallOption) (*pb.GigResponse, error)
}
