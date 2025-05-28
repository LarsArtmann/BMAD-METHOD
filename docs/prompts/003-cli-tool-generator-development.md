# CLI Tool Generator Development

## Prompt Name: CLI Tool Generator Development

## Description
Build comprehensive CLI tools with template generation capabilities, configuration management, and validation frameworks.

## When to Use
- Creating developer tools and generators
- Building template-based project scaffolding
- Need configuration-driven code generation
- Want to standardize project creation workflows

## Prompt

```
Create a comprehensive CLI tool for [TOOL_NAME] with the following capabilities:

**CLI Framework:**
- Use Cobra for command structure and flag handling
- Implement comprehensive help and usage documentation
- Support configuration files (YAML/JSON)
- Include verbose and quiet output modes
- Add dry-run capabilities for preview

**Template Generation:**
- Create template-based project scaffolding
- Support multiple complexity tiers/templates
- Implement configuration-driven customization
- Generate complete project structures
- Include all necessary build and deployment files

**Validation Framework:**
- Validate input schemas and configurations
- Check generated code compilation
- Verify template consistency
- Test integration with target platforms
- Provide detailed error reporting

**Configuration Management:**
- Support tier-based feature flags
- Enable environment-specific settings
- Allow configuration inheritance and overrides
- Validate configuration completeness
- Provide configuration examples and documentation

**Code Generation:**
- Generate working source code in target languages
- Create build scripts and automation
- Include testing frameworks and examples
- Generate deployment configurations
- Produce comprehensive documentation

**Quality Features:**
- Comprehensive error handling with helpful messages
- Progress indicators for long operations
- Colored output for better UX
- Configuration validation and suggestions
- Template validation and testing

Build the tool incrementally:
1. CLI framework and basic commands
2. Configuration system and validation
3. Template engine and generation
4. Code generation and scaffolding
5. Testing and validation framework

Ensure the tool is production-ready with comprehensive testing and documentation.
```

## Expected Outcomes
- Working CLI tool with full functionality
- Template generation system
- Configuration management
- Validation framework
- Comprehensive documentation

## Success Criteria
- CLI tool builds and runs successfully
- All commands work as expected
- Generated projects compile and run
- Configuration validation works properly
- Documentation enables easy adoption
