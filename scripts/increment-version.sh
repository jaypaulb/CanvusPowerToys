#!/bin/bash

# increment-version.sh
#
# Increments the patch version in internal/atoms/version/version.go
# Usage: ./scripts/increment-version.sh
# Example: 1.0.66 → 1.0.67

set -e

# Color codes for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Path to version file
VERSION_FILE="internal/atoms/version/version.go"

# Check if file exists
if [ ! -f "$VERSION_FILE" ]; then
    echo -e "${YELLOW}Error: $VERSION_FILE not found${NC}"
    exit 1
fi

# Extract current version (e.g., "1.0.66")
CURRENT_VERSION=$(grep -oP 'Version\s*=\s*"\K[^"]+' "$VERSION_FILE")

if [ -z "$CURRENT_VERSION" ]; then
    echo -e "${YELLOW}Error: Could not find Version in $VERSION_FILE${NC}"
    exit 1
fi

# Parse version parts
IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT_VERSION"

# Increment patch version
NEW_PATCH=$((PATCH + 1))
NEW_VERSION="$MAJOR.$MINOR.$NEW_PATCH"

# Update the version file
sed -i "s/Version     = \"$CURRENT_VERSION\"/Version     = \"$NEW_VERSION\"/" "$VERSION_FILE"

# Verify the change was made
UPDATED_VERSION=$(grep -oP 'Version\s*=\s*"\K[^"]+' "$VERSION_FILE")

if [ "$UPDATED_VERSION" = "$NEW_VERSION" ]; then
    echo -e "${GREEN}✓ Version incremented successfully${NC}"
    echo -e "${BLUE}  $CURRENT_VERSION → $NEW_VERSION${NC}"
    echo ""
    echo -e "${BLUE}File: $VERSION_FILE${NC}"
    echo -e "${BLUE}New version line:${NC}"
    grep "Version     =" "$VERSION_FILE"
    exit 0
else
    echo -e "${YELLOW}Error: Version update failed${NC}"
    echo "Expected: $NEW_VERSION"
    echo "Got: $UPDATED_VERSION"
    exit 1
fi
