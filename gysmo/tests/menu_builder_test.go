package tests

import (
	"fmt"
	"gysmo/gysmo/pkg"
	"strings"
	"testing"
)

func printByteValues(s string) {
	for i, c := range s {
		fmt.Printf("Index %d: %q (%d)\n", i, c, c)
	}
}

func normalizeString(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func TestBuildBoxMenuWithAsciiTop(t *testing.T) {
	config := GetConfigWithAscii("top")
	items := GetTestItems()
	asciiArt := "ASCII ART"
	menu := pkg.BuildBoxMenu(items, asciiArt, config)

	expectedMenu := "  ASCII ART\n" +
		"\n" +
		" ╭──────────╮\n" +
		" │ Header   │\n" +
		" ├──────────┤\n" +
		" │   user  │ testuser\n" +
		" │   shell │ zsh\n" +
		" ├──────────┤\n" +
		" │ Footer   │\n" +
		" ╰──────────╯"

	// Strip ANSI codes and normalize strings
	normalizedMenu := normalizeString(pkg.StripAnsiCodes(menu))
	normalizedExpectedMenu := normalizeString(pkg.StripAnsiCodes(expectedMenu))

	if normalizedMenu != normalizedExpectedMenu {
		t.Errorf("Expected menu:\n%s\nGot:\n%s", expectedMenu, menu)
		printByteValues(menu)
	}
}

func TestBuildBoxMenuWithAsciiBottom(t *testing.T) {
	config := GetConfigWithAscii("bottom")
	items := GetTestItems()
	asciiArt := "ASCII ART"
	menu := pkg.BuildBoxMenu(items, asciiArt, config)

	expectedMenu := " ╭──────────╮\n" +
		" │ Header   │\n" +
		" ├──────────┤\n" +
		" │   user  │ testuser\n" +
		" │   shell │ zsh\n" +
		" ├──────────┤\n" +
		" │ Footer   │\n" +
		" ╰──────────╯\n" +
		"  ASCII ART\n"

	// Strip ANSI codes and normalize strings
	normalizedMenu := normalizeString(pkg.StripAnsiCodes(menu))
	normalizedExpectedMenu := normalizeString(pkg.StripAnsiCodes(expectedMenu))

	if normalizedMenu != normalizedExpectedMenu {
		t.Errorf("Expected menu:\n%s\nGot:\n%s", expectedMenu, menu)
		printByteValues(menu)
	}
}

func TestBuildBoxMenuWithAsciiRight(t *testing.T) {
	config := GetConfigWithAscii("right")
	items := GetTestItems()
	asciiArt := "ASCII ART"
	menu := pkg.BuildBoxMenu(items, asciiArt, config)

	expectedMenu := " ╭──────────╮\n" +
		" │ Header   │             ASCII ART\n" +
		" ├──────────┤\n" +
		" │   user  │ testuser\n" +
		" │   shell │ zsh\n" +
		" ├──────────┤\n" +
		" │ Footer   │\n" +
		" ╰──────────╯\n"

	// Strip ANSI codes and normalize strings
	normalizedMenu := normalizeString(pkg.StripAnsiCodes(menu))
	normalizedExpectedMenu := normalizeString(pkg.StripAnsiCodes(expectedMenu))

	if normalizedMenu != normalizedExpectedMenu {
		t.Errorf("Expected menu:\n%s\nGot:\n%s", expectedMenu, menu)
		printByteValues(menu)
	}
}

func TestBuildBoxMenuWithAsciiLeft(t *testing.T) {
	config := GetConfigWithAscii("left")
	items := GetTestItems()
	asciiArt := "ASCII ART"
	menu := pkg.BuildBoxMenu(items, asciiArt, config)

	expectedMenu := "                        ╭──────────╮\n" +
		"          ASCII ART     │ Header   │\n" +
		"                        ├──────────┤\n" +
		"                        │   user  │ testuser\n" +
		"                        │   shell │ zsh\n" +
		"                        ├──────────┤\n" +
		"                        │ Footer   │\n" +
		"                        ╰──────────╯\n"

	// Strip ANSI codes and normalize strings
	normalizedMenu := normalizeString(pkg.StripAnsiCodes(menu))
	normalizedExpectedMenu := normalizeString(pkg.StripAnsiCodes(expectedMenu))

	if normalizedMenu != normalizedExpectedMenu {
		t.Errorf("Expected menu:\n%s\nGot:\n%s", expectedMenu, menu)
		printByteValues(menu)
	}
}

func TestBuildBoxMenuWithAsciiPadding(t *testing.T) {
	config := GetConfigWithAscii("top")
	config.Ascii.HorizontalPadding = 4
	config.Ascii.VerticalPadding = 2
	items := GetTestItems()
	asciiArt := "ASCII ART"
	menu := pkg.BuildBoxMenu(items, asciiArt, config)

	expectedMenu := "\n\n" +
		"    ASCII ART    \n\n" +
		" ╭──────────╮\n" +
		" │ Header   │\n" +
		" ├──────────┤\n" +
		" │   user  │ testuser\n" +
		" │   shell │ zsh\n" +
		" ├──────────┤\n" +
		" │ Footer   │\n" +
		" ╰──────────╯\n"

	// Strip ANSI codes and normalize strings
	normalizedMenu := normalizeString(pkg.StripAnsiCodes(menu))
	normalizedExpectedMenu := normalizeString(pkg.StripAnsiCodes(expectedMenu))

	if normalizedMenu != normalizedExpectedMenu {
		t.Errorf("Expected menu:\n%s\nGot:\n%s", expectedMenu, menu)
		printByteValues(menu)
	}
}

func TestBuildListMenuWithAsciiTop(t *testing.T) {
	config := GetConfigWithAscii("top")
	items := GetTestItems()
	asciiArt := "ASCII ART"
	menu := pkg.BuildListMenu(items, asciiArt, config)

	expectedMenu := "  ASCII ART\n" +
		"\n" +
		"Header\n" +
		"──────────\n" +
		"  user   testuser\n" +
		"  shell  zsh\n" +
		"──────────\n" +
		"Footer\n"

	// Strip ANSI codes and normalize strings
	normalizedMenu := normalizeString(pkg.StripAnsiCodes(menu))
	normalizedExpectedMenu := normalizeString(pkg.StripAnsiCodes(expectedMenu))

	if normalizedMenu != normalizedExpectedMenu {
		t.Errorf("Expected menu:\n%s\nGot:\n%s", expectedMenu, menu)
		printByteValues(menu)
	}
}

func TestBuildListMenuWithAsciiBottom(t *testing.T) {
	config := GetConfigWithAscii("bottom")
	items := GetTestItems()
	asciiArt := "ASCII ART"
	menu := pkg.BuildListMenu(items, asciiArt, config)

	expectedMenu := "Header\n" +
		"──────────\n" +
		"  user   testuser\n" +
		"  shell  zsh\n" +
		"──────────\n" +
		"Footer\n" +
		"  ASCII ART\n"

	// Strip ANSI codes and normalize strings
	normalizedMenu := normalizeString(pkg.StripAnsiCodes(menu))
	normalizedExpectedMenu := normalizeString(pkg.StripAnsiCodes(expectedMenu))

	if normalizedMenu != normalizedExpectedMenu {
		t.Errorf("Expected menu:\n%s\nGot:\n%s", expectedMenu, menu)
		printByteValues(menu)
	}
}

func TestBuildListMenuWithAsciiRight(t *testing.T) {
	config := GetConfigWithAscii("right")
	items := GetTestItems()
	asciiArt := "ASCII ART"
	menu := pkg.BuildListMenu(items, asciiArt, config)

	expectedMenu := "Header\n" +
		"──────────                      ASCII ART\n" +
		"  user   testuser\n" +
		"  shell  zsh\n" +
		"──────────\n" +
		"Footer\n"

	// Strip ANSI codes and normalize strings
	normalizedMenu := normalizeString(pkg.StripAnsiCodes(menu))
	normalizedExpectedMenu := normalizeString(pkg.StripAnsiCodes(expectedMenu))

	if normalizedMenu != normalizedExpectedMenu {
		t.Errorf("Expected menu:\n%s\nGot:\n%s", expectedMenu, menu)
		printByteValues(menu)
	}
}

func TestBuildListMenuWithAsciiLeft(t *testing.T) {
	config := GetConfigWithAscii("left")
	items := GetTestItems()
	asciiArt := "ASCII ART"
	menu := pkg.BuildListMenu(items, asciiArt, config)

	expectedMenu := "                       Header\n" +
		"          ASCII ART   ──────────\n" +
		"                         user   testuser\n" +
		"                         shell  zsh\n" +
		"                       ──────────\n" +
		"                       Footer\n"

	// Strip ANSI codes and normalize strings
	normalizedMenu := normalizeString(pkg.StripAnsiCodes(menu))
	normalizedExpectedMenu := normalizeString(pkg.StripAnsiCodes(expectedMenu))

	if normalizedMenu != normalizedExpectedMenu {
		t.Errorf("Expected menu:\n%s\nGot:\n%s", expectedMenu, menu)
		printByteValues(menu)
	}
}

func TestBuildListMenuWithColumns(t *testing.T) {
	config := GetConfigWithHeader()
	config.General.Columns = true
	items := GetTestItems()
	asciiArt := "ASCII ART"
	menu := pkg.BuildListMenu(items, asciiArt, config)

	expectedMenu := "  ASCII ART\n" +
		"Header\n" +
		"──────────\n" +
		"  user   testuser |  shell  zsh\n" +
		"──────────\n" +
		"Footer\n"

	// Strip ANSI codes and normalize strings
	normalizedMenu := normalizeString(pkg.StripAnsiCodes(menu))
	normalizedExpectedMenu := normalizeString(pkg.StripAnsiCodes(expectedMenu))

	if normalizedMenu != normalizedExpectedMenu {
		t.Errorf("Expected menu:\n%s\nGot:\n%s", expectedMenu, menu)
		printByteValues(menu)
	}
}

// Helper functions to Get configurations and items for tests

func GetConfigWithAscii(position string) pkg.Config {
	return pkg.Config{
		Items: []pkg.ConfigItem{
			{Text: "user", Keyword: "user", Icon: "", TextColor: "", ValueColor: "", IconColor: ""},
			{Text: "shell", Keyword: "shell", Icon: "", TextColor: "", ValueColor: "", IconColor: ""},
		},
		Ascii: struct {
			Path              string `json:"path"`
			Colors            string `json:"colors"`
			Enabled           bool   `json:"enabled"`
			HorizontalPadding int    `json:"horizontal_padding"`
			VerticalPadding   int    `json:"vertical_padding"`
			Position          string `json:"position"`
		}{
			Path:              "ascii/gysmo",
			Colors:            "",
			Enabled:           true,
			HorizontalPadding: 2,
			VerticalPadding:   1,
			Position:          position,
		},
		Header: struct {
			Text      string `json:"text"`
			TextColor string `json:"text_color"`
			LineColor string `json:"line_color"`
			Line      bool   `json:"line"`
			Enabled   bool   `json:"enabled"`
		}{
			Text:      "Header",
			TextColor: "",
			LineColor: "",
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
			TextColor: "",
			LineColor: "",
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
}

func GetConfigWithHeader() pkg.Config {
	return pkg.Config{
		Items: []pkg.ConfigItem{
			{Text: "user", Keyword: "user", Icon: "", TextColor: "red", ValueColor: "green", IconColor: "blue"},
			{Text: "shell", Keyword: "shell", Icon: "", TextColor: "yellow", ValueColor: "blue", IconColor: "cyan"},
		},
		Ascii: struct {
			Path              string `json:"path"`
			Colors            string `json:"colors"`
			Enabled           bool   `json:"enabled"`
			HorizontalPadding int    `json:"horizontal_padding"`
			VerticalPadding   int    `json:"vertical_padding"`
			Position          string `json:"position"`
		}{
			Path:              "ascii/gysmo",
			Colors:            "red",
			Enabled:           true,
			HorizontalPadding: 2,
			VerticalPadding:   1,
			Position:          "top",
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
}

func GetConfigWithFooter() pkg.Config {
	return pkg.Config{
		Items: []pkg.ConfigItem{
			{Text: "user", Keyword: "user", Icon: "", TextColor: "red", ValueColor: "green", IconColor: "blue"},
			{Text: "shell", Keyword: "shell", Icon: "", TextColor: "yellow", ValueColor: "blue", IconColor: "cyan"},
		},
		Ascii: struct {
			Path              string `json:"path"`
			Colors            string `json:"colors"`
			Enabled           bool   `json:"enabled"`
			HorizontalPadding int    `json:"horizontal_padding"`
			VerticalPadding   int    `json:"vertical_padding"`
			Position          string `json:"position"`
		}{
			Path:              "ascii/gysmo",
			Colors:            "red",
			Enabled:           true,
			HorizontalPadding: 2,
			VerticalPadding:   1,
			Position:          "top",
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
}

func GetTestItems() map[string]string {
	return map[string]string{
		"user":  "testuser",
		"shell": "zsh",
	}
}

func GetTestItemsWithMoreThan10Items() map[string]string {
	items := map[string]string{
		"user":   "testuser",
		"shell":  "zsh",
		"item1":  "value1",
		"item2":  "value2",
		"item3":  "value3",
		"item4":  "value4",
		"item5":  "value5",
		"item6":  "value6",
		"item7":  "value7",
		"item8":  "value8",
		"item9":  "value9",
		"item10": "value10",
		"item11": "value11",
	}
	return items
}

func GetTestItemsWithLongText() map[string]string {
	return map[string]string{
		"user":     "testuser",
		"shell":    "zsh",
		"longtext": "This is a very long text item that should be tested for proper handling in the menu",
	}
}

func GetTestItemsWithLongValue() map[string]string {
	return map[string]string{
		"user":      "testuser",
		"shell":     "zsh",
		"longvalue": "This is a very long value for the item that should be tested for proper handling in the menu",
	}
}
