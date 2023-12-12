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
