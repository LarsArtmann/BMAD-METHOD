# Repository Structure Conversion and Template Management

## Prompt Name
**Repository Structure Conversion and Template Management**

## Context
Converting CLI-embedded templates to static template repositories requires systematic restructuring while maintaining functionality. This prompt guides the conversion from embedded templates to proper template repository structure following established patterns.

## Objective
Convert CLI tools with embedded templates into dual-purpose systems that provide both static template directories (following template-* repository patterns) and CLI generation capabilities.

## Key Requirements

### 1. Target Repository Structure
Follow established template-* repository patterns:

```
template-{name}/
├── README.md                     # Comprehensive documentation
├── LICENSE                       # Appropriate license
├── templates/                    # 🆕 Static template directories
│   ├── basic/                    # Users can copy this directly
│   │   ├── cmd/server/main.go
│   │   ├── internal/handlers/health.go
│   │   ├── go.mod.tmpl
│   │   ├── README.md.tmpl
│   │   ├── template.yaml         # Template metadata
│   │   └── ...
│   ├── intermediate/             # Production-ready template
│   ├── advanced/                 # Full observability template
│   └── enterprise/               # Kubernetes & compliance template
├── examples/                     # Generated examples from templates
│   ├── go/                       # Generated Go examples
│   └── typescript/               # Generated TypeScript examples
├── scripts/                      # Template operations
│   ├── generate.sh               # Wrapper for CLI generate
│   ├── update.sh                 # Wrapper for CLI update
│   ├── validate.sh               # Template validation
│   └── install.sh                # CLI installation
├── cmd/                          # Enhanced CLI tool
│   └── generator/
│       ├── main.go               # CLI entry point
│       ├── commands/
│       │   ├── generate.go       # Generate new projects
│       │   ├── template.go       # Template management
│       │   ├── update.go         # Update existing projects
│       │   └── validate.go       # Template validation
├── pkg/                          # Core logic
│   ├── generator/                # Template generation logic
│   ├── templates/                # Template processing
│   └── validation/               # Validation logic
├── schemas/                      # Generated schemas
│   ├── openapi/                  # Generated OpenAPI specs
│   ├── json-schema/              # Generated JSON schemas
│   └── cloudevents/              # CloudEvents schemas
└── docs/                         # Comprehensive documentation
    ├── setup.md                  # Setup instructions
    ├── template-guide.md          # Template usage guide
    ├── cli-usage.md              # CLI documentation
    └── migration-guide.md        # Migration between tiers
```

### 2. Template Extraction Process

#### **Step 1: Generate Static Templates from CLI**
```bash
# Generate each tier to extract templates
./bin/cli-tool generate --name template-basic --tier basic --module github.com/template/basic --output templates/basic
./bin/cli-tool generate --name template-intermediate --tier intermediate --module github.com/template/intermediate --output templates/intermediate
./bin/cli-tool generate --name template-advanced --tier advanced --module github.com/template/advanced --output templates/advanced
./bin/cli-tool generate --name template-enterprise --tier enterprise --module github.com/template/enterprise --output templates/enterprise
```

#### **Step 2: Convert to Template Variables**
Create conversion script to replace hardcoded values:

```bash
#!/bin/bash
# convert-to-templates.sh

convert_file_to_template() {
    local file=$1
    local temp_file="${file}.tmp"
    
    # Skip binary files and directories
    if [[ -d "$file" ]] || ! [[ -f "$file" ]]; then
        return 0
    fi
    
    # Convert hardcoded values to template variables
    sed \
        -e 's/template-basic/{{.Config.Name}}/g' \
        -e 's/github\.com\/template\/basic/{{.Config.GoModule}}/g' \
        -e 's/Basic health endpoint service/{{.Config.Description}}/g' \
        -e 's/Template Health Endpoint Generator v1\.0\.0/Template Health Endpoint Generator v{{.Version}}/g' \
        -e 's/Generated at: [0-9T:-]*Z/Generated at: {{.Timestamp}}/g' \
        "$file" > "$temp_file"
    
    # Replace original with template version
    mv "$temp_file" "$file"
}

# Process all files in template directory
find templates/ -type f | while read -r file; do
    convert_file_to_template "$file"
done
```

#### **Step 3: Add Template Extensions**
```bash
add_template_extension() {
    local file=$1
    local base_name=$(basename "$file")
    
    # Add .tmpl extension to files that need templating
    case "$base_name" in
        "go.mod"|"README.md"|"package.json"|"Dockerfile"|"docker-compose.yml")
            if [[ ! "$file" =~ \.tmpl$ ]]; then
                mv "$file" "${file}.tmpl"
            fi
            ;;
    esac
}
```

### 3. Template Metadata Management

#### **Template Metadata Schema**
Create `template.yaml` for each template tier:

