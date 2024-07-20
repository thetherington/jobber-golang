package grpc

import (
	"fmt"
	"log/slog"
	"net"
	"sync"

	pb "github.com/thetherington/jobber-common/protogen/go/order"
	"github.com/thetherington/jobber-order/internal/adapters/config"
	"github.com/thetherington/jobber-order/internal/core/port"
	"go.elastic.co/apm/module/apmgrpc/v2"
	"google.golang.org/grpc"

	"github.com/thetherington/jobber-common/middleware"
)

type GrpcAdapter struct {
	orderService        port.OrderService
	notificationService port.NotificationService
	grpcPort            int
	server              *grpc.Server

	// subscribers is a concurrent map that holds mapping from a client ID to it's subscriber
	notifySubscribers sync.Map
	orderSubscribers  sync.Map

	// the gRPC server signatures are imposed
	pb.OrderServiceServer
	pb.NotificationServiceServer
}

func NewGrpcAdapter(orderService port.OrderService, notifyService port.NotificationService, grpcPort int) *GrpcAdapter {
	return &GrpcAdapter{
		orderService:        orderService,
		notificationService: notifyService,
		grpcPort:            grpcPort,
	}
}

func (a *GrpcAdapter) Run() {
	var err error

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.grpcPort))
	if err != nil {
		slog.With("error", err).Error("gRPC Failed to listen", "port", a.grpcPort)
	}

	// TODO add streaming middleware interceptor
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.CreateTokenReceiverInterceptor(config.Config.Tokens.GATEWAY, "order"), // validate gateway token
			middleware.CreateMetadataReceiverInterceptor(),                                   // extract gRPC headers insert into context
			middleware.CreatePanicRecoveryInterceptor(),                                      // panic recovery
			apmgrpc.NewUnaryServerInterceptor(),                                              // apm tracing
		),
		grpc.ChainStreamInterceptor(
			middleware.CreateStreamTokenReceiverInterceptor(config.Config.Tokens.GATEWAY, "order"),
			middleware.CreateStreamPanicRecoveryInterceptor(),
			apmgrpc.NewStreamServerInterceptor(),
		),
		grpc.MaxRecvMsgSize(1024*1024*20),
		grpc.MaxSendMsgSize(1024*1024*20),
	)

	a.server = grpcServer

	pb.RegisterOrderServiceServer(grpcServer, a)
	pb.RegisterNotificationServiceServer(grpcServer, a)

	if err = grpcServer.Serve(listen); err != nil {
		slog.With("error", err).Error("Failed to serve gRPC Server", "port", a.grpcPort)
	}
}

func (a *GrpcAdapter) Stop() {
	a.server.Stop()
}
