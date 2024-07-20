package grpc

import (
	"context"

	"github.com/thetherington/jobber-common/models/users"
	pb "github.com/thetherington/jobber-common/protogen/go/users"
)

func (s *GrpcAdapter) CreateSeller(ctx context.Context, req *pb.CreateUpdateSellerPayload) (*pb.SellerResponse, error) {
	resp, err := s.sellerService.CreateSeller(ctx, users.CreateSellerFromReqPayload(req))
	if err != nil {
		return nil, serviceError(err)
	}

	return &pb.SellerResponse{Message: resp.Message, Seller: users.CreatePayloadFromSeller(resp.Seller)}, nil
}

func (s *GrpcAdapter) UpdateSeller(ctx context.Context, req *pb.UpdateSellerRequest) (*pb.SellerResponse, error) {
	resp, err := s.sellerService.UpdateSeller(ctx, req.Id, users.CreateSellerFromReqPayload(req.Seller))
	if err != nil {
		return nil, serviceError(err)
	}

	return &pb.SellerResponse{Message: resp.Message, Seller: users.CreatePayloadFromSeller(resp.Seller)}, nil
}

func (s *GrpcAdapter) GetSellerById(ctx context.Context, req *pb.GetSellerByIdRequest) (*pb.SellerResponse, error) {
	resp, err := s.sellerService.GetSellerById(ctx, req.Id)
	if err != nil {
		return nil, serviceError(err)
	}

	return &pb.SellerResponse{Message: resp.Message, Seller: users.CreatePayloadFromSeller(resp.Seller)}, nil
}

func (s *GrpcAdapter) GetSellerByUsername(ctx context.Context, req *pb.GetSellerByUsernameRequest) (*pb.SellerResponse, error) {
	resp, err := s.sellerService.GetSellerByUsername(ctx, req.Username)
	if err != nil {
		return nil, serviceError(err)
	}

	return &pb.SellerResponse{Message: resp.Message, Seller: users.CreatePayloadFromSeller(resp.Seller)}, nil
}

func (s *GrpcAdapter) GetRandomSellers(ctx context.Context, req *pb.RandomSellersRequest) (*pb.SellersResponse, error) {
	resp, err := s.sellerService.GetRandomSellers(ctx, req.Size)
	if err != nil {
		return nil, serviceError(err)
	}

	sellers := make([]*pb.SellerPayload, 0)

	for _, s := range resp.Sellers {
		sellers = append(sellers, users.CreatePayloadFromSeller(s))
	}

	return &pb.SellersResponse{Message: resp.Message, Sellers: sellers}, nil
}

func (s *GrpcAdapter) SeedSellers(ctx context.Context, req *pb.SeedSellersRequest) (*pb.SeedSellerResponse, error) {
	resp, err := s.sellerService.SeedSellers(ctx, req.Seed)
	if err != nil {
		return nil, serviceError(err)
	}

	return &pb.SeedSellerResponse{Message: resp}, nil
}
