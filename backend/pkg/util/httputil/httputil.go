package httputil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// // UnmarshalRequestBody unmarshals the response body into the given typ without changing the request body.
func UnmarshalRequestBody[T any](r *http.Request) (T, error) {
	var v T
	buffer, err := CopyBody(r)
	if err != nil {
		return v, errors.Wrap(err, "failed to read request body")
	}

	err = json.Unmarshal(buffer.Bytes(), &v)
	return v, errors.Wrap(err, "failed to unmarshal request body")
}

// // UnmarshalRequestBody unmarshals the response body into the given typ without changing the request body.
func UnmarshalResponsetBody[T any](r *http.Response) (T, error) {
	var v T
	buf, err := readBody(r.Body)
	if err != nil {
		return v, errors.Wrap(err, "failed to read response body")
	}

	err = json.Unmarshal(buf.Bytes(), &v)
	return v, errors.Wrap(err, "failed to unmarshal request body")
}

// CopyBody reads the request body and returns it as a bytes.Buffer. The request body is then reset to its original state.
func CopyBody(r *http.Request) (*bytes.Buffer, error) {
	buf, err := readBody(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read request body")
	}

	// reset the request body
	r.Body = io.NopCloser(buf)

	return buf, nil
}

func readBody(body io.ReadCloser) (*bytes.Buffer, error) {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(bodyBytes)

	return buf, nil
}

// StatusOK returns true if the status code is between 200 and 299.
func StatusOK(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

// SanitizeURL removes query parameters and trailing slashes from the given url. If the url has no scheme, https is used.
func SanitizeURL(urlStr string) (string, error) {

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	// if parsedURL.Hostname() == "" {
	// 	return "", errors.New("missing host")
	// }

	parsedURL.RawQuery = ""
	parsedURL.Path = strings.TrimSuffix(parsedURL.Path, "/")

	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "https"
	}

	return parsedURL.String(), nil
}
