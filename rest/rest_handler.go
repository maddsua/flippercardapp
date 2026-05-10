package rest

import (
	"database/sql"
	"net/http"

	db_pkg "github.com/maddsua/flippercardapp/db"
	"github.com/maddsua/flippercardapp/rest/model"
)

func NewHandler(dbconn *sql.DB) http.Handler {

	rslv := resolver{db_pkg.NewWrapper(dbconn)}

	mux := http.NewServeMux()

	mux.Handle("GET /collections", MethodHandleFunc(func(req *http.Request) (*Page[model.CollectionMetadata], error) {
		return rslv.CollectionsPage(
			req.Context(),
			ParseUUIDSet(req.URL.Query().Get("ids")),
			Pagination(req),
		)
	}))

	mux.Handle("GET /collections/{id}", MethodHandleFunc(func(req *http.Request) (*model.Collection, error) {
		//	todo: refactor
		collectionID, err := ParseUUID(req.PathValue("id"))
		if err != nil {
			return nil, err
		}
		return rslv.LoadCollection(req.Context(), collectionID)
	}))

	mux.Handle("GET /decks", MethodHandleFunc(func(req *http.Request) (*Page[model.CardDeckMetadata], error) {
		return rslv.DecksPage(
			req.Context(),
			ParseUUIDSet(req.URL.Query().Get("ids")),
			ParseNullUUID(req.URL.Query().Get("collection_id")),
			Pagination(req),
		)
	}))

	mux.Handle("GET /decks/{id}", MethodHandleFunc(func(req *http.Request) (*model.CardDeck, error) {
		//	todo: refactor
		deckID, err := ParseUUID(req.PathValue("id"))
		if err != nil {
			return nil, err
		}
		return rslv.LoadDeck(req.Context(), deckID)
	}))

	return mux
}
