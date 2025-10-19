#!/bin/bash

# Run script for LinkedIn Connector API

set -e

echo "🚀 Starting LinkedIn Connector API..."

# Get the directory of this script
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_ROOT"

# Check if .env file exists in configs/
if [ ! -f "configs/.env" ]; then
    echo "⚠️  Warning: configs/.env not found"
    echo "📝 Creating from configs/.env.example..."
    if [ -f "configs/.env.example" ]; then
        cp configs/.env.example configs/.env
        echo "✅ Created configs/.env - please edit it with your configuration"
        echo "Press Ctrl+C to exit and edit, or Enter to continue with defaults..."
        read
    else
        echo "❌ Error: configs/.env.example not found"
        exit 1
    fi
fi

# Export environment variables from configs/.env
if [ -f "configs/.env" ]; then
    export $(cat configs/.env | grep -v '^#' | xargs)
fi

# Run the application
echo "🎯 Running API server..."
echo "📍 Port: ${PORT:-8080}"
echo ""

go run cmd/api/main.go

