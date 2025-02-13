package pkg

import (
	"encoding/json"
	"fmt"
	"os"
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
	file, err := OpenFile(filename)
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
