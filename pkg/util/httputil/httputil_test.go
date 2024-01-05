package httputil_test

import (
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
		{"invalid-url", ""}, // Error case, expecting an empty string for invalid URL

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
