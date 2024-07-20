package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/thetherington/jobber-chat/internal/adapters/config"
	"github.com/thetherington/jobber-chat/internal/adapters/handler/grpc"
	"github.com/thetherington/jobber-chat/internal/adapters/handler/http"
	"github.com/thetherington/jobber-chat/internal/adapters/handler/rabbitmq"
	"github.com/thetherington/jobber-chat/internal/adapters/storage/mongodb"
	"github.com/thetherington/jobber-chat/internal/core/service"
	"github.com/thetherington/jobber-common/cloudinary"
	"github.com/thetherington/jobber-common/logger"
)

const (
	WEB_PORT  = 5005
	GRPC_PORT = 4005
	APP_NAME  = "chat-service"
)

func main() {
	// Load environment variables
	config := config.Config

	// Set logger
	logger.Set(config.App.Env, APP_NAME, config.App.LogLevel)

	// Connect to MongoDB
	slog.Info("Connecting to MongoDB", "address", config.DB.URI)

	client, err := mongodb.New(*config.DB)
	if err != nil {
		slog.With("error", err).Error("Failed to connect to Mongodb", "uri", config.DB.URI)
		os.Exit(1)
	}
	defer mongodb.MustClose(client)

	// Create the RabbitMQ connection
	producer := rabbitmq.NewRabbitMQAdapter(config.RabbitMQ.Endpoint)
	defer producer.Close()

	// Create the Cloudinary image upload dependency
	image, err := cloudinary.New(config.Cloudinary.Name, config.Cloudinary.ApiKey, config.Cloudinary.ApiSecret)
	if err != nil {
		slog.With("error", err).Error("failed to setup cloudinary credentials")
		os.Exit(1)
	}

	chat := service.NewChatService(client.Database(config.DB.Name), producer, image)

	// Create the gRPC Server with the gig service
	grpc := grpc.NewGrpcAdapter(chat, GRPC_PORT)

	// setup the http server to listen for Ping health checks
	router := http.NewRouter()

	go func() {
		addr := fmt.Sprintf(":%d", WEB_PORT)

		slog.Info("Starting HTTP Server", "address", addr)
		if err := http.Serve(addr, router); err != nil {
			slog.With("error", err).Error("Failed to start HTTP Server")
		}
	}()

	// startup the gRPC server
	slog.Info("Starting gRPC Server", "port", GRPC_PORT)
	grpc.Run()
}