```yaml
# templates/basic/template.yaml
name: basic
description: Basic tier health endpoint template
tier: basic
features:
  kubernetes: true
  typescript: true
  docker: true
  opentelemetry: false
  cloudevents: false
  dependencies: false
version: "1.0.0"
requirements:
  go_version: "1.21+"
  node_version: "16+" # if typescript enabled
dependencies:
  - github.com/gorilla/mux
  - github.com/spf13/cobra
```

#### **Progressive Feature Matrix**
```yaml
# templates/intermediate/template.yaml
features:
  kubernetes: true
  typescript: true
  docker: true
  opentelemetry: basic
  cloudevents: false
  dependencies: true

# templates/advanced/template.yaml
features:
  kubernetes: true
  typescript: true
  docker: true
  opentelemetry: full
  cloudevents: true
  dependencies: true
  server_timing: true
  metrics: custom

# templates/enterprise/template.yaml
features:
  kubernetes: true
  typescript: true
  docker: true
  opentelemetry: full
  cloudevents: true
  dependencies: true
  server_timing: true
  metrics: custom
  security: mtls
  compliance: true
  multi_env: true
```

### 4. CLI Integration with Static Templates

#### **Template Command Implementation**
```go
// cmd/generator/commands/template.go
func runGenerateFromTemplate(cmd *cobra.Command, args []string) error {
    name, _ := cmd.Flags().GetString("name")
    tier, _ := cmd.Flags().GetString("tier")
    module, _ := cmd.Flags().GetString("module")
    output, _ := cmd.Flags().GetString("output")
    
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
    
    // Create template context
    context := map[string]interface{}{
        "Config": map[string]interface{}{
            "Name":        name,
            "Description": fmt.Sprintf("%s health endpoint service", name),
            "GoModule":    module,
            "Tier":       tier,
            "Features":   metadata.Features,
        },
        "Version":   metadata.Version,
        "Timestamp": time.Now().Format(time.RFC3339),
    }
    
    // Generate project from template
    return generateFromStaticTemplate(templateDir, output, context)
}
```

#### **Template Processing Logic**
```go
func generateFromStaticTemplate(templateDir, outputDir string, context map[string]interface{}) error {
    return filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        // Skip template.yaml metadata file
        if filepath.Base(path) == "template.yaml" {
            return nil
        }
        
        // Calculate relative path and output path
        relPath, err := filepath.Rel(templateDir, path)
        if err != nil {
            return err
        }
        outputPath := filepath.Join(outputDir, relPath)
        
        if info.IsDir() {
            return os.MkdirAll(outputPath, info.Mode())
        }
        
        return processTemplateFile(path, outputPath, context)
    })
}
```

### 5. Template Validation Framework

#### **Template Structure Validation**
```go
func validateTemplateStructure(templateDir string) error {
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
            return fmt.Errorf("missing required file: %s", file)
        }
    }
    
    return nil
}
```

#### **Template Generation Testing**
```go
func testTemplateGeneration(templateDir string) error {
    // Generate test project
    testOutput := filepath.Join(os.TempDir(), "template-test")
    defer os.RemoveAll(testOutput)
    
    context := map[string]interface{}{
        "Config": map[string]interface{}{
            "Name":        "test-service",
            "Description": "Test service",
            "GoModule":    "github.com/test/service",
        },
        "Version":   "1.0.0",
        "Timestamp": time.Now().Format(time.RFC3339),
    }
    
    if err := generateFromStaticTemplate(templateDir, testOutput, context); err != nil {
        return fmt.Errorf("template generation failed: %w", err)
    }
    
    // Test compilation
    return testProjectCompilation(testOutput)
}
```

### 6. Migration Strategy

#### **Phase 1: Extract Templates**
1. Generate static templates from existing CLI
2. Convert hardcoded values to template variables
3. Add template metadata files
4. Validate template structure

#### **Phase 2: CLI Integration**
1. Add template command group to CLI
2. Implement static template generation
3. Add template listing and validation
4. Test CLI with static templates

#### **Phase 3: Documentation and Examples**
1. Update repository documentation
2. Generate example projects
3. Create migration guides
4. Add comprehensive testing

## Success Criteria
- [ ] Static template directories follow template-* repository pattern
- [ ] CLI can generate projects from static templates
- [ ] Template variables are properly substituted
- [ ] Generated projects compile and run successfully
- [ ] Template validation passes for all tiers
- [ ] Documentation is comprehensive and accurate

## Common Pitfalls
1. **Incomplete Variable Conversion**: Missing template variables in some files
2. **File Type Coverage**: Not processing all relevant file types
3. **Path Handling**: Incorrect handling of template paths and extensions
4. **Metadata Inconsistency**: Template metadata doesn't match actual features

## Related Patterns
- Template repository design
- CLI tool architecture
- Code generation systems
- Repository structure standards
- Template processing pipelines
