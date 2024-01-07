package util_test

import (
	"testing"

	"github.com/jo-fr/activityhub/backend/pkg/util"
)

func TestTrimStringLength(t *testing.T) {
	testCases := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{input: "This is a sentence.", maxLen: 10, expected: "This is a  [...]"},
		{input: "This is a sentence.", maxLen: 20, expected: "This is a sentence."},
		{input: "This is a sentence.", maxLen: 5, expected: "This  [...]"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := util.TrimStringLength(tc.input, tc.maxLen)

			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestRemoveHTMLTags(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{input: "<p>This is a paragraph.</p>", expected: "This is a paragraph."},
		{input: "<h1>Title</h1>", expected: "Title"},
		{input: "<div><span>Content</span></div>", expected: "Content"},
		{input: "No HTML tags", expected: "No HTML tags"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := util.RemoveHTMLTags(tc.input)

			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}
