# template-health-endpoint Product Requirements Document (PRD)

## Goal, Objective and Context

**Primary Objective:** Create a comprehensive template repository for health endpoint APIs using TypeSpec for type definitions, with automatic JSON Schema & OpenAPI v3 generation, Go code generation, OpenTelemetry integration, and CloudEvents support.

**Context:** The modern microservices ecosystem lacks standardized, comprehensive health endpoint templates. Developers repeatedly implement health checks from scratch, leading to inconsistent implementations, missing observability features, and poor integration with cloud-native infrastructure. This project establishes the industry standard for health endpoint implementation through TypeSpec-driven templates.

**Success Criteria:**
- [ ] Generate valid TypeSpec definitions for all template tiers
- [ ] Produce correct JSON Schema and OpenAPI v3 specifications
- [ ] Generate working Go server code with Server Timing API integration
- [ ] Create functional TypeScript client SDKs
- [ ] Provide Kubernetes-ready health probe configurations
- [ ] Implement native OpenTelemetry instrumentation
- [ ] Enable CloudEvents integration for health status changes
- [ ] Pass all validation tests across template tiers

## Functional Requirements (MVP)

**Core Template Generation System:**
1. **TypeSpec-First Development:** All health API definitions must be written in TypeSpec with automatic generation of JSON Schema, OpenAPI v3, Go code, and TypeScript types
2. **Four-Tier Template Hierarchy:** Basic (5 min) → Intermediate (15 min) → Advanced (30 min) → Enterprise (45 min) deployment times
3. **Go-Based Generator:** Command-line tool for template generation and code scaffolding
4. **Reference Implementation Integration:** Build upon LarsArtmann/CV health.go implementation patterns

**Health Endpoint Requirements:**
1. **Basic Health Check:** `/health` endpoint returning status, timestamp, version, uptime
2. **ServerTime API:** `/health/time` endpoint with RFC3339, Unix timestamps, timezone info, formatted time
3. **Kubernetes Probes:** `/health/ready` (readiness) and `/health/live` (liveness) endpoints
4. **Dependency Checks:** External service health validation for intermediate+ tiers

**Observability Requirements:**
1. **OpenTelemetry Integration:** Native tracing, metrics, and logging for all endpoints
2. **Server Timing API:** Performance metrics in HTTP headers for browser DevTools
3. **CloudEvents Support:** Health status change events following CloudEvents specification
4. **Prometheus Metrics:** Metrics exposition for monitoring and alerting

**Code Generation Requirements:**
1. **Go Server Code:** HTTP handlers, middleware, models, and client SDKs
2. **TypeScript Types:** Complete type definitions and client SDK generation
3. **Kubernetes Manifests:** Health probe configs, ServiceMonitor, Ingress templates
4. **Documentation:** Auto-generated API documentation and integration guides

## Non Functional Requirements (MVP)

**Performance:**
- Health endpoints must respond within 100ms under normal conditions
- ServerTime API must include sub-millisecond timing accuracy
- Template generation must complete within 30 seconds for enterprise tier

**Reliability:**
- Health checks must have 99.9% uptime when dependencies are available
- Graceful degradation when external dependencies are unavailable
- Circuit breaker patterns for dependency health checks

**Scalability:**
- Templates must support horizontal scaling patterns
- Health endpoints must handle 1000+ requests per second
- OpenTelemetry integration must not impact performance by more than 5%

**Security:**
- No sensitive information exposed in health endpoints
- Optional authentication support for enterprise tier
- Rate limiting capabilities for health endpoints

**Maintainability:**
- Clear separation between template tiers
- Comprehensive test coverage for generated code
- Automated validation for TypeSpec definitions

## User Interaction and Design Goals

This project is primarily a developer-focused template system with no traditional UI. However, the "user experience" focuses on developer interaction:

**Developer Experience Goals:**
- **Simplicity:** Single command template generation with sensible defaults
- **Progressive Complexity:** Clear upgrade path from basic to enterprise tiers
- **Documentation-First:** Comprehensive guides and examples for each tier
- **IDE Integration:** Generated code should work seamlessly with Go and TypeScript tooling

**Key Interaction Paradigms:**
- **CLI-First:** Primary interaction through command-line generator tool
- **Configuration-Driven:** YAML/JSON configuration for template customization
- **Template Selection:** Clear tier selection with value proposition explanation

## Technical Assumptions

