package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/LarsArtmann/BMAD-METHOD/pkg/typespec"
)

var (
	schemasPath string
	outputPath  string
	emitters    []string
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate TypeSpec schemas and generate outputs",
	Long: `Validate TypeSpec schemas for health endpoints and optionally generate outputs.

This command validates all TypeSpec schema definitions and can generate:
- JSON Schema files for validation
- OpenAPI v3 specifications for documentation
- Code generation artifacts

Examples:
  # Validate all schemas in the default location
  template-health-endpoint validate

  # Validate schemas in a specific directory
  template-health-endpoint validate --schemas ./my-schemas

  # Validate and generate OpenAPI specifications
  template-health-endpoint validate --emit openapi3

  # Validate and generate both JSON Schema and OpenAPI
  template-health-endpoint validate --emit json-schema,openapi3 --output ./generated`,
	RunE: runValidate,
}

func init() {
	rootCmd.AddCommand(validateCmd)

	// Flags
	validateCmd.Flags().StringVarP(&schemasPath, "schemas", "s", "pkg/schemas", "path to TypeSpec schemas directory")
	validateCmd.Flags().StringVarP(&outputPath, "output", "o", "tsp-output", "output directory for generated files")
	validateCmd.Flags().StringSliceVarP(&emitters, "emit", "e", []string{}, "comma-separated list of emitters (json-schema,openapi3)")
}

func runValidate(cmd *cobra.Command, args []string) error {
	// Check if TypeSpec schemas exist
	if _, err := os.Stat(schemasPath); os.IsNotExist(err) {
		return fmt.Errorf("schemas directory not found: %s", schemasPath)
	}

	// Create TypeSpec validator
	validator, err := typespec.NewValidator(schemasPath)
	if err != nil {
		return fmt.Errorf("failed to create TypeSpec validator: %w", err)
	}

	// Validate schemas
	fmt.Println("🔍 Validating TypeSpec schemas...")
	
	result, err := validator.Validate()
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Show validation results
	if err := showValidationResults(result); err != nil {
		return fmt.Errorf("failed to show validation results: %w", err)
	}

	// If validation failed, don't proceed with generation
	if !result.Success {
		return fmt.Errorf("schema validation failed with %d errors", len(result.Errors))
	}

	// Generate outputs if emitters are specified
	if len(emitters) > 0 {
		fmt.Println("\n📊 Generating outputs...")
		
		for _, emitter := range emitters {
			if err := generateOutput(validator, emitter); err != nil {
				return fmt.Errorf("failed to generate %s output: %w", emitter, err)
			}
		}
	}

	fmt.Println("\n✅ Validation completed successfully!")
	return nil
}

func showValidationResults(result *typespec.ValidationResult) error {
	fmt.Printf("\n📋 Validation Results:\n")
	fmt.Printf("  Files validated: %d\n", result.FilesValidated)
	fmt.Printf("  Schemas found: %d\n", result.SchemasFound)
	fmt.Printf("  Errors: %d\n", len(result.Errors))
	fmt.Printf("  Warnings: %d\n", len(result.Warnings))

	// Show errors
	if len(result.Errors) > 0 {
		fmt.Println("\n❌ Errors:")
		for _, err := range result.Errors {
			fmt.Printf("  - %s:%d:%d - %s\n", err.File, err.Line, err.Column, err.Message)
		}
	}

	// Show warnings
	if len(result.Warnings) > 0 {
		fmt.Println("\n⚠️  Warnings:")
		for _, warning := range result.Warnings {
			fmt.Printf("  - %s:%d:%d - %s\n", warning.File, warning.Line, warning.Column, warning.Message)
		}
	}

	// Show schema summary
	if len(result.Schemas) > 0 {
		fmt.Println("\n📐 Schemas found:")
		for _, schema := range result.Schemas {
			fmt.Printf("  - %s (%s)\n", schema.Name, schema.Type)
		}
	}

	return nil
}

func generateOutput(validator *typespec.Validator, emitter string) error {
	switch strings.ToLower(emitter) {
	case "json-schema", "jsonschema":
		return generateJSONSchema(validator)
	case "openapi3", "openapi":
		return generateOpenAPI(validator)
	default:
		return fmt.Errorf("unknown emitter: %s (supported: json-schema, openapi3)", emitter)
	}
}

func generateJSONSchema(validator *typespec.Validator) error {
	fmt.Println("  📄 Generating JSON Schema...")
	
	output, err := validator.GenerateJSONSchema()
	if err != nil {
		return err
	}

	// Create output directory
	jsonSchemaDir := filepath.Join(outputPath, "json-schema")
	if err := os.MkdirAll(jsonSchemaDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write JSON Schema files
	for filename, content := range output.Files {
		outputFile := filepath.Join(jsonSchemaDir, filename)
		if err := os.WriteFile(outputFile, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write JSON Schema file %s: %w", filename, err)
		}
		
		if viper.GetBool("verbose") {
			fmt.Printf("    ✓ Generated: %s\n", outputFile)
		}
	}

	fmt.Printf("    ✅ JSON Schema generation complete (%d files)\n", len(output.Files))
	return nil
}

func generateOpenAPI(validator *typespec.Validator) error {
	fmt.Println("  📋 Generating OpenAPI v3...")
	
	output, err := validator.GenerateOpenAPI()
	if err != nil {
		return err
	}

	// Create output directory
	openAPIDir := filepath.Join(outputPath, "openapi")
	if err := os.MkdirAll(openAPIDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write OpenAPI files
	for filename, content := range output.Files {
		outputFile := filepath.Join(openAPIDir, filename)
		if err := os.WriteFile(outputFile, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write OpenAPI file %s: %w", filename, err)
		}
		
		if viper.GetBool("verbose") {
			fmt.Printf("    ✓ Generated: %s\n", outputFile)
		}
	}

	fmt.Printf("    ✅ OpenAPI generation complete (%d files)\n", len(output.Files))
	return nil
}
