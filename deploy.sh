#!/bin/bash

# Quick deployment script for Prometheus Metrics Activity
# This script builds the activity and provides next steps for Flogo app integration

set -e

echo "🚀 Prometheus Metrics Activity - Quick Deploy"
echo "=============================================="

# Step 1: Build the activity
echo "Step 1: Building activity..."
./build.sh

echo ""
echo "Step 2: Preparing for Flogo application integration..."

# Step 2: Clean module cache
echo "🧹 Cleaning Go module cache to ensure fresh module resolution..."
go clean -modcache

# Step 3: Check for Flogo applications that need rebuilding
echo ""
echo "Step 3: Checking for Flogo applications to rebuild..."

FLOGO_APPS_FOUND=false

# Check for .flogo files in parent directories
if ls ../../*.flogo 1> /dev/null 2>&1; then
    echo "📁 Flogo applications found:"
    ls -1 ../../*.flogo | sed 's/^/  - /'
    FLOGO_APPS_FOUND=true
fi

# Check for compiled binaries
if [ -d "../../bin" ] && [ "$(ls -A ../../bin 2>/dev/null)" ]; then
    echo "🔧 Compiled Flogo binaries found:"
    ls -1 ../../bin/ | grep -v "\\.sh$" | sed 's/^/  - /'
    FLOGO_APPS_FOUND=true
fi

echo ""
if [ "$FLOGO_APPS_FOUND" = true ]; then
    echo "⚠️  IMPORTANT: You must rebuild your Flogo application(s) to use the updated activity!"
    echo ""
    echo "🔄 Choose your rebuild method:"
    echo ""
    echo "Method A - Remove binary and rebuild:"
    echo "  rm ../../bin/your-app-name"
    echo "  # Then rebuild using your normal process"
    echo ""
    echo "Method B - If using Flogo CLI:"
    echo "  cd ../.."
    echo "  flogo build -f your-app.flogo"
    echo ""
    echo "Method C - Development Environment:"
    echo "  1. Stop your Flogo application"
    echo "  2. Rebuild in your IDE/development environment"
    echo "  3. Restart the application"
else
    echo "ℹ️  No Flogo applications detected in parent directory."
    echo "   Make sure to rebuild your Flogo application wherever it's located."
fi

echo ""
echo "🎯 Final Steps:"
echo "=============="
echo "1. Rebuild your Flogo application (see methods above)"
echo "2. Restart the Flogo application"
echo "3. Test your flow to verify the updated activity works"
echo ""
echo "✅ Expected output: Multi-line Prometheus metrics with all numeric fields"
echo ""
echo "🔍 Verification:"
echo "==============="
echo "Look for DEBUG statements in logs like:"
echo "  'DEBUG: Processing field 'cpu_usage' with value '75''"
echo "  'DEBUG: Number of metric lines generated: 3'"
echo ""
echo "📖 For more details, see README.md"
