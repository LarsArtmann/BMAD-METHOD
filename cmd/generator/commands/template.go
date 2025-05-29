package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// TemplateMetadata represents the metadata for a template tier
type TemplateMetadata struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Tier        string            `yaml:"tier"`
	Features    map[string]bool   `yaml:"features"`
	Version     string            `yaml:"version"`
}

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage template operations",
	Long:  `Manage template operations including listing, validating, and generating from static templates.`,
}

// listTemplatesCmd lists available templates
var listTemplatesCmd = &cobra.Command{
	Use:   "list",
	Short: "List available template tiers",
	Long:  `List all available template tiers with their descriptions and features.`,
	RunE:  runListTemplates,
}

// generateFromTemplateCmd generates a project from static template
var generateFromTemplateCmd = &cobra.Command{
	Use:   "from-static",
	Short: "Generate project from static template directory",
	Long:  `Generate a new project from a static template directory.`,
	RunE:  runGenerateFromTemplate,
}

// validateTemplatesCmd validates template directories
var validateTemplatesCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate template directories",
	Long:  `Validate that all template directories are properly structured and contain valid templates.`,
	RunE:  runValidateTemplates,
}

func init() {
	// Add template subcommands
	templateCmd.AddCommand(listTemplatesCmd)
	templateCmd.AddCommand(generateFromTemplateCmd)
	templateCmd.AddCommand(validateTemplatesCmd)

	// Add flags for generate-from-template
	generateFromTemplateCmd.Flags().StringP("name", "n", "", "Project name (required)")
	generateFromTemplateCmd.Flags().StringP("tier", "t", "basic", "Template tier (basic, intermediate, advanced, enterprise)")
	generateFromTemplateCmd.Flags().StringP("module", "m", "", "Go module path (required)")
	generateFromTemplateCmd.Flags().StringP("output", "o", "", "Output directory (default: project name)")
	generateFromTemplateCmd.Flags().StringP("description", "d", "", "Project description")
	generateFromTemplateCmd.MarkFlagRequired("name")
	generateFromTemplateCmd.MarkFlagRequired("module")

	// Add template command to root
	rootCmd.AddCommand(templateCmd)
}

func runListTemplates(cmd *cobra.Command, args []string) error {
	templatesDir := "templates"

	fmt.Println("📋 Available Template Tiers:")
	fmt.Println("=" + fmt.Sprintf("%*s", 50, ""))

	// Read template directories
	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		return fmt.Errorf("failed to read templates directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		tierName := entry.Name()
		metadataPath := filepath.Join(templatesDir, tierName, "template.yaml")

		// Read template metadata
		metadata, err := readTemplateMetadata(metadataPath)
		if err != nil {
			fmt.Printf("⚠️  %s: No metadata found\n", tierName)
			continue
		}

		fmt.Printf("\n🎯 **%s** (%s)\n", metadata.Name, metadata.Version)
		fmt.Printf("   %s\n", metadata.Description)
		fmt.Printf("   Features: ")

		var features []string
		for feature, enabled := range metadata.Features {
			if enabled {
				features = append(features, feature)
			}
		}

		if len(features) > 0 {
			for i, feature := range features {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Print(feature)
			}
		} else {
			fmt.Print("none")
		}
		fmt.Println()
	}

	fmt.Println("\n💡 Use 'template from-static --tier <tier>' to generate from a template")
	return nil
}

func runGenerateFromTemplate(cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")
	tier, _ := cmd.Flags().GetString("tier")
	module, _ := cmd.Flags().GetString("module")
	output, _ := cmd.Flags().GetString("output")
	description, _ := cmd.Flags().GetString("description")

	if output == "" {
		output = name
	}

	if description == "" {
		description = fmt.Sprintf("%s health endpoint service", name)
	}

	fmt.Printf("🚀 Generating project from %s template...\n", tier)

	// Validate template exists
	templateDir := filepath.Join("templates", tier)
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		return fmt.Errorf("template tier '%s' not found in templates directory", tier)
	}

	// Read template metadata
	metadataPath := filepath.Join(templateDir, "template.yaml")
	metadata, err := readTemplateMetadata(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to read template metadata: %w", err)
	}

	fmt.Printf("📋 Using template: %s (%s)\n", metadata.Description, metadata.Version)

	// Create template context
	context := map[string]interface{}{
		"Config": map[string]interface{}{
			"Name":        name,
			"Description": description,
			"GoModule":    module,
			"Tier":       tier,
			"Features":   metadata.Features,
		},
		"Version":   metadata.Version,
		"Timestamp": "2024-01-01T00:00:00Z", // TODO: Use actual timestamp
	}

	// Generate project from template
	if err := generateFromStaticTemplate(templateDir, output, context); err != nil {
		return fmt.Errorf("failed to generate project: %w", err)
	}

	fmt.Printf("✅ Successfully generated project from %s template!\n", tier)
	fmt.Printf("📁 Project created in: %s\n", output)

	return nil
}

