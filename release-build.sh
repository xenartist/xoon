#!/bin/bash

set -e

# Default version
VERSION="0.0.0"

# Parse command line arguments
while [[ "$#" -gt 0 ]]; do
    case $1 in
        -v|--version) VERSION="$2"; shift ;;
        *) echo "Unknown parameter: $1"; exit 1 ;;
    esac
    shift
done

# Define other variables
OS="linux"
ARCH="amd64"
BUILD_DIR="xoon"
BINARY_NAME="xoon"
ARCHIVE_NAME="xoon-${VERSION}-${OS}-${ARCH}.tar.gz"

echo "Building version: $VERSION"

# Create build directory if it doesn't exist
if [ ! -d "$BUILD_DIR" ]; then
    echo "Creating build directory: $BUILD_DIR"
    mkdir "$BUILD_DIR"
fi

# Build the binary
echo "Building $BINARY_NAME for Linux (amd64)..."
GOOS=$OS GOARCH=$ARCH go build -o "$BUILD_DIR/$BINARY_NAME"

if [ $? -ne 0 ]; then
    echo "Build failed"
    exit 1
fi

# Create tar.gz archive
echo "Creating archive: $ARCHIVE_NAME"
tar -czvf "$ARCHIVE_NAME" "$BUILD_DIR"

echo "Build complete. Archive created: $ARCHIVE_NAME"