**Repository & Service Architecture:** 
Standalone template repository (polyrepo) with generated microservice templates. Each generated service follows cloud-native microservice patterns with clear separation of concerns. The template repository itself uses a monorepo structure for easier maintenance of related templates, but generates polyrepo-style individual services.

**Rationale:** Template repositories benefit from monorepo structure for shared tooling and consistency, while generated services should be independently deployable microservices. This approach provides the best of both worlds - centralized template management with distributed service deployment.

### Testing Requirements

**Template Validation:**
- TypeSpec schema validation and compilation tests
- Generated code compilation and unit tests
- Integration tests for health endpoints
- Contract testing between generated clients and servers

**Kubernetes Integration Testing:**
- Health probe functionality validation
- Service discovery integration tests
- Load balancer health check validation
- OpenTelemetry Collector integration tests

**Performance Testing:**
- Health endpoint response time validation
- Server Timing API accuracy tests
- OpenTelemetry overhead measurement
- Load testing for generated services

**End-to-End Testing:**
- Complete template generation workflow
- Multi-tier template upgrade scenarios
- CloudEvents emission and consumption
- Cross-platform compatibility (Linux, macOS, Windows)

## Epic Overview

- **Epic 1: Foundation & TypeSpec Schema Design**
  - Goal: Establish the core TypeSpec schemas and basic template generation infrastructure for health endpoints.
  - Story 1: As a developer, I want to analyze the reference health.go implementation so that I can extract reusable patterns for TypeSpec definitions.
    - AC1: Complete analysis document of LarsArtmann/CV health.go implementation
    - AC2: Identify key patterns for ServerTime API, error handling, and response structures
    - AC3: Document integration points for OpenTelemetry and Server Timing API
  - Story 2: As a developer, I want comprehensive TypeSpec schemas for health endpoints so that I can generate consistent APIs across all tiers.
    - AC1: TypeSpec models for HealthStatus, ServerTime, and ServerTimingMetrics
    - AC2: TypeSpec interfaces for health, time, ready, and live endpoints
    - AC3: CloudEvents schema integration for health status changes
    - AC4: JSON Schema and OpenAPI v3 generation from TypeSpec definitions
  - Story 3: As a developer, I want a Go-based template generator so that I can create health endpoint projects from TypeSpec definitions.
    - AC1: CLI tool for template generation with tier selection
    - AC2: TypeSpec compilation and validation
    - AC3: Project scaffolding with proper directory structure
    - AC4: Configuration file generation for template customization

- **Epic 2: Go Code Generation & Basic Template**
  - Goal: Implement Go code generation from TypeSpec schemas and create the basic template tier with essential health endpoints.
  - Story 1: As a developer, I want Go code generation from TypeSpec so that I can have type-safe server implementations.
    - AC1: Go struct generation from TypeSpec models
    - AC2: HTTP handler generation for health endpoints
    - AC3: Middleware generation for OpenTelemetry integration
    - AC4: Client SDK generation for service-to-service communication
  - Story 2: As a developer, I want a basic health endpoint template so that I can quickly add health checks to my services.
    - AC1: `/health` endpoint with status, timestamp, version, uptime
    - AC2: `/health/time` endpoint with comprehensive timestamp formats
    - AC3: Basic error handling and graceful degradation
    - AC4: Docker and deployment configuration
  - Story 3: As a developer, I want TypeScript type generation so that I can integrate health endpoints with frontend applications.
    - AC1: TypeScript interface generation from TypeSpec models
    - AC2: Client SDK generation for health endpoint consumption
    - AC3: Type-safe API client with proper error handling
    - AC4: npm package structure for easy distribution

- **Epic 3: Observability Integration & Advanced Templates**
  - Goal: Implement comprehensive observability features including OpenTelemetry, Server Timing API, and CloudEvents support for intermediate and advanced template tiers.
  - Story 1: As a developer, I want native OpenTelemetry integration so that I can monitor health endpoint performance and trace requests.
    - AC1: Automatic trace generation for all health endpoints
    - AC2: Custom metrics for health status and response times
    - AC3: Structured logging with trace and span ID correlation
    - AC4: Baggage propagation for distributed tracing
  - Story 2: As a developer, I want Server Timing API integration so that I can analyze performance in browser DevTools.
    - AC1: Server-Timing headers for database queries, cache hits, external API calls
    - AC2: Integration with OpenTelemetry spans for timing correlation
    - AC3: Configurable timing metrics based on environment
    - AC4: TypeScript client support for timing consumption
  - Story 3: As a developer, I want CloudEvents support so that I can implement event-driven health monitoring.
    - AC1: CloudEvents emission for health status changes
    - AC2: Event schemas for dependency failures and recovery
    - AC3: Integration with event brokers (NATS, Kafka, etc.)
    - AC4: Event-driven alerting and monitoring capabilities
  - Story 4: As a developer, I want dependency health checks so that I can monitor external service availability.
    - AC1: Configurable dependency health check framework
    - AC2: Circuit breaker patterns for failing dependencies
    - AC3: Timeout and retry configuration
    - AC4: Dependency status aggregation in health responses

