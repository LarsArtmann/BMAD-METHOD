#!/bin/bash

# Automated Template Validation Framework
# Validates all template tiers and ensures quality standards

set -e

# Ensure we're using bash 4+ for associative arrays
if [[ ${BASH_VERSION%%.*} -lt 4 ]]; then
    echo "Error: This script requires bash 4.0 or later for associative arrays"
    echo "Current version: $BASH_VERSION"
    exit 1
fi

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
VALIDATION_DIR="$PROJECT_ROOT/validation-temp"
CLI_BINARY="$PROJECT_ROOT/bin/template-health-endpoint"

# Test configurations
TEST_TIERS=("basic")
# Note: Only testing basic tier for now as other tiers need additional templates

# Validation metrics
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Logging functions
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
    ((PASSED_TESTS++))
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
    ((FAILED_TESTS++))
}

log_header() {
    echo -e "\n${BLUE}🔍 $1${NC}"
    echo "$(printf '=%.0s' {1..60})"
}

# Test functions
test_cli_exists() {
    ((TOTAL_TESTS++))
    if [[ -f "$CLI_BINARY" ]]; then
        log_success "CLI binary exists"
    else
        log_error "CLI binary not found at $CLI_BINARY"
        return 1
    fi
}

test_cli_help() {
    ((TOTAL_TESTS++))
    if "$CLI_BINARY" --help > /dev/null 2>&1; then
        log_success "CLI help command works"
    else
        log_error "CLI help command failed"
        return 1
    fi
}

test_typespec_validation() {
    ((TOTAL_TESTS++))
    log_info "Validating TypeSpec schemas..."

    if "$CLI_BINARY" validate --verbose > /dev/null 2>&1; then
        log_success "TypeSpec validation passed"
    else
        log_error "TypeSpec validation failed"
        return 1
    fi
}

test_tier_generation() {
    local tier=$1

    log_header "Testing $tier tier generation"

    local project_dir="$VALIDATION_DIR/test-$tier-service"
    local module_path="github.com/test/$tier-service"

    # Clean up previous test
    rm -rf "$project_dir"

    # Test generation
    ((TOTAL_TESTS++))
    log_info "Generating $tier tier project..."

    if "$CLI_BINARY" generate \
        --name "test-$tier-service" \
        --tier "$tier" \
        --module "$module_path" \
        --output "$project_dir" > /dev/null 2>&1; then
        log_success "$tier tier generation completed"
    else
        log_error "$tier tier generation failed"
        return 1
    fi

    # Test project structure
    test_project_structure "$project_dir" "$tier"

    # Test compilation
    test_project_compilation "$project_dir" "$tier"

    # Test endpoints
    test_project_endpoints "$project_dir" "$tier"

    # Test documentation
    test_project_documentation "$project_dir" "$tier"

    # Test TypeScript client (if enabled)
    test_typescript_client "$project_dir" "$tier"

    # Test Kubernetes manifests
    test_kubernetes_manifests "$project_dir" "$tier"
}

test_project_structure() {
    local project_dir=$1
    local tier=$2

    log_info "Validating project structure for $tier tier..."

    # Essential files that should exist in all tiers
    local essential_files=(
        "README.md"
        "go.mod"
        "Dockerfile"
        "Makefile"
        "cmd/server/main.go"
        "internal/handlers/health.go"
        "internal/models/health.go"
        "internal/server/server.go"
        "internal/config/config.go"
        "docs/API.md"
        "deployments/kubernetes/deployment.yaml"
        "deployments/kubernetes/service.yaml"
        "deployments/kubernetes/configmap.yaml"
    )

    for file in "${essential_files[@]}"; do
        ((TOTAL_TESTS++))
        if [[ -f "$project_dir/$file" ]]; then
            log_success "Essential file exists: $file"
        else
            log_error "Missing essential file: $file"
        fi
    done

    # TypeScript files (should exist if TypeScript is enabled)
    if [[ -d "$project_dir/client/typescript" ]]; then
        local ts_files=(
            "client/typescript/src/client.ts"
            "client/typescript/src/types.ts"
            "client/typescript/package.json"
            "client/typescript/tsconfig.json"
        )

        for file in "${ts_files[@]}"; do
            ((TOTAL_TESTS++))
            if [[ -f "$project_dir/$file" ]]; then
                log_success "TypeScript file exists: $file"
            else
                log_error "Missing TypeScript file: $file"
            fi
        done
    fi
}

