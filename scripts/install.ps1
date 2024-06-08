Write-Host "Installing ShareIT for Windows..."

# Variables
$GITHUB_REPO = "harshrai654/shareit"
$CLI_EXECUTABLE = "shareit.cli.windows.exe"
$SERVER_EXECUTABLE = "shareit.server.windows.exe"
$INSTALL_DIR = "$env:ProgramFiles\ShareIT"
$RUNTIME_DIR = "$env:LOCALAPPDATA\ShareIT"

# Create necessary directories
New-Item -ItemType Directory -Path $INSTALL_DIR -Force | Out-Null
New-Item -ItemType Directory -Path $RUNTIME_DIR -Force | Out-Null

# Download latest release from GitHub
$releaseInfo = Invoke-RestMethod -Uri "https://api.github.com/repos/$GITHUB_REPO/releases/latest"

$CLI_URL = $releaseInfo.assets | Where-Object { $_.name -like "*$CLI_EXECUTABLE*" } | Select-Object -ExpandProperty browser_download_url
$SERVER_URL = $releaseInfo.assets | Where-Object { $_.name -like "*$SERVER_EXECUTABLE*" } | Select-Object -ExpandProperty browser_download_url

# Download the executables
Write-Host "Downloading executables..."
Write-Host "CLI: $CLI_URL"
Invoke-WebRequest -Uri $CLI_URL -OutFile "$INSTALL_DIR\$CLI_EXECUTABLE"

Write-Host "Server: $SERVER_URL"
Invoke-WebRequest -Uri $SERVER_URL -OutFile "$INSTALL_DIR\$SERVER_EXECUTABLE"

# Add INSTALL_DIR to PATH
Write-Host "Adding $INSTALL_DIR to PATH..."
$oldPath = [System.Environment]::GetEnvironmentVariable("Path", [System.EnvironmentVariableTarget]::Machine)
if ($oldPath -notlike "*$INSTALL_DIR*") {
    $newPath = "$oldPath;$INSTALL_DIR"
    [System.Environment]::SetEnvironmentVariable("Path", $newPath, [System.EnvironmentVariableTarget]::Machine)
}

Write-Host "Installation complete. Please restart your terminal or log off and log back in for PATH changes to take effect."