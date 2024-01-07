package util_test

import (
	"testing"

	"github.com/jo-fr/activityhub/backend/pkg/util"
)

func TestFromPointer(t *testing.T) {
	// Test case 1: Pointer is nil
	var ptr *int
	expected1 := 0
	result1 := util.FromPointer(ptr)
	if result1 != expected1 {
		t.Errorf("Expected %v, got %v", expected1, result1)
	}

	// Test case 2: Pointer is not nil
	value := 42
	ptr = &value
	expected2 := 42
	result2 := util.FromPointer(ptr)
	if result2 != expected2 {
		t.Errorf("Expected %v, got %v", expected2, result2)
	}
}
