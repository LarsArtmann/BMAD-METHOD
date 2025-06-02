package typespec

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/LarsArtmann/BMAD-METHOD/pkg/config"
)

// TypeSpecGenerator handles TypeSpec code generation
type TypeSpecGenerator struct {
	config      *config.ProjectConfig
	tspExecutor *TSPExecutor
}

// TSPExecutor handles TypeSpec compiler execution
type TSPExecutor struct {
	binaryPath   string
	timeout      time.Duration
	workingDir   string
	outputDir    string
}

// GenerationTarget represents a code generation target
type GenerationTarget struct {
	Language    string            `json:"language"`
	Emitter     string            `json:"emitter"`
	OutputDir   string            `json:"output_dir"`
	Options     map[string]string `json:"options"`
	Enabled     bool              `json:"enabled"`
}

// TypeSpecConfig holds TypeSpec generation configuration
type TypeSpecConfig struct {
	SchemaPath      string              `json:"schema_path"`
	OutputDir       string              `json:"output_dir"`
	Targets         []GenerationTarget  `json:"targets"`
	GlobalOptions   map[string]string   `json:"global_options"`
	CustomEmitters  []string            `json:"custom_emitters"`
}

// NewTypeSpecGenerator creates a new TypeSpec generator
func NewTypeSpecGenerator(cfg *config.ProjectConfig) (*TypeSpecGenerator, error) {
	executor, err := NewTSPExecutor()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize TSP executor: %w", err)
	}

	return &TypeSpecGenerator{
		config:      cfg,
		tspExecutor: executor,
	}, nil
}

// NewTSPExecutor creates a new TSP executor
func NewTSPExecutor() (*TSPExecutor, error) {
	// Try to find tsp binary
	binaryPath, err := exec.LookPath("tsp")
	if err != nil {
		return nil, fmt.Errorf("TypeSpec compiler (tsp) not found in PATH: %w", err)
	}

	return &TSPExecutor{
		binaryPath: binaryPath,
		timeout:    5 * time.Minute,
	}, nil
}

// GenerateAll generates code for all configured targets
func (tsg *TypeSpecGenerator) GenerateAll(ctx context.Context, tspConfig TypeSpecConfig) error {
	// Set working directory
	tsg.tspExecutor.workingDir = filepath.Dir(tspConfig.SchemaPath)
	tsg.tspExecutor.outputDir = tspConfig.OutputDir

	// Generate for each target
	for _, target := range tspConfig.Targets {
		if !target.Enabled {
			continue
		}

		err := tsg.generateForTarget(ctx, tspConfig.SchemaPath, target)
		if err != nil {
			return fmt.Errorf("failed to generate for target %s: %w", target.Language, err)
		}
	}

	return nil
}

// generateForTarget generates code for a specific target
func (tsg *TypeSpecGenerator) generateForTarget(ctx context.Context, schemaPath string, target GenerationTarget) error {
	args := []string{"compile", schemaPath}

	// Add emitter
	if target.Emitter != "" {
		args = append(args, "--emit", target.Emitter)
	}

	// Add output directory
	if target.OutputDir != "" {
		args = append(args, "--output-dir", target.OutputDir)
	}

	// Add target-specific options
	for key, value := range target.Options {
		args = append(args, fmt.Sprintf("--%s", key), value)
	}

	return tsg.tspExecutor.Execute(ctx, args...)
}

// GenerateOpenAPI generates OpenAPI v3 specification
func (tsg *TypeSpecGenerator) GenerateOpenAPI(ctx context.Context, schemaPath, outputDir string) error {
	return tsg.tspExecutor.Execute(ctx, "compile", schemaPath, 
		"--emit", "@typespec/openapi3", 
		"--output-dir", outputDir)
}

// GenerateJSONSchema generates JSON Schema
func (tsg *TypeSpecGenerator) GenerateJSONSchema(ctx context.Context, schemaPath, outputDir string) error {
	return tsg.tspExecutor.Execute(ctx, "compile", schemaPath,
		"--emit", "@typespec/json-schema",
		"--output-dir", outputDir)
}

