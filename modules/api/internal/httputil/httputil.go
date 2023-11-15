package httputil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// UnmarshalBody unmarshals the request body into the given type. The request body is then reset to its original state.
func UnmarshalBody[T any](r *http.Request) (T, error) {
	var v T
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return v, errors.Wrap(err, "failed to read request body")
	}

	r.Body = io.NopCloser(bytes.NewReader(body))

	err = json.Unmarshal(body, &v)
	return v, errors.Wrap(err, "failed to unmarshal request body")
}
