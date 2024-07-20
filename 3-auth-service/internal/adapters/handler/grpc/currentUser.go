package grpc

import (
	"context"

	"github.com/thetherington/jobber-common/models/auth"
	pb "github.com/thetherington/jobber-common/protogen/go/auth"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (a *GrpcAdapter) CurrentUser(ctx context.Context, req *emptypb.Empty) (*pb.AuthResponse, error) {
	resp, err := a.authService.CurrentUser(ctx)
	if err != nil {
		return nil, serviceError(err)
	}

	// send back gRPC payload
	return createAuthResponse(resp), nil
}

func (a *GrpcAdapter) RefreshRoken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.AuthResponse, error) {
	resp, err := a.authService.RefreshToken(ctx, &auth.RefreshToken{Username: req.GetUsername()})
	if err != nil {
		return nil, serviceError(err)
	}

	// send back gRPC payload
	return createAuthResponse(resp), nil
}

func (a *GrpcAdapter) ResendEmail(ctx context.Context, req *pb.ResendEmailRequest) (*pb.AuthResponse, error) {
	resp, err := a.authService.ResendEmail(ctx, &auth.ResendEmail{
		Id:    req.GetUserId(),
		Email: req.GetEmail(),
	})
	if err != nil {
		return nil, serviceError(err)
	}

	// send back gRPC payload
	return createAuthResponse(resp), nil
}
