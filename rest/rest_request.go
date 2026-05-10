package rest

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type UUIDSet map[uuid.UUID]struct{}

func (set UUIDSet) List() []uuid.UUID {
	var entries []uuid.UUID
	for item := range set {
		entries = append(entries, item)
	}
	return entries
}

func (set UUIDSet) WithPage(page PagePointers) []uuid.UUID {

	list := set.List()

	offset := page.QueryOffset()
	if int(offset) >= len(list) {
		return nil
	}

	return list[offset:min(len(list), int(page.QueryLimit()+1))]
}

func (set UUIDSet) IsEmpty() bool {
	return len(set) == 0
}

func ParseUUIDSet(val string) UUIDSet {
	set := UUIDSet{}
	for item := range strings.SplitSeq(val, ",") {
		if id, err := uuid.Parse(strings.TrimSpace(item)); err == nil {
			set[id] = struct{}{}
		}
	}
	return set
}

func ParseUUID(val string) (uuid.UUID, error) {
	id, err := uuid.Parse(val)
	if err != nil {
		return uuid.UUID{}, &APIError{Message: "Invalid resource ID"}
	}
	return id, nil
}

func ParseNullUUID(val string) uuid.NullUUID {
	id, err := uuid.Parse(val)
	if err != nil {
		return uuid.NullUUID{}
	}
	return uuid.NullUUID{UUID: id, Valid: true}
}

type PagePointers struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

func (pptr *PagePointers) QueryLimit() int64 {
	return int64(pptr.Limit + 1)
}

func (pptr *PagePointers) QueryOffset() int64 {
	return int64(pptr.Offset)
}

func Pagination(req *http.Request) PagePointers {
	limit, offset := requestQueryPointers(req)
	return PagePointers{Limit: int(limit), Offset: int(offset)}
}

func requestQueryPointers(req *http.Request) (limit int, offset int) {

	limit = parsePagePointer(req.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 25
	} else if limit < 0 {
		limit = math.MaxInt32
	}

	offset = max(0, parsePagePointer(req.URL.Query().Get("offset")))

	return
}

func parsePagePointer(val string) int {
	intVal, _ := strconv.ParseInt(val, 10, 64)
	return min(max(math.MinInt, int(intVal)), math.MaxInt32)
}
