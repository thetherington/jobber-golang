package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/thetherington/jobber-auth/internal/adapters/config"
	"github.com/thetherington/jobber-auth/internal/adapters/handler/grpc"
	"github.com/thetherington/jobber-auth/internal/adapters/handler/http"
	"github.com/thetherington/jobber-auth/internal/adapters/handler/rabbitmq"
	"github.com/thetherington/jobber-auth/internal/adapters/storage/elasticsearch"
	"github.com/thetherington/jobber-auth/internal/adapters/storage/postgres"
	"github.com/thetherington/jobber-auth/internal/adapters/storage/postgres/repository"
	"github.com/thetherington/jobber-auth/internal/core/service"
	token "github.com/thetherington/jobber-common/client-token"
	"github.com/thetherington/jobber-common/cloudinary"
	"github.com/thetherington/jobber-common/logger"
)

const (
	WEB_PORT  = 5002
	GRPC_PORT = 4002
	APP_NAME  = "auth-service"
	Index     = "gigs"
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

	// Create RabbitMQ connection
	rbbtmq := rabbitmq.NewRabbitMQAdapter(config.RabbitMQ.Endpoint)
	defer rbbtmq.Close()

	// Create TokenMaker Service for user jwt token creation
	tkmkr, err := token.NewClientJWTMaker(config.Tokens.JWT)
	if err != nil {
		slog.With("error", err).Error("failed to create token creator service")
		os.Exit(1)
	}

	// Create the Cloudinary image upload dependency
	image, err := cloudinary.New(config.Cloudinary.Name, config.Cloudinary.ApiKey, config.Cloudinary.ApiSecret)
	if err != nil {
		slog.With("error", err).Error("failed to setup cloudinary credentials")
		os.Exit(1)
	}

	// Create Elasticsearch connection
	esClient, err := elasticsearch.New(config.Elastic, Index)
	if err != nil {
		slog.With("error", err).Error("failed to setup elasticsearch client")
		os.Exit(1)
	}

	// Check Elasticsearch connection and create index in background
	go func() {
		esClient.CheckConnection(context.Background())
		esClient.CreateIndex(context.Background())
	}()

	// Create the dependency injection for db, service and grpc handler
	queries := repository.New(db)

	auth := service.NewAuthService(queries, rbbtmq, tkmkr, image)
	search := service.NewSearchService(esClient)

	grpc := grpc.NewGrpcAdapter(auth, search, GRPC_PORT)

	router := http.NewRouter()

	// startup the http server to listen for health checks
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
