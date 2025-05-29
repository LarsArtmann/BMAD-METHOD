package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/LarsArtmann/BMAD-METHOD/pkg/config"
	"github.com/LarsArtmann/BMAD-METHOD/pkg/generator"
)

var (
	customizeProfile    string
	customizeOutputDir  string
	customizeConfigFile string
	customizeInteractive bool
	customizeSaveProfile bool
)

// customizeCmd represents the customize command
var customizeCmd = &cobra.Command{
	Use:   "customize",
	Short: "Interactively customize template generation with advanced options",
	Long: `Interactively customize template generation with advanced configuration options.

This command provides an interactive interface to customize templates beyond the
standard tier options, allowing fine-grained control over features, dependencies,
and configurations.

Customization options include:
- Feature selection (enable/disable specific features)
- Dependency configuration (versions, optional dependencies)
- Environment settings (development, staging, production)
- Security settings (authentication, authorization)
- Observability configuration (metrics, tracing, logging)
- Kubernetes settings (resources, scaling, networking)

Customization profiles can be saved and reused for consistent project generation.

Examples:
  # Interactive customization
  template-health-endpoint customize --name my-service

  # Use a saved customization profile
  template-health-endpoint customize --profile my-profile --name my-service

  # Customize and save as new profile
  template-health-endpoint customize --name my-service --save-profile my-custom-profile

  # Load customization from config file
  template-health-endpoint customize --config custom-config.yaml --name my-service`,
	RunE: runCustomizeTemplate,
}

func init() {
	rootCmd.AddCommand(customizeCmd)

	// Flags
	customizeCmd.Flags().StringVar(&customizeProfile, "profile", "", "use saved customization profile")
	customizeCmd.Flags().StringVarP(&customizeOutputDir, "output", "o", ".", "output directory for generated project")
	customizeCmd.Flags().StringVar(&customizeConfigFile, "config", "", "load customization from config file")
	customizeCmd.Flags().BoolVarP(&customizeInteractive, "interactive", "i", true, "interactive customization mode")
	customizeCmd.Flags().BoolVar(&customizeSaveProfile, "save-profile", false, "save customization as profile")

	// Bind flags to viper
	viper.BindPFlag("customize.profile", customizeCmd.Flags().Lookup("profile"))
	viper.BindPFlag("customize.output", customizeCmd.Flags().Lookup("output"))
	viper.BindPFlag("customize.config", customizeCmd.Flags().Lookup("config"))
	viper.BindPFlag("customize.interactive", customizeCmd.Flags().Lookup("interactive"))
	viper.BindPFlag("customize.save-profile", customizeCmd.Flags().Lookup("save-profile"))
}

func runCustomizeTemplate(cmd *cobra.Command, args []string) error {
	fmt.Println("🎨 Template Customization Wizard")
	fmt.Println("=================================")

	// 1. Load base customization
	customization, err := loadBaseCustomization()
	if err != nil {
		return fmt.Errorf("failed to load base customization: %w", err)
	}

	// 2. Load from profile or config file if specified
	if customizeProfile != "" {
		if err := loadCustomizationProfile(customization, customizeProfile); err != nil {
			return fmt.Errorf("failed to load profile '%s': %w", customizeProfile, err)
		}
		fmt.Printf("📋 Loaded profile: %s\n\n", customizeProfile)
	}

	if customizeConfigFile != "" {
		if err := loadCustomizationConfig(customization, customizeConfigFile); err != nil {
			return fmt.Errorf("failed to load config file '%s': %w", customizeConfigFile, err)
		}
		fmt.Printf("📋 Loaded config: %s\n\n", customizeConfigFile)
	}

	// 3. Interactive customization
	if customizeInteractive {
		if err := runInteractiveCustomization(customization); err != nil {
			return fmt.Errorf("interactive customization failed: %w", err)
		}
	}

	// 4. Show customization summary
	showCustomizationSummary(customization)

	// 5. Confirm generation
	if !confirmGeneration() {
		fmt.Println("❌ Customization cancelled.")
		return nil
	}

	// 6. Save profile if requested
	if customizeSaveProfile {
		profileName := promptForInput("Enter profile name to save", "my-custom-profile")
		if err := saveCustomizationProfile(customization, profileName); err != nil {
			fmt.Printf("⚠️  Failed to save profile: %v\n", err)
		} else {
			fmt.Printf("💾 Profile saved: %s\n", profileName)
		}
	}

	// 7. Generate project with customization
	if err := generateCustomizedProject(customization); err != nil {
		return fmt.Errorf("failed to generate customized project: %w", err)
	}

	fmt.Printf("✅ Customized project generated successfully!\n")
	fmt.Printf("   Output directory: %s\n", customizeOutputDir)

	return nil
}

