package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/thetherington/jobber-common/cloudinary"
	"github.com/thetherington/jobber-common/logger"
	"github.com/thetherington/jobber-order/internal/adapters/config"
	"github.com/thetherington/jobber-order/internal/adapters/handler/grpc"
	"github.com/thetherington/jobber-order/internal/adapters/handler/http"
	"github.com/thetherington/jobber-order/internal/adapters/handler/rabbitmq"
	"github.com/thetherington/jobber-order/internal/adapters/payment/stripe"
	"github.com/thetherington/jobber-order/internal/adapters/storage/mongodb"
	"github.com/thetherington/jobber-order/internal/core/service"
)

const (
	WEB_PORT  = 5006
	GRPC_PORT = 4006
	APP_NAME  = "order-service"
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

	// Create the Cloudinary image upload dependency
	image, err := cloudinary.New(config.Cloudinary.Name, config.Cloudinary.ApiKey, config.Cloudinary.ApiSecret)
	if err != nil {
		slog.With("error", err).Error("failed to setup cloudinary credentials")
		os.Exit(1)
	}

	// Create Payment Service
	payment := stripe.NewStripePayment(config.Stripe)

	order := service.NewOrderService(client.Database(config.DB.Name), nil, image, payment)

	notification := service.NewNotificatonService(client.Database(config.DB.Name))

	// Create the gRPC Server with the gig service
	grpc := grpc.NewGrpcAdapter(order, notification, GRPC_PORT)

	// Create the RabbitMQ connection
	consumer := rabbitmq.NewRabbitMQAdapter(config.RabbitMQ.Endpoint, order, grpc)
	defer consumer.Close()

	// review messages consumer from review service
	if err := consumer.AddConsumer(consumer.ConsumerReviewFanoutMessages); err != nil {
		slog.With("error", err).Error("Failed to Add ConsumerReviewFanoutMessages")
		os.Exit(1)
	}

	// set the queue dependency to publish messages from the order service
	order.SetQueue(consumer)

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
