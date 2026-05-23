package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/maddsua/flippercardapp/auth"
	db_model "github.com/maddsua/flippercardapp/db/model"
	"github.com/maddsua/flippercardapp/rest/model"
)

func MethodHandleFunc[T any](fn func(req *http.Request) (*T, error)) http.Handler {
	return http.HandlerFunc(func(wrt http.ResponseWriter, req *http.Request) {
		NewResponse(fn(req)).Write(wrt)
	})
}

func NewResponse[T any](data *T, err error) *model.Response[T] {
	if err != nil {
		return NewErrorResponseStatus[T](err, http.StatusBadRequest)
	}
	return &model.Response[T]{Data: data}
}

func NewErrorResponseStatus[T any](err error, code int) *model.Response[T] {

	if err, ok := err.(*model.Error); ok {
		return &model.Response[T]{Error: err}
	}

	if sc, ok := err.(interface{ StatusCode() int }); ok {
		return &model.Response[T]{Error: &model.Error{
			Message: err.Error(),
			Code:    sc.StatusCode(),
		}}
	}

	return &model.Response[T]{Error: &model.Error{
		Message: err.Error(),
		Code:    code,
	}}
}

func RespondWithError(wrt http.ResponseWriter, err error) {
	NewErrorResponseStatus[any](err, http.StatusBadRequest).Write(wrt)
}

func InternalError(op string, err error) error {

	slog.Error("Internal server error",
		slog.String("op", op),
		slog.String("err", err.Error()))

	return &model.Error{Message: fmt.Sprintf("%s: %v", op, err), Code: http.StatusInternalServerError}
}

func ParseUUIDSet(val string) (uuid.UUIDs, error) {

	if val = strings.TrimSpace(val); val == "" {
		return nil, nil
	}

	idMap := map[uuid.UUID]struct{}{}

	for item := range strings.SplitSeq(val, ",") {
		id, err := ParseUUID(strings.TrimSpace(item))
		if err != nil {
			return nil, err
		}
		idMap[id] = struct{}{}
	}

	var result []uuid.UUID
	for item := range idMap {
		result = append(result, item)
	}

	return result, nil
}

func ParseUUID(val string) (uuid.UUID, error) {
	id, err := uuid.Parse(val)
	if err != nil {
		return uuid.UUID{}, &model.Error{Message: "Invalid resource ID"}
	}
	return id, nil
}

func ParseGenericJSON[T any](req *http.Request) (T, error) {

	var result T

	if req.ContentLength == 0 {
		return result, &model.Error{
			Message: "empty json payload",
			Code:    http.StatusLengthRequired,
		}
	}

	if contentType := req.Header.Get("Content-Type"); !strings.EqualFold(contentType, "application/json") {
		return result, &model.Error{
			Message: fmt.Sprintf("content type '%v' not allowed", contentType),
			Code:    http.StatusUnprocessableEntity,
		}
	}

	if err := json.NewDecoder(req.Body).Decode(&result); err != nil {
		return result, &model.Error{
			Message: fmt.Sprintf("content decoding error: %v", err),
			Code:    http.StatusBadRequest,
		}
	}

	return result, nil
}

func ResourceVisibilityFilter(ctx context.Context, idFilter uuid.UUIDs) db_model.ResourceVisibilities {

	if perms, _ := auth.For(ctx).Permissions(); perms != nil && perms.AsTeamMember() == nil {
		return db_model.ResourceVisibilities{
			db_model.ResourceVisibilityPrivate,
			db_model.ResourceVisibilityHidden,
			db_model.ResourceVisibilityPublic,
		}
	}

	if len(idFilter) > 0 {
		return db_model.ResourceVisibilities{
			db_model.ResourceVisibilityHidden,
			db_model.ResourceVisibilityPublic,
		}
	}

	return db_model.ResourceVisibilities{
		db_model.ResourceVisibilityPublic,
	}
}

func EnforceResourceVisibility(ctx context.Context, resourceVisibility db_model.ResourceVisibility) error {

	switch resourceVisibility {

	case db_model.ResourceVisibilityPrivate:

		if perms, _ := auth.For(ctx).Permissions(); perms == nil || perms.AsTeamMember() != nil {
			return &model.Error{Message: "Resource access denied", Code: http.StatusForbidden}
		}

		return nil

	default:
		return nil
	}
}

func UnwrapSearchIndex(index map[uuid.UUID]int, page PagePointers) []uuid.UUID {

	type RankedSearchEntry struct {
		ID   uuid.UUID
		Rank int
	}

	var indexSlice []RankedSearchEntry
	for id, rank := range index {
		indexSlice = append(indexSlice, RankedSearchEntry{ID: id, Rank: rank})
	}

	sort.SliceStable(indexSlice, func(i, j int) bool {
		return indexSlice[i].Rank < indexSlice[j].Rank
	})

	indexSlice = SliceQueriedPage(indexSlice, page)

	result := make([]uuid.UUID, len(indexSlice))
	for idx, val := range indexSlice {
		result[idx] = val.ID
	}

	return result
}
