$ErrorActionPreference = "Stop"

<#
.SYNOPSIS
Verify go-opencv: Go tests + optional backend wasm smoke test.

.DESCRIPTION
Runs Go unit tests. If a wasm backend is found at dist\goopencv.wasm,
also runs integration-level checks against it.

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

$wasmPath = Join-Path $repoRoot "dist\goopencv.wasm"
if (Test-Path $wasmPath) {
  Write-Host "=== Backend wasm found: $wasmPath ===" -ForegroundColor Cyan
  $size = (Get-Item $wasmPath).Length
  Write-Host "Size: $([math]::Round($size/1KB, 1)) KB"
  Write-Host "Full integration test not yet implemented (needs backend with real exports)."
} else {
  Write-Host "=== No dist\goopencv.wasm found. Skipping backend verification. ===" -ForegroundColor Yellow
  Write-Host "Run scripts\build-backend.ps1 to produce one."
}

Write-Host ""
Write-Host "=== Done ===" -ForegroundColor Green
