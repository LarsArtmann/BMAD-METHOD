package features

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/LarsArtmann/BMAD-METHOD/pkg/config"
)

// BaseFeatureGenerator provides common functionality for feature generators
type BaseFeatureGenerator struct {
	ID          string
	Name        string
	Description string
	Templates   []string
	Assets      []string
}

// Validate provides basic validation for all features
func (g *BaseFeatureGenerator) Validate(config *config.ProjectConfig, featureConfig map[string]interface{}) error {
	// Basic validation - can be overridden by specific implementations
	return nil
}

// GetTemplates returns the templates required by this feature
func (g *BaseFeatureGenerator) GetTemplates() []string {
	return g.Templates
}

// GetAssets returns the assets required by this feature
func (g *BaseFeatureGenerator) GetAssets() []string {
	return g.Assets
}

// HealthFeatureGenerator generates health check components
type HealthFeatureGenerator struct {
	BaseFeatureGenerator
	Tier string
}

// Generate creates health check implementation for the specified tier
func (g *HealthFeatureGenerator) Generate(ctx context.Context, config *config.ProjectConfig, featureConfig map[string]interface{}) (*GeneratedFeature, error) {
	result := &GeneratedFeature{
		Files:       make(map[string]string),
		Templates:   make([]string, 0),
		Assets:      make([]string, 0),
		Metadata:    make(map[string]interface{}),
		PostActions: make([]PostGenerationAction, 0),
	}

	// Add tier-specific templates
	switch g.Tier {
	case "basic":
		result.Templates = append(result.Templates, "template-health/schemas/health.tsp")
		result.Files["internal/handlers/health.go"] = "// Basic health handler implementation"
		
	case "intermediate":
		result.Templates = append(result.Templates, 
			"template-health/schemas/health.tsp",
			"template-health/schemas/health-api.tsp",
		)
		result.Files["internal/handlers/health.go"] = "// Intermediate health handler with dependencies"
		
	case "advanced":
		result.Templates = append(result.Templates,
			"template-health/schemas/health.tsp",
			"template-health/schemas/health-api.tsp",
			"template-health/templates/go-metrics-advanced.tmpl",
		)
		result.Files["internal/handlers/health.go"] = "// Advanced health handler with metrics"
		result.Files["internal/observability/metrics.go"] = "// Advanced metrics implementation"
		
	case "enterprise":
		result.Templates = append(result.Templates,
			"template-health/schemas/health.tsp",
			"template-health/schemas/health-api.tsp",
			"template-health/templates/go-metrics-advanced.tmpl",
			"template-health/templates/go-tracing-advanced.tmpl",
			"template-health/templates/sre-prometheus-alerts.tmpl",
		)
		result.Files["internal/handlers/health.go"] = "// Enterprise health handler with full observability"
		result.Files["internal/observability/metrics.go"] = "// Enterprise metrics implementation"
		result.Files["internal/observability/tracing.go"] = "// Enterprise tracing implementation"
		result.Files["deployments/monitoring/alerts.yml"] = "// Prometheus alerting rules"
	}

	// Add metadata
	result.Metadata["health_tier"] = g.Tier
	result.Metadata["endpoints"] = []string{"/health", "/health/ready", "/health/live"}

	return result, nil
}

// ObservabilityFeatureGenerator generates observability components
type ObservabilityFeatureGenerator struct {
	BaseFeatureGenerator
	Components []string // metrics, tracing, logging
}

