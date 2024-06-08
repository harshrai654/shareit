Write-Host "Uninstalling ShareIT for Windows..."

# Variables
$CLI_EXECUTABLE = "shareit.cli.exe"
$SERVER_EXECUTABLE = "shareit.server.exe"
$INSTALL_DIR = "$env:ProgramFiles\ShareIT"
$RUNTIME_DIR = "$env:LOCALAPPDATA\ShareIT"

# Function to check if a process is running
function Is-ProcessRunning {
    param (
        [string]$ProcessName
    )
    $process = Get-Process -Name $ProcessName -ErrorAction SilentlyContinue
    return $process -ne $null
}

# Check if CLI executable is running
if (Is-ProcessRunning -ProcessName "shareit.cli") {
    Write-Host "Error: CLI executable ($CLI_EXECUTABLE) is running. Please stop it before uninstalling."
    exit 1
}

# Check if Server executable is running
if (Is-ProcessRunning -ProcessName "shareit.server") {
    Write-Host "Server executable ($SERVER_EXECUTABLE) is running. Stopping it..."
    Stop-Process -Name "shareit.server" -Force
    Start-Sleep -Seconds 2
    if (Is-ProcessRunning -ProcessName "shareit.server") {
        Write-Host "Error: Unable to stop the server executable. Please stop it manually and try again."
        exit 1
    }
}

# Remove installed files
Write-Host "Removing installed files..."
Remove-Item -Path "$INSTALL_DIR\$CLI_EXECUTABLE" -Force
Remove-Item -Path "$INSTALL_DIR\$SERVER_EXECUTABLE" -Force

# Remove runtime directory
Write-Host "Removing runtime directory..."
Remove-Item -Path $RUNTIME_DIR -Recurse -Force

Write-Host "Uninstallation complete."
