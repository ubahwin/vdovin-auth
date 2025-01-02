package api

import "net/http"

type ErrorResponse struct {
	Error string `json:"error"`
}

func ErrHandler(_ *Context, err error) (interface{}, int) {
	return &ErrorResponse{
		Error: err.Error(),
	}, http.StatusBadRequest // TODO: wrap to errors types mapper with http status codes
}
