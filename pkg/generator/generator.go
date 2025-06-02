package generator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/LarsArtmann/BMAD-METHOD/pkg/config"
)

// Generator handles the generation of health endpoint projects
type Generator struct {
	config           *config.ProjectConfig
	templates        *TemplateRegistry
	cache            *TemplateCache
	parallelGen      *ParallelGenerator
	enableParallel   bool
	enableCaching    bool
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

	// Create template cache (1 hour max age)
	cache := NewTemplateCache(1 * time.Hour)

	// Create parallel generator with optimal worker count
	parallelGen := NewParallelGenerator(GetOptimalWorkerCount())

	return &Generator{
		config:         cfg,
		templates:      registry,
		cache:          cache,
		parallelGen:    parallelGen,
		enableParallel: true,  // Enable by default
		enableCaching:  true,  // Enable by default
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
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   "1.0.0",
	}

	if g.enableParallel {
		return g.generateParallel(ctx)
	}

	return g.generateSequential(ctx)
}

// generateParallel generates files using parallel processing
func (g *Generator) generateParallel(ctx *GenerationContext) error {
	// Collect all generation tasks
	tasks := g.collectGenerationTasks(ctx)

	// Use context with timeout for parallel generation
	genCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Generate files in parallel
	summary, err := g.parallelGen.GenerateFiles(genCtx, tasks)
	if err != nil {
		return fmt.Errorf("parallel generation failed: %w", err)
	}

	// Check for failures
	if summary.FailureCount > 0 {
		return fmt.Errorf("failed to generate %d out of %d files", summary.FailureCount, summary.TotalFiles)
	}

	fmt.Printf("✅ Generated %d files in %v (parallel mode)\n", summary.SuccessCount, summary.TotalDuration)
	return nil
}

