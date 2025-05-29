# Issue Alignment and CLI Architecture Guidelines

## Overview
Guidelines for ensuring project alignment with original requirements and building excellent CLI tools with hierarchical command structures, based on learnings from the template-health-endpoint project realignment.

## Core Principles

### 1. Requirement Validation and Alignment
**Guideline**: Establish systematic requirement validation throughout the project lifecycle.

**Implementation**:
- **Read Complete Issues**: Always read the entire original issue including all comments and related issues
- **Use GitHub API**: Programmatically fetch issues to ensure you have the complete context
- **Periodic Validation**: Regularly validate current work against original requirements
- **Document Gaps**: Clearly document any deviations from original requirements with justification

**Process**:
```bash
# Use GitHub API to fetch original issue
gh api repos/owner/repo/issues/123 > original-issue.json
gh api repos/owner/repo/issues/123/comments > issue-comments.json

# Review and document requirements
echo "Original Requirements:" > requirements-analysis.md
echo "Current Implementation:" >> requirements-analysis.md
echo "Gaps Identified:" >> requirements-analysis.md
echo "Alignment Strategy:" >> requirements-analysis.md
```

**Validation Checklist**:
- [ ] Original issue read completely including all comments
- [ ] Related issues reviewed and understood
- [ ] Current implementation audited against requirements
- [ ] Gaps identified and prioritized
- [ ] Alignment strategy documented

### 2. Dual-Purpose System Architecture
**Guideline**: Design systems that serve both manual and automated workflows for maximum flexibility.

**Architecture Pattern**:
```
Dual-Purpose System:
├── Static Resources (Manual Users)
│   ├── templates/           # Users can copy/fork directly
│   ├── examples/            # Reference implementations
│   └── docs/                # Comprehensive documentation
└── CLI Tools (Automated Users)
    ├── Generate              # Create new projects
    ├── Update                # Update existing projects
    ├── Validate              # Validate configurations
    └── Migrate               # Migrate between versions
```

**Benefits**:
- **Manual Users**: Direct access to templates and examples
- **Automated Users**: Programmatic generation and management
- **Ecosystem Integration**: Works with existing tools and workflows
- **Maximum Adoption**: Supports different user preferences and workflows

### 3. Hierarchical CLI Command Structure
**Guideline**: Use hierarchical command structure for complex CLI tools to improve organization and scalability.

**Command Organization**:
```bash
# Root command
tool-name

# Primary command groups
tool-name generate          # Core functionality
tool-name template          # Template management
tool-name project           # Project management
tool-name config            # Configuration management

# Subcommands
tool-name template list     # List available templates
tool-name template validate # Validate template integrity
tool-name template create   # Create new template

tool-name project update    # Update existing project
tool-name project migrate   # Migrate project between versions
tool-name project validate  # Validate project structure
```

**Implementation with Cobra**:
```go
// Root command
var rootCmd = &cobra.Command{
    Use:   "tool-name",
    Short: "Tool description",
    Long:  "Comprehensive tool description with examples",
}

// Command groups
var templateCmd = &cobra.Command{
    Use:   "template",
    Short: "Manage templates",
    Long:  "Template management operations",
}

var projectCmd = &cobra.Command{
    Use:   "project", 
    Short: "Manage projects",
    Long:  "Project management operations",
}

// Subcommands
var listTemplatesCmd = &cobra.Command{
    Use:   "list",
    Short: "List available templates",
    RunE:  runListTemplates,
}

func init() {
    // Build hierarchy
    templateCmd.AddCommand(listTemplatesCmd)
    rootCmd.AddCommand(templateCmd)
    rootCmd.AddCommand(projectCmd)
}
```

### 4. Comprehensive Template Processing
**Guideline**: Implement template variable processing that covers ALL relevant file types comprehensively.

**File Type Coverage**:
```go
func needsTemplateProcessing(filePath string) bool {
    ext := filepath.Ext(filePath)
    baseName := filepath.Base(filePath)
    
    // Source code files
    sourceExts := []string{".go", ".js", ".ts", ".py", ".java", ".cs", ".rs"}
    
    // Configuration files
    configExts := []string{".yaml", ".yml", ".json", ".toml", ".ini", ".xml"}
    
    // Script files
    scriptExts := []string{".sh", ".bat", ".ps1", ".fish"}
    
    // Documentation files
    docExts := []string{".md", ".txt", ".rst", ".adoc"}
    
    // Special files by name
    specialFiles := []string{
        "go.mod", "package.json", "requirements.txt", "Cargo.toml",
        "Makefile", "Dockerfile", "docker-compose.yml",
    }
    
    return contains(sourceExts, ext) || 
           contains(configExts, ext) || 
           contains(scriptExts, ext) || 
           contains(docExts, ext) || 
           contains(specialFiles, baseName) ||
           strings.HasSuffix(filePath, ".tmpl")
}
```

