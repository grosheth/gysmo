package src

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
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
	case "":
		return Reset
	default:
		// Check if the color is an RGB code in the format "#RRGGBB"
		if strings.HasPrefix(color, "#") && len(color) == 7 {
			r, g, b := HexToRGB(color)
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
		if IsProcessRunning(process) {
			return name
		}
	}
	return defaultConfigValue
}

func HexToRGB(hex string) (int, int, int) {
	var r, g, b int
	fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	return r, g, b
}

func Abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func CharsToString(ca [65]int8) string {
	s := make([]byte, len(ca))
	for i, v := range ca {
		if v == 0 {
			break
		}
		s[i] = byte(v)
	}
	return strings.TrimRight(string(s), "\x00")
}

func IsCommandAvailable(name string) bool {
	cmd := exec.Command("which", name)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func IsProcessRunning(processName string) bool {
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

func VisibleChars(s string) string {
	var builder strings.Builder
	for _, r := range s {
		switch r {
		case '\n':
			builder.WriteString("$\n")
		case '\t':
			builder.WriteString("^I")
		case '\r':
			builder.WriteString("^M")
		default:
			if r < 32 || r == 127 {
				builder.WriteString(fmt.Sprintf("^%c", r+64))
			} else {
				builder.WriteRune(r)
			}
		}
	}
	return builder.String()
}

func StripAnsiCodes(str string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(str, "")
}

func SaveDataToFile(data map[string]string) error {
	// Read existing data from the file
	workingPath := LoadWorkingPath()
	dataPath := filepath.Join(workingPath, "data", "data.json")

	existingData := make(map[string]string)
	file, err := os.Open(dataPath)
	if err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&existingData); err != nil {
			return err
		}
	}

	for key, value := range data {
		existingData[key] = value
	}

	file, err = os.Create(dataPath)
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

func LoadWorkingPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
	}
	workingPath := filepath.Join(homeDir, ".config", "gysmo")
	return workingPath
}

func IsHeaderOrFooter(line string, config Config, index int, totalLines int) string {
	if index == 0 && strings.Contains(line, config.Header.Text) {
		return "Header"
	}
	if index != 0 && strings.Contains(line, config.Footer.Text) {
		return "Footer"
	}
	return "Content"
}

// validate if string is a footer or header line
func IsLine(s string) bool {
	hasDash := false
	s = strings.TrimSpace(s)
	s = StripAnsiCodes(s)
	for _, char := range s {
		if char != '─' {
			return false
		}
		if char == '─' {
			hasDash = true
		}
	}
	return hasDash
}

// New functions for copying files and checking their existence
func CopyFile(gysmo, dst string) error {
	sourceFile, err := os.Open(gysmo)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	return err
}

func EnsureConfigFilesExist() error {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("Failed to get user home directory: %v", err))
	}

	configDir := filepath.Join(homeDir, ".config", "gysmo")
	if os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return err
		}
	}

	shareDir := ""
	if shareDir := os.Getenv("GYSMO_SHARE_PATH"); shareDir == "" {
		shareDir = "/usr/share/gysmo"
	}

	files := []string{"config/config.json", "config/schema/config_schema.json", "ascii/gysmo"}

	for _, file := range files {
		configFilePath := filepath.Join(configDir, file)
		shareFilePath := filepath.Join(shareDir, file)

		if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
			if err := CopyFile(shareFilePath, configFilePath); err != nil {
				return err
			}
		}
	}

	return nil
}
