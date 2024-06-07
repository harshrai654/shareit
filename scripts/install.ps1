Write-Host "Installing ShareIT for Windows..."

# Variables
$GITHUB_REPO = "harshrai654/shareit"
$CLI_EXECUTABLE = "shareit.cli.exe"
$SERVER_EXECUTABLE = "shareit.server.exe"
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