// Customization holds all customization options
type Customization struct {
	// Project settings
	ProjectName string `yaml:"project_name"`
	GoModule    string `yaml:"go_module"`
	BaseTier    string `yaml:"base_tier"`

	// Features
	Features FeatureCustomization `yaml:"features"`

	// Dependencies
	Dependencies DependencyCustomization `yaml:"dependencies"`

	// Environment settings
	Environment EnvironmentCustomization `yaml:"environment"`

	// Security settings
	Security SecurityCustomization `yaml:"security"`

	// Observability settings
	Observability ObservabilityCustomization `yaml:"observability"`

	// Kubernetes settings
	Kubernetes KubernetesCustomization `yaml:"kubernetes"`
}

type FeatureCustomization struct {
	TypeScript   bool `yaml:"typescript"`
	Docker       bool `yaml:"docker"`
	Kubernetes   bool `yaml:"kubernetes"`
	OpenTelemetry bool `yaml:"opentelemetry"`
	CloudEvents  bool `yaml:"cloudevents"`
	Dependencies bool `yaml:"dependencies"`
	ServerTiming bool `yaml:"server_timing"`
	Metrics      bool `yaml:"metrics"`
	MTLS         bool `yaml:"mtls"`
	RBAC         bool `yaml:"rbac"`
	AuditLogging bool `yaml:"audit_logging"`
}

type DependencyCustomization struct {
	GoVersion      string            `yaml:"go_version"`
	CustomPackages map[string]string `yaml:"custom_packages"`
	OptionalDeps   []string          `yaml:"optional_deps"`
}

type EnvironmentCustomization struct {
	Environments     []string          `yaml:"environments"`
	DefaultEnv       string            `yaml:"default_env"`
	ConfigFormat     string            `yaml:"config_format"` // yaml, json, toml
	EnvironmentVars  map[string]string `yaml:"environment_vars"`
}

type SecurityCustomization struct {
	EnableMTLS       bool     `yaml:"enable_mtls"`
	EnableRBAC       bool     `yaml:"enable_rbac"`
	AuthMethods      []string `yaml:"auth_methods"`
	CertificatePaths struct {
		CertFile string `yaml:"cert_file"`
		KeyFile  string `yaml:"key_file"`
		CAFile   string `yaml:"ca_file"`
	} `yaml:"certificate_paths"`
}

type ObservabilityCustomization struct {
	MetricsEnabled  bool   `yaml:"metrics_enabled"`
	TracingEnabled  bool   `yaml:"tracing_enabled"`
	LoggingLevel    string `yaml:"logging_level"`
	MetricsPort     int    `yaml:"metrics_port"`
	TracingEndpoint string `yaml:"tracing_endpoint"`
}

type KubernetesCustomization struct {
	Namespace       string            `yaml:"namespace"`
	ResourceLimits  map[string]string `yaml:"resource_limits"`
	Replicas        int               `yaml:"replicas"`
	ServiceType     string            `yaml:"service_type"`
	IngressEnabled  bool              `yaml:"ingress_enabled"`
	ServiceMonitor  bool              `yaml:"service_monitor"`
}

// loadBaseCustomization loads default customization settings
func loadBaseCustomization() (*Customization, error) {
	return &Customization{
		ProjectName: "my-health-service",
		GoModule:    "github.com/example/my-health-service",
		BaseTier:    "intermediate",
		Features: FeatureCustomization{
			TypeScript:   true,
			Docker:       true,
			Kubernetes:   true,
			OpenTelemetry: false,
			CloudEvents:  false,
			Dependencies: true,
			ServerTiming: false,
			Metrics:      true,
			MTLS:         false,
			RBAC:         false,
			AuditLogging: false,
		},
		Dependencies: DependencyCustomization{
			GoVersion:      "1.21",
			CustomPackages: make(map[string]string),
			OptionalDeps:   []string{},
		},
		Environment: EnvironmentCustomization{
			Environments: []string{"development", "staging", "production"},
			DefaultEnv:   "development",
			ConfigFormat: "yaml",
			EnvironmentVars: make(map[string]string),
		},
		Security: SecurityCustomization{
			EnableMTLS: false,
			EnableRBAC: false,
			AuthMethods: []string{},
		},
		Observability: ObservabilityCustomization{
			MetricsEnabled:  true,
			TracingEnabled:  false,
			LoggingLevel:    "info",
			MetricsPort:     9090,
			TracingEndpoint: "",
		},
		Kubernetes: KubernetesCustomization{
			Namespace:      "default",
			ResourceLimits: map[string]string{"cpu": "100m", "memory": "128Mi"},
			Replicas:       1,
			ServiceType:    "ClusterIP",
			IngressEnabled: false,
			ServiceMonitor: false,
		},
	}, nil
}

