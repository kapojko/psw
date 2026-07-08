# Build the project
Write-Host "Building psw.exe..."
go build ./cmd/psw
if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed"
    exit 1
}

# Test if the file exists
if (-not (Test-Path "psw.exe")) {
    Write-Host "psw.exe not found"
    exit 1
}

# Test if the destination directory exists
if (-not (Test-Path "C:\Programs\Util")) {
    Write-Host "Destination directory not found"
    exit 1
}

Write-Host "Copying psw.exe to C:\Programs\Util..."
Copy-Item -Path "psw.exe" -Destination "C:\Programs\Util"

Write-Host "Done"
