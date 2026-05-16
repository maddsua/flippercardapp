package rest

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	db_pkg "github.com/maddsua/flippercardapp/db"
	db_gen "github.com/maddsua/flippercardapp/db/generated"
	db_model "github.com/maddsua/flippercardapp/db/model"
	"github.com/maddsua/flippercardapp/db/types"
)

const SessionTTL = 6 * time.Hour
const SessionCookieKey = "st"

type AuthCookieJar struct {
	Values []http.Cookie
}

func (jar *AuthCookieJar) SetSession(id uuid.UUID, secret []byte, expires time.Time) {
	token := sessionToken{id: id, secret: secret}
	jar.Values = append(jar.Values, http.Cookie{
		Name:     SessionCookieKey,
		Value:    token.String(),
		Path:     "/",
		Expires:  expires,
		Secure:   true,
		HttpOnly: true,
	})
}

func (jar *AuthCookieJar) ClearSession() {
	jar.Values = append(jar.Values, http.Cookie{
		Name:     SessionCookieKey,
		Path:     "/",
		Expires:  time.Unix(0, 0),
		Secure:   true,
		HttpOnly: true,
	})
}

func parseSessionToken(val string) (*sessionToken, error) {

	idToken, secretToken, ok := strings.Cut(val, ".")
	if !ok {
		return nil, fmt.Errorf("invalid session token format")
	}

	idBytes, err := base64.RawStdEncoding.DecodeString(idToken)
	if err != nil {
		return nil, fmt.Errorf("invalid session token id: %v", err)
	}

	id, err := uuid.FromBytes(idBytes)
	if err != nil {
		return nil, fmt.Errorf("invalid session token id: %v", err)
	}

	secret, err := base64.RawStdEncoding.DecodeString(secretToken)
	if err != nil {
		return nil, fmt.Errorf("invalid session token secret: %v", err)
	}

	return &sessionToken{id: id, secret: secret}, nil
}

type sessionToken struct {
	id     uuid.UUID
	secret []byte
}

func (token *sessionToken) String() string {
	return base64.RawStdEncoding.EncodeToString(token.id[:]) + "." + base64.RawStdEncoding.EncodeToString(token.secret)
}

type authStateKey struct{}

type AuthState struct {
	Actor     *Actor        `json:"actor"`
	Session   *AuthSession  `json:"session"`
	CookieJar AuthCookieJar `json:"-"`
}

func (state *AuthState) Cookies() []http.Cookie {
	if state == nil {
		return nil
	}
	return state.CookieJar.Values
}

func (state *AuthState) Permissions() (*ActorPermissions, error) {

	if state.Actor == nil {
		return nil, &APIError{Message: "Unauthorized", Code: http.StatusUnauthorized}
	}

	return &ActorPermissions{UserPermissions: state.Actor.Permissions}, nil
}

type ActorPermissions struct {
	db_model.UserPermissions
}

func (perms *ActorPermissions) IsAdmin() error {
	if !perms.UserPermissions.Administrative {
		return &APIError{Message: "'Administrative' permission required", Code: http.StatusForbidden}
	}
	return nil
}

func (perms *ActorPermissions) CanEditContent() error {

	if err := perms.IsAdmin(); err != nil {
		return err
	}

	if !perms.ContentEdit {
		return &APIError{Message: "'ContentEdit' permission required", Code: http.StatusForbidden}
	}

	return nil
}

type AuthSession struct {
	ID      uuid.UUID `json:"id"`
	Expires time.Time `json:"expires"`
}

type Actor struct {
	ID          uuid.UUID                `json:"id"`
	Name        string                   `json:"name"`
	Permissions db_model.UserPermissions `json:"permissions"`
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

		// this path is only used for the session lifetime extension
		for _, cookie := range state.Cookies() {
			http.SetCookie(wrt, &cookie)
		}

		next.ServeHTTP(wrt, req.WithContext(context.WithValue(req.Context(), authStateKey{}, state)))
	})
}

func authStateFor(ctx context.Context) *AuthState {
	val := ctx.Value(authStateKey{})
	return val.(*AuthState)
}

type NativeAuthProvider struct {
	DB *db_pkg.Wrapper
}

func (auth *NativeAuthProvider) AuthorizeRequest(req *http.Request) (*AuthState, error) {

	cookie, _ := req.Cookie(SessionCookieKey)
	if cookie == nil || cookie.Value == "" {
		return &AuthState{}, nil
	}

	var cookies AuthCookieJar

	var invalidateRequest = func(message string) (*AuthState, error) {
		cookies.ClearSession()
		return nil, &APIError{Message: message, Code: http.StatusUnauthorized, Cookies: cookies.Values}
	}

	token, err := parseSessionToken(cookie.Value)
	if err != nil {
		return invalidateRequest("Invalid session cookie")
	}

	session, err := auth.DB.GetSession(req.Context(), token.id)
	if db_pkg.IsNull(err) || (err == nil && !sessionValid(session)) {
		return invalidateRequest("Session expired")
	} else if err != nil {
		return nil, InternalError("sqlc.GetSession", err)
	} else if subtle.ConstantTimeCompare(token.secret, session.Secret) != 1 {
		return invalidateRequest("Invalid session secret")
	}

	if thresold := SessionTTL / 10; time.Until(session.ExpiresAt.Time) < thresold {

		newExpiration := session.ExpiresAt.Add(SessionTTL)

		slog.Debug("Extending session lifetime",
			slog.String("id", session.ID.String()),
			slog.String("client", req.RemoteAddr),
			slog.String("from", session.ExpiresAt.String()),
			slog.String("to", newExpiration.String()))

		if session, err = auth.DB.SetSessionExpirationTime(req.Context(), db_gen.SetSessionExpirationTimeParams{
			ID:        session.ID,
			ExpiresAt: types.NewTime(newExpiration),
		}); err != nil {
			return nil, InternalError("sqlc.SetSessionExpirationTime", err)
		}

		cookies.SetSession(session.ID, session.Secret, session.ExpiresAt.Time)
	}

	user, err := auth.DB.GetUserByID(req.Context(), session.UserID)
	if err != nil {
		return nil, InternalError("sqlc.GetUserByID", err)
	}

	return &AuthState{
		Actor: &Actor{
			ID:          user.ID,
			Name:        user.Name,
			Permissions: user.Permissions.Permissions,
		},
		Session: &AuthSession{
			ID:      session.ID,
			Expires: session.ExpiresAt.Time,
		},
		CookieJar: cookies,
	}, nil
}

func sessionValid(session db_gen.UserSession) bool {
	return session.ExpiresAt.After(time.Now()) && len(session.Secret) > 0
}