// runInteractiveCustomization runs the interactive customization wizard
func runInteractiveCustomization(c *Customization) error {
	_ = bufio.NewReader(os.Stdin)

	// Project settings
	fmt.Println("\n📋 Project Settings")
	fmt.Println("-------------------")
	c.ProjectName = promptForInput("Project name", c.ProjectName)
	c.GoModule = promptForInput("Go module", c.GoModule)
	c.BaseTier = promptForChoice("Base tier", []string{"basic", "intermediate", "advanced", "enterprise"}, c.BaseTier)

	// Features
	fmt.Println("\n🎯 Features")
	fmt.Println("-----------")
	c.Features.TypeScript = promptForBool("Include TypeScript client SDK", c.Features.TypeScript)
	c.Features.Docker = promptForBool("Include Docker configuration", c.Features.Docker)
	c.Features.Kubernetes = promptForBool("Include Kubernetes manifests", c.Features.Kubernetes)
	c.Features.OpenTelemetry = promptForBool("Enable OpenTelemetry observability", c.Features.OpenTelemetry)
	c.Features.CloudEvents = promptForBool("Enable CloudEvents integration", c.Features.CloudEvents)
	c.Features.Dependencies = promptForBool("Include dependency health checks", c.Features.Dependencies)

	// Security (if advanced features)
	if c.BaseTier == "advanced" || c.BaseTier == "enterprise" {
		fmt.Println("\n🔒 Security")
		fmt.Println("-----------")
		c.Security.EnableMTLS = promptForBool("Enable mutual TLS (mTLS)", c.Security.EnableMTLS)
		c.Security.EnableRBAC = promptForBool("Enable role-based access control", c.Security.EnableRBAC)
	}

	// Observability
	fmt.Println("\n📊 Observability")
	fmt.Println("----------------")
	c.Observability.MetricsEnabled = promptForBool("Enable metrics collection", c.Observability.MetricsEnabled)
	c.Observability.TracingEnabled = promptForBool("Enable distributed tracing", c.Observability.TracingEnabled)
	c.Observability.LoggingLevel = promptForChoice("Logging level", []string{"debug", "info", "warn", "error"}, c.Observability.LoggingLevel)

	// Kubernetes
	if c.Features.Kubernetes {
		fmt.Println("\n☸️  Kubernetes")
		fmt.Println("--------------")
		c.Kubernetes.Namespace = promptForInput("Kubernetes namespace", c.Kubernetes.Namespace)
		c.Kubernetes.Replicas = promptForInt("Number of replicas", c.Kubernetes.Replicas)
		c.Kubernetes.ServiceType = promptForChoice("Service type", []string{"ClusterIP", "NodePort", "LoadBalancer"}, c.Kubernetes.ServiceType)
		c.Kubernetes.IngressEnabled = promptForBool("Enable Ingress", c.Kubernetes.IngressEnabled)
	}

	return nil
}

// showCustomizationSummary displays the final customization summary
func showCustomizationSummary(c *Customization) {
	fmt.Println("\n📋 Customization Summary")
	fmt.Println("========================")
	fmt.Printf("Project Name:     %s\n", c.ProjectName)
	fmt.Printf("Go Module:        %s\n", c.GoModule)
	fmt.Printf("Base Tier:        %s\n", c.BaseTier)
	fmt.Printf("Output Directory: %s\n", customizeOutputDir)

	fmt.Println("\nEnabled Features:")
	if c.Features.TypeScript { fmt.Println("  ✅ TypeScript client SDK") }
	if c.Features.Docker { fmt.Println("  ✅ Docker configuration") }
	if c.Features.Kubernetes { fmt.Println("  ✅ Kubernetes manifests") }
	if c.Features.OpenTelemetry { fmt.Println("  ✅ OpenTelemetry observability") }
	if c.Features.CloudEvents { fmt.Println("  ✅ CloudEvents integration") }
	if c.Features.Dependencies { fmt.Println("  ✅ Dependency health checks") }
	if c.Security.EnableMTLS { fmt.Println("  ✅ Mutual TLS (mTLS)") }
	if c.Security.EnableRBAC { fmt.Println("  ✅ Role-based access control") }

	fmt.Printf("\nObservability:\n")
	fmt.Printf("  Metrics:  %v\n", c.Observability.MetricsEnabled)
	fmt.Printf("  Tracing:  %v\n", c.Observability.TracingEnabled)
	fmt.Printf("  Log Level: %s\n", c.Observability.LoggingLevel)

	if c.Features.Kubernetes {
		fmt.Printf("\nKubernetes:\n")
		fmt.Printf("  Namespace: %s\n", c.Kubernetes.Namespace)
		fmt.Printf("  Replicas:  %d\n", c.Kubernetes.Replicas)
		fmt.Printf("  Service:   %s\n", c.Kubernetes.ServiceType)
		fmt.Printf("  Ingress:   %v\n", c.Kubernetes.IngressEnabled)
	}
}

