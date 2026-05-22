package rest

import (
	"math"
	"net/http"
	"strconv"
)

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

func SliceQueriedPage[T any](entres []T, page PagePointers) []T {
	offset := page.QueryOffset()
	if int(offset) >= len(entres) {
		return nil
	}
	return entres[offset:min(len(entres), int(page.QueryLimit()))]
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

type Page[T any] struct {
	PagePointers
	Entries []T  `json:"entries"`
	HasNext bool `json:"has_next"`
}

func TransformPage[E any, R any](page PagePointers, entries []E, transformer func(E) R) *Page[R] {

	resultSize := min(len(entries), page.Limit)

	result := make([]R, resultSize)
	for idx := range resultSize {
		result[idx] = transformer(entries[idx])
	}

	return WrapPage(page, result)
}

func WrapPage[T any](page PagePointers, entries []T) *Page[T] {

	if entries == nil {
		entries = make([]T, 0)
	}

	return &Page[T]{
		PagePointers: page,
		Entries:      entries[0:min(len(entries), page.Limit)],
		HasNext:      len(entries) > page.Limit,
	}
}
