package httputil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// UnmarshalRequestBody unmarshals the request body into the given type. The request body is then reset to its original state.
func UnmarshalRequestBody[T any](r *http.Request) (T, error) {
	var v T
	body, err := GetBody(r)
	if err != nil {
		return v, errors.Wrap(err, "failed to read request body")
	}

	err = json.Unmarshal(body.Bytes(), &v)
	return v, errors.Wrap(err, "failed to unmarshal request body")
}

func UnmarshalResponseBody[T any](r *http.Response) (T, error) {
	var v T
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return v, errors.Wrap(err, "failed to read request body")
	}

	r.Body = io.NopCloser(bytes.NewReader(body))

	err = json.Unmarshal(body, &v)
	return v, errors.Wrap(err, "failed to unmarshal request body")
}

// GetBody reads the request body and returns it as a bytes.Buffer. The request body is then reset to its original state.
func GetBody(r *http.Request) (*bytes.Buffer, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(body)

	r.Body = io.NopCloser(buffer)

	return buffer, nil
}

// StatusOK returns true if the status code is between 200 and 299.
func StatusOK(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}
