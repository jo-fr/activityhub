package errutil

// TEST HALLO
type ErrorType string

const (
	TypeNotFound           ErrorType = "not_found"
	TypeAlreadyExists      ErrorType = "already_exists"
	TypeInvalidRequestBody ErrorType = "invalid_request_body"
	TypeValidationError    ErrorType = "validation_error"
	TypeBadRequest         ErrorType = "bad_request"
	TypeMissingHeader      ErrorType = "missing_header"
)
