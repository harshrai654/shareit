# Determine OS
OS="$(uname -s)"
case "${OS}" in
    Linux*)     OS=Linux;;
    Darwin*)    OS=Mac;;
    *)          echo "Unsupported OS: ${OS}"; exit 1;;
esac

echo "Installing ShareIT for ${OS}..."


# Variables
GITHUB_REPO="harshrai654/shareit"
CLI_EXECUTABLE="shareit.cli.linux"
SERVER_EXECUTABLE="shareit.server.linux"
INSTALL_DIR="/usr/local/bin"
RUNTIME_DIR="$HOME/.shareit/"

if [ "${OS}" == "Mac" ]; then
    RUNTIME_DIR="$HOME/Library/Application Support/shareit"
    CLI_EXECUTABLE="shareit.cli.darwin"
    SERVER_EXECUTABLE="shareit.server.darwin"
fi

# Create necessary directories
sudo mkdir -p "$INSTALL_DIR"
mkdir -p "$RUNTIME_DIR"

# Download latest release from GitHub
CLI_URL=$(curl -s https://api.github.com/repos/$GITHUB_REPO/releases/latest | grep browser_download_url | grep $CLI_EXECUTABLE | cut -d\" -f4)
SERVER_URL=$(curl -s https://api.github.com/repos/$GITHUB_REPO/releases/latest | grep browser_download_url | grep $SERVER_EXECUTABLE | cut -d\" -f4)

# Download the executables
echo "Downloading executables..."
echo "CLI: $CLI_URL"
sudo curl -L -o "$INSTALL_DIR/$CLI_EXECUTABLE" "$CLI_URL"

echo "Server: $SERVER_URL"
sudo curl -L -o "$INSTALL_DIR/$SERVER_EXECUTABLE" "$SERVER_URL"

# Set execute permissions
sudo chmod +x "$INSTALL_DIR/$CLI_EXECUTABLE"
sudo chmod +x "$INSTALL_DIR/$SERVER_EXECUTABLE"

echo "Installation complete."
