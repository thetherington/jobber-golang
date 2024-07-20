package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"

	essearch "github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptlanguage"
	"github.com/thetherington/jobber-common/models/gig"
	"github.com/thetherington/jobber-common/models/search"
	"github.com/thetherington/jobber-common/utils"
)

func (es *ESSearch) SearchGigs(ctx context.Context, searchQuery string, paginate *search.PaginateProps, deliveryTime *string, min *float64, max *float64) (int64, []*gig.SellerGig, error) {
	req := essearch.NewRequest()

	mustBoolSlice := make([]types.Query, 0)

	// if empty string then use asterix as wildcard
	if searchQuery == "" {
		searchQuery = "*"
	}

	// searchQuery is match string to all gig fields possible
	mustBoolSlice = append(mustBoolSlice, types.Query{
		QueryString: &types.QueryStringQuery{
			Fields: []string{"username", "title", "description", "basicDescription", "basicTitle", "categories", "subCategories", "tags"},
			Query:  searchQuery,
		},
	})

	// only get active gigs
	mustBoolSlice = append(mustBoolSlice, types.Query{
		Term: map[string]types.TermQuery{
			"active": {
				Value: true,
			},
		},
	})

	// if the deliveryTime is provided then match exactly (1 Day Delivery, 2 Days Delivery....)
	if deliveryTime != nil {
		mustBoolSlice = append(mustBoolSlice, types.Query{
			MatchPhrase: map[string]types.MatchPhraseQuery{
				"expectedDelivery": {
					Query: *deliveryTime,
				},
			},
		})
	}

	// if min and max are provided (not nil) then find prices in the range
	if min != nil && max != nil {
		mustBoolSlice = append(mustBoolSlice, types.Query{
			Range: map[string]types.RangeQuery{
				"price": types.NumberRangeQuery{
					Gte: tF64(min),
					Lte: tF64(max),
				},
			},
		})
	}

	req.Query = &types.Query{
		Bool: &types.BoolQuery{
			Must: mustBoolSlice,
		},
	}

	// pagination for number of gigs, direction and page
	if paginate != nil {
		// set the return size
		req.Size = utils.PtrI(paginate.Size)

		// set the sort direction (type === forward ? ("asc") : "desc")
		req.Sort = []types.SortCombinations{
			&types.SortOptions{
				SortOptions: map[string]types.FieldSort{
					"sortId": {
						Order: sortlookup(paginate.Type),
					},
				},
			},
		}

		// set the page if From is anything but "0"
		if paginate.From != "0" {
			req.SearchAfter = []types.FieldValue{
				paginate.From,
			}
		}
	}

	resp, err := es.client.Search().Index(es.index).Request(req).Do(ctx)
	if err != nil {
		return 0, nil, fmt.Errorf("request failed: %w", err)
	}

	Gigs := make([]*gig.SellerGig, 0)

	for _, hit := range resp.Hits.Hits {
		var gig gig.SellerGig

		// need to convert the data into json.
		b, _ := hit.Source_.MarshalJSON()
		if err != nil {
			return 0, nil, fmt.Errorf("marshal failed %w", err)
		}

		// convert the json data back into a struct.
		err = json.Unmarshal(b, &gig)
		if err != nil {
			return 0, nil, fmt.Errorf("unmarshal failed into SellerGig %w", err)
		}

		Gigs = append(Gigs, &gig)
	}

	return resp.Hits.Total.Value, Gigs, nil
}

func (es *ESSearch) GigsSearchBySellerId(ctx context.Context, id string, active bool) ([]*gig.SellerGig, error) {
	req := essearch.NewRequest()

	mustBoolSlice := make([]types.Query, 0)

	// match seller id
	mustBoolSlice = append(mustBoolSlice, types.Query{
		QueryString: &types.QueryStringQuery{
			Fields: []string{"sellerId"},
			Query:  id,
		},
	})

	// get active = true / false
	mustBoolSlice = append(mustBoolSlice, types.Query{
		Term: map[string]types.TermQuery{
			"active": {
				Value: active,
			},
		},
	})

	req.Query = &types.Query{
		Bool: &types.BoolQuery{
			Must: mustBoolSlice,
		},
	}

	resp, err := es.client.Search().Index(es.index).Request(req).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	Gigs := make([]*gig.SellerGig, 0)

	for _, hit := range resp.Hits.Hits {
		var gig gig.SellerGig

		// need to convert the data into json.
		b, _ := hit.Source_.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("marshal failed %w", err)
		}

		// convert the json data back into a struct.
		err = json.Unmarshal(b, &gig)
		if err != nil {
			return nil, fmt.Errorf("unmarshal failed into SellerGig %w", err)
		}

		Gigs = append(Gigs, &gig)
	}

	return Gigs, nil
}

