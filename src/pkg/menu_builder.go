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

	if config.Ascii.Position == "top" {
		menu += fmt.Sprintf("%s%s%s\n", asciiColors, paddedAsciiArt, Reset)
		menu += fmt.Sprintf(" ╭%s╮\n", border)
	} else {
		menu += fmt.Sprintf(" ╭%s╮\n", border)
	}

	minIconLength := GetMinIconLength(config.Items)

	if config.Header.Enabled {
		menu += buildHeader(config, borderWidth, border)
	}

	menu += buildMenuItems(config, items, borderWidth, minIconLength)

	if config.Footer.Enabled {
		menu += buildFooter(config, borderWidth, border)
	}

	if config.Ascii.Position == "bottom" {
		menu += fmt.Sprintf(" ╰%s╯\n", border)
		menu += fmt.Sprintf("%s%s%s\n", asciiColors, paddedAsciiArt, Reset)
	} else {
		menu += fmt.Sprintf(" ╰%s╯\n", border)
	}

	if config.Ascii.Position == "left" {
		menu = combineAsciiAndMenuLeft(menu, paddedAsciiArt, asciiColors)
	} else if config.Ascii.Position == "right" {
		menu = combineAsciiAndMenuRightBox(menu, paddedAsciiArt, asciiColors)
	}

	return menu
}

func BuildListMenu(items map[string]string, asciiArt string, config Config, borderWidth int) string {
	asciiColors := GetColorCode(config.Ascii.Colors)

	paddedAsciiArt := AddPaddingToMultilineString(asciiArt, config.Ascii.HorizontalPadding, config.Ascii.VerticalPadding)
	menu := ""

	headerWidth := 0
	footerWidth := 0

	if config.Ascii.Position == "top" {
		menu += fmt.Sprintf("%s%s%s\n", asciiColors, paddedAsciiArt, Reset)
	}

	if config.Header.Enabled {
		headerWidth = len(config.Header.Text)
		menu += fmt.Sprintf("%s%s%s\n", GetColorCode(config.Header.TextColor), config.Header.Text, Reset)
		if config.Header.Line {
			menu += fmt.Sprintf("%s%s%s\n", GetColorCode(config.Header.LineColor), strings.Repeat("─", borderWidth), Reset)
		}
	}

	minIconLength := GetMinIconLength(config.Items)

	formattedItems := formatMenuItems(config, items, borderWidth, minIconLength)

	maxItemWidth := 0
	for _, item := range formattedItems {
		itemWidth := len(StripAnsiCodes(item))
		if itemWidth > maxItemWidth {
			maxItemWidth = itemWidth
		}
	}

	if config.General.Columns {
		menu += buildColumns(formattedItems)
	} else {
		for _, item := range formattedItems {
			menu += fmt.Sprintf("%s\n", item)
		}
	}

	if config.Footer.Enabled {
		footerWidth = len(config.Footer.Text)
		if config.Footer.Line {
			menu += fmt.Sprintf("%s%s%s\n", GetColorCode(config.Footer.LineColor), strings.Repeat("─", borderWidth), Reset)
		}
		menu += fmt.Sprintf("%s%s%s\n", GetColorCode(config.Footer.TextColor), config.Footer.Text, Reset)
	}

	if config.Ascii.Position == "left" {
		menu = combineAsciiAndMenuLeft(menu, paddedAsciiArt, asciiColors)
	} else if config.Ascii.Position == "right" {
		if config.General.Columns {
			menu = combineAsciiAndMenuRightColumns(menu, paddedAsciiArt, asciiColors, config, borderWidth, headerWidth, footerWidth)
		} else {
			menu = combineAsciiAndMenuRight(menu, paddedAsciiArt, asciiColors, config, borderWidth, headerWidth, footerWidth)
		}
	}

	if config.Ascii.Position == "bottom" {
		menu += fmt.Sprintf("%s%s%s\n", asciiColors, paddedAsciiArt, Reset)
	}
	return menu
}

func GetMaxLineWidth(lines []string) int {
	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}
	return maxWidth
}

func GetMinIconLength(items []ConfigItem) int {
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
	header += fmt.Sprintf(" │ %s%s%s%s │\n", headerColor, config.Header.Text, headerPadding, Reset)
	if config.Header.Line {
		lineColor := GetColorCode(config.Header.LineColor)
		header += fmt.Sprintf(" ├%s%s%s┤\n", lineColor, border, Reset)
	}
	return header
}

