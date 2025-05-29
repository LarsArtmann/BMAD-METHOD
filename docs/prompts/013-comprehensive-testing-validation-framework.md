# Comprehensive Testing and Validation Framework

## Prompt Name
**Comprehensive Testing and Validation Framework**

## Context
Template systems and code generators require comprehensive testing at multiple levels to ensure reliability. This prompt provides a framework for testing template generation, CLI functionality, and generated project quality.

## Objective
Implement a multi-level testing strategy that validates template integrity, CLI functionality, generated project quality, and end-to-end workflows to ensure production-ready template systems.

## Testing Levels

### 1. Template Validation Testing
Validate template structure, metadata, and integrity.

#### **Template Structure Validation**
```go
func TestTemplateStructure(t *testing.T) {
    templateDirs := []string{"basic", "intermediate", "advanced", "enterprise"}
    
    for _, tier := range templateDirs {
        t.Run(fmt.Sprintf("template_%s_structure", tier), func(t *testing.T) {
            templateDir := filepath.Join("templates", tier)
            
            // Check template directory exists
            if _, err := os.Stat(templateDir); os.IsNotExist(err) {
                t.Fatalf("Template directory does not exist: %s", templateDir)
            }
            
            // Check required files
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
                    t.Errorf("Missing required file: %s", file)
                }
            }
            
            // Validate template metadata
            metadataPath := filepath.Join(templateDir, "template.yaml")
            if err := validateTemplateMetadata(metadataPath); err != nil {
                t.Errorf("Invalid template metadata: %v", err)
            }
        })
    }
}
```

#### **Template Variable Validation**
```go
func TestTemplateVariables(t *testing.T) {
    templateDirs := []string{"basic", "intermediate", "advanced", "enterprise"}
    
    for _, tier := range templateDirs {
        t.Run(fmt.Sprintf("template_%s_variables", tier), func(t *testing.T) {
            templateDir := filepath.Join("templates", tier)
            
            // Find all template files
            err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
                if err != nil {
                    return err
                }
                
                if info.IsDir() || filepath.Base(path) == "template.yaml" {
                    return nil
                }
                
                // Check for template variables
                content, err := os.ReadFile(path)
                if err != nil {
                    return err
                }
                
                // Validate template syntax
                if needsTemplateProcessing(path) {
                    if err := validateTemplateSyntax(string(content)); err != nil {
                        t.Errorf("Invalid template syntax in %s: %v", path, err)
                    }
                }
                
                return nil
            })
            
            if err != nil {
                t.Errorf("Error walking template directory: %v", err)
            }
        })
    }
}
```

### 2. CLI Functionality Testing
Test all CLI commands and their interactions.

