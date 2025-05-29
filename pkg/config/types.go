package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// TemplateTier represents the complexity tier of the generated template
type TemplateTier string

const (
	// TierBasic provides simple health endpoints (5 min deployment)
	TierBasic TemplateTier = "basic"

	// TierIntermediate adds dependency checks and basic observability (15 min deployment)
	TierIntermediate TemplateTier = "intermediate"

	// TierAdvanced includes full observability and CloudEvents (30 min deployment)
	TierAdvanced TemplateTier = "advanced"

	// TierEnterprise provides compliance and enterprise features (45 min deployment)
	TierEnterprise TemplateTier = "enterprise"
)

// IsValid checks if the tier is a valid option
func (t TemplateTier) IsValid() bool {
	switch t {
	case TierBasic, TierIntermediate, TierAdvanced, TierEnterprise:
		return true
	default:
		return false
	}
}

// String returns the string representation of the tier
func (t TemplateTier) String() string {
	return string(t)
}

// Description returns a human-readable description of the tier
func (t TemplateTier) Description() string {
	switch t {
	case TierBasic:
		return "Basic health endpoints with ServerTime API (~5 min deployment)"
	case TierIntermediate:
		return "Production-ready with dependency checks and basic observability (~15 min deployment)"
	case TierAdvanced:
		return "Full observability with OpenTelemetry and CloudEvents (~30 min deployment)"
	case TierEnterprise:
		return "Enterprise-grade with compliance and advanced monitoring (~45 min deployment)"
	default:
		return "Unknown tier"
	}
}

// ProjectConfig represents the configuration for generating a health endpoint project
type ProjectConfig struct {
	// Project metadata
	Name        string `yaml:"name" mapstructure:"name"`
	Version     string `yaml:"version" mapstructure:"version"`
	Description string `yaml:"description" mapstructure:"description"`
	GoModule    string `yaml:"go_module" mapstructure:"go_module"`

	// Template configuration
	Tier     TemplateTier `yaml:"tier" mapstructure:"tier"`
	OutputDir string      `yaml:"output_dir" mapstructure:"output_dir"`

	// Feature flags
	Features FeatureConfig `yaml:"features" mapstructure:"features"`

	// Dependencies configuration
	Dependencies DependencyConfig `yaml:"dependencies" mapstructure:"dependencies"`

	// Kubernetes configuration
	Kubernetes KubernetesConfig `yaml:"kubernetes" mapstructure:"kubernetes"`

	// Observability configuration
	Observability ObservabilityConfig `yaml:"observability" mapstructure:"observability"`
}

// FeatureConfig controls which features are enabled
type FeatureConfig struct {
	OpenTelemetry bool `yaml:"opentelemetry" mapstructure:"opentelemetry"`
	ServerTiming  bool `yaml:"server_timing" mapstructure:"server_timing"`
	CloudEvents   bool `yaml:"cloudevents" mapstructure:"cloudevents"`
	Kubernetes    bool `yaml:"kubernetes" mapstructure:"kubernetes"`
	TypeScript    bool `yaml:"typescript" mapstructure:"typescript"`
	Docker        bool `yaml:"docker" mapstructure:"docker"`
	Security      bool `yaml:"security" mapstructure:"security"`
	Compliance    bool `yaml:"compliance" mapstructure:"compliance"`
}

// DependencyConfig configures external dependency health checks
type DependencyConfig struct {
	DatabaseChecks    bool     `yaml:"database_checks" mapstructure:"database_checks"`
	ExternalServices  []string `yaml:"external_services" mapstructure:"external_services"`
	CacheChecks       bool     `yaml:"cache_checks" mapstructure:"cache_checks"`
	FilesystemChecks  bool     `yaml:"filesystem_checks" mapstructure:"filesystem_checks"`
	MemoryChecks      bool     `yaml:"memory_checks" mapstructure:"memory_checks"`
}

// KubernetesConfig configures Kubernetes integration
type KubernetesConfig struct {
	Enabled        bool              `yaml:"enabled" mapstructure:"enabled"`
	Namespace      string            `yaml:"namespace" mapstructure:"namespace"`
	ServiceName    string            `yaml:"service_name" mapstructure:"service_name"`
	Port           int               `yaml:"port" mapstructure:"port"`
	Labels         map[string]string `yaml:"labels" mapstructure:"labels"`
	Annotations    map[string]string `yaml:"annotations" mapstructure:"annotations"`
	HealthProbes   HealthProbeConfig `yaml:"health_probes" mapstructure:"health_probes"`
	ServiceMonitor bool              `yaml:"service_monitor" mapstructure:"service_monitor"`
	Ingress        IngressConfig     `yaml:"ingress" mapstructure:"ingress"`
}

// HealthProbeConfig configures Kubernetes health probes
type HealthProbeConfig struct {
	LivenessProbe  ProbeConfig `yaml:"liveness_probe" mapstructure:"liveness_probe"`
	ReadinessProbe ProbeConfig `yaml:"readiness_probe" mapstructure:"readiness_probe"`
	StartupProbe   ProbeConfig `yaml:"startup_probe" mapstructure:"startup_probe"`
}

