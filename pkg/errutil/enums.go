package errutil

// ErrorReason is the reason for the error
type ErrorReason string

const (
	ReasonAPIError     ErrorReason = "api_error"
	ReasonRequestError ErrorReason = "request_error"
)

type ErrorType string

const (
	TypeNotFound           ErrorType = "not_found"
	TypeAlreadyExists      ErrorType = "already_exists"
	TypeInvalidRequestBody ErrorType = "invalid_request_body"
	TypeBadRequest         ErrorType = "bad_request"
	TypeMissingHeader      ErrorType = "missing_header"
)
