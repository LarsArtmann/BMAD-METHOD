#!/bin/bash

# Test script for test-health-service

set -e

echo "🧪 Running tests for test-health-service..."

# Run tests
go test -v ./...

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

echo "✅ Tests complete. Coverage report: coverage.html"
