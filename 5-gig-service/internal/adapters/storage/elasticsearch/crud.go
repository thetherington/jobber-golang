package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/thetherington/jobber-common/models/gig"
)

func (es *ESSearch) InsertGig(ctx context.Context, newGig *gig.SellerGig) (string, error) {
	resp, err := es.client.Index(es.index).Id(newGig.ID).Request(newGig).Do(ctx)
	if err != nil {
		return "", err
	}

	return resp.Id_, nil
}

func (es *ESSearch) GetGigsCount(ctx context.Context) (int32, error) {
	result, err := es.client.Count().Index(es.index).Query(&types.Query{
		MatchAll: &types.MatchAllQuery{},
	}).Do(ctx)
	if err != nil {
		slog.With("error", err).Error("GetGigsCount error")
	}

	return int32(result.Count), nil
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

func (es *ESSearch) UpdateGig(ctx context.Context, id string, updateGig *gig.SellerGig) (string, error) {
	data, err := json.Marshal(updateGig)
	if err != nil {
		return "", err
	}

	resp, err := es.client.Update(es.index, id).Request(&update.Request{
		Doc: json.RawMessage(data),
	}).Do(ctx)
	if err != nil {
		return "", err
	}

	return resp.Id_, nil
}

func (es *ESSearch) DeleteGig(ctx context.Context, id string) error {
	_, err := es.client.Delete(es.index, id).Do(ctx)
	return err
}
