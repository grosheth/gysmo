package pkg

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
)

func GetColorCode(color string) string {
	switch strings.ToLower(color) {
	case "red":
		return Red
	case "green":
		return Green
	case "yellow":
		return Yellow
	case "blue":
		return Blue
	case "purple":
		return Purple
	case "cyan":
		return Cyan
	case "white":
		return White
	default:
		// Check if the color is an RGB code in the format "#RRGGBB"
		if strings.HasPrefix(color, "#") && len(color) == 7 {
			r, g, b := hexToRGB(color)
			return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
		}
		return Reset
	}
}

func GetEnvVar(envVars []string) string {
	for _, envVar := range envVars {
		if value, exists := os.LookupEnv(envVar); exists {
			return value
		}
	}
	return defaultConfigValue
}

func GetRunningProcess(processes map[string]string) string {
	for process, name := range processes {
		if isProcessRunning(process) {
			return name
		}
	}
	return defaultConfigValue
}

// hexToRGB converts a hex color code to RGB values
func hexToRGB(hex string) (int, int, int) {
	var r, g, b int
	fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	return r, g, b
}

// Abs returns the Absolute value of an integer
func Abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

// charsToString converts a [65]int8 array to a string
func charsToString(ca [65]int8) string {
	s := make([]byte, len(ca))
	for i, v := range ca {
		if v == 0 {
			break
		}
		s[i] = byte(v)
	}
	return strings.TrimRight(string(s), "\x00")
}

// isCommandAvailable checks if a command is available on the system
func isCommandAvailable(name string) bool {
	cmd := exec.Command("which", name)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// isProcessRunning checks if a process with the given name is running
func isProcessRunning(processName string) bool {
	procDir := "/proc"
	entries, err := os.ReadDir(procDir)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if entry.IsDir() {
			cmdlinePath := fmt.Sprintf("%s/%s/cmdline", procDir, entry.Name())
			cmdline, err := os.ReadFile(cmdlinePath)
			if err == nil && strings.Contains(string(cmdline), processName) {
				return true
			}
		}
	}

	return false
}

func ValidateJsonConfig(configPath string, schemaPath string) error {
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + schemaPath)
	documentLoader := gojsonschema.NewReferenceLoader("file://" + configPath)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("error validating config: %v", err)
	}

	if !result.Valid() {
		var errorMessages []string
		for _, desc := range result.Errors() {
			errorMessages = append(errorMessages, fmt.Sprintf("- %s", desc))
		}
		return fmt.Errorf("config file is not valid:\n%s", strings.Join(errorMessages, "\n"))
	}

	return nil
}

func SaveDataToFile(data map[string]string, filename string) error {
	// Read existing data from the file
	existingData := make(map[string]string)
	file, err := os.Open(filename)
	if err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&existingData); err != nil {
			return err
		}
	}

	// Update existing data with new data
	for key, value := range data {
		existingData[key] = value
	}

	// Write updated data back to the file
	file, err = os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(existingData); err != nil {
		return err
	}

	return nil
}
