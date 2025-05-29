# CLI Command Architecture and User Experience

## Prompt Name
**CLI Command Architecture and User Experience**

## Context
Building CLI tools that are both powerful and user-friendly requires careful design of command structure, error handling, and user feedback. This prompt provides guidelines for creating professional-grade CLI tools with excellent user experience.

## Objective
Design and implement CLI tools with hierarchical command structure, comprehensive error handling, helpful user feedback, and professional user experience that scales from simple to complex operations.

## Key Principles

### 1. Hierarchical Command Structure
Design commands that group related functionality logically:

```bash
# Root command
tool-name

# Primary commands
tool-name generate          # Core functionality
tool-name validate          # Validation operations
tool-name template          # Template management

# Subcommands
tool-name template list     # List available templates
tool-name template validate # Validate template integrity
tool-name template from-static  # Generate from static template

# Advanced commands
tool-name update            # Update existing projects
tool-name migrate           # Migrate between versions/tiers
tool-name customize         # Interactive customization
```

### 2. User Experience Standards

#### **Progress Indicators**
Show users what's happening during long operations:
```
🚀 Generating project from basic template...
📋 Using template: Basic tier health endpoint template (1.0.0)
✅ Successfully generated project from basic template!
📁 Project created in: my-project
```

#### **Clear Error Messages**
Provide actionable error information:
```
❌ Template tier 'invalid' not found in templates directory
💡 Available tiers: basic, intermediate, advanced, enterprise
🔧 Use 'tool-name template list' to see all available templates
```

#### **Helpful Success Feedback**
Confirm what was accomplished and suggest next steps:
```
✅ Successfully generated project from basic template!
📁 Project created in: my-project

🚀 Next steps:
  cd my-project
  go mod tidy
  go run cmd/server/main.go
```

### 3. Command Implementation Pattern

```go
// Command structure using Cobra
var generateCmd = &cobra.Command{
    Use:   "generate",
    Short: "Generate a new project from template",
    Long:  `Generate a new health endpoint project from the specified template tier.`,
    RunE:  runGenerate,
}

func init() {
    // Add flags with validation
    generateCmd.Flags().StringP("name", "n", "", "Project name (required)")
    generateCmd.Flags().StringP("tier", "t", "basic", "Template tier")
    generateCmd.Flags().StringP("module", "m", "", "Go module path (required)")
    generateCmd.Flags().StringP("output", "o", "", "Output directory")
    generateCmd.Flags().Bool("dry-run", false, "Preview generation without creating files")
    
    // Mark required flags
    generateCmd.MarkFlagRequired("name")
    generateCmd.MarkFlagRequired("module")
    
    // Add to root command
    rootCmd.AddCommand(generateCmd)
}

func runGenerate(cmd *cobra.Command, args []string) error {
    // Extract flags with validation
    name, _ := cmd.Flags().GetString("name")
    tier, _ := cmd.Flags().GetString("tier")
    module, _ := cmd.Flags().GetString("module")
    output, _ := cmd.Flags().GetString("output")
    dryRun, _ := cmd.Flags().GetBool("dry-run")
    
    // Validate inputs
    if err := validateInputs(name, tier, module); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    
    // Show progress
    fmt.Printf("🚀 Generating project from %s template...\n", tier)
    
    // Perform operation
    if err := generateProject(name, tier, module, output, dryRun); err != nil {
        return fmt.Errorf("generation failed: %w", err)
    }
    
    // Show success
    fmt.Printf("✅ Successfully generated project from %s template!\n", tier)
    fmt.Printf("📁 Project created in: %s\n", output)
    
    return nil
}
```

### 4. Error Handling Best Practices

#### **Input Validation**
```go
func validateInputs(name, tier, module string) error {
    // Validate project name
    if name == "" {
        return fmt.Errorf("project name is required")
    }
    if !isValidProjectName(name) {
        return fmt.Errorf("invalid project name: %s (must be alphanumeric with hyphens)", name)
    }
    
    // Validate tier
    validTiers := []string{"basic", "intermediate", "advanced", "enterprise"}
    if !contains(validTiers, tier) {
        return fmt.Errorf("invalid tier: %s (available: %s)", tier, strings.Join(validTiers, ", "))
    }
    
    // Validate module path
    if module == "" {
        return fmt.Errorf("Go module path is required")
    }
    if !isValidModulePath(module) {
        return fmt.Errorf("invalid Go module path: %s", module)
    }
    
    return nil
}
```

