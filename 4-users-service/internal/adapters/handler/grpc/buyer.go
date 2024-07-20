package grpc

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/thetherington/jobber-common/error-handling/grpcerror"
	"github.com/thetherington/jobber-common/error-handling/svc"
	"github.com/thetherington/jobber-common/middleware"
	"github.com/thetherington/jobber-common/models/users"
	pb "github.com/thetherington/jobber-common/protogen/go/users"
	"github.com/thetherington/jobber-common/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func createBuyerPayload(buyer *users.Buyer) *pb.BuyerPayload {
	pbBuyer := &pb.BuyerPayload{
		BuyerId:        &buyer.Id,
		Username:       &buyer.Username,
		Email:          &buyer.Email,
		ProfilePicture: &buyer.ProfilePicture,
		Country:        &buyer.Country,
		IsSeller:       &buyer.IsSeller,
		PurchasedGigs:  make([]string, 0),
		CreatedAt:      utils.ToDateTime(buyer.CreatedAt),
		UpdatedAt:      utils.ToDateTime(buyer.UpdatedAt),
	}

	pbBuyer.PurchasedGigs = append(pbBuyer.PurchasedGigs, buyer.PurchasedGigs...)

	return pbBuyer
}

// parses service error and returns a equivelient gRPC error
func serviceError(err error) error {
	// try to cast the error to a grpcerror lookup
	if apiError, ok := grpcerror.FromError(err); ok {
		s := status.New(apiError.Status, apiError.Message)
		return s.Err()
	}

	// generic response
	s := status.New(codes.Internal, err.Error())
	return s.Err()
}

func (b *GrpcAdapter) GetBuyerByEmail(ctx context.Context, req *emptypb.Empty) (*pb.BuyerPayload, error) {
	// Get email from the user cookie session passed down into the context.
	email := ctx.Value(middleware.CtxEmailKey)
	if email == nil {
		slog.Debug("Email in context is nil")
		return nil, serviceError(svc.NewError(svc.ErrUnAuthorized, fmt.Errorf("session invalid. please login and try again")))
	}

	buyer, err := b.buyerService.GetBuyerByEmail(ctx, email.(string))
	if err != nil {
		return nil, serviceError(err)
	}

	return createBuyerPayload(buyer), nil
}

func (b *GrpcAdapter) GetBuyerByUsername(ctx context.Context, req *emptypb.Empty) (*pb.BuyerPayload, error) {
	// Get username from the user cookie session passed down into the context.
	username := ctx.Value(middleware.CtxUsernameKey)
	if username == nil {
		slog.Debug("Username in context is nil")
		return nil, serviceError(svc.NewError(svc.ErrUnAuthorized, fmt.Errorf("session invalid. please login and try again")))
	}

	buyer, err := b.buyerService.GetBuyerByUsername(ctx, username.(string))
	if err != nil {
		return nil, serviceError(err)
	}

	return createBuyerPayload(buyer), nil
}

func (b *GrpcAdapter) GetBuyerByProvidedUsername(ctx context.Context, req *pb.BuyerUsernameRequest) (*pb.BuyerPayload, error) {
	buyer, err := b.buyerService.GetBuyerByUsername(ctx, req.GetUsername())
	if err != nil {
		return nil, serviceError(err)
	}

	return createBuyerPayload(buyer), nil
}
