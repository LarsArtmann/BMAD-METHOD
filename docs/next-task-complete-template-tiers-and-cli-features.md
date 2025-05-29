# Complete Template Tiers and Advanced CLI Features

## 🎯 **Task Overview**

**Objective**: Complete the template-health-endpoint project by implementing all remaining template tiers (Intermediate, Advanced, Enterprise) and adding advanced CLI functionality (update, migrate, customize commands).

**Priority**: **CRITICAL** - Required to fulfill GitHub issue #127 requirements  
**Estimated Effort**: 4-6 hours  
**Impact**: **HIGH** - Completes the core project deliverables

## 📋 **Current State Analysis**

### ✅ **Completed (60% Complete)**
- ✅ **Basic Template Tier**: Fully implemented with static templates and CLI generation
- ✅ **Dual Template System**: Both static templates and CLI generation working
- ✅ **Core CLI Commands**: `generate`, `template list`, `template from-static`, `template validate`
- ✅ **Template Processing**: Comprehensive variable substitution across all file types
- ✅ **Testing Framework**: Unit tests, integration tests, validation scripts
- ✅ **Documentation**: Setup guides, tier comparison, migration guides
- ✅ **TypeSpec Schemas**: Complete health endpoint schemas with CloudEvents support

### 🔄 **Remaining Work (40% to Complete)**
- 🔄 **Intermediate Template Tier**: Add dependency health checks, basic OpenTelemetry
- 🔄 **Advanced Template Tier**: Full observability, CloudEvents, Server Timing API
- 🔄 **Enterprise Template Tier**: Kubernetes ServiceMonitor, compliance features
- 🔄 **CLI Update Command**: Update existing projects to newer template versions
- 🔄 **CLI Migrate Command**: Migrate projects between tiers
- 🔄 **CLI Customize Command**: Interactive template customization
- 🔄 **Template Validation**: Ensure all tiers work correctly

## 🎯 **Detailed Task Breakdown**

### **Phase 1: Complete Template Tiers (2-3 hours)**

#### **Task 1.1: Intermediate Template Tier**
**Objective**: Create production-ready template with dependency health checks

**Requirements**:
- Add dependency health check endpoints (`/health/dependencies`)
- Include basic OpenTelemetry instrumentation
- Add structured logging configuration
- Include database and cache health check examples
- Update Kubernetes manifests with additional ConfigMaps

**Implementation Steps**:
1. Generate intermediate tier from basic tier
2. Add dependency health check handlers
3. Include OpenTelemetry SDK dependencies
4. Add configuration for external dependencies
5. Update documentation and examples

**Files to Create/Modify**:
```
templates/intermediate/
├── internal/handlers/dependencies.go    # NEW: Dependency health checks
├── internal/config/dependencies.go     # NEW: Dependency configuration
├── internal/observability/otel.go      # NEW: OpenTelemetry setup
├── go.mod.tmpl                         # ADD: OpenTelemetry dependencies
├── deployments/kubernetes/
│   ├── configmap-dependencies.yaml    # NEW: Dependency configuration
│   └── deployment.yaml                # UPDATE: Environment variables
└── docs/dependencies.md               # NEW: Dependency documentation
```

#### **Task 1.2: Advanced Template Tier**
**Objective**: Create full observability template with CloudEvents

**Requirements**:
- Full OpenTelemetry instrumentation (traces, metrics, logs)
- Server Timing API implementation
- CloudEvents emission for health status changes
- Custom metrics and performance tracking
- Advanced Kubernetes configurations

**Implementation Steps**:
1. Generate advanced tier from intermediate tier
2. Add comprehensive OpenTelemetry instrumentation
3. Implement Server Timing API middleware
4. Add CloudEvents publisher for health events
5. Include Prometheus metrics exposition
6. Add Jaeger tracing configuration

**Files to Create/Modify**:
```
templates/advanced/
├── internal/observability/
│   ├── metrics.go                     # NEW: Custom metrics
│   ├── tracing.go                     # NEW: Distributed tracing
│   ├── server_timing.go               # NEW: Server Timing API
│   └── cloudevents.go                 # NEW: CloudEvents publisher
├── internal/middleware/
│   ├── observability.go               # NEW: Observability middleware
│   └── server_timing.go               # NEW: Server Timing middleware
├── deployments/kubernetes/
│   ├── servicemonitor.yaml            # NEW: Prometheus ServiceMonitor
│   └── otel-collector.yaml            # NEW: OpenTelemetry Collector
└── docs/observability.md              # NEW: Observability guide
```

#### **Task 1.3: Enterprise Template Tier**
**Objective**: Create enterprise-grade template with compliance features

**Requirements**:
- Advanced security features (mTLS, RBAC)
- Compliance logging and audit trails
- Multi-environment configuration
- Advanced monitoring and alerting
- Service mesh integration