// Generate creates observability implementation
func (g *ObservabilityFeatureGenerator) Generate(ctx context.Context, config *config.ProjectConfig, featureConfig map[string]interface{}) (*GeneratedFeature, error) {
	result := &GeneratedFeature{
		Files:       make(map[string]string),
		Templates:   make([]string, 0),
		Assets:      make([]string, 0),
		Metadata:    make(map[string]interface{}),
		PostActions: make([]PostGenerationAction, 0),
	}

	for _, component := range g.Components {
		switch component {
		case "metrics":
			result.Templates = append(result.Templates, "template-health/templates/go-metrics-advanced.tmpl")
			result.Files["internal/observability/metrics.go"] = "// OpenTelemetry metrics implementation"
			
		case "tracing":
			result.Templates = append(result.Templates, "template-health/templates/go-tracing-advanced.tmpl")
			result.Files["internal/observability/tracing.go"] = "// OpenTelemetry tracing implementation"
			
		case "logging":
			result.Templates = append(result.Templates, "template-health/templates/go-structured-logging.tmpl")
			result.Files["internal/observability/logging.go"] = "// Structured logging implementation"
			
		case "sre":
			result.Templates = append(result.Templates, 
				"template-health/templates/sre-prometheus-alerts.tmpl",
				"template-health/templates/sre-grafana-dashboard.tmpl",
				"template-health/templates/sre-sli-slo-config.tmpl",
			)
			result.Files["deployments/monitoring/alerts.yml"] = "// SRE alerting rules"
			result.Files["deployments/monitoring/dashboard.json"] = "// Grafana dashboard"
			result.Files["deployments/monitoring/slo-config.yml"] = "// SLI/SLO configuration"
		}
	}

	// Add post-generation actions
	result.PostActions = append(result.PostActions, PostGenerationAction{
		Type:        "install",
		Description: "Install OpenTelemetry dependencies",
		Command:     "go",
		Args:        []string{"get", "go.opentelemetry.io/otel"},
		WorkingDir:  ".",
	})

	result.Metadata["observability_components"] = g.Components
	return result, nil
}

// SecurityFeatureGenerator generates security components
type SecurityFeatureGenerator struct {
	BaseFeatureGenerator
	SecurityLevel string // basic, rbac, enterprise
}

// Generate creates security implementation
func (g *SecurityFeatureGenerator) Generate(ctx context.Context, config *config.ProjectConfig, featureConfig map[string]interface{}) (*GeneratedFeature, error) {
	result := &GeneratedFeature{
		Files:       make(map[string]string),
		Templates:   make([]string, 0),
		Assets:      make([]string, 0),
		Metadata:    make(map[string]interface{}),
		PostActions: make([]PostGenerationAction, 0),
	}

	switch g.SecurityLevel {
	case "basic":
		result.Files["internal/security/auth.go"] = "// Basic authentication implementation"
		
	case "rbac":
		result.Files["internal/security/auth.go"] = "// RBAC authentication implementation"
		result.Files["internal/security/rbac.go"] = "// Role-based access control"
		result.Files["configs/rbac-dev.json"] = "// RBAC development configuration"
		
	case "enterprise":
		result.Files["internal/security/auth.go"] = "// Enterprise authentication implementation"
		result.Files["internal/security/rbac.go"] = "// Enterprise RBAC"
		result.Files["internal/security/mtls.go"] = "// Mutual TLS implementation"
		result.Files["internal/compliance/audit.go"] = "// Audit logging implementation"
		result.Files["configs/rbac-production.json"] = "// Production RBAC configuration"
	}

	result.Metadata["security_level"] = g.SecurityLevel
	return result, nil
}

// StorageFeatureGenerator generates storage components
type StorageFeatureGenerator struct {
	BaseFeatureGenerator
	StorageType string // database, cache, file
}

// Generate creates storage implementation
func (g *StorageFeatureGenerator) Generate(ctx context.Context, config *config.ProjectConfig, featureConfig map[string]interface{}) (*GeneratedFeature, error) {
	result := &GeneratedFeature{
		Files:       make(map[string]string),
		Templates:   make([]string, 0),
		Assets:      make([]string, 0),
		Metadata:    make(map[string]interface{}),
		PostActions: make([]PostGenerationAction, 0),
	}

	switch g.StorageType {
	case "database":
		result.Files["internal/storage/database.go"] = "// Database connection and operations"
		result.Files["internal/storage/migrations/001_initial.sql"] = "// Initial database schema"
		result.PostActions = append(result.PostActions, PostGenerationAction{
			Type:        "install",
			Description: "Install database driver",
			Command:     "go",
			Args:        []string{"get", "github.com/lib/pq"},
		})
		
	case "cache":
		result.Files["internal/storage/cache.go"] = "// Cache implementation"
		result.PostActions = append(result.PostActions, PostGenerationAction{
			Type:        "install",
			Description: "Install Redis client",
			Command:     "go",
			Args:        []string{"get", "github.com/go-redis/redis/v8"},
		})
		
	case "file":
		result.Files["internal/storage/file.go"] = "// File storage implementation"
	}

	result.Metadata["storage_type"] = g.StorageType
	return result, nil
}

