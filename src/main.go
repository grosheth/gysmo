package main

import (
	"flag"
	"fmt"
	"gysmo/src/pkg"
	"path/filepath"
)

const version = "v0.2.0"

func main() {
	filename := flag.String("f", "config.json", "name of the config file in ~/.config/gysmo/")
	useDataFile := flag.Bool("c", false, "use data file for all values")
	showVersion := flag.Bool("v", false, "Show version of gysmo")

	flag.Parse()

	if *showVersion {
		fmt.Printf("%s\n", version)
		return
	}

	workingPath := pkg.LoadWorkingPath()
	configPath := filepath.Join(workingPath, "config", *filename)
	schemaPath := filepath.Join(workingPath, "config", "schema", "config_schema.json")
	asciiPath := filepath.Join(workingPath, "ascii")

	if err := pkg.EnsureConfigFilesExist(); err != nil {
		fmt.Println("Error ensuring config files exist:", err)
		return
	}

	err := pkg.ValidateJsonConfig(configPath, schemaPath)
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
		asciiArt, err = pkg.ReadAsciiArt(filepath.Join(asciiPath, config.Ascii.Path))
		if err != nil {
			fmt.Println("Error reading ASCII art:", err)
			return
		}
	}

	var menu string
	switch config.General.MenuType {
	case "box":
		menu = pkg.BuildBoxMenu(pkg.MenuItems(config, *useDataFile), asciiArt, config)
	case "list":
		menu = pkg.BuildListMenu(pkg.MenuItems(config, *useDataFile), asciiArt, config)
	default:
		menu = pkg.BuildBoxMenu(pkg.MenuItems(config, *useDataFile), asciiArt, config)
	}
	fmt.Println(menu)
}
