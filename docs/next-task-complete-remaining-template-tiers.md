# Complete Remaining Template Tiers and Advanced CLI Features

## 🎯 **Task Overview**

**Objective**: Complete the template-health-endpoint project by implementing the remaining template tiers (Intermediate, Advanced, Enterprise) and adding advanced CLI functionality (update, migrate, customize commands) to fulfill all GitHub issue #127 requirements.

**Priority**: **CRITICAL** - Required to complete the project and fulfill original requirements  
**Estimated Effort**: 4-6 hours  
**Impact**: **MAXIMUM** - Takes project from 60% to 100% complete  
**Context**: This is the final task to complete the template-health-endpoint project

## 📋 **Current State Analysis**

### ✅ **Completed (60% Complete)**
Based on the current chat history and repository state:

- ✅ **Basic Template Tier**: Fully implemented with static templates and CLI generation
- ✅ **Dual Template System**: Both static templates (`/templates/basic/`) and CLI generation working
- ✅ **Core CLI Commands**: `generate`, `template list`, `template from-static`, `template validate`
- ✅ **Template Processing**: Comprehensive variable substitution across all file types (.go, .yaml, .json, .ts, .sh, etc.)
- ✅ **Testing Framework**: Unit tests, integration tests, validation scripts with 100% success rate
- ✅ **Documentation**: Setup guides, tier comparison, migration guides, comprehensive PROJECT_GUIDELINES.md
- ✅ **TypeSpec Schemas**: Complete health endpoint schemas with CloudEvents support
- ✅ **Repository Structure**: Proper template-* repository pattern with static templates

### 🔄 **Remaining Work (40% to Complete)**
- 🔄 **Intermediate Template Tier**: Add dependency health checks, basic OpenTelemetry
- 🔄 **Advanced Template Tier**: Full observability, CloudEvents, Server Timing API
- 🔄 **Enterprise Template Tier**: Kubernetes ServiceMonitor, compliance features, mTLS
- 🔄 **CLI Update Command**: Update existing projects to newer template versions
- 🔄 **CLI Migrate Command**: Migrate projects between tiers (basic → intermediate → advanced → enterprise)
- 🔄 **CLI Customize Command**: Interactive template customization
- 🔄 **Final Validation**: Ensure all tiers work correctly and meet GitHub issue #127 requirements

## 🎯 **Detailed Implementation Plan**

### **Phase 1: Complete Template Tiers (2-3 hours)**

#### **Task 1.1: Intermediate Template Tier**
**Objective**: Create production-ready template with dependency health checks

**Current State**: Only basic tier exists in `/templates/basic/`  
**Target**: Create `/templates/intermediate/` with enhanced features

**Implementation Steps**:
1. **Generate Base**: Use existing CLI to generate intermediate tier project
   ```bash
   ./bin/template-health-endpoint generate --name template-intermediate --tier intermediate --module github.com/template/intermediate --output templates/intermediate
   ```

2. **Add Dependency Health Checks**:
   ```go
   // templates/intermediate/internal/handlers/dependencies.go
   func (h *HealthHandler) DependenciesCheck(w http.ResponseWriter, r *http.Request) {
       // Check database, cache, external APIs
       dependencies := map[string]DependencyStatus{
           "database": h.checkDatabase(),
           "cache":    h.checkCache(),
           "external_api": h.checkExternalAPI(),
       }
       // Return dependency status
   }
   ```

3. **Add Basic OpenTelemetry**:
   ```go
   // templates/intermediate/internal/observability/otel.go
   func InitOpenTelemetry(serviceName string) {
       // Basic OTLP exporter setup
       // Trace provider configuration
       // Metric provider configuration
   }
   ```

4. **Update Template Metadata**:
   ```yaml
   # templates/intermediate/template.yaml
   name: intermediate
   description: Intermediate tier health endpoint template
   tier: intermediate
   features:
     kubernetes: true
     typescript: true
     docker: true
     opentelemetry: basic
     dependencies: true
     cloudevents: false
   ```

5. **Convert to Template Variables**: Run conversion script to replace hardcoded values

#### **Task 1.2: Advanced Template Tier**
**Objective**: Create full observability template with CloudEvents

**Implementation Steps**:
1. **Generate Base**: Start from intermediate tier
2. **Add Full OpenTelemetry**:
   ```go
   // templates/advanced/internal/observability/
   ├── metrics.go          # Custom metrics
   ├── tracing.go          # Distributed tracing
   ├── server_timing.go    # Server Timing API
   └── cloudevents.go      # CloudEvents publisher
   ```

