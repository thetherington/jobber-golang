package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/thetherington/jobber-common/logger"
	"github.com/thetherington/jobber-review/internal/adapters/config"
	"github.com/thetherington/jobber-review/internal/adapters/handler/grpc"
	"github.com/thetherington/jobber-review/internal/adapters/handler/http"
	"github.com/thetherington/jobber-review/internal/adapters/handler/rabbitmq"
	"github.com/thetherington/jobber-review/internal/adapters/storage/postgres"
	"github.com/thetherington/jobber-review/internal/adapters/storage/postgres/repository"
	"github.com/thetherington/jobber-review/internal/core/service"
)

const (
	WEB_PORT  = 5007
	GRPC_PORT = 4007
	APP_NAME  = "review-service"
)

func main() {
	// Load environment variables
	config := config.Config

	// Set logger
	logger.Set(config.App.Env, APP_NAME, config.App.LogLevel)

	// Init database
	ctx := context.Background()
	db, err := postgres.New(ctx, config.DB.URI+"?sslmode=disable")
	if err != nil {
		slog.With("error", err).Error("Error initializing database connection")
		os.Exit(1)
	}
	defer db.Close()

	slog.With("address", config.DB.URI).Info("Connected to the database")

	// Migrate database
	if err := db.Migrate(); err != nil {
		slog.With("error", err).Error("Error migrating database")
		os.Exit(1)
	}

	slog.Info("Successfully migrated the database")

	// Create the RabbitMQ connection
	producer := rabbitmq.NewRabbitMQAdapter(config.RabbitMQ.Endpoint)
	defer producer.Close()

	// initialize the review db respository with the review service
	repo := repository.NewReviewRepository(db)
	review := service.NewReviewService(repo, producer)

	// Create the gRPC Server with the gig service
	grpc := grpc.NewGrpcAdapter(review, GRPC_PORT)

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
