package render

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jo-fr/activityhub/pkg/log"
)

// Success processes successful requests and returns a JSON response.
func Success(ctx context.Context, data any, statusCode int, w http.ResponseWriter, log *log.Logger) {
	json, err := json.Marshal(data)
	if err != nil {
		Error(ctx, err, w, log)
		return
	}

	w.WriteHeader(statusCode)
	w.Write(json) // nolint: errcheck
}
