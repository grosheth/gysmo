#!/bin/sh

RELEASE_URL="https://github.com/grosheth/gysmo/releases/download/Alpha-0.1.0/gysmo"
INSTALL_DIR="$HOME/bin"
INSTALL_PATH="$INSTALL_DIR/gysmo"
CONFIGURATION_PATH="$HOME/.config/gysmo"

command_exists() {
    command -v "$1" >/dev/null 2>&1
}

if ! command_exists wget; then
    echo "Error: wget is not installed. Please install wget and try again."
    exit 1
fi

if [ ! -d "$INSTALL_DIR" ]; then
    mkdir -p "$INSTALL_DIR"
    if [ $? -ne 0 ]; then
        echo "Error: Failed to create directory $INSTALL_DIR."
        exit 1
    fi
fi

echo ""
echo "Installation..."

if ! wget -O gysmo "$RELEASE_URL"; then
    echo "Error: Failed to download the release binary."
    exit 1
fi

chmod +x gysmo

if ! mv gysmo "$INSTALL_PATH"; then
    echo "Error: Failed to move the binary to $INSTALL_PATH."
    exit 1
fi

if [ ! -d "$CONFIGURATION_PATH" ]; then
    mkdir -p "$CONFIGURATION_PATH"
    if [ $? -ne 0 ]; then
        echo "Error: Failed to create directory $CONFIGURATION_PATH."
        exit 1
    fi
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
