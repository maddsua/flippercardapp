package rest

import (
	"context"
	"math"
	"net/http"

	"github.com/google/uuid"
	db_pkg "github.com/maddsua/flippercardapp/db"
	db_gen "github.com/maddsua/flippercardapp/db/generated"
	"github.com/maddsua/flippercardapp/db/types"
	"github.com/maddsua/flippercardapp/rest/model"
	"github.com/maddsua/flippercardapp/rest/model/transform"
)

type resolver struct {
	db *db_pkg.Wrapper
}

func (rslv *resolver) LoadCardDeck(ctx context.Context, id uuid.UUID) (*model.CardDeck, error) {

	deck, err := rslv.db.GetDeckById(ctx, id)
	if db_pkg.IsNull(err) {
		return nil, &APIError{Message: "deck not found", Code: http.StatusNotFound}
	} else if err != nil {
		return nil, InternalError("sqlc.GetDeckById", err)
	}

	cards, err := rslv.db.GetDeckCards(ctx, deck.ID)
	if err != nil {
		return nil, InternalError("sqlc.GetDeckCards", err)
	}

	result := model.CardDeck{
		CardDeckMetadata: transform.CardDeckMetadataFromRow(deck),
		Labels:           []string{deck.Name},
		Cards:            make([]model.Card, len(cards)),
	}

	result.CardDeckMetadata.Size = len(cards)

	if collection, err := rslv.db.GetCollectionById(ctx, deck.CollectionID); err != nil {
		return nil, InternalError("sqlc.GetCollectionById", err)
	} else {
		result.Labels = append(result.Labels, collection.Name)
	}

	for idx, val := range cards {
		result.Cards[idx] = transform.CardFromRow(val)
	}

	return &result, nil
}

func (rslv *resolver) ListCardDeckPage(ctx context.Context, ids UUIDSet, page PagePointers) (*Page[model.CardDeckMetadata], error) {

	if !ids.IsEmpty() {

		var entries []db_gen.GetDecksBatchRow

		for _, id := range ids.WithPage(page) {

			next, err := rslv.db.GetDecksBatch(ctx, db_gen.GetDecksBatchParams{
				ID:    types.NewNullUUID(id),
				Limit: 1,
			})

			if err != nil {
				return nil, InternalError("sqlc.GetDecksBatch", err)
			}

			entries = append(entries, next...)
		}

		return Paginate(page, entries, transform.CardDeckMetadataFromBatchRow), nil
	}

	entries, err := rslv.db.GetDecksBatch(ctx, db_gen.GetDecksBatchParams{
		Limit:  page.QueryLimit(),
		Offset: page.QueryOffset(),
	})

	if err != nil {
		return nil, InternalError("sqlc.GetDecksBatch", err)
	}

	return Paginate(page, entries, transform.CardDeckMetadataFromBatchRow), nil
}

func (rslv *resolver) LoadCollection(ctx context.Context, id uuid.UUID) (*model.Collection, error) {

	collection, err := rslv.db.GetCollectionById(ctx, id)
	if db_pkg.IsNull(err) {
		return nil, &APIError{Message: "collection not found", Code: http.StatusNotFound}
	} else if err != nil {
		return nil, InternalError("sqlc.GetCollectionById", err)
	}

	decks, err := rslv.db.GetDecksBatch(ctx, db_gen.GetDecksBatchParams{
		CollectionID: types.NewNullUUID(collection.ID),
		Limit:        math.MaxInt32,
	})

	if err != nil {
		return nil, InternalError("sqlc.GetDecksBatch", err)
	}

	result := model.Collection{
		CollectionMetadata: transform.CollectionMetadataFromRow(collection),
		Decks:              make([]model.CardDeckMetadata, len(decks)),
	}

	result.CollectionMetadata.Size = len(decks)

	for idx, val := range decks {
		result.Decks[idx] = transform.CardDeckMetadataFromBatchRow(val)
	}

	return &result, nil
}

func (rslv *resolver) ListCollectionsPage(ctx context.Context, ids UUIDSet, page PagePointers) (*Page[model.CollectionMetadata], error) {

	if !ids.IsEmpty() {

		var entries []db_gen.GetCollectionBatchRow

		for _, id := range ids.WithPage(page) {

			next, err := rslv.db.GetCollectionBatch(ctx, db_gen.GetCollectionBatchParams{
				ID:    types.NewNullUUID(id),
				Limit: 1,
			})

			if err != nil {
				return nil, InternalError("sqlc.GetCollectionBatch", err)
			}

			entries = append(entries, next...)
		}

		return Paginate(page, entries, transform.CollectionMetadataFromBatchRow), nil
	}

	entries, err := rslv.db.GetCollectionBatch(ctx, db_gen.GetCollectionBatchParams{
		Limit:  page.QueryLimit(),
		Offset: page.QueryOffset(),
	})
	if err != nil {
		return nil, InternalError("sqlc.GetCollectionBatch", err)
	}

	return Paginate(page, entries, transform.CollectionMetadataFromBatchRow), nil
}
