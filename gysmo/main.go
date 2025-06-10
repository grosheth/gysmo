package main

import (
	"flag"
	"fmt"
	"gysmo/gysmo/src"
	"path/filepath"
)

const version = "v0.2.2"

func main() {
	filename := flag.String("f", "config.json", "name of the config file in ~/.config/gysmo/")
	useDataFile := flag.Bool("c", false, "use data file for all values")
	showVersion := flag.Bool("v", false, "Show version of gysmo")

	flag.Parse()

	if *showVersion {
		fmt.Printf("%s\n", version)
		return
	}

	workingPath := src.LoadWorkingPath()
	configPath := filepath.Join(workingPath, "config", *filename)
	schemaPath := filepath.Join(workingPath, "config", "schema", "config_schema.json")
	asciiPath := filepath.Join(workingPath, "ascii")

	if err := src.EnsureConfigFilesExist(); err != nil {
		fmt.Println("Error ensuring config files exist:", err)
		return
	}

	err := src.ValidateJsonConfig(configPath, schemaPath)
	if err != nil {
		fmt.Println("Error validating config.json:", err)
		return
	}

	// Ensure config files exist
	err = src.EnsureConfigFilesExist()
	if err != nil {
		fmt.Println("No Config file found.", err)
		return
	}

	config, err := src.LoadConfig(configPath)
	if err != nil {
		fmt.Println("Error loading config.json:", err)
		return
	}

	var asciiArt string
	if config.Ascii.Enabled {
		asciiArt, err = src.ReadAsciiArt(filepath.Join(asciiPath, config.Ascii.Path))
		if err != nil {
			fmt.Println("Error reading ASCII art:", err)
			return
		}
	}

	var menu string
	switch config.General.MenuType {
	case "box":
		menu = src.BuildBoxMenu(src.MenuItems(config, *useDataFile), asciiArt, config)
	case "list":
		menu = src.BuildListMenu(src.MenuItems(config, *useDataFile), asciiArt, config)
	default:
		menu = src.BuildBoxMenu(src.MenuItems(config, *useDataFile), asciiArt, config)
	}
	fmt.Println(menu)
}
