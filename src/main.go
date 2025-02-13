package main

import (
	"flag"
	"fmt"
	"gysmo/src/pkg"
	"os"
	"path/filepath"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}

	filename := flag.String("f", "config.json", "name of the config file in ~/.config/gysmo/")
	useDataFile := flag.Bool("c", false, "use data file for all values")
	workingPath := filepath.Join(homeDir, ".config", "gysmo")

	flag.Parse()
	configPath := filepath.Join(workingPath, "config", *filename)
	schemaPath := filepath.Join(workingPath, "config", "schema", "config_schema.json")

	err = pkg.ValidateJsonConfig(configPath, schemaPath)
	if err != nil {
		fmt.Println("Error validating config.json:", err)
		return
	}

	config, err := pkg.LoadConfig(configPath)
	if err != nil {
		fmt.Println("Error loading config.json:", err)
		return
	}

	var asciiArt string
	if config.Ascii.Enabled {
		asciiArt, err = pkg.ReadAsciiArt(filepath.Join(workingPath, config.Ascii.Path))
		if err != nil {
			fmt.Println("Error reading ASCII art:", err)
			return
		}
	}

	borderWidth := pkg.DefineBoxBorder(config)
	var menu string
	switch config.General.MenuType {
	case "box":
		menu = pkg.BuildBoxMenu(pkg.MenuItems(config, *useDataFile), asciiArt, config, borderWidth)
	case "list":
		menu = pkg.BuildListMenu(pkg.MenuItems(config, *useDataFile), asciiArt, config, borderWidth)
	default:
		menu = pkg.BuildBoxMenu(pkg.MenuItems(config, *useDataFile), asciiArt, config, borderWidth)
	}
	fmt.Println(menu)
}