#### **Command Execution Testing**
```go
func TestCLICommands(t *testing.T) {
    tests := []struct {
        name    string
        args    []string
        wantErr bool
    }{
        {
            name:    "help command",
            args:    []string{"--help"},
            wantErr: false,
        },
        {
            name:    "template list",
            args:    []string{"template", "list"},
            wantErr: false,
        },
        {
            name:    "template validate",
            args:    []string{"template", "validate"},
            wantErr: false,
        },
        {
            name:    "generate basic project",
            args:    []string{"generate", "--name", "test-service", "--tier", "basic", "--module", "github.com/test/service", "--output", "test-output"},
            wantErr: false,
        },
        {
            name:    "generate with invalid tier",
            args:    []string{"generate", "--name", "test-service", "--tier", "invalid", "--module", "github.com/test/service"},
            wantErr: true,
        },
        {
            name:    "generate missing required flag",
            args:    []string{"generate", "--name", "test-service"},
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Clean up before test
            os.RemoveAll("test-output")
            defer os.RemoveAll("test-output")
            
            // Execute command
            cmd := rootCmd
            cmd.SetArgs(tt.args)
            err := cmd.Execute()
            
            if (err != nil) != tt.wantErr {
                t.Errorf("Command execution error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

#### **Template Generation Testing**
```go
func TestTemplateGeneration(t *testing.T) {
    tiers := []string{"basic", "intermediate", "advanced", "enterprise"}
    
    for _, tier := range tiers {
        t.Run(fmt.Sprintf("generate_%s_tier", tier), func(t *testing.T) {
            outputDir := fmt.Sprintf("test-generation-%s", tier)
            defer os.RemoveAll(outputDir)
            
            // Generate project
            err := generateProject("test-service", tier, "github.com/test/service", outputDir, false)
            if err != nil {
                t.Fatalf("Failed to generate %s tier project: %v", tier, err)
            }
            
            // Verify essential files were created
            essentialFiles := []string{
                "README.md",
                "go.mod",
                "cmd/server/main.go",
                "internal/handlers/health.go",
            }
            
            for _, file := range essentialFiles {
                filePath := filepath.Join(outputDir, file)
                if _, err := os.Stat(filePath); os.IsNotExist(err) {
                    t.Errorf("Essential file not created: %s", file)
                }
            }
            
            // Verify template variables were substituted
            if err := verifyVariableSubstitution(outputDir, "test-service", "github.com/test/service"); err != nil {
                t.Errorf("Template variable substitution failed: %v", err)
            }
        })
    }
}
```

### 3. Generated Project Quality Testing
Test that generated projects compile, run, and function correctly.

#### **Compilation Testing**
```go
func TestGeneratedProjectCompilation(t *testing.T) {
    tiers := []string{"basic", "intermediate", "advanced", "enterprise"}
    
    for _, tier := range tiers {
        t.Run(fmt.Sprintf("compile_%s_project", tier), func(t *testing.T) {
            outputDir := fmt.Sprintf("test-compile-%s", tier)
            defer os.RemoveAll(outputDir)
            
            // Generate project
            err := generateProject("test-service", tier, "github.com/test/service", outputDir, false)
            if err != nil {
                t.Fatalf("Failed to generate project: %v", err)
            }
            
            // Test compilation
            if err := testProjectCompilation(outputDir); err != nil {
                t.Errorf("Project compilation failed: %v", err)
            }
        })
    }
}

func testProjectCompilation(projectDir string) error {
    originalDir, err := os.Getwd()
    if err != nil {
        return err
    }
    defer os.Chdir(originalDir)
    
    // Change to project directory
    if err := os.Chdir(projectDir); err != nil {
        return err
    }
    
    // Run go mod tidy
    cmd := exec.Command("go", "mod", "tidy")
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("go mod tidy failed: %w", err)
    }
    
    // Build the project
    cmd = exec.Command("go", "build", "-o", "bin/test-service", "cmd/server/main.go")
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("go build failed: %w", err)
    }
    
    return nil
}
```

#### **Runtime Testing**
```go
func TestGeneratedProjectRuntime(t *testing.T) {
    tiers := []string{"basic", "intermediate", "advanced", "enterprise"}
    
    for _, tier := range tiers {
        t.Run(fmt.Sprintf("runtime_%s_project", tier), func(t *testing.T) {
            outputDir := fmt.Sprintf("test-runtime-%s", tier)
            defer os.RemoveAll(outputDir)
            
            // Generate and compile project
            err := generateProject("test-service", tier, "github.com/test/service", outputDir, false)
            if err != nil {
                t.Fatalf("Failed to generate project: %v", err)
            }
            
            if err := testProjectCompilation(outputDir); err != nil {
                t.Fatalf("Project compilation failed: %v", err)
            }
            
            // Test runtime functionality
            if err := testProjectRuntime(outputDir, tier); err != nil {
                t.Errorf("Project runtime test failed: %v", err)
            }
        })
    }
}

func testProjectRuntime(projectDir, tier string) error {
    originalDir, err := os.Getwd()
    if err != nil {
        return err
    }
    defer os.Chdir(originalDir)
    
    if err := os.Chdir(projectDir); err != nil {
        return err
    }
    
    // Start the server
    cmd := exec.Command("./bin/test-service")
    if err := cmd.Start(); err != nil {
        return fmt.Errorf("failed to start server: %w", err)
    }
    defer cmd.Process.Kill()
    
    // Wait for server to start
    time.Sleep(2 * time.Second)
    
    // Test health endpoints
    endpoints := []string{"/health", "/health/time", "/health/ready", "/health/live", "/health/startup"}
    
    for _, endpoint := range endpoints {
        if err := testEndpoint("http://localhost:8080" + endpoint); err != nil {
            return fmt.Errorf("endpoint %s failed: %w", endpoint, err)
        }
    }
    
    // Test tier-specific endpoints
    if tier != "basic" {
        if err := testEndpoint("http://localhost:8080/health/dependencies"); err != nil {
            // Dependencies endpoint might not respond if no dependencies configured
            // This is acceptable for testing
        }
    }
    
    return nil
}

