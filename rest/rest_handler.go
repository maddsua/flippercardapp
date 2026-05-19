package rest

import (
	"database/sql"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	db_pkg "github.com/maddsua/flippercardapp/db"
	db_gen "github.com/maddsua/flippercardapp/db/generated"
	"github.com/maddsua/flippercardapp/db/types"
	"github.com/maddsua/flippercardapp/rest/model"
	"github.com/maddsua/flippercardapp/utils"
	"golang.org/x/crypto/bcrypt"
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

	mux.Handle("GET /auth/whoami", MethodHandleFunc(func(req *http.Request) (*AuthState, error) {
		return authStateFor(req.Context()), nil
	}))

	mux.Handle("POST /auth/signin", MethodHandleFunc(func(req *http.Request) (*AuthState, error) {

		params, err := ParseGeneric[model.SignInParams](req)
		if err != nil {
			return nil, err
		}

		user, err := db.GetUserByName(req.Context(), params.Username)
		if db_pkg.IsNull(err) {
			return nil, &APIError{Message: "User not found", Code: http.StatusUnauthorized}
		} else if err != nil {
			return nil, InternalError("sqlc.GetUserByName", err)
		} else if strings.TrimSpace(params.Password) == "" {
			return nil, &APIError{Message: "Invalid password", Code: http.StatusUnauthorized}
		}

		if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(params.Password)); err != nil {

			if err != bcrypt.ErrMismatchedHashAndPassword {
				slog.Error("REST: bcrypt.CompareHashAndPassword",
					slog.String("err", err.Error()))
			}

			return nil, &APIError{Message: "Invalid password", Code: http.StatusUnauthorized}
		}

		expires := time.Now().Add(SessionTTL)

		sess, err := db.InsertSession(req.Context(), db_gen.InsertSessionParams{
			ID:        uuid.New(),
			CreatedAt: types.NewTime(time.Now()),
			ExpiresAt: types.NewTime(expires.Add(time.Minute)),
			UserID:    user.ID,
			Secret:    utils.NewRandomToken(128),
		})

		if err != nil {
			return nil, InternalError("sqlc.GetUserByName", err)
		}

		if sess := authStateFor(req.Context()).Session; sess != nil {
			if err := db.InvalidateSession(req.Context(), sess.ID); err != nil {
				slog.Warn("REST Auth: InvalidateSession",
					slog.String("op", "sqlc.InvalidateSession"),
					slog.String("err", err.Error()))
			}
		}

		state := AuthState{
			Actor: &Actor{
				ID:          user.ID,
				Name:        user.Name,
				Permissions: user.Permissions.Permissions,
			},
			Session: &AuthSession{
				ID:      sess.ID,
				Expires: sess.ExpiresAt.Time,
			},
		}

		state.CookieJar.SetSession(sess.ID, sess.Secret, expires)

		return &state, nil
	}))

	mux.Handle("POST /auth/signout", MethodHandleFunc(func(req *http.Request) (*AuthState, error) {

		auth := authStateFor(req.Context())
		if auth.Session == nil {
			return auth, nil
		}

		if err := db.InvalidateSession(req.Context(), auth.Session.ID); err != nil {
			return nil, InternalError("sqlc.InvalidateSession", err)
		}

		var state AuthState
		state.CookieJar.ClearSession()

		return &state, nil
	}))

	mux.Handle("PUT /manage/content/collection", MethodHandleFunc(func(req *http.Request) (*model.CollectionMetadata, error) {

		if perms, err := authStateFor(req.Context()).Permissions(); err != nil {
			return nil, err
		} else if err := perms.CanEditContent(); err != nil {
			return nil, err
		}

		params, err := ParseGeneric[model.CollectionPatch](req)
		if err != nil {
			return nil, err
		}

		return rslv.CreateContentCollection(req.Context(), params)
	}))

	mux.Handle("PATCH /manage/content/collection/{id}/metadata", MethodHandleFunc(func(req *http.Request) (*model.CollectionMetadata, error) {

		if perms, err := authStateFor(req.Context()).Permissions(); err != nil {
			return nil, err
		} else if err := perms.CanEditContent(); err != nil {
			return nil, err
		}

		collectionID, err := ParseUUID(req.PathValue("id"))
		if err != nil {
			return nil, err
		}

		params, err := ParseGeneric[model.CollectionPatch](req)
		if err != nil {
			return nil, err
		}

		return rslv.UpdateContentCollection(req.Context(), collectionID, params)
	}))

	mux.Handle("DELETE /manage/content/collection/{id}", MethodHandleFunc(func(req *http.Request) (*any, error) {

		if perms, err := authStateFor(req.Context()).Permissions(); err != nil {
			return nil, err
		} else if err := perms.CanEditContent(); err != nil {
			return nil, err
		}

		collectionID, err := ParseUUID(req.PathValue("id"))
		if err != nil {
			return nil, err
		}

		return nil, rslv.DeleteCollection(req.Context(), collectionID)
	}))

	mux.Handle("PUT /manage/content/deck", MethodHandleFunc(func(req *http.Request) (*model.CardDeckMetadata, error) {

		if perms, err := authStateFor(req.Context()).Permissions(); err != nil {
			return nil, err
		} else if err := perms.CanEditContent(); err != nil {
			return nil, err
		}

		params, err := ParseGeneric[model.CardDeckPatch](req)
		if err != nil {
			return nil, err
		}

		return rslv.CreateCardDeck(req.Context(), params)
	}))

	mux.Handle("PATCH /manage/content/deck/{id}", MethodHandleFunc(func(req *http.Request) (*model.CardDeckMetadata, error) {

		if perms, err := authStateFor(req.Context()).Permissions(); err != nil {
			return nil, err
		} else if err := perms.CanEditContent(); err != nil {
			return nil, err
		}

		deckID, err := ParseUUID(req.PathValue("id"))
		if err != nil {
			return nil, err
		}

		params, err := ParseGeneric[model.CardDeckPatch](req)
		if err != nil {
			return nil, err
		}

		return rslv.UpdateCardDeck(req.Context(), deckID, params)
	}))

	mux.Handle("DELETE /manage/content/deck/{id}", MethodHandleFunc(func(req *http.Request) (*any, error) {

		if perms, err := authStateFor(req.Context()).Permissions(); err != nil {
			return nil, err
		} else if err := perms.CanEditContent(); err != nil {
			return nil, err
		}

		deckID, err := ParseUUID(req.PathValue("id"))
		if err != nil {
			return nil, err
		}

		return nil, rslv.DeleteDeck(req.Context(), deckID)
	}))

	return authMiddleware(&NativeAuthProvider{DB: db}, mux)
}
