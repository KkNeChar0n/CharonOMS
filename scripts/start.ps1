# CharonOMS Project Start Script
# PowerShell version

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "     CharonOMS Project Launcher" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Prompt for database password
$securePassword = Read-Host "Please enter MySQL database password" -AsSecureString
$BSTR = [System.Runtime.InteropServices.Marshal]::SecureStringToBSTR($securePassword)
$password = [System.Runtime.InteropServices.Marshal]::PtrToStringAuto($BSTR)

# Set environment variable
$env:DB_PASSWORD = $password

Write-Host ""
Write-Host "Starting project..." -ForegroundColor Green
Write-Host "Database user: root" -ForegroundColor Gray
Write-Host "Database password: ******" -ForegroundColor Gray
Write-Host "Database name: charonoms" -ForegroundColor Gray
Write-Host ""

# Change to project directory
Set-Location "d:\claude space\CharonOMS"

# Run project
go run cmd/server/main.go

Read-Host "Press Enter to exit"