// GenerateTypeScriptTypes generates TypeScript type definitions
func (tsg *TypeSpecGenerator) GenerateTypeScriptTypes(ctx context.Context, schemaPath, outputDir string) error {
	// This would use a custom TypeScript emitter
	return tsg.tspExecutor.Execute(ctx, "compile", schemaPath,
		"--emit", "@typespec/typescript",
		"--output-dir", outputDir)
}

// GenerateGoTypes generates Go type definitions
func (tsg *TypeSpecGenerator) GenerateGoTypes(ctx context.Context, schemaPath, outputDir string, packageName string) error {
	// This would use a custom Go emitter
	return tsg.tspExecutor.Execute(ctx, "compile", schemaPath,
		"--emit", "@typespec/go",
		"--output-dir", outputDir,
		"--package-name", packageName)
}

// GeneratePythonTypes generates Python type definitions
func (tsg *TypeSpecGenerator) GeneratePythonTypes(ctx context.Context, schemaPath, outputDir string) error {
	// This would use a custom Python emitter
	return tsg.tspExecutor.Execute(ctx, "compile", schemaPath,
		"--emit", "@typespec/python",
		"--output-dir", outputDir)
}

// Execute runs the TypeSpec compiler with given arguments
func (tsp *TSPExecutor) Execute(ctx context.Context, args ...string) error {
	// Create context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, tsp.timeout)
	defer cancel()

	// Create command
	cmd := exec.CommandContext(timeoutCtx, tsp.binaryPath, args...)
	
	// Set working directory if specified
	if tsp.workingDir != "" {
		cmd.Dir = tsp.workingDir
	}

	// Set environment variables
	cmd.Env = os.Environ()

	// Run command and capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("TypeSpec compilation failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// ValidateSchema validates a TypeSpec schema
func (tsg *TypeSpecGenerator) ValidateSchema(ctx context.Context, schemaPath string) error {
	return tsg.tspExecutor.Execute(ctx, "compile", schemaPath, "--no-emit")
}

// GetDefaultConfig returns default TypeSpec configuration
func (tsg *TypeSpecGenerator) GetDefaultConfig() TypeSpecConfig {
	baseOutputDir := filepath.Join(tsg.config.OutputDir, "generated")
	
	targets := []GenerationTarget{
		{
			Language:  "openapi",
			Emitter:   "@typespec/openapi3",
			OutputDir: filepath.Join(baseOutputDir, "openapi"),
			Enabled:   true,
		},
		{
			Language:  "json-schema",
			Emitter:   "@typespec/json-schema",
			OutputDir: filepath.Join(baseOutputDir, "json-schema"),
			Enabled:   true,
		},
	}

	// Add TypeScript target if enabled
	if tsg.config.Features.TypeScript {
		targets = append(targets, GenerationTarget{
			Language:  "typescript",
			Emitter:   "@typespec/typescript",
			OutputDir: filepath.Join(baseOutputDir, "typescript"),
			Options: map[string]string{
				"packageName": tsg.config.Name + "-types",
			},
			Enabled: true,
		})
	}

	return TypeSpecConfig{
		SchemaPath: "main.tsp",
		OutputDir:  baseOutputDir,
		Targets:    targets,
		GlobalOptions: map[string]string{
			"no-emit": "false",
		},
	}
}

// InstallEmitters installs required TypeSpec emitters
func (tsg *TypeSpecGenerator) InstallEmitters(ctx context.Context, emitters []string) error {
	for _, emitter := range emitters {
		err := tsg.installEmitter(ctx, emitter)
		if err != nil {
			return fmt.Errorf("failed to install emitter %s: %w", emitter, err)
		}
	}
	return nil
}

// installEmitter installs a single TypeSpec emitter
func (tsg *TypeSpecGenerator) installEmitter(ctx context.Context, emitter string) error {
	// Use npm to install the emitter
	cmd := exec.CommandContext(ctx, "npm", "install", "-g", emitter)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("npm install failed: %w\nOutput: %s", err, string(output))
	}
	return nil
}

