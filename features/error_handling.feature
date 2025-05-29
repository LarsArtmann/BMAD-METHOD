Feature: Error Handling and Validation
  As a developer
  I want clear error messages and proper validation
  So that I can quickly understand and fix issues

  Background:
    Given the template-health-endpoint CLI is available

  Scenario: Invalid tier specification
    When I run "template-health-endpoint generate --name test-service --tier invalid-tier"
    Then the command should fail with exit code 1
    And the error message should contain "invalid tier 'invalid-tier'"
    And the error message should list valid tiers: "basic, intermediate, advanced, enterprise"

  Scenario: Missing required project name
    When I run "template-health-endpoint generate --tier basic"
    Then the command should fail with exit code 1
    And the error message should contain "project name is required"

  Scenario: Invalid project name characters
    When I run "template-health-endpoint generate --name 'invalid name with spaces' --tier basic"
    Then the command should fail with exit code 1
    And the error message should contain "invalid project name"
    And the error message should mention valid naming conventions

  Scenario: Invalid Go module name
    When I run "template-health-endpoint generate --name test-service --tier basic --go-module 'invalid module name'"
    Then the command should fail with exit code 1
    And the error message should contain "invalid Go module name"

  Scenario: Output directory permission denied
    Given I have a directory with no write permissions
    When I run "template-health-endpoint generate --name test-service --tier basic --output /no-write-permission"
    Then the command should fail with exit code 1
    And the error message should contain "permission denied"

  Scenario: Output directory already exists with files
    Given I have a directory "existing-project" with files
    When I run "template-health-endpoint generate --name existing-project --tier basic"
    Then the command should fail with exit code 1
    And the error message should contain "directory already exists and is not empty"
    And the error message should suggest using --force or a different directory

  Scenario: Invalid feature specification
    When I run "template-health-endpoint generate --name test-service --tier basic --features invalid-feature"
    Then the command should fail with exit code 1
    And the error message should contain "invalid feature 'invalid-feature'"
    And the error message should list valid features

  Scenario: Incompatible feature combination
    When I run "template-health-endpoint generate --name test-service --tier basic --features enterprise-security"
    Then the command should fail with exit code 1
    And the error message should contain "feature 'enterprise-security' is not available in 'basic' tier"
    And the error message should suggest upgrading to a higher tier

  Scenario: Network connectivity issues
    Given the network is unavailable
    When I run "template-health-endpoint generate --name test-service --tier basic"
    Then the command should handle network errors gracefully
    And provide helpful error messages about connectivity
    And suggest offline alternatives if available

  Scenario: Insufficient disk space
    Given the disk is full
    When I run "template-health-endpoint generate --name test-service --tier basic"
    Then the command should fail with exit code 1
    And the error message should contain "insufficient disk space"

  Scenario: Template corruption detection
    Given the template files are corrupted
    When I run "template-health-endpoint generate --name test-service --tier basic"
    Then the command should fail with exit code 1
    And the error message should contain "template validation failed"
    And the error message should suggest reinstalling or updating

  Scenario: Migration from unsupported project
    Given I have a project that is not a template project
    When I run "template-health-endpoint migrate --to intermediate" in the project directory
    Then the command should fail with exit code 1
    And the error message should contain "not a template-generated project"
    And the error message should suggest using generate instead

  Scenario: Migration to same tier
    Given I have an intermediate tier project
    When I run "template-health-endpoint migrate --to intermediate" in the project directory
    Then the command should succeed with exit code 0
    And the output should contain "already at tier 'intermediate'"
    And no changes should be made

  Scenario: Update on non-template project
    Given I have a regular Go project
    When I run "template-health-endpoint update" in the project directory
    Then the command should fail with exit code 1
    And the error message should contain "not a template-generated project"

  Scenario: Validation command with invalid schema
    Given I have an invalid TypeSpec schema file
    When I run "template-health-endpoint validate --schemas invalid-schema.tsp"
    Then the command should fail with exit code 1
    And the error message should contain schema validation errors
    And the error message should point to the specific line and column

  Scenario: Help and usage information
    When I run "template-health-endpoint --help"
    Then the command should succeed with exit code 0
    And the output should contain usage information
    And the output should list all available commands

  Scenario: Version information
    When I run "template-health-endpoint --version"
    Then the command should succeed with exit code 0
    And the output should contain the version number

  Scenario: Verbose output mode
    When I run "template-health-endpoint generate --name test-service --tier basic --verbose"
    Then the command should succeed
    And the output should contain detailed progress information
    And the output should show file creation steps

  Scenario: Configuration file errors
    Given I have an invalid configuration file
    When I run "template-health-endpoint generate --config invalid-config.yaml --name test-service --tier basic"
    Then the command should fail with exit code 1
    And the error message should contain "invalid configuration file"
    And the error message should specify what is wrong with the configuration

  Scenario: Graceful interruption handling
    Given I start a long-running generate command
    When I interrupt the command with Ctrl+C
    Then the command should handle the interruption gracefully
    And any partial files should be cleaned up
    And the output should indicate the operation was cancelled
