# BDD Implementation Plan for Template Health Endpoint System

## 🎯 **Overview**

This document outlines the comprehensive Behavior-Driven Development (BDD) implementation for the template-health-endpoint system, ensuring all functionality works exactly as users expect through natural language scenarios.

## 📋 **BDD Implementation Strategy**

### **Phase 1: Core BDD Framework Setup (1 hour)**

#### **Task 1.1: Install BDD Dependencies**
```bash
# Add BDD testing dependencies
go get github.com/cucumber/godog
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/require
```

#### **Task 1.2: Create BDD Test Structure**
```
features/
├── template_generation.feature      # Core generation scenarios
├── static_templates.feature         # Static template usage
├── project_migration.feature        # Tier migration scenarios
├── error_handling.feature           # Error conditions
├── performance.feature              # Performance requirements
├── kubernetes_integration.feature   # K8s integration
├── steps/
│   ├── common_steps.go             # Shared step definitions
│   ├── generation_steps.go         # Generation-specific steps
│   ├── migration_steps.go          # Migration-specific steps
│   └── validation_steps.go         # Validation steps
├── support/
│   ├── test_context.go             # Test context management
│   ├── test_data.go                # Test data helpers
│   └── cleanup.go                  # Cleanup utilities
└── main_test.go                    # BDD test runner
```

### **Phase 2: Core Feature Scenarios (2 hours)**

#### **Task 2.1: Template Generation Features**
```gherkin
# features/template_generation.feature
Feature: Template Generation
  As a developer
  I want to generate projects from different template tiers
  So that I can choose the right complexity level for my needs

  Background:
    Given I have the template-health-endpoint CLI installed
    And I have cleaned up any existing test projects

  Scenario Outline: Generate project from different tiers
    When I run "template-health-endpoint generate --name <project_name> --tier <tier> --module github.com/test/<project_name>"
    Then a new project should be created in "./<project_name>"
    And the project should compile successfully
    And the project should have all required health endpoints
    And the health endpoints should respond correctly
    And the project should have <tier> tier features

    Examples:
      | project_name | tier         |
      | basic-test   | basic        |
      | inter-test   | intermediate |
      | adv-test     | advanced     |
      | ent-test     | enterprise   |

  Scenario: Generate project with custom configuration
    Given I want to customize my project generation
    When I run "template-health-endpoint generate --name custom-service --tier basic --module github.com/org/custom-service --features kubernetes,typescript"
    Then a new project should be created in "./custom-service"
    And the project should have Kubernetes manifests
    And the project should have TypeScript client
    And the project should compile successfully

  Scenario: Generate project with dry run
    When I run "template-health-endpoint generate --name my-service --tier basic --module github.com/org/my-service --dry-run"
    Then the command should succeed
    And the output should mention "dry run mode"
    And the output should show what would be generated
    And no project directory should be created
```

#### **Task 2.2: Static Template Features**
```gherkin
# features/static_templates.feature
Feature: Static Template Usage
  As a developer
  I want to use static templates directly
  So that I can manually customize templates before generation

  Scenario: List available templates
    When I run "template-health-endpoint template list"
    Then the command should succeed
    And the output should list "basic" template
    And the output should list "intermediate" template
    And the output should list "advanced" template
    And the output should list "enterprise" template
    And each template should show its description and features

  Scenario: Validate all templates
    When I run "template-health-endpoint template validate"
    Then the command should succeed
    And all templates should pass structure validation
    And all templates should pass metadata validation
    And the output should confirm "All templates are valid"

  Scenario: Generate from static template
    When I run "template-health-endpoint template from-static --name static-test --tier basic --module github.com/test/static-test"
    Then a new project should be created in "./static-test"
    And the project should have proper template variable substitution
    And the project should compile successfully
    And all health endpoints should work correctly
```

#### **Task 2.3: Project Migration Features**
```gherkin
# features/project_migration.feature
Feature: Project Migration Between Tiers
  As a developer
  I want to migrate my project between template tiers
  So that I can add more features as my service evolves

  Background:
    Given I have a basic tier project called "migration-test"
    And the project is compiled and working
    And all basic health endpoints are responding

  Scenario: Migrate from basic to intermediate
    When I run "template-health-endpoint migrate --project ./migration-test --to intermediate"
    Then the migration should succeed
    And the project should be at intermediate tier
    And the project should have dependency health check files
    And the project should have basic OpenTelemetry configuration
    And the project should still compile successfully
    And all original endpoints should still work
    And new dependency endpoints should be available

  Scenario: Migrate with dry run shows plan
    When I run "template-health-endpoint migrate --project ./migration-test --to advanced --dry-run"
    Then the command should succeed
    And the output should show "Migration plan: basic -> advanced"
    And the output should list files that would be added
    And the output should list files that would be modified
    And the output should mention "dry run mode"
    And the project should remain unchanged

  Scenario: Progressive migration path
    Given I have migrated the project to intermediate tier
    When I run "template-health-endpoint migrate --project ./migration-test --to advanced"
    Then the migration should succeed
    And the project should have CloudEvents integration
    And the project should have Server Timing headers
    And the project should have full OpenTelemetry configuration
    When I run "template-health-endpoint migrate --project ./migration-test --to enterprise"
    Then the migration should succeed
    And the project should have mTLS configuration
    And the project should have compliance logging
    And the project should have multi-environment configs
```

