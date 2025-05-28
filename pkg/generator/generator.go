package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/LarsArtmann/BMAD-METHOD/pkg/config"
)

// Generator handles the generation of health endpoint projects
type Generator struct {
	config    *config.ProjectConfig
	templates *TemplateRegistry
}

// TemplateRegistry manages all template files and functions
type TemplateRegistry struct {
	templates map[string]*template.Template
	functions template.FuncMap
}

// GenerationContext provides context for template execution
type GenerationContext struct {
	Config    *config.ProjectConfig
	Timestamp string
	Version   string
}

// New creates a new generator instance
func New(cfg *config.ProjectConfig) (*Generator, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Create template registry
	registry, err := NewTemplateRegistry()
	if err != nil {
		return nil, fmt.Errorf("failed to create template registry: %w", err)
	}

	return &Generator{
		config:    cfg,
		templates: registry,
	}, nil
}

// Generate generates the complete health endpoint project
func (g *Generator) Generate() error {
	// Create output directory
	if err := g.createOutputDirectory(); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create generation context
	ctx := &GenerationContext{
		Config:    g.config,
		Timestamp: "2024-01-01T00:00:00Z", // TODO: Use actual timestamp
		Version:   "1.0.0",
	}

	// Generate core files
	if err := g.generateCoreFiles(ctx); err != nil {
		return fmt.Errorf("failed to generate core files: %w", err)
	}

	// Generate Go source files
	if err := g.generateGoFiles(ctx); err != nil {
		return fmt.Errorf("failed to generate Go files: %w", err)
	}

	// Generate TypeScript files (if enabled)
	if g.config.Features.TypeScript {
		if err := g.generateTypeScriptFiles(ctx); err != nil {
			return fmt.Errorf("failed to generate TypeScript files: %w", err)
		}
	}

	// Generate Kubernetes files (if enabled)
	if g.config.Features.Kubernetes {
		if err := g.generateKubernetesFiles(ctx); err != nil {
			return fmt.Errorf("failed to generate Kubernetes files: %w", err)
		}
	}

	// Generate Docker files (if enabled)
	if g.config.Features.Docker {
		if err := g.generateDockerFiles(ctx); err != nil {
			return fmt.Errorf("failed to generate Docker files: %w", err)
		}
	}

	return nil
}

