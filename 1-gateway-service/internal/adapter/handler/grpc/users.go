package grpc

import (
	"log/slog"

	"github.com/alexedwards/scs/v2"
	pb "github.com/thetherington/jobber-common/protogen/go/users"
	"go.elastic.co/apm/module/apmgrpc/v2"
	"google.golang.org/grpc"
)

type UsersServiceClient struct {
	pb.BuyerServiceClient
	pb.SellerServiceClient
}

func NewUsersAdapter(address string, session *scs.SessionManager) *UsersServiceClient {
	ctxForwardInterceptor := CreateContextForwarderInterceptor(session) // create context interceptor forwarder
	tokenForwardInterceptor := CreateTokenForwarderInterceptor("users") // create gateway token and gRPC header

	conn, err := NewClient(address, grpc.WithChainUnaryInterceptor(ctxForwardInterceptor, tokenForwardInterceptor, apmgrpc.NewUnaryClientInterceptor()))
	if err != nil {
		slog.With("error", err).Error("Failed to create Users gRPC Client")
	}

	buyer := pb.NewBuyerServiceClient(conn)   // user buyer rpc calls
	seller := pb.NewSellerServiceClient(conn) // user seller rpc calls

	return &UsersServiceClient{buyer, seller}
}