### **Phase 3: Error Handling and Edge Cases (1 hour)**

#### **Task 3.1: Error Scenario Testing**
```gherkin
# features/error_handling.feature
Feature: Error Handling and Validation
  As a developer
  I want clear error messages when something goes wrong
  So that I can quickly understand and fix issues

  Scenario: Generate with invalid tier
    When I run "template-health-endpoint generate --name test --tier invalid --module github.com/test/test"
    Then the command should fail with exit code 1
    And the error message should contain "invalid tier 'invalid'"
    And the error message should list available tiers: "basic, intermediate, advanced, enterprise"
    And the error message should suggest "Use 'template-health-endpoint template list' to see all templates"

  Scenario: Generate with missing required flags
    When I run "template-health-endpoint generate --name test"
    Then the command should fail with exit code 1
    And the error message should contain "required flag"
    And the error message should mention "--module is required"

  Scenario: Generate with invalid module path
    When I run "template-health-endpoint generate --name test --tier basic --module invalid-module-path"
    Then the command should fail with exit code 1
    And the error message should contain "invalid Go module path"

  Scenario: Migrate non-existent project
    When I run "template-health-endpoint migrate --project ./nonexistent --to intermediate"
    Then the command should fail with exit code 1
    And the error message should contain "project not found"
    And the error message should suggest checking the project path

  Scenario: Migrate with invalid tier transition
    Given I have a basic tier project called "invalid-migration-test"
    When I run "template-health-endpoint migrate --project ./invalid-migration-test --to invalid"
    Then the command should fail with exit code 1
    And the error message should contain "invalid tier 'invalid'"
    And the error message should list valid migration targets

  Scenario: Template validation with corrupted template
    Given I have corrupted the basic template metadata
    When I run "template-health-endpoint template validate"
    Then the command should fail with exit code 1
    And the error message should identify the corrupted template
    And the error message should describe the validation error
```

### **Phase 4: Performance and Integration Testing (1 hour)**

#### **Task 4.1: Performance Requirements**
```gherkin
# features/performance.feature
Feature: Performance Requirements
  As a developer
  I want the CLI to perform efficiently
  So that it doesn't slow down my development workflow

  Scenario: CLI commands complete quickly
    When I run "template-health-endpoint template list"
    Then the command should complete in less than 2 seconds

  Scenario: Project generation is reasonably fast
    When I run "template-health-endpoint generate --name perf-test --tier basic --module github.com/test/perf-test"
    Then the command should complete in less than 10 seconds
    And a working project should be created

  Scenario: Generated service handles concurrent requests
    Given I have a running basic tier service
    When I send 50 concurrent requests to "/health"
    Then all requests should return status 200
    And the average response time should be less than 100ms
    And no requests should fail

  Scenario: Template validation is fast
    When I run "template-health-endpoint template validate"
    Then the command should complete in less than 5 seconds
    And all templates should be validated
```

#### **Task 4.2: Kubernetes Integration**
```gherkin
# features/kubernetes_integration.feature
Feature: Kubernetes Integration
  As a DevOps engineer
  I want generated projects to work seamlessly with Kubernetes
  So that I can deploy them in production environments

  Background:
    Given I have kubectl available
    And I have a test Kubernetes namespace

  Scenario: Generated Kubernetes manifests are valid
    Given I have generated a basic tier project called "k8s-test"
    When I validate the Kubernetes manifests with "kubectl apply --dry-run=client -f deployments/kubernetes/"
    Then all manifests should be valid YAML
    And kubectl should report no validation errors

  Scenario: Health probes are properly configured
    Given I have generated an intermediate tier project called "k8s-health-test"
    When I examine the deployment manifest
    Then the liveness probe should point to "/health/live"
    And the readiness probe should point to "/health/ready"
    And the startup probe should point to "/health/startup"
    And all probes should have appropriate timeouts and intervals

  Scenario: Service discovery works correctly
    Given I have deployed a generated project to Kubernetes
    When I check the service endpoints
    Then the service should be accessible within the cluster
    And health endpoints should respond correctly
    And the service should be marked as ready
```

### **Phase 5: Advanced Scenarios (1 hour)**

