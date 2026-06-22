#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
BUILD_DIR="$SCRIPT_DIR/build_nix"
EMBED_DIR="$PROJECT_DIR/native/libs"

# Detect arch
ARCH="$(uname -m)"
case "$ARCH" in
    x86_64|amd64) GOARCH="amd64" ;;
    aarch64|arm64) GOARCH="arm64" ;;
    *) echo "Unknown arch: $ARCH"; exit 1 ;;
esac

# Detect OS
OS="$(uname -s)"
case "$OS" in
    Linux)
        LIB_NAME="libgolibjpeg.so"
        GOOS="linux"
        ;;
    Darwin)
        LIB_NAME="libgolibjpeg.dylib"
        GOOS="darwin"
        ;;
    *) echo "Unknown OS: $OS"; exit 1 ;;
esac

echo "=== Building for $GOOS/$GOARCH ==="

# Configure
cmake -S "$SCRIPT_DIR" -B "$BUILD_DIR" \
    -DCMAKE_BUILD_TYPE=Release \
    -DCMAKE_POSITION_INDEPENDENT_CODE=ON

# Build
cmake --build "$BUILD_DIR" --config Release

# Copy to embed location
mkdir -p "$EMBED_DIR"
EMBED_NAME="golibjpeg_${GOOS}_${GOARCH}${LIB_NAME##*libgolibjpeg}"
cp "$BUILD_DIR/libgolibjpeg.$([ $OS = Darwin ] && echo dylib || echo so)" "$EMBED_DIR/$EMBED_NAME"

echo "=== Done: $EMBED_DIR/$EMBED_NAME ==="
