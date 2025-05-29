# BMAD-METHOD Completion Guidelines

## Overview
Guidelines for completing sophisticated template generation systems with enterprise-grade features, comprehensive testing, and production deployment.

## Core Principles

### 1. Completion-Driven Development
- **Fix Critical Issues First**: Address compilation errors and test failures before adding features
- **Validate Continuously**: Run integration tests after every change
- **Document Progress**: Maintain clear status of what's working vs. what needs fixing
- **Measure Success**: Use quantitative metrics (test pass rates, compilation success)

### 2. Template System Quality
- **Clean Compilation**: All generated code must compile without warnings
- **Unused Import Detection**: Remove all unused imports from template code
- **Comprehensive Testing**: Test template generation, compilation, and runtime functionality
- **Multi-Tier Validation**: Ensure all complexity tiers work correctly

### 3. Integration Testing Excellence
- **End-to-End Validation**: Test complete workflow from generation to running application
- **Structure Validation**: Verify generated project structure matches expectations
- **Runtime Testing**: Validate that generated applications actually work
- **Performance Metrics**: Measure generation speed and resource usage

## Implementation Guidelines

### Template Debugging Process
1. **Identify Compilation Errors**
   ```bash
   # Generate project and test compilation
   ./bin/tool generate --name test --tier enterprise --output test-output
   cd test-output && go mod tidy && go build ./...
   ```

2. **Locate Template Issues**
   - Find inline templates in generator code
   - Check import statements for unused imports
   - Validate template variable substitution

3. **Fix Template Code**
   ```go
   // Remove unused imports
   import (
       "net/http"     // Used for http.Handler
       "strings"      // Used for strings.HasPrefix()
       // Remove "fmt" unless calling fmt functions
       // Remove "encoding/json" unless calling json functions
   )
   ```

4. **Update Integration Tests**
   - Check what files are actually generated
   - Update test expectations to match reality
   - Remove checks for non-existent files

### Testing Strategy
```bash
# Comprehensive testing approach
test_all_tiers() {
    for tier in basic intermediate advanced enterprise; do
        # Generate project
        ./bin/tool generate --name ${tier}-test --tier $tier --output test-output/${tier}
        
        # Test compilation
        (cd test-output/${tier} && go mod tidy && go build ./...)
        
        # Test runtime (for enterprise)
        if [ "$tier" = "enterprise" ]; then
            test_runtime_functionality test-output/${tier}
        fi
    done
}
```

### Error Resolution Patterns
- **Unused Import "encoding/json"**: Only needed for json.Marshal/Unmarshal calls, not struct tags
- **Unused Import "fmt"**: Only needed for fmt.Printf/Sprintf calls
- **Unused Import "context"**: Only needed for context.Context type or context package functions
- **Missing Files in Tests**: Update tests to check for actually generated files

## Quality Assurance

### Success Metrics
- **Test Pass Rate**: 100% of integration tests must pass
- **Compilation Success**: All generated projects compile without warnings
- **Runtime Validation**: Generated applications start and respond correctly
- **Performance**: Generation completes in under 5 seconds per tier

### Validation Checklist
- [ ] All template imports are used
- [ ] Generated projects compile cleanly
- [ ] Integration tests pass completely
- [ ] Runtime functionality verified
- [ ] Documentation updated
- [ ] Examples generated and tested

## Best Practices

### Template Code Quality
```go
// Good: Only include necessary imports
import (
    "net/http"    // Used for http.Handler
    "strings"     // Used for strings.HasPrefix()
)

// Bad: Including unused imports
import (
    "encoding/json"  // Not used - struct tags don't require import
    "fmt"           // Not used - no fmt function calls
    "net/http"      // Used
    "strings"       // Used
)
```

### Integration Test Design
```bash
# Good: Test actual generated structure
if [[ -f "project/internal/security/rbac.go" && \
      -f "project/internal/security/mtls.go" ]]; then
    log_success "Enterprise structure correct"
fi

# Bad: Test for files that aren't generated
if [[ -f "project/configs/development.yaml" ]]; then
    log_success "Config files present"  # May not exist
fi
```

### Error Handling
```bash
# Provide detailed error information
test_compilation() {
    if (cd "$project_dir" && go build ./... > /dev/null 2>&1); then
        log_success "Project compiles successfully"
    else
        log_error "Project compilation failed"
        # Show actual errors for debugging
        (cd "$project_dir" && go build ./... 2>&1 | head -10)
    fi
}
```

## Completion Workflow

### Phase 1: Issue Identification (10 minutes)
1. Run integration tests to identify failures
2. Check compilation errors in generated projects
3. Locate specific template issues

### Phase 2: Template Fixes (20 minutes)
1. Remove unused imports from templates
2. Fix template variable substitution issues
3. Ensure all generated code is valid

### Phase 3: Test Updates (10 minutes)
1. Update integration tests to match actual generation
2. Remove checks for non-existent files
3. Add validation for actually generated structure

### Phase 4: Comprehensive Validation (20 minutes)
1. Run full integration test suite
2. Test all tiers for generation and compilation
3. Verify runtime functionality
4. Update documentation and examples

## Anti-Patterns to Avoid

### Template Issues
- ❌ Including imports "just in case"
- ❌ Not testing all template tiers after changes
- ❌ Assuming struct tags require package imports
- ❌ Hardcoding test expectations

### Testing Issues
- ❌ Testing for files that may not be generated
- ❌ Not providing detailed error information
- ❌ Skipping runtime validation
- ❌ Not testing edge cases

### Process Issues
- ❌ Making multiple changes without testing
- ❌ Not documenting what was fixed
- ❌ Ignoring compilation warnings
- ❌ Not validating end-to-end workflows

## Success Indicators

### Quantitative Metrics
- **17/17 integration tests passing**
- **0 compilation warnings in generated code**
- **< 5 second generation time per tier**
- **100% template validation success**

### Qualitative Indicators
- Generated projects work immediately after creation
- Clear, actionable error messages when issues occur
- Comprehensive documentation and examples
- Smooth user experience from generation to deployment

This completion-focused approach ensures high-quality, production-ready template systems that work reliably in real-world scenarios.
