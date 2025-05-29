# Comprehensive Integration Testing for Template Systems

## Prompt Name: Comprehensive Integration Testing for Template Systems

## Context
You need to implement comprehensive integration testing for a template generation system that creates multi-tier projects. The testing must validate template generation, project compilation, runtime functionality, and CLI commands across all tiers.

## Testing Architecture

### Multi-Level Testing Strategy
```
Integration Testing Levels:
├── CLI Command Testing
│   ├── Help and version commands
│   ├── Template listing and validation
│   └── Error handling and edge cases
├── Template Generation Testing
│   ├── All tier generation (basic, intermediate, advanced, enterprise)
│   ├── Project structure validation
│   └── File content verification
├── Compilation Testing
│   ├── Go mod tidy execution
│   ├── Build success validation
│   └── Warning-free compilation
├── Runtime Testing
│   ├── Server startup validation
│   ├── Endpoint response testing
│   └── Feature functionality verification
└── Advanced Feature Testing
    ├── TypeScript client generation
    ├── Kubernetes manifest validation
    └── BDD framework execution
```

## Implementation Pattern

### 1. Test Script Structure
```bash
#!/bin/bash
# integration_test.sh

# Test configuration
TESTS_PASSED=0
TESTS_FAILED=0
TEST_OUTPUT_DIR="test-output"

# Helper functions
log_success() {
    echo -e "${GREEN}✅ $1${NC}"
    ((TESTS_PASSED++))
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
    ((TESTS_FAILED++))
}

# Test phases
test_cli_commands() {
    log_info "Phase 1: Testing CLI Commands"
    
    # Test help command
    if ./bin/tool --help > /dev/null 2>&1; then
        log_success "CLI help command works"
    else
        log_error "CLI help command failed"
    fi
}

test_template_generation() {
    log_info "Phase 2: Testing Template Generation"
    
    for tier in basic intermediate advanced enterprise; do
        if ./bin/tool generate --name ${tier}-test --tier $tier --output ${TEST_OUTPUT_DIR}/${tier}-test > /dev/null 2>&1; then
            log_success "${tier} tier project generated"
        else
            log_error "${tier} tier project generation failed"
        fi
    done
}
```

### 2. Project Structure Validation
```bash
validate_project_structure() {
    local project_dir=$1
    local tier=$2
    
    # Common files all tiers should have
    local required_files=(
        "go.mod"
        "cmd/server/main.go"
        "internal/handlers/health.go"
        "internal/models/health.go"
        "README.md"
    )
    
    # Tier-specific files
    case $tier in
        "enterprise")
            required_files+=(
                "internal/security/mtls.go"
                "internal/security/rbac.go"
                "internal/compliance/audit.go"
            )
            ;;
        "advanced")
            required_files+=(
                "internal/observability/tracing.go"
                "internal/events/emitter.go"
            )
            ;;
    esac
    
    for file in "${required_files[@]}"; do
        if [[ -f "${project_dir}/${file}" ]]; then
            continue
        else
            return 1
        fi
    done
    
    return 0
}
```

### 3. Compilation Testing
```bash
test_compilation() {
    log_info "Phase 3: Testing Project Compilation"
    
    for tier in basic intermediate advanced enterprise; do
        local project_dir="${TEST_OUTPUT_DIR}/${tier}-test"
        
        if (cd "$project_dir" && go mod tidy > /dev/null 2>&1 && go build ./... > /dev/null 2>&1); then
            log_success "${tier} project compiles successfully"
        else
            log_error "${tier} project compilation failed"
            # Show compilation errors for debugging
            (cd "$project_dir" && go build ./... 2>&1 | head -10)
        fi
    done
}
```

### 4. Runtime Testing
```bash
test_runtime_functionality() {
    log_info "Phase 4: Testing Runtime Functionality"
    
    local project_dir="${TEST_OUTPUT_DIR}/enterprise-test"
    
    # Start server in background
    (cd "$project_dir" && go run cmd/server/main.go) &
    local server_pid=$!
    
    # Wait for server to start
    sleep 3
    
    # Test health endpoint
    if curl -s http://localhost:8080/health > /dev/null; then
        log_success "Enterprise server responds to health checks"
    else
        log_error "Enterprise server health check failed"
    fi
    
    # Test additional endpoints
    if curl -s http://localhost:8080/health/time > /dev/null; then
        log_success "Server time endpoint works"
    else
        log_error "Server time endpoint failed"
    fi
    
    # Cleanup
    kill $server_pid 2>/dev/null
}
```

