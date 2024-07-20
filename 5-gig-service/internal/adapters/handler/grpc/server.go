package grpc

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/thetherington/jobber-gig/internal/adapters/config"
	"github.com/thetherington/jobber-gig/internal/core/port"
	"go.elastic.co/apm/module/apmgrpc/v2"
	"google.golang.org/grpc"

	"github.com/thetherington/jobber-common/middleware"
	pb "github.com/thetherington/jobber-common/protogen/go/gig"
)

type GrpcAdapter struct {
	gigService port.GigService
	grpcPort   int
	server     *grpc.Server

	// the gRPC server signatures are imposed
	pb.GigServiceServer
	pb.GigSearchServer
}

func NewGrpcAdapter(gigService port.GigService, grpcPort int) *GrpcAdapter {
	return &GrpcAdapter{
		gigService: gigService,
		grpcPort:   grpcPort,
	}
}

func (a *GrpcAdapter) Run() {
	var err error

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.grpcPort))
	if err != nil {
		slog.With("error", err).Error("gRPC Failed to listen", "port", a.grpcPort)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.CreateTokenReceiverInterceptor(config.Config.Tokens.GATEWAY, "gig"), // validate gateway token
			middleware.CreateMetadataReceiverInterceptor(),                                 // extract gRPC headers insert into context
			middleware.CreatePanicRecoveryInterceptor(),                                    // panic recovery
			apmgrpc.NewUnaryServerInterceptor(),                                            // apm tracing
		),
	)
	a.server = grpcServer

	pb.RegisterGigServiceServer(grpcServer, a)
	pb.RegisterGigSearchServer(grpcServer, a)

	if err = grpcServer.Serve(listen); err != nil {
		slog.With("error", err).Error("Failed to serve gRPC Server", "port", a.grpcPort)
	}
}

func (a *GrpcAdapter) Stop() {
	a.server.Stop()
}
