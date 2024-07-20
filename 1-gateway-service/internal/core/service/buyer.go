package service

import (
	"context"
	"log/slog"

	"github.com/thetherington/jobber-common/error-handling/svc"
	"github.com/thetherington/jobber-common/models/users"
	pb "github.com/thetherington/jobber-common/protogen/go/users"
	"github.com/thetherington/jobber-common/utils"
	"github.com/thetherington/jobber-gateway/internal/core/port"
	"google.golang.org/protobuf/types/known/emptypb"
)

/**
 * BuyerService implements
 */
type BuyerService struct {
	client port.BuyerRPCClient
}

// NewBuyerService creates a new buyer service instance
func NewBuyerService(rpc port.BuyerRPCClient) *BuyerService {
	return &BuyerService{
		rpc,
	}
}

func createModelBuyer(b *pb.BuyerPayload) *users.Buyer {
	buyer := &users.Buyer{
		Id:             *b.BuyerId,
		Username:       *b.Username,
		Email:          *b.Email,
		ProfilePicture: *b.ProfilePicture,
		Country:        *b.Country,
		IsSeller:       *b.IsSeller,
		PurchasedGigs:  make([]string, 0),
		CreatedAt:      utils.ToTime(b.CreatedAt),
		UpdatedAt:      utils.ToTime(b.UpdatedAt),
	}

	buyer.PurchasedGigs = append(buyer.PurchasedGigs, b.PurchasedGigs...)

	return buyer
}

func (bs *BuyerService) GetBuyerByEmail(ctx context.Context) (*users.BuyerResponse, error) {
	resp, err := bs.client.GetBuyerByEmail(ctx, &emptypb.Empty{})
	if err != nil {
		slog.With("error", err).Debug("GetBuyerByEmail error")
		return nil, svc.GrpcErrorResolve(err, "GetBuyerByEmail")
	}

	return &users.BuyerResponse{Message: "Buyer profile", Buyer: createModelBuyer(resp)}, nil
}

func (bs *BuyerService) GetBuyerByUsername(ctx context.Context) (*users.BuyerResponse, error) {
	resp, err := bs.client.GetBuyerByUsername(ctx, &emptypb.Empty{})
	if err != nil {
		slog.With("error", err).Debug("GetBuyerByUsername error")
		return nil, svc.GrpcErrorResolve(err, "GetBuyerByUsername")
	}

	return &users.BuyerResponse{Message: "Buyer profile", Buyer: createModelBuyer(resp)}, nil
}

func (bs *BuyerService) GetBuyerByProvidedUsername(ctx context.Context, username string) (*users.BuyerResponse, error) {
	resp, err := bs.client.GetBuyerByProvidedUsername(ctx, &pb.BuyerUsernameRequest{Username: username})
	if err != nil {
		slog.With("error", err).Debug("GetBuyerByProvidedUsername error")
		return nil, svc.GrpcErrorResolve(err, "GetBuyerByProvidedUsername")
	}

	return &users.BuyerResponse{Message: "Buyer profile", Buyer: createModelBuyer(resp)}, nil
}
