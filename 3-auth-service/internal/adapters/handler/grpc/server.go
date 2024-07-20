package grpc

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/thetherington/jobber-auth/internal/adapters/config"
	"github.com/thetherington/jobber-auth/internal/core/port"
	"go.elastic.co/apm/module/apmgrpc/v2"
	"google.golang.org/grpc"

	"github.com/thetherington/jobber-common/middleware"
	pb "github.com/thetherington/jobber-common/protogen/go/auth"
)

type GrpcAdapter struct {
	authService   port.AuthService
	searchService port.SearchService
	grpcPort      int
	server        *grpc.Server

	// the gRPC server signatures are imposed
	pb.AuthServiceServer
	pb.CurrentUserServiceServer
	pb.SearchServiceServer
}

func NewGrpcAdapter(authService port.AuthService, searchService port.SearchService, grpcPort int) *GrpcAdapter {
	return &GrpcAdapter{
		authService:   authService,
		searchService: searchService,
		grpcPort:      grpcPort,
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
			middleware.CreateTokenReceiverInterceptor(config.Config.Tokens.GATEWAY, "auth"), // validate gateway token
			middleware.CreateMetadataReceiverInterceptor(),                                  // extract gRPC headers insert into context
			middleware.CreatePanicRecoveryInterceptor(),                                     // panic recovery
			apmgrpc.NewUnaryServerInterceptor(),                                             // apm tracing
		),
	)
	a.server = grpcServer

	pb.RegisterAuthServiceServer(grpcServer, a)
	pb.RegisterCurrentUserServiceServer(grpcServer, a)
	pb.RegisterSearchServiceServer(grpcServer, a)

	if err = grpcServer.Serve(listen); err != nil {
		slog.With("error", err).Error("Failed to serve gRPC Server", "port", a.grpcPort)
	}
}

func (a *GrpcAdapter) Stop() {
	a.server.Stop()
}
