# Website Downloader - Build Script for Windows
# This script builds the production executable

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Website Downloader - Build Script" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Check if Wails CLI is installed
Write-Host "Checking for Wails CLI..." -ForegroundColor Yellow
$wailsVersion = wails version 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: Wails CLI not found!" -ForegroundColor Red
    Write-Host "Please install Wails: go install github.com/wailsapp/wails/v2/cmd/wails@latest" -ForegroundColor Red
    exit 1
}
Write-Host "Wails CLI found: $wailsVersion" -ForegroundColor Green
Write-Host ""

# Clean previous builds
Write-Host "Cleaning previous builds..." -ForegroundColor Yellow
if (Test-Path "build") {
    Remove-Item -Path "build" -Recurse -Force
    Write-Host "Previous build cleaned" -ForegroundColor Green
}
Write-Host ""

# Build for Windows
Write-Host "Building for Windows..." -ForegroundColor Yellow
Write-Host "This may take a few minutes..." -ForegroundColor Gray
Write-Host ""

wails build -platform windows/amd64

if ($LASTEXITCODE -eq 0) {
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Green
    Write-Host "  BUILD SUCCESSFUL!" -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Green
    Write-Host ""
    Write-Host "Executable location: build\bin\website-downloader-app.exe" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "You can now distribute this single executable file!" -ForegroundColor Yellow
    Write-Host "No dependencies required - just run the .exe file" -ForegroundColor Yellow
} else {
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Red
    Write-Host "  BUILD FAILED!" -ForegroundColor Red
    Write-Host "========================================" -ForegroundColor Red
    Write-Host ""
    Write-Host "Please check the error messages above" -ForegroundColor Red
    exit 1
}
