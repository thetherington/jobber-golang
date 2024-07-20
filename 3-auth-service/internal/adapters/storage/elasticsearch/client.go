package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	essearch "github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/thetherington/jobber-auth/internal/adapters/config"
	"github.com/thetherington/jobber-common/models/gig"
	"github.com/thetherington/jobber-common/models/search"
	"github.com/thetherington/jobber-common/utils"
	"go.elastic.co/apm/module/apmelasticsearch/v2"
)

/**
 * Elasticsearch implements port.AuthSearch interface
 * and provides an access to the elasticsearch library
 */
type ESSearch struct {
	client *elasticsearch.TypedClient
	index  string
}

// New creates a new instance of Elasticsearch
func New(config *config.Elastic, index string) (*ESSearch, error) {
	typedClient, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{config.SearchUrl},
		Transport: apmelasticsearch.WrapRoundTripper(http.DefaultTransport),
	})
	if err != nil {
		return nil, err
	}

	return &ESSearch{typedClient, index}, nil
}

func (es *ESSearch) CheckConnection(ctx context.Context) {
	for {
		resp, err := es.client.Cluster.Health().Do(ctx)
		if err != nil {
			slog.With("error", err).Error("AuthService ES CheckConnection()")
			slog.Warn("Connection to Elasticsearch failed. Retrying...")

			time.Sleep(5 * time.Second)
			continue
		}

		slog.Info("Elasticsearch health status", "status", resp.Status.String())
		break
	}
}

func (es *ESSearch) CreateIndex(ctx context.Context) {
	exists, err := es.client.Indices.Exists(es.index).Do(ctx)
	if err != nil {
		slog.With("error", err).Error("AuthService ES Indices.Exists()")
		return
	}

	if exists {
		slog.Info("Index already exists", "index", es.index)
		return
	}

	_, err = es.client.Indices.Create(es.index).Request(&create.Request{
		Mappings: &types.TypeMapping{
			Properties: map[string]types.Property{
				"price":        types.NewDoubleNumberProperty(),
				"ratingsCount": types.NewDoubleNumberProperty(),
				"ratingSum":    types.NewDoubleNumberProperty(),
			},
		},
	}).Do(ctx)
	if err != nil {
		slog.With("error", err).Error("AuthService ES Indices.Create()")
		return
	}

	es.client.Indices.Refresh().Do(ctx)
	slog.Info("Created index", "index", es.index)
}

func (es *ESSearch) GetDocumentById(ctx context.Context, id string) (*gig.SellerGig, error) {
	resp, err := es.client.Get(es.index, id).Do(ctx)

	if resp == nil && err == nil {
		return nil, ErrIndexNotExist
	}
	if err != nil {
		return nil, err
	}
	if !resp.Found {
		return nil, ErrGigNotFound
	}

	b, err := resp.Source_.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var gig gig.SellerGig

	err = json.Unmarshal(b, &gig)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response into SellerGig %w", err)
	}

	return &gig, nil
}

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