func runValidateTemplates(cmd *cobra.Command, args []string) error {
	templatesDir := "templates"

	fmt.Println("🔍 Validating template directories...")

	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		return fmt.Errorf("failed to read templates directory: %w", err)
	}

	valid := true

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		tierName := entry.Name()
		templateDir := filepath.Join(templatesDir, tierName)

		fmt.Printf("\n📋 Validating %s template...\n", tierName)

		// Check for required files
		requiredFiles := []string{
			"template.yaml",
			"cmd/server/main.go",
			"internal/handlers/health.go",
			"go.mod.tmpl",
			"README.md.tmpl",
		}

		for _, file := range requiredFiles {
			filePath := filepath.Join(templateDir, file)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				fmt.Printf("❌ Missing required file: %s\n", file)
				valid = false
			} else {
				fmt.Printf("✅ Found: %s\n", file)
			}
		}

		// Validate metadata
		metadataPath := filepath.Join(templateDir, "template.yaml")
		if _, err := readTemplateMetadata(metadataPath); err != nil {
			fmt.Printf("❌ Invalid metadata: %v\n", err)
			valid = false
		} else {
			fmt.Printf("✅ Valid metadata\n")
		}
	}

	if valid {
		fmt.Println("\n🎉 All templates are valid!")
	} else {
		fmt.Println("\n❌ Some templates have issues")
		return fmt.Errorf("template validation failed")
	}

	return nil
}

func readTemplateMetadata(path string) (*TemplateMetadata, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var metadata TemplateMetadata
	if err := yaml.Unmarshal(data, &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

func generateFromStaticTemplate(templateDir, outputDir string, context map[string]interface{}) error {
	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Walk through template directory
	return filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip template.yaml metadata file
		if filepath.Base(path) == "template.yaml" {
			return nil
		}

		// Calculate relative path
		relPath, err := filepath.Rel(templateDir, path)
		if err != nil {
			return err
		}

		// Calculate output path
		outputPath := filepath.Join(outputDir, relPath)

		if info.IsDir() {
			// Create directory
			return os.MkdirAll(outputPath, info.Mode())
		}

		// Process file
		return processTemplateFile(path, outputPath, context)
	})
}

func processTemplateFile(inputPath, outputPath string, context map[string]interface{}) error {
	// Read input file
	content, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	// Create output directory if needed
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	// Remove .tmpl extension from output path
	if filepath.Ext(outputPath) == ".tmpl" {
		outputPath = outputPath[:len(outputPath)-5]
	}

	// Check if file needs template processing
	if filepath.Ext(inputPath) == ".tmpl" ||
	   filepath.Base(inputPath) == "go.mod" ||
	   filepath.Base(inputPath) == "README.md" ||
	   filepath.Ext(inputPath) == ".go" ||
	   filepath.Ext(inputPath) == ".yaml" ||
	   filepath.Ext(inputPath) == ".yml" ||
	   filepath.Ext(inputPath) == ".json" ||
	   filepath.Ext(inputPath) == ".ts" ||
	   filepath.Ext(inputPath) == ".sh" {

		// Process as template
		tmpl, err := template.New(filepath.Base(inputPath)).Parse(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", inputPath, err)
		}

		// Create output file
		outputFile, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer outputFile.Close()

		// Execute template
		return tmpl.Execute(outputFile, context)
	} else {
		// Copy file as-is
		return os.WriteFile(outputPath, content, 0644)
	}
}
