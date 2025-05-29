package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/LarsArtmann/BMAD-METHOD/pkg/config"
	"github.com/LarsArtmann/BMAD-METHOD/pkg/generator"
)

var (
	projectName   string
	tier          string
	outputDir     string
	goModule      string
	features      []string
	dryRun        bool
	configFile    string
	interactive   bool
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a health endpoint project from TypeSpec templates",
	Long: `Generate a complete health endpoint project with the specified tier and features.

The generator creates a production-ready health endpoint implementation including:
- TypeSpec API definitions
- Go server implementation
- TypeScript client SDK (optional)
- Kubernetes deployment manifests (optional)
- Docker configuration
- Comprehensive documentation

Available tiers:
  basic        - Simple health endpoints (~5 min deployment)
  intermediate - Production-ready with dependency checks (~15 min deployment)
  advanced     - Full observability with OpenTelemetry (~30 min deployment)
  enterprise   - Enterprise-grade with compliance features (~45 min deployment)

Examples:
  # Generate a basic health service
  template-health-endpoint generate --name my-service --tier basic

  # Generate an advanced service with specific features
  template-health-endpoint generate --name my-service --tier advanced \
    --features opentelemetry,cloudevents,kubernetes

  # Generate from a configuration file
  template-health-endpoint generate --config my-config.yaml

  # Preview what would be generated (dry run)
  template-health-endpoint generate --name my-service --tier basic --dry-run`,
	RunE: runGenerate,
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Required flags
	generateCmd.Flags().StringVarP(&projectName, "name", "n", "", "project name (required)")
	generateCmd.Flags().StringVarP(&tier, "tier", "t", "basic", "template tier (basic|intermediate|advanced|enterprise)")

	// Optional flags
	generateCmd.Flags().StringVarP(&outputDir, "output", "o", "", "output directory (default: project name)")
	generateCmd.Flags().StringVarP(&goModule, "module", "m", "", "Go module path (default: github.com/example/{name})")
	generateCmd.Flags().StringSliceVarP(&features, "features", "f", []string{}, "comma-separated list of features to enable")
	generateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview what would be generated without creating files")
	generateCmd.Flags().StringVarP(&configFile, "config", "c", "", "configuration file path")
	generateCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "interactive mode with prompts")

	// Mark required flags
	generateCmd.MarkFlagRequired("name")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := loadConfiguration()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	// Apply tier-specific defaults
	cfg.ApplyTierDefaults()

	// Show configuration summary
	if viper.GetBool("verbose") || dryRun {
		if err := showConfigurationSummary(cfg); err != nil {
			return fmt.Errorf("failed to show configuration summary: %w", err)
		}
	}

	// Dry run mode - just show what would be generated
	if dryRun {
		fmt.Println("\n🔍 Dry run mode - no files will be created")
		return showGenerationPlan(cfg)
	}

	// Create the generator
	gen, err := generator.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to create generator: %w", err)
	}

	// Generate the project
	fmt.Printf("🚀 Generating %s tier health endpoint project: %s\n", cfg.Tier, cfg.Name)

	if err := gen.Generate(); err != nil {
		return fmt.Errorf("generation failed: %w", err)
	}

	// Show success message with next steps
	return showSuccessMessage(cfg)
}