test_project_compilation() {
    local project_dir=$1
    local tier=$2

    log_info "Testing compilation for $tier tier..."

    cd "$project_dir"

    # Test go mod tidy
    ((TOTAL_TESTS++))
    if go mod tidy > /dev/null 2>&1; then
        log_success "go mod tidy succeeded"
    else
        log_error "go mod tidy failed"
        cd "$PROJECT_ROOT"
        return 1
    fi

    # Test go build
    ((TOTAL_TESTS++))
    if go build -o "bin/test-$tier-service" cmd/server/main.go > /dev/null 2>&1; then
        log_success "Go build succeeded"
    else
        log_error "Go build failed"
        cd "$PROJECT_ROOT"
        return 1
    fi

    cd "$PROJECT_ROOT"
}

test_project_endpoints() {
    local project_dir=$1
    local tier=$2

    log_info "Testing endpoints for $tier tier..."

    cd "$project_dir"

    # Start the server in background
    ./bin/test-$tier-service &
    local server_pid=$!

    # Wait for server to start
    sleep 3

    # Test basic endpoints (all tiers)
    local basic_endpoints=(
        "/health"
        "/health/time"
        "/health/ready"
        "/health/live"
        "/health/startup"
    )

    for endpoint in "${basic_endpoints[@]}"; do
        ((TOTAL_TESTS++))
        if curl -f "http://localhost:8080$endpoint" > /dev/null 2>&1; then
            log_success "Endpoint working: $endpoint"
        else
            log_error "Endpoint failed: $endpoint"
        fi
    done

    # Test tier-specific endpoints
    case "$tier" in
        "intermediate"|"advanced"|"enterprise")
            ((TOTAL_TESTS++))
            if curl -f "http://localhost:8080/health/dependencies" > /dev/null 2>&1; then
                log_success "Dependencies endpoint working"
            else
                log_warning "Dependencies endpoint not responding (may be expected if no deps configured)"
            fi
            ;;
    esac

    case "$tier" in
        "advanced"|"enterprise")
            ((TOTAL_TESTS++))
            if curl -f "http://localhost:8080/health/metrics" > /dev/null 2>&1; then
                log_success "Metrics endpoint working"
            else
                log_warning "Metrics endpoint not responding"
            fi
            ;;
    esac

    # Stop the server
    kill $server_pid 2>/dev/null || true
    wait $server_pid 2>/dev/null || true

    cd "$PROJECT_ROOT"
}

test_project_documentation() {
    local project_dir=$1
    local tier=$2

    log_info "Testing documentation for $tier tier..."

    # Check README content
    ((TOTAL_TESTS++))
    if grep -q "/health/startup" "$project_dir/README.md"; then
        log_success "README includes startup endpoint"
    else
        log_error "README missing startup endpoint documentation"
    fi

    # Check API documentation
    ((TOTAL_TESTS++))
    if grep -q "GET /health/startup" "$project_dir/docs/API.md"; then
        log_success "API docs include startup endpoint"
    else
        log_error "API docs missing startup endpoint"
    fi

    # Check tier-specific documentation
    case "$tier" in
        "intermediate"|"advanced"|"enterprise")
            ((TOTAL_TESTS++))
            if grep -q "dependencies" "$project_dir/README.md"; then
                log_success "README mentions dependencies (appropriate for $tier tier)"
            else
                log_warning "README doesn't mention dependencies for $tier tier"
            fi
            ;;
    esac
}

test_typescript_client() {
    local project_dir=$1
    local tier=$2

    if [[ ! -d "$project_dir/client/typescript" ]]; then
        log_info "TypeScript client not enabled for $tier tier"
        return 0
    fi

    log_info "Testing TypeScript client for $tier tier..."

    # Check if client includes startup method
    ((TOTAL_TESTS++))
    if grep -q "checkStartup" "$project_dir/client/typescript/src/client.ts"; then
        log_success "TypeScript client includes checkStartup method"
    else
        log_error "TypeScript client missing checkStartup method"
    fi

    # Check if types are properly defined
    ((TOTAL_TESTS++))
    if grep -q "HealthReport" "$project_dir/client/typescript/src/types.ts"; then
        log_success "TypeScript types properly defined"
    else
        log_error "TypeScript types missing or incomplete"
    fi

    # Test TypeScript compilation (if Node.js is available)
    if command -v npm > /dev/null 2>&1; then
        cd "$project_dir/client/typescript"

        ((TOTAL_TESTS++))
        if npm install > /dev/null 2>&1 && npm run build > /dev/null 2>&1; then
            log_success "TypeScript client compiles successfully"
        else
            log_warning "TypeScript client compilation failed (may need dependencies)"
        fi

        cd "$PROJECT_ROOT"
    else
        log_info "Node.js not available, skipping TypeScript compilation test"
    fi
}

