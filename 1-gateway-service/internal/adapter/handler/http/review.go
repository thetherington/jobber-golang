package http

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/thetherington/jobber-common/models/review"
	"github.com/thetherington/jobber-gateway/internal/core/port"
)

// ReviewHandler represents the HTTP handler for review service requests
type ReviewHandler struct {
	svc port.ReviewService
}

// /api/gateway/v1/review
func (rh ReviewHandler) Routes(router chi.Router) {
	router.Get("/gig/{gigId}", rh.GetReviewsByGigId)
	router.Get("/seller/{sellerId}", rh.GetReviewsBySellerId)

	router.Post("/", rh.CreateReview)
}

// NewReviewHandler creates a new ReviewHandler instance
func NewReviewHandler(svc port.ReviewService) *ReviewHandler {
	return &ReviewHandler{
		svc,
	}
}

func (rh *ReviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	var review *review.ReviewDocument

	// unmarshal the request body into the order
	if err := ReadJSON(w, r, &review); err != nil {
		ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	resp, err := rh.svc.CreateReview(r.Context(), review)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusCreated, &resp); err != nil {
		slog.With("error", err).Error("CreateReview: failed to write http response")
	}
}

func (rh *ReviewHandler) GetReviewsByGigId(w http.ResponseWriter, r *http.Request) {
	resp, err := rh.svc.GetReviewsByGigId(r.Context(), chi.URLParam(r, "gigId"))
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("GetReviewsByGigId: failed to write http response")
	}
}

func (rh *ReviewHandler) GetReviewsBySellerId(w http.ResponseWriter, r *http.Request) {
	resp, err := rh.svc.GetReviewsBySellerId(r.Context(), chi.URLParam(r, "sellerId"))
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("GetReviewsBySellerId: failed to write http response")
	}
}
