package rest

import (
	"context"
	"math"
	"net/http"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/lithammer/fuzzysearch/fuzzy"
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

	if idList := ids.List(); len(idList) > 0 {

		var entries []model.CardDeckMetadata

		for _, id := range idList.WithPage(page) {

			next, err := rslv.db.GetDecksBatch(ctx, db_gen.GetDecksBatchParams{
				ID:    types.NewNullUUID(id),
				Limit: 1,
			})

			if err != nil {
				return nil, InternalError("sqlc.GetDecksBatch", err)
			} else if len(next) == 0 {
				continue
			}

			entries = append(entries, transform.CardDeckMetadataFromBatchRow(next[0]))
		}

		return WrapPage(page, entries), nil
	}

	entries, err := rslv.db.GetDecksBatch(ctx, db_gen.GetDecksBatchParams{
		Limit:  page.QueryLimit(),
		Offset: page.QueryOffset(),
	})

	if err != nil {
		return nil, InternalError("sqlc.GetDecksBatch", err)
	}

	return TransformPage(page, entries, transform.CardDeckMetadataFromBatchRow), nil
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

	if idList := ids.List(); len(idList) > 0 {

		var entries []model.CollectionMetadata

		for _, id := range idList.WithPage(page) {

			next, err := rslv.db.GetCollectionBatch(ctx, db_gen.GetCollectionBatchParams{
				ID:    types.NewNullUUID(id),
				Limit: 1,
			})

			if err != nil {
				return nil, InternalError("sqlc.GetCollectionBatch", err)
			} else if len(next) == 0 {
				continue
			}

			entries = append(entries, transform.CollectionMetadataFromBatchRow(next[0]))
		}

		return WrapPage(page, entries), nil
	}

	entries, err := rslv.db.GetCollectionBatch(ctx, db_gen.GetCollectionBatchParams{
		Limit:  page.QueryLimit(),
		Offset: page.QueryOffset(),
	})
	if err != nil {
		return nil, InternalError("sqlc.GetCollectionBatch", err)
	}

	return TransformPage(page, entries, transform.CollectionMetadataFromBatchRow), nil
}

func (rslv *resolver) SearchCollections(ctx context.Context, term string, page PagePointers) (*Page[model.CollectionSearchResult], error) {

	matched, err := rslv.matchFuzzyCollections(ctx, term, page)
	if err != nil {
		return nil, err
	}

	var entries []model.CollectionSearchResult

	for _, item := range matched {

		next, err := rslv.db.GetCollectionBatch(ctx, db_gen.GetCollectionBatchParams{
			ID:    types.NewNullUUID(item.id),
			Limit: 1,
		})

		if err != nil {
			return nil, InternalError("sqlc.GetCollectionBatch", err)
		} else if len(next) == 0 {
			continue
		}

		entries = append(entries, model.CollectionSearchResult{
			CollectionMetadata: transform.CollectionMetadataFromBatchRow(next[0]),
			Rank:               item.rank,
		})
	}

	return WrapPage(page, entries), nil
}

func (rslv *resolver) matchFuzzyCollections(ctx context.Context, term string, page PagePointers) ([]rankedSearchEntry, error) {

	if term = strings.ToLower(strings.TrimSpace(term)); len(term) < 2 {
		return nil, &APIError{Message: "Search term too short"}
	} else if len(term) > math.MaxUint8 {
		return nil, &APIError{Message: "Search term too long", Code: http.StatusRequestEntityTooLarge}
	}

	tx, err := rslv.db.BeginTx(ctx)
	if err != nil {
		return nil, InternalError("sqlc.BeginTx", err)
	}
	defer tx.Rollback()

	const indexBatchSize = 100

	index := map[uuid.UUID]int{}

	for offset := 0; offset < math.MaxInt; offset += indexBatchSize {

		next, err := tx.GetCollectionSearchBatch(ctx, db_gen.GetCollectionSearchBatchParams{
			Offset: int64(offset),
			Limit:  indexBatchSize,
		})

		if err != nil {
			return nil, InternalError("sqlc.GetCollectionSearchBatch", err)
		} else if len(next) == 0 {
			break
		}

		for _, item := range next {
			doc := strings.ToLower(item.Name)
			if rank := fuzzy.RankMatch(term, doc); rank >= 0 {
				index[item.ID] = rank
			}
		}
	}

	var indexSlice []rankedSearchEntry
	for id, rank := range index {
		indexSlice = append(indexSlice, rankedSearchEntry{id: id, rank: rank})
	}

	sort.SliceStable(indexSlice, func(i, j int) bool {
		return indexSlice[i].rank < indexSlice[j].rank
	})

	priority := make(UUIDList, len(indexSlice))
	for idx, item := range indexSlice {
		priority[idx] = item.id
	}

	return SlicePage(indexSlice, page), nil
}

type rankedSearchEntry struct {
	id   uuid.UUID
	rank int
}