3. **Add Server Timing API**:
   ```go
   // templates/advanced/internal/middleware/server_timing.go
   func ServerTimingMiddleware(next http.Handler) http.Handler {
       return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
           start := time.Now()
           next.ServeHTTP(w, r)
           duration := time.Since(start)
           w.Header().Set("Server-Timing", fmt.Sprintf("total;dur=%.1f", float64(duration.Nanoseconds())/1e6))
       })
   }
   ```

4. **Add CloudEvents Integration**:
   ```go
   // templates/advanced/internal/events/health_events.go
   func EmitHealthStatusChange(status string) {
       event := cloudevents.NewEvent()
       event.SetType("health.status.changed")
       event.SetSource("health-endpoint")
       event.SetData(cloudevents.ApplicationJSON, map[string]string{"status": status})
       // Send event
   }
   ```

5. **Update Kubernetes Manifests**:
   ```yaml
   # templates/advanced/deployments/kubernetes/servicemonitor.yaml
   apiVersion: monitoring.coreos.com/v1
   kind: ServiceMonitor
   metadata:
     name: {{.Config.Name}}
   spec:
     selector:
       matchLabels:
         app: {{.Config.Name}}
     endpoints:
     - port: http
       path: /health/metrics
   ```

#### **Task 1.3: Enterprise Template Tier**
**Objective**: Create enterprise-grade template with compliance features

**Implementation Steps**:
1. **Generate Base**: Start from advanced tier
2. **Add mTLS Support**:
   ```go
   // templates/enterprise/internal/security/mtls.go
   func SetupMTLS(certFile, keyFile, caFile string) *tls.Config {
       // mTLS certificate configuration
       // Client certificate validation
       // Certificate rotation support
   }
   ```

3. **Add RBAC Integration**:
   ```go
   // templates/enterprise/internal/security/rbac.go
   func ValidateRBAC(r *http.Request) error {
       // Extract user identity
       // Check permissions
       // Audit access
   }
   ```

4. **Add Compliance Logging**:
   ```go
   // templates/enterprise/internal/compliance/audit.go
   func LogAuditEvent(event AuditEvent) {
       // Structured audit logging
       // Compliance data retention
       // Audit trail integrity
   }
   ```

5. **Add Multi-Environment Configuration**:
   ```yaml
   # templates/enterprise/configs/
   ├── development.yaml    # Dev environment config
   ├── staging.yaml        # Staging environment config
   └── production.yaml     # Production environment config
   ```

### **Phase 2: Advanced CLI Features (2-3 hours)**

#### **Task 2.1: CLI Update Command**
**Objective**: Enable updating existing projects to newer template versions

**Implementation**:
```go
// cmd/generator/commands/update.go
var updateCmd = &cobra.Command{
    Use:   "update",
    Short: "Update existing project to newer template version",
    Long:  "Update an existing project to a newer template version with selective updates.",
    RunE:  runUpdateProject,
}

func runUpdateProject(cmd *cobra.Command, args []string) error {
    projectPath, _ := cmd.Flags().GetString("project")
    templateVersion, _ := cmd.Flags().GetString("template-version")
    selective, _ := cmd.Flags().GetStringSlice("selective")
    dryRun, _ := cmd.Flags().GetBool("dry-run")
    
    // 1. Detect current project template version
    currentVersion, err := detectProjectVersion(projectPath)
    if err != nil {
        return fmt.Errorf("failed to detect project version: %w", err)
    }
    
    // 2. Compare with target template version
    changes, err := compareTemplateVersions(currentVersion, templateVersion)
    if err != nil {
        return fmt.Errorf("failed to compare versions: %w", err)
    }
    
    // 3. Show diff of changes
    fmt.Printf("🔍 Changes from %s to %s:\n", currentVersion, templateVersion)
    for _, change := range changes {
        fmt.Printf("  %s %s\n", change.Type, change.File)
    }
    
    if dryRun {
        fmt.Println("🔍 Dry run mode - no changes applied")
        return nil
    }
    
    // 4. Apply updates with user confirmation
    if !confirmUpdate() {
        fmt.Println("❌ Update cancelled")
        return nil
    }
    
    return applyUpdates(projectPath, changes, selective)
}
```

**CLI Usage**:
```bash
# Update project to latest template version
template-health-endpoint update --project ./my-service --template-version v1.2.0

# Selective update (only Kubernetes configs)
template-health-endpoint update --project ./my-service --selective kubernetes,docs

# Dry run to preview changes
template-health-endpoint update --project ./my-service --template-version v1.2.0 --dry-run
```

