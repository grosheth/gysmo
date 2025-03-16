package pkg

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
)

const padding = 2

func GetMaxIconLength(items []ConfigItem) int {
	maxIconLength := 0
	for _, item := range items {
		iconLength := len(StripAnsiCodes(item.Icon))
		if iconLength > maxIconLength {
			maxIconLength = iconLength
		}
	}
	return maxIconLength
}

func AddPaddingToMultilineString(s string, horizontalPadding int, verticalPadding int) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.Repeat(" ", horizontalPadding) + line + strings.Repeat(" ", horizontalPadding)
	}
	paddingLines := strings.Repeat("\n", verticalPadding)
	return paddingLines + strings.Join(lines, "\n") + paddingLines
}

func DefineBoxBorder(config Config) int {
	borderWidth := 0

	maxIconLength := GetMaxIconLength(config.Items)

	for _, item := range config.Items {
		maxitemLength := len(StripAnsiCodes(item.Text)) + maxIconLength + padding
		if maxitemLength > borderWidth {
			borderWidth = maxitemLength
		}
	}

	if config.Header.Enabled {
		headerLength := len(config.Header.Text)
		if headerLength > borderWidth {
			borderWidth = headerLength + padding
		}
	}

	if config.Footer.Enabled {
		footerLength := len(config.Footer.Text)
		if footerLength > borderWidth {
			borderWidth = footerLength + padding
		}
	}

	return borderWidth
}

func BuildBoxMenu(items map[string]string, asciiArt string, config Config) string {

	menuPadding := strings.Repeat(" ", config.General.MenuPadding)
	borderWidth := DefineBoxBorder(config)
	border := strings.Repeat("─", borderWidth)
	menu := ""

	asciiColors := GetColorCode(config.Ascii.Colors)
	paddedAsciiArt := AddPaddingToMultilineString(asciiArt, config.Ascii.HorizontalPadding, config.Ascii.VerticalPadding)

	if config.Ascii.Position == "top" {
		menu += fmt.Sprintf("%s%s%s\n", asciiColors, paddedAsciiArt, Reset)
		menu += fmt.Sprintf("%s╭%s╮\n", menuPadding, border)
	} else {
		menu += fmt.Sprintf("%s╭%s╮\n", menuPadding, border)
	}

	maxIconLength := GetMaxIconLength(config.Items)

	if config.Header.Enabled {
		menu += buildHeader(config, borderWidth, border)
	}

	menu += buildMenuItems(config, items, borderWidth, maxIconLength)

	if config.Footer.Enabled {
		menu += buildFooter(config, borderWidth, border)
	}

	if config.Ascii.Position == "bottom" {
		menu += fmt.Sprintf("%s╰%s╯\n", menuPadding, border)
		menu += fmt.Sprintf("%s%s%s\n", asciiColors, paddedAsciiArt, Reset)
	} else {
		menu += fmt.Sprintf("%s╰%s╯\n", menuPadding, border)
	}

	if config.Ascii.Position == "left" {
		menu = combineAsciiAndMenuLeft(menu, paddedAsciiArt, asciiColors)
	} else if config.Ascii.Position == "right" {
		menu = combineAsciiAndMenuRightBox(menu, paddedAsciiArt, asciiColors)
	}

	return menu
}

