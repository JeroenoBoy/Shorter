package api

import (
	"fmt"
	"net/http"
)

type apiError struct {
	StatusCode int    `json:"code"`
	Message    string `json:"message"`
}

var (
	ErrorNotAuthenticated    = NewApiError(http.StatusUnauthorized, "you are not logged in")
	ErrorBadRequest          = NewApiError(http.StatusBadRequest, "malformed request")
	ErrorNoPermissions       = NewApiError(http.StatusForbidden, "you are not permitted to execute this action")
	ErrorInternalServerError = NewApiError(http.StatusInternalServerError, "internal server error")
)

func NewApiError(code int, message string) apiError {
	return apiError{
		StatusCode: code,
		Message:    message,
	}
}

func (err apiError) Error() string {
	return fmt.Sprintf("Code %v: %v", err.StatusCode, err.Message)
}

func WriteError(w http.ResponseWriter, err error) error {
	if apiErr, ok := err.(apiError); ok {
		return WriteResponse(w, apiErr.StatusCode, apiErr)
	} else {
		WriteError(w, ErrorInternalServerError)
		panic(err)
	}
}
