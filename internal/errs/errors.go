package errs

import (
	"fmt"
	"net/http"
)

type Error struct {
	Kind   string `json:"kind"`
	Status int    `json:"status"`
	Reason string `json:"reason"`
}

func (e *Error) Error() string {
	return e.Kind
}
func (e *Error) String() string {
	return fmt.Sprintf("%s : %s : %s", e.Kind, http.StatusText(e.Status), e.Reason)
}

const (
	ErrNotFound         = "not_found"
	ErrNotAuthenticated = "not_authenticated"
	ErrInvalidID        = "invalid_id"
	ErrBadRequest       = "bad_request"
	ErrValidationFailed = "validation_failed"
	ErrInternal         = "internal_error"
)

func NotFound(reason string) *Error {
	return &Error{Kind: ErrNotFound, Status: http.StatusNotFound, Reason: reason}
}

func NotAuthenticated(reason string) *Error {
	return &Error{Kind: ErrNotAuthenticated, Status: http.StatusUnauthorized, Reason: reason}
}

func InvalidID(reason string) *Error {
	return &Error{Kind: ErrInvalidID, Status: http.StatusBadRequest, Reason: reason}
}

func BadRequest(reason string) *Error {
	return &Error{Kind: ErrBadRequest, Status: http.StatusBadRequest, Reason: reason}
}

func ValidationFailed(reason string) *Error {
	return &Error{Kind: ErrValidationFailed, Status: http.StatusUnprocessableEntity, Reason: reason}
}

func Internal(reason string) *Error {
	return &Error{Kind: ErrInternal, Status: http.StatusInternalServerError, Reason: reason}
}