### 5. Feature-Specific Testing
```bash
test_advanced_features() {
    log_info "Phase 5: Testing Advanced Features"
    
    # Test TypeScript generation
    if [[ -d "${TEST_OUTPUT_DIR}/advanced-test/client/typescript" ]]; then
        log_success "TypeScript client generated"
    else
        log_error "TypeScript client generation failed"
    fi
    
    # Test Kubernetes manifests
    if [[ -d "${TEST_OUTPUT_DIR}/enterprise-test/deployments/kubernetes" ]]; then
        log_success "Kubernetes manifests generated"
    else
        log_error "Kubernetes manifest generation failed"
    fi
    
    # Test BDD framework
    if (cd features/steps && go test -c > /dev/null 2>&1); then
        log_success "BDD framework compiles"
    else
        log_error "BDD framework compilation failed"
    fi
}
```

## Test Configuration

### 1. Environment Setup
```bash
# Test environment configuration
setup_test_environment() {
    # Clean previous test artifacts
    rm -rf "$TEST_OUTPUT_DIR"
    mkdir -p "$TEST_OUTPUT_DIR"
    
    # Ensure CLI is built
    if [[ ! -f "bin/tool" ]]; then
        go build -o bin/tool ./cmd/generator
    fi
    
    # Set test timeouts
    export TEST_TIMEOUT=30
    export SERVER_START_TIMEOUT=5
}
```

### 2. Cleanup and Reporting
```bash
cleanup_and_report() {
    # Kill any remaining processes
    pkill -f "go run cmd/server/main.go" 2>/dev/null || true
    
    # Generate test report
    echo ""
    echo "🏁 Integration Test Summary"
    echo "=========================="
    echo -e "Tests Passed: ${GREEN}$TESTS_PASSED${NC}"
    echo -e "Tests Failed: ${RED}$TESTS_FAILED${NC}"
    
    # Determine exit code
    if [[ $TESTS_FAILED -eq 0 ]]; then
        echo -e "\n${GREEN}🎉 All tests passed!${NC}"
        exit 0
    elif [[ $TESTS_FAILED -le 3 ]]; then
        echo -e "\n${YELLOW}⚠️  Most tests passed with minor issues.${NC}"
        exit 0
    else
        echo -e "\n${RED}❌ Multiple tests failed.${NC}"
        exit 1
    fi
}
```

## Advanced Testing Patterns

### 1. Parallel Testing
```bash
# Run tier tests in parallel
test_all_tiers_parallel() {
    local pids=()
    
    for tier in basic intermediate advanced enterprise; do
        test_single_tier "$tier" &
        pids+=($!)
    done
    
    # Wait for all tests to complete
    for pid in "${pids[@]}"; do
        wait $pid
    done
}
```

### 2. Performance Testing
```bash
# Measure generation performance
test_generation_performance() {
    local start_time=$(date +%s.%N)
    
    ./bin/tool generate --name perf-test --tier enterprise --output test-output/perf-test
    
    local end_time=$(date +%s.%N)
    local duration=$(echo "$end_time - $start_time" | bc)
    
    if (( $(echo "$duration < 5.0" | bc -l) )); then
        log_success "Generation completed in ${duration}s (under 5s threshold)"
    else
        log_error "Generation took ${duration}s (over 5s threshold)"
    fi
}
```

### 3. Error Condition Testing
```bash
# Test error handling
test_error_conditions() {
    # Test invalid tier
    if ./bin/tool generate --name test --tier invalid --output test-output 2>/dev/null; then
        log_error "Should fail with invalid tier"
    else
        log_success "Correctly handles invalid tier"
    fi
    
    # Test missing required flags
    if ./bin/tool generate --name test 2>/dev/null; then
        log_error "Should fail with missing required flags"
    else
        log_success "Correctly handles missing flags"
    fi
}
```

## Success Criteria

### Quantitative Metrics
- **Test Coverage**: All tiers and features tested
- **Pass Rate**: 95%+ test pass rate
- **Performance**: Generation under 5 seconds per tier
- **Compilation**: 100% compilation success rate

### Qualitative Metrics
- **Reliability**: Consistent test results across runs
- **Maintainability**: Easy to add new tests
- **Debuggability**: Clear error messages and logging
- **Automation**: Fully automated with minimal manual intervention

## Best Practices

1. **Isolation**: Each test should be independent
2. **Cleanup**: Always clean up test artifacts
3. **Timeouts**: Set appropriate timeouts for operations
4. **Logging**: Provide clear, actionable feedback
5. **Parallelization**: Run independent tests in parallel
6. **Error Handling**: Test both success and failure cases
7. **Performance**: Monitor and validate performance metrics

This comprehensive testing approach ensures template systems are reliable, performant, and maintainable.