#### **Task 5.1: Template Customization**
```gherkin
# features/template_customization.feature
Feature: Template Customization
  As a developer
  I want to customize templates for my organization's needs
  So that generated projects follow our standards

  Scenario: Interactive template customization
    Given I want to customize a template interactively
    When I run "template-health-endpoint customize --tier basic --interactive"
    Then I should be prompted for customization options
    And I should be able to specify custom features
    And I should be able to save my customization profile
    And the customized template should generate correctly

  Scenario: Use saved customization profile
    Given I have a saved customization profile called "my-org"
    When I run "template-health-endpoint customize --tier basic --profile my-org"
    Then the template should be customized according to the profile
    And the generated project should include organization-specific configurations

  Scenario: Batch customization with config file
    Given I have a customization config file "custom.yaml"
    When I run "template-health-endpoint customize --tier advanced --config custom.yaml"
    Then the template should be customized according to the config file
    And the generated project should reflect all customizations
```

#### **Task 5.2: Update and Maintenance**
```gherkin
# features/project_updates.feature
Feature: Project Updates and Maintenance
  As a developer
  I want to update my existing projects to newer template versions
  So that I can benefit from improvements and security updates

  Scenario: Update project to newer template version
    Given I have a basic tier project at version "1.0.0"
    And a newer template version "1.1.0" is available
    When I run "template-health-endpoint update --project ./my-project --template-version 1.1.0"
    Then the project should be updated to version "1.1.0"
    And new features should be available
    And existing functionality should remain intact
    And the project should still compile and run

  Scenario: Selective update of project components
    Given I have an intermediate tier project
    When I run "template-health-endpoint update --project ./my-project --selective kubernetes,docs"
    Then only Kubernetes manifests should be updated
    And only documentation should be updated
    And source code should remain unchanged
    And the project should still work correctly

  Scenario: Update with conflict resolution
    Given I have a modified basic tier project
    And I have made custom changes to health handlers
    When I run "template-health-endpoint update --project ./my-project --template-version 1.1.0"
    Then the update should detect conflicts
    And I should be prompted for conflict resolution
    And I should be able to choose how to handle each conflict
    And the final project should work correctly
```

## 🔧 **Implementation Details**

### **Step Definition Implementation**
```go
// features/steps/common_steps.go
package steps

import (
    "context"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "time"
    "net/http"
    "github.com/cucumber/godog"
    "github.com/stretchr/testify/assert"
)

type TestContext struct {
    LastCommand    *exec.Cmd
    LastOutput     string
    LastError      error
    LastExitCode   int
    ProjectPaths   []string
    ServerProcess  *os.Process
    StartTime      time.Time
}

func (tc *TestContext) Cleanup() {
    // Stop any running servers
    if tc.ServerProcess != nil {
        tc.ServerProcess.Kill()
    }
    
    // Clean up generated projects
    for _, path := range tc.ProjectPaths {
        os.RemoveAll(path)
    }
}

func InitializeCommonSteps(ctx *godog.ScenarioContext) {
    ctx.Given(`^I have the template-health-endpoint CLI installed$`, iHaveTheCLIInstalled)
    ctx.Given(`^I have cleaned up any existing test projects$`, iHaveCleanedUpTestProjects)
    ctx.When(`^I run "([^"]*)"$`, iRunCommand)
    ctx.Then(`^the command should succeed$`, theCommandShouldSucceed)
    ctx.Then(`^the command should fail with exit code (\d+)$`, theCommandShouldFailWithExitCode)
    ctx.Then(`^the command should complete in less than (\d+) seconds$`, theCommandShouldCompleteInLessThanSeconds)
    ctx.Then(`^the output should contain "([^"]*)"$`, theOutputShouldContain)
    ctx.Then(`^the error message should contain "([^"]*)"$`, theErrorMessageShouldContain)
}

func iHaveTheCLIInstalled(ctx context.Context) error {
    if _, err := os.Stat("./bin/template-health-endpoint"); os.IsNotExist(err) {
        return fmt.Errorf("CLI binary not found. Please run: go build -o bin/template-health-endpoint cmd/generator/main.go")
    }
    return nil
}

func iRunCommand(ctx context.Context, command string) error {
    testCtx := getTestContext(ctx)
    testCtx.StartTime = time.Now()
    
    // Parse command
    parts := strings.Fields(command)
    if len(parts) == 0 {
        return fmt.Errorf("empty command")
    }
    
    // Execute command
    cmd := exec.Command(parts[0], parts[1:]...)
    output, err := cmd.CombinedOutput()
    
    testCtx.LastCommand = cmd
    testCtx.LastOutput = string(output)
    testCtx.LastError = err
    
    if cmd.ProcessState != nil {
        testCtx.LastExitCode = cmd.ProcessState.ExitCode()
    }
    
    return nil
}

