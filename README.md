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

- For Windows

Run following commands in Powershell with admin access
In windows 11 you can open terminal, right click the powershell tab ti run it is administrator and then paste the below commands":

```powershell
$repo = "harshrai654/shareit"
$destDir = "%USERPROFILE%\AppData\Local\shareit\bin"
$assets = @("shareit.cli.windows.exe", "shareit.server.windows.exe")

for %%A in (%ASSETS%) do (
    for /f "usebackq tokens=1,* delims=: " %%G in (`curl -s https://api.github.com/repos/%REPO%/releases/latest ^| findstr "browser_download_url.*%%~A"`) do (
        curl -L -o "%DEST_DIR%\%%~A" "%%H"
    )
)

echo Executables have been downloaded to %DEST_DIR%

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

New-NetFirewallRule -DisplayName "Allow Port 8965" -Direction Inbound -Protocol TCP -LocalPort 8965 -Action Allow

# Add $destDir to the PATH environment variable if it's not already there
$env:PATH += ";$env:USERPROFILE\AppData\Local\shareit\bin"
```

## Usage

To share a file to devices in the local network you will need the Absolute address of the file in the local machine.

```sh
shareit.cli.<darwin|linux|windows>[.exe] -filepath "/path/to/file"
```
