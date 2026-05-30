#!/bin/bash
# build-goopencv-linux-docker.sh — Build goopencv.so via Docker (run from Windows or Linux)
# This runs the Linux build inside an Ubuntu 22.04 container.
# No need to install gcc or opencv-mobile on the host.
#
# Usage:
#   bash build-tools/build-goopencv-linux-docker.sh

set -euo pipefail

cd "$(dirname "$0")/.."

echo "=== Building goopencv.so via Docker ==="

docker run --rm -v "$(pwd):/workspace" -w /workspace ubuntu:22.04 bash -c '
    set -euo pipefail

    # Install build tools
    apt-get update -qq
    apt-get install -y -qq wget unzip g++ > /dev/null

    OPENCV_VERSION="4.13.0"
    OPENCV_PKG="opencv-mobile-${OPENCV_VERSION}-ubuntu-2204"
    OPENCV_URL="https://github.com/nihui/opencv-mobile/releases/latest/download/${OPENCV_PKG}.zip"

    # Download opencv-mobile
    if [ ! -d "build-tools/${OPENCV_PKG}" ]; then
        echo "Downloading ${OPENCV_PKG}..."
        wget -q "${OPENCV_URL}" -O "/tmp/${OPENCV_PKG}.zip"
        unzip -q "/tmp/${OPENCV_PKG}.zip" -d "build-tools/"
        rm "/tmp/${OPENCV_PKG}.zip"
    fi

    OPENCV_ROOT="build-tools/${OPENCV_PKG}"
    OUTPUT="dist/goopencv.so"
    SOURCE="backend/goopencv_abi.cpp"

    mkdir -p dist

    echo "Compiling..."
    g++ -shared -O2 -fPIC -std=c++11 \
        -I"${OPENCV_ROOT}/include/opencv4" \
        "${SOURCE}" \
        -o "${OUTPUT}" \
        -L"${OPENCV_ROOT}/lib" \
        -lopencv_core \
        -lopencv_imgproc \
        -lopencv_features2d \
        -lopencv_photo \
        -lopencv_highgui \
        -lpthread

    echo ""
    echo "=== SUCCESS ==="
    ls -lh "${OUTPUT}"
'
