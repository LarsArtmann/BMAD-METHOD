# PROJECT GUIDELINES: BMAD Method & Template Health Endpoint

## Overview

This document consolidates all learnings, best practices, and guidelines from the BMAD Method implementation for the template-health-endpoint project. It serves as the definitive guide for AI agents and developers working on complex software development projects.

## Table of Contents

1. [BMAD Method Workflow](#bmad-method-workflow)
2. [Development Principles](#development-principles)
3. [Technical Architecture Guidelines](#technical-architecture-guidelines)
4. [Code Quality Standards](#code-quality-standards)
5. [Testing and Validation](#testing-and-validation)
6. [Documentation Standards](#documentation-standards)
7. [AI Agent Collaboration](#ai-agent-collaboration)
8. [Technology Stack Guidelines](#technology-stack-guidelines)
9. [Project Organization](#project-organization)
10. [Quality Assurance](#quality-assurance)

---

## BMAD Method Workflow

### Systematic Development Methodology

For complex projects, follow the BMAD Method workflow:

1. **Analyst Phase (Larry)**: Create comprehensive project brief with problem analysis
2. **Product Manager Phase (John)**: Develop detailed PRD with epics and user stories  
3. **Architect Phase (Mo)**: Design technical architecture and component structure
4. **Product Owner Phase (PO)**: Validate requirements and create acceptance criteria
5. **Scrum Master Phase**: Break epics into 5-task manageable stories
6. **Developer Phase**: Implement with incremental progress and validation

### Phase Responsibilities

**Analyst Agent:**
- Analyze problem statement and create project brief
- Define vision, goals, and success metrics
- Identify target audience and user personas
- Research relevant technologies and patterns
- Document constraints and preferences

**Product Manager:**
- Transform project brief into detailed PRD
- Create 4 comprehensive epics with user stories
- Define functional and non-functional requirements
- Establish clear acceptance criteria
- Plan MVP scope and post-MVP features

**Architect:**
- Design technical architecture based on PRD
- Create component diagrams and data models
- Define technology stack with rationale
- Plan integration patterns and workflows
- Document security and scalability considerations

**Product Owner:**
- Validate all documents against master checklist
- Ensure proper sequencing and dependencies
- Verify MVP scope alignment
- Check risk management and feasibility
- Approve for development phase

**Scrum Master:**
- Break epics into exactly 5 small, manageable stories
- Create detailed task breakdowns
- Ensure stories follow BMAD template structure
- Plan sprint organization and dependencies

**Developer:**
- Implement each story with 5 focused changes
- Follow architecture and design patterns
- Create working, tested code
- Generate comprehensive documentation
- Validate against acceptance criteria

### Quality Gates

- Complete each phase before proceeding to ensure proper foundation
- Validate deliverables against phase-specific checklists
- Ensure all stakeholders approve before phase transition
- Maintain comprehensive documentation throughout

---

## Development Principles

### Incremental Development Approach

**5-Change Rule**: Each story must have exactly 5 focused changes:
1. Each change should be completable in 10-15 minutes
2. Focus on single responsibility per change
3. Include testing and validation in each change
4. Document rationale for complex decisions
5. Maintain working state after each change

**Quality First Approach**:
- "Keep going until everything works and you think you did a great job"
- Validate each change before proceeding
- Test in real environments, not just theory
- Ensure production-ready quality at each step

### Schema-First Development

For API and data-driven projects:
1. **Define Schemas First**: Use TypeSpec, JSON Schema, or similar for API definitions
2. **Generate Code**: Use schema-driven code generation for consistency
3. **Multi-Language Support**: Generate clients and servers from single schema
4. **Validation Integration**: Include schema validation in generated code
5. **Documentation Generation**: Auto-generate API documentation from schemas

**Benefits**: Ensures consistency, reduces errors, enables multi-language support

### Progressive Complexity Implementation

For systems with multiple use cases:
1. **Tier Structure**: Design Basic → Intermediate → Advanced → Enterprise tiers
2. **Clear Value Proposition**: Each tier should provide obvious value over previous
3. **Feature Flags**: Use configuration-driven feature enablement
4. **Upgrade Paths**: Provide clear migration between tiers
5. **Time Targets**: Set realistic deployment time goals for each tier

**Implementation Strategy**:
- Start with simplest tier that provides value
- Build complexity incrementally
- Maintain backward compatibility
- Test upgrade paths thoroughly

---

## Technical Architecture Guidelines

### Template Generation Architecture

**Core Components**:
- TypeSpec Schema Registry for canonical definitions
- Go-based Template Generator with CLI interface
- Multi-language Code Generators (Go, TypeScript)
- Kubernetes Manifest Generator
- Observability Integration Layer

**Design Patterns**:
- Template Generation Pattern for consistent implementations
- Schema-First Development for type safety
- Progressive Complexity Tiers for scalability
- Observability-First Design for monitoring
- Cloud-Native Patterns for deployment

### TypeSpec Integration

**Schema Organization**:
- Place schemas in `pkg/schemas/` directory structure
- Use separate files for logical grouping
- Follow TypeSpec naming conventions and import patterns
- Include comprehensive documentation and examples

**Code Generation Pipeline**:
- TypeSpec Schemas → JSON Schema/OpenAPI → Generated Code
- Validate schemas before code generation
- Generate complete project structures
- Include testing frameworks and examples

### CLI Tool Development

**Framework Standards**:
- Use Cobra for Go CLI tools, similar frameworks for other languages
- Support YAML/JSON config files with environment overrides
- Include help, verbose mode, dry-run capabilities
- Provide clear, actionable error messages
- Use language-native templating for code generation

**User Experience Requirements**:
- Comprehensive help documentation
- Progress indicators for long operations
- Colored output for better UX
- Configuration validation with helpful suggestions

---

## Code Quality Standards

### Generated Code Quality

**Requirements**:
- Generated code must compile without errors
- Include comprehensive error handling
- Follow language-specific best practices
- Include testing frameworks and examples
- Generate complete project structures

**Validation**:
- Test generated code in real environments
- Validate against actual use cases
- Ensure performance meets requirements
- Check security and compliance standards

### Performance Requirements

**API Performance**:
- Health endpoints must respond within 100ms
- ServerTime API must include sub-millisecond timing accuracy
- Template generation must complete within 30 seconds for enterprise tier

**Resource Optimization**:
- Generated applications should meet performance benchmarks
- Resource usage should be optimized
- Container images should be minimal and secure

### Security Standards

**Generated Code Security**:
- No sensitive information exposed in health endpoints
- Optional authentication support for enterprise tier
- Rate limiting capabilities for health endpoints
- Follow security best practices for each language

---

## Testing and Validation

### Comprehensive Testing Strategy

**Multi-Level Testing**:
1. **Unit Tests**: Test individual components and functions
2. **Integration Tests**: Validate component interactions
3. **End-to-End Tests**: Test complete user workflows
4. **Generated Code Tests**: Validate generated code compiles and runs
5. **Real-World Validation**: Test in actual deployment environments

**Validation Requirements**:
- Schema validation for all definitions
- Performance benchmarking against requirements
- Security scanning and compliance checking
- Cross-platform compatibility testing
- Documentation accuracy verification

### Testing Framework

**Automated Testing**:
- Continuous integration pipeline setup
- Automated test execution on code changes
- Performance regression testing
- Security vulnerability scanning
- Dependency vulnerability checking

**Quality Metrics**:
- Code coverage requirements (minimum 80%)
- Performance benchmarks and SLAs
- Security compliance validation
- Documentation completeness checks
- API response time validation

---

## Documentation Standards

### Comprehensive Documentation

**Required Documentation**:
- README with quick start guide
- Architecture documentation with diagrams
- API documentation with examples
- Setup and deployment guides
- Troubleshooting and FAQ sections

**Documentation Quality**:
- Include working examples and verification steps
- Provide clear setup and usage instructions
- Document architectural decisions and rationale
- Include troubleshooting guides for common issues
- Maintain up-to-date content

### AI Agent Handoff Documentation

**Handoff Document Structure**:
- Mission statement and objectives
- Current status with completion percentages
- Verified working features with examples
- Prioritized next steps with acceptance criteria
- Technical reference and troubleshooting guide

**Requirements**:
- Complete context preservation
- Immediate actionability for new team member
- Comprehensive technical reference
- Clear quality standards and success criteria

---

## AI Agent Collaboration

### Collaboration Patterns

**Effective Practices**:
- Break work into 5 small changes per story
- Maintain comprehensive context documentation
- Use "Keep going until everything works" approach
- Create detailed handoff procedures
- Validate each increment before proceeding

**Context Management**:
- Preserve project context across conversations
- Document decisions and rationale
- Maintain working state information
- Track progress against objectives

### Quality vs Speed Balance

**Guidelines**:
- Prioritize working solutions over perfect code
- Validate in real environments, not just theory
- Focus on user value and practical outcomes
- Maintain high standards while making progress

---

## Technology Stack Guidelines

### Proven Technology Combinations

**API Development**:
- TypeSpec + Go + TypeScript for multi-language consistency
- OpenTelemetry + Prometheus + Grafana for observability
- CloudEvents for event-driven architecture

**CLI Tools**:
- Go + Cobra + Viper for robust command-line interfaces
- Template-based code generation
- Configuration-driven customization

**Container Deployment**:
- Docker + Kubernetes with proper health probes
- Multi-stage builds for minimal images
- Security-hardened containers

**Testing**:
- Language-native frameworks + integration testing
- Real-world validation scenarios
- Automated CI/CD pipelines

### Selection Criteria

**Technology Evaluation**:
- Maturity and community support
- Integration capabilities
- Performance characteristics
- Learning curve and documentation
- Long-term maintenance considerations

---

## Project Organization

### Directory Structure

```
project/
├── cmd/                    # CLI applications and entry points
├── pkg/                    # Reusable packages and libraries
├── internal/               # Private application code
├── docs/                   # Comprehensive documentation
│   ├── prompts/           # Reusable prompt templates
│   ├── learnings/         # Project learnings and insights
│   ├── stories/           # BMAD Method user stories
│   └── guideline-suggestions/ # Process improvement suggestions
├── examples/              # Working examples and demonstrations
├── scripts/               # Build and utility scripts
├── tests/                 # Test suites and fixtures
├── pkg/schemas/           # TypeSpec schema definitions
├── deployments/           # Kubernetes and deployment configs
└── generated/             # Generated code and artifacts
```

### File Organization

**Naming Conventions**:
- Use clear, descriptive names
- Follow language-specific conventions
- Organize by domain/feature over technical layers
- Separate concerns between components

**Documentation Organization**:
- Keep documentation close to relevant code
- Use consistent formatting and structure
- Include examples and verification steps
- Maintain comprehensive cross-references

---

## Quality Assurance

### Success Metrics

**Technical Metrics**:
- All generated code compiles and runs successfully
- Health endpoints respond within specified timeframes
- TypeScript clients integrate properly
- Kubernetes deployments work correctly
- Documentation enables target deployment times

**Process Metrics**:
- BMAD Method phases completed systematically
- Stories broken into exactly 5 manageable tasks
- Quality gates passed at each phase
- Comprehensive documentation maintained
- Real-world validation completed

### Continuous Improvement

**Learning Capture**:
- Document lessons learned from each project
- Create reusable prompts for common patterns
- Update guidelines based on experience
- Share knowledge across team members

**Process Refinement**:
- Regular review of methodology effectiveness
- Update templates and checklists
- Improve tooling and automation
- Enhance documentation standards

---

## Implementation Checklist

### Project Initiation
- [ ] Define clear project objectives and scope
- [ ] Identify target users and use cases
- [ ] Establish success criteria and metrics
- [ ] Set up project structure and documentation

### BMAD Method Execution
- [ ] Complete Analyst phase with comprehensive brief
- [ ] Develop detailed PRD with epics and stories
- [ ] Design technical architecture and components
- [ ] Validate requirements and acceptance criteria
- [ ] Break epics into 5-task stories
- [ ] Implement with incremental validation

### Quality Validation
- [ ] Test generated code in real environments
- [ ] Validate performance against requirements
- [ ] Check security and compliance standards
- [ ] Verify documentation accuracy and completeness
- [ ] Confirm deployment procedures work correctly

### Project Completion
- [ ] Create comprehensive handoff documentation
- [ ] Document lessons learned and improvements
- [ ] Archive project artifacts and knowledge
- [ ] Plan maintenance and evolution strategy

---

## Conclusion

These guidelines represent the consolidated wisdom from successful implementation of complex software projects using the BMAD Method. They provide a systematic approach to software development that ensures quality, maintainability, and successful outcomes.

**Key Success Factors**:
1. **Systematic Approach**: Follow BMAD Method for complex projects
2. **Quality Focus**: Maintain high standards throughout development
3. **Incremental Progress**: Break work into manageable pieces
4. **Comprehensive Testing**: Validate in real-world scenarios
5. **Documentation Excellence**: Enable knowledge transfer and continuity

By following these guidelines, teams can achieve consistent, high-quality results while maintaining development velocity and ensuring long-term project success.
