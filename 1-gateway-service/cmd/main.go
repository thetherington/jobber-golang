package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	token "github.com/thetherington/jobber-common/client-token"
	"github.com/thetherington/jobber-common/logger"
	"github.com/thetherington/jobber-gateway/internal/adapter/config"
	"github.com/thetherington/jobber-gateway/internal/adapter/handler/grpc"
	"github.com/thetherington/jobber-gateway/internal/adapter/handler/http"
	"github.com/thetherington/jobber-gateway/internal/adapter/handler/websocket"
	"github.com/thetherington/jobber-gateway/internal/adapter/storage/redis"
	"github.com/thetherington/jobber-gateway/internal/core/service"
)

const (
	WEB_PORT = 4000
	APP_NAME = "gateway-service"
)

var sessionManager *scs.SessionManager

func init() {
	gob.Register(token.Payload{})
}

func main() {
	// Load environment variables
	config := config.Config

	// Set logger
	logger.Set(config.App.Env, APP_NAME, config.App.LogLevel)

	slog.Info("Creating Redis Cache Connection", "address", config.Redis.Host)

	// Init Cache service
	cache, err := redis.New(context.Background(), config.Redis)
	if err != nil {
		slog.With("error", err).Error("Failed to make redis client connection")
		os.Exit(1)
	}
	defer cache.Close()

	// Initialize a new session manager and configure the session lifetime.
	sessionManager = scs.New()

	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Name = "jobber"
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.Secure = false

	// use redis as a cache to store the session cookie
	sessionManager.Store = cache

	// if config.App.Env == "Production" {
	// 	sessionManager.Cookie.Secure = true
	// 	sessionManager.Cookie.SameSite = 0
	// }

	// Init Token Service
	clientTokenHandler, err := token.NewClientJWTMaker(config.Tokens.JWT)
	if err != nil {
		slog.With("error", err).Error("failed to create token maker service")
		os.Exit(1)
	}

	// Init Middlware handler for token verifying
	middleware := http.NewMiddlewareHandler(sessionManager, clientTokenHandler)

	// Init Websocket handler for frontend. redis cache is injected as a dependency
	// web socket manager is injected into the chat grpc client
	// web socket manager is provided to the http.NewRouter
	socketHandler := websocket.NewManager(cache)

	// Dependency injection
	// --------------------
	//
	// Init Auth Service
	slog.Info("Creating gRPC Auth Service Client", "address", config.Services.Auth)
	authClient := grpc.NewAuthAdapter(config.Services.Auth, sessionManager)
	authService := service.NewAuthService(authClient, cache, socketHandler)
	authHandler := http.NewAuthHandler(authService, sessionManager)

	// Init Search service
	searchService := service.NewSearchService(authClient)
	searchHandler := http.NewSearchHandler(searchService)

	// Init Buyer Service
	slog.Info("Creating gRPC Users Service Client", "address", config.Services.Users)
	usersClient := grpc.NewUsersAdapter(config.Services.Users, sessionManager)
	buyerService := service.NewBuyerService(usersClient)
	buyerHandler := http.NewBuyerHandler(buyerService)

	// Init Seller Service
	sellerService := service.NewSellerService(usersClient)
	sellerHandler := http.NewSellerHandler(sellerService)

	// Init Gig Service
	slog.Info("Creating gRPC Gig Service Client", "address", config.Services.Gig)
	gigClient := grpc.NewGigAdapter(config.Services.Gig, sessionManager)
	gigService := service.NewGigService(gigClient)
	gigHandler := http.NewGigHandler(gigService)

	// Init Chat Service
	slog.Info("Creating gRPC Chat Service Client", "address", config.Services.Message)
	chatClient := grpc.NewChatAdapter(config.Services.Message, sessionManager, socketHandler)
	chatService := service.NewChatService(chatClient)
	chatHandler := http.NewChatHandler(chatService)

	// Init Order Service
	slog.Info("Creating gRPC Order Service Client", "address", config.Services.Order)
	orderClient := grpc.NewOrderAdapter(config.Services.Order, sessionManager, socketHandler)
	orderService := service.NewOrderService(orderClient)
	orderHandler := http.NewOrderHandler(orderService)

	// Init Review Service
	slog.Info("Creating gRPC Review Service Client", "address", config.Services.Order)
	reviewClient := grpc.NewReviewAdapter(config.Services.Review, sessionManager)
	reviewService := service.NewReviewService(reviewClient)
	reviewHandler := http.NewReviewHandler(reviewService)

	// subscribe to chat service message stream
	go chatClient.SubscribeStream()
	defer chatClient.DisconnectStream()

	// subscribe to notification service message stream
	go orderClient.SubscribeNotifyStream()
	defer orderClient.DisconnectStream()

	// subscribe to order service order stream
	go orderClient.SubscribeOrderStream()
	defer orderClient.DisconnectOrderStream()

	// Init HTTP Router
	router := http.NewRouter(socketHandler, *middleware, *authHandler, *searchHandler,
		*buyerHandler, *sellerHandler, *gigHandler, *chatHandler, *orderHandler, *reviewHandler,
	)

	addr := fmt.Sprintf(":%d", WEB_PORT)

	slog.Info("Starting HTTP Gateway Server", "address", addr)
	if err := http.Serve(addr, sessionManager.LoadAndSave(router)); err != nil {
		slog.With("error", err).Error("Failed to start HTTP Server")
	}
}
