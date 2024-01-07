package render

import (
	"context"
	"encoding/json"

	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"

	"github.com/jo-fr/activityhub/backend/pkg/errutil"
	"github.com/jo-fr/activityhub/backend/pkg/log"
)

// ErrorReason is the reason for the error
type ErrorReason string

const (
	ReasonAPIError     ErrorReason = "api_error"
	ReasonRequestError ErrorReason = "request_error"
)

type errorResponse struct {
	Reason    ErrorReason              `json:"reason"`
	RequestID string                   `json:"requestID,omitempty"`
	Errors    []errutil.AnnotatedError `json:"errors"`
}

// Error processes internal errors and returns a JSON response.
func Error(ctx context.Context, err error, w http.ResponseWriter, log *log.Logger) {
	var errorResponse errorResponse
	var statusCode int
	switch errType := err.(type) {
	case validator.ValidationErrors:
		annotatedErrors := errutil.ValidationErrorsToAnnotatedError(errType)
		errorResponse = toErrorResponse(ReasonRequestError, middleware.GetReqID(ctx), annotatedErrors...)
		statusCode = http.StatusBadRequest
	case errutil.AnnotatedError:
		errorResponse = toErrorResponse(ReasonRequestError, middleware.GetReqID(ctx), errType)
		statusCode = errType.HTTPStatusCode()
	default:
		annerr := errutil.InternalServerError
		errorResponse = toErrorResponse(ReasonAPIError, middleware.GetReqID(ctx), annerr)
		statusCode = annerr.HTTPStatusCode()
		log.Error(err)
	}

	json, _ := json.Marshal(errorResponse)
	w.WriteHeader(statusCode)
	w.Write(json) // nolint: errcheck
}

func toErrorResponse(reason ErrorReason, reqID string, err ...errutil.AnnotatedError) errorResponse {
	return errorResponse{
		Reason:    reason,
		RequestID: reqID,
		Errors:    err,
	}
}
