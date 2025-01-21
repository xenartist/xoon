#!/bin/bash

# Default version
VERSION=""

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

# Check if version is provided
if [ -z "$VERSION" ]; then
    echo "Error: Version is required"
    echo "Usage: $0 -v VERSION"
    echo "Example: $0 -v 0.1.0"
    exit 1
fi

RELEASE_NAME="xoon-linux-x64-v${VERSION}"
RELEASE_FILE="${RELEASE_NAME}.tar.gz"

# Check if release file exists
if [ ! -f "$RELEASE_FILE" ]; then
    echo "Error: Release file $RELEASE_FILE not found!"
    echo "Please run release-build.sh first to create the release file."
    exit 1
fi

# Create GitHub release and upload file
echo "Creating GitHub release v${VERSION}..."
echo "Uploading file: $RELEASE_FILE"

gh release create "v${VERSION}" \
    --title "xoon v${VERSION}" \
    --notes "Release version ${VERSION}" \
    "$RELEASE_FILE"

if [ $? -eq 0 ]; then
    echo "Success! Release v${VERSION} has been created and $RELEASE_FILE has been uploaded."
else
    echo "Error: Failed to create release or upload file."
    exit 1
fi