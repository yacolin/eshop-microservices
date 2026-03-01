# Install Go dependencies for eshop-microservices
# Run from project root: .\scripts\install-deps.ps1

$env:GOPROXY = "https://proxy.golang.org,direct"
Set-Location $PSScriptRoot\..

Write-Host "Running go mod tidy..."
go mod tidy
if ($LASTEXITCODE -ne 0) {
    Write-Host "If you see 'git not found': install Git and add it to PATH."
    Write-Host "If you see proxy errors: ensure GOPROXY is not set to a broken proxy."
    exit 1
}
Write-Host "Done. You can now build: go build ./cmd/order-service"
