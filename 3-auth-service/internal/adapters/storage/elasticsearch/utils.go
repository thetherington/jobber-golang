package elasticsearch

import (
	"errors"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
)

var (
	ErrGigNotFound   = errors.New("gig not found")
	ErrIndexNotExist = errors.New("index doesn't exist")
)

// type === forward ? ("asc") : ("desc")
func sortlookup(sort string) *sortorder.SortOrder {
	if sort == "forward" {
		return &sortorder.Asc
	}

	return &sortorder.Desc
}

// returns a pointer of types.Float64 from float64
func tF64(f *float64) *types.Float64 {
	t := types.Float64(*f)
	return &t
}