// ProbeConfig represents a single health probe configuration
type ProbeConfig struct {
	Enabled             bool          `yaml:"enabled" mapstructure:"enabled"`
	Path                string        `yaml:"path" mapstructure:"path"`
	InitialDelaySeconds int           `yaml:"initial_delay_seconds" mapstructure:"initial_delay_seconds"`
	PeriodSeconds       int           `yaml:"period_seconds" mapstructure:"period_seconds"`
	TimeoutSeconds      int           `yaml:"timeout_seconds" mapstructure:"timeout_seconds"`
	FailureThreshold    int           `yaml:"failure_threshold" mapstructure:"failure_threshold"`
	SuccessThreshold    int           `yaml:"success_threshold" mapstructure:"success_threshold"`
}

// IngressConfig configures Kubernetes Ingress
type IngressConfig struct {
	Enabled     bool              `yaml:"enabled" mapstructure:"enabled"`
	Host        string            `yaml:"host" mapstructure:"host"`
	Path        string            `yaml:"path" mapstructure:"path"`
	TLS         bool              `yaml:"tls" mapstructure:"tls"`
	Annotations map[string]string `yaml:"annotations" mapstructure:"annotations"`
}

// ObservabilityConfig configures observability features
type ObservabilityConfig struct {
	OpenTelemetry OpenTelemetryConfig `yaml:"opentelemetry" mapstructure:"opentelemetry"`
	ServerTiming  ServerTimingConfig  `yaml:"server_timing" mapstructure:"server_timing"`
	CloudEvents   CloudEventsConfig   `yaml:"cloudevents" mapstructure:"cloudevents"`
	Metrics       MetricsConfig       `yaml:"metrics" mapstructure:"metrics"`
}

// OpenTelemetryConfig configures OpenTelemetry integration
type OpenTelemetryConfig struct {
	Enabled     bool   `yaml:"enabled" mapstructure:"enabled"`
	ServiceName string `yaml:"service_name" mapstructure:"service_name"`
	Endpoint    string `yaml:"endpoint" mapstructure:"endpoint"`
	Tracing     bool   `yaml:"tracing" mapstructure:"tracing"`
	Metrics     bool   `yaml:"metrics" mapstructure:"metrics"`
	Logging     bool   `yaml:"logging" mapstructure:"logging"`
}

// ServerTimingConfig configures Server Timing API
type ServerTimingConfig struct {
	Enabled     bool     `yaml:"enabled" mapstructure:"enabled"`
	Metrics     []string `yaml:"metrics" mapstructure:"metrics"`
	MaxDuration string   `yaml:"max_duration" mapstructure:"max_duration"`
}

// CloudEventsConfig configures CloudEvents integration
type CloudEventsConfig struct {
	Enabled       bool   `yaml:"enabled" mapstructure:"enabled"`
	BrokerURL     string `yaml:"broker_url" mapstructure:"broker_url"`
	TopicPattern  string `yaml:"topic_pattern" mapstructure:"topic_pattern"`
	Source        string `yaml:"source" mapstructure:"source"`
	EventTypes    []string `yaml:"event_types" mapstructure:"event_types"`
}

// MetricsConfig configures metrics collection
type MetricsConfig struct {
	Enabled    bool   `yaml:"enabled" mapstructure:"enabled"`
	Prometheus bool   `yaml:"prometheus" mapstructure:"prometheus"`
	Port       int    `yaml:"port" mapstructure:"port"`
	Path       string `yaml:"path" mapstructure:"path"`
}

// Validate validates the project configuration
func (c *ProjectConfig) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("project name is required")
	}

	if !c.Tier.IsValid() {
		return fmt.Errorf("invalid tier: %s (must be one of: basic, intermediate, advanced, enterprise)", c.Tier)
	}

	if c.GoModule == "" {
		return fmt.Errorf("go module path is required")
	}

	if c.OutputDir == "" {
		c.OutputDir = c.Name
	}

	if c.Version == "" {
		c.Version = "1.0.0"
	}

	return nil
}

