#!/bin/bash

# Build script for test-health-service

set -e

echo "🔨 Building test-health-service..."

# Clean previous builds
rm -rf bin/
mkdir -p bin/

# Build the application
go build -o bin/test-health-service cmd/server/main.go

echo "✅ Build complete: bin/test-health-service"
