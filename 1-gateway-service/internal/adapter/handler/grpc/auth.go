package grpc

import (
	"log/slog"

	"github.com/alexedwards/scs/v2"
	pb "github.com/thetherington/jobber-common/protogen/go/auth"
	"go.elastic.co/apm/module/apmgrpc/v2"
	"google.golang.org/grpc"
)

type AuthServiceClient struct {
	pb.AuthServiceClient
	pb.CurrentUserServiceClient
	pb.SearchServiceClient
}

func NewAuthAdapter(address string, session *scs.SessionManager) *AuthServiceClient {
	ctxForwardInterceptor := CreateContextForwarderInterceptor(session) // create context interceptor forwarder
	tokenForwardInterceptor := CreateTokenForwarderInterceptor("auth")  // create gateway token and gRPC header

	conn, err := NewClient(address, grpc.WithChainUnaryInterceptor(ctxForwardInterceptor, tokenForwardInterceptor, apmgrpc.NewUnaryClientInterceptor()))
	if err != nil {
		slog.With("error", err).Error("Failed to create Auth gRPC Client")
	}

	c1 := pb.NewAuthServiceClient(conn)        // user signup/login/password reset rpc calls
	c2 := pb.NewCurrentUserServiceClient(conn) // current user, token refresh, email resend calls
	c3 := pb.NewSearchServiceClient(conn)      // elastic search GetGigById SearchGigs

	return &AuthServiceClient{c1, c2, c3}
}
