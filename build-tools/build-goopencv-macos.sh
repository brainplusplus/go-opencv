#!/bin/bash
# build-goopencv-macos.sh — Build goopencv.dylib for macOS (universal: arm64 + x86_64)
# Statically links OpenCV into the dylib — same approach as Windows build.
# Result: standalone goopencv.dylib, no runtime framework dependency.
# Requires: clang++, wget, unzip
# Usage: bash build-tools/build-goopencv-macos.sh

set -euo pipefail

cd "$(dirname "$0")/.."

OPENCV_VERSION="4.13.0"
OPENCV_PKG="opencv-mobile-${OPENCV_VERSION}-macos"
OPENCV_URL="https://github.com/nihui/opencv-mobile/releases/latest/download/${OPENCV_PKG}.zip"

# Download opencv-mobile if not already present
if [ ! -d "build-tools/${OPENCV_PKG}" ]; then
    echo "=== Downloading ${OPENCV_PKG} ==="
    wget -q "${OPENCV_URL}" -O "/tmp/${OPENCV_PKG}.zip"
    unzip -q -o "/tmp/${OPENCV_PKG}.zip" -d "build-tools/"
    rm "/tmp/${OPENCV_PKG}.zip"
    echo "Downloaded and extracted."
fi

OPENCV_ROOT="build-tools/${OPENCV_PKG}"
FRAMEWORK="${OPENCV_ROOT}/opencv2.framework"
FRAMEWORK_BIN="${FRAMEWORK}/opencv2"
OUTPUT="dist/goopencv.dylib"
SOURCE="backend/goopencv_abi.cpp"

echo "=== Compiling goopencv.dylib ==="
echo "Source:    ${SOURCE}"
echo "Framework: ${FRAMEWORK}"
echo "Output:    ${OUTPUT}"

mkdir -p dist

# Detect framework type: static archive vs dynamic library
FRAMEWORK_TYPE=$(file "${FRAMEWORK_BIN}" 2>/dev/null || echo "unknown")
echo "Framework binary type: ${FRAMEWORK_TYPE}"

if echo "${FRAMEWORK_TYPE}" | grep -q "archive\|current ar archive"; then
    echo ">>> Static framework detected — using -force_load for standalone dylib"
    clang++ -shared -O2 -fPIC -std=c++11 \
        -I"${FRAMEWORK}/Headers" \
        -force_load "${FRAMEWORK_BIN}" \
        "${SOURCE}" \
        -o "${OUTPUT}" \
        -lpthread \
        -framework Cocoa \
        -install_name @rpath/goopencv.dylib
else
    echo ">>> Dynamic framework detected — using -force_load (may work on universal Mach-O)"
    # -force_load works on both static archives and Mach-O universal binaries.
    # It forces all symbols to be loaded, effectively static-linking the framework.
    clang++ -shared -O2 -fPIC -std=c++11 \
        -I"${FRAMEWORK}/Headers" \
        -force_load "${FRAMEWORK_BIN}" \
        "${SOURCE}" \
        -o "${OUTPUT}" \
        -lpthread \
        -framework Cocoa \
        -install_name @rpath/goopencv.dylib
fi

# Verify: check that goopencv.dylib does NOT depend on opencv2 at runtime
echo ""
echo "=== Checking dynamic dependencies ==="
DYLIB_DEPS=$(otool -L "${OUTPUT}" 2>/dev/null || true)
if echo "${DYLIB_DEPS}" | grep -q "opencv2"; then
    echo "WARNING: goopencv.dylib still has dynamic dependency on opencv2!"
    echo "  ${DYLIB_DEPS}" | grep "opencv2"
    echo ""
    echo "Falling back to @loader_path approach..."
    # Change the opencv2 reference to @loader_path so it looks next to the dylib
    INSTALL_NAME=$(otool -L "${OUTPUT}" | grep "opencv2" | awk '{print $1}' | tr -d '[:space:]')
    if [ -n "${INSTALL_NAME}" ]; then
        install_name_tool -change "${INSTALL_NAME}" @loader_path/opencv2 "${OUTPUT}"
        # Copy the framework binary next to goopencv.dylib
        cp "${FRAMEWORK_BIN}" "dist/opencv2"
        echo "Copied opencv2 binary to dist/opencv2"
        echo "You must embed BOTH dist/goopencv.dylib and dist/opencv2 in the Go binary."
    fi
else
    echo "✅ goopencv.dylib is standalone — no opencv2 runtime dependency"
fi

echo ""
echo "=== SUCCESS ==="
lipo -info "${OUTPUT}" 2>/dev/null || echo "(lipo not available, skipping arch info)"
ls -lh "${OUTPUT}"
