Feature: Performance and Scalability
  As a developer
  I want the template generator to perform efficiently
  So that I can generate projects quickly even for large templates

  Background:
    Given the template-health-endpoint CLI is available
    And I have a clean working directory

  Scenario: Basic project generation performance
    When I run "template-health-endpoint generate --name perf-basic --tier basic"
    Then the command should complete in less than 10 seconds
    And the generated project should be functional
    And memory usage should be reasonable

  Scenario: Enterprise project generation performance
    When I run "template-health-endpoint generate --name perf-enterprise --tier enterprise"
    Then the command should complete in less than 30 seconds
    And the generated project should be functional
    And all enterprise features should be properly configured

  Scenario: Multiple project generation
    When I generate 5 projects simultaneously with different names
    Then all commands should complete successfully
    And the total time should be reasonable
    And there should be no resource conflicts

  Scenario: Large project template processing
    Given I have a template with many files and complex structure
    When I run "template-health-endpoint generate --name large-project --tier advanced"
    Then the command should handle the large template efficiently
    And memory usage should remain stable
    And the generation should complete successfully

  Scenario: Template validation performance
    When I run "template-health-endpoint validate --schemas pkg/schemas/"
    Then the validation should complete in less than 5 seconds
    And all schemas should be validated correctly
    And memory usage should be minimal

  Scenario: Migration performance
    Given I have a basic tier project with many files
    When I run "template-health-endpoint migrate --to enterprise" in the project directory
    Then the migration should complete in less than 20 seconds
    And all files should be processed correctly
    And the migrated project should be functional

  Scenario: Update performance with many changes
    Given I have a project that needs many updates
    When I run "template-health-endpoint update" in the project directory
    Then the update should complete in less than 15 seconds
    And all changes should be applied correctly
    And the updated project should be functional

  Scenario: Concurrent CLI operations
    When I run multiple CLI commands concurrently
    Then all commands should complete without interference
    And there should be no file locking issues
    And all generated projects should be valid

  Scenario: Memory usage during generation
    When I run "template-health-endpoint generate --name memory-test --tier enterprise"
    Then the peak memory usage should be less than 100MB
    And memory should be released properly after completion
    And there should be no memory leaks

  Scenario: CPU usage optimization
    When I run "template-health-endpoint generate --name cpu-test --tier advanced"
    Then the CPU usage should be reasonable
    And the command should not consume excessive CPU resources
    And the system should remain responsive

  Scenario: File I/O performance
    When I generate a project with many template files
    Then file operations should be efficient
    And there should be no unnecessary file reads or writes
    And the file system should not be overwhelmed

  Scenario: Network operation performance
    Given the CLI needs to download dependencies or templates
    When I run generation commands
    Then network operations should be optimized
    And there should be appropriate caching
    And timeouts should be reasonable

  Scenario: Cleanup performance
    Given I have generated multiple test projects
    When I clean up the test projects
    Then the cleanup should be fast and thorough
    And no temporary files should be left behind
    And system resources should be properly released

  Scenario: Scalability with project size
    When I generate projects of increasing complexity
    Then the generation time should scale reasonably
    And resource usage should scale predictably
    And there should be no performance cliffs

  Scenario: Template caching effectiveness
    When I generate multiple projects with the same tier
    Then template caching should improve performance
    And subsequent generations should be faster
    And cache invalidation should work correctly

  Scenario: Progress reporting accuracy
    When I run a long-running generation command with verbose output
    Then progress reporting should be accurate
    And progress updates should be timely
    And the user should have clear feedback on operation status
