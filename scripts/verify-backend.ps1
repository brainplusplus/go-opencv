$ErrorActionPreference = "Stop"

<#
.SYNOPSIS
Verify go-opencv: Go tests + optional native backend check.

.DESCRIPTION
Runs Go unit tests. If a native shared library is found in dist\,
also shows its size.

.EXAMPLE
.\scripts\verify-backend.ps1
#>

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")

Write-Host "=== Go tests ===" -ForegroundColor Cyan
& go test -p 1 ./...
if ($LASTEXITCODE -ne 0) {
  Write-Error "Go tests failed."
  exit 1
}

Write-Host ""

$nativeCandidates = @(
  (Join-Path $repoRoot "dist\goopencv.dll"),
  (Join-Path $repoRoot "dist\goopencv.so"),
  (Join-Path $repoRoot "dist\goopencv-linux-arm64.so"),
  (Join-Path $repoRoot "dist\goopencv.dylib")
)

$nativePath = $nativeCandidates | Where-Object { Test-Path $_ } | Select-Object -First 1
if ($nativePath) {
  Write-Host "=== Native backend found: $nativePath ===" -ForegroundColor Cyan
  $size = (Get-Item $nativePath).Length
  Write-Host "Size: $([math]::Round($size/1MB, 2)) MB"
} else {
  Write-Host "=== No native library found in dist\. Skipping backend verification. ===" -ForegroundColor Yellow
  Write-Host "Run the platform build script in build-tools\ to produce one."
}

Write-Host ""
Write-Host "=== Done ===" -ForegroundColor Green
