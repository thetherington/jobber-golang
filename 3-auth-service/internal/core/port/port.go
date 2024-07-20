package port

import (
	"context"

	"github.com/thetherington/jobber-auth/internal/adapters/storage/postgres/repository"
	"github.com/thetherington/jobber-common/models/auth"
	"github.com/thetherington/jobber-common/models/gig"
	"github.com/thetherington/jobber-common/models/search"
	pb "github.com/thetherington/jobber-common/protogen/go/auth"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthRepository interface {
	repository.Querier
}

type ImageUploader interface {
	UploadImage(ctx context.Context, file string, public_id string, overwrite bool, invalidate bool) (string, error)
	UploadVideo(ctx context.Context, file string, public_id string, overwrite bool, invalidate bool) (string, error)
}

type AuthService interface {
	SignUp(ctx context.Context, req *auth.SignUpPayload) (*auth.AuthResponse, error)
	SignIn(ctx context.Context, req *auth.SignInPayload) (*auth.AuthResponse, error)
	VerifyEmail(ctx context.Context, req *auth.VerifyEmail) (*auth.AuthResponse, error)
	ForgotPassword(ctx context.Context, req *auth.ForgotPassword) (*auth.AuthResponse, error)
	ResetPassword(ctx context.Context, req *auth.ResetPassword) (*auth.AuthResponse, error)
	ChangePassword(ctx context.Context, req *auth.ChangePassword) (*auth.AuthResponse, error)
	Seed(ctx context.Context, count int) (string, error)
	CurrentUser(ctx context.Context) (*auth.AuthResponse, error)
	RefreshToken(ctx context.Context, req *auth.RefreshToken) (*auth.AuthResponse, error)
	ResendEmail(ctx context.Context, req *auth.ResendEmail) (*auth.AuthResponse, error)
}

type AuthRPCServer interface {
	SignUp(context.Context, *pb.SignUpRequest) (*pb.AuthResponse, error)
	SignIn(context.Context, *pb.SignInRequest) (*pb.AuthResponse, error)
	VerifyEmail(context.Context, *pb.VerifyEmailRequest) (*pb.AuthResponse, error)
	ForgotPassword(context.Context, *pb.ForgotPasswordRequest) (*pb.AuthResponse, error)
	ResetPassword(context.Context, *pb.ResetPasswordRequest) (*pb.AuthResponse, error)
	ChangePassword(context.Context, *pb.ChangePasswordRequest) (*pb.AuthResponse, error)
	Seed(context.Context, *pb.Request) (*pb.Response, error)
	CurrentUser(context.Context, *emptypb.Empty) (*pb.AuthResponse, error)
	RefreshRoken(context.Context, *pb.RefreshTokenRequest) (*pb.AuthResponse, error)
	ResendEmail(context.Context, *pb.ResendEmailRequest) (*pb.AuthResponse, error)
	GetGigById(context.Context, *pb.GetGigRequest) (*pb.GigResponse, error)
	SearchGig(context.Context, *pb.SearchRequest) (*pb.SearchResponse, error)
}

type AuthProducer interface {
	PublishDirectMessage(exchangeName string, routingKey string, data []byte) error
}

type SearchClient interface {
	GetDocumentById(ctx context.Context, id string) (*gig.SellerGig, error)
	SearchGigs(ctx context.Context, searchQuery string, paginate *search.PaginateProps, deliveryTime *string, min *float64, max *float64) (int64, []*gig.SellerGig, error)
}

type SearchService interface {
	GetGigByID(ctx context.Context, id string) (*gig.SellerGig, error)
	SearchGigs(ctx context.Context, req search.SearchRequest) (*search.SearchResponse, error)
}
