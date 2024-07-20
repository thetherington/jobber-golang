package grpc

import (
	"log/slog"

	"github.com/alexedwards/scs/v2"
	pb "github.com/thetherington/jobber-common/protogen/go/gig"
	"go.elastic.co/apm/module/apmgrpc/v2"
	"google.golang.org/grpc"
)

type GigServiceClient struct {
	pb.GigServiceClient
	pb.GigSearchClient
}

func NewGigAdapter(address string, session *scs.SessionManager) *GigServiceClient {
	ctxForwardInterceptor := CreateContextForwarderInterceptor(session) // create context interceptor forwarder
	tokenForwardInterceptor := CreateTokenForwarderInterceptor("gig")   // create gateway token and gRPC header

	conn, err := NewClient(address, grpc.WithChainUnaryInterceptor(ctxForwardInterceptor, tokenForwardInterceptor, apmgrpc.NewUnaryClientInterceptor()))
	if err != nil {
		slog.With("error", err).Error("Failed to create Gig gRPC Client")
	}

	gig := pb.NewGigServiceClient(conn)   // gig crud rpc calls
	search := pb.NewGigSearchClient(conn) // gig search rpc calls

	return &GigServiceClient{
		gig,
		search,
	}
}
