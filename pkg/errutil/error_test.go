package errutil_test

import (
	"net/http"
	"testing"

	"github.com/jo-fr/activityhub/pkg/errutil"
)

func TestHTTPStatusCode(t *testing.T) {
	testCases := []struct {
		err      errutil.AnnotatedError
		expected int
	}{
		{err: errutil.AnnotatedError{Type: errutil.TypeNotFound}, expected: http.StatusNotFound},
		{err: errutil.AnnotatedError{Type: errutil.TypeInvalidRequestBody}, expected: http.StatusBadRequest},
		{err: errutil.AnnotatedError{Type: errutil.TypeMissingHeader}, expected: http.StatusBadRequest},
		{err: errutil.AnnotatedError{Type: errutil.TypeBadRequest}, expected: http.StatusBadRequest},
		{err: errutil.AnnotatedError{Type: errutil.TypeAlreadyExists}, expected: http.StatusBadRequest},
		{err: errutil.AnnotatedError{Type: errutil.TypeValidationError}, expected: http.StatusBadRequest},
		{err: errutil.AnnotatedError{Type: "unknown"}, expected: http.StatusInternalServerError},
	}

	for _, tc := range testCases {
		t.Run(tc.err.Message, func(t *testing.T) {
			result := tc.err.HTTPStatusCode()

			if result != tc.expected {
				t.Errorf("Expected %d, got %d", tc.expected, result)
			}
		})
	}
}
