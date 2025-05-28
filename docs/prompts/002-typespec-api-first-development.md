# TypeSpec API-First Development

## Prompt Name: TypeSpec API-First Development

## Description
Create comprehensive API definitions using TypeSpec with automatic code generation for multiple languages and deployment targets.

## When to Use
- Building APIs that need multiple client SDKs
- Want type-safe API development across languages
- Need automatic OpenAPI and JSON Schema generation
- Building microservices with consistent API patterns

## Prompt

```
Create a comprehensive TypeSpec-first API system for [API_NAME] with the following requirements:

**TypeSpec Schema Design:**
- Design complete TypeSpec models for all data structures
- Create HTTP interfaces with proper decorators
- Include comprehensive documentation and examples
- Support progressive complexity (basic → advanced tiers)
- Implement proper error handling and status codes

**Code Generation Targets:**
- Generate Go server implementation with HTTP handlers
- Create TypeScript client SDK with full type safety
- Generate JSON Schema for validation
- Produce OpenAPI v3 specification for documentation
- Include Kubernetes manifests with health probes

**Quality Requirements:**
- All schemas must compile without errors
- Generated code must be production-ready
- Include comprehensive test coverage
- Support multiple deployment environments
- Follow enterprise security patterns

**Architecture Patterns:**
- Use schema-first development approach
- Implement proper separation of concerns
- Include observability and monitoring
- Support graceful error handling
- Enable easy testing and validation

**Deliverables:**
- Complete TypeSpec schema definitions
- Working Go server implementation
- TypeScript client SDK with npm package
- Docker containerization
- Kubernetes deployment manifests
- Comprehensive API documentation

Ensure all generated code compiles, runs, and passes tests. Create working examples and integration guides.
```

## Expected Outcomes
- Complete TypeSpec schema suite
- Multi-language code generation
- Production-ready API implementation
- Client SDKs and documentation
- Deployment configurations

## Success Criteria
- TypeSpec schemas compile successfully
- Generated code works without modification
- API endpoints respond correctly
- Client SDKs integrate properly
- Documentation is comprehensive and accurate
