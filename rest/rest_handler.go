package rest

import (
	"database/sql"
	"net/http"
	"strings"

	db_pkg "github.com/maddsua/flippercardapp/db"
	"github.com/maddsua/flippercardapp/rest/model"
)

func NewHandler(dbconn *sql.DB) http.Handler {

	rslv := resolver{db_pkg.NewWrapper(dbconn)}

	mux := http.NewServeMux()

	mux.Handle("GET /collections", MethodHandleFunc(func(req *http.Request) (*Page[model.CollectionMetadata], error) {
		idSet, err := ParseUUIDSet(req.URL.Query().Get("ids"))
		if err != nil {
			return nil, err
		}
		return rslv.ListCollectionsPage(req.Context(), idSet, Pagination(req))
	}))

	mux.Handle("GET /collections/search", MethodHandleFunc(func(req *http.Request) (*Page[model.CollectionSearchResult], error) {
		term := strings.TrimSpace(req.URL.Query().Get("term"))
		if term == "" {
			return nil, &APIError{Message: "Search term cannot be empty"}
		}
		return rslv.SearchCollections(req.Context(), term, Pagination(req))
	}))

	mux.Handle("GET /collections/{id}", MethodHandleFunc(func(req *http.Request) (*model.Collection, error) {
		collectionID, err := ParseUUID(req.PathValue("id"))
		if err != nil {
			return nil, err
		}
		return rslv.LoadCollection(req.Context(), collectionID)
	}))

	mux.Handle("GET /decks", MethodHandleFunc(func(req *http.Request) (*Page[model.CardDeckMetadata], error) {
		idSet, err := ParseUUIDSet(req.URL.Query().Get("ids"))
		if err != nil {
			return nil, err
		}
		return rslv.ListCardDeckPage(req.Context(), idSet, Pagination(req))
	}))

	mux.Handle("GET /decks/{id}", MethodHandleFunc(func(req *http.Request) (*model.CardDeck, error) {
		deckID, err := ParseUUID(req.PathValue("id"))
		if err != nil {
			return nil, err
		}
		return rslv.LoadCardDeck(req.Context(), deckID)
	}))

	mux.Handle("GET /whoami", MethodHandleFunc(func(req *http.Request) (*AuthState, error) {
		return authStateFor(req.Context()), nil
	}))

	mux.Handle("PUT /manage/content/collection", MethodHandleFunc(func(req *http.Request) (*model.CollectionMetadata, error) {

		if state := authStateFor(req.Context()); state.Actor == nil {
			return nil, &APIError{Message: "Unauthorized", Code: http.StatusUnauthorized}
		} else if !state.Permissions.ContentEdit {
			return nil, &APIError{Message: "ContentEdit permission required", Code: http.StatusForbidden}
		}

		params, err := ParseGeneric[model.CollectionPatch](req)
		if err != nil {
			return nil, err
		}

		return rslv.CreateContentCollection(req.Context(), params)
	}))

	mux.Handle("PATCH /manage/content/collection/{id}/metadata", MethodHandleFunc(func(req *http.Request) (*model.CollectionMetadata, error) {

		if state := authStateFor(req.Context()); state.Actor == nil {
			return nil, &APIError{Message: "Unauthorized", Code: http.StatusUnauthorized}
		} else if !state.Permissions.ContentEdit {
			return nil, &APIError{Message: "ContentEdit permission required", Code: http.StatusForbidden}
		}

		collectionID, err := ParseUUID(req.PathValue("id"))
		if err != nil {
			return nil, err
		}

		params, err := ParseGeneric[model.CollectionPatch](req)
		if err != nil {
			return nil, err
		}

		return rslv.UpdateContentCollectionMetadata(req.Context(), collectionID, params)
	}))

	mux.Handle("DELETE /manage/content/collection/{id}", MethodHandleFunc(func(req *http.Request) (*any, error) {

		if state := authStateFor(req.Context()); state.Actor == nil {
			return nil, &APIError{Message: "Unauthorized", Code: http.StatusUnauthorized}
		} else if !state.Permissions.ContentEdit {
			return nil, &APIError{Message: "ContentEdit permission required", Code: http.StatusForbidden}
		}

		collectionID, err := ParseUUID(req.PathValue("id"))
		if err != nil {
			return nil, err
		}

		return nil, rslv.DeleteCollection(req.Context(), collectionID)
	}))

	mux.Handle("PUT /manage/content/deck", MethodHandleFunc(func(req *http.Request) (*model.CardDeckMetadata, error) {

		if state := authStateFor(req.Context()); state.Actor == nil {
			return nil, &APIError{Message: "Unauthorized", Code: http.StatusUnauthorized}
		} else if !state.Permissions.ContentEdit {
			return nil, &APIError{Message: "ContentEdit permission required", Code: http.StatusForbidden}
		}

		params, err := ParseGeneric[model.CardDeckPatch](req)
		if err != nil {
			return nil, err
		}

		return rslv.CreateCardDeck(req.Context(), params)
	}))

	mux.Handle("PATCH /manage/content/deck/{id}/metadata", MethodHandleFunc(func(req *http.Request) (*model.CardDeckMetadata, error) {

		if state := authStateFor(req.Context()); state.Actor == nil {
			return nil, &APIError{Message: "Unauthorized", Code: http.StatusUnauthorized}
		} else if !state.Permissions.ContentEdit {
			return nil, &APIError{Message: "ContentEdit permission required", Code: http.StatusForbidden}
		}

		deckID, err := ParseUUID(req.PathValue("id"))
		if err != nil {
			return nil, err
		}

		params, err := ParseGeneric[model.CardDeckMetadataPatch](req)
		if err != nil {
			return nil, err
		}

		return rslv.UpdateCardDeckMetadata(req.Context(), deckID, params)
	}))

	mux.Handle("PATCH /manage/content/deck/{id}/content", MethodHandleFunc(func(req *http.Request) (*model.CardDeckMetadata, error) {

		if state := authStateFor(req.Context()); state.Actor == nil {
			return nil, &APIError{Message: "Unauthorized", Code: http.StatusUnauthorized}
		} else if !state.Permissions.ContentEdit {
			return nil, &APIError{Message: "ContentEdit permission required", Code: http.StatusForbidden}
		}

		deckID, err := ParseUUID(req.PathValue("id"))
		if err != nil {
			return nil, err
		}

		params, err := ParseGeneric[model.CardDeckContentPatch](req)
		if err != nil {
			return nil, err
		}

		return rslv.UpdateCardDeckContent(req.Context(), deckID, params)
	}))

	mux.Handle("DELETE /manage/content/deck/{id}", MethodHandleFunc(func(req *http.Request) (*any, error) {

		if state := authStateFor(req.Context()); state.Actor == nil {
			return nil, &APIError{Message: "Unauthorized", Code: http.StatusUnauthorized}
		} else if !state.Permissions.ContentEdit {
			return nil, &APIError{Message: "ContentEdit permission required", Code: http.StatusForbidden}
		}

		deckID, err := ParseUUID(req.PathValue("id"))
		if err != nil {
			return nil, err
		}

		return nil, rslv.DeleteDeck(req.Context(), deckID)
	}))

	return authMiddleware(&StaticAuthProvider{}, mux)
}
