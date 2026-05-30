#!/bin/bash
# build-goopencv-linux.sh — Build goopencv.so for Linux amd64
# Requires: g++, wget, unzip
# Usage: bash build-tools/build-goopencv-linux.sh

set -euo pipefail

cd "$(dirname "$0")/.."

OPENCV_VERSION="4.13.0"
OPENCV_PKG="opencv-mobile-${OPENCV_VERSION}-ubuntu-2204"
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
OUTPUT="dist/goopencv.so"
SOURCE="backend/goopencv_abi.cpp"

echo "=== Compiling goopencv.so ==="
echo "Source:  ${SOURCE}"
echo "Include: ${OPENCV_ROOT}/include"
echo "Lib:     ${OPENCV_ROOT}/lib"
echo "Output:  ${OUTPUT}"

mkdir -p dist

g++ -shared -O2 -fPIC -std=c++11 \
    -I"${OPENCV_ROOT}/include/opencv4" \
    "${SOURCE}" \
    -o "${OUTPUT}" \
    -L"${OPENCV_ROOT}/lib" \
    -Wl,--start-group \
    -lopencv_core \
    -lopencv_imgproc \
    -lopencv_features2d \
    -lopencv_photo \
    -lopencv_highgui \
    -lopencv_video \
    -Wl,--end-group \
    -lpthread

echo ""
echo "=== SUCCESS ==="
ls -lh "${OUTPUT}"
