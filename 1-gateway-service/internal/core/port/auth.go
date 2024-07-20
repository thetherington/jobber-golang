package port

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/thetherington/jobber-common/models/auth"
	pb "github.com/thetherington/jobber-common/protogen/go/auth"
)

type AuthService interface {
	SignUp(ctx context.Context, req *auth.SignUpPayload) (*auth.AuthResponse, error)
	SignIn(ctx context.Context, req *auth.SignInPayload) (*auth.AuthResponse, error)
	SignOut(ctx context.Context)
	VerifyEmail(ctx context.Context, req *auth.VerifyEmail) (*auth.AuthResponse, error)
	ForgotPassword(ctx context.Context, req *auth.ForgotPassword) (*auth.AuthResponse, error)
	ResetPassword(ctx context.Context, req *auth.ResetPassword) (*auth.AuthResponse, error)
	ChangePassword(ctx context.Context, req *auth.ChangePassword) (*auth.AuthResponse, error)
	Seed(ctx context.Context, count int) (*auth.AuthResponse, error)
	CurrentUser(ctx context.Context) (*auth.AuthResponse, error)
	RefreshToken(ctx context.Context, req *auth.RefreshToken) (*auth.AuthResponse, error)
	ResendEmail(ctx context.Context, req *auth.ResendEmail) (*auth.AuthResponse, error)
	GetLoggedInUsers(ctx context.Context) (*auth.AuthResponse, error)
	RemoveLoggedInUser(ctx context.Context, username string) error
}

type AuthRPCClient interface {
	SignUp(ctx context.Context, in *pb.SignUpRequest, opts ...grpc.CallOption) (*pb.AuthResponse, error)
	SignIn(ctx context.Context, in *pb.SignInRequest, opts ...grpc.CallOption) (*pb.AuthResponse, error)
	VerifyEmail(ctx context.Context, in *pb.VerifyEmailRequest, opts ...grpc.CallOption) (*pb.AuthResponse, error)
	ForgotPassword(ctx context.Context, in *pb.ForgotPasswordRequest, opts ...grpc.CallOption) (*pb.AuthResponse, error)
	ResetPassword(ctx context.Context, in *pb.ResetPasswordRequest, opts ...grpc.CallOption) (*pb.AuthResponse, error)
	ChangePassword(ctx context.Context, in *pb.ChangePasswordRequest, opts ...grpc.CallOption) (*pb.AuthResponse, error)
	Seed(ctx context.Context, in *pb.SeedRequest, opts ...grpc.CallOption) (*pb.Response, error)
	CurrentUser(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*pb.AuthResponse, error)
	RefreshRoken(ctx context.Context, in *pb.RefreshTokenRequest, opts ...grpc.CallOption) (*pb.AuthResponse, error)
	ResendEmail(ctx context.Context, in *pb.ResendEmailRequest, opts ...grpc.CallOption) (*pb.AuthResponse, error)
}

type AuthSocketManager interface {
	PushLoggedInUsers(users []string) error
}
