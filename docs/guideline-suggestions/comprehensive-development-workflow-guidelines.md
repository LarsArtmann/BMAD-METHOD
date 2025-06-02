# Comprehensive Development Workflow Guidelines for .augment-guidelines

## Overview
Based on extensive work with the BMAD Method project, these suggestions aim to improve the .augment-guidelines file to better support complex software development workflows, architectural decisions, and enterprise-grade implementations.

## Suggested Additions to .augment-guidelines

### 1. Architecture & Design Patterns
```markdown
## Architecture Guidelines

### Domain-Driven Design (DDD)
- Always identify core domain concepts before implementation
- Use ubiquitous language consistently across code and documentation
- Implement proper aggregate boundaries with clear business invariants
- Separate domain logic from infrastructure concerns

### Clean Architecture Principles
- Maintain dependency inversion: inner layers should not depend on outer layers
- Use interface segregation for loose coupling
- Implement repository patterns for data access abstraction
- Keep business logic independent of frameworks and external systems

### Event-Driven Architecture
- Use domain events for loose coupling between bounded contexts
- Implement event sourcing for audit trails and temporal queries
- Design idempotent event handlers for reliability
- Version events properly for backward compatibility
```

### 2. Feature Composition & Plugin Architecture
```markdown
## Feature Composition Guidelines

### Dynamic Feature Systems
- Design features as composable units with clear interfaces
- Implement dependency resolution with conflict detection
- Provide intelligent feature recommendations based on project context
- Support tier-based feature availability (basic, intermediate, advanced, enterprise)

### Plugin Architecture
- Define clear plugin interfaces with versioning support
- Implement plugin discovery and registration mechanisms
- Provide plugin lifecycle management (load, unload, update)
- Establish security boundaries for plugin execution
```

### 3. Enterprise-Grade Observability
```markdown
## Observability & SRE Guidelines

### OpenTelemetry Integration
- Always implement structured logging with trace correlation
- Use OpenTelemetry for metrics, tracing, and logging standardization
- Implement proper sampling strategies for production environments
- Ensure minimal performance overhead (<5% impact)

### SLO/SLI Management
- Define clear Service Level Indicators (SLIs) for critical user journeys
- Implement Service Level Objectives (SLOs) with error budget tracking
- Create actionable alerts based on burn rate analysis
- Establish incident response procedures with runbooks

### Monitoring Stack
- Use Prometheus for metrics collection with proper labeling
- Implement Grafana dashboards for operational visibility
- Create tier-appropriate alerting (basic to enterprise complexity)
- Include security event monitoring for enterprise deployments
```

### 4. Multi-Platform Distribution Strategy
```markdown
## Distribution & Packaging Guidelines

### Single-Line Installation
- Provide universal installation script with OS/architecture detection
- Support multiple package managers (Homebrew, NPM, Nix, Docker)
- Implement proper version management and rollback capabilities
- Ensure security verification (checksums, code signing)

### Cross-Platform Compatibility
- Build for all major platforms (Windows, macOS, Linux arm64/amd64)
- Use statically linked binaries for zero dependencies
- Optimize binary size (<50MB) with fast startup (<2s)
- Implement proper error handling for platform-specific issues

### Container & Cloud Integration
- Create optimized Docker images with security scanning
- Support Kubernetes deployment with Helm charts
- Integrate with major cloud providers (AWS, GCP, Azure)
- Provide Infrastructure as Code templates (Terraform, Pulumi)
```

### 5. Development Workflow Optimization
```markdown
## Development Workflow Guidelines

### TodoWrite Integration
- Use TodoWrite tool for task tracking and progress visibility
- Break complex tasks into manageable, trackable subtasks
- Mark todos as completed immediately upon finishing work
- Provide clear task descriptions with acceptance criteria

### AI Agent Collaboration
- Create comprehensive prompts for reusable development patterns
- Document learnings and patterns for future reference
- Establish clear handoff procedures between different AI agents
- Maintain context and continuity across development sessions

### Code Generation & Templates
- Implement incremental compilation for faster development cycles
- Use template caching with TTL for performance optimization
- Support template inheritance and composition
- Provide real-time validation and error feedback
```

### 6. Quality Assurance & Testing
```markdown
## Quality Assurance Guidelines

### Testing Strategy
- Implement comprehensive test suites with >90% coverage
- Use contract testing for API validation
- Include performance testing with baseline comparisons
- Implement chaos engineering for reliability validation

### Security & Compliance
- Integrate SAST/DAST tools in CI/CD pipeline
- Implement dependency vulnerability scanning
- Provide compliance frameworks (SOC2, GDPR, HIPAA)
- Include audit trails and security event logging

### Code Quality
- Enforce consistent code formatting and linting
- Use static analysis tools for code quality metrics
- Implement automated code review processes
- Maintain architectural decision records (ADRs)
```

### 7. Performance & Scalability
```markdown
## Performance Guidelines

### Optimization Strategies
- Implement parallel processing for CPU-intensive tasks
- Use worker pools with proper resource management
- Implement caching strategies with appropriate TTL
- Profile applications for bottleneck identification

### Scalability Patterns
- Design for horizontal scaling from the beginning
- Implement proper load balancing strategies
- Use event-driven architectures for loose coupling
- Support multi-tenancy with proper isolation
```

### 8. Documentation & Knowledge Management
```markdown
## Documentation Guidelines

### Comprehensive Documentation
- Maintain up-to-date API documentation with examples
- Create detailed installation and setup guides
- Provide troubleshooting guides for common issues
- Document architectural decisions and trade-offs

### Knowledge Preservation
- Create reusable prompts for common development patterns
- Document learnings and patterns for future reference
- Maintain comprehensive issue tracking for transparency
- Establish feedback loops for continuous improvement
```

## Implementation Priority

### High Priority (Immediate Implementation)
1. TodoWrite integration for task tracking
2. Architecture guidelines for DDD and Clean Architecture
3. Basic observability with OpenTelemetry
4. Single-line installation support

### Medium Priority (Next Phase)
1. Advanced SRE monitoring with SLO management
2. Feature composition and plugin architecture
3. Comprehensive testing and quality assurance
4. Multi-platform distribution strategy

### Low Priority (Future Enhancement)
1. AI-powered optimization and recommendations
2. Advanced performance profiling and optimization
3. Enterprise compliance and audit frameworks
4. Community marketplace and ecosystem development

## Expected Benefits

### Developer Experience
- Faster onboarding with clear guidelines and templates
- Reduced cognitive load through established patterns
- Consistent code quality across team members
- Automated workflows for common development tasks

### System Quality
- Higher reliability through comprehensive monitoring
- Better performance through optimization guidelines
- Enhanced security through security-first development
- Improved maintainability through clean architecture

### Business Impact
- Faster time-to-market through reusable patterns
- Reduced technical debt through quality guidelines
- Better scalability through proper architecture
- Enhanced user adoption through superior user experience

## Measurement and Feedback

### Success Metrics
- Development velocity (story points per sprint)
- Code quality metrics (test coverage, complexity)
- System reliability (SLO achievement, incident frequency)
- User satisfaction (adoption rates, retention metrics)

### Continuous Improvement
- Regular guideline reviews and updates
- Feedback collection from development teams
- Performance benchmarking and optimization
- Industry best practice integration

## Related Documentation
- Domain-Driven Design patterns and practices
- Clean Architecture implementation guides
- OpenTelemetry integration best practices
- Enterprise security and compliance frameworks