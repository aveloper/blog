#!/bin/sh

echo "Fetching latest version information"

# Fetch the latest to ensure we have the latest tag
git fetch origin

# Get the version from the git tag
# If not available, default to v0.0.0
# https://stackoverflow.com/a/33217295/4257791
VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")

# All the ldflags for the build
# Refer here to learn about ldflags -> https://www.digitalocean.com/community/tutorials/using-ldflags-to-set-version-information-for-go-applications
LDFLAGS="-X main.Version=$VERSION"

OUTPUT='bin/'

echo "Building app $VERSION"

go build -ldflags="$LDFLAGS" -o "$OUTPUT" ./cmd/blog || exit 1

echo "Build successful. Output -> ${OUTPUT}blog"
