#!/bin/bash

# Build script for demo-enterprise

set -e

echo "🔨 Building demo-enterprise..."

# Clean previous builds
rm -rf bin/
mkdir -p bin/

# Build the application
go build -o bin/demo-enterprise cmd/server/main.go

echo "✅ Build complete: bin/demo-enterprise"
