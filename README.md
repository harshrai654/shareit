# ShareIT

A file sharing CLI based on QR codes and local area network.

## Download and Install CLI and Server Executables

You can download and install the latest CLI and Server executables using the following command:

### For Linux

```bash
bash -c 'REPO="harshrai654/shareit"; DEST_DIR="$HOME/.local/bin/shareit"; ASSETS=("shareit.cli.linux" "shareit.server.linux"); mkdir -p $DEST_DIR; chmod -R 775 $DEST_DIR; for ASSET in "${ASSETS[@]}"; do URL=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep browser_download_url | grep $ASSET | cut -d\" -f4); curl -L -o "$DEST_DIR/$ASSET" "$URL"; chmod +x "$DEST_DIR/$ASSET"; export PATH=$PATH:$DEST_DIR; done'
```

### For MacOS

```sh
sh -c 'REPO="harshrai654/shareit"; DEST_DIR="$HOME/.local/bin/shareit"; ASSETS=("shareit.cli.darwin" "shareit.server.darwin"); mkdir -p $DEST_DIR; chmod -R 775 $DEST_DIR; for ASSET in "${ASSETS[@]}"; do URL=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep browser_download_url | grep $ASSET | cut -d\" -f4); curl -L -o "$DEST_DIR/$ASSET" "$URL"; chmod +x "$DEST_DIR/$ASSET"; export PATH=$PATH:$DEST_DIR; done'
```

### For Windows (Powershell)

Run following commands in Powershell with admin access
In windows 11 you can open terminal, right click the powershell tab and run it as administrator and then paste the below commands":

```powershell
$repo = "harshrai654/shareit"
$destDir = "$env:USERPROFILE\AppData\Local\shareit\bin"
$assets = @("shareit.cli.windows.exe", "shareit.server.windows.exe")

if (-Not (Test-Path -Path $destDir)) {
    New-Item -ItemType Directory -Path $destDir -Force
}

try {
    $releaseInfo = Invoke-RestMethod -Uri "https://api.github.com/repos/$repo/releases/latest"
} catch {
    Write-Error "Failed to get release information from GitHub: $_"
    exit 1
}

foreach ($asset in $assets) {
    $url = $releaseInfo.assets | Where-Object { $_.name -eq $asset } | Select-Object -ExpandProperty browser_download_url
    if ($url) {
        $destPath = Join-Path -Path $destDir -ChildPath $asset
        Invoke-WebRequest -Uri $url -OutFile $destPath
        Write-Host "Downloaded $asset to $destPath"
    } else {
        Write-Host "Asset $asset not found in the latest release."
    }
}
try {
    New-NetFirewallRule -DisplayName "Allow Port 8965" -Direction Inbound -Protocol TCP -LocalPort 8965 -Action Allow -ErrorAction Stop
    Write-Host "Firewall rule created to allow port 8965"
} catch {
    Write-Error "Failed to create firewall rule: $_"
}

# Add $destDir to the PATH environment variable if it's not already there
$currentPath = [System.Environment]::GetEnvironmentVariable("PATH", [System.EnvironmentVariableTarget]::User)
if (-Not ($currentPath -split ";" | ForEach-Object { $_.Trim() } | Where-Object { $_ -eq $destDir })) {
    $newPath = "$currentPath;$destDir"
    [System.Environment]::SetEnvironmentVariable("PATH", $newPath, [System.EnvironmentVariableTarget]::User)
    Write-Host "Added $destDir to the PATH environment variable"
} else {
    Write-Host "$destDir is already in the PATH environment variable"
}
```

## Usage

To share a file to devices in the local network you will need the Absolute address of the file in the local machine.

```sh
shareit.cli.<darwin|linux|windows>[.exe] -filepath "/path/to/file"
```
