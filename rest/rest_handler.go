package rest

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	db_pkg "github.com/maddsua/flippercardapp/db"
	db_gen "github.com/maddsua/flippercardapp/db/generated"
	"github.com/maddsua/flippercardapp/rest/model"
	"github.com/maddsua/flippercardapp/rest/model/transform"
)

func NewHandler(dbconn *sql.DB) http.Handler {

	rslv := resolver{db_pkg.NewWrapper(dbconn)}

	mux := http.NewServeMux()

	mux.Handle("GET /collections", MethodHandleFunc(func(req *http.Request) (*Page[model.CollectionMetadata], error) {
		return rslv.ListCollections(
			req.Context(),
			ReqParamUUID(req.URL.Query().Get("id")),
			Pagination(req),
		)
	}))

	mux.Handle("GET /collections/{id}", MethodHandleFunc(func(req *http.Request) (*model.Collection, error) {
		collectionID, err := ReqPathUUID(req.PathValue("id"))
		if err != nil {
			return nil, err
		}
		return rslv.LoadCollection(req.Context(), collectionID)
	}))

	mux.Handle("GET /decks", MethodHandleFunc(func(req *http.Request) (*Page[model.CardDeckMetadata], error) {
		return rslv.ListDecks(
			req.Context(),
			ReqParamUUID(req.URL.Query().Get("id")),
			ReqParamUUID(req.URL.Query().Get("collection_id")),
			Pagination(req),
		)
	}))

	mux.Handle("GET /decks/{id}", MethodHandleFunc(func(req *http.Request) (*model.CardDeck, error) {
		deckID, err := ReqPathUUID(req.PathValue("id"))
		if err != nil {
			return nil, err
		}
		return rslv.LoadDeck(req.Context(), deckID)
	}))

	return mux
}

type resolver struct {
	db *db_pkg.Wrapper
}

func (rslv *resolver) LoadDeck(ctx context.Context, id uuid.UUID) (*model.CardDeck, error) {

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
		Cards:            make([]model.Card, len(cards)),
	}

	result.CardDeckMetadata.Size = len(cards)

	for idx, val := range cards {
		result.Cards[idx] = transform.CardFromRow(val)
	}

	return &result, nil
}

func (rslv *resolver) ListDecks(ctx context.Context, id uuid.NullUUID, collectionID uuid.NullUUID, page PagePointers) (*Page[model.CardDeckMetadata], error) {

	entries, err := rslv.db.GetDecksBatch(ctx, db_gen.GetDecksBatchParams{
		ID:           id,
		CollectionID: collectionID,
		Limit:        page.QueryLimit(),
		Offset:       page.QueryOffset(),
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
		CollectionID: uuid.NullUUID{
			UUID:  collection.ID,
			Valid: true,
		},
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

func (rslv *resolver) ListCollections(ctx context.Context, id uuid.NullUUID, page PagePointers) (*Page[model.CollectionMetadata], error) {

	entries, err := rslv.db.GetCollectionBatch(ctx, db_gen.GetCollectionBatchParams{
		ID:     id,
		Limit:  page.QueryLimit(),
		Offset: page.QueryOffset(),
	})

	if err != nil {
		return nil, InternalError("sqlc.GetCollectionBatch", err)
	}

	return Paginate(page, entries, transform.CollectionMetadataFromBatchRow), nil
}