#### **Task 2.2: CLI Migrate Command**
**Objective**: Enable migrating projects between template tiers

**Implementation**:
```go
// cmd/generator/commands/migrate.go
var migrateCmd = &cobra.Command{
    Use:   "migrate",
    Short: "Migrate project between template tiers",
    Long:  "Migrate a project from one template tier to another (e.g., basic to intermediate).",
    RunE:  runMigrateProject,
}

func runMigrateProject(cmd *cobra.Command, args []string) error {
    projectPath, _ := cmd.Flags().GetString("project")
    fromTier, _ := cmd.Flags().GetString("from")
    toTier, _ := cmd.Flags().GetString("to")
    dryRun, _ := cmd.Flags().GetBool("dry-run")
    
    // 1. Detect current project tier
    if fromTier == "" {
        detectedTier, err := detectProjectTier(projectPath)
        if err != nil {
            return fmt.Errorf("failed to detect project tier: %w", err)
        }
        fromTier = detectedTier
    }
    
    // 2. Validate migration path
    if !isValidMigrationPath(fromTier, toTier) {
        return fmt.Errorf("invalid migration path: %s -> %s", fromTier, toTier)
    }
    
    // 3. Plan migration steps
    migrationPlan, err := createMigrationPlan(fromTier, toTier)
    if err != nil {
        return fmt.Errorf("failed to create migration plan: %w", err)
    }
    
    // 4. Show migration plan
    fmt.Printf("🚀 Migration plan: %s -> %s\n", fromTier, toTier)
    for _, step := range migrationPlan.Steps {
        fmt.Printf("  %s %s\n", step.Action, step.Description)
    }
    
    if dryRun {
        fmt.Println("🔍 Dry run mode - no changes applied")
        return nil
    }
    
    // 5. Execute migration
    return executeMigration(projectPath, migrationPlan)
}
```

**CLI Usage**:
```bash
# Migrate from basic to intermediate
template-health-endpoint migrate --project ./my-service --from basic --to intermediate

# Auto-detect current tier and migrate to advanced
template-health-endpoint migrate --project ./my-service --to advanced

# Dry run migration
template-health-endpoint migrate --project ./my-service --to enterprise --dry-run
```

#### **Task 2.3: CLI Customize Command**
**Objective**: Enable interactive template customization

**Implementation**:
```go
// cmd/generator/commands/customize.go
var customizeCmd = &cobra.Command{
    Use:   "customize",
    Short: "Customize template before generation",
    Long:  "Interactively customize a template before generating a project.",
    RunE:  runCustomizeTemplate,
}

func runCustomizeTemplate(cmd *cobra.Command, args []string) error {
    tier, _ := cmd.Flags().GetString("tier")
    interactive, _ := cmd.Flags().GetBool("interactive")
    configFile, _ := cmd.Flags().GetString("config")
    saveProfile, _ := cmd.Flags().GetString("save-profile")
    
    // 1. Load template
    template, err := loadTemplate(tier)
    if err != nil {
        return fmt.Errorf("failed to load template: %w", err)
    }
    
    // 2. Load customization config
    var customization *TemplateCustomization
    if configFile != "" {
        customization, err = loadCustomizationConfig(configFile)
        if err != nil {
            return fmt.Errorf("failed to load customization config: %w", err)
        }
    } else if interactive {
        customization, err = runInteractiveCustomization(template)
        if err != nil {
            return fmt.Errorf("interactive customization failed: %w", err)
        }
    }
    
    // 3. Apply customizations
    customizedTemplate, err := applyCustomizations(template, customization)
    if err != nil {
        return fmt.Errorf("failed to apply customizations: %w", err)
    }
    
    // 4. Save profile if requested
    if saveProfile != "" {
        if err := saveCustomizationProfile(saveProfile, customization); err != nil {
            return fmt.Errorf("failed to save profile: %w", err)
        }
        fmt.Printf("✅ Customization profile saved: %s\n", saveProfile)
    }
    
    // 5. Generate customized project
    return generateFromCustomizedTemplate(customizedTemplate)
}
```

**CLI Usage**:
```bash
# Interactive customization
template-health-endpoint customize --tier basic --interactive

# Use customization config file
template-health-endpoint customize --tier advanced --config custom.yaml

# Save customization profile
template-health-endpoint customize --tier enterprise --interactive --save-profile my-org
```

### **Phase 3: Final Validation and Testing (1 hour)**

#### **Task 3.1: Comprehensive Testing**
**Objective**: Ensure all template tiers and CLI commands work correctly