- **Epic 4: Kubernetes Integration & Enterprise Template**
  - Goal: Provide complete Kubernetes integration with enterprise-grade features including compliance, security, and advanced monitoring capabilities.
  - Story 1: As a DevOps engineer, I want Kubernetes health probe configurations so that I can deploy services with proper health monitoring.
    - AC1: Liveness, readiness, and startup probe configurations
    - AC2: Service discovery and load balancer integration
    - AC3: Ingress routing with health check endpoints
    - AC4: ConfigMap templates for health check configuration
  - Story 2: As a platform engineer, I want Prometheus integration so that I can monitor health metrics across my infrastructure.
    - AC1: Prometheus metrics exposition endpoint
    - AC2: ServiceMonitor configuration for Prometheus Operator
    - AC3: Grafana dashboard templates for health monitoring
    - AC4: Alerting rules for health status changes
  - Story 3: As an enterprise architect, I want compliance-ready health endpoints so that I can meet regulatory requirements.
    - AC1: Audit logging for health endpoint access
    - AC2: Authentication and authorization support
    - AC3: Rate limiting and DDoS protection
    - AC4: Compliance documentation (SOC2, HIPAA, PCI)
  - Story 4: As a developer, I want comprehensive documentation and examples so that I can effectively use the health endpoint templates.
    - AC1: Complete setup and usage documentation
    - AC2: Integration guides for major platforms (Kubernetes, Docker, etc.)
    - AC3: Best practices documentation for each template tier
    - AC4: Example implementations and migration guides

## Key Reference Documents

- `docs/project-brief.md` - Project vision and requirements
- `docs/architecture.md` - Technical architecture and design decisions
- `docs/epic-1.md` - Foundation & TypeSpec Schema Design
- `docs/epic-2.md` - Go Code Generation & Basic Template
- `docs/epic-3.md` - Observability Integration & Advanced Templates
- `docs/epic-4.md` - Kubernetes Integration & Enterprise Template
- `docs/typespec-schemas/` - TypeSpec schema definitions
- `docs/examples/` - Example implementations and usage guides

## Out of Scope Ideas Post MVP

**Extended Language Support:**
- Rust, Python, Java, C# code generation (evaluate post-MVP based on community demand)
- Additional framework integrations (Spring Boot, FastAPI, etc.)

**Advanced Monitoring Features:**
- Custom health check plugin system with marketplace
- Machine learning-based anomaly detection
- Historical health data analysis and trending
- Multi-region health status aggregation

**Enterprise Extensions:**
- Service mesh integration (Istio, Linkerd, Consul Connect)
- Advanced security features (mTLS, JWT validation, RBAC)
- Multi-environment configuration management
- Enterprise support and SLA guarantees

**Developer Experience Enhancements:**
- VS Code extension for template generation
- Web-based template configurator
- Integration with popular development frameworks
- Automated migration tools between template tiers

## [OPTIONAL: For Simplified PM-to-Development Workflow Only] Core Technical Decisions & Application Structure

### Technology Stack Selections

**Core Technologies:**
- **TypeSpec**: Primary API definition language for schema-first development
- **Go 1.21+**: Server implementation and template generator
- **TypeScript 5.0+**: Client SDK and type definitions
- **OpenTelemetry Go SDK**: Native observability and tracing
- **CloudEvents Go SDK**: Event-driven architecture support

**Development Tools:**
- **TypeSpec Compiler**: Schema validation and code generation
- **Go Modules**: Dependency management
- **npm/yarn**: TypeScript package management
- **Docker**: Containerization and deployment
- **Kubernetes**: Orchestration and health probe integration

