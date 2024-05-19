# ShareIT

A file sharing CLI based on QR codes and local area network.

## Download and Install CLI and Server Executables

You can download and install the latest CLI and Server executables using the following command:

- For Linux

```sh
sudo sh -c 'REPO="harshrai654/shareit"; DEST_DIR="/usr/local/bin"; ASSETS=("shareit.cli.linux" "shareit.server.linux"); for ASSET in "${ASSETS[@]}"; do URL=$(curl -s https://api.github.com/repos/$REPO/releases/tags/release-v-latest | jq -r ".assets[] | select(.name == \"$ASSET\") | .browser_download_url"); curl -L -o "$DEST_DIR/$ASSET" "$URL"; chmod +x "$DEST_DIR/$ASSET"; done'
```

- For MacOS

```sh
sudo sh -c 'REPO="harshrai654/shareit"; DEST_DIR="/usr/local/bin"; ASSETS=("shareit.cli.darwin" "shareit.server.darwin"); for ASSET in "${ASSETS[@]}"; do URL=$(curl -s https://api.github.com/repos/$REPO/releases/tags/release-v-latest | jq -r ".assets[] | select(.name == \"$ASSET\") | .browser_download_url"); curl -L -o "$DEST_DIR/$ASSET" "$URL"; chmod +x "$DEST_DIR/$ASSET"; done'
```

- For Windows

```cmd
@echo off
set REPO=harshrai654/shareit
set DEST_DIR=%SystemDrive%\bin
set ASSETS=shareit.cli.windows.exe,shareit.server.windows.exe

for %%A in (%ASSETS%) do (
    for /f "usebackq tokens=1,* delims=: " %%G in (`curl -s https://api.github.com/repos/%REPO%/releases/latest ^| findstr "browser_download_url.*%%~A"`) do (
        curl -L -o "%DEST_DIR%\%%~A" "%%H"
    )
)

echo Executables have been downloaded to %DEST_DIR%

```
