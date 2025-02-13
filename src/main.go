package main

import (
	"flag"
	"fmt"
	"gysmo/src/pkg"
	"path/filepath"
)

func main() {

	filename := flag.String("f", "config.json", "name of the config file in /.config/gysmo/")
	useDataFile := flag.Bool("c", false, "use data file for all values")

	flag.Parse()
	configPath := filepath.Join("config", *filename)
	err := pkg.ValidateJsonConfig(configPath, "config/config_schema.json")
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
		asciiArt, err = pkg.ReadAsciiArt(config.Ascii.Path)
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
	case "wtf":
		menu = pkg.BuildWTFMenu(pkg.MenuItems(config, *useDataFile), asciiArt, config, borderWidth)
	default:
		menu = pkg.BuildBoxMenu(pkg.MenuItems(config, *useDataFile), asciiArt, config, borderWidth)
	}
	fmt.Println(menu)
}
