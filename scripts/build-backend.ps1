$ErrorActionPreference = "Stop"

param(
    [ValidateSet("amd64", "arm64")]
    [string]$TargetArch = "amd64"
)

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
$outDir = Join-Path $repoRoot "dist"

Write-Host "TargetArch: $TargetArch"

# Delegate to the canonical native build script
$buildScript = Join-Path $repoRoot "build-tools\build-goopencv.bat"
if (Test-Path $buildScript) {
    Write-Host "`n=== Running native build ==="
    & cmd /c "set TARGET_ARCH=$TargetArch && `"$buildScript`""
    if ($LASTEXITCODE -ne 0) { Write-Error "Build failed"; exit $LASTEXITCODE }
} else {
    Write-Error "build-tools/build-goopencv.bat not found"
    exit 1
}

# Show result
$dllPath = Join-Path $outDir "goopencv.dll"
if (Test-Path $dllPath) {
    $size = (Get-Item $dllPath).Length
    Write-Host "`n=== Done ==="
    Write-Host "Output: dist\goopencv.dll ($([math]::Round($size/1MB, 2)) MB)"
} else {
    Write-Error "Build completed but dist\goopencv.dll not found"
    exit 1
}
