package http

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/thetherington/jobber-gateway/internal/core/port"
)

// BuyerHandler represents the HTTP handler for buyer service requests
type BuyerHandler struct {
	svc port.BuyerService
}

// /api/gateway/v1/buyer
func (bh BuyerHandler) Routes(router chi.Router) {
	router.Get("/email", bh.GetBuyerByEmail)
	router.Get("/username", bh.GetBuyerByUsername)
	router.Get("/{username}", bh.GetBuyerByProvidedUsername)
}

// NewBuyerHandler creates a new BuyerHandler instance
func NewBuyerHandler(svc port.BuyerService) *BuyerHandler {
	return &BuyerHandler{
		svc,
	}
}

func (bh *BuyerHandler) GetBuyerByEmail(w http.ResponseWriter, r *http.Request) {
	resp, err := bh.svc.GetBuyerByEmail(r.Context())
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("GetBuyerByEmail: failed to write http response")
	}
}

func (bh *BuyerHandler) GetBuyerByUsername(w http.ResponseWriter, r *http.Request) {
	resp, err := bh.svc.GetBuyerByUsername(r.Context())
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("GetBuyerByUsername: failed to write http response")
	}
}

func (bh *BuyerHandler) GetBuyerByProvidedUsername(w http.ResponseWriter, r *http.Request) {
	var username string

	// get the token from the url
	if username = chi.URLParam(r, "username"); username == "" {
		ErrorJSON(w, fmt.Errorf("missing username"), http.StatusBadRequest)
		return
	}

	resp, err := bh.svc.GetBuyerByProvidedUsername(r.Context(), username)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("GetBuyerByProvidedUsername: failed to write http response")
	}
}
