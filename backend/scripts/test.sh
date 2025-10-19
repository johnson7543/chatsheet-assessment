#!/bin/bash

# Test script for LinkedIn Connector API

set -e

echo "🧪 Running tests for LinkedIn Connector API..."

# Get the directory of this script
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_ROOT"

# Run tests with coverage
echo "📊 Running tests with coverage..."
go test -v -cover ./...

# Generate coverage report
echo ""
echo "📈 Generating coverage report..."
go test -coverprofile=coverage.out ./...

if [ -f "coverage.out" ]; then
    echo "✅ Coverage report generated: coverage.out"
    echo ""
    echo "View coverage in browser:"
    echo "  go tool cover -html=coverage.out"
else
    echo "⚠️  Coverage report not generated"
fi

echo ""
echo "✅ Tests completed!"

