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

# Find the framework — may be at different locations depending on zip structure
FRAMEWORK=""
for candidate in \
    "build-tools/${OPENCV_PKG}/opencv2.framework" \
    "build-tools/${OPENCV_PKG}/${OPENCV_PKG}/opencv2.framework" \
    "build-tools/${OPENCV_PKG}/opencv-mobile-${OPENCV_VERSION}-macos/opencv2.framework" \
    "build-tools/opencv2.framework"
do
    if [ -d "$candidate" ]; then
        FRAMEWORK="$candidate"
        break
    fi
done

if [ -z "$FRAMEWORK" ]; then
    echo "ERROR: opencv2.framework not found in build-tools/${OPENCV_PKG}/"
    echo "Directory contents:"
    find "build-tools/${OPENCV_PKG}/" -maxdepth 3 -type d 2>/dev/null | head -20
    exit 1
fi

# Find the framework binary — could be at Versions/A/opencv2 or directly opencv2
FRAMEWORK_BIN=""
for bin_candidate in \
    "${FRAMEWORK}/Versions/A/opencv2" \
    "${FRAMEWORK}/Versions/Current/opencv2" \
    "${FRAMEWORK}/opencv2"
do
    if [ -f "$bin_candidate" ] || [ -L "$bin_candidate" ]; then
        FRAMEWORK_BIN="$bin_candidate"
        break
    fi
done

# Find headers — prefer Versions/A/Headers (real directory) over symlinks
HEADER_DIR=""
for h_candidate in \
    "${FRAMEWORK}/Versions/A/Headers" \
    "${FRAMEWORK}/Versions/Current/Headers" \
    "${FRAMEWORK}/Headers"
do
    if [ -d "$h_candidate" ] && [ -f "${h_candidate}/core.hpp" ]; then
        HEADER_DIR="$h_candidate"
        break
    fi
done

OUTPUT="dist/goopencv.dylib"
SOURCE="backend/goopencv_abi.cpp"

echo "=== Compiling goopencv.dylib ==="
echo "Source:      ${SOURCE}"
echo "Framework:   ${FRAMEWORK}"
echo "Binary:      ${FRAMEWORK_BIN}"
echo "Headers:     ${HEADER_DIR}"
echo "Output:      ${OUTPUT}"

mkdir -p dist

if [ -z "$HEADER_DIR" ]; then
    echo "ERROR: Framework headers not found"
    find "${FRAMEWORK}" -maxdepth 3 -type d 2>/dev/null
    exit 1
fi

# Create opencv2 symlink in Headers dir so #include <opencv2/core.hpp> resolves.
# Framework headers are directly in Headers/ (e.g., Headers/core.hpp) but
# the code uses #include <opencv2/core.hpp>, so we need Headers/opencv2 -> .
if [ ! -d "${HEADER_DIR}/opencv2" ]; then
    ln -s . "${HEADER_DIR}/opencv2"
fi

if [ -n "$FRAMEWORK_BIN" ]; then
    FRAMEWORK_TYPE=$(file "${FRAMEWORK_BIN}" 2>/dev/null || echo "unknown")
    echo "Framework binary type: ${FRAMEWORK_TYPE}"

    # Use -force_load to statically link all symbols (standalone dylib)
    echo ">>> Using -force_load for standalone dylib"
    clang++ -shared -O2 -fPIC -std=c++11 \
        -I"${HEADER_DIR}" \
        -force_load "${FRAMEWORK_BIN}" \
        "${SOURCE}" \
        -o "${OUTPUT}" \
        -lpthread \
        -framework Cocoa \
        -framework Accelerate \
        -install_name @rpath/goopencv.dylib
else
    # No binary found — try dynamic framework link
    echo ">>> No framework binary found, trying dynamic -framework link"
    FRAMEWORK_DIR=$(dirname "${FRAMEWORK}")
    clang++ -shared -O2 -fPIC -std=c++11 \
        -I"${HEADER_DIR}" \
        -F"${FRAMEWORK_DIR}" \
        -framework opencv2 \
        "${SOURCE}" \
        -o "${OUTPUT}" \
        -lpthread \
        -framework Cocoa \
        -framework Accelerate \
        -install_name @rpath/goopencv.dylib
fi

# Verify: check that goopencv.dylib does NOT depend on opencv2 at runtime
echo ""
echo "=== Checking dynamic dependencies ==="
DYLIB_DEPS=$(otool -L "${OUTPUT}" 2>/dev/null || true)
if echo "${DYLIB_DEPS}" | grep -q "opencv2"; then
    echo "WARNING: goopencv.dylib still has dynamic dependency on opencv2!"
    echo "  ${DYLIB_DEPS}" | grep "opencv2"
else
    echo "✅ goopencv.dylib is standalone — no opencv2 runtime dependency"
fi

echo ""
echo "=== SUCCESS ==="
lipo -info "${OUTPUT}" 2>/dev/null || echo "(lipo not available, skipping arch info)"
ls -lh "${OUTPUT}"
