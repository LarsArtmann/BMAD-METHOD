package commands

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/LarsArtmann/BMAD-METHOD/pkg/config"
)

// InteractiveWizard runs the interactive project configuration wizard
func InteractiveWizard() (*config.ProjectConfig, error) {
	fmt.Println("🧙 Welcome to the BMAD Method Health Endpoint Generator!")
	fmt.Println("Let's create your perfect health endpoint project step by step.\n")

	var cfg config.ProjectConfig

	// Step 1: Project basics
	if err := askProjectBasics(&cfg); err != nil {
		return nil, err
	}

	// Step 2: Tier selection with smart recommendations
	if err := askTierSelection(&cfg); err != nil {
		return nil, err
	}

	// Step 3: Feature selection based on tier
	if err := askFeatureSelection(&cfg); err != nil {
		return nil, err
	}

	// Step 4: Advanced configuration (optional)
	if err := askAdvancedConfiguration(&cfg); err != nil {
		return nil, err
	}

	// Apply tier-specific defaults
	cfg.ApplyTierDefaults()

	return &cfg, nil
}

func askProjectBasics(cfg *config.ProjectConfig) error {
	var qs = []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: "What's your project name?",
				Help:    "This will be used for the directory name and Go module. Use kebab-case (e.g., my-health-service)",
			},
			Validate: survey.Required,
		},
		{
			Name: "description",
			Prompt: &survey.Input{
				Message: "Project description (optional):",
				Help:    "A brief description of what your service does",
			},
		},
		{
			Name: "gomodule",
			Prompt: &survey.Input{
				Message: "Go module path:",
				Help:    "The import path for your Go module",
			},
		},
	}

	answers := struct {
		Name        string
		Description string
		GoModule    string
	}{}

	if err := survey.Ask(qs, &answers); err != nil {
		return err
	}

	cfg.Name = answers.Name
	cfg.Description = answers.Description
	
	if answers.GoModule != "" {
		cfg.GoModule = answers.GoModule
	} else {
		cfg.GoModule = fmt.Sprintf("github.com/example/%s", answers.Name)
	}

	// Set output directory
	cfg.OutputDir = answers.Name

	return nil
}

func askTierSelection(cfg *config.ProjectConfig) error {
	fmt.Println("\n🎯 Choose your project tier:")
	fmt.Println("Each tier builds upon the previous one with additional features and complexity.")

	tierOptions := []string{
		"basic - Simple health endpoints (~5 min deployment)",
		"intermediate - Production-ready with dependency checks (~15 min deployment)",
		"advanced - Full observability with OpenTelemetry (~30 min deployment)",
		"enterprise - Enterprise-grade with compliance features (~45 min deployment)",
	}

	var selectedTier string
	prompt := &survey.Select{
		Message: "Select your project tier:",
		Options: tierOptions,
		Help:    "Higher tiers include all features from lower tiers plus additional capabilities",
	}

	if err := survey.AskOne(prompt, &selectedTier); err != nil {
		return err
	}

	// Extract tier name from selection
	tierName := strings.Split(selectedTier, " ")[0]
	cfg.Tier = config.TemplateTier(tierName)

	return nil
}

func askFeatureSelection(cfg *config.ProjectConfig) error {
	fmt.Printf("\n⚙️  Configure features for %s tier:\n", cfg.Tier)

	// Define available features based on tier
	availableFeatures := getAvailableFeatures(cfg.Tier)
	
	if len(availableFeatures) == 0 {
		fmt.Println("No additional features to configure for this tier.")
		return nil
	}

	var selectedFeatures []string
	prompt := &survey.MultiSelect{
		Message: "Select additional features to enable:",
		Options: availableFeatures,
		Help:    "You can enable/disable these features. Some may be required for higher tiers.",
	}

	if err := survey.AskOne(prompt, &selectedFeatures); err != nil {
		return err
	}

	// Apply selected features
	applyFeatureSelections(cfg, selectedFeatures)

	return nil
}

func askAdvancedConfiguration(cfg *config.ProjectConfig) error {
	fmt.Println("\n🔧 Advanced configuration (optional):")

	var wantAdvanced bool
	prompt := &survey.Confirm{
		Message: "Would you like to configure advanced settings?",
		Help:    "This includes Kubernetes, deployment, and custom configurations",
		Default: false,
	}

	if err := survey.AskOne(prompt, &wantAdvanced); err != nil {
		return err
	}

	if !wantAdvanced {
		return nil
	}

	return askKubernetesConfiguration(cfg)
}

