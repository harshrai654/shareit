#!/bin/bash

# Determine OS
OS="$(uname -s)"
case "${OS}" in
    Linux*)     OS=Linux;;
    Darwin*)    OS=Mac;;
    *)          echo "Unsupported OS: ${OS}"; exit 1;;
esac

echo "Uninstalling ShareIT for ${OS}..."

# Variables
CLI_EXECUTABLE="shareit.cli.linux"
SERVER_EXECUTABLE="shareit.server.linux"
INSTALL_DIR="/usr/local/bin"
RUNTIME_DIR="$HOME/.shareit/"

if [ "${OS}" == "Mac" ]; then
    CLI_EXECUTABLE="shareit.cli.darwin"
    SERVER_EXECUTABLE="shareit.server.darwin"
fi

# Function to check if a process is running
is_running() {
    pgrep -f "$1" > /dev/null 2>&1
}

# Check if CLI executable is running
if is_running "$INSTALL_DIR/$CLI_EXECUTABLE"; then
    echo "Error: CLI executable ($CLI_EXECUTABLE) is running. Please stop it before uninstalling."
    exit 1
fi

# Check if Server executable is running
if is_running "$INSTALL_DIR/$SERVER_EXECUTABLE"; then
    echo "Server executable ($SERVER_EXECUTABLE) is running. Stopping it..."
    pkill -f "$INSTALL_DIR/$SERVER_EXECUTABLE"
    # Wait for the process to stop
    sleep 2
    if is_running "$INSTALL_DIR/$SERVER_EXECUTABLE"; then
        echo "Error: Unable to stop the server executable. Please stop it manually and try again."
        exit 1
    fi
fi

# Remove installed files
echo "Removing installed files..."
sudo rm -f "$INSTALL_DIR/$CLI_EXECUTABLE"
sudo rm -f "$INSTALL_DIR/$SERVER_EXECUTABLE"

# Remove runtime directory
echo "Removing runtime directory..."
rm -rf "$RUNTIME_DIR"

echo "Uninstallation complete."
