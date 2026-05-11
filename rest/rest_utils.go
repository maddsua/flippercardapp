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

	if err == nil {
		return &Response[T]{Data: data}
	}

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
		Code:    http.StatusBadRequest,
	}}
}

type StatusCoder interface {
	StatusCode() int
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

	if resp.Error != nil {
		wrt.WriteHeader(resp.Error.StatusCode())
	}

	json.NewEncoder(wrt).Encode(resp)
}

type APIError struct {
	Message string `json:"message"`
	Code    int    `json:"-"`
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
