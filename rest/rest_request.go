package rest

import (
	"math"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

func ReqPathUUID(val string) (uuid.UUID, error) {
	id, err := uuid.Parse(val)
	if err != nil {
		return uuid.UUID{}, &APIError{Message: "Invalid resource ID"}
	}
	return id, nil
}

func ReqParamUUID(val string) uuid.NullUUID {
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
		limit = math.MaxInt
	}

	offset = max(0, parsePagePointer(req.URL.Query().Get("offset")))

	return
}

func parsePagePointer(val string) int {
	intVal, _ := strconv.ParseInt(val, 10, 64)
	return min(max(math.MinInt, int(intVal)), math.MaxInt)
}
