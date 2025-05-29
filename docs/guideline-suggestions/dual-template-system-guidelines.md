# Dual Template System Development Guidelines

## Overview
Guidelines for developing template systems that serve both as static template repositories and CLI generator/updater tools, based on learnings from the template-health-endpoint project.

## Core Principles

### 1. Dual-Purpose Architecture
**Guideline**: Always design template systems to serve both manual and automated workflows
- **Static Templates**: Provide copyable template directories for manual customization
- **CLI Tools**: Provide programmatic generation and update capabilities
- **Integration**: Ensure both approaches work seamlessly together

**Rationale**: Different users have different preferences and workflows. Some prefer manual control, others prefer automation.

### 2. Repository Structure Standards
**Guideline**: Follow established template-* repository patterns
```
template-{name}/
├── templates/           # Static template directories (NEW)
│   ├── basic/
│   ├── intermediate/
│   ├── advanced/
│   └── enterprise/
├── examples/            # Generated examples
├── scripts/             # Template operations
├── cmd/                 # CLI tool
├── pkg/                 # Core logic
└── docs/                # Documentation
```

**Rationale**: Consistency with existing ecosystem enables better integration and user adoption.

### 3. Template Variable Processing
**Guideline**: Process template variables comprehensively across all relevant file types
- **File Types**: `.tmpl`, `.go`, `.yaml`, `.yml`, `.json`, `.ts`, `.sh`, `go.mod`, `README.md`
- **Variables**: Use Go template syntax with clear, descriptive names
- **Validation**: Ensure all variables are properly substituted

**Implementation**:
```go
// Check if file needs template processing
if needsTemplateProcessing(inputPath) {
    // Process as template with variable substitution
} else {
    // Copy file as-is
}
```

### 4. Template Metadata Management
**Guideline**: Include structured metadata for each template tier
- **Format**: YAML files (`template.yaml`) in each template directory
- **Content**: Name, description, tier, features, version
- **Usage**: Enable CLI listing, validation, and processing

**Schema**:
```yaml
name: string
description: string
tier: string
features:
  feature_name: boolean
version: string
```

### 5. Progressive Complexity Design
**Guideline**: Design template tiers with clear progression and differentiation
- **Basic**: Core functionality, minimal dependencies
- **Intermediate**: Production features, basic observability
- **Advanced**: Full observability, event-driven features
- **Enterprise**: Compliance, security, multi-environment support

**Feature Matrix**: Clearly document which features are available in each tier.

### 6. CLI Command Structure
**Guideline**: Use hierarchical command structure for complex template operations
```bash
tool-name generate          # Generate new project (embedded templates)
tool-name template list     # List available static templates
tool-name template from-static  # Generate from static template
tool-name template validate # Validate template integrity
tool-name update            # Update existing project
tool-name migrate           # Migrate between tiers
```

**Rationale**: Hierarchical structure scales better and provides clearer organization.

### 7. Comprehensive Testing Strategy
**Guideline**: Implement multi-level testing for template systems
- **Template Validation**: Structure, metadata, required files
- **Generation Testing**: Templates generate working projects
- **Compilation Testing**: Generated projects compile successfully
- **Runtime Testing**: Generated applications work correctly
- **Integration Testing**: End-to-end workflow validation

### 8. User Experience Standards
**Guideline**: Provide clear, helpful feedback throughout template operations
- **Progress Indicators**: Show what's happening during generation
- **Success Confirmation**: Confirm what was created and next steps
- **Error Messages**: Clear, actionable error descriptions
- **Help Text**: Comprehensive help and examples

**Example Output**:
```
🚀 Generating project from basic template...
📋 Using template: Basic tier health endpoint template (1.0.0)
✅ Successfully generated project from basic template!
📁 Project created in: my-project
```

### 9. Documentation Requirements
**Guideline**: Provide comprehensive documentation for template systems
- **Setup Guide**: Installation and basic usage
- **Template Guide**: Available tiers and their differences
- **Migration Guide**: How to upgrade between tiers
- **API Documentation**: Generated code documentation
- **Examples**: Working examples for each tier

### 10. Ecosystem Integration
**Guideline**: Design for integration with existing development ecosystems
- **Package Managers**: Work with language-specific package managers
- **Build Tools**: Integrate with existing build and deployment tools
- **CI/CD**: Support automated workflows and pipelines
- **IDEs**: Provide IDE-friendly project structures

## Implementation Guidelines

### Template Conversion Process
1. **Extract**: Generate projects using existing CLI/templates
2. **Convert**: Replace hardcoded values with template variables
3. **Validate**: Ensure templates work correctly
4. **Document**: Add metadata and documentation
5. **Test**: Comprehensive testing of generated projects

### CLI Development Standards
- Use established CLI frameworks (Cobra for Go)
- Implement proper flag validation and error handling
- Support both interactive and non-interactive modes
- Provide dry-run options for safe testing
- Include comprehensive help text and examples

### Template Variable Naming
- Use descriptive, clear variable names
- Follow consistent naming conventions
- Group related variables logically
- Provide sensible defaults where possible

### Error Handling
- Provide clear, actionable error messages
- Include context about what went wrong and how to fix it
- Validate inputs early and provide immediate feedback
- Log detailed information for debugging while keeping user output clean

## Quality Assurance

### Validation Checklist
- [ ] Template metadata is valid and complete
- [ ] All required files exist in template directories
- [ ] Template variables are properly defined and used
- [ ] Generated projects compile without errors
- [ ] All endpoints/functionality work as expected
- [ ] Documentation is comprehensive and accurate
- [ ] Tests cover all template tiers and scenarios

### Performance Considerations
- Template generation should complete in reasonable time
- CLI operations should provide progress feedback for long operations
- Template processing should be efficient for large projects
- Memory usage should be reasonable for template operations

### Security Considerations
- Validate all user inputs to prevent injection attacks
- Sanitize template variables to prevent code injection
- Use secure defaults for generated configurations
- Include security best practices in generated code

## Maintenance Guidelines

### Version Management
- Use semantic versioning for template versions
- Maintain compatibility matrices between template versions
- Provide migration guides for breaking changes
- Test template updates thoroughly before release

### Community Contributions
- Provide clear contribution guidelines
- Include templates for new template tiers
- Maintain consistent code quality standards
- Document the template development process

### Monitoring and Feedback
- Collect usage analytics (with user consent)
- Monitor common error patterns
- Gather user feedback on template quality
- Continuously improve based on real-world usage

## Conclusion
These guidelines provide a framework for developing high-quality, dual-purpose template systems that serve both manual and automated workflows. Following these principles ensures consistency, quality, and user satisfaction while enabling ecosystem integration and long-term maintainability.
