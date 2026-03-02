$ErrorActionPreference = "Stop"

if (!(Test-Path "./go.mod")) {
    Write-Host "Error: go.mod not found!" -ForegroundColor Red
    return
}

if (!(Test-Path "./dist")) { New-Item -ItemType Directory -Path "./dist" | Out-Null }

Write-Host "Checking dependencies..." -ForegroundColor Yellow
go mod tidy

$source = "./cmd/flibgolite"
$flags = "-s -w"

Write-Host "`n--- Building for Keenetic ---" -ForegroundColor Cyan

Write-Host "Hero (arm64)..."
$env:GOOS="linux"; $env:GOARCH="arm64"; $env:GOMIPS=""
go build -ldflags="-s -w" -o ./dist/flibgolite-hero $source

Write-Host "`n--- Building for Windows ---" -ForegroundColor Cyan

Write-Host "Windows x64..."
$env:GOOS="windows"; $env:GOARCH="amd64"
go build -ldflags="-s -w" -o ./dist/flibgolite-win64.exe $source

Write-Host "Windows x86..."
$env:GOOS="windows"; $env:GOARCH="386"
go build -ldflags="-s -w" -o ./dist/flibgolite-win32.exe $source

Write-Host "`nSuccessfully finished! Check 'dist' folder." -ForegroundColor Green