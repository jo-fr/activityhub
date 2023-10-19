package errutil

import (
	"fmt"
	"net/http"
)

// predefined errors

var (
	InternalServerError = AnnotatedError{
		Reason:  ReasonAPIError,
		Message: "internal server error",
	}
)

// AnnotatedError is an error that contains additional information
type AnnotatedError struct {
	Reason    ErrorReason `json:"reason"`
	Type      ErrorType   `json:"type,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
	Message   string      `json:"message"`
}

func (e AnnotatedError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// HTTPStatusCode returns the HTTP status code for error type
func (e AnnotatedError) HTTPStatusCode() int {
	switch e.Type {
	case TypeNotFound:
		return http.StatusNotFound
	case TypeInvalidRequestBody, TypeMissingHeader:
		return http.StatusBadRequest

	default:
		return http.StatusInternalServerError
	}
}

// NewError returns a new AnnotatedError with the given message
func NewError(errType ErrorType, msg string) error {
	return AnnotatedError{
		Type:    errType,
		Message: msg,
	}
}

// ExtractAnnotedError returns the AnnotatedError from the given error, if it exists
func ExtractAnnotedError(err error) (AnnotatedError, bool) {
	annotatedErr, ok := err.(AnnotatedError)
	return annotatedErr, ok
}
