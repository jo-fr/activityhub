package errutil

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jo-fr/activityhub/pkg/validate"
)

// predefined errors

var (
	InternalServerError = AnnotatedError{
		Message: "internal server error",
	}
)

// AnnotatedError is an error that contains additional information
type AnnotatedError struct {
	Type    ErrorType `json:"type,omitempty"`
	Message string    `json:"message"`
}

func (e AnnotatedError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// HTTPStatusCode returns the HTTP status code for error type
func (e AnnotatedError) HTTPStatusCode() int {
	switch e.Type {
	case TypeNotFound:
		return http.StatusNotFound
	case TypeInvalidRequestBody, TypeMissingHeader, TypeBadRequest, TypeAlreadyExists, TypeValidationError:
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

func ValidationErrorsToAnnotatedError(err validator.ValidationErrors) []AnnotatedError {
	var annotatedErrors []AnnotatedError
	for _, fieldError := range err {
		annotatedError := AnnotatedError{
			Type:    TypeValidationError,
			Message: validate.GetErrorMessage(fieldError),
		}
		annotatedErrors = append(annotatedErrors, annotatedError)
	}
	return annotatedErrors

}
