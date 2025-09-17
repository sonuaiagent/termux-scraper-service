#!/bin/bash

echo "Starting cleanup..."

# Clean NPM cache
echo "Cleaning NPM cache..."
npm cache clean --force 2>/dev/null || echo "NPM not available"

# Remove common cache directories
echo "Removing cache files..."
rm -rf ~/.cache/node-gyp/
rm -rf ~/.npm/_cacache/
rm -rf ~/.cache/matplotlib/

# Remove temporary files
echo "Removing temporary files..."
find ~ -name "*.tmp" -type f -delete 2>/dev/null
find ~ -name "*~" -type f -delete 2>/dev/null
find ~ -name ".DS_Store" -type f -delete 2>/dev/null

# Show disk space saved
echo "Cleanup completed!"