// ApplyTierDefaults applies default configuration based on the selected tier
func (c *ProjectConfig) ApplyTierDefaults() {
	switch c.Tier {
	case TierBasic:
		c.Features.OpenTelemetry = false
		c.Features.ServerTiming = false
		c.Features.CloudEvents = false
		c.Features.Kubernetes = true
		c.Features.TypeScript = true
		c.Features.Docker = true

		c.Dependencies.DatabaseChecks = false
		c.Dependencies.CacheChecks = false
		c.Dependencies.FilesystemChecks = true
		c.Dependencies.MemoryChecks = true

	case TierIntermediate:
		c.Features.OpenTelemetry = true
		c.Features.ServerTiming = false
		c.Features.CloudEvents = false
		c.Features.Kubernetes = true
		c.Features.TypeScript = true
		c.Features.Docker = true

		c.Dependencies.DatabaseChecks = true
		c.Dependencies.CacheChecks = true
		c.Dependencies.FilesystemChecks = true
		c.Dependencies.MemoryChecks = true

		c.Observability.OpenTelemetry.Enabled = true
		c.Observability.OpenTelemetry.Tracing = true
		c.Observability.OpenTelemetry.Metrics = true

	case TierAdvanced:
		c.Features.OpenTelemetry = true
		c.Features.ServerTiming = true
		c.Features.CloudEvents = true
		c.Features.Kubernetes = true
		c.Features.TypeScript = true
		c.Features.Docker = true

		c.Dependencies.DatabaseChecks = true
		c.Dependencies.CacheChecks = true
		c.Dependencies.FilesystemChecks = true
		c.Dependencies.MemoryChecks = true

		c.Observability.OpenTelemetry.Enabled = true
		c.Observability.OpenTelemetry.Tracing = true
		c.Observability.OpenTelemetry.Metrics = true
		c.Observability.OpenTelemetry.Logging = true

		c.Observability.ServerTiming.Enabled = true
		c.Observability.CloudEvents.Enabled = true

	case TierEnterprise:
		c.Features.OpenTelemetry = true
		c.Features.ServerTiming = true
		c.Features.CloudEvents = true
		c.Features.Kubernetes = true
		c.Features.TypeScript = true
		c.Features.Docker = true
		c.Features.Security = true
		c.Features.Compliance = true

		c.Dependencies.DatabaseChecks = true
		c.Dependencies.CacheChecks = true
		c.Dependencies.FilesystemChecks = true
		c.Dependencies.MemoryChecks = true

		c.Observability.OpenTelemetry.Enabled = true
		c.Observability.OpenTelemetry.Tracing = true
		c.Observability.OpenTelemetry.Metrics = true
		c.Observability.OpenTelemetry.Logging = true

		c.Observability.ServerTiming.Enabled = true
		c.Observability.CloudEvents.Enabled = true
		c.Observability.Metrics.Enabled = true
		c.Observability.Metrics.Prometheus = true

		c.Kubernetes.Enabled = true
		c.Kubernetes.ServiceMonitor = true
		c.Kubernetes.Ingress.Enabled = true
	}

	// Apply default Kubernetes health probe configurations
	if c.Features.Kubernetes {
		c.applyDefaultHealthProbes()
	}
}

// applyDefaultHealthProbes sets up default health probe configurations
func (c *ProjectConfig) applyDefaultHealthProbes() {
	c.Kubernetes.HealthProbes.LivenessProbe = ProbeConfig{
		Enabled:             true,
		Path:                "/health/live",
		InitialDelaySeconds: 30,
		PeriodSeconds:       10,
		TimeoutSeconds:      5,
		FailureThreshold:    3,
		SuccessThreshold:    1,
	}

	c.Kubernetes.HealthProbes.ReadinessProbe = ProbeConfig{
		Enabled:             true,
		Path:                "/health/ready",
		InitialDelaySeconds: 5,
		PeriodSeconds:       5,
		TimeoutSeconds:      3,
		FailureThreshold:    3,
		SuccessThreshold:    1,
	}

	c.Kubernetes.HealthProbes.StartupProbe = ProbeConfig{
		Enabled:             true,
		Path:                "/health/startup",
		InitialDelaySeconds: 10,
		PeriodSeconds:       10,
		TimeoutSeconds:      5,
		FailureThreshold:    30,
		SuccessThreshold:    1,
	}
}

// TemplateConfig represents the configuration for a template tier
type TemplateConfig struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Tier        string            `yaml:"tier"`
	Version     string            `yaml:"version"`
	Features    map[string]bool   `yaml:"features"`
	Metadata    map[string]string `yaml:"metadata,omitempty"`
}

// LoadTemplateConfig loads a template configuration from a YAML file
func LoadTemplateConfig(path string) (*TemplateConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read template config: %w", err)
	}

	var config TemplateConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse template config: %w", err)
	}

	return &config, nil
}

// GeneratorConfig represents the configuration for the generator
type GeneratorConfig struct {
	ProjectName string            `yaml:"project_name"`
	GoModule    string            `yaml:"go_module"`
	Tier        string            `yaml:"tier"`
	OutputDir   string            `yaml:"output_dir"`
	Features    map[string]bool   `yaml:"features"`
	Variables   map[string]string `yaml:"variables,omitempty"`
}

// Validate validates the generator configuration
func (g *GeneratorConfig) Validate() error {
	if g.ProjectName == "" {
		return fmt.Errorf("project name is required")
	}
	if g.GoModule == "" {
		return fmt.Errorf("go module is required")
	}
	if g.Tier == "" {
		return fmt.Errorf("tier is required")
	}
	if g.OutputDir == "" {
		g.OutputDir = g.ProjectName
	}
	return nil
}