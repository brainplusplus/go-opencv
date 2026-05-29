$ErrorActionPreference = "Stop"

<#
.SYNOPSIS
Regenerate contract-derived artifacts from internal/contract.

.EXAMPLE
.\scripts\generate-artifacts.ps1
#>

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")

Write-Host "Generating docs/abi.md..." -ForegroundColor Cyan
& go run "./cmd/goopencv-gen" -format abi-md -out "docs/abi.md"
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Write-Host "Generating backend/goopencv_abi.cpp..." -ForegroundColor Cyan
& go run "./cmd/goopencv-gen" -format cpp -out "backend/goopencv_abi.cpp"
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Write-Host "Done. Review changes before committing." -ForegroundColor Green