**Implementation Steps**:
1. Generate enterprise tier from advanced tier
2. Add mTLS certificate management
3. Implement RBAC integration
4. Add compliance audit logging
5. Include multi-environment configuration
6. Add advanced security middleware

**Files to Create/Modify**:
```
templates/enterprise/
├── internal/security/
│   ├── mtls.go                        # NEW: mTLS implementation
│   ├── rbac.go                        # NEW: RBAC integration
│   └── audit.go                       # NEW: Audit logging
├── internal/compliance/
│   ├── logging.go                     # NEW: Compliance logging
│   └── reporting.go                   # NEW: Compliance reporting
├── configs/
│   ├── development.yaml               # NEW: Dev environment config
│   ├── staging.yaml                   # NEW: Staging environment config
│   └── production.yaml                # NEW: Production environment config
├── deployments/kubernetes/
│   ├── rbac.yaml                      # NEW: RBAC configuration
│   ├── certificates.yaml              # NEW: Certificate management
│   └── ingress-mtls.yaml              # NEW: mTLS Ingress
└── docs/
    ├── security.md                    # NEW: Security guide
    └── compliance.md                  # NEW: Compliance guide
```

### **Phase 2: Advanced CLI Features (2-3 hours)**

#### **Task 2.1: CLI Update Command**
**Objective**: Enable updating existing projects to newer template versions

**Requirements**:
- Compare current project with newer template version
- Selective updates (e.g., only Kubernetes configs)
- Backup existing files before updates
- Merge conflicts resolution guidance

**Implementation**:
```go
// cmd/generator/commands/update.go
func runUpdateProject(cmd *cobra.Command, args []string) error {
    // 1. Detect current project template version
    // 2. Compare with target template version
    // 3. Show diff of changes
    // 4. Apply updates with user confirmation
    // 5. Provide rollback instructions
}
```

**CLI Usage**:
```bash
template-health-endpoint update --project ./my-service --template-version v1.2.0
template-health-endpoint update --project ./my-service --selective kubernetes,docs
template-health-endpoint update --project ./my-service --dry-run
```

#### **Task 2.2: CLI Migrate Command**
**Objective**: Enable migrating projects between template tiers

**Requirements**:
- Migrate from basic → intermediate → advanced → enterprise
- Add new dependencies and configurations
- Update existing code with new features
- Provide migration validation

**Implementation**:
```go
// cmd/generator/commands/migrate.go
func runMigrateProject(cmd *cobra.Command, args []string) error {
    // 1. Detect current project tier
    // 2. Validate migration path
    // 3. Add new dependencies and files
    // 4. Update existing configurations
    // 5. Validate migrated project
}
```

**CLI Usage**:
```bash
template-health-endpoint migrate --project ./my-service --from basic --to intermediate
template-health-endpoint migrate --project ./my-service --to advanced --dry-run
```

#### **Task 2.3: CLI Customize Command**
**Objective**: Enable interactive template customization

**Requirements**:
- Interactive prompts for common customizations
- Template variable file support
- Preview changes before applying
- Save customization profiles

**Implementation**:
```go
// cmd/generator/commands/customize.go
func runCustomizeTemplate(cmd *cobra.Command, args []string) error {
    // 1. Load template tier
    // 2. Present customization options
    // 3. Collect user preferences
    // 4. Generate customized template
    // 5. Save customization profile
}
```

**CLI Usage**:
```bash
template-health-endpoint customize --tier basic --interactive
template-health-endpoint customize --tier advanced --config custom.yaml
template-health-endpoint customize --tier enterprise --save-profile my-org
```

### **Phase 3: Validation and Testing (1 hour)**

#### **Task 3.1: Comprehensive Template Testing**
**Objective**: Ensure all template tiers work correctly

**Requirements**:
- Generate projects from all template tiers
- Compile and run all generated projects
- Test all health endpoints for each tier
- Validate Kubernetes deployments
- Test CLI commands end-to-end

**Implementation**:
```bash
# Update validation script to test all tiers
./scripts/validate-templates.sh --all-tiers
./scripts/test-cli-commands.sh
./scripts/test-generated-projects.sh
```

#### **Task 3.2: Documentation Updates**
**Objective**: Update documentation for all new features

**Requirements**:
- Update tier comparison with all 4 tiers
- Add CLI command documentation
- Update migration guides
- Add troubleshooting sections

## 🔧 **Technical Implementation Details**

### **Template Tier Differentiation**

**Basic Tier**:
```yaml
features:
  kubernetes: true
  typescript: true
  docker: true
  opentelemetry: false
  cloudevents: false
  dependencies: false
```

**Intermediate Tier**:
```yaml
features:
  kubernetes: true
  typescript: true
  docker: true
  opentelemetry: basic
  cloudevents: false
  dependencies: true
```

