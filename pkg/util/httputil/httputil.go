package httputil

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/pkg/errors"
)

// // UnmarshaBody unmarshals the response body into the given typ without changing the request body.
func UnmarshaBody[T any](body io.ReadCloser) (T, error) {
	var v T
	buffer, err := GetBody(body)
	if err != nil {
		return v, errors.Wrap(err, "failed to read request body")
	}

	err = json.Unmarshal(buffer.Bytes(), &v)
	return v, errors.Wrap(err, "failed to unmarshal request body")
}

// GetBody reads the request body and returns it as a bytes.Buffer. The request body is then reset to its original state.
func GetBody(body io.ReadCloser) (*bytes.Buffer, error) {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(bodyBytes)
	body = io.NopCloser(buffer)

	return buffer, nil
}

// StatusOK returns true if the status code is between 200 and 299.
func StatusOK(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}
