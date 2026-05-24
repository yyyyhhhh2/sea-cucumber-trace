$ErrorActionPreference = "Stop"

$root = Split-Path -Parent $PSScriptRoot

Write-Host "Building frontend dist..."
Push-Location (Join-Path $root "frontend")
try {
  npm run build
} finally {
  Pop-Location
}

Write-Host "Building Linux backend binary..."
$backendBuildDir = Join-Path $root "backend\build"
New-Item -ItemType Directory -Force -Path $backendBuildDir | Out-Null

Push-Location (Join-Path $root "backend")
try {
  $env:CGO_ENABLED = "0"
  $env:GOOS = "linux"
  $env:GOARCH = "amd64"
  go build -o (Join-Path $backendBuildDir "server-linux-amd64") ./cmd/server
} finally {
  Remove-Item Env:CGO_ENABLED -ErrorAction SilentlyContinue
  Remove-Item Env:GOOS -ErrorAction SilentlyContinue
  Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
  Pop-Location
}

Write-Host "Docker assets ready."