**Template Variable Standards**:
```go
// Use hierarchical, descriptive variable names
type TemplateContext struct {
    Config struct {
        Name        string
        Description string
        GoModule    string
        Version     string
    }
    Features struct {
        Kubernetes    bool
        TypeScript    bool
        OpenTelemetry bool
        CloudEvents   bool
    }
    Metadata struct {
        Timestamp string
        Generator string
        Version   string
    }
}
```

### 5. Excellent CLI User Experience
**Guideline**: Provide exceptional user experience with clear feedback, helpful errors, and intuitive interactions.

**User Experience Standards**:
```bash
# Progress indicators with emojis
🚀 Generating project from basic template...
📋 Using template: Basic tier health endpoint template (1.0.0)
✅ Successfully generated project from basic template!
📁 Project created in: my-project

# Clear, actionable error messages
❌ Template tier 'invalid' not found in templates directory
💡 Available tiers: basic, intermediate, advanced, enterprise
🔧 Use 'tool-name template list' to see all available templates

# Helpful success feedback with next steps
✅ Project generated successfully!
📁 Location: ./my-project
🚀 Next steps:
   cd my-project
   go mod tidy
   go run cmd/server/main.go
```

**Implementation Pattern**:
```go
func provideFeedback(operation, status, details string) {
    switch status {
    case "start":
        fmt.Printf("🚀 %s...\n", operation)
    case "success":
        fmt.Printf("✅ %s completed successfully!\n", operation)
        if details != "" {
            fmt.Printf("📁 %s\n", details)
        }
    case "error":
        fmt.Printf("❌ %s failed\n", operation)
        if details != "" {
            fmt.Printf("💡 %s\n", details)
        }
    }
}
```

### 6. Template Metadata Management
**Guideline**: Use structured metadata to enable powerful template management and validation.

**Metadata Schema**:
```yaml
# template.yaml
name: string              # Template identifier
description: string       # Human-readable description
tier: string             # Complexity tier
version: string          # Template version
features:                # Feature flags
  feature_name: boolean
requirements:            # System requirements
  go_version: string
  node_version: string
dependencies:            # Template dependencies
  - package_name
maintainers:             # Template maintainers
  - name: string
    email: string
```

**Usage in CLI**:
```go
func listTemplates() error {
    templates, err := discoverTemplates("templates/")
    if err != nil {
        return err
    }
    
    fmt.Println("📋 Available Templates:")
    for _, tmpl := range templates {
        fmt.Printf("🎯 %s (%s)\n", tmpl.Name, tmpl.Version)
        fmt.Printf("   %s\n", tmpl.Description)
        fmt.Printf("   Features: %s\n", formatFeatures(tmpl.Features))
    }
    
    return nil
}
```

### 7. Multi-Level Testing Strategy
**Guideline**: Implement comprehensive testing that covers all levels from templates to running applications.

**Testing Levels**:
1. **Template Validation**: Structure, metadata, syntax
2. **CLI Testing**: Command execution, error handling
3. **Generation Testing**: Template processing, variable substitution
4. **Compilation Testing**: Generated projects compile
5. **Runtime Testing**: Generated applications work
6. **Integration Testing**: End-to-end workflows

**Automated Testing Framework**:
```bash
#!/bin/bash
# comprehensive-test.sh

echo "🧪 Running Comprehensive Tests"

# Level 1: Template validation
echo "📋 Validating templates..."
./bin/tool-name template validate

# Level 2: CLI testing
echo "🖥️ Testing CLI commands..."
go test -v ./cmd/...

# Level 3: Generation testing
echo "🏗️ Testing template generation..."
for tier in basic intermediate advanced; do
    ./bin/tool-name generate --name "test-$tier" --tier "$tier" --module "github.com/test/$tier"
    
    # Level 4: Compilation testing
    cd "test-$tier"
    go mod tidy && go build ./...
    cd ..
    
    # Level 5: Runtime testing (if applicable)
    # Start service and test endpoints
    
    # Cleanup
    rm -rf "test-$tier"
done

echo "✅ All tests passed!"
```

