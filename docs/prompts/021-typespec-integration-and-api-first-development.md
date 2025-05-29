# TypeSpec Integration and API-First Development

## Prompt Name: TypeSpec Integration and API-First Development

## Context
You need to integrate TypeSpec into an existing Go-based template system to enable API-first development with automatic schema generation, OpenAPI documentation, and multi-language client generation.

## TypeSpec Integration Strategy

### 1. Hybrid Architecture Approach
```
Hybrid TypeSpec Integration:
├── Existing Go Templates (Preserve)
│   ├── Working CLI tool
│   ├── Proven template system
│   └── Production-ready generation
├── TypeSpec Layer (Add)
│   ├── API definitions
│   ├── Schema generation
│   └── Client generation
└── Unified Workflow (Combine)
    ├── TypeSpec → OpenAPI → Go types
    ├── Go templates use generated types
    └── Multi-language client support
```

### 2. TypeSpec Definition Structure
```typescript
// typespec/health.tsp
import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

using TypeSpec.Http;
using TypeSpec.Rest;

@service({
  title: "Health API",
  version: "1.0.0",
})
namespace HealthAPI;

// Core health status model
model HealthStatus {
  /** Service health status */
  status: "healthy" | "unhealthy" | "degraded";
  
  /** Timestamp of health check */
  timestamp: utcDateTime;
  
  /** Service version */
  version: string;
  
  /** Service uptime duration */
  uptime: duration;
  
  /** Server timing metrics */
  serverTiming?: ServerTimingMetrics;
  
  /** OpenTelemetry trace ID */
  traceId?: string;
  
  /** Dependency health status */
  dependencies?: DependencyStatus[];
}

// Server timing metrics
model ServerTimingMetrics {
  /** Database query time in milliseconds */
  dbQuery?: float64;
  
  /** Cache lookup time in milliseconds */
  cacheHit?: float64;
  
  /** Total processing time in milliseconds */
  total?: float64;
  
  /** External API call time in milliseconds */
  externalApi?: float64;
}

// Dependency health status
model DependencyStatus {
  /** Dependency name */
  name: string;
  
  /** Dependency status */
  status: "healthy" | "unhealthy" | "timeout";
  
  /** Response time in milliseconds */
  responseTime?: float64;
  
  /** Error message if unhealthy */
  error?: string;
}

// Server time response
model ServerTime {
  /** Current server timestamp in RFC3339 format */
  timestamp: utcDateTime;
  
  /** Server timezone */
  timezone: string;
  
  /** Unix timestamp in seconds */
  unix: int64;
  
  /** Unix timestamp in milliseconds */
  unixMilli: int64;
  
  /** ISO 8601 formatted timestamp */
  iso8601: string;
  
  /** Human-readable format */
  formatted: string;
  
  /** Server timing metrics */
  serverTiming?: ServerTimingMetrics;
  
  /** OpenTelemetry trace ID */
  traceId?: string;
}

// Health API endpoints
@route("/health")
interface Health {
  /** Check overall service health */
  @get check(): HealthStatus;
  
  /** Get current server time */
  @get @route("/time") serverTime(): ServerTime;
  
  /** Check service readiness */
  @get @route("/ready") readiness(): HealthStatus;
  
  /** Check service liveness */
  @get @route("/live") liveness(): HealthStatus;
  
  /** Check startup status */
  @get @route("/startup") startup(): HealthStatus;
}

// CloudEvents integration
model HealthEvent {
  /** CloudEvents specification version */
  specversion: "1.0";
  
  /** Event type */
  type: "com.health.status.changed";
  
  /** Event source */
  source: string;
  
  /** Event ID */
  id: string;
  
  /** Event time */
  time: utcDateTime;
  
  /** Event data */
  data: HealthStatus;
}
```

