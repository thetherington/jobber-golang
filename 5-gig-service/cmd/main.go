package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/thetherington/jobber-common/cloudinary"
	"github.com/thetherington/jobber-common/logger"
	"github.com/thetherington/jobber-gig/internal/adapters/config"
	"github.com/thetherington/jobber-gig/internal/adapters/handler/grpc"
	"github.com/thetherington/jobber-gig/internal/adapters/handler/http"
	"github.com/thetherington/jobber-gig/internal/adapters/handler/rabbitmq"
	"github.com/thetherington/jobber-gig/internal/adapters/storage/elasticsearch"
	"github.com/thetherington/jobber-gig/internal/adapters/storage/mongodb"
	"github.com/thetherington/jobber-gig/internal/adapters/storage/redis"
	"github.com/thetherington/jobber-gig/internal/core/service"
)

const (
	WEB_PORT  = 5004
	GRPC_PORT = 4004
	APP_NAME  = "gig-service"
	Index     = "gigs"
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

	// Init Cache service
	slog.Info("Creating Redis Cache Connection", "address", config.Redis.Host)
	cache, err := redis.New(context.Background(), config.Redis)
	if err != nil {
		slog.With("error", err).Error("Failed to make redis client connection")
		os.Exit(1)
	}
	defer cache.Close()

	// Create the Cloudinary image upload dependency
	image, err := cloudinary.New(config.Cloudinary.Name, config.Cloudinary.ApiKey, config.Cloudinary.ApiSecret)
	if err != nil {
		slog.With("error", err).Error("failed to setup cloudinary credentials")
		os.Exit(1)
	}

	// Create the Gig service with dependencies (mongodb, elasticsearch, rabbitmq, redis, cloudinary)
	gig := service.NewGigService(client.Database(config.DB.Name), esClient, nil, cache, image)

	// Create the RabbitMQ connection and add the consumers for x service and x service
	consumer := rabbitmq.NewRabbitMQAdapter(config.RabbitMQ.Endpoint, gig)
	defer consumer.Close()

	// gig update from review and order services
	if err := consumer.AddConsumer(consumer.ConsumeGigDirectMessage); err != nil {
		slog.With("error", err).Error("Failed to Add ConsumeGigDirectMessage")
		os.Exit(1)
	}

	// gig seed from users services (random users feed)
	if err := consumer.AddConsumer(consumer.ConsumeSeedDirectMessage); err != nil {
		slog.With("error", err).Error("Failed to Add ConsumeSeedDirectMessage")
		os.Exit(1)
	}

	// set the queue dependency to publish messages
	gig.SetQueue(consumer)

	// Create the gRPC Server with the gig service
	grpc := grpc.NewGrpcAdapter(gig, GRPC_PORT)

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
