#!/bin/bash

# Build script for LinkedIn Connector API

set -e

echo "🔨 Building LinkedIn Connector API..."

# Create bin directory if it doesn't exist
mkdir -p bin

# Get the directory of this script
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_ROOT"

# Build the application
echo "📦 Compiling Go binary..."
go build -o bin/api -ldflags="-s -w" cmd/api/main.go

# Check if build was successful
if [ -f "bin/api" ]; then
    echo "✅ Build successful!"
    echo "📍 Binary location: bin/api"
    echo ""
    echo "Run with: ./bin/api"
else
    echo "❌ Build failed!"
    exit 1
fi

