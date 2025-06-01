package tests

import (
	"gysmo/gysmo/src"
	"os"
	"testing"
)

func TestGetColorCode(t *testing.T) {
	tests := []struct {
		color    string
		expected string
	}{
		{"red", src.Red},
		{"green", src.Green},
		{"yellow", src.Yellow},
		{"blue", src.Blue},
		{"purple", src.Purple},
		{"cyan", src.Cyan},
		{"white", src.White},
		{"#FF5733", "\033[38;2;255;87;51m"},
		{"unknown", src.Reset},
	}

	for _, test := range tests {
		result := src.GetColorCode(test.color)
		if result != test.expected {
			t.Errorf("For color %s, expected %s, but got %s", test.color, test.expected, result)
		}
	}
}

func TestGetEnvVar(t *testing.T) {
	os.Setenv("TEST_ENV_VAR", "test_value")
	defer os.Unsetenv("TEST_ENV_VAR")

	tests := []struct {
		envVars  []string
		expected string
	}{
		{[]string{"TEST_ENV_VAR"}, "test_value"},
		{[]string{"NON_EXISTENT_VAR"}, "Not Found"},
	}

	for _, test := range tests {
		result := src.GetEnvVar(test.envVars)
		if result != test.expected {
			t.Errorf("For envVars %v, expected %s, but got %s", test.envVars, test.expected, result)
		}
	}
}

func TestHexToRGB(t *testing.T) {
	tests := []struct {
		hex      string
		expected [3]int
	}{
		{"#FF5733", [3]int{255, 87, 51}},
		{"#00FF00", [3]int{0, 255, 0}},
		{"#0000FF", [3]int{0, 0, 255}},
	}

	for _, test := range tests {
		r, g, b := src.HexToRGB(test.hex)
		if r != test.expected[0] || g != test.expected[1] || b != test.expected[2] {
			t.Errorf("For hex %s, expected (%d, %d, %d), but got (%d, %d, %d)", test.hex, test.expected[0], test.expected[1], test.expected[2], r, g, b)
		}
	}
}

func TestAbs(t *testing.T) {
	tests := []struct {
		input    int64
		expected int64
	}{
		{-5, 5},
		{5, 5},
		{0, 0},
	}

	for _, test := range tests {
		result := src.Abs(test.input)
		if result != test.expected {
			t.Errorf("For input %d, expected %d, but got %d", test.input, test.expected, result)
		}
	}
}

func TestCharsToString(t *testing.T) {
	input := [65]int8{'H', 'e', 'l', 'l', 'o', 0}
	expected := "Hello"
	result := src.CharsToString(input)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}
