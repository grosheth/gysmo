package tests

import (
	"fmt"
	"gysmo/gysmo/pkg"
	"os"
	"strings"
	"testing"
)

// Mock implementation of os.Open for testing
func mockOpenAscii(name string) (*os.File, error) {
	if name == "ascii/gysmo" {
		return os.Open("tests/ascii/gysmo")
	}
	return nil, fmt.Errorf("file not found")
}

func TestReadAsciiArt(t *testing.T) {
	// Override the actual functions with the mock functions
	pkg.OpenFile = mockOpenAscii

	asciiArt, err := pkg.ReadAsciiArt("ascii/gysmo")
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	if !strings.Contains(asciiArt, "_____") {
		t.Errorf("Expected ASCII art to contain '_____', but got %s", asciiArt)
	}
}
