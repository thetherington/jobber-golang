package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/thetherington/jobber-common/logger"
	"github.com/thetherington/jobber-users/internal/adapters/config"
	"github.com/thetherington/jobber-users/internal/adapters/handler/grpc"
	"github.com/thetherington/jobber-users/internal/adapters/handler/http"
	"github.com/thetherington/jobber-users/internal/adapters/handler/rabbitmq"
	"github.com/thetherington/jobber-users/internal/adapters/storage/mongodb"
	"github.com/thetherington/jobber-users/internal/core/service"
)

const (
	WEB_PORT  = 5003
	GRPC_PORT = 4003
	APP_NAME  = "users-service"
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

	// Initialize buyer and seller services with mongo database
	buyer := service.NewBuyerService(client.Database(config.DB.Name))
	seller := service.NewSellerService(client.Database(config.DB.Name), buyer)

	// Create the RabbitMQ connection and add the consumers for authentication service and order service
	consumer := rabbitmq.NewRabbitMQAdapter(config.RabbitMQ.Endpoint, seller, buyer)
	defer consumer.Close()

	// buyer creation from auth service queue
	if err := consumer.AddConsumer(consumer.ConsumeBuyerDirectMessage); err != nil {
		slog.With("error", err).Error("Failed to Add ConsumeBuyerDirectMessage")
		os.Exit(1)
	}

	// Seller order update from gig, order service queue
	if err := consumer.AddConsumer(consumer.ConsumeSellerDirectMessage); err != nil {
		slog.With("error", err).Error("Failed to Add ConsumeSellerDirectMessage")
		os.Exit(1)
	}

	// Seller / Gig seeding returns random sellers to gig service
	if err := consumer.AddConsumer(consumer.ConsumeSeedGigDirectMessages); err != nil {
		slog.With("error", err).Error("Failed to Add ConsumeSeedGigDirectMessages")
		os.Exit(1)
	}

	// Review fanout exchange from Review microservice
	if err := consumer.AddConsumer(consumer.ConsumeReviewFanoutMessages); err != nil {
		slog.With("error", err).Error("Failed to Add ConsumeReviewFanoutMessages")
		os.Exit(1)
	}

	// Relay offer messages to email service (get buyer email)
	if err := consumer.AddConsumer(consumer.ConsumeNotificationRelay); err != nil {
		slog.With("error", err).Error("Failed to Add ConsumeNotificationRelay")
		os.Exit(1)
	}

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
	grpc := grpc.NewGrpcAdapter(buyer, seller, GRPC_PORT)

	slog.Info("Starting gRPC Server", "port", GRPC_PORT)
	grpc.Run()
}