**Testing Strategy**:
```bash
#!/bin/bash
# test-complete-system.sh

echo "🧪 Testing Complete Template Health Endpoint System"

# Test all template tiers
for tier in basic intermediate advanced enterprise; do
    echo "Testing $tier tier..."
    
    # Generate from static template
    ./bin/template-health-endpoint template from-static \
        --name "test-$tier" \
        --tier "$tier" \
        --module "github.com/test/$tier" \
        --output "test-$tier"
    
    # Test compilation
    cd "test-$tier"
    go mod tidy
    go build -o "bin/test-$tier" cmd/server/main.go
    
    # Test runtime (start server and test endpoints)
    ./bin/test-$tier &
    SERVER_PID=$!
    sleep 3
    
    # Test all health endpoints
    curl -f http://localhost:8080/health
    curl -f http://localhost:8080/health/time
    curl -f http://localhost:8080/health/ready
    curl -f http://localhost:8080/health/live
    curl -f http://localhost:8080/health/startup
    
    # Test tier-specific endpoints
    if [[ "$tier" != "basic" ]]; then
        curl -f http://localhost:8080/health/dependencies || echo "Dependencies endpoint not configured (expected)"
    fi
    
    if [[ "$tier" == "advanced" || "$tier" == "enterprise" ]]; then
        curl -f http://localhost:8080/health/metrics || echo "Metrics endpoint not configured (expected)"
    fi
    
    # Stop server and clean up
    kill $SERVER_PID
    cd ..
    rm -rf "test-$tier"
    
    echo "✅ $tier tier test passed"
done

# Test CLI commands
echo "Testing CLI commands..."
./bin/template-health-endpoint template list
./bin/template-health-endpoint template validate
./bin/template-health-endpoint --help

echo "🎉 All tests passed!"
```

#### **Task 3.2: Documentation Updates**
**Objective**: Update all documentation to reflect completed features

**Updates Required**:
1. **Update tier-comparison.md**: Add all 4 tiers with complete feature matrix
2. **Update setup-guide.md**: Include all CLI commands and examples
3. **Update migration-guide.md**: Add migration paths between all tiers
4. **Update README.md**: Reflect completed project status
5. **Create CLI reference documentation**: Complete command reference

#### **Task 3.3: GitHub Issue #127 Validation**
**Objective**: Ensure all original requirements are fulfilled

**Validation Checklist**:
- [ ] Template repository following template-* pattern ✅
- [ ] Static template directories users can copy/fork ✅
- [ ] 4 tiers: Basic ✅, Intermediate 🔄, Advanced 🔄, Enterprise 🔄
- [ ] TypeSpec-first API definitions ✅
- [ ] Go server implementations with health endpoints ✅
- [ ] TypeScript client SDKs ✅
- [ ] Kubernetes deployment configurations ✅
- [ ] OpenTelemetry integration (Advanced/Enterprise tiers) 🔄
- [ ] CloudEvents support (Advanced/Enterprise tiers) 🔄
- [ ] Progressive complexity with clear upgrade paths 🔄
- [ ] CLI tool for generation and management ✅
- [ ] Comprehensive documentation ✅

## 🔧 **Technical Implementation Details**

### **Template Tier Feature Matrix**

```yaml
Basic:
  features: {kubernetes: true, typescript: true, docker: true}
  endpoints: ["/health", "/health/time", "/health/ready", "/health/live", "/health/startup"]
  
Intermediate:
  features: {kubernetes: true, typescript: true, docker: true, opentelemetry: basic, dependencies: true}
  endpoints: [...basic, "/health/dependencies"]
  
Advanced:
  features: {kubernetes: true, typescript: true, docker: true, opentelemetry: full, cloudevents: true, server_timing: true}
  endpoints: [...intermediate, "/health/metrics"]
  
Enterprise:
  features: {kubernetes: true, typescript: true, docker: true, opentelemetry: full, cloudevents: true, server_timing: true, security: mtls, compliance: true}
  endpoints: [...advanced, "/health/compliance"]
```

### **CLI Command Architecture**

```bash
template-health-endpoint
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
    metadata, err := tp.loadMetadata()
    if err != nil {
        return err
    }
    
    // 2. Process template variables
    context := tp.createTemplateContext(metadata)
    
    // 3. Apply tier-specific features
    if err := tp.applyTierFeatures(context); err != nil {
        return err
    }
    
    // 4. Generate output files
    if err := tp.generateFiles(context); err != nil {
        return err
    }
    
    // 5. Validate generated project
    return tp.validateOutput()
}
```

## 📊 **Success Criteria**

