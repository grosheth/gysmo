package pkg

import (
	"fmt"
	"strings"
)

const padding = 2

func AddPaddingToMultilineString(s string, horizontalPadding int, verticalPadding int) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.Repeat(" ", horizontalPadding) + line + strings.Repeat(" ", horizontalPadding)
	}
	paddingLines := strings.Repeat("\n", verticalPadding)
	return paddingLines + strings.Join(lines, "\n") + paddingLines
}

// DefineBorder calculates the border width for the menu
func DefineBoxBorder(config Config) int {
	maxLength := 0

	// Find the maximum image length can be weird
	minIconLength := 1000
	for _, item := range config.Items {
		imageLength := len(item.Icon)
		if imageLength < minIconLength {
			minIconLength = imageLength
		}
	}

	for _, item := range config.Items {
		itemLength := len(item.Text) + minIconLength + padding
		if itemLength > maxLength {
			maxLength = itemLength
		}
	}

	if config.Header.Enabled {
		headerLength := len(config.Header.Text)
		if headerLength > maxLength {
			maxLength = headerLength + padding
		}
	}

	if config.Footer.Enabled {
		footerLength := len(config.Footer.Text)
		if footerLength > maxLength {
			maxLength = footerLength + padding
		}
	}

	return maxLength
}

func BuildBoxMenu(items map[string]string, asciiArt string, config Config, borderWidth int) string {
	border := strings.Repeat("─", borderWidth)
	menu := ""

	asciiColors := GetColorCode(config.Ascii.Colors)
	paddedAsciiArt := AddPaddingToMultilineString(asciiArt, config.Ascii.HorizontalPadding, config.Ascii.VerticalPadding)

	asciiLines := strings.Split(paddedAsciiArt, "\n")
	maxAsciiWidth := getMaxLineWidth(asciiLines)

	if config.Ascii.Position == "top" {
		menu += fmt.Sprintf("%s%s%s\n", asciiColors, paddedAsciiArt, Reset)
		menu += fmt.Sprintf("  ╭%s╮\n", border)
	} else {
		menu += fmt.Sprintf("  ╭%s╮\n", border)
	}

	minIconLength := getMinIconLength(config.Items)

	if config.Header.Enabled {
		menu += buildHeader(config, borderWidth, border)
	}

	menu += buildMenuItems(config, items, borderWidth, minIconLength)

	if config.Footer.Enabled {
		menu += buildFooter(config, borderWidth, border)
	}

	if config.Ascii.Position == "bottom" {
		menu += fmt.Sprintf("  ╰%s╯\n", border)
		menu += fmt.Sprintf("%s%s%s\n", asciiColors, paddedAsciiArt, Reset)
	} else {
		menu += fmt.Sprintf("  ╰%s╯\n", border)
	}

	if config.Ascii.Position == "left" || config.Ascii.Position == "right" {
		menu = combineAsciiAndMenu(menu, paddedAsciiArt, asciiColors, config.Ascii.Position, maxAsciiWidth)
	}

	return menu
}

func BuildListMenu(items map[string]string, asciiArt string, config Config, borderWidth int) string {
	asciiColors := GetColorCode(config.Ascii.Colors)

	paddedAsciiArt := AddPaddingToMultilineString(asciiArt, config.Ascii.HorizontalPadding, config.Ascii.VerticalPadding)
	border := strings.Repeat("─", borderWidth)

	asciiLines := strings.Split(paddedAsciiArt, "\n")
	maxAsciiWidth := getMaxLineWidth(asciiLines)

	menu := ""

	if config.Ascii.Position == "top" {
		menu += fmt.Sprintf("%s%s%s\n", asciiColors, paddedAsciiArt, Reset)
	}

	if config.Header.Enabled {
		menu += buildHeader(config, borderWidth, border)
	}

	minIconLength := getMinIconLength(config.Items)

	formattedItems := formatMenuItems(config, items, borderWidth, minIconLength)

	if config.General.Columns {
		menu += buildColumns(formattedItems)
	} else {
		for _, item := range formattedItems {
			menu += fmt.Sprintf("%s\n", item)
		}
	}

	if config.Footer.Enabled {
		menu += buildFooter(config, borderWidth, border)
	}

	if config.Ascii.Position == "bottom" {
		menu += fmt.Sprintf("%s%s%s\n", asciiColors, paddedAsciiArt, Reset)
	}

	menu = combineAsciiAndMenu(menu, paddedAsciiArt, asciiColors, config.Ascii.Position, maxAsciiWidth)

	return menu
}

func getMaxLineWidth(lines []string) int {
	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}
	return maxWidth
}

func getMinIconLength(items []ConfigItem) int {
	minIconLength := 1000
	for _, item := range items {
		imageLength := len(item.Icon)
		if imageLength < minIconLength {
			minIconLength = imageLength
		}
	}
	return minIconLength
}

func buildHeader(config Config, borderWidth int, border string) string {
	header := ""
	headerColor := GetColorCode(config.Header.TextColor)
	headerPaddingLength := max(0, borderWidth-len(config.Header.Text)-padding)
	headerPadding := strings.Repeat(" ", headerPaddingLength)
	header += fmt.Sprintf("  │ %s%s%s%s │\n", headerColor, config.Header.Text, headerPadding, Reset)
	if config.Header.Line {
		lineColor := GetColorCode(config.Header.LineColor)
		header += fmt.Sprintf("  ├%s%s%s┤\n", lineColor, border, Reset)
	}
	return header
}

