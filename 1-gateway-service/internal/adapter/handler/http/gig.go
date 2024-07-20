package http

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/thetherington/jobber-common/models/gig"
	"github.com/thetherington/jobber-common/models/search"
	"github.com/thetherington/jobber-common/utils"
	"github.com/thetherington/jobber-gateway/internal/core/port"
)

// GigHandler represents the HTTP handler for gig service requests
type GigHandler struct {
	svc port.GigService
}

// /api/gateway/v1/gig
func (gh GigHandler) Routes(router chi.Router) {
	router.Get("/{gigId}", gh.GetGigByID)
	router.Get("/seller/{sellerId}", gh.GetSellerGigs)
	router.Get("/seller/pause/{sellerId}", gh.GetSellerPausedGigs)

	router.Get("/search/{from}/{size}/{type}", gh.SearchGigs)
	router.Get("/category/{username}", gh.GetGigsByCategory)
	router.Get("/top/{username}", gh.GetTopRatedGigsByCategory)
	router.Get("/similar/{gigId}", gh.GetSimiliarGigs)

	router.Post("/create", gh.CreateGig)

	router.Put("/{gigId}", gh.UpdateGig)
	router.Put("/active/{gigId}", gh.UpdateActiveGig)
	router.Put("/seed/{count}", gh.SeedGigs)

	router.Delete("/{gigId}/{sellerId}", gh.DeleteGig)
}

// NewGigHandlerr creates a new GigHandler instance
func NewGigHandler(svc port.GigService) *GigHandler {
	return &GigHandler{
		svc,
	}
}

func (gh *GigHandler) CreateGig(w http.ResponseWriter, r *http.Request) {
	var gigPayload *gig.SellerGig

	// unmarshal the request body into the signup
	if err := ReadJSON(w, r, &gigPayload); err != nil {
		ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// send a request to the gig microservice via grpc client
	resp, err := gh.svc.CreateGig(r.Context(), gigPayload)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusCreated, &resp); err != nil {
		slog.With("error", err).Error("CreateGig: failed to write http response")
	}
}

func (gh *GigHandler) UpdateGig(w http.ResponseWriter, r *http.Request) {
	var gigPayload *gig.SellerGig

	// unmarshal the request body into the signup
	if err := ReadJSON(w, r, &gigPayload); err != nil {
		ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// send a request to the gig microservice via grpc client
	resp, err := gh.svc.UpdateGig(r.Context(), chi.URLParam(r, "gigId"), gigPayload)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("UpdateGig: failed to write http response")
	}
}

func (gh *GigHandler) DeleteGig(w http.ResponseWriter, r *http.Request) {
	// send a request to the gig microservice via grpc client
	msg, err := gh.svc.DeleteGig(r.Context(), chi.URLParam(r, "gigId"), chi.URLParam(r, "sellerId"))
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
		slog.With("error", err).Error("DeleteGig: failed to write http response")
	}
}

func (gh *GigHandler) GetGigByID(w http.ResponseWriter, r *http.Request) {
	// send a request to the gig microservice via grpc client
	resp, err := gh.svc.GetGigById(r.Context(), chi.URLParam(r, "gigId"))
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("GetGigByID: failed to write http response")
	}
}

func (gh *GigHandler) GetSellerGigs(w http.ResponseWriter, r *http.Request) {
	// send a request to the gig microservice via grpc client
	resp, err := gh.svc.GetSellerGigs(r.Context(), chi.URLParam(r, "sellerId"))
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("GetSellerGigs: failed to write http response")
	}
}

func (gh *GigHandler) GetSellerPausedGigs(w http.ResponseWriter, r *http.Request) {
	// send a request to the gig microservice via grpc client
	resp, err := gh.svc.GetSellerPausedGigs(r.Context(), chi.URLParam(r, "sellerId"))
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("GetSellerPausedGigs: failed to write http response")
	}
}

func (gh *GigHandler) SearchGigs(w http.ResponseWriter, r *http.Request) {
	// compute the size into a number
	size, err := strconv.Atoi(chi.URLParam(r, "size"))
	if err != nil {
		ErrorJSON(w, fmt.Errorf("size is not a number"), http.StatusBadRequest)
		return
	}

	// create the paginate struct
	paginate := search.PaginateProps{
		From: chi.URLParam(r, "from"),
		Size: size,
		Type: chi.URLParam(r, "type"),
	}

	// start with a search request struct with nil for pointers
	req := search.SearchRequest{
		SearchQuery:   r.URL.Query().Get("query"), // empty string is ok here
		PaginateProps: &paginate,
		DeliveryTime:  nil,
		Min:           nil,
		Max:           nil,
	}

	// validate each query param option known and get the pointer of the value
	if v := r.URL.Query().Get("delivery_time"); v != "" {
		req.DeliveryTime = utils.Ptr(v)
	}

	// min/max parse into a float
	if v := r.URL.Query().Get("minPrice"); v != "" {
		min, err := strconv.ParseFloat(v, 64)
		if err != nil {
			ErrorJSON(w, fmt.Errorf("min price is not a number"), http.StatusBadRequest)
			return
		}

		req.Min = utils.PtrF64(min)
	}

	if v := r.URL.Query().Get("maxPrice"); v != "" {
		max, err := strconv.ParseFloat(v, 64)
		if err != nil {
			ErrorJSON(w, fmt.Errorf("max price is not a number"), http.StatusBadRequest)
			return
		}

		req.Max = utils.PtrF64(max)
	}

	// send a request to the gig microservice via grpc client
	resp, err := gh.svc.SearchGig(r.Context(), req)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("SearchGigs: failed to write http response")
	}
}

func (gh *GigHandler) GetGigsByCategory(w http.ResponseWriter, r *http.Request) {
	// send a request to the gig microservice via grpc client
	resp, err := gh.svc.SearchGigCategory(r.Context(), chi.URLParam(r, "username"))
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("GetGigsByCategory: failed to write http response")
	}
}

func (gh *GigHandler) GetSimiliarGigs(w http.ResponseWriter, r *http.Request) {
	// send a request to the gig microservice via grpc client
	resp, err := gh.svc.SearchGigSimilar(r.Context(), chi.URLParam(r, "gigId"))
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("GetSimiliarGigs: failed to write http response")
	}
}

func (gh *GigHandler) GetTopRatedGigsByCategory(w http.ResponseWriter, r *http.Request) {
	// send a request to the gig microservice via grpc client
	resp, err := gh.svc.SearchGigTop(r.Context(), chi.URLParam(r, "username"))
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("GetTopRatedGigsByCategory: failed to write http response")
	}
}

func (gh *GigHandler) UpdateActiveGig(w http.ResponseWriter, r *http.Request) {
	var active struct{ Active bool }

	// unmarshal the request body into the bool struct
	if err := ReadJSON(w, r, &active); err != nil {
		ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	resp, err := gh.svc.UpdateActiveGig(r.Context(), chi.URLParam(r, "gigId"), active.Active)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("UpdateActiveGig: failed to write http response")
	}
}

func (gh *GigHandler) SeedGigs(w http.ResponseWriter, r *http.Request) {
	count, err := strconv.Atoi(chi.URLParam(r, "count"))
	if err != nil {
		ErrorJSON(w, fmt.Errorf("count is not a number"), http.StatusBadRequest)
		return
	}

	// send a request to the gig microservice via grpc client
	resp, err := gh.svc.SeedGigs(r.Context(), int32(count))
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("SeedGigs: failed to write http response")
	}
}
