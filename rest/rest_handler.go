package rest

import (
	"compress/gzip"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/maddsua/flippercardapp/auth"
	db_pkg "github.com/maddsua/flippercardapp/db"
	"github.com/maddsua/flippercardapp/rest/model"
)

func NewHandler(dbconn *sql.DB) http.Handler {

	db := db_pkg.NewWrapper(dbconn)
	rslv := resolver{db: db}

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

	mux.Handle("PUT /manage/content/collection", MethodHandleFunc(func(req *http.Request) (*model.CollectionMetadata, error) {
		params, err := ParseGenericJSON[model.CollectionPatch](req)
		if err != nil {
			return nil, err
		}
		return rslv.CreateContentCollection(req.Context(), params)
	}))

	mux.Handle("POST /manage/content/collections/import", MethodHandleFunc(func(req *http.Request) (*model.CollectionMetadata, error) {

		var bundle model.CollectionBundle

		gz, err := gzip.NewReader(req.Body)
		if err != nil {
			return nil, &model.Error{Message: "Invalid bundle: " + err.Error()}
		}

		defer gz.Close()

		if err := json.NewDecoder(gz).Decode(&bundle); err != nil {
			return nil, &model.Error{Message: "Invalid bundle: " + err.Error()}
		}

		return rslv.ImportCollectionBundle(req.Context(), &bundle)
	}))

	mux.Handle("PATCH /manage/content/collection/{id}/metadata", MethodHandleFunc(func(req *http.Request) (*model.CollectionMetadata, error) {

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

	mux.Handle("POST /manage/content/collection/{id}/export", http.HandlerFunc(func(wrt http.ResponseWriter, req *http.Request) {

		collectionID, err := ParseUUID(req.PathValue("id"))
		if err != nil {
			NewErrorResponseStatus[any](err, http.StatusBadRequest).Write(wrt)
			return
		}

		bundle, err := rslv.ExportCollectionBundle(req.Context(), collectionID)
		if err != nil {
			NewErrorResponseStatus[any](err, http.StatusBadRequest).Write(wrt)
			return
		}

		wrt.Header().Set("Content-Type", "application/json+gzip")
		wrt.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", url.PathEscape(bundle.Name)))

		gz := gzip.NewWriter(wrt)
		defer gz.Close()

		if err := json.NewEncoder(gz).Encode(bundle); err != nil {
			slog.Error("REST: Encode and compress collection bundle",
				slog.String("type", "json"),
				slog.String("client", req.RemoteAddr),
				slog.String("err", err.Error()))
			return
		}

	}))

	mux.Handle("DELETE /manage/content/collection/{id}", MethodHandleFunc(func(req *http.Request) (*any, error) {
		collectionID, err := ParseUUID(req.PathValue("id"))
		if err != nil {
			return nil, err
		}
		recursive := strings.EqualFold(req.URL.Query().Get("recursive"), "true")
		return nil, rslv.DeleteCollection(req.Context(), collectionID, recursive)
	}))

	mux.Handle("PUT /manage/content/deck", MethodHandleFunc(func(req *http.Request) (*model.CardDeckMetadata, error) {
		params, err := ParseGenericJSON[model.CardDeckPatch](req)
		if err != nil {
			return nil, err
		}
		return rslv.CreateCardDeck(req.Context(), params)
	}))

	mux.Handle("PATCH /manage/content/deck/{id}", MethodHandleFunc(func(req *http.Request) (*model.CardDeckMetadata, error) {

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

	mux.Handle("DELETE /manage/content/deck/{id}", MethodHandleFunc(func(req *http.Request) (*any, error) {
		deckID, err := ParseUUID(req.PathValue("id"))
		if err != nil {
			return nil, err
		}
		return nil, rslv.DeleteDeck(req.Context(), deckID)
	}))

	mux.Handle("PUT /manage/content/images/upload", MethodHandleFunc(func(req *http.Request) (*model.ImageMetadata, error) {

		if req.ContentLength <= 0 {
			return nil, &model.Error{Message: "Image uploads must be of a known size", Code: http.StatusLengthRequired}
		} else if req.ContentLength > 1_200_000 {
			return nil, &model.Error{Message: "Image upload size limited to 1.2MB", Code: http.StatusRequestEntityTooLarge}
		}

		return rslv.UploadImage(req.Context(), req.URL.Query().Get("name"), req.Body)
	}))

	mux.Handle("GET /manage/content/images/{id}/metadata", MethodHandleFunc(func(req *http.Request) (*model.ImageMetadata, error) {
		return rslv.ImageMetadata(req.Context(), req.PathValue("id"))
	}))

	return auth.Middleware(db, mux)
}
