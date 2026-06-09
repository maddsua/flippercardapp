package rest

import (
	"database/sql"
	"io"
	"net/http"
	"strings"

	"github.com/maddsua/flippercardapp/auth"
	db_pkg "github.com/maddsua/flippercardapp/db"
	"github.com/maddsua/flippercardapp/rest/model"
)

func NewHandler(dbconn *sql.DB) http.Handler {

	db := db_pkg.NewWrapper(dbconn)
	rslv := resolver{db: db}

	mux := http.NewServeMux()

	mux.Handle("GET /auth/whoami", MethodHandleFunc(func(req *http.Request) (*auth.RequestAuth, error) {
		return auth.For(req.Context()), nil
	}))

	mux.Handle("POST /auth/signin", MethodHandleFunc(func(req *http.Request) (*auth.RequestAuth, error) {
		params, err := ParseGenericJSON[model.SignInParams](req)
		if err != nil {
			return nil, err
		}
		return auth.NewWebSessionWithPassword(req.Context(), params.Username, params.Password)
	}))

	mux.Handle("POST /auth/signout", MethodHandleFunc(func(req *http.Request) (*auth.RequestAuth, error) {
		return auth.TerminateWebSession(req.Context())
	}))

	mux.Handle("GET /collections", MethodHandleFunc(func(req *http.Request) (*Page[model.CollectionMetadata], error) {
		idSet, err := ParseUUIDSet(req.URL.Query().Get("ids"))
		if err != nil {
			return nil, err
		}
		return rslv.ListCollectionsBatch(req.Context(), idSet, Pagination(req))
	}))

	mux.Handle("GET /collections/search", MethodHandleFunc(func(req *http.Request) (*Page[model.CollectionSearchResult], error) {
		term := strings.TrimSpace(req.URL.Query().Get("term"))
		if term == "" {
			return nil, &model.Error{Message: "Search term cannot be empty"}
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

	mux.Handle("PUT /collections/new", MethodHandleFunc(func(req *http.Request) (*model.CollectionMetadata, error) {
		params, err := ParseGenericJSON[model.CollectionPatch](req)
		if err != nil {
			return nil, err
		}
		return rslv.CreateContentCollection(req.Context(), params)
	}))

	mux.Handle("PATCH /collections/{id}", MethodHandleFunc(func(req *http.Request) (*model.CollectionMetadata, error) {

		collectionID, err := ParseUUID(req.PathValue("id"))
		if err != nil {
			return nil, err
		}

		params, err := ParseGenericJSON[model.CollectionPatch](req)
		if err != nil {
			return nil, err
		}

		return rslv.UpdateContentCollection(req.Context(), collectionID, params)
	}))

	mux.Handle("DELETE /collections/{id}", MethodHandleFunc(func(req *http.Request) (*any, error) {
		collectionID, err := ParseUUID(req.PathValue("id"))
		if err != nil {
			return nil, err
		}
		recursive := strings.EqualFold(req.URL.Query().Get("recursive"), "true")
		return nil, rslv.DeleteCollection(req.Context(), collectionID, recursive)
	}))

	mux.Handle("GET /decks", MethodHandleFunc(func(req *http.Request) (*Page[model.CardDeckMetadata], error) {
		idSet, err := ParseUUIDSet(req.URL.Query().Get("ids"))
		if err != nil {
			return nil, err
		}
		return rslv.ListCardDeckBatch(req.Context(), idSet, Pagination(req))
	}))

	mux.Handle("GET /decks/{id}", MethodHandleFunc(func(req *http.Request) (*model.CardDeck, error) {
		deckID, err := ParseUUID(req.PathValue("id"))
		if err != nil {
			return nil, err
		}
		return rslv.LoadCardDeck(req.Context(), deckID)
	}))

	mux.Handle("GET /decks/{id}/versions", MethodHandleFunc(func(req *http.Request) (*Page[model.CardDeckVersionMetadata], error) {
		deckID, err := ParseUUID(req.PathValue("id"))
		if err != nil {
			return nil, err
		}
		return rslv.ListCardDeckVersions(req.Context(), deckID, Pagination(req))
	}))

	mux.Handle("GET /decks/{id}/version/{vid}", MethodHandleFunc(func(req *http.Request) (*model.CardDeckVersion, error) {
		deckID, err := ParseUUID(req.PathValue("id"))
		if err != nil {
			return nil, err
		}
		versionID, err := ParseUUID(req.PathValue("vid"))
		if err != nil {
			return nil, err
		}
		return rslv.LoadCardDeckVersion(req.Context(), deckID, versionID)
	}))

	mux.Handle("POST /decks/{id}/version/{vid}/rollback", MethodHandleFunc(func(req *http.Request) (*model.CardDeckMetadata, error) {
		deckID, err := ParseUUID(req.PathValue("id"))
		if err != nil {
			return nil, err
		}
		versionID, err := ParseUUID(req.PathValue("vid"))
		if err != nil {
			return nil, err
		}
		return rslv.RollbackCardDeckVersion(req.Context(), deckID, versionID)
	}))

	mux.Handle("PUT /decks/new", MethodHandleFunc(func(req *http.Request) (*model.CardDeckMetadata, error) {
		params, err := ParseGenericJSON[model.CardDeckPatch](req)
		if err != nil {
			return nil, err
		}
		return rslv.CreateCardDeck(req.Context(), params)
	}))

	mux.Handle("PATCH /decks/{id}", MethodHandleFunc(func(req *http.Request) (*model.CardDeckMetadata, error) {

		deckID, err := ParseUUID(req.PathValue("id"))
		if err != nil {
			return nil, err
		}

		params, err := ParseGenericJSON[model.CardDeckPatch](req)
		if err != nil {
			return nil, err
		}

		return rslv.UpdateCardDeck(req.Context(), deckID, params)
	}))

	mux.Handle("DELETE /decks/{id}", MethodHandleFunc(func(req *http.Request) (*any, error) {
		deckID, err := ParseUUID(req.PathValue("id"))
		if err != nil {
			return nil, err
		}
		return nil, rslv.DeleteDeck(req.Context(), deckID)
	}))

	mux.Handle("PUT /images/upload", MethodHandleFunc(func(req *http.Request) (*model.ImageMetadata, error) {

		if req.ContentLength <= 0 {
			return nil, &model.Error{Message: "Image uploads must be of a known size", Code: http.StatusLengthRequired}
		} else if req.ContentLength > 1_200_000 {
			return nil, &model.Error{Message: "Image upload size limited to 1.2MB", Code: http.StatusRequestEntityTooLarge}
		}

		return rslv.UploadImage(req.Context(), req.URL.Query().Get("name"), req.Body)
	}))

	mux.Handle("GET /images/{id}/metadata", MethodHandleFunc(func(req *http.Request) (*model.ImageMetadata, error) {
		return rslv.ImageMetadata(req.Context(), req.PathValue("id"))
	}))

	mux.Handle("GET /images/{id}/blob", http.HandlerFunc(func(wrt http.ResponseWriter, req *http.Request) {

		blob, err := rslv.ImageBlob(req.Context(), req.PathValue("id"))
		if err != nil {
			NewErrorResponseStatus[any](err, http.StatusBadRequest).Write(wrt)
			return
		}

		wrt.Header().Set("Content-Type", "application/octet-stream")

		io.Copy(wrt, blob)
	}))

	return auth.Middleware(db, mux)
}
