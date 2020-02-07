package main

import (
	"testing"
)

// Test that parser returns error when template cannot be found
func TestTemplate_TemplateNotFound(t *testing.T) {
	_, err := NewTemplate("death-star-was-an-inside-job.txt")
	if err == nil {
		t.Fatal("expected an error")
	}
}
