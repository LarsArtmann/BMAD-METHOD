#!/bin/bash

# Test script for {{.Config.Name}}

set -e

echo "🧪 Running tests for {{.Config.Name}}..."

# Run tests
go test -v ./...

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

echo "✅ Tests complete. Coverage report: coverage.html"
