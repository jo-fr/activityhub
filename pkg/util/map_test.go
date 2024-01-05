package util_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/jo-fr/activityhub/pkg/util"
)

func TestMap(t *testing.T) {
	testCases := []struct {
		input    []int
		mapFunc  func(item int, index int) string
		expected []string
	}{
		{
			input:    []int{1, 2, 3},
			mapFunc:  func(item int, index int) string { return strconv.Itoa(item) },
			expected: []string{"1", "2", "3"},
		},
		{
			input:    []int{4, 5, 6},
			mapFunc:  func(item int, index int) string { return strconv.Itoa(item * 2) },
			expected: []string{"8", "10", "12"},
		},
		// Add more test cases here...
	}

	for _, tc := range testCases {
		result := util.Map(tc.input, tc.mapFunc)

		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("Expected %v, got %v", tc.expected, result)
		}
	}
}