func BuildListMenu(items map[string]string, asciiArt string, config Config) string {

	borderWidth := DefineBoxBorder(config)
	asciiColors := GetColorCode(config.Ascii.Colors)
	menuPadding := strings.Repeat(" ", config.General.MenuPadding)

	paddedAsciiArt := AddPaddingToMultilineString(asciiArt, config.Ascii.HorizontalPadding, config.Ascii.VerticalPadding)
	menu := ""

	headerWidth := 0
	footerWidth := 0

	if config.Ascii.Position == "top" {
		menu += fmt.Sprintf("%s%s%s\n", asciiColors, paddedAsciiArt, Reset)
	}

	if config.Header.Enabled {
		headerWidth = len(config.Header.Text)
		menu += fmt.Sprintf("%s%s%s%s\n", menuPadding, GetColorCode(config.Header.TextColor), config.Header.Text, Reset)
		if config.Header.Line {
			menu += fmt.Sprintf("%s%s%s%s\n", menuPadding, GetColorCode(config.Header.LineColor), strings.Repeat("─", borderWidth), Reset)
		}
	}

	maxIconLength := GetMaxIconLength(config.Items)
	formattedItems := formatMenuItems(config, items, borderWidth, maxIconLength)

	maxItemWidth := 0
	for _, item := range formattedItems {
		itemWidth := len(StripAnsiCodes(item))
		if itemWidth > maxItemWidth {
			maxItemWidth = itemWidth
		}
	}

	if config.General.Columns {
		menu += buildColumns(formattedItems, config)
	} else {
		for _, item := range formattedItems {
			menu += fmt.Sprintf("%s%s\n", menuPadding, item)
		}
	}

	if config.Footer.Enabled {
		footerWidth = len(config.Footer.Text)
		if config.Footer.Line {
			menu += fmt.Sprintf("%s%s%s%s\n", menuPadding, GetColorCode(config.Footer.LineColor), strings.Repeat("─", borderWidth), Reset)
		}
		menu += fmt.Sprintf("%s%s%s%s\n", menuPadding, GetColorCode(config.Footer.TextColor), config.Footer.Text, Reset)
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

func buildHeader(config Config, borderWidth int, border string) string {
	menuPadding := strings.Repeat(" ", config.General.MenuPadding)
	header := ""
	headerColor := GetColorCode(config.Header.TextColor)
	headerPaddingLength := max(0, borderWidth-len(config.Header.Text)-padding)
	headerPadding := strings.Repeat(" ", headerPaddingLength)
	header += fmt.Sprintf("%s│ %s%s%s%s │\n", menuPadding, headerColor, config.Header.Text, headerPadding, Reset)
	if config.Header.Line {
		lineColor := GetColorCode(config.Header.LineColor)
		header += fmt.Sprintf("%s├%s%s%s┤\n", menuPadding, lineColor, border, Reset)
	}
	return header
}

func buildFooter(config Config, borderWidth int, border string) string {
	menuPadding := strings.Repeat(" ", config.General.MenuPadding)
	footer := ""
	if config.Footer.Line {
		lineColor := GetColorCode(config.Footer.LineColor)
		footer += fmt.Sprintf("%s├%s%s%s┤\n", menuPadding, lineColor, border, Reset)
	}
	footerColor := GetColorCode(config.Footer.TextColor)
	footerPaddingLength := max(0, borderWidth-len(config.Footer.Text)-padding)
	footerPadding := strings.Repeat(" ", footerPaddingLength)
	footer += fmt.Sprintf("%s│ %s%s%s%s │\n", menuPadding, footerColor, config.Footer.Text, footerPadding, Reset)
	return footer
}

func buildMenuItems(config Config, items map[string]string, borderWidth int, IconLength int) string {
	menuItems := ""
	for _, item := range config.Items {
		value, exists := items[item.Keyword]
		if exists {
			value = items[item.Keyword]
		} else {
			value = item.Value
		}

		menuPadding := strings.Repeat(" ", config.General.MenuPadding)
		textLength := len(StripAnsiCodes(item.Text))
		fixedLength := IconLength + textLength + padding
		paddingLength := max(0, borderWidth-fixedLength)
		padding := strings.Repeat(" ", paddingLength)
		itemTextColor := GetColorCode(item.TextColor)
		itemIconColor := GetColorCode(item.IconColor)
		itemValueColor := GetColorCode(item.ValueColor)

		fixedIconSpace := fmt.Sprintf("%-*s", IconLength, item.Icon)
		itemString := fmt.Sprintf("%s%s%s", itemIconColor, fixedIconSpace, Reset)
		itemTextString := fmt.Sprintf("%s%s%s", itemTextColor, item.Text, Reset)
		itemValueString := fmt.Sprintf("%s%s%s", itemValueColor, value, Reset)

		if IconLength > 0 {
			menuItems += fmt.Sprintf("%s│ %s%s%s │ %s\n", menuPadding, itemString, itemTextString, padding, itemValueString)
		} else {
			menuItems += fmt.Sprintf("%s│ %s%s │ %s\n", menuPadding, itemTextString, padding, itemValueString)
		}
	}
	return menuItems
}

func formatMenuItems(config Config, items map[string]string, borderWidth int, IconLength int) []string {
	formattedItems := []string{}
	for _, item := range config.Items {
		value, exists := items[item.Keyword]
		if exists {
			value = items[item.Keyword]
		} else {
			value = item.Value
		}

		fixedLength := IconLength + len(item.Text) + padding
		paddingLength := borderWidth - fixedLength

		padding := strings.Repeat(" ", paddingLength)
		itemTextColor := GetColorCode(item.TextColor)
		itemIconColor := GetColorCode(item.IconColor)
		itemValueColor := GetColorCode(item.ValueColor)

		iconString := fmt.Sprintf("%s%s%s", itemIconColor, item.Icon, Reset)
		valueString := fmt.Sprintf("%s%s%s", itemValueColor, value, Reset)
		textString := fmt.Sprintf("%s%s%s", itemTextColor, item.Text, Reset)

		formattedItem := fmt.Sprintf("%s %s  %s%s", iconString, textString, padding, valueString)
		formattedItems = append(formattedItems, formattedItem)
	}

	return formattedItems
}

func getLengthLeftColumn(leftColumn []string) int {
	maxLeftLength := 0
	for _, item := range leftColumn {
		itemLength := len(StripAnsiCodes(item))
		if itemLength > maxLeftLength {
			maxLeftLength = itemLength
		}
	}
	return maxLeftLength
}

func buildColumns(formattedItems []string, config Config) string {
	menu := ""
	half := (len(formattedItems) + 1) / 2
	leftColumn := formattedItems[:half]
	rightColumn := formattedItems[half:]
	menuPadding := strings.Repeat(" ", config.General.MenuPadding)

	maxDisplayWidth := 0
	for _, item := range leftColumn {
		stripped := StripAnsiCodes(item)
		displayWidth := runewidth.StringWidth(stripped)
		if displayWidth > maxDisplayWidth {
			maxDisplayWidth = displayWidth
		}
	}

	for i := range len(leftColumn) {
		if i < len(rightColumn) {
			leftItem := leftColumn[i]
			rightItem := rightColumn[i]
			strippedLeftItem := StripAnsiCodes(leftItem)
			leftItemDisplayWidth := runewidth.StringWidth(strippedLeftItem)

			paddingLength := maxDisplayWidth - leftItemDisplayWidth
			padding := strings.Repeat(" ", paddingLength)

			menuLine := fmt.Sprintf("%s%s%s |  %s\n", menuPadding, leftItem, padding, rightItem)
			menu += menuLine

		} else {
			menuLine := fmt.Sprintf("%s%s\n", menuPadding, leftColumn[i])
			menu += menuLine
		}
	}
	return menu
}

func menuLinesEqualsAsciiLines(menuLines []string, asciiLines []string, maxLines int) ([]string, []string) {
	for len(menuLines) < maxLines {
		menuLines = append(menuLines, "")
	}
	for len(asciiLines) < maxLines {
		asciiLines = append(asciiLines, "")
	}

	return menuLines, asciiLines
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
	menuLines, asciiLines = menuLinesEqualsAsciiLines(menuLines, asciiLines, maxLines)

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
		asciiString := fmt.Sprintf("%s%s%s", asciiColors, asciiLine, Reset)
		padding := strings.Repeat(" ", max(2, maxAsciiLineWidth-len(StripAnsiCodes(asciiLine))+2))
		combinedLines = append(combinedLines, fmt.Sprintf("%s%s%s", asciiString, padding, menuLine))
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

	menuLines, asciiLines = menuLinesEqualsAsciiLines(menuLines, asciiLines, maxLines)

	combinedLines := []string{}

	menuPadding := config.General.MenuPadding
	for i := range maxLines {
		menuLine := menuLines[i]
		asciiLine := asciiLines[i]
		padding := 0
		line := IsLine(menuLine)
		headerOrFooter := IsHeaderOrFooter(menuLine, config, i, maxLines)
		if len(menuLine) == 0 {
			padding = longestMenuLineWidth
		} else if line {
			padding = longestMenuLineWidth - borderWidth - menuPadding
		} else if headerOrFooter == "Header" {
			padding = longestMenuLineWidth - headerWidth - menuPadding
		} else if headerOrFooter == "Footer" {
			padding = longestMenuLineWidth - footerWidth - menuPadding
		} else {
			padding = longestMenuLineWidth - len(StripAnsiCodes(menuLine)) + 2
		}
		asciiString := fmt.Sprintf("%s%s%s", asciiColors, asciiLine, Reset)
		combinedLines = append(combinedLines, fmt.Sprintf("%s%s%s", menuLine, strings.Repeat(" ", padding), asciiString))
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

	menuLines, asciiLines = menuLinesEqualsAsciiLines(menuLines, asciiLines, maxLines)

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

		asciiString := fmt.Sprintf("%s%s%s", asciiColors, asciiLine, Reset)
		combinedLines = append(combinedLines, fmt.Sprintf("%s%s%s", menuLine, strings.Repeat(" ", padding), asciiString))
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
		asciiString := fmt.Sprintf("%s%s%s", asciiColors, asciiLine, Reset)
		combinedLines = append(combinedLines, fmt.Sprintf("%s%s%s", menuLine, strings.Repeat(" ", padding), asciiString))
	}

	return strings.Join(combinedLines, "\n")
}
