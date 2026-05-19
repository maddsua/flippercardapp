package rest

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
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

func InternalError(op string, err error) error {

	slog.Error("Internal server error",
		slog.String("op", op),
		slog.String("err", err.Error()))

	return &model.Error{Message: err.Error(), Code: http.StatusInternalServerError}
}

type UUIDSet map[uuid.UUID]struct{}

func (set UUIDSet) List() UUIDList {

	if set == nil {
		return nil
	}

	var entries []uuid.UUID
	for item := range set {
		entries = append(entries, item)
	}
	return entries
}

type UUIDList uuid.UUIDs

func (list UUIDList) WithPage(page PagePointers) []uuid.UUID {
	return SlicePage([]uuid.UUID(list), page)
}

func ParseUUIDSet(val string) (UUIDSet, error) {

	if val = strings.TrimSpace(val); val == "" {
		return nil, nil
	}

	set := UUIDSet{}

	for item := range strings.SplitSeq(val, ",") {
		id, err := ParseUUID(strings.TrimSpace(item))
		if err != nil {
			return nil, err
		}
		set[id] = struct{}{}
	}
	return set, nil
}

func ParseUUID(val string) (uuid.UUID, error) {
	id, err := uuid.Parse(val)
	if err != nil {
		return uuid.UUID{}, &model.Error{Message: "Invalid resource ID"}
	}
	return id, nil
}

func ParseGeneric[T any](req *http.Request) (T, error) {

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
