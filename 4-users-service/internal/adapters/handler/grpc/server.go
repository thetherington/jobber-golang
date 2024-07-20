package grpc

import (
	"fmt"
	"log/slog"
	"net"

	"go.elastic.co/apm/module/apmgrpc/v2"
	"google.golang.org/grpc"

	"github.com/thetherington/jobber-common/middleware"
	"github.com/thetherington/jobber-users/internal/adapters/config"
	"github.com/thetherington/jobber-users/internal/core/port"

	pb "github.com/thetherington/jobber-common/protogen/go/users"
)

type GrpcAdapter struct {
	buyerService  port.BuyerService
	sellerService port.SellerService
	grpcPort      int
	server        *grpc.Server

	// the gRPC server signatures are imposed
	pb.BuyerServiceServer
	pb.SellerServiceServer
}

func NewGrpcAdapter(buyerService port.BuyerService, sellerService port.SellerService, grpcPort int) *GrpcAdapter {
	return &GrpcAdapter{
		buyerService:  buyerService,
		sellerService: sellerService,
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
			middleware.CreateTokenReceiverInterceptor(config.Config.Tokens.GATEWAY, "users"), // validate gateway token
			middleware.CreateMetadataReceiverInterceptor(),                                   // extract gRPC headers insert into context
			middleware.CreatePanicRecoveryInterceptor(),                                      // panic recovery
			apmgrpc.NewUnaryServerInterceptor(),                                              // apm tracing
		),
	)
	a.server = grpcServer

	pb.RegisterBuyerServiceServer(grpcServer, a)
	pb.RegisterSellerServiceServer(grpcServer, a)

	if err = grpcServer.Serve(listen); err != nil {
		slog.With("error", err).Error("Failed to serve gRPC Server", "port", a.grpcPort)
	}
}

func (a *GrpcAdapter) Stop() {
	a.server.Stop()
}