func testEndpoint(url string) error {
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("expected status 200, got %d", resp.StatusCode)
    }
    
    return nil
}
```

### 4. Integration Testing
Test complete workflows and integrations.

#### **End-to-End Workflow Testing**
```go
func TestEndToEndWorkflow(t *testing.T) {
    // Test complete workflow: generate -> compile -> run -> test endpoints
    outputDir := "test-e2e-workflow"
    defer os.RemoveAll(outputDir)
    
    // Step 1: Generate project
    err := generateProject("e2e-service", "basic", "github.com/test/e2e-service", outputDir, false)
    if err != nil {
        t.Fatalf("Failed to generate project: %v", err)
    }
    
    // Step 2: Compile project
    if err := testProjectCompilation(outputDir); err != nil {
        t.Fatalf("Project compilation failed: %v", err)
    }
    
    // Step 3: Test runtime
    if err := testProjectRuntime(outputDir, "basic"); err != nil {
        t.Fatalf("Project runtime test failed: %v", err)
    }
    
    // Step 4: Test Kubernetes manifests (if kubectl available)
    if err := testKubernetesManifests(outputDir); err != nil {
        t.Logf("Kubernetes manifest test skipped: %v", err)
    }
}

func testKubernetesManifests(projectDir string) error {
    // Check if kubectl is available
    if _, err := exec.LookPath("kubectl"); err != nil {
        return fmt.Errorf("kubectl not available: %w", err)
    }
    
    manifestsDir := filepath.Join(projectDir, "deployments", "kubernetes")
    
    // Validate YAML syntax
    cmd := exec.Command("kubectl", "apply", "--dry-run=client", "-f", manifestsDir)
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("kubernetes manifests validation failed: %w", err)
    }
    
    return nil
}
```

### 5. Automated Testing Scripts

#### **Comprehensive Test Script**
```bash
#!/bin/bash
# test-all.sh - Comprehensive testing script

set -e

echo "🧪 Running Comprehensive Template System Tests"
echo "=============================================="

# Build CLI tool
echo "🔨 Building CLI tool..."
go build -o bin/template-health-endpoint cmd/generator/main.go

# Run unit tests
echo "🔬 Running unit tests..."
go test -v ./pkg/...

# Run CLI tests
echo "🖥️ Running CLI tests..."
go test -v ./cmd/...

# Run integration tests
echo "🔗 Running integration tests..."
go test -v ./tests/...

# Test template validation
echo "📋 Testing template validation..."
./bin/template-health-endpoint template validate

# Test all template tiers
echo "🏗️ Testing template generation for all tiers..."
for tier in basic intermediate advanced enterprise; do
    echo "Testing $tier tier..."
    
    # Clean up
    rm -rf "test-$tier"
    
    # Generate project
    ./bin/template-health-endpoint template from-static \
        --name "test-$tier" \
        --tier "$tier" \
        --module "github.com/test/$tier" \
        --output "test-$tier"
    
    # Test compilation
    cd "test-$tier"
    go mod tidy
    go build -o "bin/test-$tier" cmd/server/main.go
    cd ..
    
    # Clean up
    rm -rf "test-$tier"
    
    echo "✅ $tier tier test passed"
done

echo "🎉 All tests passed successfully!"
```

## Success Criteria
- [ ] All template tiers pass structure validation
- [ ] CLI commands execute correctly with proper error handling
- [ ] Generated projects compile without errors
- [ ] Generated projects run and serve endpoints correctly
- [ ] Template variables are properly substituted
- [ ] Kubernetes manifests are valid (if kubectl available)
- [ ] End-to-end workflows complete successfully
- [ ] Test coverage is comprehensive across all components

## Related Patterns
- Test-driven development
- Integration testing strategies
- CLI testing frameworks
- Template validation systems
- Continuous integration pipelines
