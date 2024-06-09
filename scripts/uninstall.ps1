Write-Host "Uninstalling ShareIT for Windows..."

# Variables
$CLI_EXECUTABLE = "shareit.cli.windows.exe"
$SERVER_EXECUTABLE = "shareit.server.windows.exe"
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

# Remove installed files and directories
Write-Host "Removing installed files..."
Remove-Item -Path "$INSTALL_DIR\$CLI_EXECUTABLE" -Force
Remove-Item -Path "$INSTALL_DIR\$SERVER_EXECUTABLE" -Force
Remove-Item -Path "$INSTALL_DIR" -Recurse -Force


# Remove runtime directory
Write-Host "Removing runtime directory..."
Remove-Item -Path $RUNTIME_DIR -Recurse -Force

# Remove INSTALL_DIR from PATH
Write-Host "Removing $INSTALL_DIR from PATH..."
$oldPath = [System.Environment]::GetEnvironmentVariable("Path", [System.EnvironmentVariableTarget]::Machine)
$newPath = ($oldPath -split ';') -notmatch [Regex]::Escape($INSTALL_DIR) -join ';'
[System.Environment]::SetEnvironmentVariable("Path", $newPath, [System.EnvironmentVariableTarget]::Machine)

# Define variables
$contextMenuName = "Share File with ShareIT"  # The name that appears in the context menu

# Registry key paths to delete
$shellKeyPath = "HKEY_CLASSES_ROOT\*\shell\$contextMenuName"
$commandKeyPath = "$shellKeyPath\command"

# Debugging output
Write-Host "Removing registry entry at $commandKeyPath"
Write-Host "Removing registry entry at $shellKeyPath"

try {
    # Remove the command key
    reg.exe DELETE "$commandKeyPath" /f | Out-Null

    # Remove the shell key
    reg.exe DELETE "$shellKeyPath" /f | Out-Null

    Write-Host "Context menu item '$contextMenuName' removed successfully."
} catch {
    Write-Error "Failed to remove registry entry: $_"
}


Write-Host "Uninstallation complete."
