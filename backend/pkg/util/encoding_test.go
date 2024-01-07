package util_test

import (
	"testing"

	"github.com/jo-fr/activityhub/backend/pkg/util"
)

func TestDecodeBase64(t *testing.T) {
	testCases := []struct {
		input    string
		expected []byte
	}{
		{input: "SGVsbG8gd29ybGQ=", expected: []byte("Hello world")},
		{input: "dGVzdCBtZXNzYWdl", expected: []byte("test message")},
		{input: "MTIzNDU2Nzg5MA==", expected: []byte("1234567890")},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result, err := util.DecodeBase64(tc.input)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if string(result) != string(tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}