// APIFeatureGenerator generates API components
type APIFeatureGenerator struct {
	BaseFeatureGenerator
	APIType string // rest, graphql, grpc
}

// Generate creates API implementation
func (g *APIFeatureGenerator) Generate(ctx context.Context, config *config.ProjectConfig, featureConfig map[string]interface{}) (*GeneratedFeature, error) {
	result := &GeneratedFeature{
		Files:       make(map[string]string),
		Templates:   make([]string, 0),
		Assets:      make([]string, 0),
		Metadata:    make(map[string]interface{}),
		PostActions: make([]PostGenerationAction, 0),
	}

	switch g.APIType {
	case "rest":
		result.Files["internal/api/rest/router.go"] = "// REST API router implementation"
		result.Files["internal/api/rest/handlers.go"] = "// REST API handlers"
		result.Templates = append(result.Templates, "template-health/schemas/health-api.tsp")
		
	case "graphql":
		result.Files["internal/api/graphql/schema.go"] = "// GraphQL schema implementation"
		result.Files["internal/api/graphql/resolvers.go"] = "// GraphQL resolvers"
		result.PostActions = append(result.PostActions, PostGenerationAction{
			Type:        "install",
			Description: "Install GraphQL library",
			Command:     "go",
			Args:        []string{"get", "github.com/99designs/gqlgen"},
		})
		
	case "grpc":
		result.Files["internal/api/grpc/server.go"] = "// gRPC server implementation"
		result.Files["api/proto/service.proto"] = "// Protocol buffer definitions"
		result.PostActions = append(result.PostActions, PostGenerationAction{
			Type:        "install",
			Description: "Install gRPC libraries",
			Command:     "go",
			Args:        []string{"get", "google.golang.org/grpc"},
		})
	}

	result.Metadata["api_type"] = g.APIType
	return result, nil
}