func loadConfiguration() (*config.ProjectConfig, error) {
	var cfg config.ProjectConfig

	// Load from config file if specified
	if configFile != "" {
		data, err := os.ReadFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	// Override with command line flags
	if projectName != "" {
		cfg.Name = projectName
	}

	if tier != "" {
		cfg.Tier = config.TemplateTier(tier)
	}

	if outputDir != "" {
		cfg.OutputDir = outputDir
	} else if cfg.OutputDir == "" {
		cfg.OutputDir = cfg.Name
	}

	if goModule != "" {
		cfg.GoModule = goModule
	} else if cfg.GoModule == "" {
		cfg.GoModule = fmt.Sprintf("github.com/example/%s", cfg.Name)
	}

	// Parse features flag
	if len(features) > 0 {
		for _, feature := range features {
			switch strings.ToLower(feature) {
			case "opentelemetry", "otel":
				cfg.Features.OpenTelemetry = true
			case "server-timing", "servertiming":
				cfg.Features.ServerTiming = true
			case "cloudevents", "events":
				cfg.Features.CloudEvents = true
			case "kubernetes", "k8s":
				cfg.Features.Kubernetes = true
			case "typescript", "ts":
				cfg.Features.TypeScript = true
			case "docker":
				cfg.Features.Docker = true
			default:
				return nil, fmt.Errorf("unknown feature: %s", feature)
			}
		}
	}

	return &cfg, nil
}

func showConfigurationSummary(cfg *config.ProjectConfig) error {
	fmt.Println("\n📋 Configuration Summary:")
	fmt.Printf("  Project Name: %s\n", cfg.Name)
	fmt.Printf("  Tier: %s\n", cfg.Tier)
	fmt.Printf("  Description: %s\n", cfg.Tier.Description())
	fmt.Printf("  Go Module: %s\n", cfg.GoModule)
	fmt.Printf("  Output Directory: %s\n", cfg.OutputDir)

	fmt.Println("\n🎛️  Features:")
	fmt.Printf("  OpenTelemetry: %v\n", cfg.Features.OpenTelemetry)
	fmt.Printf("  Server Timing: %v\n", cfg.Features.ServerTiming)
	fmt.Printf("  CloudEvents: %v\n", cfg.Features.CloudEvents)
	fmt.Printf("  Kubernetes: %v\n", cfg.Features.Kubernetes)
	fmt.Printf("  TypeScript: %v\n", cfg.Features.TypeScript)
	fmt.Printf("  Docker: %v\n", cfg.Features.Docker)

	if cfg.Features.Kubernetes {
		fmt.Println("\n☸️  Kubernetes Configuration:")
		fmt.Printf("  Service Monitor: %v\n", cfg.Kubernetes.ServiceMonitor)
		fmt.Printf("  Ingress: %v\n", cfg.Kubernetes.Ingress.Enabled)
		fmt.Printf("  Health Probes: Liveness=%v, Readiness=%v, Startup=%v\n",
			cfg.Kubernetes.HealthProbes.LivenessProbe.Enabled,
			cfg.Kubernetes.HealthProbes.ReadinessProbe.Enabled,
			cfg.Kubernetes.HealthProbes.StartupProbe.Enabled)
	}

	return nil
}

func showGenerationPlan(cfg *config.ProjectConfig) error {
	fmt.Println("\n📁 Files that would be generated:")

	// Core files
	files := []string{
		"README.md",
		"go.mod",
		"go.sum",
		"Dockerfile",
		"docker-compose.yml",
		".gitignore",
		"Makefile",
	}

	// Go source files
	goFiles := []string{
		"cmd/server/main.go",
		"internal/handlers/health.go",
		"internal/handlers/server_time.go",
		"internal/models/health.go",
		"internal/server/server.go",
		"internal/config/config.go",
	}

	// TypeScript files (if enabled)
	var tsFiles []string
	if cfg.Features.TypeScript {
		tsFiles = []string{
			"client/typescript/src/client.ts",
			"client/typescript/src/types.ts",
			"client/typescript/package.json",
			"client/typescript/tsconfig.json",
		}
	}

	// Kubernetes files (if enabled)
	var k8sFiles []string
	if cfg.Features.Kubernetes {
		k8sFiles = []string{
			"deployments/kubernetes/deployment.yaml",
			"deployments/kubernetes/service.yaml",
			"deployments/kubernetes/configmap.yaml",
		}

		if cfg.Kubernetes.ServiceMonitor {
			k8sFiles = append(k8sFiles, "deployments/kubernetes/servicemonitor.yaml")
		}

		if cfg.Kubernetes.Ingress.Enabled {
			k8sFiles = append(k8sFiles, "deployments/kubernetes/ingress.yaml")
		}
	}

	// Print file lists
	printFileList("Core Files", files)
	printFileList("Go Source Files", goFiles)

	if len(tsFiles) > 0 {
		printFileList("TypeScript Client Files", tsFiles)
	}

	if len(k8sFiles) > 0 {
		printFileList("Kubernetes Manifests", k8sFiles)
	}

	// Show estimated deployment time
	fmt.Printf("\n⏱️  Estimated deployment time: %s\n", cfg.Tier.Description())

	return nil
}

func printFileList(title string, files []string) {
	fmt.Printf("\n  %s:\n", title)
	for _, file := range files {
		fmt.Printf("    - %s\n", file)
	}
}

func showSuccessMessage(cfg *config.ProjectConfig) error {
	fmt.Printf("\n✅ Successfully generated %s tier health endpoint project!\n", cfg.Tier)
	fmt.Printf("\n📁 Project created in: %s\n", cfg.OutputDir)

	fmt.Println("\n🚀 Next steps:")
	fmt.Printf("  1. cd %s\n", cfg.OutputDir)
	fmt.Println("  2. go mod tidy")
	fmt.Println("  3. go run cmd/server/main.go")
	fmt.Println("  4. curl http://localhost:8080/health")

	if cfg.Features.TypeScript {
		fmt.Println("\n📦 TypeScript client:")
		fmt.Printf("  cd %s/client/typescript && npm install\n", cfg.OutputDir)
	}

	if cfg.Features.Kubernetes {
		fmt.Println("\n☸️  Kubernetes deployment:")
		fmt.Printf("  kubectl apply -f %s/deployments/kubernetes/\n", cfg.OutputDir)
	}

	fmt.Printf("\n📚 Documentation: %s/README.md\n", cfg.OutputDir)

	return nil
}
