package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/thetherington/jobber-common/error-handling/svc"
	"github.com/thetherington/jobber-common/models/auth"
	pb "github.com/thetherington/jobber-common/protogen/go/auth"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (as *AuthService) CurrentUser(ctx context.Context) (*auth.AuthResponse, error) {
	// gRPC request
	resp, err := as.client.CurrentUser(ctx, &emptypb.Empty{})
	if err != nil {
		slog.With("error", err).Debug("current user error")
		return nil, svc.GrpcErrorResolve(err, "currentuser")
	}

	// return back to http handler with service response
	return createAuthResponse(resp), nil
}

func (as *AuthService) RefreshToken(ctx context.Context, req *auth.RefreshToken) (*auth.AuthResponse, error) {
	// gRPC request
	resp, err := as.client.RefreshRoken(ctx, &pb.RefreshTokenRequest{Username: req.Username})
	if err != nil {
		slog.With("error", err).Debug("refresh token error")
		return nil, svc.GrpcErrorResolve(err, "refreshtoken")
	}

	// return back to http handler with service response
	return createAuthResponse(resp), nil
}

func (as *AuthService) GetLoggedInUsers(ctx context.Context) (*auth.AuthResponse, error) {
	users, err := as.cache.GetLoggedInUsersFromCache(ctx)
	if err != nil {
		slog.With("error", err).Error("failed to get logged in users from cache")
		return nil, err
	}

	// TODO send users to socket io

	return &auth.AuthResponse{Message: fmt.Sprintf("Users online: [%s]", strings.Join(users, ", "))}, nil
}

func (as *AuthService) ResendEmail(ctx context.Context, req *auth.ResendEmail) (*auth.AuthResponse, error) {
	// gRPC request
	resp, err := as.client.ResendEmail(ctx, &pb.ResendEmailRequest{
		UserId: req.Id,
		Email:  req.Email,
	})
	if err != nil {
		slog.With("error", err).Debug("resend email error")
		return nil, svc.GrpcErrorResolve(err, "resendemail")
	}

	// return back to http handler with service response
	return createAuthResponse(resp), nil
}

func (as *AuthService) RemoveLoggedInUser(ctx context.Context, username string) error {
	users, err := as.cache.RemoveLoggedInUserFromCache(ctx, username)
	if err != nil {
		return err
	}

	// TODO send users to socket io
	return as.socket.PushLoggedInUsers(users)
}
