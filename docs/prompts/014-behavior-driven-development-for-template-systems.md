# Behavior-Driven Development (BDD) for Template Systems

## Prompt Name
**Behavior-Driven Development (BDD) for Template Systems**

## Context
Template systems and CLI tools require testing from the user's perspective to ensure they behave correctly in real-world scenarios. BDD provides a framework for writing tests that describe expected behavior in natural language, making them accessible to both technical and non-technical stakeholders.

## Objective
Implement comprehensive BDD testing for template systems using Gherkin syntax to describe user scenarios, validate expected behaviors, and ensure the system works as intended from the user's perspective.

## Key Principles

### 1. User-Centric Test Design
Write tests from the user's perspective, focusing on what they want to accomplish:

```gherkin
Feature: Generate health endpoint project
  As a developer
  I want to generate a health endpoint project from a template
  So that I can quickly bootstrap a production-ready service

  Scenario: Generate basic tier project
    Given I have the template-health-endpoint CLI installed
    When I run "template-health-endpoint generate --name my-service --tier basic --module github.com/org/my-service"
    Then a new project should be created in "./my-service"
    And the project should compile successfully
    And the project should have all required health endpoints
    And the health endpoints should respond correctly
```

### 2. Comprehensive Scenario Coverage
Cover all user workflows and edge cases:

```gherkin
Feature: Template tier progression
  As a developer
  I want to migrate my project between template tiers
  So that I can add more features as my service grows

  Background:
    Given I have a basic tier project called "my-service"
    And the project is compiled and working

  Scenario: Migrate from basic to intermediate
    When I run "template-health-endpoint migrate --project ./my-service --to intermediate"
    Then the project should be upgraded to intermediate tier
    And new dependency health check endpoints should be available
    And basic OpenTelemetry should be configured
    And the project should still compile and run

  Scenario: Migrate from intermediate to advanced
    Given the project is at intermediate tier
    When I run "template-health-endpoint migrate --project ./my-service --to advanced"
    Then the project should be upgraded to advanced tier
    And CloudEvents integration should be available
    And Server Timing headers should be present
    And full OpenTelemetry should be configured

  Scenario: Migrate from advanced to enterprise
    Given the project is at advanced tier
    When I run "template-health-endpoint migrate --project ./my-service --to enterprise"
    Then the project should be upgraded to enterprise tier
    And mTLS configuration should be available
    And compliance logging should be enabled
    And multi-environment configs should be present
```

### 3. Error Scenario Testing
Test error conditions and edge cases:

```gherkin
Feature: Error handling and validation
  As a developer
  I want clear error messages when something goes wrong
  So that I can quickly understand and fix issues

  Scenario: Generate project with invalid tier
    When I run "template-health-endpoint generate --name my-service --tier invalid --module github.com/org/my-service"
    Then the command should fail with exit code 1
    And the error message should mention "invalid tier"
    And the error message should list available tiers
    And the error message should suggest using "template list"

  Scenario: Generate project with missing required flags
    When I run "template-health-endpoint generate --name my-service"
    Then the command should fail with exit code 1
    And the error message should mention "required flag"
    And the error message should specify which flags are missing

  Scenario: Migrate project with invalid path
    When I run "template-health-endpoint migrate --project ./nonexistent --to intermediate"
    Then the command should fail with exit code 1
    And the error message should mention "project not found"
    And the error message should suggest checking the path
```

## Implementation Framework

### 1. BDD Test Structure
Organize BDD tests using the standard Given-When-Then structure:

```go
// features/generate_project_test.go
package features

import (
    "context"
    "testing"
    "github.com/cucumber/godog"
)

func TestFeatures(t *testing.T) {
    suite := godog.TestSuite{
        ScenarioInitializer: InitializeScenario,
        Options: &godog.Options{
            Format:   "pretty",
            Paths:    []string{"features"},
            TestingT: t,
        },
    }

    if suite.Run() != 0 {
        t.Fatal("non-zero status returned, failed to run feature tests")
    }
}

func InitializeScenario(ctx *godog.ScenarioContext) {
    // Step definitions
    ctx.Given(`^I have the template-health-endpoint CLI installed$`, iHaveTheCLIInstalled)
    ctx.When(`^I run "([^"]*)"$`, iRunCommand)
    ctx.Then(`^a new project should be created in "([^"]*)"$`, aNewProjectShouldBeCreated)
    ctx.Then(`^the project should compile successfully$`, theProjectShouldCompile)
    ctx.Then(`^the project should have all required health endpoints$`, theProjectShouldHaveHealthEndpoints)
    ctx.Then(`^the health endpoints should respond correctly$`, theHealthEndpointsShouldRespond)
}
```

### 2. Step Definitions
Implement step definitions that map Gherkin steps to Go code:

```go
// features/steps.go
package features

