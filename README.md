# ShareIT

A file sharing CLI based on QR codes and local area network.

## Download and Install CLI and Server Executables

You can download and install the latest CLI and Server executables using the following command:

- For Linux

```bash
sudo apt-get install -y jq
sudo bash -c 'REPO="harshrai654/shareit"; DEST_DIR="/usr/local/bin"; ASSETS=("shareit.cli.linux" "shareit.server.linux"); for ASSET in "${ASSETS[@]}"; do URL=$(curl -s https://api.github.com/repos/$REPO/releases/tags/release-v-latest | jq -r ".assets[] | select(.name == \"$ASSET\") | .browser_download_url"); curl -L -o "$DEST_DIR/$ASSET" "$URL"; chmod +x "$DEST_DIR/$ASSET"; done'
```

- For MacOS

```sh
sudo sh -c 'REPO="harshrai654/shareit"; DEST_DIR="/usr/local/bin"; ASSETS=("shareit.cli.darwin" "shareit.server.darwin"); for ASSET in "${ASSETS[@]}"; do URL=$(curl -s https://api.github.com/repos/$REPO/releases/tags/release-v-latest | jq -r ".assets[] | select(.name == \"$ASSET\") | .browser_download_url"); curl -L -o "$DEST_DIR/$ASSET" "$URL"; chmod +x "$DEST_DIR/$ASSET"; done'
```

- For Windows

## Download and Install CLI and Server Executables (PowerShell)

You can download and install the latest CLI and Server executables using the following PowerShell command:

```powershell
$repo = "harshrai654/shareit"
$destDir = "$env:SystemDrive\bin"
$assets = @("shareit.cli.windows.exe", "shareit.server.windows.exe")

if (-Not (Test-Path -Path $destDir)) {
    New-Item -ItemType Directory -Path $destDir
}

$releaseInfo = Invoke-RestMethod -Uri "https://api.github.com/repos/$repo/releases/latest"

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

# Optionally, add $destDir to the PATH if it's not already there
$envPath = [System.Environment]::G
