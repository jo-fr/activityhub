package httputil_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/jo-fr/activityhub/pkg/util/httputil"
)

func TestSanitizeURL(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		// Test cases with different URLs
		{"http://example.com/path?query=param", "http://example.com/path"},
		{"https://example.org/path/", "https://example.org/path"},
		{"ftp://example.net/", "ftp://example.net"},

		// Test cases with missing schemes
		{"www.example.com", "https://www.example.com"},
		{"example.org", "https://example.org"},
		{"ftpserver", "https://ftpserver"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result, err := httputil.SanitizeURL(tc.input)

			// Check for error if expected result is an empty string
			if err != nil && tc.expected != "" {
				t.Errorf("Unexpected error: %v", err)
			}

			// Check if the result matches the expected value
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}
func TestUnmarshalBody(t *testing.T) {
	type TestData struct {
		Foo string `json:"foo"`
		Bar int    `json:"bar"`
	}

	testCases := []struct {
		name        string
		body        string
		expected    TestData
		errExpected bool
	}{
		{
			name: "Valid JSON",
			body: `{"foo": "hello", "bar": 42}`,
			expected: TestData{
				Foo: "hello",
				Bar: 42,
			},
		},
		{
			name:        "Empty body",
			body:        "",
			expected:    TestData{},
			errExpected: true,
		},
		{
			name:        "Invalid JSON",
			body:        `{"foo": "hello", "bar": "invalid"}`,
			expected:    TestData{},
			errExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body := io.NopCloser(strings.NewReader(tc.body))
			result, err := httputil.UnmarshalBody[TestData](body)

			if tc.errExpected && err != nil {
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}
func TestStatusOK(t *testing.T) {
	testCases := []struct {
		statusCode int
		expected   bool
	}{
		{200, true},
		{201, true},
		{299, true},
		{300, false},
		{400, false},
		{500, false},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("StatusCode %d", tc.statusCode), func(t *testing.T) {
			result := httputil.StatusOK(tc.statusCode)
			if result != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}
func TestGetBody(t *testing.T) {

	testCases := []struct {
		name     string
		body     io.ReadCloser
		expected *bytes.Buffer
		err      error
	}{
		{
			name:     "Valid body",
			body:     io.NopCloser(strings.NewReader("test")),
			expected: bytes.NewBufferString("test"),
			err:      nil,
		},
		{
			name:     "Empty body",
			body:     io.NopCloser(strings.NewReader("")),
			expected: bytes.NewBufferString(""),
			err:      nil,
		},
		{
			name:     "Error reading body",
			body:     io.NopCloser(errorReader{}),
			expected: nil,
			err:      errors.New("error reading body"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := httputil.GetBody(tc.body)

			if err != nil && tc.err == nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if err == nil && tc.err != nil {
				t.Errorf("Expected error: %v, got nil", tc.err)
			}

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}

type errorReader struct{}

func (er errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("error reading body")
}
