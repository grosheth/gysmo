package tests

import (
	"encoding/json"
	"fmt"
	"gysmo/gysmo/pkg"
	"io"
	"os"
	"testing"
)

// Mock implementation of os.Open for testing
func mockOpen(name string) (*os.File, error) {
	if name == "config/config.json" {
		return os.Open("tests/config/config.json")
	}
	return nil, fmt.Errorf("file not found")
}

// Mock implementation of json.NewDecoder for testing
type mockReader struct {
	data []byte
	pos  int
}

func (r *mockReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

func mockNewDecoder(file *os.File) *json.Decoder {
	data, _ := os.ReadFile(file.Name())
	return json.NewDecoder(&mockReader{data: data})
}

func TestLoadConfig(t *testing.T) {
	// Override the actual functions with the mock functions
	pkg.OpenFile = mockOpen
	pkg.NewDecoder = mockNewDecoder

	config, err := pkg.LoadConfig("config/config.json")
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	if len(config.Items) == 0 {
		t.Errorf("Expected items to be loaded, but got none")
	}

	expectedAsciiPath := "ascii/gysmo"
	if config.Ascii.Path != expectedAsciiPath {
		t.Errorf("Expected ASCII path to be %s, but got %s", expectedAsciiPath, config.Ascii.Path)
	}
}
