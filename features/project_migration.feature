Feature: Project Migration Between Tiers
  As a developer
  I want to migrate my project between template tiers
  So that I can add more features as my service evolves

  Background:
    Given the template-health-endpoint CLI is available
    And I have a clean working directory

  Scenario: Migrate from basic to intermediate
    Given I have a basic tier project named "migrate-test"
    When I run "template-health-endpoint migrate --to intermediate" in the project directory
    Then the command should succeed
    And the project should be upgraded to intermediate tier
    And dependency health check endpoints should be available
    And the project should compile successfully
    And all existing functionality should still work

  Scenario: Migrate from intermediate to advanced
    Given I have an intermediate tier project named "migrate-advanced"
    When I run "template-health-endpoint migrate --to advanced" in the project directory
    Then the command should succeed
    And the project should be upgraded to advanced tier
    And OpenTelemetry observability should be enabled
    And CloudEvents integration should be available
    And the project should compile successfully

  Scenario: Migrate from advanced to enterprise
    Given I have an advanced tier project named "migrate-enterprise"
    When I run "template-health-endpoint migrate --to enterprise" in the project directory
    Then the command should succeed
    And the project should be upgraded to enterprise tier
    And mTLS security features should be available
    And RBAC configuration should be present
    And audit logging should be enabled
    And multi-environment configs should be created

  Scenario: Migration with backup
    Given I have a basic tier project named "backup-test"
    When I run "template-health-endpoint migrate --to intermediate --backup" in the project directory
    Then the command should succeed
    And a backup directory should be created
    And the backup should contain the original project files
    And the project should be successfully migrated

  Scenario: Dry run migration
    Given I have a basic tier project named "dry-run-migrate"
    When I run "template-health-endpoint migrate --to advanced --dry-run" in the project directory
    Then the command should succeed
    And no files should be modified
    And the migration plan should be displayed
    And the plan should show what changes would be made

  Scenario: Migration path validation
    Given I have an intermediate tier project named "path-test"
    When I run "template-health-endpoint migrate --to enterprise" in the project directory
    Then the command should succeed
    And the migration path should be "intermediate → advanced → enterprise"
    And all intermediate steps should be applied correctly

  Scenario: Downgrade migration with warning
    Given I have an advanced tier project named "downgrade-test"
    When I run "template-health-endpoint migrate --to basic" in the project directory
    Then the command should display a downgrade warning
    And I should be prompted for confirmation
    When I confirm the downgrade
    Then the project should be downgraded to basic tier
    And advanced features should be removed
    And the basic functionality should still work

  Scenario: Migration preserves customizations
    Given I have a basic tier project named "custom-test"
    And I have made custom modifications to the health handler
    When I run "template-health-endpoint migrate --to intermediate" in the project directory
    Then the command should succeed
    And my custom modifications should be preserved
    And the new intermediate features should be added
    And the project should compile successfully

  Scenario: Migration updates dependencies
    Given I have a basic tier project named "deps-test"
    When I run "template-health-endpoint migrate --to advanced" in the project directory
    Then the command should succeed
    And the go.mod file should include OpenTelemetry dependencies
    And the go.mod file should include CloudEvents dependencies
    And "go mod tidy" should run successfully

  Scenario: Migration error handling
    When I run "template-health-endpoint migrate --to invalid-tier" in a project directory
    Then the command should fail
    And the error message should mention "invalid target tier"

  Scenario: Migration from non-template project
    Given I have a regular Go project that is not from a template
    When I run "template-health-endpoint migrate --to intermediate" in the project directory
    Then the command should fail
    And the error message should mention "not a template project"

  Scenario: Force migration without confirmation
    Given I have an advanced tier project named "force-test"
    When I run "template-health-endpoint migrate --to basic --force" in the project directory
    Then the command should succeed without prompting
    And the project should be downgraded to basic tier

  Scenario: Migration updates project metadata
    Given I have a basic tier project named "metadata-test"
    When I run "template-health-endpoint migrate --to intermediate" in the project directory
    Then the command should succeed
    And the .template-metadata.yaml file should be updated
    And the metadata should reflect the new tier
    And the migration timestamp should be recorded
