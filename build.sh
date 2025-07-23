#!/bin/bash

# Build script for TIBCO Flogo Prometheus Metrics Activity

set -e  # Exit on any error

echo "🚀 Building Prometheus Metrics Activity..."
echo "============================================"

# Navigate to the activity directory
cd "$(dirname "$0")"

# Download dependencies
echo "📦 Downloading dependencies..."
go mod tidy

# Run tests (if test dependencies are available)
echo "🧪 Running tests..."
if go test -v ./...; then
    echo "✅ All tests passed!"
else
    echo "⚠️  Tests failed or skipped - continuing with build"
fi

# Build the activity
echo "🔨 Building activity..."
go build -v .

# Check for successful build
if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Activity build successful!"
    echo ""
    echo "📋 Next Steps:"
    echo "=============="
    echo "1. Clear Go module cache:     go clean -modcache"
    echo "2. Rebuild Flogo application: See instructions below"
    echo "3. Restart Flogo application"
    echo ""
    echo "🔄 Flogo Application Rebuild Options:"
    echo "======================================"
    echo ""
    echo "Option A - Using Flogo CLI:"
    echo "  cd /path/to/your/flogo/app"
    echo "  flogo build -f your-app.flogo"
    echo ""
    echo "Option B - Development Environment:"
    echo "  1. Stop current Flogo application"
    echo "  2. Clear build cache in your IDE/environment"
    echo "  3. Rebuild/restart the application"
    echo ""
    echo "Option C - Manual Binary Rebuild:"
    echo "  1. Remove old binary: rm bin/your-flogo-app"
    echo "  2. Rebuild using your build process"
    echo ""
    echo "🔧 Integration Details:"
    echo "======================"
    echo "  Module Path: github.com/kulbhushanbhalerao/activity/prometheus-metrics"
    echo "  Activity Reference: #prometheus-metrics"
    echo ""
    echo "📖 For detailed usage examples, see README.md"
    echo ""
    
    # Check if we're in a Flogo workspace and offer to clean cache
    if [ -f "../../go.mod" ] || [ -f "../go.mod" ]; then
        echo "🧹 Flogo workspace detected. Clean module cache? (y/n)"
        read -r response
        if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
            echo "🧹 Cleaning Go module cache..."
            go clean -modcache
            echo "✅ Module cache cleaned!"
        fi
    fi
    
    # Check if there are Flogo binaries to rebuild
    if [ -d "../../bin" ] && [ "$(ls -A ../../bin)" ]; then
        echo ""
        echo "🔍 Flogo binaries found in ../../bin/"
        echo "📋 Consider rebuilding these applications:"
        ls -1 ../../bin/ | grep -v "\\.sh$" | sed 's/^/  - /'
        echo ""
        echo "💡 To rebuild, remove the binary and rebuild your Flogo application"
    fi
else
    echo ""
    echo "❌ Build failed!"
    echo ""
    echo "🔍 Troubleshooting:"
    echo "=================="
    echo "1. Check Go version: go version (requires Go 1.19+)"
    echo "2. Verify dependencies: go mod tidy"
    echo "3. Check for syntax errors in activity.go"
    echo "4. Ensure all imports are available"
    echo ""
    exit 1
fi
