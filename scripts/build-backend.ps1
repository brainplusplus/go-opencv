$ErrorActionPreference = "Stop"

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
$backendDir = Join-Path $repoRoot "backend"
$buildDir = Join-Path $repoRoot "build-wasi"
$outDir = Join-Path $repoRoot "dist"

# Locate WASI SDK
$wasiSdk = $env:WASI_SDK_PATH
if (-not $wasiSdk -or -not (Test-Path $wasiSdk)) {
  $candidate = Get-ChildItem (Join-Path $repoRoot "build-tools") -Directory -Filter "wasi-sdk-*" -ErrorAction SilentlyContinue | Select-Object -First 1
  if ($candidate) { $wasiSdk = $candidate.FullName }
}
if (-not $wasiSdk -or -not (Test-Path $wasiSdk)) {
  Write-Error "WASI SDK not found. Set WASI_SDK_PATH or install to build-tools/wasi-sdk-*/"
  exit 1
}

$clang = Join-Path $wasiSdk "bin\clang++.exe"
$ld = Join-Path $wasiSdk "bin\wasm-ld.exe"
if (-not (Test-Path $clang)) { Write-Error "clang++ not found at $clang"; exit 1 }
if (-not (Test-Path $ld)) { Write-Error "wasm-ld not found at $ld"; exit 1 }

Write-Host "WASI SDK: $wasiSdk"

# Prepare build dir
if (Test-Path $buildDir) { Remove-Item -Recurse -Force $buildDir }
New-Item -ItemType Directory -Path $buildDir -Force | Out-Null

# Compile
Write-Host "`n=== Compiling ==="
$objFile = Join-Path $buildDir "goopencv_abi.o"
& $clang --target=wasm32-wasip1 -O2 -std=c++17 -c (Join-Path $backendDir "goopencv_abi.cpp") -o $objFile
if ($LASTEXITCODE -ne 0) { Write-Error "Compilation failed"; exit $LASTEXITCODE }
Write-Host "Object: $((Get-Item $objFile).Length) bytes"

# Link (bare wasm, no CRT)
Write-Host "`n=== Linking ==="
$wasmOut = Join-Path $buildDir "goopencv_wasm"
& $ld --no-entry --no-gc-sections --export-dynamic `
    --export=goopencv_mat_new `
    --export=goopencv_mat_delete `
    --export=goopencv_mat_rows `
    --export=goopencv_mat_cols `
    --export=goopencv_mat_type `
    --export=goopencv_mat_clone `
    --export=goopencv_imgproc_cvt_color `
    --export=goopencv_imgproc_resize `
    --export-memory `
    -o $wasmOut $objFile
if ($LASTEXITCODE -ne 0) { Write-Error "Linking failed"; exit $LASTEXITCODE }

# Copy to dist
if (-not (Test-Path $outDir)) { New-Item -ItemType Directory -Path $outDir -Force | Out-Null }
Copy-Item $wasmOut -Destination (Join-Path $repoRoot "dist\goopencv.wasm") -Force
$size = (Get-Item (Join-Path $repoRoot "dist\goopencv.wasm")).Length
Write-Host "`n=== Done ==="
Write-Host "Output: dist\goopencv.wasm ($([math]::Round($size/1KB, 1)) KB)"