// GenerateClientSDK generates complete client SDK
func (tsg *TypeSpecGenerator) GenerateClientSDK(ctx context.Context, languages []string) error {
	baseConfig := tsg.GetDefaultConfig()
	
	for _, lang := range languages {
		switch strings.ToLower(lang) {
		case "typescript", "ts":
			err := tsg.GenerateTypeScriptTypes(ctx, baseConfig.SchemaPath, 
				filepath.Join(baseConfig.OutputDir, "typescript"))
			if err != nil {
				return err
			}
			
		case "go":
			err := tsg.GenerateGoTypes(ctx, baseConfig.SchemaPath,
				filepath.Join(baseConfig.OutputDir, "go"), tsg.config.Name)
			if err != nil {
				return err
			}
			
		case "python", "py":
			err := tsg.GeneratePythonTypes(ctx, baseConfig.SchemaPath,
				filepath.Join(baseConfig.OutputDir, "python"))
			if err != nil {
				return err
			}
		}
	}
	
	return nil
}

// MultiLanguageConfig holds configuration for multi-language generation
type MultiLanguageConfig struct {
	Languages    []string          `json:"languages"`
	OutputBase   string            `json:"output_base"`
	PackageNames map[string]string `json:"package_names"`
	Options      map[string]map[string]string `json:"options"`
}

// GenerateMultiLanguage generates code for multiple languages
func (tsg *TypeSpecGenerator) GenerateMultiLanguage(ctx context.Context, config MultiLanguageConfig) error {
	for _, lang := range config.Languages {
		outputDir := filepath.Join(config.OutputBase, lang)
		
		var target GenerationTarget
		switch strings.ToLower(lang) {
		case "typescript":
			target = GenerationTarget{
				Language:  "typescript",
				Emitter:   "@typespec/typescript",
				OutputDir: outputDir,
				Options:   getLanguageOptions(config, lang),
				Enabled:   true,
			}
		case "go":
			target = GenerationTarget{
				Language:  "go",
				Emitter:   "@typespec/go",
				OutputDir: outputDir,
				Options:   getLanguageOptions(config, lang),
				Enabled:   true,
			}
		case "python":
			target = GenerationTarget{
				Language:  "python",
				Emitter:   "@typespec/python",
				OutputDir: outputDir,
				Options:   getLanguageOptions(config, lang),
				Enabled:   true,
			}
		case "csharp":
			target = GenerationTarget{
				Language:  "csharp",
				Emitter:   "@typespec/csharp",
				OutputDir: outputDir,
				Options:   getLanguageOptions(config, lang),
				Enabled:   true,
			}
		case "java":
			target = GenerationTarget{
				Language:  "java",
				Emitter:   "@typespec/java",
				OutputDir: outputDir,
				Options:   getLanguageOptions(config, lang),
				Enabled:   true,
			}
		default:
			continue // Skip unsupported languages
		}
		
		err := tsg.generateForTarget(ctx, "main.tsp", target)
		if err != nil {
			return fmt.Errorf("failed to generate %s code: %w", lang, err)
		}
	}
	
	return nil
}

// getLanguageOptions gets language-specific options
func getLanguageOptions(config MultiLanguageConfig, lang string) map[string]string {
	options := make(map[string]string)
	
	// Add package name if specified
	if packageName, exists := config.PackageNames[lang]; exists {
		options["packageName"] = packageName
	}
	
	// Add language-specific options
	if langOptions, exists := config.Options[lang]; exists {
		for k, v := range langOptions {
			options[k] = v
		}
	}
	
	return options
}

// GetRequiredEmitters returns list of required emitters for the configuration
func (tsg *TypeSpecGenerator) GetRequiredEmitters() []string {
	emitters := []string{
		"@typespec/openapi3",
		"@typespec/json-schema",
	}
	
	if tsg.config.Features.TypeScript {
		emitters = append(emitters, "@typespec/typescript")
	}
	
	// Add other language emitters based on configuration
	// These would be configurable based on user preferences
	
	return emitters
}

// WatchAndRegenerate watches for schema changes and regenerates code
func (tsg *TypeSpecGenerator) WatchAndRegenerate(ctx context.Context, schemaPath string) error {
	// This would implement file watching and automatic regeneration
	// For now, it's a placeholder for future implementation
	return fmt.Errorf("watch mode not yet implemented")
}