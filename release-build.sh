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

# Output names with version
RELEASE_NAME="xoon-linux-x64-v${VERSION}"
DIR_NAME="xoon-${VERSION}"

echo "Building version ${VERSION}..."

# Backup the original main.rs
echo "Backing up main.rs..."
cp src/main.rs src/main.rs.bak

# Update version in main.rs
echo "Updating version in main.rs..."
sed -i "s/title(\"xoon\")/title(\"xoon-${VERSION}\")/" src/main.rs

# Clean previous builds
rm -rf ${DIR_NAME}
rm -f ${RELEASE_NAME}.tar.gz

# Build release version
echo "Building release version..."
cargo build --release --target x86_64-unknown-linux-gnu

# Create release directory
echo "Creating release directory..."
mkdir ${DIR_NAME}

# Copy files
echo "Copying files..."
cp target/x86_64-unknown-linux-gnu/release/xoon ${DIR_NAME}/
if [ -f README.md ]; then
    cp README.md ${DIR_NAME}/
fi

# Set permissions
echo "Setting permissions..."
chmod +x ${DIR_NAME}/xoon

# Create tar.gz
echo "Creating archive..."
tar -czf ${RELEASE_NAME}.tar.gz ${DIR_NAME}/

# Clean up
echo "Cleaning up..."
rm -rf ${DIR_NAME}

# Restore the original main.rs
echo "Restoring main.rs..."
mv src/main.rs.bak src/main.rs

echo "Done! Archive created: ${RELEASE_NAME}.tar.gz"