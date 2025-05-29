Feature: Project Update
  As a developer
  I want to update my existing project to newer template versions
  So that I can get bug fixes and improvements while preserving my customizations

  Background:
    Given the template-health-endpoint CLI is available
    And I have a clean working directory

  Scenario: Update project to latest template version
    Given I have a project generated from an older template version
    When I run "template-health-endpoint update" in the project directory
    Then the command should succeed
    And the project should be updated to the latest template version
    And my customizations should be preserved
    And the project should compile successfully

  Scenario: Selective component update
    Given I have a project that needs updates
    When I run "template-health-endpoint update --components kubernetes,docker" in the project directory
    Then the command should succeed
    And only the Kubernetes and Docker files should be updated
    And other components should remain unchanged
    And the project should compile successfully

  Scenario: Update with diff preview
    Given I have a project that needs updates
    When I run "template-health-endpoint update --show-diff" in the project directory
    Then the command should succeed
    And the diff of changes should be displayed
    And I should see what files will be modified
    And I should see the specific changes to be made

  Scenario: Dry run update
    Given I have a project that needs updates
    When I run "template-health-endpoint update --dry-run" in the project directory
    Then the command should succeed
    And no files should be modified
    And the update plan should be displayed
    And I should see what would be updated

  Scenario: Update with backup
    Given I have a project that needs updates
    When I run "template-health-endpoint update --backup" in the project directory
    Then the command should succeed
    And a backup directory should be created
    And the backup should contain the original files
    And the project should be successfully updated

  Scenario: Force update without confirmation
    Given I have a project that needs updates
    When I run "template-health-endpoint update --force" in the project directory
    Then the command should succeed without prompting
    And the project should be updated
    And the project should compile successfully

  Scenario: Update already up-to-date project
    Given I have a project that is already up-to-date
    When I run "template-health-endpoint update" in the project directory
    Then the command should succeed
    And the output should indicate "already up to date"
    And no files should be modified

  Scenario: Update preserves custom configurations
    Given I have a project with custom configuration files
    And the project needs template updates
    When I run "template-health-endpoint update" in the project directory
    Then the command should succeed
    And my custom configuration should be preserved
    And the template updates should be applied
    And the project should compile successfully

  Scenario: Update handles merge conflicts
    Given I have a project with conflicting customizations
    And the project needs template updates
    When I run "template-health-endpoint update" in the project directory
    Then the command should detect conflicts
    And I should be prompted to resolve conflicts
    And conflict markers should be added to affected files
    And I should be able to manually resolve the conflicts

  Scenario: Update specific target directory
    Given I have a project in a custom directory
    When I run "template-health-endpoint update --target /path/to/project"
    Then the command should succeed
    And the project in the specified directory should be updated
    And the project should compile successfully

  Scenario: Update error handling for non-template project
    Given I have a regular Go project that is not from a template
    When I run "template-health-endpoint update" in the project directory
    Then the command should fail
    And the error message should mention "not a template project"

  Scenario: Update with custom components
    Given I have a project with custom components defined
    When I run "template-health-endpoint update --components custom-component" in the project directory
    Then the command should handle the custom component appropriately
    And the update should complete successfully

  Scenario: Update rollback on failure
    Given I have a project that needs updates
    And the update process will fail partway through
    When I run "template-health-endpoint update" in the project directory
    Then the command should fail gracefully
    And any partial changes should be rolled back
    And the project should be in its original state
    And an error message should explain what went wrong
