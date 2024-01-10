package httputil_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/jo-fr/activityhub/backend/pkg/util/httputil"
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

func TestUnmarshalRequestBody(t *testing.T) {
	type SampleStruct struct {
		Message string `json:"message"`
	}

	tests := []struct {
		name             string
		requestBody      string
		expectedResult   SampleStruct
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name:        "Valid JSON",
			requestBody: `{"message": "Hello, World!"}`,
			expectedResult: SampleStruct{
				Message: "Hello, World!",
			},
			expectError:      false,
			expectedErrorMsg: "",
		},
		{
			name:             "Invalid JSON",
			requestBody:      `invalid json`,
			expectedResult:   SampleStruct{},
			expectError:      true,
			expectedErrorMsg: "failed to unmarshal request body: invalid character 'i' looking for beginning of value",
		},
		{
			name:             "Empty Body",
			requestBody:      "",
			expectedResult:   SampleStruct{},
			expectError:      true,
			expectedErrorMsg: "failed to unmarshal request body: unexpected end of JSON input",
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a sample request with the specified body
			request, err := http.NewRequest("POST", "http://example.com", bytes.NewBufferString(tt.requestBody))
			if err != nil {
				t.Fatalf("Error creating request: %v", err)
			}

			// Call the function
			result, err := httputil.UnmarshalRequestBody[SampleStruct](request)

			if err != nil {
				if tt.expectError {
					if err.Error() != tt.expectedErrorMsg {
						t.Errorf("Unexpected error message. Got %s, want %s", err.Error(), tt.expectedErrorMsg)
					}
					return
				}

				t.Errorf("Unexpected error: %v", err)
			}

			// Check the unmarshaled result
			if result != tt.expectedResult {
				t.Errorf("Unexpected result. Got %+v, want %+v", result, tt.expectedResult)
			}
		})
	}
}

func TestUnmarshalResponseBody(t *testing.T) {
	type SampleStruct struct {
		Message string `json:"message"`
	}

	tests := []struct {
		name             string
		requestBody      string
		expectedResult   SampleStruct
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name:        "Valid JSON",
			requestBody: `{"message": "Hello, World!"}`,
			expectedResult: SampleStruct{
				Message: "Hello, World!",
			},
			expectError:      false,
			expectedErrorMsg: "",
		},
		{
			name:             "Invalid JSON",
			requestBody:      `invalid json`,
			expectedResult:   SampleStruct{},
			expectError:      true,
			expectedErrorMsg: "failed to unmarshal request body: invalid character 'i' looking for beginning of value",
		},
		{
			name:             "Empty Body",
			requestBody:      "",
			expectedResult:   SampleStruct{},
			expectError:      true,
			expectedErrorMsg: "failed to unmarshal request body: unexpected end of JSON input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			resp := &http.Response{
				Body: io.NopCloser(bytes.NewBufferString(tt.requestBody)),
			}

			// Call the function
			result, err := httputil.UnmarshalResponsetBody[SampleStruct](resp)

			if err != nil {
				if tt.expectError {
					if err.Error() != tt.expectedErrorMsg {
						t.Errorf("Unexpected error message. Got %s, want %s", err.Error(), tt.expectedErrorMsg)
					}
					return
				}

				t.Errorf("Unexpected error: %v", err)
			}

			// Check the unmarshaled result
			if result != tt.expectedResult {
				t.Errorf("Unexpected result. Got %+v, want %+v", result, tt.expectedResult)
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

func TestCopyBody(t *testing.T) {
	tests := []struct {
		name          string
		requestBody   string
		expectedError error
		expectedBody  string
	}{
		{
			name:          "Valid Case",
			requestBody:   "Hello, World!",
			expectedError: nil,
			expectedBody:  "Hello, World!",
		},
		{
			name:          "Empty Body",
			requestBody:   "",
			expectedError: nil,
			expectedBody:  "",
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a sample request with the specified body
			request, err := http.NewRequest("POST", "http://example.com", bytes.NewBufferString(tt.requestBody))
			if err != nil {
				t.Fatalf("Error creating request: %v", err)
			}

			// Call the function
			resultBuffer, err := httputil.CopyBody(request)

			// Check for errors
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) || (err != nil && err.Error() != tt.expectedError.Error()) {
				t.Errorf("Unexpected error. Got %v, want %v", err, tt.expectedError)
			}

			// Check the content of the body
			if resultBuffer != nil && resultBuffer.String() != tt.expectedBody {
				t.Errorf("Unexpected body content. Got %s, want %s", resultBuffer.String(), tt.expectedBody)
			}
		})
	}
}
