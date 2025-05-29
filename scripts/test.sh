#!/bin/bash

# Test script for BMAD Method Template Health Endpoint Generator

set -e

echo "🧪 Running BMAD Method Template Health Endpoint Generator Tests"
echo "=============================================================="

# Ensure we're in the project root
cd "$(dirname "$0")/.."

# Build the CLI tool first
echo "🔨 Building CLI tool..."
go build -o bin/template-health-endpoint cmd/generator/main.go

# Run unit tests
echo ""
echo "🔬 Running unit tests..."
go test -v ./pkg/generator/...

# Run integration tests
echo ""
echo "🔗 Running integration tests..."
go test -v ./tests/...

# Test CLI functionality
echo ""
echo "🖥️  Testing CLI functionality..."

# Test help command
echo "Testing CLI help..."
./bin/template-health-endpoint --help > /dev/null

# Test validate command
echo "Testing TypeSpec validation..."
./bin/template-health-endpoint validate --verbose

# Test dry run generation
echo "Testing dry run generation..."
./bin/template-health-endpoint generate \
  --name test-dry-run \
  --tier basic \
  --module github.com/example/test-dry-run \
  --dry-run

# Test actual generation and compilation
echo ""
echo "🏗️  Testing full project generation and compilation..."

# Clean up any existing test project
rm -rf test-full-generation

# Generate a test project
./bin/template-health-endpoint generate \
  --name test-full-generation \
  --tier basic \
  --module github.com/example/test-full-generation

# Test that it compiles
cd test-full-generation
echo "Running go mod tidy..."
go mod tidy

echo "Building generated project..."
go build -o bin/test-full-generation cmd/server/main.go

echo "Starting server for endpoint testing..."
./bin/test-full-generation &
SERVER_PID=$!

# Wait for server to start
sleep 3

# Test all endpoints
echo "Testing health endpoints..."
curl -f http://localhost:8080/health > /dev/null
curl -f http://localhost:8080/health/time > /dev/null  
curl -f http://localhost:8080/health/ready > /dev/null
curl -f http://localhost:8080/health/live > /dev/null
curl -f http://localhost:8080/health/startup > /dev/null

echo "All endpoints responded successfully!"

# Stop the server
kill $SERVER_PID

# Return to project root
cd ..

# Clean up test project
rm -rf test-full-generation

echo ""
echo "✅ All tests passed successfully!"
echo ""
echo "📊 Test Summary:"
echo "  ✅ Unit tests passed"
echo "  ✅ Integration tests passed" 
echo "  ✅ CLI functionality verified"
echo "  ✅ Project generation and compilation verified"
echo "  ✅ All health endpoints working (including /health/startup)"
echo ""
echo "🎉 BMAD Method Template Health Endpoint Generator is ready for production!"
