package rest

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/maddsua/flippercardapp/auth"
	db_pkg "github.com/maddsua/flippercardapp/db"
	db_gen "github.com/maddsua/flippercardapp/db/generated"
	"github.com/maddsua/flippercardapp/db/types"
	"github.com/maddsua/flippercardapp/rest/model"
)

type resolver struct {
	db *db_pkg.Wrapper
}

func (rslv *resolver) LoadCardDeck(ctx context.Context, id uuid.UUID) (*model.CardDeck, error) {

	deck, err := rslv.db.GetDeckById(ctx, id)
	if db_pkg.IsNull(err) {
		return nil, &model.Error{Message: "deck not found", Code: http.StatusNotFound}
	} else if err != nil {
		return nil, InternalError("sqlc.GetDeckById", err)
	}

	cards, err := rslv.db.GetDeckCards(ctx, deck.ID)
	if err != nil {
		return nil, InternalError("sqlc.GetDeckCards", err)
	}

	result := model.CardDeck{
		CardDeckMetadata: db_pkg.TransformRow[model.CardDeckMetadata](deck),
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
		result.Cards[idx].FromRow(val)
	}

	return &result, nil
}

func (rslv *resolver) ListCardDeckPage(ctx context.Context, ids UUIDSet, page PagePointers) (*Page[model.CardDeckMetadata], error) {

	if idList := ids.List(); len(idList) > 0 {

		var entries []model.CardDeckMetadata

		for _, id := range idList.WithPage(page) {

			nextEntry, err := rslv.db.GetDecksBatch(ctx, db_gen.GetDecksBatchParams{
				ID:    types.NewNullUUID(id),
				Limit: 1,
			})

			if err != nil {
				return nil, InternalError("sqlc.GetDecksBatch", err)
			} else if len(nextEntry) == 0 {
				continue
			}

			entries = append(entries, db_pkg.TransformBatchRow[model.CardDeckMetadata](nextEntry[0]))
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

	return TransformPage(page, entries, db_pkg.TransformBatchRow[model.CardDeckMetadata, db_gen.GetDecksBatchRow]), nil
}

func (rslv *resolver) LoadCollection(ctx context.Context, id uuid.UUID) (*model.Collection, error) {

	collection, err := rslv.db.GetCollectionById(ctx, id)
	if db_pkg.IsNull(err) {
		return nil, &model.Error{Message: "collection not found", Code: http.StatusNotFound}
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
		CollectionMetadata: db_pkg.TransformRow[model.CollectionMetadata](collection),
		Decks:              make([]model.CardDeckMetadata, len(decks)),
	}

	result.CollectionMetadata.Size = len(decks)

	for idx, val := range decks {
		result.Decks[idx].FromBatchRow(val)
	}

	return &result, nil
}

func (rslv *resolver) ListCollectionsPage(ctx context.Context, ids UUIDSet, page PagePointers) (*Page[model.CollectionMetadata], error) {

	if idList := ids.List(); len(idList) > 0 {

		var entries []model.CollectionMetadata

		for _, id := range idList.WithPage(page) {

			nextEntry, err := rslv.db.GetCollectionBatch(ctx, db_gen.GetCollectionBatchParams{
				ID:    types.NewNullUUID(id),
				Limit: 1,
			})

			if err != nil {
				return nil, InternalError("sqlc.GetCollectionBatch", err)
			} else if len(nextEntry) == 0 {
				continue
			}

			entries = append(entries, db_pkg.TransformBatchRow[model.CollectionMetadata](nextEntry[0]))
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

	return TransformPage(page, entries, db_pkg.TransformBatchRow[model.CollectionMetadata, db_gen.GetCollectionBatchRow]), nil
}

func (rslv *resolver) SearchCollections(ctx context.Context, term string, page PagePointers) (*Page[model.CollectionSearchResult], error) {

	matched, err := rslv.matchFuzzyCollections(ctx, term, page)
	if err != nil {
		return nil, err
	}

	var entries []model.CollectionSearchResult

	for _, item := range matched {

		nextEntry, err := rslv.db.GetCollectionBatch(ctx, db_gen.GetCollectionBatchParams{
			ID:    types.NewNullUUID(item.id),
			Limit: 1,
		})

		if err != nil {
			return nil, InternalError("sqlc.GetCollectionBatch", err)
		} else if len(nextEntry) == 0 {
			continue
		}

		entries = append(entries, db_pkg.TransformBatchRow[model.CollectionSearchResult](nextEntry[0]))
	}

	return WrapPage(page, entries), nil
}

func (rslv *resolver) matchFuzzyCollections(ctx context.Context, term string, page PagePointers) ([]rankedSearchEntry, error) {

	if term = strings.ToLower(strings.TrimSpace(term)); len(term) < 2 {
		return nil, &model.Error{Message: "Search term too short"}
	} else if len(term) > math.MaxUint8 {
		return nil, &model.Error{Message: "Search term too long", Code: http.StatusRequestEntityTooLarge}
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

func (rslv *resolver) CreateContentCollection(ctx context.Context, params model.CollectionPatch) (*model.CollectionMetadata, error) {

	if perms, err := auth.For(ctx).Permissions(); err != nil {
		return nil, err
	} else if err := perms.CanEditContent(); err != nil {
		return nil, err
	}

	tx, err := rslv.db.BeginTx(ctx)
	if err != nil {
		return nil, InternalError("sqlc.BeginTx", err)
	}
	defer tx.Rollback()

	if err := params.Valid(); err != nil {
		return nil, &model.Error{Message: fmt.Sprintf("Invalid collection details: %v", err)}
	}

	if exists, err := tx.CollectionNameExists(ctx, params.Name); err != nil {
		return nil, InternalError("sqlc.CollectionNameExists", err)
	} else if exists > 0 {
		return nil, &model.Error{
			Message: "Collection name already exists",
			Code:    http.StatusConflict,
		}
	}

	entry, err := tx.InsertCollection(ctx, db_gen.InsertCollectionParams{
		ID:          uuid.New(),
		CreatedAt:   types.NewTime(time.Now()),
		UpdatedAt:   types.NewTime(time.Now()),
		Name:        params.Name,
		Description: types.NewNullString(params.Description),
	})

	if err != nil {
		return nil, InternalError("sqlc.InsertCollection", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, InternalError("sqlc.Commit", err)
	}

	result := db_pkg.TransformRow[model.CollectionMetadata](entry)
	return &result, nil
}

func (rslv *resolver) UpdateContentCollection(ctx context.Context, id uuid.UUID, patch model.CollectionPatch) (*model.CollectionMetadata, error) {

	if perms, err := auth.For(ctx).Permissions(); err != nil {
		return nil, err
	} else if err := perms.CanEditContent(); err != nil {
		return nil, err
	}

	if err := patch.Valid(); err != nil {
		return nil, &model.Error{Message: fmt.Sprintf("Invalid collection details: %v", err)}
	}

	entry, err := rslv.db.UpdateCollection(ctx, db_gen.UpdateCollectionParams{
		ID:          id,
		Name:        patch.Name,
		Description: types.NewNullString(patch.Description),
		UpdatedAt:   types.NewTime(time.Now()),
	})
	if db_pkg.IsNull(err) {
		return nil, &model.Error{Message: "Collection not found", Code: http.StatusNotFound}
	} else if err != nil {
		return nil, InternalError("sqlc.UpdateCollection", err)
	}

	result := db_pkg.TransformRow[model.CollectionMetadata](entry)
	return &result, nil
}

func (rslv *resolver) DeleteCollection(ctx context.Context, id uuid.UUID) error {

	if perms, err := auth.For(ctx).Permissions(); err != nil {
		return err
	} else if err := perms.CanEditContent(); err != nil {
		return err
	}

	tx, err := rslv.db.BeginTx(ctx)
	if err != nil {
		return InternalError("sqlc.BeginTx", err)
	}
	defer tx.Rollback()

	if count, err := tx.CollectionSize(ctx, id); err != nil {

		if db_pkg.IsNull(err) {
			return &model.Error{Message: "Collection not found", Code: http.StatusNotFound}
		}

		return InternalError("sqlc.CollectionSize", err)

	} else if count > 0 {
		return &model.Error{Message: "Collection is not empty", Code: http.StatusConflict}
	}

	if count, err := tx.DeleteCollection(ctx, id); err != nil {
		return InternalError("sqlc.DeleteCollection", err)
	} else if count == 0 {
		return &model.Error{Message: "Collectcion not found", Code: http.StatusNotFound}
	}

	if err := tx.Commit(); err != nil {
		return InternalError("sqlc.Commit", err)
	}

	return nil
}

func (rslv *resolver) CreateCardDeck(ctx context.Context, params model.CardDeckPatch) (*model.CardDeckMetadata, error) {

	if perms, err := auth.For(ctx).Permissions(); err != nil {
		return nil, err
	} else if err := perms.CanEditContent(); err != nil {
		return nil, err
	}

	if params.Details == nil {
		return nil, &model.Error{Message: "Deck details must be provided"}
	} else if err := params.Details.Valid(); err != nil {
		return nil, &model.Error{Message: fmt.Sprintf("Invalid deck details: %v", err)}
	} else if params.Content == nil || len(params.Content.Cards) == 0 {
		return nil, &model.Error{Message: "Deck has no cards in it"}
	} else if !params.CollectionID.Valid {
		return nil, &model.Error{Message: "Collection ID must be provided"}
	}

	tx, err := rslv.db.BeginTx(ctx)
	if err != nil {
		return nil, InternalError("sqlc.BeginTx", err)
	}
	defer tx.Rollback()

	if count, err := tx.CollectionIDExists(ctx, params.CollectionID.UUID); err != nil {
		return nil, InternalError("sqlc.CollectionIDExists", err)
	} else if count == 0 {
		return nil, &model.Error{Message: "Collection ID not found", Code: http.StatusNotFound}
	}

	deck, err := tx.InsertDeck(ctx, db_gen.InsertDeckParams{
		ID:           uuid.New(),
		CollectionID: params.CollectionID.UUID,
		CreatedAt:    types.NewTime(time.Now()),
		UpdatedAt:    types.NewTime(time.Now()),
		Name:         params.Details.Name,
		Description:  types.NewNullString(params.Details.Description),
	})

	if err != nil {
		return nil, InternalError("sqlc.InsertDeck", err)
	}

	for _, card := range params.Content.Cards {
		if err := tx.InsertCard(ctx, db_gen.InsertCardParams{
			Content:   card.CardNodeContent,
			ID:        uuid.New(),
			DeckID:    deck.ID,
			CreatedAt: types.NewTime(time.Now()),
			UpdatedAt: types.NewTime(time.Now()),
		}); err != nil {
			return nil, InternalError("sqlc.InsertCard", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, InternalError("sqlc.Commit", err)
	}

	result := db_pkg.TransformRow[model.CardDeckMetadata](deck)
	result.Size = len(params.Content.Cards)
	return &result, nil
}

func (rslv *resolver) UpdateCardDeck(ctx context.Context, id uuid.UUID, patch model.CardDeckPatch) (*model.CardDeckMetadata, error) {

	if perms, err := auth.For(ctx).Permissions(); err != nil {
		return nil, err
	} else if err := perms.CanEditContent(); err != nil {
		return nil, err
	}

	tx, err := rslv.db.BeginTx(ctx)
	if err != nil {
		return nil, InternalError("sqlc.BeginTx", err)
	}
	defer tx.Rollback()

	deck, err := tx.SetDeckUpdateTime(ctx, db_gen.SetDeckUpdateTimeParams{
		ID:        id,
		UpdatedAt: types.NewTime(time.Now()),
	})

	if db_pkg.IsNull(err) {
		return nil, &model.Error{Message: "Deck ID not found", Code: http.StatusNotFound}
	} else if err != nil {
		return nil, InternalError("sqlc.GetDeckById", err)
	}

	if patch.CollectionID.Valid {
		if count, err := tx.CollectionIDExists(ctx, patch.CollectionID.UUID); err != nil {
			return nil, InternalError("sqlc.CollectionIDExists", err)
		} else if count == 0 {
			return nil, &model.Error{Message: "Collection ID not found", Code: http.StatusNotFound}
		}
	}

	if patch.Details != nil {

		if err := patch.Details.Valid(); err != nil {
			return nil, &model.Error{Message: fmt.Sprintf("Invalid deck details: %v", err)}
		}

		if deck, err = tx.UpdateDeckMetadata(ctx, db_gen.UpdateDeckMetadataParams{
			ID:           id,
			CollectionID: patch.CollectionID,
			Name:         patch.Details.Name,
			Description:  types.NewNullString(patch.Details.Description),
		}); err != nil {
			return nil, InternalError("sqlc.UpdateDeckMetadata", err)
		}
	}

	result := db_pkg.TransformRow[model.CardDeckMetadata](deck)

	if patch.Content != nil {

		staleCards := map[uuid.UUID]struct{}{}
		if cardIds, err := tx.DeckCardSet(ctx, deck.ID); err != nil {
			return nil, InternalError("sqlc.DeckCardSet", err)
		} else {
			for _, id := range cardIds {
				staleCards[id] = struct{}{}
			}
		}

		for _, card := range patch.Content.Cards {

			if card.ID.Valid {
				if count, err := tx.UpdateCardContent(ctx, db_gen.UpdateCardContentParams{
					Content:   card.CardNodeContent,
					ID:        card.ID.UUID,
					DeckID:    deck.ID,
					UpdatedAt: types.NewTime(time.Now()),
				}); err != nil {
					return nil, InternalError("sqlc.UpdateCardContent", err)
				} else if count > 0 {
					delete(staleCards, card.ID.UUID)
					continue
				}
			}

			if err := tx.InsertCard(ctx, db_gen.InsertCardParams{
				Content:   card.CardNodeContent,
				ID:        uuid.New(),
				DeckID:    deck.ID,
				CreatedAt: types.NewTime(time.Now()),
				UpdatedAt: types.NewTime(time.Now()),
			}); err != nil {
				return nil, InternalError("sqlc.InsertCard", err)
			}
		}

		for id := range staleCards {
			if err := tx.DeleteCard(ctx, id); err != nil {
				return nil, InternalError("sqlc.DeleteCard", err)
			}
		}

		result.Size = len(patch.Content.Cards)
	}

	if err := tx.Commit(); err != nil {
		return nil, InternalError("sqlc.Commit", err)
	}

	return &result, nil
}

func (rslv *resolver) DeleteDeck(ctx context.Context, id uuid.UUID) error {

	if perms, err := auth.For(ctx).Permissions(); err != nil {
		return err
	} else if err := perms.CanEditContent(); err != nil {
		return err
	}

	if count, err := rslv.db.DeleteDeck(ctx, id); err != nil {
		return InternalError("sqlc.DeleteDeck", err)
	} else if count == 0 {
		return &model.Error{Message: "Deck not found", Code: http.StatusNotFound}
	}

	return nil
}

func (rslv *resolver) ExportCollectionBundle(ctx context.Context, id uuid.UUID) (*model.CollectionBundle, error) {

	if perms, err := auth.For(ctx).Permissions(); err != nil {
		return nil, err
	} else if err := perms.IsTeamMember(); err != nil {
		return nil, err
	}

	tx, err := rslv.db.BeginTx(ctx)
	if err != nil {
		return nil, InternalError("sqlc.BeginTx", err)
	}
	defer tx.Rollback()

	collection, err := tx.GetCollectionById(ctx, id)
	if db_pkg.IsNull(err) {
		return nil, &model.Error{Message: "collection not found", Code: http.StatusNotFound}
	} else if err != nil {
		return nil, InternalError("sqlc.GetCollectionById", err)
	}

	decks, err := tx.GetDecksBatch(ctx, db_gen.GetDecksBatchParams{
		CollectionID: types.NewNullUUID(collection.ID),
		Limit:        math.MaxInt,
	})

	if err != nil {
		return nil, InternalError("sqlc.GetDecksBatch", err)
	}

	bundle := model.CollectionBundle{
		CollectionMetadata: db_pkg.TransformRow[model.CollectionMetadata](collection),
		Decks:              make([]model.CardDeckBundle, len(decks)),
	}

	for idx, deck := range decks {

		cards, err := tx.GetDeckCards(ctx, deck.ID)
		if err != nil {
			return nil, InternalError("sqlc.GetDeckCards", err)
		}

		deckBundle := model.CardDeckBundle{
			CardDeckMetadata: db_pkg.TransformBatchRow[model.CardDeckMetadata](deck),
			Cards:            make([]model.Card, len(cards)),
		}

		for idx, val := range cards {
			deckBundle.Cards[idx].FromRow(val)
		}

		bundle.Decks[idx] = deckBundle
	}

	return &bundle, nil
}

func (rslv *resolver) ImportCollectionBundle(ctx context.Context, bundle *model.CollectionBundle) (*model.CollectionMetadata, error) {

	if perms, err := auth.For(ctx).Permissions(); err != nil {
		return nil, err
	} else if err := perms.CanEditContent(); err != nil {
		return nil, err
	}

	if err := bundle.Valid(); err != nil {
		return nil, &model.Error{Message: fmt.Sprintf("Invalid collection bundle: %v", err)}
	}

	tx, err := rslv.db.BeginTx(ctx)
	if err != nil {
		return nil, InternalError("sqlc.BeginTx", err)
	}
	defer tx.Rollback()

	collectionEntry, err := tx.InsertCollection(ctx, db_gen.InsertCollectionParams{
		ID:          uuid.New(),
		CreatedAt:   types.NewTime(time.Now()),
		UpdatedAt:   types.NewTime(time.Now()),
		Name:        fmt.Sprintf("%v - import %d", bundle.Name, time.Now().Unix()),
		Description: types.NewNullString(bundle.Description),
	})

	if err != nil {
		return nil, InternalError("sqlc.InsertCollection", err)
	}

	for _, deck := range bundle.Decks {

		deckEntry, err := tx.InsertDeck(ctx, db_gen.InsertDeckParams{
			ID:           uuid.New(),
			CollectionID: collectionEntry.ID,
			CreatedAt:    types.NewTime(deck.Created),
			UpdatedAt:    types.NewTime(deck.Updated),
			Name:         deck.Name,
			Description:  types.NewNullString(deck.Description),
		})

		if err != nil {
			return nil, InternalError("sqlc.InsertDeck", err)
		}

		for _, card := range deck.Cards {
			if err := tx.InsertCard(ctx, db_gen.InsertCardParams{
				Content:   card.CardNodeContent,
				ID:        uuid.New(),
				DeckID:    deckEntry.ID,
				CreatedAt: types.NewTime(card.Created),
				UpdatedAt: types.NewTime(card.Updated),
			}); err != nil {
				return nil, InternalError("sqlc.InsertCard", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, InternalError("sqlc.Commit", err)
	}

	result := db_pkg.TransformRow[model.CollectionMetadata](collectionEntry)
	result.Size = len(bundle.Decks)

	return &result, nil
}
