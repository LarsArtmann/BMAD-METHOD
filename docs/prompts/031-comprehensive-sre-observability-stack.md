# Comprehensive SRE Observability Stack Implementation

## Objective
Implement a complete Site Reliability Engineering (SRE) observability stack including metrics, tracing, logging, alerting, dashboards, and SLI/SLO management for production-ready applications across all complexity tiers.

## Context
This prompt focuses on creating enterprise-grade observability infrastructure that follows SRE best practices. The system should provide comprehensive monitoring, alerting, and performance tracking suitable for production environments.

## Task Description

### Phase 1: OpenTelemetry Integration
1. **Advanced Metrics System**
   - Create `template-health/templates/go-metrics-advanced.tmpl` with OpenTelemetry metrics
   - Implement HTTP request tracking, health check monitoring, dependency monitoring
   - Build runtime metrics collection (memory, goroutines, GC cycles)
   - Support multiple exporters: Prometheus, OTLP, stdout for flexibility

2. **Distributed Tracing**
   - Create `template-health/templates/go-tracing-advanced.tmpl` with OpenTelemetry tracing
   - Implement context propagation across service boundaries
   - Build specialized spans for HTTP requests, database operations, external services
   - Support Jaeger, OTLP, and custom exporters with sampling configuration

3. **Structured Logging**
   - Create `template-health/templates/go-structured-logging.tmpl` with slog integration
   - Implement trace correlation for log-trace integration
   - Build specialized logging for HTTP requests, business events, security events
   - Support JSON and text formats with configurable outputs

### Phase 2: SRE Monitoring Infrastructure
1. **Prometheus Alerting Rules**
   - Create `template-health/templates/sre-prometheus-alerts.tmpl` with comprehensive alerting
   - Implement service availability, error rate, response time, and resource usage alerts
   - Build tier-specific alerting (basic, intermediate, advanced, enterprise)
   - Include SLO-based alerting with error budget tracking and burn rate analysis

2. **Grafana Dashboards**
   - Create `template-health/templates/sre-grafana-dashboard.tmpl` with visual monitoring
   - Build service health overview, request rate metrics, response time percentiles
   - Implement SLO availability gauges and error budget visualization
   - Include tier-specific panels for dependency health, resource usage, security events

3. **SLI/SLO Configuration**
   - Create `template-health/templates/sre-sli-slo-config.tmpl` with comprehensive SLO management
   - Define availability, latency, quality, and dependency health SLIs
   - Implement error budget policies with automated actions
   - Build compliance and audit framework for enterprise requirements

### Phase 3: Tier-Specific Implementation
1. **Basic Tier**
   - Essential health checks and basic logging
   - Minimal resource overhead with core functionality
   - Simple alerting for critical issues only

2. **Intermediate Tier**
   - Add dependency monitoring and structured logging
   - Include basic metrics collection and simple dashboards
   - Implement warning-level alerting for proactive monitoring

3. **Advanced Tier**
   - Full observability stack with distributed tracing
   - Advanced metrics including runtime and performance monitoring
   - Comprehensive dashboards with resource usage tracking

4. **Enterprise Tier**
   - Complete SRE stack with security event monitoring
   - Audit logging and compliance reporting
   - Advanced SLO management with automated error budget policies
   - Full integration with enterprise monitoring tools

### Phase 4: Integration and Automation
1. **Template Integration**
   - Seamless integration with existing template generation system
   - Configuration-driven feature selection based on project tier
   - Automatic dependency injection and setup

2. **CI/CD Integration**
   - Automated monitoring stack deployment
   - Configuration validation and testing
   - Performance regression detection

## Technical Requirements
- **OpenTelemetry**: Latest SDK with full instrumentation support
- **Prometheus**: Production-ready alerting rules and recording rules
- **Grafana**: Professional dashboards with proper visualization
- **Performance**: Minimal overhead (<5% performance impact)
- **Security**: Secure defaults and enterprise compliance support
- **Configuration**: Environment-based configuration with secure secrets management

## Expected Deliverables
1. Complete OpenTelemetry integration templates for metrics, tracing, and logging
2. Production-ready Prometheus alerting rules with SLO-based monitoring
3. Professional Grafana dashboards with comprehensive visualizations
4. SLI/SLO configuration framework with error budget management
5. Tier-specific implementations with appropriate feature sets
6. Integration with existing template system and CLI tools
7. Documentation and best practices guide

## Success Criteria
- Generated services have comprehensive observability out-of-the-box
- Monitoring stack follows SRE best practices and industry standards
- Alerting provides actionable insights without excessive noise
- Dashboards enable quick troubleshooting and performance analysis
- SLO framework supports business reliability requirements
- System scales from development to enterprise production environments

## Template Structure
```
template-health/templates/
├── go-metrics-advanced.tmpl           # OpenTelemetry metrics
├── go-tracing-advanced.tmpl           # Distributed tracing
├── go-structured-logging.tmpl         # Structured logging with trace correlation
├── sre-prometheus-alerts.tmpl         # Prometheus alerting rules
├── sre-grafana-dashboard.tmpl         # Grafana dashboards
└── sre-sli-slo-config.tmpl           # SLI/SLO configuration
```

## Key Features
- **Zero-config observability** for new projects
- **Production-ready defaults** with security best practices
- **Tier-appropriate complexity** scaling from basic to enterprise
- **Industry standard tools** (OpenTelemetry, Prometheus, Grafana)
- **SRE methodology** with proper SLI/SLO implementation
- **Extensible architecture** for custom monitoring requirements

## Related Files
- Template generation system in `/pkg/generator/`
- Configuration management in `/pkg/config/`
- Health template system in `/template-health/`
- TypeSpec schemas for API definitions