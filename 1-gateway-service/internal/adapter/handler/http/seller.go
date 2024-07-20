package http

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/thetherington/jobber-common/models/users"
	"github.com/thetherington/jobber-gateway/internal/core/port"
)

// SellerHandler represents the HTTP handler for seller service requests
type SellerHandler struct {
	svc port.SellerService
}

// /api/gateway/v1/seller
func (sh SellerHandler) Routes(router chi.Router) {
	router.Get("/id/{sellerId}", sh.GetSellerById)
	router.Get("/username/{username}", sh.GetSellerByUsername)
	router.Get("/random/{size}", sh.GetRandomSellers)

	router.Post("/create", sh.CreateSeller)

	router.Put("/{sellerId}", sh.UpdateSeller)
	router.Put("/seed/{count}", sh.SeedSellers)
}

// NewSellerHandler creates a new SellerHandler instance
func NewSellerHandler(svc port.SellerService) *SellerHandler {
	return &SellerHandler{
		svc,
	}
}

func (sh *SellerHandler) CreateSeller(w http.ResponseWriter, r *http.Request) {
	var sellerPayload *users.Seller

	// unmarshal the request body into the sellerPayload
	if err := ReadJSON(w, r, &sellerPayload); err != nil {
		ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	resp, err := sh.svc.CreateSeller(r.Context(), sellerPayload)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusCreated, &resp); err != nil {
		slog.With("error", err).Error("CreateSeller: failed to write http response")
	}
}

func (sh *SellerHandler) UpdateSeller(w http.ResponseWriter, r *http.Request) {
	var (
		sellerId      string
		sellerPayload *users.Seller
	)

	// get the sellerId from the url
	if sellerId = chi.URLParam(r, "sellerId"); sellerId == "" {
		ErrorJSON(w, fmt.Errorf("missing id from url"), http.StatusBadRequest)
		return
	}

	// unmarshal the request body into the sellerPayload
	if err := ReadJSON(w, r, &sellerPayload); err != nil {
		ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	resp, err := sh.svc.UpdateSeller(r.Context(), sellerId, sellerPayload)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("UpdateSeller: failed to write http response")
	}
}

func (sh *SellerHandler) GetSellerById(w http.ResponseWriter, r *http.Request) {
	var sellerId string

	// get the sellerId from the url
	if sellerId = chi.URLParam(r, "sellerId"); sellerId == "" {
		ErrorJSON(w, fmt.Errorf("missing id from url"), http.StatusBadRequest)
		return
	}

	resp, err := sh.svc.GetSellerById(r.Context(), sellerId)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("GetSellerById: failed to write http response")
	}
}

func (sh *SellerHandler) GetSellerByUsername(w http.ResponseWriter, r *http.Request) {
	var username string

	// get the sellerId from the url
	if username = chi.URLParam(r, "username"); username == "" {
		ErrorJSON(w, fmt.Errorf("missing username from url"), http.StatusBadRequest)
		return
	}

	resp, err := sh.svc.GetSellerByUsername(r.Context(), username)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("GetSellerByUsername: failed to write http response")
	}
}

func (sh *SellerHandler) GetRandomSellers(w http.ResponseWriter, r *http.Request) {
	var param string

	// get the size from the url
	if param = chi.URLParam(r, "size"); param == "" {
		ErrorJSON(w, fmt.Errorf("missing size from url"), http.StatusBadRequest)
		return
	}

	size, err := strconv.Atoi(param)
	if err != nil {
		ErrorJSON(w, fmt.Errorf("size is not a number"), http.StatusBadRequest)
		return
	}

	resp, err := sh.svc.GetRandomSellers(r.Context(), int32(size))
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("GetRandomSellers: failed to write http response")
	}
}

func (sh *SellerHandler) SeedSellers(w http.ResponseWriter, r *http.Request) {
	var param string

	// get the count from the url
	if param = chi.URLParam(r, "count"); param == "" {
		ErrorJSON(w, fmt.Errorf("missing count from url"), http.StatusBadRequest)
		return
	}

	count, err := strconv.Atoi(param)
	if err != nil {
		ErrorJSON(w, fmt.Errorf("count is not a number"), http.StatusBadRequest)
		return
	}

	msg, err := sh.svc.SeedSellers(r.Context(), int32(count))
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	resp := &struct {
		Message string `json:"message"`
	}{
		Message: msg,
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("SeedSellers: failed to write http response")
	}
}