**Advanced Tier**:
```yaml
features:
  kubernetes: true
  typescript: true
  docker: true
  opentelemetry: full
  cloudevents: true
  dependencies: true
  server_timing: true
  metrics: custom
```

**Enterprise Tier**:
```yaml
features:
  kubernetes: true
  typescript: true
  docker: true
  opentelemetry: full
  cloudevents: true
  dependencies: true
  server_timing: true
  metrics: custom
  security: mtls
  compliance: true
  multi_env: true
```

### **CLI Command Architecture**

```go
// Enhanced CLI structure
rootCmd
├── generate              # Generate new project (embedded templates)
├── template
│   ├── list             # List available static templates
│   ├── from-static      # Generate from static template
│   └── validate         # Validate template integrity
├── update               # Update existing project
├── migrate              # Migrate between tiers
├── customize            # Interactive template customization
└── validate             # Validate TypeSpec schemas
```

### **Template Processing Pipeline**

```go
type TemplateProcessor struct {
    SourceTemplate string
    TargetTier     string
    Variables      map[string]interface{}
    Features       FeatureConfig
}

func (tp *TemplateProcessor) Process() error {
    // 1. Load template metadata
    // 2. Process template variables
    // 3. Apply tier-specific features
    // 4. Generate output files
    // 5. Validate generated project
}
```

## 📊 **Success Criteria**

### **Functional Requirements**
- [ ] All 4 template tiers generate working projects
- [ ] All generated projects compile and run successfully
- [ ] All health endpoints respond correctly for each tier
- [ ] CLI update command works for version upgrades
- [ ] CLI migrate command works for tier transitions
- [ ] CLI customize command enables template customization
- [ ] Kubernetes deployments work for all tiers
- [ ] OpenTelemetry integration works in advanced/enterprise tiers
- [ ] CloudEvents emission works in advanced/enterprise tiers

### **Quality Requirements**
- [ ] Comprehensive test coverage for all tiers
- [ ] Documentation is complete and accurate
- [ ] Error handling is robust and user-friendly
- [ ] Performance is acceptable for all operations
- [ ] Security best practices are implemented

### **Integration Requirements**
- [ ] Works with existing template-* ecosystem
- [ ] Integrates with CI/CD pipelines
- [ ] Compatible with Kubernetes environments
- [ ] Supports OpenTelemetry observability stack
- [ ] Works with CloudEvents infrastructure

## 🚀 **Getting Started**

### **Prerequisites**
- Current BMAD-METHOD repository with basic tier completed
- Go 1.21+ development environment
- Docker (for testing containerization)
- kubectl (for testing Kubernetes deployments)

### **Implementation Order**
1. **Start with Intermediate Tier**: Build upon basic tier foundation
2. **Add Advanced Tier**: Implement full observability features
3. **Create Enterprise Tier**: Add security and compliance features
4. **Implement CLI Commands**: Add update, migrate, customize functionality
5. **Comprehensive Testing**: Validate all tiers and CLI commands
6. **Documentation Updates**: Ensure all documentation is current

### **Validation Steps**
1. Generate projects from all template tiers
2. Compile and run all generated projects
3. Test all health endpoints for each tier
4. Test CLI commands with real projects
5. Validate Kubernetes deployments
6. Run comprehensive test suite

## 📚 **Context and Background**

### **Original GitHub Issue #127 Requirements**
- Create comprehensive template repository for health endpoint APIs
- Support 4 template tiers: Basic → Intermediate → Advanced → Enterprise
- TypeSpec-first API definitions with code generation
- Go server implementations with OpenTelemetry integration
- TypeScript client SDKs
- Kubernetes deployment configurations
- Progressive complexity with clear upgrade paths

### **Current Achievement**
- ✅ 60% complete with solid foundation
- ✅ Basic tier fully functional
- ✅ Dual template system (static + CLI) working
- ✅ Comprehensive testing and documentation framework
- ✅ TypeSpec schemas and code generation pipeline

### **Remaining Work**
- 🔄 40% remaining to complete all requirements
- 🔄 3 additional template tiers
- 🔄 Advanced CLI functionality
- 🔄 Final validation and testing

## 🎯 **Expected Outcomes**

Upon completion of this task:

1. **Complete Template Repository**: All 4 tiers available as static templates
2. **Full CLI Functionality**: Generate, update, migrate, customize commands
3. **Production-Ready**: Enterprise-grade templates with security and compliance
4. **Comprehensive Documentation**: Complete guides for all tiers and features
5. **Validated Quality**: All templates tested and working correctly
6. **GitHub Issue #127 Fulfilled**: All original requirements completed

This task represents the final 40% of work needed to complete the template-health-endpoint project and fulfill all requirements from the original GitHub issue #127. The foundation is solid, and this task will deliver a production-ready, comprehensive template system that serves both manual and automated workflows.

---

**Ready to implement? Start with Phase 1, Task 1.1: Intermediate Template Tier!** 🚀
