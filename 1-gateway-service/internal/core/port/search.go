package port

import (
	"context"

	"github.com/thetherington/jobber-common/models/gig"
	"github.com/thetherington/jobber-common/models/search"
	pb "github.com/thetherington/jobber-common/protogen/go/auth"
	"google.golang.org/grpc"
)

type SearchService interface {
	GetGigByID(ctx context.Context, id string) (*gig.SellerGig, error)
	SearchGigs(ctx context.Context, req search.SearchRequest) (*search.SearchResponse, error)
}

type SearchRPCClient interface {
	GetGigById(ctx context.Context, in *pb.GetGigRequest, opts ...grpc.CallOption) (*pb.GigResponse, error)
	SearchGig(ctx context.Context, in *pb.SearchRequest, opts ...grpc.CallOption) (*pb.SearchResponse, error)
}