// createOutputDirectory creates the output directory structure
func (g *Generator) createOutputDirectory() error {
	dirs := []string{
		g.config.OutputDir,
		filepath.Join(g.config.OutputDir, "cmd", "server"),
		filepath.Join(g.config.OutputDir, "internal", "handlers"),
		filepath.Join(g.config.OutputDir, "internal", "models"),
		filepath.Join(g.config.OutputDir, "internal", "server"),
		filepath.Join(g.config.OutputDir, "internal", "config"),
		filepath.Join(g.config.OutputDir, "docs"),
		filepath.Join(g.config.OutputDir, "scripts"),
	}

	if g.config.Features.TypeScript {
		dirs = append(dirs,
			filepath.Join(g.config.OutputDir, "client", "typescript", "src"),
		)
	}

	if g.config.Features.Kubernetes {
		dirs = append(dirs,
			filepath.Join(g.config.OutputDir, "deployments", "kubernetes"),
		)
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// generateCoreFiles generates core project files
func (g *Generator) generateCoreFiles(ctx *GenerationContext) error {
	files := map[string]string{
		"README.md":           "readme",
		"go.mod":              "go-mod",
		".gitignore":          "gitignore",
		"Makefile":            "makefile",
		"docs/API.md":         "api-docs",
		"scripts/build.sh":    "build-script",
		"scripts/test.sh":     "test-script",
	}

	for filename, templateName := range files {
		if err := g.generateFile(filename, templateName, ctx); err != nil {
			return fmt.Errorf("failed to generate %s: %w", filename, err)
		}
	}

	return nil
}

// generateGoFiles generates Go source files
func (g *Generator) generateGoFiles(ctx *GenerationContext) error {
	files := map[string]string{
		"cmd/server/main.go":           "go-main",
		"internal/handlers/health.go":  "go-health-handler",
		"internal/models/health.go":    "go-health-models",
		"internal/server/server.go":    "go-server",
		"internal/config/config.go":    "go-config",
	}

	// Add tier-specific files
	switch g.config.Tier {
	case config.TierBasic:
		// Basic tier files already included
	case config.TierIntermediate, config.TierAdvanced, config.TierEnterprise:
		files["internal/handlers/server_time.go"] = "go-server-time-handler"
		files["internal/handlers/dependencies.go"] = "go-dependencies-handler"
	}

	// Add observability files for advanced tiers
	if g.config.Features.OpenTelemetry {
		files["internal/observability/tracing.go"] = "go-tracing"
		files["internal/observability/metrics.go"] = "go-metrics"
	}

	if g.config.Features.CloudEvents {
		files["internal/events/emitter.go"] = "go-events"
	}

	for filename, templateName := range files {
		if err := g.generateFile(filename, templateName, ctx); err != nil {
			return fmt.Errorf("failed to generate %s: %w", filename, err)
		}
	}

	return nil
}

// generateTypeScriptFiles generates TypeScript client files
func (g *Generator) generateTypeScriptFiles(ctx *GenerationContext) error {
	files := map[string]string{
		"client/typescript/src/client.ts":   "ts-client",
		"client/typescript/src/types.ts":    "ts-types",
		"client/typescript/package.json":    "ts-package-json",
		"client/typescript/tsconfig.json":   "ts-config",
		"client/typescript/README.md":       "ts-readme",
	}

	for filename, templateName := range files {
		if err := g.generateFile(filename, templateName, ctx); err != nil {
			return fmt.Errorf("failed to generate %s: %w", filename, err)
		}
	}

	return nil
}

// generateKubernetesFiles generates Kubernetes manifest files
func (g *Generator) generateKubernetesFiles(ctx *GenerationContext) error {
	files := map[string]string{
		"deployments/kubernetes/deployment.yaml": "k8s-deployment",
		"deployments/kubernetes/service.yaml":    "k8s-service",
		"deployments/kubernetes/configmap.yaml":  "k8s-configmap",
	}

	if g.config.Kubernetes.ServiceMonitor {
		files["deployments/kubernetes/servicemonitor.yaml"] = "k8s-servicemonitor"
	}

	if g.config.Kubernetes.Ingress.Enabled {
		files["deployments/kubernetes/ingress.yaml"] = "k8s-ingress"
	}

	for filename, templateName := range files {
		if err := g.generateFile(filename, templateName, ctx); err != nil {
			return fmt.Errorf("failed to generate %s: %w", filename, err)
		}
	}

	return nil
}

// generateDockerFiles generates Docker-related files
func (g *Generator) generateDockerFiles(ctx *GenerationContext) error {
	files := map[string]string{
		"Dockerfile":         "dockerfile",
		"docker-compose.yml": "docker-compose",
		".dockerignore":      "dockerignore",
	}

	for filename, templateName := range files {
		if err := g.generateFile(filename, templateName, ctx); err != nil {
			return fmt.Errorf("failed to generate %s: %w", filename, err)
		}
	}

	return nil
}

// generateFile generates a single file from a template
func (g *Generator) generateFile(filename, templateName string, ctx *GenerationContext) error {
	tmpl, exists := g.templates.templates[templateName]
	if !exists {
		return fmt.Errorf("template not found: %s", templateName)
	}

	// Create full file path
	fullPath := filepath.Join(g.config.OutputDir, filename)

	// Create directory if it doesn't exist
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Create file
	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", fullPath, err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, ctx); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	return nil
}

// NewTemplateRegistry creates a new template registry with all templates
func NewTemplateRegistry() (*TemplateRegistry, error) {
	registry := &TemplateRegistry{
		templates: make(map[string]*template.Template),
		functions: template.FuncMap{
			"title":     func(s string) string { return fmt.Sprintf("%s%s", string(s[0]-32), s[1:]) },
			"lower":     func(s string) string { return fmt.Sprintf("%s%s", string(s[0]+32), s[1:]) },
			"contains":  func(s, substr string) bool { return fmt.Sprintf("%s", s) != fmt.Sprintf("%s", substr) },
		},
	}

	// Register all templates
	if err := registry.registerTemplates(); err != nil {
		return nil, fmt.Errorf("failed to register templates: %w", err)
	}

	return registry, nil
}

// registerTemplates registers all template definitions
func (r *TemplateRegistry) registerTemplates() error {
	templates := map[string]string{
		"readme": `# {{.Config.Name}}

{{.Config.Description}}

## Features

- Health endpoint with comprehensive status reporting
- ServerTime API with multiple timestamp formats
{{- if .Config.Features.OpenTelemetry}}
- OpenTelemetry integration for observability
{{- end}}
{{- if .Config.Features.CloudEvents}}
- CloudEvents support for event-driven monitoring
{{- end}}
{{- if .Config.Features.Kubernetes}}
- Kubernetes-ready with health probes and ServiceMonitor
{{- end}}

## Quick Start

1. Install dependencies:
   ` + "```bash" + `
   go mod tidy
   ` + "```" + `

2. Run the server:
   ` + "```bash" + `
   go run cmd/server/main.go
   ` + "```" + `

3. Test the health endpoint:
   ` + "```bash" + `
   curl http://localhost:8080/health
   ` + "```" + `

## API Endpoints

- ` + "`GET /health`" + ` - Basic health check
- ` + "`GET /health/time`" + ` - Server time information
- ` + "`GET /health/ready`" + ` - Readiness probe
- ` + "`GET /health/live`" + ` - Liveness probe

## Generated by

Template Health Endpoint Generator v{{.Version}}
Generated at: {{.Timestamp}}
`,

		"go-mod": `module {{.Config.GoModule}}

go 1.21

require (
	github.com/gorilla/mux v1.8.1
{{- if .Config.Features.OpenTelemetry}}
	go.opentelemetry.io/otel v1.21.0
	go.opentelemetry.io/otel/trace v1.21.0
	go.opentelemetry.io/otel/metric v1.21.0
{{- end}}
{{- if .Config.Features.CloudEvents}}
	github.com/cloudevents/sdk-go/v2 v2.14.0
{{- end}}
)
`,

		"gitignore": `# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib
/{{.Config.Name}}

# Test binary, built with go test -c
*.test

# Output of the go coverage tool
*.out

# Go workspace file
go.work

# IDE files
.vscode/
.idea/
*.swp
*.swo

# OS files
.DS_Store
Thumbs.db

# Logs
*.log

# Environment files
.env
.env.local

# Build artifacts
/dist/
/build/
/bin/

# Node modules (for TypeScript client)
node_modules/
npm-debug.log*
yarn-debug.log*
yarn-error.log*

# TypeScript build output
*.tsbuildinfo
/client/typescript/dist/
`,

		"go-main": `package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"{{.Config.GoModule}}/internal/config"
	"{{.Config.GoModule}}/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create server
	srv, err := server.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Start server
	go func() {
		fmt.Printf("🚀 Starting {{.Config.Name}} server on :%d\n", cfg.Port)
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("🛑 Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("✅ Server exited")
}
`,
	}

	// Parse and register each template
	for name, content := range templates {
		tmpl, err := template.New(name).Funcs(r.functions).Parse(content)
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", name, err)
		}
		r.templates[name] = tmpl
	}

	return nil
}