### **Functional Requirements**
- [ ] All 4 template tiers (Basic, Intermediate, Advanced, Enterprise) generate working projects
- [ ] All generated projects compile and run successfully
- [ ] All health endpoints respond correctly for each tier
- [ ] CLI update command works for version upgrades
- [ ] CLI migrate command works for tier transitions
- [ ] CLI customize command enables template customization
- [ ] Kubernetes deployments work for all tiers
- [ ] OpenTelemetry integration works in intermediate/advanced/enterprise tiers
- [ ] CloudEvents emission works in advanced/enterprise tiers
- [ ] mTLS and compliance features work in enterprise tier

### **Quality Requirements**
- [ ] Comprehensive test coverage for all tiers and CLI commands
- [ ] Documentation is complete and accurate for all features
- [ ] Error handling is robust and user-friendly
- [ ] Performance is acceptable for all operations
- [ ] Security best practices are implemented

### **GitHub Issue #127 Compliance**
- [ ] All original requirements from GitHub issue #127 are fulfilled
- [ ] Template repository structure follows template-* pattern
- [ ] Progressive complexity is implemented correctly
- [ ] CLI tool provides both generation and management capabilities
- [ ] Documentation is comprehensive and user-friendly

## 🚀 **Getting Started**

### **Prerequisites**
- Current BMAD-METHOD repository with basic tier completed ✅
- Go 1.21+ development environment ✅
- Docker (for testing containerization) ✅
- kubectl (for testing Kubernetes deployments) ✅
- Node.js 16+ (for TypeScript client testing) ✅

### **Implementation Order**
1. **Start with Intermediate Tier**: Build upon solid basic tier foundation
2. **Add Advanced Tier**: Implement full observability features
3. **Create Enterprise Tier**: Add security and compliance features
4. **Implement CLI Commands**: Add update, migrate, customize functionality
5. **Comprehensive Testing**: Validate all tiers and CLI commands
6. **Final Documentation**: Ensure all documentation is current and complete

### **Validation Steps**
1. Generate projects from all template tiers
2. Compile and run all generated projects
3. Test all health endpoints for each tier
4. Test CLI commands with real projects
5. Validate Kubernetes deployments
6. Run comprehensive test suite
7. Verify GitHub issue #127 requirements are met

## 📚 **Context and Background**

### **Project Status**
- **60% Complete**: Solid foundation with basic tier, dual template system, comprehensive testing
- **40% Remaining**: 3 additional template tiers + advanced CLI functionality
- **High Quality Foundation**: Excellent architecture, testing, and documentation already in place

### **Original GitHub Issue #127 Requirements**
The original issue specifically requested:
- Comprehensive template repository for health endpoint APIs
- 4 template tiers with progressive complexity
- TypeSpec-first API definitions with code generation
- Go server implementations with OpenTelemetry integration
- TypeScript client SDKs
- Kubernetes deployment configurations
- Both static templates AND CLI tool functionality

### **Current Achievement**
- ✅ Dual template system (static + CLI) architecture implemented
- ✅ Basic tier fully functional with all health endpoints including `/health/startup`
- ✅ Comprehensive testing framework with 100% success rate
- ✅ Complete documentation and guidelines
- ✅ TypeSpec schemas and code generation pipeline
- ✅ Repository structure following template-* pattern

### **Why This Task is Critical**
1. **Completes Original Requirements**: Fulfills all GitHub issue #127 specifications
2. **Maximum User Value**: Provides complete progression path from basic to enterprise
3. **Production Ready**: Delivers enterprise-grade templates with security and compliance
4. **CLI Completeness**: Adds essential update and migration functionality
5. **Project Completion**: Takes project from foundation to fully delivered product

## 🎯 **Expected Outcomes**

Upon completion of this task:

1. **Complete Template Repository**: All 4 tiers available as static templates
2. **Full CLI Functionality**: Generate, update, migrate, customize commands working
3. **Production-Ready Templates**: Enterprise-grade templates with security, compliance, and observability
4. **Comprehensive Documentation**: Complete guides for all tiers and features
5. **Validated Quality**: All templates tested and working correctly
6. **GitHub Issue #127 Fulfilled**: All original requirements completed
7. **Ready for Dedicated Repository**: Project ready to be merged into new dedicated repo

This task represents the final 40% of work needed to complete the template-health-endpoint project and fulfill all requirements from the original GitHub issue #127. The foundation is excellent, and this task will deliver a production-ready, comprehensive template system that serves both manual and automated workflows.

---

**Ready to implement? Start with Phase 1, Task 1.1: Intermediate Template Tier!** 🚀

**This is the most important and impactful task to complete the entire project and fulfill all original requirements.**