### 3. TypeSpec Tier Definitions
```typescript
// typespec/tiers/basic.tsp
import "../health.tsp";

using HealthAPI;

// Basic tier - core functionality only
@service({
  title: "Basic Health API",
  version: "1.0.0",
})
namespace BasicHealthAPI;

// Simplified health status for basic tier
model BasicHealthStatus {
  status: "healthy" | "unhealthy";
  timestamp: utcDateTime;
  version: string;
  uptime: duration;
}

@route("/health")
interface BasicHealth {
  @get check(): BasicHealthStatus;
  @get @route("/time") serverTime(): ServerTime;
}
```

```typescript
// typespec/tiers/enterprise.tsp
import "../health.tsp";

using HealthAPI;

// Enterprise tier - full feature set
@service({
  title: "Enterprise Health API",
  version: "1.0.0",
})
namespace EnterpriseHealthAPI;

// Enterprise health status with security context
model EnterpriseHealthStatus extends HealthStatus {
  /** Security context */
  securityContext?: SecurityContext;
  
  /** Compliance status */
  compliance?: ComplianceStatus;
  
  /** Audit information */
  audit?: AuditInfo;
}

// Security context
model SecurityContext {
  /** Client identity */
  clientId?: string;
  
  /** User roles */
  roles?: string[];
  
  /** Permissions */
  permissions?: string[];
  
  /** mTLS certificate info */
  certificateInfo?: CertificateInfo;
}

// Compliance status
model ComplianceStatus {
  /** SOC2 compliance */
  soc2: boolean;
  
  /** HIPAA compliance */
  hipaa: boolean;
  
  /** GDPR compliance */
  gdpr: boolean;
  
  /** Last audit date */
  lastAudit?: utcDateTime;
}

@route("/health")
interface EnterpriseHealth extends Health {
  /** Enterprise health check with security */
  @get check(): EnterpriseHealthStatus;
  
  /** Security status check */
  @get @route("/security") security(): SecurityContext;
  
  /** Compliance status check */
  @get @route("/compliance") compliance(): ComplianceStatus;
}
```

## Integration Implementation

### 1. TypeSpec Build Pipeline
```bash
#!/bin/bash
# scripts/build-typespec.sh

set -e

echo "🔧 Building TypeSpec definitions..."

# Install TypeSpec if not present
if ! command -v tsp &> /dev/null; then
    echo "Installing TypeSpec..."
    npm install -g @typespec/compiler
    npm install -g @typespec/http
    npm install -g @typespec/rest
    npm install -g @typespec/openapi3
fi

# Create output directories
mkdir -p generated/openapi
mkdir -p generated/json-schema
mkdir -p generated/go-types
mkdir -p generated/typescript

# Compile TypeSpec for each tier
for tier in basic intermediate advanced enterprise; do
    echo "📋 Compiling $tier tier TypeSpec..."
    
    # Generate OpenAPI spec
    tsp compile typespec/tiers/${tier}.tsp \
        --emit @typespec/openapi3 \
        --output-dir generated/openapi/${tier}
    
    # Generate JSON Schema
    tsp compile typespec/tiers/${tier}.tsp \
        --emit @typespec/json-schema \
        --output-dir generated/json-schema/${tier}
    
    echo "✅ $tier tier compiled successfully"
done

echo "🎉 TypeSpec compilation complete!"
```

### 2. Go Type Generation
```go
// pkg/codegen/typespec.go
package codegen

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "text/template"
)

// TypeSpecGenerator generates Go types from TypeSpec-generated schemas
type TypeSpecGenerator struct {
    schemaDir string
    outputDir string
}

// GenerateGoTypes generates Go types from JSON Schema
func (g *TypeSpecGenerator) GenerateGoTypes(tier string) error {
    schemaPath := filepath.Join(g.schemaDir, tier, "schema.json")
    
    // Read JSON Schema
    schemaData, err := os.ReadFile(schemaPath)
    if err != nil {
        return fmt.Errorf("failed to read schema: %w", err)
    }
    
    var schema JSONSchema
    if err := json.Unmarshal(schemaData, &schema); err != nil {
        return fmt.Errorf("failed to parse schema: %w", err)
    }
    
    // Generate Go types
    goTypes := g.generateGoTypesFromSchema(schema)
    
    // Write Go file
    outputPath := filepath.Join(g.outputDir, tier, "types.go")
    return g.writeGoFile(outputPath, goTypes, tier)
}

// Go type template
const goTypeTemplate = `// Code generated by TypeSpec. DO NOT EDIT.

