package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func MethodHandleFunc[T any](fn func(req *http.Request) (*T, error)) http.Handler {
	return http.HandlerFunc(func(wrt http.ResponseWriter, req *http.Request) {
		NewResponse(fn(req)).Write(wrt)
	})
}

func NewResponse[T any](data *T, err error) *Response[T] {
	if err != nil {
		return NewErrorResponseStatus[T](err, http.StatusBadRequest)
	}
	return &Response[T]{Data: data}
}

func NewErrorResponseStatus[T any](err error, code int) *Response[T] {

	if err, ok := err.(*APIError); ok {
		return &Response[T]{Error: err}
	}

	if sc, ok := err.(StatusCoder); ok {
		return &Response[T]{Error: &APIError{
			Message: err.Error(),
			Code:    sc.StatusCode(),
		}}
	}

	return &Response[T]{Error: &APIError{
		Message: err.Error(),
		Code:    code,
	}}
}

type StatusCoder interface {
	StatusCode() int
}

type ReadableCookieJar interface {
	Cookies() []http.Cookie
}

type Response[T any] struct {
	Data  *T        `json:"data"`
	Error *APIError `json:"error"`
}

func (resp *Response[T]) Write(wrt http.ResponseWriter) {

	if resp.Data == nil && resp.Error == nil {
		wrt.WriteHeader(http.StatusNoContent)
		return
	}

	wrt.Header().Set("Content-Type", "application/json")

	if jar, ok := any(resp.Data).(ReadableCookieJar); ok {
		for _, cookie := range jar.Cookies() {
			http.SetCookie(wrt, &cookie)
		}
	}

	if err := resp.Error; err != nil {
		for _, cookie := range err.Cookies {
			http.SetCookie(wrt, &cookie)
		}
		wrt.WriteHeader(err.StatusCode())
	}

	json.NewEncoder(wrt).Encode(resp)
}

type APIError struct {
	Message string        `json:"message"`
	Code    int           `json:"-"`
	Cookies []http.Cookie `json:"-"`
}

func (err *APIError) StatusCode() int {
	// min-max the error code to avoid whoopsie-daisies with invalid statuses
	return min(max(http.StatusBadRequest, err.Code), http.StatusNetworkAuthenticationRequired)
}

func (err *APIError) Error() string {
	return err.Message
}

func InternalError(op string, err error) error {

	slog.Error("Internal server error",
		slog.String("op", op),
		slog.String("err", err.Error()))

	return &APIError{Message: err.Error(), Code: http.StatusInternalServerError}
}
