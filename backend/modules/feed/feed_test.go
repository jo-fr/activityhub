package feed

import (
	"testing"
)

func TestCamelToSnake(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"camelCaseString", "camel_case_string"},
		{"anotherExample", "another_example"},
		{"mixed123Case", "mixed123_case"},
		{"This isA Mixed Example", "this_is_a_mixed_example"},
		{"single", "single"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := CamelToSnake(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, but got %s", tt.expected, result)
			}
		})
	}
}

func TestCamelToSnake_NoUppercase(t *testing.T) {
	input := "nouppercasestring"
	expected := "nouppercasestring"

	result := CamelToSnake(input)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestCamelToSnake_WithDigits(t *testing.T) {
	input := "mixed123Case"
	expected := "mixed123_case"

	result := CamelToSnake(input)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestBuiltPost(t *testing.T) {
	testCases := []struct {
		title, description, link, expectedOutput string
	}{
		{
			title:          "Test Title",
			description:    "This is a test description.\nIt has multiple lines.",
			link:           "https://example.com/test",
			expectedOutput: "<p><strong>Test Title</strong><br/>This is a test description. It has multiple lines.</br><a href=\"https://example.com/test\" target=\"_blank\" rel=\"nofollow noopener noreferrer\" translate=\"no\">https://example.com/test</a></p>",
		},
		// Add more test cases for different scenarios
		{
			title:          "Another Title",
			description:    "Another description.",
			link:           "https://example.com/another",
			expectedOutput: "<p><strong>Another Title</strong><br/>Another description.</br><a href=\"https://example.com/another\" target=\"_blank\" rel=\"nofollow noopener noreferrer\" translate=\"no\">https://example.com/another</a></p>",
		},
	}

	for _, testCase := range testCases {
		result := builtPost(testCase.title, testCase.description, testCase.link)
		if result != testCase.expectedOutput {
			t.Errorf("For input (%s, %s, %s), expected %s, but got %s", testCase.title, testCase.description, testCase.link, testCase.expectedOutput, result)
		}
	}
}
