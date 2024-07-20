package port

import (
	"context"

	"github.com/thetherington/jobber-common/models/users"
	pb "github.com/thetherington/jobber-common/protogen/go/users"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type BuyerRPCClient interface {
	GetBuyerByEmail(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*pb.BuyerPayload, error)
	GetBuyerByUsername(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*pb.BuyerPayload, error)
	GetBuyerByProvidedUsername(ctx context.Context, in *pb.BuyerUsernameRequest, opts ...grpc.CallOption) (*pb.BuyerPayload, error)
}

type BuyerService interface {
	GetBuyerByEmail(context.Context) (*users.BuyerResponse, error)
	GetBuyerByUsername(context.Context) (*users.BuyerResponse, error)
	GetBuyerByProvidedUsername(ctx context.Context, username string) (*users.BuyerResponse, error)
}

type SellerRPCClient interface {
	CreateSeller(ctx context.Context, in *pb.CreateUpdateSellerPayload, opts ...grpc.CallOption) (*pb.SellerResponse, error)
	UpdateSeller(ctx context.Context, in *pb.UpdateSellerRequest, opts ...grpc.CallOption) (*pb.SellerResponse, error)
	GetSellerById(ctx context.Context, in *pb.GetSellerByIdRequest, opts ...grpc.CallOption) (*pb.SellerResponse, error)
	GetSellerByUsername(ctx context.Context, in *pb.GetSellerByUsernameRequest, opts ...grpc.CallOption) (*pb.SellerResponse, error)
	GetRandomSellers(ctx context.Context, in *pb.RandomSellersRequest, opts ...grpc.CallOption) (*pb.SellersResponse, error)
	SeedSellers(ctx context.Context, in *pb.SeedSellersRequest, opts ...grpc.CallOption) (*pb.SeedSellerResponse, error)
}

type SellerService interface {
	CreateSeller(ctx context.Context, seller *users.Seller) (*users.SellerResponse, error)
	UpdateSeller(ctx context.Context, id string, req *users.Seller) (*users.SellerResponse, error)
	GetSellerById(ctx context.Context, id string) (*users.SellerResponse, error)
	GetSellerByUsername(ctx context.Context, username string) (*users.SellerResponse, error)
	GetRandomSellers(ctx context.Context, count int32) (*users.SellersResponse, error)
	SeedSellers(ctx context.Context, count int32) (string, error)
}
