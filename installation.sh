#!/bin/sh

REPO="grosheth/gysmo"
INSTALL_DIR="$HOME/bin"
INSTALL_PATH="$INSTALL_DIR/gysmo"
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
        echo "Error: wget or curl is not installed. Please install wget or curl and try again."
        exit 1
    fi
}

if ! command_exists wget && ! command_exists curl; then
    echo "Error: wget or curl is not installed. Please install wget or curl and try again."
    exit 1
fi

# Fetch the list of releases
echo "Fetching available versions..."
VERSIONS=$(fetch_releases)
if [ -z "$VERSIONS" ]; then
    echo "Error: Failed to fetch available versions."
    exit 1
fi

# Prompt the user to select a version
echo "Available versions:"
echo "$VERSIONS" | nl -w 2 -s '. '
read -rp "Enter the number of the version you want to install: " version_number

# Get the selected version
VERSION=$(echo "$VERSIONS" | sed -n "${version_number}p")
if [ -z "$VERSION" ]; then
    echo "Error: Invalid version selected."
    exit 1
fi

RELEASE_URL="https://github.com/$REPO/releases/download/$VERSION/gysmo"
CONFIG_URL="https://raw.githubusercontent.com/$REPO/refs/tags/$VERSION/gysmo/config/config.json"
SCHEMA_URL="https://raw.githubusercontent.com/$REPO/refs/tags/$VERSION/gysmo/config/schema/config_schema.json"
ASCII_URL="https://raw.githubusercontent.com/$REPO/refs/tags/$VERSION/gysmo/ascii/gysmo"

# Create the necessary directory structure
mkdir -p "$CONFIGURATION_PATH"
mkdir -p "$ASCII_PATH"
mkdir -p "$DATA_PATH"
mkdir -p "$SCHEMA_PATH"

# Download the config.json file if it doesn't already exist
if [ ! -f "$CONFIG_FILE" ]; then
    echo "Downloading configuration file..."
    if command_exists wget; then
        if ! wget -O "$CONFIG_FILE" "$CONFIG_URL"; then
            echo "Error: Failed to download configuration file using wget."
            exit 1
        fi
    elif command_exists curl; then
        if ! curl -o "$CONFIG_FILE" "$CONFIG_URL"; then
            echo "Error: Failed to download configuration file using curl."
            exit 1
        fi
    fi
    echo "Configuration file downloaded successfully."
else
    echo "Configuration file already exists. Skipping Config download."
fi

# Always download the schema file
echo "Downloading schema file..."
if command_exists wget; then
    if ! wget -O "$SCHEMA_FILE" "$SCHEMA_URL"; then
        echo "Error: Failed to download schema file using wget."
        exit 1
    fi
elif command_exists curl; then
    if ! curl -o "$SCHEMA_FILE" "$SCHEMA_URL"; then
        echo "Error: Failed to download schema file using curl."
        exit 1
    fi
fi
echo "Schema file downloaded successfully."

# Download the ASCII file if it doesn't already exist
if [ ! -f "$ASCII_FILE" ]; then
    echo "Downloading ASCII file..."
    if command_exists wget; then
        if ! wget -O "$ASCII_FILE" "$ASCII_URL"; then
            echo "Error: Failed to download ASCII file using wget."
            exit 1
        fi
    elif command_exists curl; then
        if ! curl -o "$ASCII_FILE" "$ASCII_URL"; then
            echo "Error: Failed to download ASCII file using curl."
            exit 1
        fi
    fi
    echo "ASCII file downloaded successfully."
else
    echo "ASCII file already exists. Skipping download."
fi

if [ ! -d "$INSTALL_DIR" ]; then
    if ! mkdir -p "$INSTALL_DIR"; then
        echo "Error: Failed to create directory $INSTALL_DIR."
        exit 1
    fi
fi

# Download the binary
echo "Downloading binary..."
if command_exists wget; then
    if ! wget -O gysmo "$RELEASE_URL"; then
        echo "Error: Failed to download the release binary using wget."
        exit 1
    fi
elif command_exists curl; then
    if ! curl -o gysmo "$RELEASE_URL"; then
        echo "Error: Failed to download the release binary using curl."
        exit 1
    fi
fi

chmod +x gysmo

if ! mv gysmo "$INSTALL_PATH"; then
    echo "Error: Failed to move the binary to $INSTALL_PATH."
    exit 1
fi

read -rp "Do you want to add $INSTALL_DIR to your PATH? (y/n): " response
case "$response" in
    [yY][eE][sS]|[yY])
        if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
            echo "export PATH=\$PATH:$INSTALL_DIR" >> "$HOME/.profile"
            echo "Added $INSTALL_DIR to PATH. Please restart your terminal or run 'source ~/.profile' to update your PATH."
        else
            echo "$INSTALL_DIR is already in your PATH."
        fi
        ;;
    *)
        echo "Skipping adding $INSTALL_DIR to PATH."
        ;;
esac

echo "Installation completed successfully."
echo ""
