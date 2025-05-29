#!/bin/bash

# Build script for {{.Config.Name}}

set -e

echo "🔨 Building {{.Config.Name}}..."

# Clean previous builds
rm -rf bin/
mkdir -p bin/

# Build the application
go build -o bin/{{.Config.Name}} cmd/server/main.go

echo "✅ Build complete: bin/{{.Config.Name}}"
