package service

import (
	"context"
	"log/slog"

	"github.com/thetherington/jobber-common/error-handling/svc"
	"github.com/thetherington/jobber-common/models/users"
	pb "github.com/thetherington/jobber-common/protogen/go/users"
	"github.com/thetherington/jobber-gateway/internal/core/port"
)

/**
 * SellerService implements
 */
type SellerService struct {
	client port.SellerRPCClient
}

// NewBuyerService creates a new buyer service instance
func NewSellerService(rpc port.SellerRPCClient) *SellerService {
	return &SellerService{
		rpc,
	}
}

func (s *SellerService) CreateSeller(ctx context.Context, seller *users.Seller) (*users.SellerResponse, error) {
	resp, err := s.client.CreateSeller(ctx, users.CreateReqPayload(seller))
	if err != nil {
		slog.With("error", err).Debug("CreateSeller error")
		return nil, svc.GrpcErrorResolve(err, "CreateSeller")
	}

	return &users.SellerResponse{
		Message: resp.Message,
		Seller:  users.CreateSellerFromPayload(resp.Seller),
	}, nil
}

func (s *SellerService) UpdateSeller(ctx context.Context, id string, seller *users.Seller) (*users.SellerResponse, error) {
	resp, err := s.client.UpdateSeller(ctx, &pb.UpdateSellerRequest{
		Id:     id,
		Seller: users.CreateReqPayload(seller),
	})
	if err != nil {
		slog.With("error", err).Debug("UpdateSeller error")
		return nil, svc.GrpcErrorResolve(err, "UpdateSeller")
	}

	return &users.SellerResponse{
		Message: resp.Message,
		Seller:  users.CreateSellerFromPayload(resp.Seller),
	}, nil
}

func (s *SellerService) GetSellerById(ctx context.Context, id string) (*users.SellerResponse, error) {
	resp, err := s.client.GetSellerById(ctx, &pb.GetSellerByIdRequest{Id: id})
	if err != nil {
		slog.With("error", err).Debug("GetSellerById error")
		return nil, svc.GrpcErrorResolve(err, "GetSellerById")
	}

	return &users.SellerResponse{
		Message: resp.Message,
		Seller:  users.CreateSellerFromPayload(resp.Seller),
	}, nil
}

func (s *SellerService) GetSellerByUsername(ctx context.Context, username string) (*users.SellerResponse, error) {
	resp, err := s.client.GetSellerByUsername(ctx, &pb.GetSellerByUsernameRequest{Username: username})
	if err != nil {
		slog.With("error", err).Debug("GetSellerByUsername error")
		return nil, svc.GrpcErrorResolve(err, "GetSellerByUsername")
	}

	return &users.SellerResponse{
		Message: resp.Message,
		Seller:  users.CreateSellerFromPayload(resp.Seller),
	}, nil
}

func (s *SellerService) GetRandomSellers(ctx context.Context, count int32) (*users.SellersResponse, error) {
	resp, err := s.client.GetRandomSellers(ctx, &pb.RandomSellersRequest{Size: count})
	if err != nil {
		slog.With("error", err).Debug("GetRandomSellers error")
		return nil, svc.GrpcErrorResolve(err, "GetRandomSellers")
	}

	sellers := make([]*users.Seller, 0)

	for _, s := range resp.Sellers {
		sellers = append(sellers, users.CreateSellerFromPayload(s))
	}

	return &users.SellersResponse{
		Message: resp.Message,
		Sellers: sellers,
	}, nil
}

func (s *SellerService) SeedSellers(ctx context.Context, count int32) (string, error) {
	resp, err := s.client.SeedSellers(ctx, &pb.SeedSellersRequest{Seed: count})
	if err != nil {
		slog.With("error", err).Debug("SeedSellers error")
		return "", svc.GrpcErrorResolve(err, "SeedSellers")
	}

	return resp.Message, nil
}