test_kubernetes_manifests() {
    local project_dir=$1
    local tier=$2

    log_info "Testing Kubernetes manifests for $tier tier..."

    # Check if startup probe is configured
    ((TOTAL_TESTS++))
    if grep -q "/health/startup" "$project_dir/deployments/kubernetes/deployment.yaml"; then
        log_success "Kubernetes deployment includes startup probe"
    else
        log_error "Kubernetes deployment missing startup probe"
    fi

    # Validate YAML syntax (if kubectl is available)
    if command -v kubectl > /dev/null 2>&1; then
        ((TOTAL_TESTS++))
        if kubectl apply --dry-run=client -f "$project_dir/deployments/kubernetes/" > /dev/null 2>&1; then
            log_success "Kubernetes manifests are valid"
        else
            log_error "Kubernetes manifests have syntax errors"
        fi
    else
        log_info "kubectl not available, skipping manifest validation"
    fi
}

test_dry_run_mode() {
    log_header "Testing dry run mode"

    ((TOTAL_TESTS++))
    log_info "Testing dry run generation..."

    if "$CLI_BINARY" generate \
        --name "dry-run-test" \
        --tier "basic" \
        --module "github.com/test/dry-run" \
        --dry-run > /dev/null 2>&1; then
        log_success "Dry run mode works"
    else
        log_error "Dry run mode failed"
    fi

    # Ensure no files were created
    ((TOTAL_TESTS++))
    if [[ ! -d "dry-run-test" ]]; then
        log_success "Dry run mode didn't create files"
    else
        log_error "Dry run mode created files when it shouldn't"
        rm -rf "dry-run-test"
    fi
}

test_error_handling() {
    log_header "Testing error handling"

    # Test invalid tier
    ((TOTAL_TESTS++))
    if ! "$CLI_BINARY" generate --name "test" --tier "invalid" --module "test" > /dev/null 2>&1; then
        log_success "Invalid tier properly rejected"
    else
        log_error "Invalid tier was accepted"
    fi

    # Test missing required flags
    ((TOTAL_TESTS++))
    if ! "$CLI_BINARY" generate --name "test" > /dev/null 2>&1; then
        log_success "Missing required flags properly rejected"
    else
        log_error "Missing required flags were accepted"
    fi
}

generate_validation_report() {
    log_header "Validation Report"

    echo "📊 Test Results:"
    echo "  Total Tests: $TOTAL_TESTS"
    echo "  Passed: $PASSED_TESTS"
    echo "  Failed: $FAILED_TESTS"
    echo "  Success Rate: $(( PASSED_TESTS * 100 / TOTAL_TESTS ))%"

    if [[ $FAILED_TESTS -eq 0 ]]; then
        echo -e "\n${GREEN}🎉 All validations passed! Templates are ready for production.${NC}"
        return 0
    else
        echo -e "\n${RED}❌ $FAILED_TESTS validation(s) failed. Please review and fix issues.${NC}"
        return 1
    fi
}

cleanup() {
    log_info "Cleaning up validation artifacts..."
    rm -rf "$VALIDATION_DIR"

    # Kill any remaining test servers
    pkill -f "test-.*-service" 2>/dev/null || true
}

main() {
    log_header "BMAD Method Template Validation Framework"

    # Ensure we're in the project root
    cd "$PROJECT_ROOT"

    # Create validation directory
    mkdir -p "$VALIDATION_DIR"

    # Trap cleanup on exit
    trap cleanup EXIT

    # Build CLI if it doesn't exist
    if [[ ! -f "$CLI_BINARY" ]]; then
        log_info "Building CLI tool..."
        go build -o "$CLI_BINARY" cmd/generator/main.go
    fi

    # Run validation tests
    test_cli_exists
    test_cli_help
    test_typespec_validation
    test_dry_run_mode
    test_error_handling

    # Test each tier
    for tier in "${TEST_TIERS[@]}"; do
        test_tier_generation "$tier"
    done

    # Generate final report
    generate_validation_report
}

# Run main function if script is executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
