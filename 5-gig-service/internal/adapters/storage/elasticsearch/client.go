package elasticsearch

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/thetherington/jobber-gig/internal/adapters/config"
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