### 8. Repository Structure Standards
**Guideline**: Follow established ecosystem patterns for repository organization.

**Template Repository Structure**:
```
template-{name}/
├── README.md                     # Comprehensive documentation
├── LICENSE                       # Appropriate license
├── templates/                    # Static template directories
│   ├── basic/                    # Basic tier template
│   │   ├── template.yaml         # Template metadata
│   │   ├── cmd/server/main.go
│   │   ├── go.mod.tmpl
│   │   └── README.md.tmpl
│   ├── intermediate/             # Intermediate tier
│   ├── advanced/                 # Advanced tier
│   └── enterprise/               # Enterprise tier
├── examples/                     # Generated examples
│   ├── basic-example/
│   ├── intermediate-example/
│   └── advanced-example/
├── scripts/                      # Automation scripts
│   ├── generate.sh               # Template generation
│   ├── validate.sh               # Template validation
│   └── test.sh                   # Comprehensive testing
├── cmd/                          # CLI tool
│   └── generator/
│       ├── main.go
│       └── commands/
├── pkg/                          # Core libraries
│   ├── templates/
│   ├── validation/
│   └── cli/
└── docs/                         # Documentation
    ├── setup.md
    ├── template-guide.md
    ├── cli-reference.md
    └── migration-guide.md
```

### 9. Error Handling and Validation
**Guideline**: Implement robust error handling with clear, actionable error messages.

**Error Handling Pattern**:
```go
func validateAndExecute(operation func() error, context string) error {
    // Pre-validation
    if err := preValidate(context); err != nil {
        return fmt.Errorf("validation failed for %s: %w\n💡 %s", 
            context, err, getSuggestion(err))
    }
    
    // Execute operation
    if err := operation(); err != nil {
        return fmt.Errorf("operation failed for %s: %w\n🔧 %s", 
            context, err, getFixSuggestion(err))
    }
    
    // Post-validation
    if err := postValidate(context); err != nil {
        return fmt.Errorf("post-validation failed for %s: %w\n⚠️ %s", 
            context, err, getWarning(err))
    }
    
    return nil
}
```

**Error Message Standards**:
- **Context**: What operation was being performed
- **Cause**: What specifically went wrong
- **Suggestion**: How to fix the problem
- **Visual Indicators**: Use emojis for quick recognition

### 10. Documentation and Help Systems
**Guideline**: Provide comprehensive documentation and help systems that scale with complexity.

**Documentation Structure**:
```
docs/
├── README.md                     # Quick start and overview
├── setup.md                      # Installation and setup
├── user-guide.md                 # User guide with examples
├── cli-reference.md              # Complete CLI reference
├── template-guide.md             # Template usage and customization
├── migration-guide.md            # Migration between versions
├── troubleshooting.md            # Common issues and solutions
└── contributing.md               # Contribution guidelines
```

**CLI Help Standards**:
```go
var rootCmd = &cobra.Command{
    Use:   "tool-name",
    Short: "Brief description",
    Long: `Comprehensive description with context and purpose.

This tool provides [main functionality] with support for [key features].
It follows [relevant standards] and integrates with [ecosystem tools].`,
    Example: `  # Basic usage
  tool-name generate --name my-project --tier basic
  
  # Advanced usage with options
  tool-name generate --name my-project --tier advanced --features opentelemetry,cloudevents
  
  # List available templates
  tool-name template list`,
}
```

## Quality Assurance

### Validation Checklist
- [ ] Original requirements fully understood and documented
- [ ] Current implementation audited against requirements
- [ ] Dual-purpose architecture serves both manual and automated workflows
- [ ] CLI commands follow hierarchical structure
- [ ] Template processing covers all relevant file types
- [ ] User experience provides clear feedback and helpful errors
- [ ] Template metadata enables powerful management features
- [ ] Multi-level testing validates entire pipeline
- [ ] Repository structure follows established patterns
- [ ] Documentation is comprehensive and accessible

### Performance Standards
- CLI operations should complete in reasonable time
- Template processing should be efficient for large projects
- Error handling should be fast and not impact user experience
- Help systems should be responsive and comprehensive

### Maintenance Guidelines
- Regular validation against original requirements
- Continuous improvement of user experience
- Template updates and version management
- Documentation updates with new features
- Community feedback integration

## Conclusion
These guidelines provide a framework for building high-quality, user-friendly CLI tools and template systems that align with original requirements and provide excellent user experience. Following these principles ensures consistency, quality, and long-term maintainability while serving diverse user needs and workflows.