func (es *ESSearch) SearchGigsByCategory(ctx context.Context, category string) (int64, []*gig.SellerGig, error) {
	req := essearch.NewRequest()

	mustBoolSlice := make([]types.Query, 0)

	// if empty string then use asterix as wildcard
	if category == "" {
		category = "*"
	}

	// match seller id
	mustBoolSlice = append(mustBoolSlice, types.Query{
		QueryString: &types.QueryStringQuery{
			Fields: []string{"categories"},
			Query:  category,
		},
	})

	// get active = true / false
	mustBoolSlice = append(mustBoolSlice, types.Query{
		Term: map[string]types.TermQuery{
			"active": {
				Value: true,
			},
		},
	})

	req.Query = &types.Query{
		Bool: &types.BoolQuery{
			Must: mustBoolSlice,
		},
	}

	resp, err := es.client.Search().Index(es.index).Request(req).Do(ctx)
	if err != nil {
		return 0, nil, fmt.Errorf("request failed: %w", err)
	}

	Gigs := make([]*gig.SellerGig, 0)

	for _, hit := range resp.Hits.Hits {
		var gig gig.SellerGig

		// need to convert the data into json.
		b, _ := hit.Source_.MarshalJSON()
		if err != nil {
			return 0, nil, fmt.Errorf("marshal failed %w", err)
		}

		// convert the json data back into a struct.
		err = json.Unmarshal(b, &gig)
		if err != nil {
			return 0, nil, fmt.Errorf("unmarshal failed into SellerGig %w", err)
		}

		Gigs = append(Gigs, &gig)
	}

	return resp.Hits.Total.Value, Gigs, nil
}

func (es *ESSearch) SearchSimiliarGigs(ctx context.Context, id string) (int64, []*gig.SellerGig, error) {
	req := essearch.NewRequest()

	req.Query = &types.Query{
		MoreLikeThis: &types.MoreLikeThisQuery{
			Fields: []string{"username", "title", "description", "basicDescription", "basicTitle", "categories", "subCategories", "tags"},
			Like: []types.Like{
				types.LikeDocument{Index_: &es.index, Id_: &id},
			},
		},
	}

	req.Size = utils.PtrI(5)

	resp, err := es.client.Search().Index(es.index).Request(req).Do(ctx)
	if err != nil {
		return 0, nil, fmt.Errorf("request failed: %w", err)
	}

	Gigs := make([]*gig.SellerGig, 0)

	for _, hit := range resp.Hits.Hits {
		var gig gig.SellerGig

		// need to convert the data into json.
		b, _ := hit.Source_.MarshalJSON()
		if err != nil {
			return 0, nil, fmt.Errorf("marshal failed %w", err)
		}

		// convert the json data back into a struct.
		err = json.Unmarshal(b, &gig)
		if err != nil {
			return 0, nil, fmt.Errorf("unmarshal failed into SellerGig %w", err)
		}

		Gigs = append(Gigs, &gig)
	}

	return resp.Hits.Total.Value, Gigs, nil
}

func (es *ESSearch) SearchTopRatedGigsbyCategory(ctx context.Context, category string) (int64, []*gig.SellerGig, error) {
	req := essearch.NewRequest()

	// if empty string then use asterix as wildcard
	if category == "" {
		category = "*"
	}

	req.Query = &types.Query{
		Bool: &types.BoolQuery{
			Must: []types.Query{
				{
					QueryString: &types.QueryStringQuery{
						Fields: []string{"categories"},
						Query:  category,
					},
				},
			},
			Filter: []types.Query{
				{
					Script: &types.ScriptQuery{
						QueryName_: utils.Ptr("script"),
						Script: types.InlineScript{
							Lang:   &scriptlanguage.Painless,
							Source: "doc['ratingSum'].value != 0 && (doc['ratingSum'].value / doc['ratingsCount'].value == params['threshold'])",
							Params: map[string]json.RawMessage{
								"threshold": json.RawMessage(`5`),
							},
						},
					},
				},
			},
		},
	}

	req.Size = utils.PtrI(10)

	resp, err := es.client.Search().Index(es.index).Request(req).Do(ctx)
	if err != nil {
		return 0, nil, fmt.Errorf("request failed: %w", err)
	}

	Gigs := make([]*gig.SellerGig, 0)

	for _, hit := range resp.Hits.Hits {
		var gig gig.SellerGig

		// need to convert the data into json.
		b, _ := hit.Source_.MarshalJSON()
		if err != nil {
			return 0, nil, fmt.Errorf("marshal failed %w", err)
		}

		// convert the json data back into a struct.
		err = json.Unmarshal(b, &gig)
		if err != nil {
			return 0, nil, fmt.Errorf("unmarshal failed into SellerGig %w", err)
		}

		Gigs = append(Gigs, &gig)
	}

	return resp.Hits.Total.Value, Gigs, nil
}