**Testing Framework:**
- **Go testing**: Unit and integration tests
- **Jest**: TypeScript testing
- **Testcontainers**: Integration testing with real dependencies
- **K3s**: Kubernetes integration testing

### Application Structure

```
template-health-endpoint/
├── cmd/
│   └── generator/           # Go-based template generator
│       ├── main.go
│       ├── typespec.go      # TypeSpec processing
│       ├── codegen.go       # Code generation logic
│       └── templates.go     # Template management
├── pkg/
│   ├── schemas/             # TypeSpec schema definitions
│   │   ├── health.tsp
│   │   ├── server-time.tsp
│   │   └── cloudevents.tsp
│   ├── generator/           # Code generation packages
│   │   ├── golang/          # Go code generation
│   │   ├── typescript/      # TypeScript generation
│   │   └── kubernetes/      # K8s manifest generation
│   └── templates/           # Template definitions
│       ├── basic/
│       ├── intermediate/
│       ├── advanced/
│       └── enterprise/
├── examples/                # Generated example projects
│   ├── basic-go/
│   ├── intermediate-ts/
│   └── enterprise-k8s/
├── docs/                    # Documentation and guides
├── scripts/                 # Build and utility scripts
├── tests/                   # Integration and E2E tests
├── Dockerfile
├── go.mod
└── README.md
```

**Key Modules/Components and Responsibilities:**
- **cmd/generator**: CLI tool for template generation and project scaffolding
- **pkg/schemas**: TypeSpec schema definitions for all health endpoint models
- **pkg/generator**: Code generation engines for Go, TypeScript, and Kubernetes manifests
- **pkg/templates**: Template definitions and tier-specific configurations
- **examples**: Generated example projects demonstrating each template tier
- **tests**: Comprehensive testing suite for validation and integration

**Data Flow Overview (Conceptual):**
TypeSpec Schemas → Template Generator → Generated Code (Go/TS) → Kubernetes Manifests → Deployed Health Endpoints → Observability Data (OpenTelemetry/CloudEvents) → Monitoring Systems

## Change Log

| Change | Date | Version | Description | Author |
| ------ | ---- | ------- | ----------- | ------ |
| Initial PRD Creation | 2025-01-XX | 1.0.0 | Complete PRD with 4 epics and comprehensive requirements | BMAD PM Agent |

----- END PRD START CHECKLIST OUTPUT ------

## Checklist Results Report

*To be populated by PO Agent during validation phase*

----- END Checklist START Design Architect `UI/UX Specification Mode` Prompt ------

*This project is primarily a developer-focused template system with no traditional UI requirements. Design Architect phase can be skipped.*

----- END Design Architect `UI/UX Specification Mode` Prompt START Architect Prompt ------

## Initial Architect Prompt

Based on our discussions and requirements analysis for the template-health-endpoint project, I've compiled the following technical guidance to inform your architecture analysis and decisions to kick off Architecture Creation Mode:

### Technical Infrastructure

**Primary Technology Stack:**
- TypeSpec for API-first schema definitions
- Go 1.21+ for server implementation and template generator
- TypeScript 5.0+ for client SDKs and type definitions
- OpenTelemetry Go SDK for native observability
- CloudEvents Go SDK for event-driven architecture

**Key Architectural Decisions Needed:**
1. **Template Generation Architecture**: Design the Go-based generator system for processing TypeSpec schemas and generating multi-tier templates
2. **Code Generation Pipeline**: Define the workflow from TypeSpec → JSON Schema/OpenAPI → Go/TypeScript code
3. **Observability Integration**: Architect the OpenTelemetry and Server Timing API integration patterns
4. **Kubernetes Integration**: Design the health probe and ServiceMonitor template generation
5. **CloudEvents Architecture**: Define event schemas and emission patterns for health status changes

**Critical Integration Points:**
- Reference implementation patterns from LarsArtmann/CV health.go
- Server Timing API integration for browser DevTools compatibility
- Native OpenTelemetry instrumentation without performance degradation
- CloudEvents specification compliance for event-driven monitoring
- Kubernetes health probe standards and best practices

**Performance Requirements:**
- Health endpoints must respond within 100ms
- Template generation within 30 seconds for enterprise tier
- OpenTelemetry overhead limited to 5% performance impact
- Support for 1000+ requests per second on health endpoints

Please proceed with creating the comprehensive architecture document that addresses these requirements and provides detailed technical specifications for the development team.