// generateCustomizedProject generates the project with customizations
func generateCustomizedProject(c *Customization) error {
	// Convert customization to project config
	projectConfig := &config.ProjectConfig{
		Name:        c.ProjectName,
		GoModule:    c.GoModule,
		Tier:        config.TemplateTier(c.BaseTier),
		OutputDir:   customizeOutputDir,
		Version:     "1.0.0",
		Description: fmt.Sprintf("Customized %s tier health endpoint service", c.BaseTier),
	}

	// Apply feature configuration
	projectConfig.Features.TypeScript = c.Features.TypeScript
	projectConfig.Features.Docker = c.Features.Docker
	projectConfig.Features.Kubernetes = c.Features.Kubernetes
	projectConfig.Features.OpenTelemetry = c.Features.OpenTelemetry
	projectConfig.Features.CloudEvents = c.Features.CloudEvents
	projectConfig.Features.ServerTiming = c.Features.ServerTiming

	// Apply observability configuration
	projectConfig.Observability.Metrics.Enabled = c.Observability.MetricsEnabled
	projectConfig.Observability.OpenTelemetry.Enabled = c.Observability.TracingEnabled

	// Apply Kubernetes configuration
	if c.Features.Kubernetes {
		projectConfig.Kubernetes.Enabled = true
		projectConfig.Kubernetes.Namespace = c.Kubernetes.Namespace
		projectConfig.Kubernetes.ServiceMonitor = c.Kubernetes.ServiceMonitor
		projectConfig.Kubernetes.Ingress.Enabled = c.Kubernetes.IngressEnabled
	}

	// Validate configuration
	if err := projectConfig.Validate(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	// Apply tier defaults
	projectConfig.ApplyTierDefaults()

	// Create generator
	gen, err := generator.New(projectConfig)
	if err != nil {
		return err
	}

	// Generate project
	return gen.Generate()
}

// Helper functions for interactive prompts

func promptForInput(prompt, defaultValue string) string {
	fmt.Printf("%s [%s]: ", prompt, defaultValue)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue
	}
	return input
}

func promptForBool(prompt string, defaultValue bool) bool {
	defaultStr := "n"
	if defaultValue {
		defaultStr = "y"
	}

	response := promptForInput(fmt.Sprintf("%s (y/n)", prompt), defaultStr)
	return strings.ToLower(response) == "y" || strings.ToLower(response) == "yes"
}

func promptForInt(prompt string, defaultValue int) int {
	response := promptForInput(prompt, strconv.Itoa(defaultValue))
	if val, err := strconv.Atoi(response); err == nil {
		return val
	}
	return defaultValue
}

func promptForChoice(prompt string, choices []string, defaultValue string) string {
	fmt.Printf("%s (%s) [%s]: ", prompt, strings.Join(choices, "/"), defaultValue)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultValue
	}

	for _, choice := range choices {
		if strings.ToLower(input) == strings.ToLower(choice) {
			return choice
		}
	}

	return defaultValue
}

func confirmGeneration() bool {
	return promptForBool("\nGenerate project with these settings?", true)
}

// Profile management functions

func loadCustomizationProfile(c *Customization, profileName string) error {
	profilePath := filepath.Join(os.Getenv("HOME"), ".template-health-endpoint", "profiles", profileName+".yaml")
	data, err := os.ReadFile(profilePath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, c)
}

func saveCustomizationProfile(c *Customization, profileName string) error {
	profileDir := filepath.Join(os.Getenv("HOME"), ".template-health-endpoint", "profiles")
	if err := os.MkdirAll(profileDir, 0755); err != nil {
		return err
	}

	profilePath := filepath.Join(profileDir, profileName+".yaml")
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(profilePath, data, 0644)
}

func loadCustomizationConfig(c *Customization, configFile string) error {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, c)
}

func convertFeatures(f FeatureCustomization) map[string]bool {
	return map[string]bool{
		"typescript":    f.TypeScript,
		"docker":        f.Docker,
		"kubernetes":    f.Kubernetes,
		"opentelemetry": f.OpenTelemetry,
		"cloudevents":   f.CloudEvents,
		"dependencies":  f.Dependencies,
		"server_timing": f.ServerTiming,
		"metrics":       f.Metrics,
		"mtls":          f.MTLS,
		"rbac":          f.RBAC,
		"audit_logging": f.AuditLogging,
	}
}
