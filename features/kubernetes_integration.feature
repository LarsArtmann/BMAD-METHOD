Feature: Kubernetes Integration
  As a DevOps engineer
  I want generated projects to work seamlessly with Kubernetes
  So that I can deploy health endpoints in production environments

  Background:
    Given the template-health-endpoint CLI is available
    And I have a Kubernetes cluster available
    And kubectl is configured and working

  Scenario: Basic Kubernetes deployment
    Given I have generated a "basic" tier project with Kubernetes enabled
    When I apply the Kubernetes manifests
    Then the deployment should be created successfully
    And the service should be created successfully
    And the pods should be running and healthy
    And the health endpoints should be accessible through the service

  Scenario: Health check configuration
    Given I have generated a project with Kubernetes health checks
    When I apply the Kubernetes manifests
    Then the deployment should have proper liveness probes
    And the deployment should have proper readiness probes
    And the deployment should have proper startup probes
    And Kubernetes should use the health endpoints for pod management

  Scenario: Service discovery and networking
    Given I have deployed a health endpoint service
    When I create a test pod in the same namespace
    Then the test pod should be able to reach the health service
    And DNS resolution should work correctly
    And the service endpoints should be properly configured

  Scenario: Resource limits and requests
    Given I have generated a project with resource specifications
    When I apply the Kubernetes manifests
    Then the pods should have appropriate resource requests
    And the pods should have appropriate resource limits
    And the pods should not exceed their resource allocations

  Scenario: ConfigMap integration
    Given I have generated a project with ConfigMap support
    When I apply the Kubernetes manifests
    Then the ConfigMap should be created with correct configuration
    And the deployment should mount the ConfigMap correctly
    And the application should use the configuration from the ConfigMap

  Scenario: Secret management
    Given I have generated an enterprise project with secrets
    When I create the required secrets
    And I apply the Kubernetes manifests
    Then the deployment should mount the secrets correctly
    And the application should use the secrets for authentication
    And the secrets should not be exposed in logs or environment

  Scenario: Horizontal Pod Autoscaling
    Given I have generated a project with HPA configuration
    When I apply the Kubernetes manifests including HPA
    Then the HPA should be created successfully
    And the HPA should monitor the correct metrics
    And the HPA should scale pods based on load

  Scenario: ServiceMonitor for Prometheus
    Given I have generated an advanced project with metrics
    When I apply the Kubernetes manifests
    Then the ServiceMonitor should be created correctly
    And Prometheus should discover the service endpoints
    And metrics should be scraped successfully

  Scenario: Ingress configuration
    Given I have generated a project with Ingress enabled
    When I apply the Kubernetes manifests
    Then the Ingress should be created successfully
    And the Ingress should route traffic to the service
    And the health endpoints should be accessible through the Ingress

  Scenario: Network policies
    Given I have generated an enterprise project with network policies
    When I apply the Kubernetes manifests
    Then the NetworkPolicy should be created correctly
    And traffic should be restricted according to the policy
    And authorized traffic should still flow correctly

  Scenario: Pod security context
    Given I have generated an enterprise project with security context
    When I apply the Kubernetes manifests
    Then the pods should run with the correct security context
    And the pods should not run as root
    And the pods should have appropriate security constraints

  Scenario: Multi-environment deployment
    Given I have generated an enterprise project with multiple environments
    When I apply manifests for different environments
    Then each environment should have its own namespace
    And each environment should have appropriate configurations
    And environments should be isolated from each other

  Scenario: Rolling updates
    Given I have a deployed health endpoint service
    When I update the deployment with a new image
    Then the rolling update should proceed smoothly
    And there should be no service downtime
    And health checks should pass throughout the update

  Scenario: Persistent volume claims
    Given I have generated a project that requires persistent storage
    When I apply the Kubernetes manifests
    Then the PVC should be created successfully
    And the PVC should be bound to a persistent volume
    And the application should be able to write to the persistent storage

  Scenario: RBAC configuration
    Given I have generated an enterprise project with RBAC
    When I apply the Kubernetes manifests
    Then the ServiceAccount should be created correctly
    And the Role and RoleBinding should be configured properly
    And the pods should have the correct permissions

  Scenario: Cluster-level monitoring integration
    Given I have deployed multiple health endpoint services
    When I check the cluster monitoring
    Then all services should be visible in the monitoring dashboard
    And health metrics should be aggregated correctly
    And alerts should be configured for service health

  Scenario: Disaster recovery and backup
    Given I have a deployed health endpoint service with persistent data
    When I simulate a node failure
    Then the service should be rescheduled to another node
    And persistent data should be preserved
    And the service should recover automatically

  Scenario: Load testing in Kubernetes
    Given I have deployed a health endpoint service
    When I run load tests against the service
    Then the service should handle the load appropriately
    And Kubernetes should scale resources as needed
    And the health endpoints should remain responsive

  Scenario: Logging and observability
    Given I have deployed an advanced tier service
    When I check the logs and metrics
    Then logs should be properly formatted and accessible
    And metrics should be exported correctly
    And distributed tracing should work across pod boundaries
