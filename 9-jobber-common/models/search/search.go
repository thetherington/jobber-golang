package search

import "github.com/thetherington/jobber-common/models/gig"

type PaginateProps struct {
	From string `json:"from"`
	Size int    `json:"size"`
	Type string `json:"type"`
}

type SearchResponse struct {
	Total int64
	Hits  []*gig.SellerGig
}

type SearchRequest struct {
	SearchQuery   string         `json:"searchQuery"`
	PaginateProps *PaginateProps `json:"PaginateProps,omitempty"`
	DeliveryTime  *string        `json:"deliveryTime,omitempty"`
	Min           *float64       `json:"min,omitempty"`
	Max           *float64       `json:"max,omitempty"`
}
