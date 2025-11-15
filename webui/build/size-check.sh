#!/bin/bash
# Size monitoring script for WebUI assets

set -e

PUBLIC_DIR="webui/public"
MAX_SIZE_KB=500
MAX_SIZE_BYTES=$((MAX_SIZE_KB * 1024))

echo "Checking WebUI asset sizes..."

total_size=0
file_count=0

while IFS= read -r file; do
    size=$(stat -f%z "$file" 2>/dev/null || stat -c%s "$file" 2>/dev/null)
    total_size=$((total_size + size))
    file_count=$((file_count + 1))
    
    size_kb=$((size / 1024))
    if [ $size -gt $MAX_SIZE_BYTES ]; then
        echo "⚠️  WARNING: $file is ${size_kb}KB (exceeds ${MAX_SIZE_KB}KB limit)"
    fi
done < <(find "$PUBLIC_DIR" -type f)

total_kb=$((total_size / 1024))
echo ""
echo "Total assets: ${file_count} files, ${total_kb}KB"

if [ $total_size -gt $MAX_SIZE_BYTES ]; then
    echo "❌ Total size ${total_kb}KB exceeds ${MAX_SIZE_KB}KB limit"
    exit 1
else
    echo "✅ Total size ${total_kb}KB is within ${MAX_SIZE_KB}KB limit"
fi

