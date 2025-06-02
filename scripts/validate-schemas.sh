#!/bin/bash

# Schema validation script for TypeSpec definitions
# This script validates all TypeSpec schemas and generates JSON Schema + OpenAPI v3

set -e

echo "🔍 Validating TypeSpec schemas..."

# Check if TypeSpec is installed
if ! command -v tsp &> /dev/null; then
    echo "❌ TypeSpec compiler not found. Installing..."
    npm install -g @typespec/compiler
fi

# Create output directory for generated schemas
mkdir -p generated/schemas
mkdir -p generated/openapi

echo "📋 Validating core health schemas..."

# Validate health.tsp
echo "  ✓ Validating health.tsp..."
tsp compile template-health/schemas/health.tsp --output-dir generated/schemas/health

# Validate server-time.tsp
echo "  ✓ Validating server-time.tsp..."
tsp compile template-health/schemas/server-time.tsp --output-dir generated/schemas/server-time

# Validate health-api.tsp
echo "  ✓ Validating health-api.tsp..."
tsp compile template-health/schemas/health-api.tsp --output-dir generated/schemas/health-api

# Validate cloudevents.tsp
echo "  ✓ Validating cloudevents.tsp..."
tsp compile template-health/schemas/cloudevents.tsp --output-dir generated/schemas/cloudevents

# Validate basic tier
echo "  ✓ Validating basic tier..."
tsp compile pkg/schemas/tiers/basic.tsp --output-dir generated/schemas/basic

echo "📊 Generating OpenAPI v3 specifications..."

# Generate OpenAPI for health API
tsp compile template-health/schemas/health-api.tsp --emit @typespec/openapi3 --output-dir generated/openapi/health-api

# Generate OpenAPI for basic tier
tsp compile pkg/schemas/tiers/basic.tsp --emit @typespec/openapi3 --output-dir generated/openapi/basic

echo "📋 Generating JSON Schemas..."

# Generate JSON Schema for models
tsp compile template-health/schemas/health.tsp --emit @typespec/json-schema --output-dir generated/schemas/json

echo "✅ All schemas validated successfully!"

# Display generated files
echo ""
echo "📁 Generated files:"
find generated -name "*.json" -o -name "*.yaml" -o -name "*.yml" | sort

echo ""
echo "🎉 Schema validation complete!"
echo "   - TypeSpec schemas: ✅ Valid"
echo "   - JSON Schema generation: ✅ Working"
echo "   - OpenAPI v3 generation: ✅ Working"
echo "   - Basic tier validation: ✅ Complete"
