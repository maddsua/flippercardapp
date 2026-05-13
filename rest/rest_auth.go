package rest

import (
	"context"
	"crypto/subtle"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

type authStateKey struct{}

type AuthState struct {
	Actor       *Actor           `json:"actor"`
	Permissions ActorPermissions `json:"permissions"`
}

type Actor struct {
	ID string `json:"id"`
}

type ActorPermissions struct {
	Administrative bool `json:"administrative"`
	ContentEdit    bool `json:"content_edit"`
}

type AuthProvider interface {
	AuthorizeRequest(req *http.Request) (*AuthState, error)
}

func authMiddleware(auth AuthProvider, next http.Handler) http.Handler {
	return http.HandlerFunc(func(wrt http.ResponseWriter, req *http.Request) {

		state, err := auth.AuthorizeRequest(req)
		if err != nil {

			slog.Warn("HTTP Auth",
				slog.String("client", req.RemoteAddr),
				slog.String("err", err.Error()))

			NewErrorResponseStatus[any](err, http.StatusUnauthorized).Write(wrt)

			return
		}

		next.ServeHTTP(wrt, req.WithContext(context.WithValue(req.Context(), authStateKey{}, state)))
	})
}

func authStateFor(ctx context.Context) *AuthState {
	val := ctx.Value(authStateKey{})
	return val.(*AuthState)
}

type StaticAuthProvider struct {
}

func (auth *StaticAuthProvider) AuthorizeRequest(req *http.Request) (*AuthState, error) {

	bearer := auth.requestBearerToken(req)
	if bearer == "" {
		return &AuthState{}, nil
	}

	if edtorToken := auth.envToken("EDITOR_TOKEN"); edtorToken != "" && subtle.ConstantTimeCompare([]byte(bearer), []byte(edtorToken)) == 1 {
		return &AuthState{
			Actor: &Actor{
				ID: "editor",
			},
			Permissions: ActorPermissions{
				Administrative: true,
				ContentEdit:    true,
			},
		}, nil
	}

	return &AuthState{}, &APIError{Message: "Invalid auth token", Code: http.StatusForbidden}
}

func (auth *StaticAuthProvider) requestBearerToken(req *http.Request) string {
	token, value, ok := strings.Cut(req.Header.Get("Authorization"), " ")
	if !ok || !strings.EqualFold(token, "Bearer") {
		return ""
	}
	return strings.TrimSpace(value)
}

func (auth *StaticAuthProvider) envToken(name string) string {
	return strings.TrimSpace(os.Getenv(name))
}
