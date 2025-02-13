package pkg

import (
	"fmt"
	"strings"
)

// making padding fixed cuz it keeps making the menu look weird and too lazy
// to figure out a clean way to fix it
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
	maxAsciiWidth := 0
	for _, line := range asciiLines {
		if len(line) > maxAsciiWidth {
			maxAsciiWidth = len(line)
		}
	}

	if config.Ascii.Position == "top" {
		menu += fmt.Sprintf("%s%s%s\n", asciiColors, paddedAsciiArt, Reset)
		menu += fmt.Sprintf("  ╭%s╮\n", border)
	} else {
		menu += fmt.Sprintf("  ╭%s╮\n", border)
	}

	// Icon length can be weird so that's why this exists
	minIconLength := 1000
	for _, item := range config.Items {
		imageLength := len(item.Icon)
		if imageLength < minIconLength {
			minIconLength = imageLength
		}
	}

	if config.Header.Enabled {
		headerColor := GetColorCode(config.Header.TextColor)
		headerPaddingLength := max(0, borderWidth-len(config.Header.Text)-padding)
		headerPadding := strings.Repeat(" ", headerPaddingLength)
		menu += fmt.Sprintf("  │ %s%s%s%s │\n", headerColor, config.Header.Text, headerPadding, Reset)
		if config.Header.Line {
			lineColor := GetColorCode(config.Header.LineColor)
			menu += fmt.Sprintf("  ├%s%s%s┤\n", lineColor, border, Reset)
		}
	}

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
		menu += fmt.Sprintf("  │ %s%s%s  %s%s%s%s │ %s%s%s\n", itemIconColor, item.Icon, Reset, itemTextColor, item.Text, Reset, padding, itemValueColor, value, Reset)
	}

	if config.Footer.Enabled {
		if config.Footer.Line {
			lineColor := GetColorCode(config.Footer.LineColor)
			menu += fmt.Sprintf("  ├%s%s%s┤\n", lineColor, border, Reset)
		}
		footerColor := GetColorCode(config.Footer.TextColor)
		footerPaddingLength := max(0, borderWidth-len(config.Footer.Text)-padding)
		footerPadding := strings.Repeat(" ", footerPaddingLength)
		menu += fmt.Sprintf("  │ %s%s%s%s │\n", footerColor, config.Footer.Text, footerPadding, Reset)
	}

	if config.Ascii.Position == "bottom" {
		menu += fmt.Sprintf("  ╰%s╯\n", border)
		menu += fmt.Sprintf("%s%s%s\n", asciiColors, paddedAsciiArt, Reset)
	} else {
		menu += fmt.Sprintf("  ╰%s╯\n", border)
	}

	if config.Ascii.Position == "left" || config.Ascii.Position == "right" {
		menuLines := strings.Split(menu, "\n")
		asciiLines := strings.Split(paddedAsciiArt, "\n")

		// Calculate the length of the longest menu line
		maxMenuWidth := 0
		for _, line := range menuLines {
			if len(line) > maxMenuWidth {
				maxMenuWidth = len(line)
			}
		}

		maxLines := max(len(menuLines), len(asciiLines))
		for len(menuLines) < maxLines {
			menuLines = append(menuLines, "")
		}
		for len(asciiLines) < maxLines {
			asciiLines = append(asciiLines, "")
		}

		combinedLines := []string{}
		for i := 0; i < maxLines; i++ {
			if config.Ascii.Position == "left" {
				combinedLines = append(combinedLines, fmt.Sprintf("%s%s%s %s%s", asciiColors, asciiLines[i], Reset, strings.Repeat(" ", maxAsciiWidth-len(asciiLines[i])), menuLines[i]))
			} else if config.Ascii.Position == "right" {
				menuLine := menuLines[i]
				padding := strings.Repeat(" ", maxMenuWidth-len(menuLine))
				combinedLines = append(combinedLines, fmt.Sprintf("%s%s %s%s%s", menuLine, padding, asciiColors, asciiLines[i], Reset))
			} else {
				combinedLines = append(combinedLines, menuLines[i])
			}
		}

		menu = strings.Join(combinedLines, "\n")
	}

	return menu
}