import (
    "context"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "time"
    "net/http"
)

type TestContext struct {
    LastCommand    *exec.Cmd
    LastOutput     string
    LastError      error
    LastExitCode   int
    ProjectPath    string
    ServerProcess  *os.Process
}

func iHaveTheCLIInstalled(ctx context.Context) error {
    // Verify CLI binary exists
    if _, err := os.Stat("./bin/template-health-endpoint"); os.IsNotExist(err) {
        return fmt.Errorf("CLI binary not found, please run: go build -o bin/template-health-endpoint cmd/generator/main.go")
    }
    return nil
}

func iRunCommand(ctx context.Context, command string) error {
    testCtx := getTestContext(ctx)
    
    // Parse command and arguments
    args := parseCommand(command)
    
    // Execute command
    cmd := exec.Command(args[0], args[1:]...)
    output, err := cmd.CombinedOutput()
    
    testCtx.LastCommand = cmd
    testCtx.LastOutput = string(output)
    testCtx.LastError = err
    testCtx.LastExitCode = cmd.ProcessState.ExitCode()
    
    return nil
}

func aNewProjectShouldBeCreated(ctx context.Context, projectPath string) error {
    testCtx := getTestContext(ctx)
    testCtx.ProjectPath = projectPath
    
    // Check if project directory exists
    if _, err := os.Stat(projectPath); os.IsNotExist(err) {
        return fmt.Errorf("project directory %s was not created", projectPath)
    }
    
    // Check for essential files
    essentialFiles := []string{
        "go.mod",
        "cmd/server/main.go",
        "internal/handlers/health.go",
        "README.md",
    }
    
    for _, file := range essentialFiles {
        fullPath := filepath.Join(projectPath, file)
        if _, err := os.Stat(fullPath); os.IsNotExist(err) {
            return fmt.Errorf("essential file %s not found in generated project", file)
        }
    }
    
    return nil
}

func theProjectShouldCompile(ctx context.Context) error {
    testCtx := getTestContext(ctx)
    
    // Change to project directory
    originalDir, err := os.Getwd()
    if err != nil {
        return err
    }
    defer os.Chdir(originalDir)
    
    if err := os.Chdir(testCtx.ProjectPath); err != nil {
        return fmt.Errorf("failed to change to project directory: %w", err)
    }
    
    // Run go mod tidy
    cmd := exec.Command("go", "mod", "tidy")
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("go mod tidy failed: %w", err)
    }
    
    // Build the project
    cmd = exec.Command("go", "build", "-o", "bin/test-service", "cmd/server/main.go")
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("project compilation failed: %w", err)
    }
    
    return nil
}

func theProjectShouldHaveHealthEndpoints(ctx context.Context) error {
    testCtx := getTestContext(ctx)
    
    // Check health handler file
    handlerPath := filepath.Join(testCtx.ProjectPath, "internal/handlers/health.go")
    content, err := os.ReadFile(handlerPath)
    if err != nil {
        return fmt.Errorf("failed to read health handler: %w", err)
    }
    
    handlerContent := string(content)
    
    // Verify required endpoints are present
    requiredMethods := []string{
        "CheckHealth",
        "ServerTime",
        "ReadinessCheck",
        "LivenessCheck",
        "StartupCheck",
    }
    
    for _, method := range requiredMethods {
        if !contains(handlerContent, method) {
            return fmt.Errorf("required health endpoint method %s not found", method)
        }
    }
    
    return nil
}

