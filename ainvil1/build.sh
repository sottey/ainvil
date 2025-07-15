#!/usr/bin/env bash

set -e

# Set version manually or read from file/tag/env
VERSION="1.0.0"

# Get latest Git commit hash (short form)
COMMIT=$(git rev-parse --short HEAD)

# Get UTC build date
BUILD_DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)

echo "Building ainvil..."
echo "Version:    $VERSION"
echo "Commit:     $COMMIT"
echo "Build Date: $BUILD_DATE"

# Build the binary with ldflags injected
go build -o ainvil -ldflags "\
  -X github.com/sottey/ainvil/cmd.Version=$VERSION \
  -X github.com/sottey/ainvil/cmd.Commit=$COMMIT \
  -X github.com/sottey/ainvil/cmd.BuildDate=$BUILD_DATE"
