package auth

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
	rest_model "github.com/maddsua/flippercardapp/rest/model"
	"github.com/maddsua/flippercardapp/utils"
	"golang.org/x/crypto/bcrypt"
)

const SessionBaseTTL = 7 * 24 * time.Hour
const SessionMaxTTL = 28 * 24 * time.Hour
const SessionTTLExtenderThreshold = min(SessionBaseTTL, 2*time.Hour)
const SessionCookieKey = "st"
const SessionCookieExpOverlap = time.Minute

func ParseSessionToken(val string) (*SessionToken, error) {

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

	return &SessionToken{ID: id, Secret: secret}, nil
}

type SessionToken struct {
	ID     uuid.UUID
	Secret []byte
}

func (token *SessionToken) String() string {
	return base64.RawStdEncoding.EncodeToString(token.ID[:]) + "." + base64.RawStdEncoding.EncodeToString(token.Secret)
}

type RequestAuth struct {
	Actor   *Actor       `json:"actor"`
	Session *AuthSession `json:"session"`
}

func (state *RequestAuth) Permissions() (*ActorPermissions, error) {

	if state.Actor == nil {
		return nil, &AuthError{}
	}

	return &ActorPermissions{UserPermissions: state.Actor.Permissions}, nil
}

type ActorPermissions struct {
	db_model.UserPermissions
}

func (perms *ActorPermissions) AsTeamMember() error {
	if !perms.UserPermissions.TeamMember {
		return &PermissionError{Permission: "TeamMember"}
	}
	return nil
}

func (perms *ActorPermissions) AsContentEditor() error {

	if err := perms.AsTeamMember(); err != nil {
		return err
	}

	if !perms.ContentEdit {
		return &PermissionError{Permission: "ContentEdit"}
	}

	return nil
}

type AuthError struct {
	Message  string
	Internal bool
}

func (err *AuthError) Error() string {

	if err.Message == "" {
		if err.Internal {
			return "Auth module error"
		}
		return "Unauthorized"
	}

	if err.Internal {
		return "Auth module: " + err.Message
	}

	return err.Message
}

func (err *AuthError) StatusCode() int {

	if err.Internal {
		return http.StatusInternalServerError
	}

	return http.StatusUnauthorized
}

type PermissionError struct {
	Permission string
}

func (err *PermissionError) Error() string {
	return fmt.Sprintf("'%s' permission required", err.Permission)
}