#### **Contextual Error Messages**
```go
func handleTemplateError(err error, tier string) error {
    if os.IsNotExist(err) {
        return fmt.Errorf(`template tier '%s' not found
💡 Available tiers: basic, intermediate, advanced, enterprise
🔧 Use 'template-health-endpoint template list' to see all templates`, tier)
    }
    
    if strings.Contains(err.Error(), "permission denied") {
        return fmt.Errorf(`permission denied accessing template
🔧 Try running with appropriate permissions or check file ownership`)
    }
    
    return fmt.Errorf("template processing failed: %w", err)
}
```

### 5. Advanced CLI Features

#### **Dry Run Support**
```go
func generateProject(name, tier, module, output string, dryRun bool) error {
    if dryRun {
        fmt.Println("🔍 Dry run mode - showing what would be generated:")
        return previewGeneration(name, tier, module, output)
    }
    
    return performGeneration(name, tier, module, output)
}
```

#### **Interactive Mode**
```go
func runInteractiveGeneration() error {
    // Prompt for project name
    name, err := promptForInput("Project name", "")
    if err != nil {
        return err
    }
    
    // Prompt for tier selection
    tier, err := promptForSelection("Template tier", []string{"basic", "intermediate", "advanced", "enterprise"})
    if err != nil {
        return err
    }
    
    // Continue with generation...
    return generateProject(name, tier, module, output, false)
}
```

#### **Configuration File Support**
```go
func loadConfigFile(configPath string) (*Config, error) {
    if configPath == "" {
        return &Config{}, nil // Use defaults
    }
    
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }
    
    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse config file: %w", err)
    }
    
    return &config, nil
}
```

### 6. Help and Documentation

#### **Comprehensive Help Text**
```go
var rootCmd = &cobra.Command{
    Use:   "template-health-endpoint",
    Short: "Generate health endpoint projects from templates",
    Long: `Template Health Endpoint Generator

Generate production-ready health endpoint projects using TypeSpec-first API definitions,
with support for multiple complexity tiers and both static templates and CLI generation.

Examples:
  # Generate basic project
  template-health-endpoint generate --name my-service --tier basic --module github.com/org/my-service
  
  # Generate from static template
  template-health-endpoint template from-static --name my-service --tier advanced --module github.com/org/my-service
  
  # List available templates
  template-health-endpoint template list
  
  # Validate templates
  template-health-endpoint template validate`,
}
```

#### **Command Examples**
```go
var generateCmd = &cobra.Command{
    Use:   "generate",
    Short: "Generate a new project from template",
    Example: `  # Basic project generation
  template-health-endpoint generate --name user-service --tier basic --module github.com/myorg/user-service
  
  # Advanced project with custom output directory
  template-health-endpoint generate --name payment-service --tier advanced --module github.com/myorg/payment-service --output ./services/payment
  
  # Dry run to preview generation
  template-health-endpoint generate --name test-service --tier basic --module github.com/test/service --dry-run`,
}
```

### 7. Testing CLI Commands

#### **Command Testing**
```go
func TestGenerateCommand(t *testing.T) {
    tests := []struct {
        name    string
        args    []string
        wantErr bool
    }{
        {
            name: "valid basic generation",
            args: []string{"generate", "--name", "test-service", "--tier", "basic", "--module", "github.com/test/service"},
            wantErr: false,
        },
        {
            name: "missing required flag",
            args: []string{"generate", "--name", "test-service"},
            wantErr: true,
        },
        {
            name: "invalid tier",
            args: []string{"generate", "--name", "test-service", "--tier", "invalid", "--module", "github.com/test/service"},
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            cmd := rootCmd
            cmd.SetArgs(tt.args)
            err := cmd.Execute()
            
            if (err != nil) != tt.wantErr {
                t.Errorf("command execution error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Success Criteria
- [ ] Commands are logically organized in hierarchical structure
- [ ] Error messages are clear and actionable
- [ ] Success feedback includes next steps
- [ ] Help text is comprehensive and includes examples
- [ ] CLI supports both interactive and non-interactive modes
- [ ] Dry run options are available for destructive operations
- [ ] Performance is acceptable for all operations
- [ ] CLI follows established conventions and best practices

## Related Patterns
- Command-line interface design
- User experience design
- Error handling and validation
- Configuration management
- Testing strategies for CLI tools
