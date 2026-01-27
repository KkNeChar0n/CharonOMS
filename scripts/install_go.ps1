# Go Language Auto Installation Script for Windows
# PowerShell script to download and install Go

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Go Language Auto Installation Script" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Check if Go is already installed
$goInstalled = Get-Command go -ErrorAction SilentlyContinue
if ($goInstalled) {
    Write-Host "Go is already installed:" -ForegroundColor Yellow
    go version
    $continue = Read-Host "Continue installation? This will overwrite existing version (y/n)"
    if ($continue -ne 'y') {
        Write-Host "Installation cancelled" -ForegroundColor Red
        exit
    }
}

# Configuration
$goVersion = "1.21.6"
$goFileName = "go$goVersion.windows-amd64.msi"
$downloadUrl = "https://golang.google.cn/dl/$goFileName"
$installerPath = "$env:TEMP\$goFileName"

Write-Host "1. Preparing to download Go $goVersion..." -ForegroundColor Green

# Download Go installer
try {
    Write-Host "   Download URL: $downloadUrl" -ForegroundColor Gray
    Write-Host "   Save location: $installerPath" -ForegroundColor Gray
    Write-Host "   Downloading... (Size: ~140MB, please wait)" -ForegroundColor Yellow

    # Download with WebClient
    $webClient = New-Object System.Net.WebClient
    $webClient.DownloadFile($downloadUrl, $installerPath)

    Write-Host "   Download completed!" -ForegroundColor Green
}
catch {
    Write-Host "   Download failed, trying official mirror..." -ForegroundColor Yellow
    $downloadUrl = "https://go.dev/dl/$goFileName"
    try {
        $webClient = New-Object System.Net.WebClient
        $webClient.DownloadFile($downloadUrl, $installerPath)
        Write-Host "   Download completed!" -ForegroundColor Green
    }
    catch {
        Write-Host "   Download failed: $_" -ForegroundColor Red
        Write-Host "   Please download manually from: $downloadUrl" -ForegroundColor Yellow
        exit 1
    }
}

Write-Host ""
Write-Host "2. Installing Go $goVersion..." -ForegroundColor Green
Write-Host "   Installation wizard will open in a new window" -ForegroundColor Yellow
Write-Host "   Recommended installation path: C:\Program Files\Go" -ForegroundColor Gray

# Run installer
Start-Process -FilePath "msiexec.exe" -ArgumentList "/i `"$installerPath`" /qb" -Wait

Write-Host "   Installation completed!" -ForegroundColor Green
Write-Host ""

# Configure environment variables
Write-Host "3. Configuring environment variables..." -ForegroundColor Green

# Set GOPATH
$goPath = "$env:USERPROFILE\go"
Write-Host "   Setting GOPATH = $goPath" -ForegroundColor Gray
[System.Environment]::SetEnvironmentVariable("GOPATH", $goPath, "User")

# Add GOPATH\bin to PATH
$currentPath = [System.Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -notlike "*$goPath\bin*") {
    Write-Host "   Adding GOPATH\bin to PATH" -ForegroundColor Gray
    [System.Environment]::SetEnvironmentVariable("Path", "$currentPath;$goPath\bin", "User")
}

# Configure Go module proxy (China mirror)
Write-Host "   Configuring Go module proxy (China mirror)" -ForegroundColor Gray
[System.Environment]::SetEnvironmentVariable("GOPROXY", "https://goproxy.cn,direct", "User")
[System.Environment]::SetEnvironmentVariable("GO111MODULE", "on", "User")

Write-Host "   Environment variables configured!" -ForegroundColor Green
Write-Host ""

# Refresh current session environment variables
$machinePath = [System.Environment]::GetEnvironmentVariable("Path", "Machine")
$userPath = [System.Environment]::GetEnvironmentVariable("Path", "User")
$env:Path = "$machinePath;$userPath"
$env:GOPATH = $goPath
$env:GOPROXY = "https://goproxy.cn,direct"
$env:GO111MODULE = "on"

Write-Host "4. Verifying installation..." -ForegroundColor Green

# Wait for environment variables to take effect
Start-Sleep -Seconds 2

# Verify installation
try {
    $goVersionOutput = & "C:\Program Files\Go\bin\go.exe" version 2>&1
    Write-Host "   Go version: $goVersionOutput" -ForegroundColor Cyan

    Write-Host ""
    Write-Host "========================================" -ForegroundColor Green
    Write-Host "  Installation Successful!" -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Green
    Write-Host ""
    Write-Host "Important:" -ForegroundColor Yellow
    Write-Host "1. Please close current command line window" -ForegroundColor Yellow
    Write-Host "2. Open a new command line window" -ForegroundColor Yellow
    Write-Host "3. Run 'go version' to verify installation" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Next steps:" -ForegroundColor Cyan
    Write-Host "cd `"d:/claude space/CharonOMS`"" -ForegroundColor Gray
    Write-Host "go mod download" -ForegroundColor Gray
    Write-Host "go run cmd/server/main.go" -ForegroundColor Gray
    Write-Host ""
}
catch {
    Write-Host "   Cannot verify installation: $_" -ForegroundColor Yellow
    Write-Host "   Please close and reopen command line, then run 'go version'" -ForegroundColor Yellow
}

# Cleanup installer file
Write-Host "Cleaning up temporary files..." -ForegroundColor Gray
Remove-Item $installerPath -Force -ErrorAction SilentlyContinue

Write-Host "Done!" -ForegroundColor Green
Write-Host ""
Read-Host "Press Enter to exit"
