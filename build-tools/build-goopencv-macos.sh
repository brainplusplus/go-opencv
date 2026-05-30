#!/bin/bash
# build-goopencv-macos.sh - Build a universal goopencv.dylib for macOS.
# Produces a standalone shared library with both arm64 and x86_64 slices.
# Requires: clang++, wget, unzip
# Usage: bash build-tools/build-goopencv-macos.sh

set -euo pipefail

cd "$(dirname "$0")/.."

OPENCV_VERSION="4.13.0"
OPENCV_PKG="opencv-mobile-${OPENCV_VERSION}-macos"
OPENCV_URL="https://github.com/nihui/opencv-mobile/releases/latest/download/${OPENCV_PKG}.zip"
OUTPUT="dist/goopencv.dylib"
SOURCE="backend/goopencv_abi.cpp"
ARCHS=(-arch arm64 -arch x86_64)

if [ ! -d "build-tools/${OPENCV_PKG}" ]; then
    echo "=== Downloading ${OPENCV_PKG} ==="
    wget -q "${OPENCV_URL}" -O "/tmp/${OPENCV_PKG}.zip"
    unzip -q -o "/tmp/${OPENCV_PKG}.zip" -d "build-tools/"
    rm "/tmp/${OPENCV_PKG}.zip"
    echo "Downloaded and extracted."
fi

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
    find "build-tools/${OPENCV_PKG}/" -maxdepth 3 -type d 2>/dev/null | head -20
    exit 1
fi

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

echo "=== Compiling goopencv.dylib ==="
echo "Source:        ${SOURCE}"
echo "Framework:     ${FRAMEWORK}"
echo "Binary:        ${FRAMEWORK_BIN}"
echo "Headers:       ${HEADER_DIR}"
echo "Output:        ${OUTPUT}"
echo "Architectures: arm64 + x86_64"

mkdir -p dist

if [ -z "$HEADER_DIR" ]; then
    echo "ERROR: Framework headers not found"
    find "${FRAMEWORK}" -maxdepth 3 -type d 2>/dev/null
    exit 1
fi

if [ ! -d "${HEADER_DIR}/opencv2" ]; then
    ln -s . "${HEADER_DIR}/opencv2"
fi

if [ -n "$FRAMEWORK_BIN" ]; then
    echo "Framework binary type: $(file "${FRAMEWORK_BIN}" 2>/dev/null || echo "unknown")"
    echo ">>> Using -force_load for standalone dylib"
    clang++ -shared -O2 -fPIC -std=c++11 \
        "${ARCHS[@]}" \
        -I"${HEADER_DIR}" \
        -force_load "${FRAMEWORK_BIN}" \
        "${SOURCE}" \
        -o "${OUTPUT}" \
        -lpthread \
        -framework Cocoa \
        -framework Accelerate \
        -install_name @rpath/goopencv.dylib
else
    echo ">>> No framework binary found, trying dynamic -framework link"
    FRAMEWORK_DIR=$(dirname "${FRAMEWORK}")
    clang++ -shared -O2 -fPIC -std=c++11 \
        "${ARCHS[@]}" \
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

echo ""
echo "=== Checking dynamic dependencies ==="
DYLIB_DEPS=$(otool -L "${OUTPUT}" 2>/dev/null || true)
if echo "${DYLIB_DEPS}" | grep -q "opencv2"; then
    echo "WARNING: goopencv.dylib still has dynamic dependency on opencv2"
    echo "${DYLIB_DEPS}" | grep "opencv2"
else
    echo "OK: goopencv.dylib is standalone"
fi

echo ""
echo "=== Verifying architectures ==="
if command -v lipo >/dev/null 2>&1; then
    LIPO_INFO=$(lipo -info "${OUTPUT}")
    echo "${LIPO_INFO}"
    if [[ "${LIPO_INFO}" != *"x86_64"* ]] || [[ "${LIPO_INFO}" != *"arm64"* ]]; then
        echo "ERROR: ${OUTPUT} is not universal"
        exit 1
    fi
else
    echo "(lipo not available, skipping architecture verification)"
fi

echo ""
echo "=== SUCCESS ==="
ls -lh "${OUTPUT}"