func buildFooter(config Config, borderWidth int, border string) string {
	footer := ""
	if config.Footer.Line {
		lineColor := GetColorCode(config.Footer.LineColor)
		footer += fmt.Sprintf(" ├%s%s%s┤\n", lineColor, border, Reset)
	}
	footerColor := GetColorCode(config.Footer.TextColor)
	footerPaddingLength := max(0, borderWidth-len(config.Footer.Text)-padding)
	footerPadding := strings.Repeat(" ", footerPaddingLength)
	footer += fmt.Sprintf(" │ %s%s%s%s │\n", footerColor, config.Footer.Text, footerPadding, Reset)
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
		menuItems += fmt.Sprintf(" │ %s%s%s  %s%s%s%s │ %s%s%s\n", itemIconColor, item.Icon, Reset, itemTextColor, item.Text, Reset, padding, itemValueColor, value, Reset)
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
		formattedItem := fmt.Sprintf(" %s%s%s %s%s%s%s  %s%s%s", itemIconColor, item.Icon, Reset, itemTextColor, item.Text, Reset, padding, itemValueColor, value, Reset)
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

	for i := range len(leftColumn) {
		if i < len(rightColumn) {
			menu += fmt.Sprintf("%-*s | %s\n", maxLeftLength, leftColumn[i], rightColumn[i])
		} else {
			menu += fmt.Sprintf("%s\n", leftColumn[i])
		}
	}
	return menu
}

func combineAsciiAndMenu(menu string, paddedAsciiArt string, asciiColors string, position string, config Config, borderWidth int, headerWidth int, footerWidth int) string {
	if position == "left" {
		return combineAsciiAndMenuLeft(menu, paddedAsciiArt, asciiColors)
	} else if position == "right" {
		return combineAsciiAndMenuRight(menu, paddedAsciiArt, asciiColors, config, borderWidth, headerWidth, footerWidth)
	}

	menuLines := strings.Split(menu, "\n")
	asciiLines := strings.Split(paddedAsciiArt, "\n")
	maxLines := max(len(menuLines), len(asciiLines))

	// Ensure both menuLines and asciiLines have the same number of lines
	for len(menuLines) < maxLines {
		menuLines = append(menuLines, "")
	}
	for len(asciiLines) < maxLines {
		asciiLines = append(asciiLines, "")
	}

	combinedLines := []string{}
	for i := range maxLines {
		combinedLines = append(combinedLines, menuLines[i])
	}

	if position == "top" {
		return strings.Join(combinedLines, "\n")
	} else if position == "bottom" {
		return strings.Join(menuLines, "\n") + "\n" + asciiColors + strings.Join(asciiLines, "\n") + Reset
	}

	return strings.Join(combinedLines, "\n")
}

func combineAsciiAndMenuLeft(menu string, paddedAsciiArt string, asciiColors string) string {
	menuLines := strings.Split(menu, "\n")
	asciiLines := strings.Split(paddedAsciiArt, "\n")
	maxLines := max(len(menuLines), len(asciiLines))

	maxAsciiLineWidth := GetMaxLineWidth(asciiLines)

	for len(menuLines) < maxLines {
		menuLines = append(menuLines, "")
	}
	for len(asciiLines) < maxLines {
		asciiLines = append(asciiLines, "")
	}

	combinedLines := []string{}

	for i := range maxLines {
		asciiLine := asciiLines[i]
		menuLine := menuLines[i]
		padding := max(2, maxAsciiLineWidth-len(StripAnsiCodes(asciiLine))+2)
		combinedLines = append(combinedLines, fmt.Sprintf("%s%s%s%s%s", asciiColors, asciiLine, Reset, strings.Repeat(" ", padding), menuLine))
	}

	return strings.Join(combinedLines, "\n")
}

func combineAsciiAndMenuRight(menu string, paddedAsciiArt string, asciiColors string, config Config, borderWidth int, headerWidth int, footerWidth int) string {
	menuLines := strings.Split(menu, "\n")
	asciiLines := strings.Split(paddedAsciiArt, "\n")
	maxLines := max(len(menuLines), len(asciiLines))

	// Calculate the maximum width of the menu lines
	longestMenuLineWidth := 0
	for _, line := range menuLines {
		currentLine := StripAnsiCodes(line)
		if longestMenuLineWidth < len(currentLine) {
			longestMenuLineWidth = len(currentLine)
		}
	}

	// Ensure both menuLines and asciiLines have the same number of lines
	for len(menuLines) < maxLines {
		menuLines = append(menuLines, "")
	}
	for len(asciiLines) < maxLines {
		asciiLines = append(asciiLines, "")
	}

	combinedLines := []string{}

	for i := range maxLines {
		menuLine := menuLines[i]
		asciiLine := asciiLines[i]
		padding := 0
		line := IsLine(menuLine)
		headerOrFooter := IsHeaderOrFooter(menuLine, config, i, maxLines)
		if len(menuLine) == 0 {
			padding = longestMenuLineWidth
		} else if line {
			padding = longestMenuLineWidth - borderWidth
		} else if headerOrFooter == "Header" {
			padding = longestMenuLineWidth - headerWidth
		} else if headerOrFooter == "Footer" {
			padding = longestMenuLineWidth - footerWidth
		} else {
			padding = longestMenuLineWidth - len(StripAnsiCodes(menuLine)) + 2
		}
		combinedLines = append(combinedLines, fmt.Sprintf("%s%s%s%s%s", menuLine, strings.Repeat(" ", padding), asciiColors, asciiLine, Reset))
	}
	return strings.Join(combinedLines, "\n")
}

