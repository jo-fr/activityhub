package render

import (
	"context"
	"encoding/json"

	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/jo-fr/activityhub/pkg/errutil"
	"github.com/jo-fr/activityhub/pkg/log"
)

type errorResponse struct {
	Errors []errutil.AnnotatedError `json:"errors"`
}

// Error processes internal errors and returns a JSON response.
func Error(ctx context.Context, err error, w http.ResponseWriter, log *log.Logger) {
	annerr, ok := errutil.ExtractAnnotedError(err)
	if !ok {
		annerr = errutil.InternalServerError
		log.Error(err)
	} else {
		annerr.Reason = errutil.ReasonRequestError
	}

	annerr.RequestID = middleware.GetReqID(ctx)

	json, _ := json.Marshal(toErrorResponse(annerr))
	w.WriteHeader(annerr.HTTPStatusCode())
	w.Write(json)
}

func toErrorResponse(err ...errutil.AnnotatedError) errorResponse {
	return errorResponse{
		Errors: err,
	}
}