// CreatePredefinedFeatures creates and registers all predefined features
func CreatePredefinedFeatures() []*Feature {
	var features []*Feature

	// Health features for different tiers
	healthTiers := []string{"basic", "intermediate", "advanced", "enterprise"}
	for _, tier := range healthTiers {
		feature := &Feature{
			ID:           fmt.Sprintf("health-%s", tier),
			Name:         fmt.Sprintf("Health Checks (%s)", tier),
			Description:  fmt.Sprintf("Health check implementation for %s tier", tier),
			Type:         FeatureTypeCore,
			Version:      "1.0.0",
			Category:     "health",
			Priority:     100,
			MinTier:      tier,
			Dependencies: []string{},
			Generator: &HealthFeatureGenerator{
				BaseFeatureGenerator: BaseFeatureGenerator{
					ID:          fmt.Sprintf("health-%s", tier),
					Name:        fmt.Sprintf("Health-%s", tier),
					Description: fmt.Sprintf("Health checks for %s tier", tier),
				},
				Tier: tier,
			},
		}
		features = append(features, feature)
	}

	// Observability features
	observabilityConfigs := []struct {
		id         string
		name       string
		components []string
		minTier    string
	}{
		{"observability-basic", "Basic Observability", []string{"logging"}, "basic"},
		{"observability-metrics", "Metrics & Logging", []string{"metrics", "logging"}, "intermediate"},
		{"observability-full", "Full Observability", []string{"metrics", "tracing", "logging"}, "advanced"},
		{"observability-enterprise", "Enterprise Observability", []string{"metrics", "tracing", "logging", "sre"}, "enterprise"},
	}

	for _, cfg := range observabilityConfigs {
		feature := &Feature{
			ID:          cfg.id,
			Name:        cfg.name,
			Description: fmt.Sprintf("Observability stack with %v", cfg.components),
			Type:        FeatureTypeObservability,
			Version:     "1.0.0",
			Category:    "observability",
			Priority:    90,
			MinTier:     cfg.minTier,
			Generator: &ObservabilityFeatureGenerator{
				BaseFeatureGenerator: BaseFeatureGenerator{
					ID:          cfg.id,
					Name:        cfg.name,
					Description: fmt.Sprintf("Observability with %v", cfg.components),
				},
				Components: cfg.components,
			},
		}
		features = append(features, feature)
	}

	// Security features
	securityConfigs := []struct {
		id    string
		name  string
		level string
		tier  string
	}{
		{"security-basic", "Basic Security", "basic", "basic"},
		{"security-rbac", "RBAC Security", "rbac", "intermediate"},
		{"security-enterprise", "Enterprise Security", "enterprise", "enterprise"},
	}

	for _, cfg := range securityConfigs {
		feature := &Feature{
			ID:          cfg.id,
			Name:        cfg.name,
			Description: fmt.Sprintf("Security implementation - %s level", cfg.level),
			Type:        FeatureTypeSecurity,
			Version:     "1.0.0",
			Category:    "security",
			Priority:    80,
			MinTier:     cfg.tier,
			Generator: &SecurityFeatureGenerator{
				BaseFeatureGenerator: BaseFeatureGenerator{
					ID:          cfg.id,
					Name:        cfg.name,
					Description: fmt.Sprintf("Security level: %s", cfg.level),
				},
				SecurityLevel: cfg.level,
			},
		}
		features = append(features, feature)
	}

	// Storage features
	storageTypes := []string{"database", "cache", "file"}
	for _, storageType := range storageTypes {
		feature := &Feature{
			ID:          fmt.Sprintf("storage-%s", storageType),
			Name:        fmt.Sprintf("%s Storage", storageType),
			Description: fmt.Sprintf("%s storage implementation", storageType),
			Type:        FeatureTypeStorage,
			Version:     "1.0.0",
			Category:    "storage",
			Priority:    70,
			Conflicts:   []string{}, // Will be populated based on type conflicts
			Generator: &StorageFeatureGenerator{
				BaseFeatureGenerator: BaseFeatureGenerator{
					ID:          fmt.Sprintf("storage-%s", storageType),
					Name:        fmt.Sprintf("Storage-%s", storageType),
					Description: fmt.Sprintf("%s storage", storageType),
				},
				StorageType: storageType,
			},
		}
		
		// Add conflicts for mutually exclusive storage types
		if storageType == "cache" {
			feature.Conflicts = []string{"storage-file"}
		}
		
		features = append(features, feature)
	}

	// API features
	apiTypes := []string{"rest", "graphql", "grpc"}
	for _, apiType := range apiTypes {
		feature := &Feature{
			ID:          fmt.Sprintf("api-%s", apiType),
			Name:        fmt.Sprintf("%s API", apiType),
			Description: fmt.Sprintf("%s API implementation", apiType),
			Type:        FeatureTypeAPI,
			Version:     "1.0.0",
			Category:    "api",
			Priority:    60,
			Generator: &APIFeatureGenerator{
				BaseFeatureGenerator: BaseFeatureGenerator{
					ID:          fmt.Sprintf("api-%s", apiType),
					Name:        fmt.Sprintf("API-%s", apiType),
					Description: fmt.Sprintf("%s API", apiType),
				},
				APIType: apiType,
			},
		}
		features = append(features, feature)
	}

	return features
}

// InitializeFeatureRegistry creates and populates a feature registry with predefined features
func InitializeFeatureRegistry() (*FeatureComposer, error) {
	composer := NewFeatureComposer()
	
	// Register all predefined features
	predefinedFeatures := CreatePredefinedFeatures()
	for _, feature := range predefinedFeatures {
		if err := composer.RegisterFeature(feature); err != nil {
			return nil, fmt.Errorf("failed to register feature %s: %w", feature.ID, err)
		}
	}
	
	return composer, nil
}