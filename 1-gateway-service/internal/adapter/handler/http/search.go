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

// SearchHandler represents the HTTP handler for search service requests
type SearchHandler struct {
	svc port.SearchService
}

// /api/gateway/v1/auth
func (sh SearchHandler) Routes(router chi.Router) {
	router.Get("/search/gig/{from}/{size}/{type}", sh.SearchGigs)

	router.Get("/search/gig/{gigId}", sh.GetGigByID)
}

// NewAuthHandler creates a new AuthHandler instance
func NewSearchHandler(svc port.SearchService) *SearchHandler {
	return &SearchHandler{
		svc,
	}
}

func (sh *SearchHandler) GetGigByID(w http.ResponseWriter, r *http.Request) {
	var id string

	// get the gig id from the url
	if id = chi.URLParam(r, "gigId"); id == "" {
		ErrorJSON(w, fmt.Errorf("missing gig id"), http.StatusBadRequest)
		return
	}

	// send request to auth microservice via gRPC
	resp, err := sh.svc.GetGigByID(r.Context(), id)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	payload := &struct {
		Message string         `json:"message"`
		Gig     *gig.SellerGig `json:"gig"`
	}{
		Message: "Single Gig Result",
		Gig:     resp,
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &payload); err != nil {
		slog.With("error", err).Error("GetGigByID: failed to write http response")
	}
}

// auth/search/gig/0/2/forward?query=Music&minPrice=0&maxPrice=99&delivery_time=1
func (sh *SearchHandler) SearchGigs(w http.ResponseWriter, r *http.Request) {
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

	// send request to auth microservice via gRPC
	resp, err := sh.svc.SearchGigs(r.Context(), req)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	payload := &struct {
		Message string           `json:"message"`
		Total   int64            `json:"total"`
		Gigs    []*gig.SellerGig `json:"gigs"`
	}{
		Message: "Search gig results",
		Total:   resp.Total,
		Gigs:    resp.Hits,
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &payload); err != nil {
		slog.With("error", err).Error("SearchGigs: failed to write http response")
	}
}