package {{.Package}}

import (
    "time"
)

{{range .Types}}
// {{.Name}} {{.Description}}
type {{.Name}} struct {
{{range .Fields}}    {{.Name}} {{.Type}} ` + "`json:\"{{.JSONName}}{{if .Optional}},omitempty{{end}}\"`" + `{{if .Comment}} // {{.Comment}}{{end}}
{{end}}}

{{end}}
`

func (g *TypeSpecGenerator) writeGoFile(path string, types []GoType, pkg string) error {
    tmpl, err := template.New("types").Parse(goTypeTemplate)
    if err != nil {
        return err
    }
    
    data := struct {
        Package string
        Types   []GoType
    }{
        Package: pkg,
        Types:   types,
    }
    
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()
    
    return tmpl.Execute(file, data)
}
```

### 3. Template Integration
```go
// Update existing templates to use generated types
const healthHandlerTemplate = `package handlers

import (
    "encoding/json"
    "net/http"
    "time"
    
    {{if hasFeature "opentelemetry"}}
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
    {{end}}
    
    // Import generated types
    "{{.Config.GoModule}}/internal/types"
)

type HealthHandler struct {
    {{if hasFeature "opentelemetry"}}
    tracer trace.Tracer
    {{end}}
}

func (h *HealthHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
    {{if hasFeature "opentelemetry"}}
    ctx, span := h.tracer.Start(r.Context(), "health.check")
    defer span.End()
    {{else}}
    ctx := r.Context()
    {{end}}
    
    // Use TypeSpec-generated types
    status := &types.{{.TierConfig.HealthStatusType}}{
        Status:    "healthy",
        Timestamp: time.Now(),
        Version:   "1.0.0",
        Uptime:    time.Since(startTime),
    }
    
    {{if hasFeature "opentelemetry"}}
    if span.SpanContext().HasTraceID() {
        traceID := span.SpanContext().TraceID().String()
        status.TraceID = &traceID
    }
    {{end}}
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(status)
}
`
```

### 4. Client Generation
```bash
#!/bin/bash
# scripts/generate-clients.sh

set -e

echo "🔧 Generating API clients..."

# Generate TypeScript client
for tier in basic intermediate advanced enterprise; do
    echo "📋 Generating TypeScript client for $tier tier..."
    
    # Use OpenAPI Generator
    openapi-generator generate \
        -i generated/openapi/${tier}/openapi.yaml \
        -g typescript-fetch \
        -o generated/typescript/${tier} \
        --additional-properties=typescriptThreePlus=true,supportsES6=true
    
    # Generate Go client
    echo "📋 Generating Go client for $tier tier..."
    
    openapi-generator generate \
        -i generated/openapi/${tier}/openapi.yaml \
        -g go \
        -o generated/go-client/${tier} \
        --additional-properties=packageName=${tier}client
    
    echo "✅ $tier tier clients generated"
done

echo "🎉 Client generation complete!"
```

## Workflow Integration

### 1. Enhanced CLI Commands
```go
// Add TypeSpec support to CLI
var generateCmd = &cobra.Command{
    Use:   "generate",
    Short: "Generate project from template",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Check if TypeSpec mode is enabled
        useTypeSpec, _ := cmd.Flags().GetBool("typespec")
        
        if useTypeSpec {
            return generateFromTypeSpec(cmd, args)
        }
        
        // Use existing template generation
        return generateFromTemplate(cmd, args)
    },
}