func theHealthEndpointsShouldRespond(ctx context.Context) error {
    testCtx := getTestContext(ctx)
    
    // Start the server
    originalDir, err := os.Getwd()
    if err != nil {
        return err
    }
    defer os.Chdir(originalDir)
    
    if err := os.Chdir(testCtx.ProjectPath); err != nil {
        return err
    }
    
    cmd := exec.Command("./bin/test-service")
    if err := cmd.Start(); err != nil {
        return fmt.Errorf("failed to start server: %w", err)
    }
    
    testCtx.ServerProcess = cmd.Process
    defer func() {
        if testCtx.ServerProcess != nil {
            testCtx.ServerProcess.Kill()
        }
    }()
    
    // Wait for server to start
    time.Sleep(3 * time.Second)
    
    // Test health endpoints
    endpoints := []string{
        "http://localhost:8080/health",
        "http://localhost:8080/health/time",
        "http://localhost:8080/health/ready",
        "http://localhost:8080/health/live",
        "http://localhost:8080/health/startup",
    }
    
    for _, endpoint := range endpoints {
        resp, err := http.Get(endpoint)
        if err != nil {
            return fmt.Errorf("failed to call endpoint %s: %w", endpoint, err)
        }
        resp.Body.Close()
        
        if resp.StatusCode != http.StatusOK {
            return fmt.Errorf("endpoint %s returned status %d, expected 200", endpoint, resp.StatusCode)
        }
    }
    
    return nil
}
```

### 3. Feature Files
Create comprehensive feature files covering all user scenarios:

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

  Scenario: Generate project with custom output directory
    When I run "template-health-endpoint generate --name my-service --tier basic --module github.com/org/my-service --output custom-dir"
    Then a new project should be created in "./custom-dir"
    And the project should compile successfully

  Scenario: Generate project with dry run
    When I run "template-health-endpoint generate --name my-service --tier basic --module github.com/org/my-service --dry-run"
    Then the command should succeed
    And the output should mention "dry run"
    And no project directory should be created
```

```gherkin
# features/static_templates.feature
Feature: Static Template Usage
  As a developer
  I want to use static templates directly
  So that I can manually customize templates before generation

  Background:
    Given I have the template-health-endpoint CLI installed

  Scenario: List available templates
    When I run "template-health-endpoint template list"
    Then the command should succeed
    And the output should list all available template tiers
    And each tier should show its description and features

  Scenario: Validate template integrity
    When I run "template-health-endpoint template validate"
    Then the command should succeed
    And all templates should pass validation
    And the output should confirm template integrity

  Scenario: Generate from static template
    When I run "template-health-endpoint template from-static --name static-test --tier basic --module github.com/test/static-test"
    Then a new project should be created in "./static-test"
    And the project should compile successfully
    And the project should have all required health endpoints
```

```gherkin
# features/project_migration.feature
Feature: Project Migration
  As a developer
  I want to migrate my project between template tiers
  So that I can add more features as my service evolves

  Background:
    Given I have the template-health-endpoint CLI installed
    And I have a basic tier project called "migration-test"
    And the project is compiled and working

  Scenario: Migrate from basic to intermediate
    When I run "template-health-endpoint migrate --project ./migration-test --to intermediate"
    Then the migration should succeed
    And the project should be at intermediate tier
    And dependency health check endpoints should be available
    And basic OpenTelemetry should be configured
    And the project should still compile and run

  Scenario: Migrate with dry run
    When I run "template-health-endpoint migrate --project ./migration-test --to advanced --dry-run"
    Then the command should succeed
    And the output should show migration plan
    And the output should mention "dry run"
    And the project should remain unchanged

  Scenario: Invalid migration path
    When I run "template-health-endpoint migrate --project ./migration-test --to invalid"
    Then the command should fail
    And the error message should mention "invalid tier"
    And the error message should list valid tiers
```

### 4. Performance and Load Testing Scenarios
Include performance-focused BDD scenarios:

