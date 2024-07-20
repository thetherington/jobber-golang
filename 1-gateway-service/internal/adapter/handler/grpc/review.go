package grpc

import (
	"log/slog"

	"github.com/alexedwards/scs/v2"
	pb "github.com/thetherington/jobber-common/protogen/go/review"
	"go.elastic.co/apm/module/apmgrpc/v2"
	"google.golang.org/grpc"
)

type ReviewServiceClient struct {
	pb.ReviewServiceClient
}

func NewReviewAdapter(address string, session *scs.SessionManager) *ReviewServiceClient {
	ctxForwardInterceptor := CreateContextForwarderInterceptor(session)  // create context interceptor forwarder
	tokenForwardInterceptor := CreateTokenForwarderInterceptor("review") // create gateway token and gRPC header

	conn, err := NewClient(address, grpc.WithChainUnaryInterceptor(ctxForwardInterceptor, tokenForwardInterceptor, apmgrpc.NewUnaryClientInterceptor()))
	if err != nil {
		slog.With("error", err).Error("Failed to create Users gRPC Client")
	}

	review := pb.NewReviewServiceClient(conn)

	return &ReviewServiceClient{review}
}
