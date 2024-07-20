package http

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	slogchi "github.com/samber/slog-chi"
	"github.com/thetherington/jobber-common/logger"
	"github.com/thetherington/jobber-gateway/internal/adapter/config"
	"github.com/thetherington/jobber-gateway/internal/adapter/handler/websocket"
	"go.elastic.co/apm/module/apmchiv5/v2"
)

const BASE_PATH string = "/api/gateway/v1"

// NewRouter creates a new HTTP router
func NewRouter(
	socketHandler *websocket.Manager,
	middlwareHandler MiddlewareHandler,
	authHandler AuthHandler,
	searchHandler SearchHandler,
	buyerHandler BuyerHandler,
	sellerHandler SellerHandler,
	gigHandler GigHandler,
	chatHandler ChatHandler,
	orderHandler OrderHandler,
	reviewHandler ReviewHandler,
) http.Handler {
	router := chi.NewRouter()

	router.Use(slogchi.NewWithConfig(logger.Logger, slogchi.Config{
		ClientErrorLevel: slog.LevelDebug,
		DefaultLevel:     slog.LevelDebug,
		ServerErrorLevel: slog.LevelError,
	}))

	router.Use(middleware.Recoverer)

	router.Use(apmchiv5.Middleware())

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{config.Config.App.ClientUrl},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.Use(middleware.Heartbeat("/gateway-health"))

	// Web socket endpoint where ServeWS upgrades the connection
	router.Get("/ws", socketHandler.ServeWS)

	// Auth Service Routes /api/gateway/v1/auth
	router.Route(BASE_PATH+"/auth", func(auth chi.Router) {
		// unauthenticated routes for auth service (signup/login/password reset)
		auth.Group(func(r chi.Router) {
			authHandler.Routes(r)
		})

		// authenticated routes for auth service (current user, token refresh, email resend)
		auth.Group(func(r chi.Router) {
			r.Use(middlwareHandler.VerifyUser)
			r.Use(middlwareHandler.CheckAuthentication)

			authHandler.CurrentUserRoutes(r)
		})

		// unauthenticated routes for searching gigs and getting gigbyid
		auth.Group(func(r chi.Router) {
			searchHandler.Routes(r)
		})
	})

	// Users Service Routes /api/gateway/v1/buyer
	router.Route(BASE_PATH+"/buyer", func(r chi.Router) {
		// Buyer routes for get buyer by email, username and supplied username
		// api/gateway/v1/buyer
		r.Use(middlwareHandler.VerifyUser)
		r.Use(middlwareHandler.CheckAuthentication)

		buyerHandler.Routes(r)
	})

	// Users Service Routes /api/gateway/v1/seller
	router.Route(BASE_PATH+"/seller", func(r chi.Router) {
		// Seller routes for create seller get seller by id, username
		// seed and get random sellers

		// api/gateway/v1/seller
		r.Use(middlwareHandler.VerifyUser)
		r.Use(middlwareHandler.CheckAuthentication)

		sellerHandler.Routes(r)
	})

	// Gig Service Routes /api/gateway/v1/gig
	router.Route(BASE_PATH+"/gig", func(r chi.Router) {
		// gig routes for crud and search

		// api/gateway/v1/gig
		r.Use(middlwareHandler.VerifyUser)
		r.Use(middlwareHandler.CheckAuthentication)

		gigHandler.Routes(r)
	})

	// Chat Service Routes /api/gateway/v1/message
	router.Route(BASE_PATH+"/message", func(r chi.Router) {
		// chat routes for crud

		// api/gateay/v1/message
		r.Use(middlwareHandler.VerifyUser)
		r.Use(middlwareHandler.CheckAuthentication)

		chatHandler.Routes(r)
	})

	// Order Service Routes /api/gateway/v1/order
	router.Route(BASE_PATH+"/order", func(r chi.Router) {
		// order routes for crud

		// api/gateay/v1/order
		r.Use(middlwareHandler.VerifyUser)
		r.Use(middlwareHandler.CheckAuthentication)

		orderHandler.Routes(r)
	})

	// Review Service Routes /api/gateway/v1/review
	router.Route(BASE_PATH+"/review", func(r chi.Router) {
		// review routes for crud

		// api/gateay/v1/order
		r.Use(middlwareHandler.VerifyUser)
		r.Use(middlwareHandler.CheckAuthentication)

		reviewHandler.Routes(r)
	})

	return router
}

// Serve starts the HTTP server
func Serve(listenAddr string, handler http.Handler) error {
	return http.ListenAndServe(listenAddr, handler)
}
