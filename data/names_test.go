package data

import (
	"testing"
)

func TestGenerateName(t *testing.T) {
	// Test the GenerateName function
	name := GenerateName()
	if len(name) <= 1 {
		t.Errorf("Expected generated name to be non-empty")
	}
}
