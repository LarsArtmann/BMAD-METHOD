package typespec

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Validator handles TypeSpec schema validation and code generation
type Validator struct {
	schemasPath string
	tspPath     string
}

// ValidationResult contains the results of schema validation
type ValidationResult struct {
	Success        bool
	FilesValidated int
	SchemasFound   int
	Errors         []ValidationError
	Warnings       []ValidationWarning
	Schemas        []SchemaInfo
}

// ValidationError represents a validation error
type ValidationError struct {
	File    string
	Line    int
	Column  int
	Message string
	Code    string
}

// ValidationWarning represents a validation warning
type ValidationWarning struct {
	File    string
	Line    int
	Column  int
	Message string
	Code    string
}

// SchemaInfo contains information about a discovered schema
type SchemaInfo struct {
	Name      string
	Type      string
	File      string
	Namespace string
}

// GenerationOutput contains the results of code generation
type GenerationOutput struct {
	Files map[string]string
	Stats GenerationStats
}

// GenerationStats contains statistics about code generation
type GenerationStats struct {
	FilesGenerated int
	ModelsGenerated int
	InterfacesGenerated int
}

// NewValidator creates a new TypeSpec validator
func NewValidator(schemasPath string) (*Validator, error) {
	// Check if TypeSpec compiler is available
	tspPath, err := exec.LookPath("tsp")
	if err != nil {
		// Try npx tsp as fallback
		npxPath, npxErr := exec.LookPath("npx")
		if npxErr != nil {
			return nil, fmt.Errorf("TypeSpec compiler not found. Please install @typespec/compiler: %w", err)
		}
		tspPath = npxPath
	}

	// Verify schemas path exists
	if _, err := os.Stat(schemasPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("schemas path does not exist: %s", schemasPath)
	}

	return &Validator{
		schemasPath: schemasPath,
		tspPath:     tspPath,
	}, nil
}

// Validate validates all TypeSpec schemas in the configured path
func (v *Validator) Validate() (*ValidationResult, error) {
	result := &ValidationResult{
		Errors:   make([]ValidationError, 0),
		Warnings: make([]ValidationWarning, 0),
		Schemas:  make([]SchemaInfo, 0),
	}

	// Find all TypeSpec files
	tspFiles, err := v.findTypeSpecFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to find TypeSpec files: %w", err)
	}

	result.FilesValidated = len(tspFiles)

	// Validate each file
	for _, file := range tspFiles {
		if err := v.validateFile(file, result); err != nil {
			return nil, fmt.Errorf("failed to validate file %s: %w", file, err)
		}
	}

	// Set success status
	result.Success = len(result.Errors) == 0

	return result, nil
}

// GenerateJSONSchema generates JSON Schema from TypeSpec definitions
func (v *Validator) GenerateJSONSchema() (*GenerationOutput, error) {
	output := &GenerationOutput{
		Files: make(map[string]string),
	}

	// Run TypeSpec compiler with JSON Schema emitter
	cmd := v.createCompileCommand([]string{"@typespec/json-schema"})
	
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to generate JSON Schema: %w", err)
	}

	// Read generated files
	jsonSchemaDir := filepath.Join("tsp-output", "@typespec", "json-schema")
	if err := v.readGeneratedFiles(jsonSchemaDir, output); err != nil {
		return nil, fmt.Errorf("failed to read generated JSON Schema files: %w", err)
	}

	return output, nil
}

// GenerateOpenAPI generates OpenAPI v3 specification from TypeSpec definitions
func (v *Validator) GenerateOpenAPI() (*GenerationOutput, error) {
	output := &GenerationOutput{
		Files: make(map[string]string),
	}

	// Run TypeSpec compiler with OpenAPI emitter
	cmd := v.createCompileCommand([]string{"@typespec/openapi3"})
	
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to generate OpenAPI: %w", err)
	}

	// Read generated files
	openAPIDir := filepath.Join("tsp-output", "@typespec", "openapi3")
	if err := v.readGeneratedFiles(openAPIDir, output); err != nil {
		return nil, fmt.Errorf("failed to read generated OpenAPI files: %w", err)
	}

	return output, nil
}

// findTypeSpecFiles finds all .tsp files in the schemas directory
func (v *Validator) findTypeSpecFiles() ([]string, error) {
	var files []string

	err := filepath.Walk(v.schemasPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".tsp") {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// validateFile validates a single TypeSpec file
func (v *Validator) validateFile(file string, result *ValidationResult) error {
	// For now, we'll do basic file validation
	// In a real implementation, this would parse the TypeSpec file and validate syntax
	
	content, err := os.ReadFile(file)
	if err != nil {
		result.Errors = append(result.Errors, ValidationError{
			File:    file,
			Line:    1,
			Column:  1,
			Message: fmt.Sprintf("Failed to read file: %v", err),
			Code:    "file-read-error",
		})
		return nil
	}

	// Basic syntax checks
	if len(content) == 0 {
		result.Warnings = append(result.Warnings, ValidationWarning{
			File:    file,
			Line:    1,
			Column:  1,
			Message: "File is empty",
			Code:    "empty-file",
		})
	}

	// Look for basic TypeSpec constructs
	contentStr := string(content)
	
	// Count models
	modelCount := strings.Count(contentStr, "model ")
	if modelCount > 0 {
		result.SchemasFound += modelCount
		// Add schema info (simplified)
		result.Schemas = append(result.Schemas, SchemaInfo{
			Name:      filepath.Base(file),
			Type:      "model",
			File:      file,
			Namespace: extractNamespace(contentStr),
		})
	}

	// Count interfaces
	interfaceCount := strings.Count(contentStr, "interface ")
	if interfaceCount > 0 {
		result.SchemasFound += interfaceCount
		result.Schemas = append(result.Schemas, SchemaInfo{
			Name:      filepath.Base(file),
			Type:      "interface",
			File:      file,
			Namespace: extractNamespace(contentStr),
		})
	}

	return nil
}

// createCompileCommand creates a TypeSpec compile command
func (v *Validator) createCompileCommand(emitters []string) *exec.Cmd {
	args := []string{"compile", "main.tsp"}
	
	if len(emitters) > 0 {
		for _, emitter := range emitters {
			args = append(args, "--emit", emitter)
		}
	} else {
		args = append(args, "--no-emit")
	}

	if strings.Contains(v.tspPath, "npx") {
		return exec.Command(v.tspPath, append([]string{"tsp"}, args...)...)
	}
	
	return exec.Command(v.tspPath, args...)
}

// readGeneratedFiles reads all files from a generated output directory
func (v *Validator) readGeneratedFiles(dir string, output *GenerationOutput) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil // No files generated
	}

	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Use relative path as key
		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			relPath = filepath.Base(path)
		}

		output.Files[relPath] = string(content)
		output.Stats.FilesGenerated++

		return nil
	})
}

// extractNamespace extracts the namespace from TypeSpec content (simplified)
func extractNamespace(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "namespace ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return strings.TrimSuffix(parts[1], ";")
			}
		}
	}
	return "default"
}
