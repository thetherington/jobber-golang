package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter creates a new HTTP router
func NewRouter() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)

	router.Get("/auth-health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		w.Write([]byte("Auth service is healthy and OK."))
	})

	return router
}

// Serve starts the HTTP server
func Serve(listenAddr string, handler http.Handler) error {
	return http.ListenAndServe(listenAddr, handler)
}
