#!/bin/bash
# build-goopencv-macos.sh — Build goopencv.dylib for macOS (arm64 + amd64)
# Requires: Xcode Command Line Tools (clang++), wget, unzip
# Usage: bash build-tools/build-goopencv-macos.sh

set -euo pipefail

cd "$(dirname "$0")/.."

OPENCV_VERSION="4.13.0"
OPENCV_PKG="opencv-mobile-${OPENCV_VERSION}-macos"
OPENCV_URL="https://github.com/nihui/opencv-mobile/releases/latest/download/${OPENCV_PKG}.zip"

# Detect architecture
ARCH=$(uname -m)
if [ "$ARCH" = "arm64" ]; then
    TARGET="arm64"
    ARCH_FLAG="-arch arm64"
elif [ "$ARCH" = "x86_64" ]; then
    TARGET="x86_64"
    ARCH_FLAG="-arch x86_64"
else
    echo "Unsupported architecture: $ARCH"
    exit 1
fi

echo "=== Building for macOS ${TARGET} ==="

# Download opencv-mobile if not already present
if [ ! -d "build-tools/${OPENCV_PKG}" ]; then
    echo "=== Downloading ${OPENCV_PKG} ==="
    wget -q "${OPENCV_URL}" -O "/tmp/${OPENCV_PKG}.zip"
    unzip -q "/tmp/${OPENCV_PKG}.zip" -d "build-tools/"
    rm "/tmp/${OPENCV_PKG}.zip"
    echo "Downloaded and extracted."
fi

OPENCV_ROOT="build-tools/${OPENCV_PKG}"
OUTPUT="dist/goopencv.dylib"
SOURCE="backend/goopencv_abi.cpp"

echo "=== Compiling goopencv.dylib ==="
echo "Source:  ${SOURCE}"
echo "Include: ${OPENCV_ROOT}/include"
echo "Lib:     ${OPENCV_ROOT}/lib"
echo "Output:  ${OUTPUT}"

mkdir -p dist

clang++ -shared -O2 -std=c++11 \
    ${ARCH_FLAG} \
    -I"${OPENCV_ROOT}/include" \
    "${SOURCE}" \
    -o "${OUTPUT}" \
    -L"${OPENCV_ROOT}/lib" \
    -lopencv_core \
    -lopencv_imgproc \
    -lopencv_features2d \
    -lopencv_photo \
    -lopencv_highgui \
    -lpthread \
    -lc++

# Create a universal binary if both architectures are available
if [ "$TARGET" = "x86_64" ] && [ -f "${OUTPUT}" ]; then
    echo "=== Creating universal binary (universal) ==="
    lipo -create -output "${OUTPUT}" "${OUTPUT}" 2>/dev/null || true
fi

echo ""
echo "=== SUCCESS ==="
ls -lh "${OUTPUT}"