// generateSequential generates files sequentially (fallback)
func (g *Generator) generateSequential(ctx *GenerationContext) error {
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

// collectGenerationTasks collects all files that need to be generated
func (g *Generator) collectGenerationTasks(ctx *GenerationContext) []GenerationTask {
	var tasks []GenerationTask

	// Core files
	coreFiles := map[string]string{
		"README.md":           "readme",
		"go.mod":              "go-mod",
		".gitignore":          "gitignore",
		"Makefile":            "makefile",
		"docs/API.md":         "api-docs",
		"scripts/build.sh":    "build-script",
		"scripts/test.sh":     "test-script",
	}

	for filename, templateName := range coreFiles {
		tasks = append(tasks, GenerationTask{
			Filename:     filename,
			TemplateName: templateName,
			Context:      ctx,
			Generator:    g,
		})
	}

	// Go source files
	goFiles := map[string]string{
		"cmd/server/main.go":           "go-main",
		"internal/handlers/health.go":  "go-health-handler",
		"internal/models/health.go":    "go-health-models",
		"internal/server/server.go":    "go-server",
		"internal/config/config.go":    "go-config",
	}

	// Add tier-specific files
	switch g.config.Tier {
	case config.TierIntermediate, config.TierAdvanced, config.TierEnterprise:
		goFiles["internal/handlers/server_time.go"] = "go-server-time-handler"
		goFiles["internal/handlers/dependencies.go"] = "go-dependencies-handler"
	}

	// Add feature-specific files
	if g.config.Features.OpenTelemetry {
		goFiles["internal/observability/tracing.go"] = "go-tracing"
		goFiles["internal/observability/metrics.go"] = "go-metrics"
	}

	if g.config.Features.Security {
		goFiles["internal/security/mtls.go"] = "go-security-mtls"
		goFiles["internal/security/rbac.go"] = "go-security-rbac"
		goFiles["internal/security/context.go"] = "go-security-context"
	}

	if g.config.Features.Compliance {
		goFiles["internal/compliance/audit.go"] = "go-compliance-audit"
	}

	if g.config.Features.CloudEvents {
		goFiles["internal/events/emitter.go"] = "go-events"
	}

	for filename, templateName := range goFiles {
		tasks = append(tasks, GenerationTask{
			Filename:     filename,
			TemplateName: templateName,
			Context:      ctx,
			Generator:    g,
		})
	}

	// TypeScript files
	if g.config.Features.TypeScript {
		tsFiles := map[string]string{
			"client/typescript/src/client.ts":   "ts-client",
			"client/typescript/src/types.ts":    "ts-types",
			"client/typescript/package.json":    "ts-package-json",
			"client/typescript/tsconfig.json":   "ts-config",
			"client/typescript/README.md":       "ts-readme",
		}

		for filename, templateName := range tsFiles {
			tasks = append(tasks, GenerationTask{
				Filename:     filename,
				TemplateName: templateName,
				Context:      ctx,
				Generator:    g,
			})
		}
	}

	// Kubernetes files
	if g.config.Features.Kubernetes {
		k8sFiles := map[string]string{
			"deployments/kubernetes/deployment.yaml": "k8s-deployment",
			"deployments/kubernetes/service.yaml":    "k8s-service",
			"deployments/kubernetes/configmap.yaml":  "k8s-configmap",
		}

		if g.config.Kubernetes.ServiceMonitor {
			k8sFiles["deployments/kubernetes/servicemonitor.yaml"] = "k8s-servicemonitor"
		}

		if g.config.Kubernetes.Ingress.Enabled {
			k8sFiles["deployments/kubernetes/ingress.yaml"] = "k8s-ingress"
		}

		for filename, templateName := range k8sFiles {
			tasks = append(tasks, GenerationTask{
				Filename:     filename,
				TemplateName: templateName,
				Context:      ctx,
				Generator:    g,
			})
		}
	}

	// Docker files
	if g.config.Features.Docker {
		dockerFiles := map[string]string{
			"Dockerfile":         "dockerfile",
			"docker-compose.yml": "docker-compose",
			".dockerignore":      "dockerignore",
		}

		for filename, templateName := range dockerFiles {
			tasks = append(tasks, GenerationTask{
				Filename:     filename,
				TemplateName: templateName,
				Context:      ctx,
				Generator:    g,
			})
		}
	}

	return tasks
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

	if g.config.Features.OpenTelemetry {
		dirs = append(dirs,
			filepath.Join(g.config.OutputDir, "internal", "observability"),
		)
	}

	if g.config.Features.CloudEvents {
		dirs = append(dirs,
			filepath.Join(g.config.OutputDir, "internal", "events"),
		)
	}

	if g.config.Features.Security {
		dirs = append(dirs,
			filepath.Join(g.config.OutputDir, "internal", "security"),
		)
	}

	if g.config.Features.Compliance {
		dirs = append(dirs,
			filepath.Join(g.config.OutputDir, "internal", "compliance"),
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

	// Add security files for enterprise tier
	if g.config.Features.Security {
		files["internal/security/mtls.go"] = "go-security-mtls"
		files["internal/security/rbac.go"] = "go-security-rbac"
		files["internal/security/context.go"] = "go-security-context"
	}

	// Add compliance files for enterprise tier
	if g.config.Features.Compliance {
		files["internal/compliance/audit.go"] = "go-compliance-audit"
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
- ` + "`GET /health/startup`" + ` - Startup probe

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

		"go-config": `package config

import (
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	Port    int    ` + "`json:\"port\" yaml:\"port\"`" + `
	Version string ` + "`json:\"version\" yaml:\"version\"`" + `
	Name    string ` + "`json:\"name\" yaml:\"name\"`" + `
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Port:    8080,
		Version: "{{.Config.Version}}",
		Name:    "{{.Config.Name}}",
	}

	// Override with environment variables
	if port := os.Getenv("PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Port = p
		}
	}

	if version := os.Getenv("VERSION"); version != "" {
		cfg.Version = version
	}

	return cfg, nil
}
`,

		"go-server": `package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"{{.Config.GoModule}}/internal/config"
	"{{.Config.GoModule}}/internal/handlers"
)

// Server represents the HTTP server
type Server struct {
	config  *config.Config
	server  *http.Server
	handler *handlers.HealthHandler
}

// New creates a new server instance
func New(cfg *config.Config) (*Server, error) {
	// Create health handler
	healthHandler := handlers.NewHealthHandler(cfg)

	// Create router
	router := mux.NewRouter()

	// Health endpoints
	health := router.PathPrefix("/health").Subrouter()
	health.HandleFunc("", healthHandler.CheckHealth).Methods("GET")
	health.HandleFunc("/", healthHandler.CheckHealth).Methods("GET")
	health.HandleFunc("/time", healthHandler.ServerTime).Methods("GET")
	health.HandleFunc("/ready", healthHandler.ReadinessCheck).Methods("GET")
	health.HandleFunc("/live", healthHandler.LivenessCheck).Methods("GET")
	health.HandleFunc("/startup", healthHandler.StartupCheck).Methods("GET")

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		config:  cfg,
		server:  srv,
		handler: healthHandler,
	}, nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
`,

		"go-health-handler": `package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"{{.Config.GoModule}}/internal/config"
	"{{.Config.GoModule}}/internal/models"
)

// HealthHandler handles health-related HTTP requests
type HealthHandler struct {
	config    *config.Config
	startTime time.Time
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(cfg *config.Config) *HealthHandler {
	return &HealthHandler{
		config:    cfg,
		startTime: time.Now(),
	}
}

// CheckHealth handles GET /health requests
func (h *HealthHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	h.setJSONContentType(w)

	uptime := time.Since(h.startTime)
	status := models.HealthReport{
		Status:      "healthy",
		Timestamp:   time.Now(),
		Version:     h.config.Version,
		Uptime:      uptime,
		UptimeHuman: h.formatUptime(uptime),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}

// ServerTime handles GET /health/time requests
func (h *HealthHandler) ServerTime(w http.ResponseWriter, r *http.Request) {
	h.setJSONContentType(w)

	now := time.Now()
	location := now.Location()

	serverTime := models.ServerTime{
		Timestamp:   now,
		Timezone:    location.String(),
		Unix:        now.Unix(),
		UnixMilli:   now.UnixMilli(),
		ISO8601:     now.Format(time.RFC3339),
		Formatted:   now.Format("Monday, January 2, 2006 at 3:04:05 PM MST"),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(serverTime)
}

// ReadinessCheck handles GET /health/ready requests
func (h *HealthHandler) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	h.setJSONContentType(w)

	// For basic tier, readiness is same as health
	uptime := time.Since(h.startTime)
	status := models.HealthReport{
		Status:      "healthy",
		Timestamp:   time.Now(),
		Version:     h.config.Version,
		Uptime:      uptime,
		UptimeHuman: h.formatUptime(uptime),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}

// LivenessCheck handles GET /health/live requests
func (h *HealthHandler) LivenessCheck(w http.ResponseWriter, r *http.Request) {
	h.setJSONContentType(w)

	// For basic tier, liveness is same as health
	uptime := time.Since(h.startTime)
	status := models.HealthReport{
		Status:      "healthy",
		Timestamp:   time.Now(),
		Version:     h.config.Version,
		Uptime:      uptime,
		UptimeHuman: h.formatUptime(uptime),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}

// StartupCheck handles GET /health/startup requests
func (h *HealthHandler) StartupCheck(w http.ResponseWriter, r *http.Request) {
	h.setJSONContentType(w)

	// For basic tier, startup is same as health
	uptime := time.Since(h.startTime)
	status := models.HealthReport{
		Status:      "healthy",
		Timestamp:   time.Now(),
		Version:     h.config.Version,
		Uptime:      uptime,
		UptimeHuman: h.formatUptime(uptime),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}

// setJSONContentType sets the JSON content type header
func (h *HealthHandler) setJSONContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

// formatUptime formats a duration into human-readable format
func (h *HealthHandler) formatUptime(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.1f seconds", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.1f minutes", d.Minutes())
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%.1f hours", d.Hours())
	}
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	return fmt.Sprintf("%d days, %d hours", days, hours)
}
`,

		"go-health-models": `package models

import "time"

// HealthReport represents the overall health status of the service
type HealthReport struct {
	Status      string        ` + "`json:\"status\"`" + `
	Timestamp   time.Time     ` + "`json:\"timestamp\"`" + `
	Version     string        ` + "`json:\"version\"`" + `
	Uptime      time.Duration ` + "`json:\"uptime\"`" + `
	UptimeHuman string        ` + "`json:\"uptime_human\"`" + `
}

// ServerTime represents server time information with multiple formats
type ServerTime struct {
	Timestamp time.Time ` + "`json:\"timestamp\"`" + `
	Timezone  string    ` + "`json:\"timezone\"`" + `
	Unix      int64     ` + "`json:\"unix\"`" + `
	UnixMilli int64     ` + "`json:\"unix_milli\"`" + `
	ISO8601   string    ` + "`json:\"iso8601\"`" + `
	Formatted string    ` + "`json:\"formatted\"`" + `
}
`,

		"go-server-time-handler": `package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

// ServerTimeHandler handles server time requests
type ServerTimeHandler struct{}

// NewServerTimeHandler creates a new server time handler
func NewServerTimeHandler() *ServerTimeHandler {
	return &ServerTimeHandler{}
}

// GetServerTime handles GET /health/time requests
func (h *ServerTimeHandler) GetServerTime(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	response := map[string]interface{}{
		"timestamp":    now,
		"unix":         now.Unix(),
		"unix_milli":   now.UnixMilli(),
		"rfc3339":      now.Format(time.RFC3339),
		"timezone":     now.Location().String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
`,

		"go-tracing": `package observability

// Simplified tracing for intermediate tier
type TracingProvider struct{}

func NewTracingProvider(serviceName string) *TracingProvider {
	return &TracingProvider{}
}
`,

		"go-metrics": `package observability

// Simplified metrics for intermediate tier
type MetricsProvider struct{}

func NewMetricsProvider(serviceName string) *MetricsProvider {
	return &MetricsProvider{}
}
`,

		"go-events": `package events

// Simplified events for intermediate tier
type EventEmitter struct{}

func NewEventEmitter(serviceName, sinkURL string) *EventEmitter {
	return &EventEmitter{}
}
`,

		"go-security-mtls": `package security

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// MTLSConfig holds the configuration for mutual TLS
type MTLSConfig struct {
	CertFile   string
	KeyFile    string
	CAFile     string
	ClientAuth tls.ClientAuthType
}

// SetupMTLS configures mutual TLS for the server
func SetupMTLS(config MTLSConfig) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load server certificate: %w", err)
	}

	caCert, err := ioutil.ReadFile(config.CAFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %w", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to parse CA certificate")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   config.ClientAuth,
		ClientCAs:    caCertPool,
		MinVersion:   tls.VersionTLS12,
	}

	return tlsConfig, nil
}

// MTLSMiddleware validates client certificates
func MTLSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
			http.Error(w, "Client certificate required", http.StatusUnauthorized)
			return
		}

		clientCert := r.TLS.PeerCertificates[0]
		clientID := clientCert.Subject.CommonName
		if clientID == "" {
			http.Error(w, "Invalid client certificate", http.StatusUnauthorized)
			return
		}

		ctx := WithClientIdentity(r.Context(), clientID)
		r = r.WithContext(ctx)

		log.Printf("mTLS: Client authenticated: %s", clientID)
		next.ServeHTTP(w, r)
	})
}
`,

		"go-security-rbac": `package security

import (
	"net/http"
	"strings"
)

// Permission represents a specific permission
type Permission string

const (
	PermissionHealthRead     Permission = "health:read"
	PermissionHealthWrite    Permission = "health:write"
	PermissionMetricsRead    Permission = "metrics:read"
	PermissionDependencyRead Permission = "dependency:read"
	PermissionAdminAccess    Permission = "admin:access"
)

// Role represents a user role with associated permissions
type Role struct {
	Name        string       ` + "`json:\"name\"`" + `
	Permissions []Permission ` + "`json:\"permissions\"`" + `
}

// User represents a user with roles
type User struct {
	ID    string ` + "`json:\"id\"`" + `
	Roles []Role ` + "`json:\"roles\"`" + `
}

// RBACPolicy holds the role-based access control policy
type RBACPolicy struct {
	Users map[string]User ` + "`json:\"users\"`" + `
	Roles map[string]Role ` + "`json:\"roles\"`" + `
}

// DefaultRBACPolicy returns a default RBAC policy
func DefaultRBACPolicy() *RBACPolicy {
	return &RBACPolicy{
		Users: map[string]User{
			"admin": {
				ID: "admin",
				Roles: []Role{
					{Name: "admin", Permissions: []Permission{
						PermissionHealthRead,
						PermissionHealthWrite,
						PermissionMetricsRead,
						PermissionDependencyRead,
						PermissionAdminAccess,
					}},
				},
			},
			"service": {
				ID: "service",
				Roles: []Role{
					{Name: "service", Permissions: []Permission{
						PermissionHealthRead,
					}},
				},
			},
		},
		Roles: map[string]Role{
			"admin": {
				Name: "admin",
				Permissions: []Permission{
					PermissionHealthRead,
					PermissionHealthWrite,
					PermissionMetricsRead,
					PermissionDependencyRead,
					PermissionAdminAccess,
				},
			},
			"service": {
				Name: "service",
				Permissions: []Permission{
					PermissionHealthRead,
				},
			},
		},
	}
}

// RBACMiddleware validates user permissions for requests
func RBACMiddleware(policy *RBACPolicy) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientID := GetClientIdentity(r.Context())
			if clientID == "" {
				http.Error(w, "Client identity required", http.StatusUnauthorized)
				return
			}

			requiredPermission := getRequiredPermission(r)
			if requiredPermission == "" {
				next.ServeHTTP(w, r)
				return
			}

			if !policy.HasPermission(clientID, requiredPermission) {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// HasPermission checks if a user has a specific permission
func (p *RBACPolicy) HasPermission(userID string, permission Permission) bool {
	user, exists := p.Users[userID]
	if !exists {
		return false
	}

	for _, role := range user.Roles {
		for _, perm := range role.Permissions {
			if perm == permission {
				return true
			}
		}
	}

	return false
}

// getRequiredPermission determines the required permission based on the request
func getRequiredPermission(r *http.Request) Permission {
	path := strings.TrimPrefix(r.URL.Path, "/")
	method := r.Method

	switch {
	case strings.HasPrefix(path, "health"):
		if method == "GET" {
			return PermissionHealthRead
		}
		return PermissionHealthWrite
	case strings.HasPrefix(path, "metrics"):
		return PermissionMetricsRead
	case strings.HasPrefix(path, "dependencies"):
		return PermissionDependencyRead
	case strings.HasPrefix(path, "admin"):
		return PermissionAdminAccess
	default:
		return ""
	}
}
`,

		"go-security-context": `package security

import "context"

type contextKey string

const (
	clientIdentityKey contextKey = "client_identity"
	auditContextKey   contextKey = "audit_context"
)

// WithClientIdentity adds client identity to context
func WithClientIdentity(ctx context.Context, clientID string) context.Context {
	return context.WithValue(ctx, clientIdentityKey, clientID)
}

// GetClientIdentity retrieves client identity from context
func GetClientIdentity(ctx context.Context) string {
	if clientID, ok := ctx.Value(clientIdentityKey).(string); ok {
		return clientID
	}
	return ""
}

// AuditContext holds audit information
type AuditContext struct {
	UserID    string
	Action    string
	Resource  string
	Timestamp int64
	RequestID string
}

// WithAuditContext adds audit context
func WithAuditContext(ctx context.Context, auditCtx *AuditContext) context.Context {
	return context.WithValue(ctx, auditContextKey, auditCtx)
}

// GetAuditContext retrieves audit context
func GetAuditContext(ctx context.Context) *AuditContext {
	if auditCtx, ok := ctx.Value(auditContextKey).(*AuditContext); ok {
		return auditCtx
	}
	return nil
}
`,

		"go-compliance-audit": `package compliance

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"{{.Config.GoModule}}/internal/security"
)

// AuditEvent represents an audit log event
type AuditEvent struct {
	Timestamp   time.Time ` + "`json:\"timestamp\"`" + `
	EventID     string    ` + "`json:\"event_id\"`" + `
	UserID      string    ` + "`json:\"user_id\"`" + `
	Action      string    ` + "`json:\"action\"`" + `
	Resource    string    ` + "`json:\"resource\"`" + `
	Method      string    ` + "`json:\"method\"`" + `
	Path        string    ` + "`json:\"path\"`" + `
	StatusCode  int       ` + "`json:\"status_code\"`" + `
	Duration    int64     ` + "`json:\"duration_ms\"`" + `
	RequestID   string    ` + "`json:\"request_id\"`" + `
	ClientIP    string    ` + "`json:\"client_ip\"`" + `
	UserAgent   string    ` + "`json:\"user_agent\"`" + `
	Success     bool      ` + "`json:\"success\"`" + `
	ErrorMsg    string    ` + "`json:\"error_message,omitempty\"`" + `
}

// AuditLogger handles audit logging
type AuditLogger struct {
	logger   *log.Logger
	file     *os.File
	enabled  bool
}

// NewAuditLogger creates a new audit logger
func NewAuditLogger(logFile string, enabled bool) (*AuditLogger, error) {
	if !enabled {
		return &AuditLogger{enabled: false}, nil
	}

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open audit log file: %w", err)
	}

	logger := log.New(file, "", 0)

	return &AuditLogger{
		logger:   logger,
		file:     file,
		enabled:  true,
	}, nil
}

// LogEvent logs an audit event
func (a *AuditLogger) LogEvent(event AuditEvent) error {
	if !a.enabled {
		return nil
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal audit event: %w", err)
	}

	a.logger.Println(string(eventJSON))
	return nil
}

// LogHTTPRequest logs an HTTP request audit event
func (a *AuditLogger) LogHTTPRequest(r *http.Request, statusCode int, duration time.Duration, err error) {
	if !a.enabled {
		return
	}

	userID := security.GetClientIdentity(r.Context())
	if userID == "" {
		userID = "anonymous"
	}

	event := AuditEvent{
		Timestamp:  time.Now().UTC(),
		EventID:    fmt.Sprintf("audit_%d", time.Now().UnixNano()),
		UserID:     userID,
		Action:     "http_request",
		Resource:   r.URL.Path,
		Method:     r.Method,
		Path:       r.URL.Path,
		StatusCode: statusCode,
		Duration:   duration.Milliseconds(),
		RequestID:  r.Header.Get("X-Request-ID"),
		ClientIP:   r.RemoteAddr,
		UserAgent:  r.UserAgent(),
		Success:    statusCode < 400,
	}

	if err != nil {
		event.ErrorMsg = err.Error()
	}

	if logErr := a.LogEvent(event); logErr != nil {
		log.Printf("Failed to log audit event: %v", logErr)
	}
}

// AuditMiddleware creates middleware for HTTP request auditing
func AuditMiddleware(auditLogger *AuditLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapped := &responseWriter{ResponseWriter: w, statusCode: 200}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)
			auditLogger.LogHTTPRequest(r, wrapped.statusCode, duration, nil)
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Close closes the audit logger
func (a *AuditLogger) Close() error {
	if a.file != nil {
		return a.file.Close()
	}
	return nil
}
`,

		"go-dependencies-handler": `package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

// DependenciesHandler handles dependency health checks
type DependenciesHandler struct{}

// NewDependenciesHandler creates a new dependencies handler
func NewDependenciesHandler() *DependenciesHandler {
	return &DependenciesHandler{}
}

// DependencyStatus represents the status of a dependency
type DependencyStatus struct {
	Name      string        ` + "`json:\"name\"`" + `
	Status    string        ` + "`json:\"status\"`" + `
	Latency   time.Duration ` + "`json:\"latency\"`" + `
	Error     string        ` + "`json:\"error,omitempty\"`" + `
	Timestamp time.Time     ` + "`json:\"timestamp\"`" + `
}

// DependenciesResponse represents the dependencies health response
type DependenciesResponse struct {
	Status       string              ` + "`json:\"status\"`" + `
	Dependencies []DependencyStatus  ` + "`json:\"dependencies\"`" + `
	Timestamp    time.Time           ` + "`json:\"timestamp\"`" + `
}

// CheckDependencies handles GET /health/dependencies requests
func (h *DependenciesHandler) CheckDependencies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Simulate dependency checks
	dependencies := []DependencyStatus{
		{
			Name:      "database",
			Status:    "healthy",
			Latency:   5 * time.Millisecond,
			Timestamp: time.Now(),
		},
		{
			Name:      "cache",
			Status:    "healthy",
			Latency:   2 * time.Millisecond,
			Timestamp: time.Now(),
		},
	}

	// Determine overall status
	overallStatus := "healthy"
	for _, dep := range dependencies {
		if dep.Status != "healthy" {
			overallStatus = "degraded"
			break
		}
	}

	response := DependenciesResponse{
		Status:       overallStatus,
		Dependencies: dependencies,
		Timestamp:    time.Now(),
	}

	if overallStatus == "healthy" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
`,

		"ts-types": `// Generated TypeScript types for {{.Config.Name}}
// Generated at: {{.Timestamp}}

export interface HealthReport {
  status: string;
  timestamp: string;
  version: string;
  uptime: number;
  uptime_human: string;
}

export interface ServerTime {
  timestamp: string;
  timezone: string;
  unix: number;
  unix_milli: number;
  iso8601: string;
  formatted: string;
}

export type HealthStatus = 'healthy' | 'degraded' | 'unhealthy';
`,

		"ts-client": `// Generated TypeScript client for {{.Config.Name}}
// Generated at: {{.Timestamp}}

import { HealthReport, ServerTime } from './types';

export interface HealthClientConfig {
  baseURL: string;
  timeout?: number;
  headers?: Record<string, string>;
}

export class HealthClient {
  private baseURL: string;
  private timeout: number;
  private headers: Record<string, string>;

  constructor(config: HealthClientConfig) {
    this.baseURL = config.baseURL.replace(/\/$/, '');
    this.timeout = config.timeout || 5000;
    this.headers = config.headers || {};
  }

  /**
   * Check the health status of the service
   */
  async checkHealth(): Promise<HealthReport> {
    return this.request<HealthReport>('/health');
  }

  /**
   * Get server time information
   */
  async getServerTime(): Promise<ServerTime> {
    return this.request<ServerTime>('/health/time');
  }

  /**
   * Check readiness status
   */
  async checkReadiness(): Promise<HealthReport> {
    return this.request<HealthReport>('/health/ready');
  }

  /**
   * Check liveness status
   */
  async checkLiveness(): Promise<HealthReport> {
    return this.request<HealthReport>('/health/live');
  }

  /**
   * Check startup status
   */
  async checkStartup(): Promise<HealthReport> {
    return this.request<HealthReport>('/health/startup');
  }

  private async request<T>(path: string): Promise<T> {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), this.timeout);

    try {
      const response = await fetch(` + "`${this.baseURL}${path}`" + `, {
        method: 'GET',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          ...this.headers,
        },
        signal: controller.signal,
      });

      clearTimeout(timeoutId);

      if (!response.ok) {
        throw new Error(` + "`HTTP ${response.status}: ${response.statusText}`" + `);
      }

      return await response.json();
    } catch (error) {
      clearTimeout(timeoutId);
      if (error instanceof Error && error.name === 'AbortError') {
        throw new Error(` + "`Request timeout after ${this.timeout}ms`" + `);
      }
      throw error;
    }
  }
}

// Default export for convenience
export default HealthClient;
`,

		"ts-package-json": `{
  "name": "{{.Config.Name}}-client",
  "version": "{{.Config.Version}}",
  "description": "TypeScript client for {{.Config.Name}} health endpoints",
  "main": "dist/index.js",
  "types": "dist/index.d.ts",
  "scripts": {
    "build": "tsc",
    "build:watch": "tsc --watch",
    "clean": "rm -rf dist",
    "prepublishOnly": "npm run clean && npm run build"
  },
  "files": [
    "dist/**/*",
    "src/**/*"
  ],
  "keywords": [
    "health-check",
    "monitoring",
    "typescript",
    "client"
  ],
  "author": "Generated by template-health-endpoint",
  "license": "MIT",
  "devDependencies": {
    "typescript": "^5.0.0",
    "@types/node": "^20.0.0"
  },
  "engines": {
    "node": ">=16.0.0"
  }
}
`,

		"ts-config": `{
  "compilerOptions": {
    "target": "ES2020",
    "module": "commonjs",
    "lib": ["ES2020", "DOM"],
    "outDir": "./dist",
    "rootDir": "./src",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "declaration": true,
    "declarationMap": true,
    "sourceMap": true,
    "removeComments": false,
    "noImplicitAny": true,
    "strictNullChecks": true,
    "strictFunctionTypes": true,
    "noImplicitThis": true,
    "noImplicitReturns": true,
    "noFallthroughCasesInSwitch": true,
    "moduleResolution": "node",
    "allowSyntheticDefaultImports": true,
    "experimentalDecorators": true,
    "emitDecoratorMetadata": true
  },
  "include": [
    "src/**/*"
  ],
  "exclude": [
    "node_modules",
    "dist"
  ]
}
`,

		"ts-readme": `# {{.Config.Name}} TypeScript Client

TypeScript client library for {{.Config.Name}} health endpoints.

## Installation

` + "```bash" + `
npm install {{.Config.Name}}-client
` + "```" + `

## Usage

` + "```typescript" + `
import { HealthClient } from '{{.Config.Name}}-client';

const client = new HealthClient({
  baseURL: 'http://localhost:8080',
  timeout: 5000,
});

// Check health status
const health = await client.checkHealth();
console.log('Health status:', health.status);

// Get server time
const serverTime = await client.getServerTime();
console.log('Server time:', serverTime.formatted);

// Check readiness
const readiness = await client.checkReadiness();
console.log('Readiness:', readiness.status);

// Check liveness
const liveness = await client.checkLiveness();
console.log('Liveness:', liveness.status);
` + "```" + `

## API

### HealthClient

#### Constructor

` + "```typescript" + `
new HealthClient(config: HealthClientConfig)
` + "```" + `

- ` + "`config.baseURL`" + ` - Base URL of the health service
- ` + "`config.timeout`" + ` - Request timeout in milliseconds (default: 5000)
- ` + "`config.headers`" + ` - Additional headers to send with requests

#### Methods

- ` + "`checkHealth(): Promise<HealthReport>`" + ` - Get overall health status
- ` + "`getServerTime(): Promise<ServerTime>`" + ` - Get server time information
- ` + "`checkReadiness(): Promise<HealthReport>`" + ` - Check if service is ready
- ` + "`checkLiveness(): Promise<HealthReport>`" + ` - Check if service is alive
- ` + "`checkStartup(): Promise<HealthReport>`" + ` - Check if service has started up

## Types

See ` + "`src/types.ts`" + ` for complete type definitions.

## Generated by

Template Health Endpoint Generator v{{.Version}}
Generated at: {{.Timestamp}}
`,

		"dockerfile": `# Multi-stage build for {{.Config.Name}}
FROM golang:1.21-alpine AS builder

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN adduser -D -s /bin/sh appuser

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Change ownership to appuser
RUN chown appuser:appuser main

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \\
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"]
`,

		"makefile": `# Makefile for {{.Config.Name}}

.PHONY: build run test clean docker-build docker-run help

# Variables
APP_NAME={{.Config.Name}}
VERSION={{.Config.Version}}
GO_VERSION=1.21
DOCKER_IMAGE=$(APP_NAME):$(VERSION)

# Default target
all: build

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	go build -o bin/$(APP_NAME) cmd/server/main.go

# Run the application
run: build
	@echo "Running $(APP_NAME)..."
	./bin/$(APP_NAME)

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Build Docker image
docker-build:
	@echo "Building Docker image $(DOCKER_IMAGE)..."
	docker build -t $(DOCKER_IMAGE) .

# Run Docker container
docker-run: docker-build
	@echo "Running Docker container..."
	docker run -p 8080:8080 --rm $(DOCKER_IMAGE)

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  test          - Run tests"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Install dependencies"
	@echo "  fmt           - Format code"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  help          - Show this help"
`,

		"api-docs": `# {{.Config.Name}} API Documentation

This document describes the health endpoints provided by {{.Config.Name}}.

## Base URL

` + "```" + `
http://localhost:8080
` + "```" + `

## Endpoints

### GET /health

Returns the overall health status of the service.

**Response:**
` + "```json" + `
{
  "status": "healthy",
  "timestamp": "2024-01-01T12:00:00Z",
  "version": "{{.Config.Version}}",
  "uptime": 3600000000000,
  "uptime_human": "1.0 hours"
}
` + "```" + `

### GET /health/time

Returns server time information in multiple formats.

**Response:**
` + "```json" + `
{
  "timestamp": "2024-01-01T12:00:00Z",
  "timezone": "UTC",
  "unix": 1704110400,
  "unix_milli": 1704110400000,
  "iso8601": "2024-01-01T12:00:00Z",
  "formatted": "Monday, January 1, 2024 at 12:00:00 PM UTC"
}
` + "```" + `

### GET /health/ready

Kubernetes readiness probe endpoint.

**Response:** Same as /health

### GET /health/live

Kubernetes liveness probe endpoint.

**Response:** Same as /health

### GET /health/startup

Kubernetes startup probe endpoint.

**Response:** Same as /health

## Status Codes

- ` + "`200 OK`" + ` - Service is healthy
- ` + "`503 Service Unavailable`" + ` - Service is unhealthy

## Generated by

Template Health Endpoint Generator v{{.Version}}
Generated at: {{.Timestamp}}
`,

		"build-script": `#!/bin/bash

# Build script for {{.Config.Name}}

set -e

echo "🔨 Building {{.Config.Name}}..."

# Clean previous builds
rm -rf bin/
mkdir -p bin/

# Build the application
go build -o bin/{{.Config.Name}} cmd/server/main.go

echo "✅ Build complete: bin/{{.Config.Name}}"
`,

		"test-script": `#!/bin/bash

# Test script for {{.Config.Name}}

set -e

echo "🧪 Running tests for {{.Config.Name}}..."

# Run tests
go test -v ./...

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

echo "✅ Tests complete. Coverage report: coverage.html"
`,

		"k8s-deployment": `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Config.Name}}
  labels:
    app: {{.Config.Name}}
    version: {{.Config.Version}}
spec:
  replicas: 3
  selector:
    matchLabels:
      app: {{.Config.Name}}
  template:
    metadata:
      labels:
        app: {{.Config.Name}}
        version: {{.Config.Version}}
    spec:
      containers:
      - name: {{.Config.Name}}
        image: {{.Config.Name}}:{{.Config.Version}}
        ports:
        - containerPort: 8080
          name: http
        env:
        - name: PORT
          value: "8080"
        - name: VERSION
          value: {{.Config.Version}}
        resources:
          requests:
            memory: "64Mi"
            cpu: "50m"
          limits:
            memory: "128Mi"
            cpu: "100m"
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 3
        startupProbe:
          httpGet:
            path: /health/startup
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 30
      restartPolicy: Always
`,

		"k8s-service": `apiVersion: v1
kind: Service
metadata:
  name: {{.Config.Name}}
  labels:
    app: {{.Config.Name}}
spec:
  selector:
    app: {{.Config.Name}}
  ports:
  - name: http
    port: 80
    targetPort: 8080
    protocol: TCP
  type: ClusterIP
`,

		"k8s-configmap": `apiVersion: v1
kind: ConfigMap
metadata:
  name: {{.Config.Name}}-config
  labels:
    app: {{.Config.Name}}
data:
  PORT: "8080"
  VERSION: "{{.Config.Version}}"
  SERVICE_NAME: "{{.Config.Name}}"
`,

		"k8s-servicemonitor": `apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: {{.Config.Name}}
  labels:
    app: {{.Config.Name}}
spec:
  selector:
    matchLabels:
      app: {{.Config.Name}}
  endpoints:
  - port: http
    path: /metrics
    interval: 30s
    scrapeTimeout: 10s
`,

		"k8s-ingress": `apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {{.Config.Name}}
  labels:
    app: {{.Config.Name}}
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
  - host: {{.Config.Name}}.local
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: {{.Config.Name}}
            port:
              number: 80
`,

		"docker-compose": `version: '3.8'

services:
  {{.Config.Name}}:
    build: .
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - VERSION={{.Config.Version}}
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 5s
    restart: unless-stopped
    networks:
      - health-network

networks:
  health-network:
    driver: bridge
`,

		"dockerignore": `# Git
.git
.gitignore

# Documentation
*.md
README*

# Build artifacts
bin/
dist/
build/

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
.env*

# Node modules
node_modules/

# Test files
*_test.go
*.test

# Coverage
*.out
coverage/
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