func theCommandShouldSucceed(ctx context.Context) error {
    testCtx := getTestContext(ctx)
    if testCtx.LastExitCode != 0 {
        return fmt.Errorf("command failed with exit code %d. Output: %s", 
            testCtx.LastExitCode, testCtx.LastOutput)
    }
    return nil
}

func theCommandShouldCompleteInLessThanSeconds(ctx context.Context, maxSeconds int) error {
    testCtx := getTestContext(ctx)
    duration := time.Since(testCtx.StartTime)
    maxDuration := time.Duration(maxSeconds) * time.Second
    
    if duration > maxDuration {
        return fmt.Errorf("command took %v, expected less than %v", duration, maxDuration)
    }
    return nil
}
```

### **Test Data Management**
```go
// features/support/test_data.go
package support

import (
    "fmt"
    "os"
    "path/filepath"
    "os/exec"
)

type TestDataManager struct {
    TempDir      string
    ProjectPaths []string
}

func NewTestDataManager() *TestDataManager {
    tempDir, _ := os.MkdirTemp("", "bdd-test-*")
    return &TestDataManager{
        TempDir: tempDir,
    }
}

func (tdm *TestDataManager) CreateBasicProject(name string) (string, error) {
    projectPath := filepath.Join(tdm.TempDir, name)
    
    cmd := exec.Command("./bin/template-health-endpoint", "generate",
        "--name", name,
        "--tier", "basic",
        "--module", fmt.Sprintf("github.com/test/%s", name),
        "--output", projectPath)
    
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("failed to create test project: %w", err)
    }
    
    tdm.ProjectPaths = append(tdm.ProjectPaths, projectPath)
    return projectPath, nil
}

func (tdm *TestDataManager) CompileProject(projectPath string) error {
    originalDir, err := os.Getwd()
    if err != nil {
        return err
    }
    defer os.Chdir(originalDir)
    
    if err := os.Chdir(projectPath); err != nil {
        return err
    }
    
    // Run go mod tidy
    if err := exec.Command("go", "mod", "tidy").Run(); err != nil {
        return fmt.Errorf("go mod tidy failed: %w", err)
    }
    
    // Build project
    if err := exec.Command("go", "build", "-o", "bin/test-service", "cmd/server/main.go").Run(); err != nil {
        return fmt.Errorf("build failed: %w", err)
    }
    
    return nil
}

func (tdm *TestDataManager) Cleanup() {
    os.RemoveAll(tdm.TempDir)
}
```

### **CI/CD Integration**
```yaml
# .github/workflows/bdd-tests.yml
name: BDD Tests

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  bdd-tests:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
    
    - name: Install dependencies
      run: |
        go mod download
        go get github.com/cucumber/godog
    
    - name: Build CLI
      run: go build -o bin/template-health-endpoint cmd/generator/main.go
    
    - name: Run BDD tests
      run: |
        cd features
        go test -v -godog.format=pretty
    
    - name: Generate BDD report
      run: |
        cd features
        go test -v -godog.format=cucumber > ../bdd-results.json
    
    - name: Upload BDD results
      uses: actions/upload-artifact@v3
      if: always()
      with:
        name: bdd-results
        path: bdd-results.json
```

## 📊 **Success Criteria**

### **BDD Coverage Requirements**
- [ ] All CLI commands have BDD scenarios
- [ ] All template tiers have generation and validation scenarios
- [ ] All migration paths have comprehensive BDD coverage
- [ ] Error conditions are thoroughly tested with BDD
- [ ] Performance requirements are validated through BDD
- [ ] Integration points (Kubernetes, etc.) are tested
- [ ] User workflows are covered end-to-end

### **Quality Standards**
- [ ] All BDD scenarios pass consistently
- [ ] Test execution time is reasonable (< 5 minutes total)
- [ ] Test reports are generated and accessible
- [ ] Test data is properly managed and cleaned up
- [ ] CI/CD pipeline includes BDD tests
- [ ] BDD scenarios serve as living documentation

### **User Experience Validation**
- [ ] All user stories are covered by BDD scenarios
- [ ] Error messages are validated for clarity and helpfulness
- [ ] Performance expectations are met and tested
- [ ] Integration scenarios work in realistic environments
- [ ] Documentation scenarios ensure accuracy

## 🎯 **Implementation Timeline**

**Total Estimated Time: 6 hours**

1. **Phase 1: BDD Framework Setup** (1 hour)
2. **Phase 2: Core Feature Scenarios** (2 hours)
3. **Phase 3: Error Handling** (1 hour)
4. **Phase 4: Performance & Integration** (1 hour)
5. **Phase 5: Advanced Scenarios** (1 hour)

This BDD implementation will ensure that our template system works exactly as users expect, providing comprehensive validation of all functionality through natural language scenarios that serve as both tests and living documentation.

---

**Ready to implement comprehensive BDD testing for the template-health-endpoint system!** 🚀
