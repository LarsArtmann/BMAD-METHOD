Feature: Template Generation
  As a developer
  I want to generate projects from different template tiers
  So that I can choose the right complexity level for my needs

  Background:
    Given the template-health-endpoint CLI is available
    And I have a clean working directory

  Scenario Outline: Generate project from different tiers
    When I run "template-health-endpoint generate --name <project_name> --tier <tier> --output <output_dir>"
    Then the command should succeed
    And a new project should be created in "<output_dir>/<project_name>"
    And the project should have the correct directory structure for "<tier>" tier
    And the project should compile successfully
    And all health endpoints should respond correctly

    Examples:
      | project_name    | tier         | output_dir |
      | basic-service   | basic        | ./test-output |
      | inter-service   | intermediate | ./test-output |
      | advanced-service| advanced     | ./test-output |
      | enterprise-svc  | enterprise   | ./test-output |

  Scenario: Generate project with custom features
    When I run "template-health-endpoint generate --name custom-service --tier intermediate --features typescript,kubernetes"
    Then the command should succeed
    And the project should include TypeScript client SDK
    And the project should include Kubernetes manifests
    And the project should compile successfully

  Scenario: Generate project with custom Go module
    When I run "template-health-endpoint generate --name my-service --tier basic --go-module github.com/myorg/my-service"
    Then the command should succeed
    And the go.mod file should contain "module github.com/myorg/my-service"
    And the project should compile successfully

  Scenario: Validate generated project structure
    Given I have generated a "intermediate" tier project named "test-service"
    Then the project should contain these files:
      | file                                    |
      | go.mod                                  |
      | cmd/server/main.go                      |
      | internal/config/config.go               |
      | internal/server/server.go               |
      | internal/handlers/health.go             |
      | internal/handlers/dependencies.go       |
      | internal/models/health.go               |
      | deployments/kubernetes/deployment.yaml  |
      | deployments/kubernetes/service.yaml     |
      | Dockerfile                              |
      | README.md                               |

  Scenario: Test health endpoints in generated project
    Given I have generated a "basic" tier project named "health-test"
    And I have started the server
    When I make a GET request to "/health"
    Then the response status should be 200
    And the response should contain valid health status
    When I make a GET request to "/health/ready"
    Then the response status should be 200
    When I make a GET request to "/health/live"
    Then the response status should be 200

  Scenario: Validate TypeScript client generation
    Given I have generated an "advanced" tier project with TypeScript enabled
    Then the client/typescript directory should exist
    And the TypeScript client should compile successfully
    And the generated types should match the TypeSpec definitions

  Scenario: Validate Kubernetes manifests
    Given I have generated a project with Kubernetes enabled
    Then the Kubernetes manifests should be valid YAML
    And kubectl should validate the manifests successfully
    And the manifests should include proper health check configurations

  Scenario: Error handling for invalid parameters
    When I run "template-health-endpoint generate --name invalid-service --tier nonexistent"
    Then the command should fail
    And the error message should mention "invalid tier"

  Scenario: Dry run mode
    When I run "template-health-endpoint generate --name dry-run-service --tier basic --dry-run"
    Then the command should succeed
    And no files should be created
    And the output should show what would be generated