func generateFromTypeSpec(cmd *cobra.Command, args []string) error {
    tier, _ := cmd.Flags().GetString("tier")
    
    // Generate TypeSpec definitions
    if err := buildTypeSpec(tier); err != nil {
        return err
    }
    
    // Generate Go types from schema
    if err := generateGoTypes(tier); err != nil {
        return err
    }
    
    // Generate project using enhanced templates
    return generateProjectWithTypeSpec(tier)
}
```

### 2. Template Enhancement
```yaml
# Template metadata with TypeSpec support
# templates/advanced/template.yaml
name: advanced
description: Advanced tier with full TypeSpec integration
tier: advanced
typespec:
  enabled: true
  schema_path: "typespec/tiers/advanced.tsp"
  generates:
    - openapi
    - json-schema
    - go-types
    - typescript-client
features:
  opentelemetry: true
  cloudevents: true
  kubernetes: true
  typescript: true
  typespec: true
version: "1.0.0"
```

### 3. Validation Pipeline
```bash
#!/bin/bash
# scripts/validate-typespec.sh

set -e

echo "🧪 Validating TypeSpec integration..."

# Validate TypeSpec syntax
for tier in basic intermediate advanced enterprise; do
    echo "📋 Validating $tier TypeSpec..."
    
    tsp compile typespec/tiers/${tier}.tsp --no-emit
    
    if [ $? -eq 0 ]; then
        echo "✅ $tier TypeSpec syntax valid"
    else
        echo "❌ $tier TypeSpec syntax invalid"
        exit 1
    fi
done

# Validate generated schemas
echo "📋 Validating generated schemas..."
for tier in basic intermediate advanced enterprise; do
    schema_file="generated/json-schema/${tier}/schema.json"
    
    if [ -f "$schema_file" ]; then
        # Validate JSON Schema syntax
        python3 -c "import json; json.load(open('$schema_file'))"
        echo "✅ $tier schema valid"
    else
        echo "❌ $tier schema missing"
        exit 1
    fi
done

# Test generated Go types compile
echo "📋 Testing generated Go types..."
for tier in basic intermediate advanced enterprise; do
    if [ -d "generated/go-types/${tier}" ]; then
        (cd "generated/go-types/${tier}" && go build .)
        echo "✅ $tier Go types compile"
    fi
done

echo "🎉 TypeSpec validation complete!"
```

## Migration Strategy

### 1. Gradual Integration
```
Phase 1: Add TypeSpec Layer
├── Keep existing Go templates working
├── Add TypeSpec definitions alongside
└── Generate schemas and documentation

Phase 2: Enhanced Templates
├── Update templates to use generated types
├── Add client generation capabilities
└── Maintain backward compatibility

Phase 3: Full Integration
├── TypeSpec as primary source of truth
├── Generated types drive template logic
└── Multi-language client support
```

### 2. Backward Compatibility
```go
// Support both modes in configuration
type ProjectConfig struct {
    // Existing fields
    Name     string
    Tier     string
    GoModule string
    
    // TypeSpec integration
    TypeSpec TypeSpecConfig `yaml:"typespec"`
}

type TypeSpecConfig struct {
    Enabled     bool   `yaml:"enabled"`
    SchemaPath  string `yaml:"schema_path"`
    GenerateGo  bool   `yaml:"generate_go"`
    GenerateTS  bool   `yaml:"generate_ts"`
}
```

## Success Criteria

### Technical Integration
- ✅ TypeSpec compiles without errors
- ✅ Generated schemas are valid
- ✅ Go types compile successfully
- ✅ TypeScript clients work correctly

### Workflow Enhancement
- ✅ CLI supports both modes
- ✅ Templates use generated types
- ✅ Documentation auto-generated
- ✅ Client SDKs available

### Developer Experience
- ✅ API-first development enabled
- ✅ Schema validation automated
- ✅ Multi-language support
- ✅ Backward compatibility maintained

This TypeSpec integration enables API-first development while preserving the existing working system, providing the best of both approaches.
