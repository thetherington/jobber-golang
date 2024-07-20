package service

import (
	"context"
	"log/slog"

	"github.com/thetherington/jobber-common/error-handling/svc"
	"github.com/thetherington/jobber-common/models/auth"
	"github.com/thetherington/jobber-common/utils"
	"github.com/thetherington/jobber-gateway/internal/core/port"

	pb "github.com/thetherington/jobber-common/protogen/go/auth"
)

/**
 * AuthService implements
 */
type AuthService struct {
	client port.AuthRPCClient
	cache  port.CacheRepository
	socket port.AuthSocketManager
}

// NewAuthService creates a new auth service instance
func NewAuthService(rpc port.AuthRPCClient, cache port.CacheRepository, socket port.AuthSocketManager) *AuthService {
	return &AuthService{
		rpc,
		cache,
		socket,
	}
}

func createAuthResponse(payload *pb.AuthResponse) *auth.AuthResponse {
	resp := &auth.AuthResponse{
		Message: payload.Message,
		User: &auth.AuthDocument{
			Id:              payload.User.GetId(),
			ProfilePublicId: payload.User.GetProfilePublicId(),
			Username:        payload.User.GetUsername(),
			Email:           payload.User.GetEmail(),
			Country:         payload.User.GetCountry(),
			EmailVerified:   payload.User.GetEmailVerified(),
			ProfilePicture:  payload.User.GetProfilePicture(),
			CreatedAt:       utils.ToTime(payload.User.GetCreatedAt()),
			UpdatedAt:       utils.ToTime(payload.User.GetUpdatedAt()),
		},
	}

	if payload.GetToken() != "" {
		resp.Token = payload.GetToken()
	}

	return resp
}

func (as *AuthService) SignUp(ctx context.Context, req *auth.SignUpPayload) (*auth.AuthResponse, error) {
	// create a gRPC request payload
	protoRequest := &pb.SignUpRequest{
		Username:       req.Username,
		Password:       req.Password,
		Email:          req.Email,
		Country:        req.Country,
		ProfilePicture: req.ProfilePicture,
	}

	// gRPC request
	resp, err := as.client.SignUp(ctx, protoRequest)
	if err != nil {
		slog.With("error", err).Debug("authentication signup error")
		return nil, svc.GrpcErrorResolve(err, "signup")
	}

	// temp
	as.cache.SaveLoggedInUserToCache(ctx, *resp.User.Username)

	// return back to http handler with service response
	return createAuthResponse(resp), nil
}

func (as *AuthService) SignIn(ctx context.Context, req *auth.SignInPayload) (*auth.AuthResponse, error) {
	// create a gRPC request payload
	protoRequest := &pb.SignInRequest{
		Username: req.Username,
		Password: req.Password,
	}

	// gRPC request
	resp, err := as.client.SignIn(ctx, protoRequest)
	if err != nil {
		slog.With("error", err).Debug("authentication signin error")
		return nil, svc.GrpcErrorResolve(err, "signup")
	}

	// temp
	as.cache.SaveLoggedInUserToCache(ctx, *resp.User.Username)

	// return back to http handler with service response
	return createAuthResponse(resp), nil
}

func (as *AuthService) SignOut(ctx context.Context) {
}

func (as *AuthService) VerifyEmail(ctx context.Context, req *auth.VerifyEmail) (*auth.AuthResponse, error) {
	// gRPC request
	resp, err := as.client.VerifyEmail(ctx, &pb.VerifyEmailRequest{
		Token: req.Token,
	})
	if err != nil {
		slog.With("error", err).Debug("verify email error")
		return nil, svc.GrpcErrorResolve(err, "verifyEmail")
	}

	// return back to http handler with service response
	return createAuthResponse(resp), nil
}

func (as *AuthService) ForgotPassword(ctx context.Context, req *auth.ForgotPassword) (*auth.AuthResponse, error) {
	// gRPC request
	resp, err := as.client.ForgotPassword(ctx, &pb.ForgotPasswordRequest{
		Email: req.Email,
	})
	if err != nil {
		slog.With("error", err).Debug("forgot password error")
		return nil, svc.GrpcErrorResolve(err, "forgotPassword")
	}

	// return back to http handler with service response
	return &auth.AuthResponse{Message: resp.GetMessage()}, nil
}

func (as *AuthService) ResetPassword(ctx context.Context, req *auth.ResetPassword) (*auth.AuthResponse, error) {
	// gRPC request
	resp, err := as.client.ResetPassword(ctx, &pb.ResetPasswordRequest{
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
		Token:           req.Token,
	})
	if err != nil {
		slog.With("error", err).Debug("reset password error")
		return nil, svc.GrpcErrorResolve(err, "resetPassword")
	}

	// return back to http handler with service response
	return &auth.AuthResponse{Message: resp.GetMessage()}, nil
}

func (as *AuthService) ChangePassword(ctx context.Context, req *auth.ChangePassword) (*auth.AuthResponse, error) {
	// gRPC request
	resp, err := as.client.ChangePassword(ctx, &pb.ChangePasswordRequest{
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	})
	if err != nil {
		slog.With("error", err).Debug("change password error")
		return nil, svc.GrpcErrorResolve(err, "changePassword")
	}

	// return back to http handler with service response
	return &auth.AuthResponse{Message: resp.GetMessage()}, nil
}

func (as *AuthService) Seed(ctx context.Context, count int) (*auth.AuthResponse, error) {
	resp, err := as.client.Seed(ctx, &pb.SeedRequest{Count: int32(count)})
	if err != nil {
		slog.With("error", err).Debug("change password error")
		return nil, svc.GrpcErrorResolve(err, "seed")
	}

	return &auth.AuthResponse{Message: resp.GetPayload()}, nil
}
