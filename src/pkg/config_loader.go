package pkg

import (
	"encoding/json"
	"fmt"
	"os"

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
		MenuType    string `json:"menu_type"`
		Columns     bool   `json:"columns"`
		MenuPadding int    `json:"menu_padding"`
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

func ValidateJsonConfig(configPath, schemaPath string) error {
	configLoader := gojsonschema.NewReferenceLoader("file://" + configPath)
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + schemaPath)

	result, err := gojsonschema.Validate(schemaLoader, configLoader)
	if err != nil {
		return fmt.Errorf("failed to load schema or config file: %v", err)
	}

	if !result.Valid() {
		var errorMessages string
		for _, desc := range result.Errors() {
			switch desc.Type() {
			case "required":
				errorMessages += fmt.Sprintf("Missing required field: %s\n", desc.Field())
			case "number_one_of":
				errorMessages += fmt.Sprintf("Field %s You need to specify either Keyword or Value for an item.\n", desc.Field())
			default:
				errorMessages += fmt.Sprintf("Validation error on field %s: %s\n", desc.Field(), desc.Description())
			}
		}
		return fmt.Errorf("config file is not valid:\n%s", errorMessages)
	}

	return nil
}