func (err *PermissionError) StatusCode() int {
	return http.StatusForbidden
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

type authContextKey struct{}

type authContext struct {
	req *RequestAuth
	wrt http.ResponseWriter
	db  *db_pkg.Wrapper
}

func Middleware(db *db_pkg.Wrapper, next http.Handler) http.Handler {
	return http.HandlerFunc(func(wrt http.ResponseWriter, req *http.Request) {

		reqAuth, err := authorizeRequest(db, wrt, req)
		if err != nil {
			slog.Warn("WEB AUTH",
				slog.String("client", req.RemoteAddr),
				slog.String("err", err.Error()))
			writeErrorResponse(wrt, err)
			return
		}

		state := &authContext{
			req: reqAuth,
			wrt: wrt,
			db:  db,
		}

		next.ServeHTTP(wrt, req.WithContext(context.WithValue(req.Context(), authContextKey{}, state)))
	})
}

func contextValue(ctx context.Context) *authContext {
	return ctx.Value(authContextKey{}).(*authContext)
}

func For(ctx context.Context) *RequestAuth {
	return contextValue(ctx).req
}

func authorizeRequest(db *db_pkg.Wrapper, wrt http.ResponseWriter, req *http.Request) (*RequestAuth, error) {

	cookie, _ := req.Cookie(SessionCookieKey)
	if cookie == nil || cookie.Value == "" {
		return &RequestAuth{}, nil
	}

	token, err := ParseSessionToken(cookie.Value)
	if err != nil {
		return invalidateRequest(wrt, "Invalid session cookie")
	}

	session, err := db.GetSession(req.Context(), token.ID)
	if db_pkg.IsNull(err) || (err == nil && !sessionValid(session)) {
		return invalidateRequest(wrt, "Session expired")
	} else if err != nil {
		return nil, internalError("sqlc.GetSession", err)
	} else if subtle.ConstantTimeCompare(token.Secret, session.Secret) != 1 {
		return invalidateRequest(wrt, "Invalid session secret")
	}

	if totalDuration := time.Since(session.CreatedAt.Time) + SessionBaseTTL; totalDuration < SessionMaxTTL {

		if expiresIn := time.Until(session.ExpiresAt.Time) - SessionCookieExpOverlap; expiresIn < SessionTTLExtenderThreshold {

			if session, err = db.SetSessionExpirationTime(req.Context(), db_gen.SetSessionExpirationTimeParams{
				ID:        session.ID,
				ExpiresAt: types.NewTime(session.ExpiresAt.Add(SessionBaseTTL)),
			}); err != nil {
				return nil, internalError("sqlc.SetSessionExpirationTime", err)
			}

			slog.Debug("WEB AUTH: User still active; Extending their session",
				slog.String("session_id", session.ID.String()),
				slog.String("client", req.RemoteAddr),
				slog.String("expires", session.ExpiresAt.String()))

			setSessionCookie(wrt, session)
		}
	}

	user, err := db.GetUserByID(req.Context(), session.UserID)
	if err != nil {
		return nil, internalError("sqlc.GetUserByID", err)
	}

	return &RequestAuth{
		Actor: &Actor{
			ID:          user.ID,
			Name:        user.Name,
			Permissions: user.Permissions.Permissions,
		},
		Session: &AuthSession{
			ID:      session.ID,
			Expires: session.ExpiresAt.Time,
		},
	}, nil
}

func invalidateRequest(wrt http.ResponseWriter, message string) (*RequestAuth, error) {
	clearSessionCookie(wrt)
	return nil, &AuthError{Message: message}
}

func internalError(op string, err error) error {

	slog.Error("WEB AUTH INTERNAL",
		slog.String("op", op),
		slog.String("err", err.Error()))

	return &AuthError{Message: "Auth module error", Internal: true}
}

func sessionValid(session db_gen.UserSession) bool {
	return session.ExpiresAt.After(time.Now()) && len(session.Secret) > 0
}

func setSessionCookie(wrt http.ResponseWriter, session db_gen.UserSession) {
	token := SessionToken{ID: session.ID, Secret: session.Secret}
	http.SetCookie(wrt, &http.Cookie{
		Name:     SessionCookieKey,
		Value:    token.String(),
		Path:     "/",
		Expires:  session.ExpiresAt.Time.Add(-SessionCookieExpOverlap),
		Secure:   true,
		HttpOnly: true,
	})
}

func clearSessionCookie(wrt http.ResponseWriter) {
	http.SetCookie(wrt, &http.Cookie{
		Name:     SessionCookieKey,
		Path:     "/",
		Expires:  time.Unix(0, 0),
		Secure:   true,
		HttpOnly: true,
	})
}

func writeErrorResponse(wrt http.ResponseWriter, err error) {

	restError := rest_model.Error{
		Message: err.Error(),
		Code:    http.StatusUnauthorized,
	}

	if sc, ok := err.(interface{ StatusCode() int }); ok {
		restError.Code = sc.StatusCode()
	}

	resp := rest_model.Response[any]{Error: &restError}

	resp.Write(wrt)
}

func NewWebSessionWithPassword(ctx context.Context, username, password string) (*RequestAuth, error) {

	state := contextValue(ctx)

	tx, err := state.db.BeginTx(ctx)
	if err != nil {
		return nil, internalError("sqlc.BeginTx", err)
	}

	defer tx.Rollback()

	user, err := tx.GetUserByName(ctx, username)
	if db_pkg.IsNull(err) {
		return nil, &AuthError{Message: "User not found"}
	} else if err != nil {
		return nil, internalError("sqlc.GetUserByName", err)
	} else if strings.TrimSpace(password) == "" {
		return nil, &AuthError{Message: "Invalid password"}
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {

		if err != bcrypt.ErrMismatchedHashAndPassword {
			slog.Error("REST: bcrypt.CompareHashAndPassword",
				slog.String("err", err.Error()))
		}

		return nil, &AuthError{Message: "Invalid password"}
	}

	expires := time.Now().Add(SessionBaseTTL)

	sess, err := tx.InsertSession(ctx, db_gen.InsertSessionParams{
		ID:        uuid.New(),
		CreatedAt: types.NewTime(time.Now()),
		ExpiresAt: types.NewTime(expires.Add(time.Minute)),
		UserID:    user.ID,
		Secret:    utils.NewRandomToken(128),
	})

	if err != nil {
		return nil, internalError("sqlc.GetUserByName", err)
	}

	if oldSess := state.req.Session; oldSess != nil {
		if err := tx.InvalidateSession(ctx, oldSess.ID); err != nil {
			slog.Warn("WEB AUTH: InvalidateSession",
				slog.String("op", "sqlc.InvalidateSession"),
				slog.String("err", err.Error()))
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, internalError("sqlc.Commit", err)
	}

	slog.Debug("WEB AUTH: Sign-in: New password session",
		slog.String("username", username),
		slog.String("session_id", sess.ID.String()),
		slog.String("expires", sess.ExpiresAt.String()))

	setSessionCookie(state.wrt, sess)

	return &RequestAuth{
		Actor: &Actor{
			ID:          user.ID,
			Name:        user.Name,
			Permissions: user.Permissions.Permissions,
		},
		Session: &AuthSession{
			ID:      sess.ID,
			Expires: sess.ExpiresAt.Time,
		},
	}, nil
}

func TerminateWebSession(ctx context.Context) (*RequestAuth, error) {

	state := contextValue(ctx)

	if state.req.Session == nil {
		return state.req, nil
	}

	if err := state.db.InvalidateSession(ctx, state.req.Session.ID); err != nil {
		return nil, internalError("sqlc.InvalidateSession", err)
	}

	slog.Debug("WEB AUTH: Sign-out",
		slog.String("session_id", state.req.Session.ID.String()))

	clearSessionCookie(state.wrt)

	return &RequestAuth{}, nil
}
