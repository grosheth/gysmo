package pkg

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadAsciiArt(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("error opening ASCII art file: %w", err)
	}
	defer file.Close()

	var asciiArt strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		asciiArt.WriteString(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading ASCII art file: %w", err)
	}

	return asciiArt.String(), nil
}
