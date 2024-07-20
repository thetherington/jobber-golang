package grpc

import (
	"context"

	"github.com/thetherington/jobber-common/error-handling/grpcerror"
	"github.com/thetherington/jobber-common/models/auth"
	pb "github.com/thetherington/jobber-common/protogen/go/auth"
	"github.com/thetherington/jobber-common/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func createAuthResponse(payload *auth.AuthResponse) *pb.AuthResponse {
	resp := &pb.AuthResponse{
		Message: payload.Message,
		User: &pb.User{
			Id:              &payload.User.Id,
			ProfilePublicId: &payload.User.ProfilePublicId,
			Username:        &payload.User.Username,
			Email:           &payload.User.Email,
			Country:         &payload.User.Country,
			EmailVerified:   &payload.User.EmailVerified,
			ProfilePicture:  &payload.User.ProfilePicture,
			CreatedAt:       utils.ToDateTime(payload.User.CreatedAt),
			UpdatedAt:       utils.ToDateTime(payload.User.UpdatedAt),
		},
	}

	if payload.Token != "" {
		resp.Token = &payload.Token
	}

	return resp
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

func (a *GrpcAdapter) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.AuthResponse, error) {
	// create a service payload
	payload := &auth.SignUpPayload{
		Username:       req.Username,
		Email:          req.Email,
		Password:       req.Password,
		Country:        req.Country,
		ProfilePicture: req.ProfilePicture,
	}

	resp, err := a.authService.SignUp(ctx, payload)
	if err != nil {
		return nil, serviceError(err)
	}

	// send back gRPC payload
	return createAuthResponse(resp), nil
}

func (a *GrpcAdapter) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.AuthResponse, error) {
	// create a service payload
	payload := &auth.SignInPayload{
		Username: req.Username,
		Password: req.Password,
	}

	resp, err := a.authService.SignIn(ctx, payload)
	if err != nil {
		return nil, serviceError(err)
	}

	// send back gRPC payload
	return createAuthResponse(resp), nil
}

func (a *GrpcAdapter) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.AuthResponse, error) {
	// create a service payload
	payload := &auth.VerifyEmail{
		Token: req.Token,
	}

	resp, err := a.authService.VerifyEmail(ctx, payload)
	if err != nil {
		return nil, serviceError(err)
	}

	// send back gRPC payload
	return createAuthResponse(resp), nil
}

func (a *GrpcAdapter) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (*pb.AuthResponse, error) {
	// create a service payload
	payload := &auth.ForgotPassword{
		Email: req.Email,
	}

	resp, err := a.authService.ForgotPassword(ctx, payload)
	if err != nil {
		return nil, serviceError(err)
	}

	// send back gRPC payload
	return &pb.AuthResponse{
		Message: resp.Message,
	}, nil
}

func (a *GrpcAdapter) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.AuthResponse, error) {
	// create a service payload
	payload := &auth.ResetPassword{
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
		Token:           req.Token,
	}

	resp, err := a.authService.ResetPassword(ctx, payload)
	if err != nil {
		return nil, serviceError(err)
	}

	// send back gRPC payload
	return &pb.AuthResponse{
		Message: resp.Message,
	}, nil
}

func (a *GrpcAdapter) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.AuthResponse, error) {
	// create a service payload
	payload := &auth.ChangePassword{
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	}

	resp, err := a.authService.ChangePassword(ctx, payload)
	if err != nil {
		return nil, serviceError(err)
	}

	// send back gRPC payload
	return &pb.AuthResponse{
		Message: resp.Message,
	}, nil
}

func (a *GrpcAdapter) Seed(ctx context.Context, req *pb.SeedRequest) (*pb.Response, error) {
	resp, err := a.authService.Seed(ctx, int(req.GetCount()))
	if err != nil {
		return nil, serviceError(err)
	}

	return &pb.Response{Payload: resp}, nil
}
