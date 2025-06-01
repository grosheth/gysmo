#!/bin/sh

# Define color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RESET='\033[0m'

REPO="grosheth/gysmo"

CONFIGURATION_PATH="$HOME/.config/gysmo/config"
ASCII_PATH="$HOME/.config/gysmo/ascii"
DATA_PATH="$HOME/.config/gysmo/data"
SCHEMA_PATH="$CONFIGURATION_PATH/schema"
CONFIG_FILE="$CONFIGURATION_PATH/config.json"
SCHEMA_FILE="$SCHEMA_PATH/config_schema.json"
ASCII_FILE="$ASCII_PATH/gysmo"

command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Get latest release tag
get_latest_release() {
    if command_exists curl; then
        curl --silent "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    elif command_exists wget; then
        wget -qO- "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    else
        echo -e "${RED}Error: wget or curl is not installed. Please install wget or curl and try again.${RESET}"
        exit 1
    fi
}

# Ensure wget or curl is installed
if ! command_exists wget && ! command_exists curl; then
    echo -e "${RED}Error: wget or curl is not installed. Please install wget or curl and try again.${RESET}"
    exit 1
fi

# Fetch the latest release tag
echo "Fetching the latest release..."
LATEST_TAG=$(get_latest_release)
if [ -z "$LATEST_TAG" ]; then
    echo -e "${RED}Error: Failed to fetch the latest release.${RESET}"
    exit 1
fi

RELEASE_URL="https://github.com/$REPO/releases/download/$LATEST_TAG/gysmo"
CONFIG_URL="https://raw.githubusercontent.com/$REPO/$LATEST_TAG/gysmo/config/config.json"
SCHEMA_URL="https://raw.githubusercontent.com/$REPO/$LATEST_TAG/gysmo/config/schema/config_schema.json"
ASCII_URL="https://raw.githubusercontent.com/$REPO/$LATEST_TAG/gysmo/ascii/gysmo"

# Create the necessary directory structure in ~/.config/gysmo
echo "Setting up default templates in $HOME/.config/gysmo..."
mkdir -p "$CONFIGURATION_PATH"
mkdir -p "$ASCII_PATH"
mkdir -p "$DATA_PATH"
mkdir -p "$SCHEMA_PATH"

# Download default templates to ~/.config/gysmo
if [ ! -f "$CONFIG_FILE" ]; then
    echo "Downloading default configuration file..."
    if command_exists wget; then
        wget -O "$CONFIG_FILE" "$CONFIG_URL"
    elif command_exists curl; then
        curl -o "$CONFIG_FILE" "$CONFIG_URL"
    fi
fi

if [ ! -f "$SCHEMA_FILE" ]; then
    echo "Downloading default schema file..."
    if command_exists wget; then
        wget -O "$SCHEMA_FILE" "$SCHEMA_URL"
    elif command_exists curl; then
        curl -o "$SCHEMA_FILE" "$SCHEMA_URL"
    fi
fi

if [ ! -f "$ASCII_FILE" ]; then
    echo "Downloading default ASCII file..."
    if command_exists wget; then
        wget -O "$ASCII_FILE" "$ASCII_URL"
    elif command_exists curl; then
        curl -o "$ASCII_FILE" "$ASCII_URL"
    fi
fi

echo -e "${GREEN}Installation completed successfully.${RESET}"
