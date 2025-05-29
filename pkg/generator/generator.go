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
