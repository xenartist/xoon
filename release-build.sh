#!/bin/bash

# Default version
VERSION="0.0.0"

# Parse command line arguments
while getopts "v:" opt; do
  case $opt in
    v)
      VERSION="$OPTARG"
      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      exit 1
      ;;
  esac
done

# Output file name with version
RELEASE_NAME="xoon-linux-x64-v${VERSION}"

echo "Building version ${VERSION}..."

# Clean previous builds
rm -rf xoon-release
rm -f ${RELEASE_NAME}.tar.gz

# Build release version
echo "Building release version..."
cargo build --release --target x86_64-unknown-linux-gnu

# Create release directory
echo "Creating release directory..."
mkdir xoon-release

# Copy files
echo "Copying files..."
cp target/x86_64-unknown-linux-gnu/release/xoon xoon-release/
if [ -f README.md ]; then
    cp README.md xoon-release/
fi

# Set permissions
echo "Setting permissions..."
chmod +x xoon-release/xoon

# Create tar.gz
echo "Creating archive..."
tar -czf ${RELEASE_NAME}.tar.gz xoon-release/

# Clean up
echo "Cleaning up..."
rm -rf xoon-release

echo "Done! Archive created: ${RELEASE_NAME}.tar.gz"