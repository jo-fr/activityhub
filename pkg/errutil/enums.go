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
	TypeInvalidRequestBody ErrorType = "invalid_request_body"
	TypeMissingHeader      ErrorType = "missing_header"
)
