package pkg

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

var (
	NewDecoder = func(file *os.File) *json.Decoder {
		return json.NewDecoder(file)
	}
)

type ConfigItem struct {
	Text       string `json:"text"`
	Keyword    string `json:"keyword"`
	Icon       string `json:"icon"`
	TextColor  string `json:"text_color"`
	ValueColor string `json:"value_color"`
	IconColor  string `json:"icon_color"`
	Value      string `json:"value"`
}

type Config struct {
	Items []ConfigItem `json:"items"`
	Ascii struct {
		Path              string `json:"path"`
		Colors            string `json:"colors"`
		Enabled           bool   `json:"enabled"`
		HorizontalPadding int    `json:"horizontal_padding"`
		VerticalPadding   int    `json:"vertical_padding"`
		Position          string `json:"position"`
	} `json:"ascii"`
	Header struct {
		Text      string `json:"text"`
		TextColor string `json:"text_color"`
		LineColor string `json:"line_color"`
		Line      bool   `json:"line"`
		Enabled   bool   `json:"enabled"`
	} `json:"header"`
	Footer struct {
		Text      string `json:"text"`
		TextColor string `json:"text_color"`
		LineColor string `json:"line_color"`
		Line      bool   `json:"line"`
		Enabled   bool   `json:"enabled"`
	} `json:"footer"`
	General struct {
		MenuType string `json:"menu_type"`
		Columns  bool   `json:"columns"`
	} `json:"general"`
}

func LoadConfig(filename string) (Config, error) {
	var config Config
	file, err := os.Open(filename)
	if err != nil {
		return config, fmt.Errorf("error opening config file: %w", err)
	}
	defer file.Close()

	decoder := NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, fmt.Errorf("error decoding config file: %w", err)
	}

	return config, nil
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
