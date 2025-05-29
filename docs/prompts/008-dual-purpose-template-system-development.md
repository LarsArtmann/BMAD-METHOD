# Dual-Purpose Template System Development

## Prompt Name
**Dual-Purpose Template System Development**

## Context
Build a template system that serves both as a static template repository (following template-* repository patterns) AND as a CLI generator/updater tool. This provides maximum flexibility for users who want either manual customization or automated workflows.

## Objective
Create a comprehensive template system that:
1. **Static Templates**: Users can fork/copy template directories directly
2. **CLI Generator**: Users can generate fresh projects programmatically  
3. **CLI Updater**: Users can update existing projects to newer template versions
4. **Template Repository Structure**: Follows established template-* repository patterns

## Key Requirements

### Repository Structure
```
template-{name}/
├── README.md                     # Comprehensive documentation
├── LICENSE                       # EUPL-1.2
├── templates/                    # 🆕 Static template directories
│   ├── basic/                    # Users can copy this directly
│   │   ├── cmd/server/main.go
│   │   ├── internal/handlers/health.go
│   │   ├── go.mod.tmpl
│   │   ├── README.md.tmpl
│   │   ├── template.yaml         # Template metadata
│   │   └── ...
│   ├── intermediate/             # Production-ready template
│   ├── advanced/                 # Full observability template
│   └── enterprise/               # Kubernetes & compliance template
├── examples/                     # Generated examples from templates
├── scripts/                      # Template operations
├── cmd/                          # Enhanced CLI tool
│   └── generator/
│       ├── main.go               # CLI entry point
│       ├── generate.go           # Generate new projects
│       ├── update.go             # Update existing projects
│       ├── migrate.go            # Migrate between tiers
│       └── validate.go           # Template validation
├── pkg/                          # Core logic
└── docs/                         # Comprehensive documentation
```

### CLI Commands Design
```bash
# Generate new project from template
template-{name} generate --name my-service --tier basic --output ./my-service

# Generate from static template directory
template-{name} template from-static --name my-service --tier basic --module github.com/org/service

# Update existing project to newer template version
template-{name} update --project ./my-service --template-version v1.2.0

# Migrate project to different tier
template-{name} migrate --project ./my-service --from basic --to intermediate

# List available templates
template-{name} template list

# Validate template directories
template-{name} template validate
```

### Template Processing Requirements
1. **Template Variables**: Support Go template syntax with proper variable substitution
2. **File Extensions**: Process `.tmpl` files, `.go` files, `.yaml/.yml`, `.json`, `.ts`, `.sh`, and specific files like `go.mod`, `README.md`
3. **Template Metadata**: Each template tier has `template.yaml` with metadata
4. **Progressive Complexity**: Basic → Intermediate → Advanced → Enterprise tiers

### Template Metadata Format
```yaml
name: basic
description: Basic tier health endpoint template
tier: basic
features:
  kubernetes: true
  typescript: true
  docker: true
  opentelemetry: false
  cloudevents: false
version: "1.0.0"
```

## Implementation Steps

### Phase 1: Repository Restructuring
1. Create static template directories from existing CLI templates
2. Convert hardcoded values to template variables
3. Add template metadata files
4. Modify CLI to read from template directories instead of embedded templates

### Phase 2: CLI Enhancement
1. Add `template` command group with subcommands
2. Implement `from-static` generation from template directories
3. Add template listing and validation
4. Ensure proper template variable substitution

### Phase 3: Advanced Features
1. Add `update` command for existing projects
2. Implement `migrate` command for tier transitions
3. Add comprehensive testing and validation
4. Create documentation and examples

## Success Criteria
- [ ] Static template directories can be copied/forked directly
- [ ] CLI generates working projects from static templates
- [ ] All template variables are properly substituted
- [ ] Generated projects compile and run successfully
- [ ] Template validation passes for all tiers
- [ ] Documentation is comprehensive and clear

## Technical Considerations
- Use Go templates for variable substitution
- Support YAML metadata for template configuration
- Implement proper error handling and validation
- Follow established CLI patterns with Cobra
- Ensure cross-platform compatibility

## Benefits
- **For Manual Users**: Can fork/copy template directories directly
- **For Automated Users**: CLI generation for CI/CD pipelines
- **For Ecosystem Integration**: Works with existing template-* repositories
- **Maximum Flexibility**: Supports both manual customization and automated workflows

## Related Patterns
- Template repository patterns (template-*)
- CLI tool development with Cobra
- Go template processing
- Progressive complexity systems
- BMAD Method implementation