func BuildListMenu(items map[string]string, asciiArt string, config Config, borderWidth int) string {
	asciiColors := GetColorCode(config.Ascii.Colors)

	paddedAsciiArt := AddPaddingToMultilineString(asciiArt, config.Ascii.HorizontalPadding, config.Ascii.VerticalPadding)
	border := strings.Repeat("─", borderWidth)

	asciiLines := strings.Split(paddedAsciiArt, "\n")
	maxAsciiWidth := 0
	for _, line := range asciiLines {
		if len(line) > maxAsciiWidth {
			maxAsciiWidth = len(line)
		}
	}

	menu := ""

	if config.Ascii.Position == "top" {
		menu += fmt.Sprintf("%s%s%s\n", asciiColors, paddedAsciiArt, Reset)
	}

	if config.Header.Enabled {
		headerColor := GetColorCode(config.Header.TextColor)
		headerPaddingLength := max(0, borderWidth-len(config.Header.Text)-padding)
		headerPadding := strings.Repeat(" ", headerPaddingLength)
		menu += fmt.Sprintf("   %s%s%s%s \n", headerColor, config.Header.Text, headerPadding, Reset)
		if config.Header.Line {
			lineColor := GetColorCode(config.Header.LineColor)
			menu += fmt.Sprintf("  %s%s%s\n", lineColor, border, Reset)
		}
	}

	// Icon length can be weird so that's why this exists
	minIconLength := 1000
	for _, item := range config.Items {
		imageLength := len(item.Icon)
		if imageLength < minIconLength {
			minIconLength = imageLength
		}
	}

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

	if config.General.Columns {
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
	} else {
		for _, item := range formattedItems {
			menu += fmt.Sprintf("%s\n", item)
		}
	}

	if config.Footer.Enabled {
		if config.Footer.Line {
			lineColor := GetColorCode(config.Footer.LineColor)
			menu += fmt.Sprintf("  %s%s%s\n", lineColor, border, Reset)
		}
		footerColor := GetColorCode(config.Footer.TextColor)
		footerPaddingLength := max(0, borderWidth-len(config.Footer.Text)-padding)
		footerPadding := strings.Repeat(" ", footerPaddingLength)
		menu += fmt.Sprintf("   %s%s%s%s \n", footerColor, config.Footer.Text, footerPadding, Reset)
	}

	if config.Ascii.Position == "bottom" {
		menu += fmt.Sprintf("%s%s%s\n", asciiColors, paddedAsciiArt, Reset)
	}

	menuLines := strings.Split(menu, "\n")
	maxLines := max(len(menuLines), len(asciiLines))

	// Calculate the length of the longest menu line
	maxMenuWidth := 0
	for _, line := range menuLines {
		if len(line) > maxMenuWidth {
			maxMenuWidth = len(line)
		}
	}

	for len(menuLines) < maxLines {
		menuLines = append(menuLines, "")
	}
	for len(asciiLines) < maxLines {
		asciiLines = append(asciiLines, "")
	}

	combinedLines := []string{}
	for i := 0; i < maxLines; i++ {
		if config.Ascii.Position == "left" {
			combinedLines = append(combinedLines, fmt.Sprintf("%s%s%s %s%s", asciiColors, asciiLines[i], Reset, strings.Repeat(" ", maxAsciiWidth-len(asciiLines[i])), menuLines[i]))
		} else if config.Ascii.Position == "right" {
			menuLine := menuLines[i]
			padding := strings.Repeat(" ", maxMenuWidth-len(menuLine))
			println(len(padding))
			combinedLines = append(combinedLines, fmt.Sprintf("%s%s %s%s%s", menuLine, padding, asciiColors, asciiLines[i], Reset))
		} else {
			combinedLines = append(combinedLines, menuLines[i])
		}
	}

	return strings.Join(combinedLines, "\n")
}

func BuildWTFMenu(items map[string]string, asciiArt string, config Config, borderWidth int) string {
	println("WTF")

	return ""
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
