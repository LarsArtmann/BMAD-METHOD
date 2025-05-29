#!/bin/bash

# Simple Template Validation Script
# Compatible with older bash versions

set -e

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
    PASSED_TESTS=$((PASSED_TESTS + 1))
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
    FAILED_TESTS=$((FAILED_TESTS + 1))
}

log_header() {
    echo -e "\n${BLUE}🔍 $1${NC}"
    echo "$(printf '=%.0s' {1..60})"
}

# Test functions
test_cli_exists() {
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    if [ -f "$CLI_BINARY" ]; then
        log_success "CLI binary exists"
    else
        log_error "CLI binary not found at $CLI_BINARY"
        return 1
    fi
}

test_basic_tier() {
    log_header "Testing Basic Tier Generation"
    
    local project_dir="$VALIDATION_DIR/test-basic-service"
    
    # Clean up previous test
    rm -rf "$project_dir"
    
    # Test generation
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    log_info "Generating basic tier project..."
    
    if "$CLI_BINARY" generate \
        --name "test-basic-service" \
        --tier "basic" \
        --module "github.com/test/basic-service" \
        --output "$project_dir" > /dev/null 2>&1; then
        log_success "Basic tier generation completed"
    else
        log_error "Basic tier generation failed"
        return 1
    fi
    
    # Test essential files
    local essential_files="README.md go.mod cmd/server/main.go internal/handlers/health.go"
    
    for file in $essential_files; do
        TOTAL_TESTS=$((TOTAL_TESTS + 1))
        if [ -f "$project_dir/$file" ]; then
            log_success "Essential file exists: $file"
        else
            log_error "Missing essential file: $file"
        fi
    done
    
    # Test startup endpoint in health handler
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    if grep -q "StartupCheck" "$project_dir/internal/handlers/health.go"; then
        log_success "Health handler includes StartupCheck method"
    else
        log_error "Health handler missing StartupCheck method"
    fi
    
    # Test startup route in server
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    if grep -q "/startup" "$project_dir/internal/server/server.go"; then
        log_success "Server includes startup route"
    else
        log_error "Server missing startup route"
    fi
    
    # Test documentation includes startup endpoint
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    if grep -q "/health/startup" "$project_dir/README.md"; then
        log_success "README includes startup endpoint"
    else
        log_error "README missing startup endpoint"
    fi
    
    # Test compilation
    cd "$project_dir"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    if go mod tidy > /dev/null 2>&1; then
        log_success "go mod tidy succeeded"
    else
        log_error "go mod tidy failed"
        cd "$PROJECT_ROOT"
        return 1
    fi
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    if go build -o "bin/test-basic-service" cmd/server/main.go > /dev/null 2>&1; then
        log_success "Go build succeeded"
    else
        log_error "Go build failed"
        cd "$PROJECT_ROOT"
        return 1
    fi
    
    # Test endpoints
    log_info "Testing endpoints..."
    
    # Start the server in background
    ./bin/test-basic-service &
    local server_pid=$!
    
    # Wait for server to start
    sleep 3
    
    # Test all health endpoints
    local endpoints="/health /health/time /health/ready /health/live /health/startup"
    
    for endpoint in $endpoints; do
        TOTAL_TESTS=$((TOTAL_TESTS + 1))
        if curl -f "http://localhost:8080$endpoint" > /dev/null 2>&1; then
            log_success "Endpoint working: $endpoint"
        else
            log_error "Endpoint failed: $endpoint"
        fi
    done
    
    # Stop the server
    kill $server_pid 2>/dev/null || true
    wait $server_pid 2>/dev/null || true
    
    cd "$PROJECT_ROOT"
}

test_typespec_validation() {
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    log_info "Validating TypeSpec schemas..."
    
    if "$CLI_BINARY" validate --verbose > /dev/null 2>&1; then
        log_success "TypeSpec validation passed"
    else
        log_error "TypeSpec validation failed"
        return 1
    fi
}

test_dry_run() {
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    log_info "Testing dry run mode..."
    
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
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    if [ ! -d "dry-run-test" ]; then
        log_success "Dry run mode didn't create files"
    else
        log_error "Dry run mode created files when it shouldn't"
        rm -rf "dry-run-test"
    fi
}

generate_report() {
    log_header "Validation Report"
    
    echo "📊 Test Results:"
    echo "  Total Tests: $TOTAL_TESTS"
    echo "  Passed: $PASSED_TESTS"
    echo "  Failed: $FAILED_TESTS"
    
    if [ $TOTAL_TESTS -gt 0 ]; then
        local success_rate=$((PASSED_TESTS * 100 / TOTAL_TESTS))
        echo "  Success Rate: ${success_rate}%"
    fi
    
    if [ $FAILED_TESTS -eq 0 ]; then
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
    log_header "BMAD Method Template Validation"
    
    # Ensure we're in the project root
    cd "$PROJECT_ROOT"
    
    # Create validation directory
    mkdir -p "$VALIDATION_DIR"
    
    # Trap cleanup on exit
    trap cleanup EXIT
    
    # Build CLI if it doesn't exist
    if [ ! -f "$CLI_BINARY" ]; then
        log_info "Building CLI tool..."
        go build -o "$CLI_BINARY" cmd/generator/main.go
    fi
    
    # Run validation tests
    test_cli_exists
    test_typespec_validation
    test_dry_run
    test_basic_tier
    
    # Generate final report
    generate_report
}

# Run main function if script is executed directly
if [ "${BASH_SOURCE[0]}" = "${0}" ]; then
    main "$@"
fi