func combineAsciiAndMenuRightColumns(menu string, paddedAsciiArt string, asciiColors string, config Config, borderWidth int, headerWidth int, footerWidth int) string {
	menuLines := strings.Split(menu, "\n")
	asciiLines := strings.Split(paddedAsciiArt, "\n")
	maxLines := max(len(menuLines), len(asciiLines))

	// Calculate the maximum width of the menu lines, including header and footer
	longestMenuLineWidth := 0
	for _, line := range menuLines {
		currentLine := StripAnsiCodes(line)
		if longestMenuLineWidth < len(currentLine) {
			longestMenuLineWidth = len(currentLine)
		}
	}

	// Ensure both menuLines and asciiLines have the same number of lines
	for len(menuLines) < maxLines {
		menuLines = append(menuLines, "")
	}
	for len(asciiLines) < maxLines {
		asciiLines = append(asciiLines, "")
	}

	combinedLines := []string{}

	for i := range maxLines {
		menuLine := menuLines[i]
		asciiLine := asciiLines[i]

		// Check if the current line has two columns by looking for the | delimiter
		hasTwoColumns := strings.Contains(menuLine, " | ")

		padding := 0
		line := IsLine(menuLine)
		headerOrFooter := IsHeaderOrFooter(menuLine, config, i, maxLines)
		if len(menuLine) == 0 {
			padding = longestMenuLineWidth
		} else if hasTwoColumns {
			padding = longestMenuLineWidth - len(StripAnsiCodes(menuLine)) + 4
		} else if line {
			padding = longestMenuLineWidth - borderWidth
		} else if headerOrFooter == "Header" {
			padding = longestMenuLineWidth - headerWidth
		} else if headerOrFooter == "Footer" {
			padding = longestMenuLineWidth - footerWidth
		} else {
			padding = longestMenuLineWidth - len(StripAnsiCodes(menuLine)) + 2
		}

		// Ensure padding is not negative
		if padding < 0 {
			padding = 0
		}

		combinedLines = append(combinedLines, fmt.Sprintf("%s%s%s%s%s", menuLine, strings.Repeat(" ", padding), asciiColors, asciiLine, Reset))
	}

	return strings.Join(combinedLines, "\n")
}

func combineAsciiAndMenuRightBox(menu string, paddedAsciiArt string, asciiColors string) string {
	menuLines := strings.Split(menu, "\n")
	asciiLines := strings.Split(paddedAsciiArt, "\n")
	maxLines := max(len(menuLines), len(asciiLines))

	menuLength := 0
	longestValueLength := 0
	for _, line := range menuLines {
		if strings.Contains(line, " │ ") {
			parts := strings.SplitN(line, " │ ", 4)
			currentLine := StripAnsiCodes(parts[1])
			if menuLength < len(currentLine) {
				menuLength = len(currentLine)
			}
			if len(parts) > 2 {
				valueLength := len(StripAnsiCodes(parts[2]))
				if valueLength > longestValueLength {
					longestValueLength = valueLength
				}
			}
		}
	}

	for len(menuLines) < maxLines {
		menuLines = append(menuLines, "")
	}
	for len(asciiLines) < maxLines {
		asciiLines = append(asciiLines, "")
	}

	combinedLines := []string{}

	for i := range maxLines {
		menuLine := menuLines[i]
		asciiLine := asciiLines[i]
		padding := 0

		hasValue := strings.Contains(menuLine, " │ ")

		part := strings.SplitN(menuLine, " │ ", 4)
		if len(menuLine) == 0 {
			padding = menuLength + longestValueLength + 4
		} else if hasValue && len(part) > 2 {
			valueLine := len(StripAnsiCodes(part[2]))
			padding = longestValueLength - valueLine + 2
		} else {
			padding = longestValueLength + 3
		}
		combinedLines = append(combinedLines, fmt.Sprintf("%s%s%s%s%s", menuLine, strings.Repeat(" ", padding), asciiColors, asciiLine, Reset))
	}

	return strings.Join(combinedLines, "\n")
}
