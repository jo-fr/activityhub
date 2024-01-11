package errutil

type ErrorType string

const (
	TypeNotFound        ErrorType = "not_found"
	TypeAlreadyExists   ErrorType = "already_exists"
	TypeValidationError ErrorType = "validation_error"
	TypeBadRequest      ErrorType = "bad_request"
	TypeMissingHeader   ErrorType = "missing_header"
)
