# Kubernetes Cloud-Native Integration

## Prompt Name: Kubernetes Cloud-Native Integration

## Description
Create comprehensive Kubernetes-ready applications with health probes, observability, and cloud-native patterns.

## When to Use
- Building cloud-native applications
- Need Kubernetes deployment configurations
- Want comprehensive health monitoring
- Building microservices for container orchestration

## Prompt

```
Create a comprehensive Kubernetes-native application for [APPLICATION_NAME] with full cloud-native integration:

**Health Monitoring:**
- Implement /health, /health/ready, /health/live endpoints
- Create comprehensive health check framework
- Support dependency health validation
- Include performance metrics and timing
- Add graceful degradation patterns

**Kubernetes Integration:**
- Generate Deployment with proper resource limits
- Create Service and ConfigMap manifests
- Configure liveness, readiness, and startup probes
- Include ServiceMonitor for Prometheus
- Add Ingress configuration with health routing

**Observability Stack:**
- Integrate OpenTelemetry for tracing and metrics
- Implement Server Timing API for performance
- Add structured logging with correlation IDs
- Support CloudEvents for event-driven monitoring
- Include Prometheus metrics exposition

**Container Optimization:**
- Multi-stage Docker builds for minimal images
- Non-root user security patterns
- Health check integration in containers
- Proper signal handling for graceful shutdown
- Resource optimization and limits

**Deployment Patterns:**
- Rolling update strategies
- Blue-green deployment support
- Canary deployment configurations
- Auto-scaling based on metrics
- Network policies and security

**Monitoring and Alerting:**
- Grafana dashboard templates
- Prometheus alerting rules
- Log aggregation patterns
- Distributed tracing setup
- Performance monitoring

**Security Patterns:**
- Pod security standards
- Network segmentation
- Secret management
- RBAC configurations
- Security scanning integration

Ensure all components work together seamlessly and follow Kubernetes best practices.
```

## Expected Outcomes
- Kubernetes-ready application
- Complete observability integration
- Health monitoring framework
- Security-hardened deployment
- Monitoring and alerting setup

## Success Criteria
- Application deploys successfully to Kubernetes
- Health probes work correctly
- Observability data flows properly
- Security policies are enforced
- Monitoring dashboards show metrics