func buildFooter(config Config, borderWidth int, border string) string {
	footer := ""
	if config.Footer.Line {
		lineColor := GetColorCode(config.Footer.LineColor)
		footer += fmt.Sprintf("  ├%s%s%s┤\n", lineColor, border, Reset)
	}
	footerColor := GetColorCode(config.Footer.TextColor)
	footerPaddingLength := max(0, borderWidth-len(config.Footer.Text)-padding)
	footerPadding := strings.Repeat(" ", footerPaddingLength)
	footer += fmt.Sprintf("  │ %s%s%s%s │\n", footerColor, config.Footer.Text, footerPadding, Reset)
	return footer
}

func buildMenuItems(config Config, items map[string]string, borderWidth int, minIconLength int) string {
	menuItems := ""
	for _, item := range config.Items {
		value, exists := items[item.Keyword]
		if exists {
			if item.Value != "" {
				value = item.Value
			} else if items[item.Keyword] != "" {
				value = items[item.Keyword]
			}
		}

		fixedLength := minIconLength + len(item.Text) + padding
		paddingLength := borderWidth - fixedLength
		if paddingLength < 0 {
			paddingLength = 0
		}

		padding := strings.Repeat(" ", paddingLength)
		itemTextColor := GetColorCode(item.TextColor)
		itemIconColor := GetColorCode(item.IconColor)
		itemValueColor := GetColorCode(item.ValueColor)
		menuItems += fmt.Sprintf("  │ %s%s%s  %s%s%s%s │ %s%s%s\n", itemIconColor, item.Icon, Reset, itemTextColor, item.Text, Reset, padding, itemValueColor, value, Reset)
	}
	return menuItems
}

func formatMenuItems(config Config, items map[string]string, borderWidth int, minIconLength int) []string {
	formattedItems := []string{}
	for _, item := range config.Items {
		value, exists := items[item.Keyword]
		if exists {
			if item.Value != "" {
				value = item.Value
			} else if items[item.Keyword] != "" {
				value = items[item.Keyword]
			}
		}

		fixedLength := minIconLength + len(item.Text) + padding
		paddingLength := borderWidth - fixedLength

		padding := strings.Repeat(" ", paddingLength)
		itemTextColor := GetColorCode(item.TextColor)
		itemIconColor := GetColorCode(item.IconColor)
		itemValueColor := GetColorCode(item.ValueColor)
		formattedItem := fmt.Sprintf("   %s%s%s %s%s%s%s  %s%s%s", itemIconColor, item.Icon, Reset, itemTextColor, item.Text, Reset, padding, itemValueColor, value, Reset)
		formattedItems = append(formattedItems, formattedItem)
	}
	return formattedItems
}

func buildColumns(formattedItems []string) string {
	menu := ""
	half := (len(formattedItems) + 1) / 2
	leftColumn := formattedItems[:half]
	rightColumn := formattedItems[half:]

	maxLeftLength := 0
	for _, item := range leftColumn {
		if len(item) > maxLeftLength {
			maxLeftLength = len(item)
		}
	}

	for i := 0; i < len(leftColumn); i++ {
		if i < len(rightColumn) {
			menu += fmt.Sprintf("%-*s %s\n", maxLeftLength, leftColumn[i], rightColumn[i])
		} else {
			menu += fmt.Sprintf("%s\n", leftColumn[i])
		}
	}
	return menu
}

func combineAsciiAndMenu(menu string, paddedAsciiArt string, asciiColors string, position string, maxAsciiWidth int) string {
	menuLines := strings.Split(menu, "\n")
	asciiLines := strings.Split(paddedAsciiArt, "\n")

	maxLines := max(len(menuLines), len(asciiLines))

	// Calculate the length of the longest ASCII art line
	maxAsciiLineWidth := getMaxLineWidth(asciiLines)
	println("ascii", maxAsciiLineWidth)

	// Calculate the length of the longest menu line
	maxMenuWidth := getMaxLineWidth(menuLines)
	println("menu", maxMenuWidth)

	for len(menuLines) < maxLines {
		menuLines = append(menuLines, "")
	}
	for len(asciiLines) < maxLines {
		asciiLines = append(asciiLines, "")
	}

	paddingBetween := 2 // Number of spaces between the menu and the ASCII art
	combinedLines := []string{}

	for i := 0; i < maxLines; i++ {
		if position == "left" {
			asciiLine := asciiLines[i]
			menuLine := menuLines[i]
			padding := strings.Repeat(" ", maxAsciiLineWidth-len(asciiLine)+paddingBetween)
			combinedLines = append(combinedLines, fmt.Sprintf("%s%s%s%s%s", asciiColors, asciiLine, Reset, padding, menuLine))
		} else if position == "right" {
			menuLine := menuLines[i]
			padding := strings.Repeat(" ", maxMenuWidth-len(menuLine)+paddingBetween)
			combinedLines = append(combinedLines, fmt.Sprintf("%s%s%s%s%s", menuLine, padding, asciiColors, asciiLines[i], Reset))
		} else {
			combinedLines = append(combinedLines, menuLines[i])
		}
	}

	if position == "top" {
		return strings.Join(combinedLines, "\n")
	} else if position == "bottom" {
		// When the ASCII art is below the menu, simply concatenate the menu and ASCII art without padding
		return strings.Join(menuLines, "\n") + "\n" + asciiColors + strings.Join(asciiLines, "\n") + Reset
	}

	return strings.Join(combinedLines, "\n")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
