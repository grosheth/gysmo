#!/bin/sh

# Define color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RESET='\033[0m'

REPO="grosheth/gysmo"
INSTALL_DIR="$HOME/.local/bin"
INSTALL_PATH="$INSTALL_DIR/gysmo"

SHARE_PATH="/usr/share/gysmo"
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

fetch_releases() {
    if command_exists curl; then
        curl -s "https://api.github.com/repos/$REPO/releases" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    elif command_exists wget; then
        wget -qO- "https://api.github.com/repos/$REPO/releases" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
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

# Fetch the list of releases
echo "Fetching available versions..."
VERSIONS=$(fetch_releases)
if [ -z "$VERSIONS" ]; then
    echo -e "${RED}Error: Failed to fetch available versions.${RESET}"
    exit 1
fi

echo "Available versions:"
echo "$VERSIONS" | nl -w 2 -s '. '
echo -e "${BLUE}"
read -rp "Enter the number of the version you want to install: " version_number
echo -e "${RESET}"

# Get the selected version
VERSION=$(echo "$VERSIONS" | sed -n "${version_number}p")
if [ -z "$VERSION" ]; then
    echo -e "${RED}Error: Invalid version selected.${RESET}"
    exit 1
fi

RELEASE_URL="https://github.com/$REPO/releases/download/$VERSION/gysmo"
CONFIG_URL="https://raw.githubusercontent.com/$REPO/refs/tags/$VERSION/src/config/config.json"
SCHEMA_URL="https://raw.githubusercontent.com/$REPO/refs/tags/$VERSION/src/config/schema/config_schema.json"
ASCII_URL="https://raw.githubusercontent.com/$REPO/refs/tags/$VERSION/src/ascii/gysmo"

# Create the necessary directory structure in /usr/share/gysmo
echo "Setting up default templates in $SHARE_PATH..."
sudo mkdir -p "$SHARE_PATH/config"
sudo mkdir -p "$SHARE_PATH/ascii"
sudo mkdir -p "$SHARE_PATH/config/schema"

# Download default templates to /usr/share/gysmo
if [ ! -f "$SHARE_PATH/config/config.json" ]; then
    echo "Downloading default configuration file..."
    if command_exists wget; then
        sudo wget -O "$SHARE_PATH/config/config.json" "$CONFIG_URL"
    elif command_exists curl; then
        sudo curl -o "$SHARE_PATH/config/config.json" "$CONFIG_URL"
    fi
fi

if [ ! -f "$SHARE_PATH/config/schema/config_schema.json" ]; then
    echo "Downloading default schema file..."
    if command_exists wget; then
        sudo wget -O "$SHARE_PATH/config/schema/config_schema.json" "$SCHEMA_URL"
    elif command_exists curl; then
        sudo curl -o "$SHARE_PATH/config/schema/config_schema.json" "$SCHEMA_URL"
    fi
fi

if [ ! -f "$SHARE_PATH/ascii/gysmo" ]; then
    echo "Downloading default ASCII file..."
    if command_exists wget; then
        sudo wget -O "$SHARE_PATH/ascii/gysmo" "$ASCII_URL"
    elif command_exists curl; then
        sudo curl -o "$SHARE_PATH/ascii/gysmo" "$ASCII_URL"
    fi
fi

# Copy missing files from /usr/share/gysmo to ~/.config/gysmo
echo "Checking for missing files in $HOME/.config/gysmo..."
mkdir -p "$CONFIGURATION_PATH"
mkdir -p "$ASCII_PATH"
mkdir -p "$DATA_PATH"
mkdir -p "$SCHEMA_PATH"

if [ ! -f "$CONFIG_FILE" ]; then
    echo "Copying default configuration file to $CONFIG_FILE..."
    cp "$SHARE_PATH/config/config.json" "$CONFIG_FILE"
fi

if [ ! -f "$SCHEMA_FILE" ]; then
    echo "Copying default schema file to $SCHEMA_FILE..."
    cp "$SHARE_PATH/config/schema/config_schema.json" "$SCHEMA_FILE"
fi

if [ ! -f "$ASCII_FILE" ]; then
    echo "Copying default ASCII file to $ASCII_FILE..."
    cp "$SHARE_PATH/ascii/gysmo" "$ASCII_FILE"
fi

# Download and install the binary
echo "Downloading binary $RELEASE_URL..."
if command_exists wget; then
    wget -O gysmo "$RELEASE_URL"
elif command_exists curl; then
    curl -o gysmo "$RELEASE_URL"
fi

chmod +x gysmo

# Ensure the local bin directory exists
mkdir -p "$INSTALL_DIR"

if ! mv gysmo "$INSTALL_PATH"; then
    echo -e "${RED}Error: Failed to move the binary to $INSTALL_PATH.${RESET}"
    exit 1
fi

echo -e "${GREEN}Installation completed successfully.${RESET}"
echo -e "${BLUE}Please ensure $INSTALL_DIR is in your PATH.${RESET}"