```gherkin
# features/performance.feature
Feature: Performance Requirements
  As a developer
  I want the CLI to perform well
  So that I can use it efficiently in my workflow

  Scenario: CLI commands should complete quickly
    Given I have the template-health-endpoint CLI installed
    When I run "template-health-endpoint template list"
    Then the command should complete in less than 2 seconds

  Scenario: Project generation should be reasonably fast
    Given I have the template-health-endpoint CLI installed
    When I run "template-health-endpoint generate --name perf-test --tier basic --module github.com/test/perf-test"
    Then the command should complete in less than 10 seconds
    And a new project should be created in "./perf-test"

  Scenario: Generated service should handle load
    Given I have a generated basic tier project
    And the service is running
    When I send 100 concurrent requests to "/health"
    Then all requests should succeed
    And the average response time should be less than 100ms
```

### 5. Integration Testing Scenarios
Test integration with external systems:

```gherkin
# features/kubernetes_integration.feature
Feature: Kubernetes Integration
  As a DevOps engineer
  I want generated projects to work with Kubernetes
  So that I can deploy them in production

  Background:
    Given I have the template-health-endpoint CLI installed
    And I have kubectl available
    And I have a Kubernetes cluster available

  Scenario: Generated Kubernetes manifests are valid
    Given I have generated a basic tier project called "k8s-test"
    When I validate the Kubernetes manifests with kubectl
    Then all manifests should be valid YAML
    And all manifests should pass kubectl validation

  Scenario: Health probes work in Kubernetes
    Given I have deployed a generated project to Kubernetes
    When Kubernetes performs health checks
    Then the liveness probe should succeed
    And the readiness probe should succeed
    And the startup probe should succeed
    And the pod should be marked as ready
```

## BDD Test Execution Framework

### 1. Test Runner Setup
```go
// Makefile
.PHONY: test-bdd
test-bdd:
	@echo "Running BDD tests..."
	go test -v ./features/...

.PHONY: test-bdd-verbose
test-bdd-verbose:
	@echo "Running BDD tests with verbose output..."
	go test -v ./features/... -godog.format=pretty

.PHONY: test-bdd-json
test-bdd-json:
	@echo "Running BDD tests with JSON output..."
	go test -v ./features/... -godog.format=cucumber > bdd-results.json
```

### 2. CI/CD Integration
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
      run: go mod download
    
    - name: Build CLI
      run: go build -o bin/template-health-endpoint cmd/generator/main.go
    
    - name: Run BDD tests
      run: make test-bdd-verbose
    
    - name: Generate BDD report
      run: make test-bdd-json
    
    - name: Upload BDD results
      uses: actions/upload-artifact@v3
      with:
        name: bdd-results
        path: bdd-results.json
```

### 3. Test Data Management
```go
// features/testdata.go
package features

import (
    "os"
    "path/filepath"
)

type TestDataManager struct {
    TempDir string
}

func NewTestDataManager() *TestDataManager {
    tempDir, _ := os.MkdirTemp("", "bdd-test-*")
    return &TestDataManager{TempDir: tempDir}
}

func (tdm *TestDataManager) CreateTestProject(name, tier string) (string, error) {
    projectPath := filepath.Join(tdm.TempDir, name)
    
    // Generate test project
    cmd := exec.Command("./bin/template-health-endpoint", "generate",
        "--name", name,
        "--tier", tier,
        "--module", fmt.Sprintf("github.com/test/%s", name),
        "--output", projectPath)
    
    return projectPath, cmd.Run()
}

func (tdm *TestDataManager) Cleanup() {
    os.RemoveAll(tdm.TempDir)
}
```

## Success Criteria

### BDD Implementation Checklist
- [ ] Feature files cover all user scenarios
- [ ] Step definitions implement all Gherkin steps
- [ ] Error scenarios are thoroughly tested
- [ ] Performance requirements are validated
- [ ] Integration scenarios are covered
- [ ] CI/CD pipeline includes BDD tests
- [ ] Test reports are generated and accessible
- [ ] Test data is properly managed and cleaned up

### Coverage Requirements
- [ ] All CLI commands have BDD scenarios
- [ ] All template tiers have generation scenarios
- [ ] All migration paths have BDD coverage
- [ ] Error conditions are tested with BDD
- [ ] Performance requirements are validated
- [ ] Integration points are tested

## Related Patterns
- Behavior-Driven Development (BDD)
- Acceptance Test-Driven Development (ATDD)
- Specification by Example
- Living Documentation
- Test Automation Pyramid
- Continuous Integration/Continuous Deployment
