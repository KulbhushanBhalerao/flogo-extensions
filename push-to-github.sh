#!/bin/bash

# Script to push prometheus-metrics to GitHub
# Usage: ./push-to-github.sh

set -e

echo "🚀 Pushing Prometheus Metrics Activity to GitHub..."
echo "=================================================="

# Ensure we're in the right directory
cd "$(dirname "$0")"

# Check if git is already initialized
if [ ! -d ".git" ]; then
    echo "📦 Initializing Git repository..."
    git init
    
    echo "📝 Creating .gitignore..."
    cat > .gitignore << 'EOF'
# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with \`go test -c\`
*.test

# Output of the go coverage tool
*.out

# Go workspace file
go.work

# IDE files
.vscode/
.idea/
*.swp
*.swo
*~

# OS files
.DS_Store
Thumbs.db

# Build artifacts
dist/
build/

# Temporary files
*.tmp
*.temp
EOF

    echo "➕ Adding files to Git..."
    git add .
    
    echo "💾 Creating initial commit..."
    git commit -m "Initial commit: TIBCO Flogo Prometheus Metrics Activity

- Custom activity for converting JSON to Prometheus metrics format
- Multi-line output with proper formatting
- Support for multiple numeric fields as separate metrics
- Configurable metric types (gauge, counter, histogram, summary)
- Automatic label generation and sanitization
- Comprehensive documentation and build scripts"

else
    echo "📂 Git repository already exists"
    
    # Add any new or modified files
    echo "➕ Adding new/modified files..."
    git add .
    
    # Check if there are changes to commit
    if git diff --staged --quiet; then
        echo "ℹ️  No changes to commit"
    else
        echo "💾 Committing changes..."
        git commit -m "Update: $(date '+%Y-%m-%d %H:%M:%S')

- Updated activity code and documentation
- Enhanced build and deployment scripts
- Improved multi-line output formatting"
    fi
fi

# Ask user which repository structure they prefer
echo ""
echo "🤔 Choose repository structure:"
echo "1) Push to flogo-extensions repository as prometheus-metrics branch"
echo "2) Push to flogo-extensions repository as main branch"
echo "3) Just prepare git (don't push yet)"
echo ""
read -p "Enter choice (1/2/3): " choice

case $choice in
    1)
        echo "🌿 Setting up as branch in flogo-extensions repository..."
        if ! git remote | grep -q "origin"; then
            git remote add origin https://github.com/kulbhushanbhalerao/flogo-extensions.git
        else
            git remote set-url origin https://github.com/kulbhushanbhalerao/flogo-extensions.git
        fi
        fi
        
        echo "🌿 Creating prometheus-metrics branch..."
        git checkout -b prometheus-metrics 2>/dev/null || git checkout prometheus-metrics
        
        echo "🚀 Pushing to GitHub..."
        git push -u origin prometheus-metrics
        
        echo "✅ Successfully pushed to: https://github.com/kulbhushanbhalerao/flogo-extensions/tree/prometheus-metrics"
        ;;
    2)
        echo "📦 Setting up as main branch in flogo-extensions repository..."
        if ! git remote | grep -q "origin"; then
            git remote add origin https://github.com/kulbhushanbhalerao/flogo-extensions.git
        else
            git remote set-url origin https://github.com/kulbhushanbhalerao/flogo-extensions.git
        fi
        fi
        
        echo "🌿 Setting main branch..."
        git branch -M main
        
        echo "🚀 Pushing to GitHub..."
        git push -u origin main
        
        echo "✅ Successfully pushed to: https://github.com/kulbhushanbhalerao/flogo-extensions"
        ;;
    3)
        echo "📋 Git repository prepared. Manual push commands:"
        echo ""
        echo "For flogo-extensions repository (as branch):"
        echo "  git remote add origin https://github.com/kulbhushanbhalerao/flogo-extensions.git"
        echo "  git checkout -b prometheus-metrics"
        echo "  git push -u origin prometheus-metrics"
        echo ""
        echo "For flogo-extensions repository (as main):"
        echo "  git remote add origin https://github.com/kulbhushanbhalerao/flogo-extensions.git"
        echo "  git branch -M main"
        echo "  git push -u origin main"
        ;;
    *)
        echo "❌ Invalid choice. Repository prepared but not pushed."
        exit 1
        ;;
esac

echo ""
echo "🎉 Done! Your Prometheus Metrics Activity is now on GitHub!"
echo ""
echo "📖 Next steps:"
echo "1. Update your Flogo applications to use the GitHub URL"
echo "2. Update import statements to reference the GitHub repository"
echo "3. Test the activity from the GitHub source"
