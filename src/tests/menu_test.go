package tests

import (
	"fmt"
	"gysmo/src/pkg"
	"strings"
	"testing"
)

func printByteValues(s string) {
	for i, c := range s {
		fmt.Printf("Index %d: %q (%d)\n", i, c, c)
	}
}

func TestBuildBoxMenu(t *testing.T) {
	config := pkg.Config{
		Items: []pkg.ConfigItem{
			{Text: "user", Keyword: "user", Icon: "", TextColor: "red", ValueColor: "green", ImageColor: "blue"},
			{Text: "shell", Keyword: "shell", Icon: "", TextColor: "yellow", ValueColor: "blue", ImageColor: "cyan"},
		},
		Ascii: struct {
			Path    string `json:"path"`
			Colors  string `json:"colors"`
			Enabled bool   `json:"enabled"`
			Padding int    `json:"padding"`
		}{
			Path:    "ascii/gysmo",
			Colors:  "red",
			Enabled: true,
			Padding: 2,
		},
		Header: struct {
			Text      string `json:"text"`
			TextColor string `json:"text_color"`
			LineColor string `json:"line_color"`
			Line      bool   `json:"line"`
			Enabled   bool   `json:"enabled"`
		}{
			Text:      "Header",
			TextColor: "white",
			LineColor: "red",
			Line:      true,
			Enabled:   true,
		},
		Footer: struct {
			Text      string `json:"text"`
			TextColor string `json:"text_color"`
			LineColor string `json:"line_color"`
			Line      bool   `json:"line"`
			Enabled   bool   `json:"enabled"`
		}{
			Text:      "Footer",
			TextColor: "white",
			LineColor: "red",
			Line:      true,
			Enabled:   true,
		},
		General: struct {
			MenuType string `json:"menu_type"`
			Columns  bool   `json:"columns"`
		}{
			MenuType: "box",
			Columns:  false,
		},
	}

	items := map[string]string{
		"user":  "testuser",
		"shell": "zsh",
	}

	asciiArt := "ASCII ART"
	borderWidth := pkg.DefineBoxBorder(config)
	asciiPadding := strings.Repeat(" ", config.Ascii.Padding)
	paddedAsciiArt := pkg.AddPaddingToMultilineString(asciiArt, asciiPadding)

	menu := pkg.BuildBoxMenu(items, paddedAsciiArt, config, borderWidth)

	expectedMenu := "    ASCII ART\x1b[31m\x1b[0m\n" +
		"  ╭──────────╮\n" +
		"  │ \x1b[37mHeader  \x1b[0m │\n" +
		"  ├\x1b[31m──────────\x1b[0m┤\n" +
		"  │ \x1b[34m\x1b[0m  \x1b[31muser\x1b[0m  │ \x1b[32mtestuser\x1b[0m\n" +
		"  │ \x1b[36m\x1b[0m  \x1b[33mshell\x1b[0m │ \x1b[34mzsh\x1b[0m\n" +
		"  ├\x1b[31m──────────\x1b[0m┤\n" +
		"  │ \x1b[37mFooter  \x1b[0m │\n" +
		"  ╰──────────╯\n"

	if strings.TrimSpace(menu) != strings.TrimSpace(expectedMenu) {
		t.Errorf("Expected menu:\n%s\nGot:\n%s", expectedMenu, menu)
		printByteValues(expectedMenu)
		printByteValues(menu)
	}
}

func TestBuildListMenu(t *testing.T) {
	config := pkg.Config{
		Items: []pkg.ConfigItem{
			{Text: "user", Keyword: "user", Icon: "", TextColor: "red", ValueColor: "green", ImageColor: "blue"},
			{Text: "shell", Keyword: "shell", Icon: "", TextColor: "yellow", ValueColor: "blue", ImageColor: "cyan"},
		},
		Ascii: struct {
			Path    string `json:"path"`
			Colors  string `json:"colors"`
			Enabled bool   `json:"enabled"`
			Padding int    `json:"padding"`
		}{
			Path:    "ascii/gysmo",
			Colors:  "red",
			Enabled: true,
			Padding: 2,
		},
		Header: struct {
			Text      string `json:"text"`
			TextColor string `json:"text_color"`
			LineColor string `json:"line_color"`
			Line      bool   `json:"line"`
			Enabled   bool   `json:"enabled"`
		}{
			Text:      "Header",
			TextColor: "white",
			LineColor: "red",
			Line:      true,
			Enabled:   true,
		},
		Footer: struct {
			Text      string `json:"text"`
			TextColor string `json:"text_color"`
			LineColor string `json:"line_color"`
			Line      bool   `json:"line"`
			Enabled   bool   `json:"enabled"`
		}{
			Text:      "Footer",
			TextColor: "white",
			LineColor: "red",
			Line:      true,
			Enabled:   true,
		},
		General: struct {
			MenuType string `json:"menu_type"`
			Columns  bool   `json:"columns"`
		}{
			MenuType: "list",
			Columns:  false,
		},
	}

	items := map[string]string{
		"user":  "testuser",
		"shell": "zsh",
	}

	asciiArt := "    ASCII ART"
	borderWidth := pkg.DefineBoxBorder(config)
	menu := pkg.BuildListMenu(items, asciiArt, config, borderWidth)

	expectedMenu := "    ASCII ART\x1b[31m\x1b[0m\n" +
		"   \x1b[37mHeader  \x1b[0m \n" +
		"  \x1b[31m──────────\x1b[0m\n" +
		"   \x1b[34m\x1b[0m \x1b[31muser\x1b[0m   \x1b[32mtestuser\x1b[0m\n" +
		"   \x1b[36m\x1b[0m \x1b[33mshell\x1b[0m  \x1b[34mzsh\x1b[0m\n" +
		"  \x1b[31m──────────\x1b[0m\n" +
		"   \x1b[37mFooter  \x1b[0m \n"

	if strings.TrimSpace(menu) != strings.TrimSpace(expectedMenu) {
		t.Errorf("Expected menu:\n%s\nGot:\n%s", expectedMenu, menu)
		printByteValues(expectedMenu)
		printByteValues(menu)
	}
}
