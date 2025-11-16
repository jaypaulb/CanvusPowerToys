#!/bin/bash
# Auto-increment version script
# Increments the patch version (1.0.3 -> 1.0.4)

set -e

VERSION_FILE="internal/atoms/version/version.go"
VERSIONINFO_FILE="versioninfo.json"

# Get current version
CURRENT_VERSION=$(grep 'Version.*=' "$VERSION_FILE" | sed -n 's/.*"\([^"]*\)".*/\1/p')

if [ -z "$CURRENT_VERSION" ]; then
    echo "Error: Could not find version in $VERSION_FILE"
    exit 1
fi

# Parse version components (assumes format: MAJOR.MINOR.PATCH)
IFS='.' read -r -a VERSION_PARTS <<< "$CURRENT_VERSION"
MAJOR=${VERSION_PARTS[0]}
MINOR=${VERSION_PARTS[1]}
PATCH=${VERSION_PARTS[2]}

# Increment patch version
NEW_PATCH=$((PATCH + 1))
NEW_VERSION="$MAJOR.$MINOR.$NEW_PATCH"

echo "Incrementing version: $CURRENT_VERSION -> $NEW_VERSION"

# Update version.go
sed -i "s/Version     = \"$CURRENT_VERSION\"/Version     = \"$NEW_VERSION\"/" "$VERSION_FILE"

# Update versioninfo.json - update Patch in both FileVersion and ProductVersion
sed -i "s/\"Patch\": $PATCH,/\"Patch\": $NEW_PATCH,/" "$VERSIONINFO_FILE"

# Update versioninfo.json - update string versions
sed -i "s/\"FileVersion\": \"$CURRENT_VERSION.0\"/\"FileVersion\": \"$NEW_VERSION.0\"/" "$VERSIONINFO_FILE"
sed -i "s/\"ProductVersion\": \"$CURRENT_VERSION.0\"/\"ProductVersion\": \"$NEW_VERSION.0\"/" "$VERSIONINFO_FILE"

echo "Version updated to $NEW_VERSION"