func askKubernetesConfiguration(cfg *config.ProjectConfig) error {
	if !cfg.Features.Kubernetes {
		return nil
	}

	fmt.Println("\n☸️  Kubernetes Configuration:")

	var qs = []*survey.Question{
		{
			Name: "namespace",
			Prompt: &survey.Input{
				Message: "Kubernetes namespace:",
				Default: "default",
				Help:    "The namespace where your service will be deployed",
			},
		},
		{
			Name: "servicemonitor",
			Prompt: &survey.Confirm{
				Message: "Enable Prometheus ServiceMonitor?",
				Default: true,
				Help:    "Automatically discovers and scrapes metrics for Prometheus",
			},
		},
		{
			Name: "ingress",
			Prompt: &survey.Confirm{
				Message: "Enable Ingress for external access?",
				Default: false,
				Help:    "Exposes your service externally through an ingress controller",
			},
		},
	}

	if cfg.Tier == "enterprise" {
		qs = append(qs, &survey.Question{
			Name: "rbac",
			Prompt: &survey.Confirm{
				Message: "Enable RBAC security policies?",
				Default: true,
				Help:    "Creates ServiceAccount, Role, and RoleBinding for enhanced security",
			},
		})
	}

	answers := struct {
		Namespace      string
		ServiceMonitor bool
		Ingress        bool
		RBAC           bool
	}{}

	if err := survey.Ask(qs, &answers); err != nil {
		return err
	}

	cfg.Kubernetes.Namespace = answers.Namespace
	cfg.Kubernetes.ServiceMonitor = answers.ServiceMonitor
	cfg.Kubernetes.Ingress.Enabled = answers.Ingress
	
	if cfg.Tier == "enterprise" {
		cfg.Security.RBAC = answers.RBAC
	}

	if answers.Ingress {
		var host string
		hostPrompt := &survey.Input{
			Message: "Ingress hostname:",
			Help:    "The hostname for external access (e.g., myservice.example.com)",
		}
		if err := survey.AskOne(hostPrompt, &host); err != nil {
			return err
		}
		cfg.Kubernetes.Ingress.Host = host
	}

	return nil
}

func getAvailableFeatures(tier config.TemplateTier) []string {
	var features []string

	switch tier {
	case "basic":
		features = []string{
			"typescript - Generate TypeScript client SDK",
			"docker - Include Docker configuration",
		}
	case "intermediate":
		features = []string{
			"typescript - Generate TypeScript client SDK",
			"docker - Include Docker configuration",
			"kubernetes - Include Kubernetes manifests",
			"server-timing - Add Server-Timing headers",
		}
	case "advanced":
		features = []string{
			"typescript - Generate TypeScript client SDK",
			"docker - Include Docker configuration",
			"kubernetes - Include Kubernetes manifests",
			"server-timing - Add Server-Timing headers",
			"opentelemetry - Full OpenTelemetry integration (recommended)",
			"cloudevents - CloudEvents for health notifications",
		}
	case "enterprise":
		features = []string{
			"typescript - Generate TypeScript client SDK",
			"docker - Include Docker configuration",
			"kubernetes - Include Kubernetes manifests (recommended)",
			"server-timing - Add Server-Timing headers",
			"opentelemetry - Full OpenTelemetry integration (included)",
			"cloudevents - CloudEvents for health notifications",
		}
	}

	return features
}

func applyFeatureSelections(cfg *config.ProjectConfig, selections []string) {
	for _, selection := range selections {
		feature := strings.Split(selection, " ")[0]
		
		switch feature {
		case "typescript":
			cfg.Features.TypeScript = true
		case "docker":
			cfg.Features.Docker = true
		case "kubernetes":
			cfg.Features.Kubernetes = true
		case "server-timing":
			cfg.Features.ServerTiming = true
		case "opentelemetry":
			cfg.Features.OpenTelemetry = true
		case "cloudevents":
			cfg.Features.CloudEvents = true
		}
	}
}

// GetSmartDefaults provides intelligent defaults based on project name and context
func GetSmartDefaults(projectName string) *config.ProjectConfig {
	cfg := &config.ProjectConfig{
		Name:      projectName,
		OutputDir: projectName,
		GoModule:  fmt.Sprintf("github.com/example/%s", projectName),
		Tier:      "basic", // Safe default
	}

	// Intelligent tier suggestions based on project name patterns
	name := strings.ToLower(projectName)
	
	if strings.Contains(name, "enterprise") || strings.Contains(name, "prod") {
		cfg.Tier = "enterprise"
	} else if strings.Contains(name, "monitor") || strings.Contains(name, "observ") {
		cfg.Tier = "advanced"
	} else if strings.Contains(name, "api") || strings.Contains(name, "service") {
		cfg.Tier = "intermediate"
	}

	// Smart feature defaults
	if strings.Contains(name, "k8s") || strings.Contains(name, "kube") {
		cfg.Features.Kubernetes = true
	}
	
	if strings.Contains(name, "client") || strings.Contains(name, "sdk") {
		cfg.Features.TypeScript = true
	}

	return cfg
}