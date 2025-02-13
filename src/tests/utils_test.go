package tests

import (
	"gysmo/src/pkg"
	"os"
	"strings"
	"testing"
)

type MockUtils struct{}

func (m MockUtils) GetColorCode(color string) string {
	switch color {
	case "red":
		return pkg.Red
	case "green":
		return pkg.Green
	case "yellow":
		return pkg.Yellow
	case "blue":
		return pkg.Blue
	case "purple":
		return pkg.Purple
	case "cyan":
		return pkg.Cyan
	case "white":
		return pkg.White
	case "#FF5733":
		return "\033[38;2;255;87;51m"
	default:
		return pkg.Reset
	}
}

func (m MockUtils) GetEnvVar(envVars []string) string {
	for _, envVar := range envVars {
		if value, exists := os.LookupEnv(envVar); exists {
			return value
		}
	}
	return "Replace This"
}

func (m MockUtils) GetRunningProcess(processes map[string]string) string {
	for process, name := range processes {
		if process == "mock_process" {
			return name
		}
	}
	return "Replace This"
}

func (m MockUtils) GetRunningProcessesCount() int {
	return 2
}

func (m MockUtils) HexToRGB(hex string) (int, int, int) {
	switch hex {
	case "#FF5733":
		return 255, 87, 51
	case "#00FF00":
		return 0, 255, 0
	case "#0000FF":
		return 0, 0, 255
	default:
		return 0, 0, 0
	}
}

func (m MockUtils) Abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func (m MockUtils) CharsToString(ca [65]int8) string {
	s := make([]byte, len(ca))
	for i, v := range ca {
		if v == 0 {
			break
		}
		s[i] = byte(v)
	}
	return strings.TrimRight(string(s), "\x00")
}

func (m MockUtils) IsCommandAvailable(name string) bool {
	return name == "mock_command"
}

func (m MockUtils) IsProcessRunning(processName string) bool {
	return processName == "mock_process"
}

func (m MockUtils) ValidateJsonConfig(configPath string, schemaPath string) error {
	return nil
}

func TestGetColorCode(t *testing.T) {
	utils := MockUtils{}
	tests := []struct {
		color    string
		expected string
	}{
		{"red", pkg.Red},
		{"green", pkg.Green},
		{"yellow", pkg.Yellow},
		{"blue", pkg.Blue},
		{"purple", pkg.Purple},
		{"cyan", pkg.Cyan},
		{"white", pkg.White},
		{"#FF5733", "\033[38;2;255;87;51m"},
		{"unknown", pkg.Reset},
	}

	for _, test := range tests {
		result := utils.GetColorCode(test.color)
		if result != test.expected {
			t.Errorf("For color %s, expected %s, but got %s", test.color, test.expected, result)
		}
	}
}

func TestGetEnvVar(t *testing.T) {
	utils := MockUtils{}
	os.Setenv("TEST_ENV_VAR", "test_value")
	defer os.Unsetenv("TEST_ENV_VAR")

	tests := []struct {
		envVars  []string
		expected string
	}{
		{[]string{"TEST_ENV_VAR"}, "test_value"},
		{[]string{"NON_EXISTENT_VAR"}, "Replace This"},
	}

	for _, test := range tests {
		result := utils.GetEnvVar(test.envVars)
		if result != test.expected {
			t.Errorf("For envVars %v, expected %s, but got %s", test.envVars, test.expected, result)
		}
	}
}

func TestGetRunningProcess(t *testing.T) {
	utils := MockUtils{}
	tests := []struct {
		processes map[string]string
		expected  string
	}{
		{map[string]string{"mock_process": "Mock Process"}, "Mock Process"},
		{map[string]string{"non_existent_process": "Non Existent Process"}, "Replace This"},
	}

	for _, test := range tests {
		result := utils.GetRunningProcess(test.processes)
		if result != test.expected {
			t.Errorf("For processes %v, expected %s, but got %s", test.processes, test.expected, result)
		}
	}
}

func TestGetRunningProcessesCount(t *testing.T) {
	utils := MockUtils{}
	expected := 2
	result := utils.GetRunningProcessesCount()
	if result != expected {
		t.Errorf("Expected %d, but got %d", expected, result)
	}
}

func TestHexToRGB(t *testing.T) {
	utils := MockUtils{}
	tests := []struct {
		hex      string
		expected [3]int
	}{
		{"#FF5733", [3]int{255, 87, 51}},
		{"#00FF00", [3]int{0, 255, 0}},
		{"#0000FF", [3]int{0, 0, 255}},
	}

	for _, test := range tests {
		r, g, b := utils.HexToRGB(test.hex)
		if r != test.expected[0] || g != test.expected[1] || b != test.expected[2] {
			t.Errorf("For hex %s, expected (%d, %d, %d), but got (%d, %d, %d)", test.hex, test.expected[0], test.expected[1], test.expected[2], r, g, b)
		}
	}
}

func TestAbs(t *testing.T) {
	utils := MockUtils{}
	tests := []struct {
		input    int64
		expected int64
	}{
		{-5, 5},
		{5, 5},
		{0, 0},
	}

	for _, test := range tests {
		result := utils.Abs(test.input)
		if result != test.expected {
			t.Errorf("For input %d, expected %d, but got %d", test.input, test.expected, result)
		}
	}
}

func TestCharsToString(t *testing.T) {
	utils := MockUtils{}
	input := [65]int8{'H', 'e', 'l', 'l', 'o', 0}
	expected := "Hello"
	result := utils.CharsToString(input)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestIsCommandAvailable(t *testing.T) {
	utils := MockUtils{}
	result := utils.IsCommandAvailable("mock_command")
	if !result {
		t.Errorf("Expected command to be available, but it was not")
	}
}

func TestIsProcessRunning(t *testing.T) {
	utils := MockUtils{}
	result := utils.IsProcessRunning("mock_process")
	if !result {
		t.Errorf("Expected process to be running, but it was not")
	}
}

// func TestValidateJsonConfig(t *testing.T) {
// 	utils := MockUtils{}
// 	err := utils.ValidateJsonConfig("config/config.json", "config/config_schema.json")
// 	if err != nil {
// 		t.Errorf("Expected no error, but got %v", err)
// 	}
// }

// Mock DirEntry implementation
type mockDirEntry struct {
	name  string
	isDir bool
}

func (m mockDirEntry) Name() string               { return m.name }
func (m mockDirEntry) IsDir() bool                { return m.isDir }
func (m mockDirEntry) Type() os.FileMode          { return 0 }
func (m mockDirEntry) Info() (os.FileInfo, error) { return nil, nil }
