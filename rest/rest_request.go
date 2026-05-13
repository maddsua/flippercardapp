package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

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
		return uuid.UUID{}, &APIError{Message: "Invalid resource ID"}
	}
	return id, nil
}

func ParseGeneric[T any](req *http.Request) (T, error) {

	var result T

	if req.ContentLength == 0 {
		return result, &APIError{
			Message: "empty patch",
			Code:    http.StatusLengthRequired,
		}
	}

	if contentType := req.Header.Get("Content-Type"); !strings.EqualFold(contentType, "application/json") {
		return result, &APIError{
			Message: fmt.Sprintf("content type '%v' not allowed", contentType),
			Code:    http.StatusUnprocessableEntity,
		}
	}

	if err := json.NewDecoder(req.Body).Decode(&result); err != nil {
		return result, &APIError{
			Message: fmt.Sprintf("content decoding error: %v", err),
			Code:    http.StatusBadRequest,
		}
	}

	return result, nil
}
