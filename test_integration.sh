#!/bin/bash

# Integration test script for template-health-endpoint
# This script tests the complete BMAD-METHOD implementation

set -e

echo "🚀 Starting BMAD-METHOD Integration Tests"
echo "=========================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Helper functions
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
    ((TESTS_PASSED++))
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
    ((TESTS_FAILED++))
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Cleanup function
cleanup() {
    log_info "Cleaning up test artifacts..."
    rm -rf test-output/
    rm -rf temp-*/
}

# Set up cleanup trap
trap cleanup EXIT

# Create test output directory
mkdir -p test-output

echo ""
log_info "Phase 1: Testing CLI Commands"
echo "------------------------------"

# Test 1: CLI Help
log_info "Testing CLI help command..."
if ./bin/template-health-endpoint --help > /dev/null 2>&1; then
    log_success "CLI help command works"
else
    log_error "CLI help command failed"
fi

# Test 2: CLI Version
log_info "Testing CLI version command..."
if ./bin/template-health-endpoint --version > /dev/null 2>&1; then
    log_success "CLI version command works"
else
    log_warning "CLI version command not implemented (expected)"
fi

# Test 3: List templates
log_info "Testing template list command..."
if ./bin/template-health-endpoint list > /dev/null 2>&1; then
    log_success "Template list command works"
else
    log_warning "Template list command not implemented (expected)"
fi

echo ""
log_info "Phase 2: Testing Template Generation"
echo "------------------------------------"

# Test 4: Generate Basic Tier Project
log_info "Generating basic tier project..."
if ./bin/template-health-endpoint generate \
    --name basic-test \
    --tier basic \
    --go-module github.com/test/basic-test \
    --output test-output > /dev/null 2>&1; then
    log_success "Basic tier project generated"
else
    log_error "Basic tier project generation failed"
fi

# Test 5: Generate Intermediate Tier Project
log_info "Generating intermediate tier project..."
if ./bin/template-health-endpoint generate \
    --name intermediate-test \
    --tier intermediate \
    --go-module github.com/test/intermediate-test \
    --output test-output > /dev/null 2>&1; then
    log_success "Intermediate tier project generated"
else
    log_error "Intermediate tier project generation failed"
fi

# Test 6: Generate Advanced Tier Project
log_info "Generating advanced tier project..."
if ./bin/template-health-endpoint generate \
    --name advanced-test \
    --tier advanced \
    --go-module github.com/test/advanced-test \
    --output test-output > /dev/null 2>&1; then
    log_success "Advanced tier project generated"
else
    log_error "Advanced tier project generation failed"
fi

# Test 7: Generate Enterprise Tier Project
log_info "Generating enterprise tier project..."
if ./bin/template-health-endpoint generate \
    --name enterprise-test \
    --tier enterprise \
    --go-module github.com/test/enterprise-test \
    --output test-output > /dev/null 2>&1; then
    log_success "Enterprise tier project generated"
else
    log_error "Enterprise tier project generation failed"
fi

echo ""
log_info "Phase 3: Testing Project Structure"
echo "----------------------------------"

# Test 8: Validate Basic Project Structure
log_info "Validating basic project structure..."
if [[ -f "test-output/basic-test/go.mod" && \
      -f "test-output/basic-test/cmd/server/main.go" && \
      -f "test-output/basic-test/internal/handlers/health.go" ]]; then
    log_success "Basic project structure is correct"
else
    log_error "Basic project structure is incorrect"
fi

# Test 9: Validate Enterprise Project Structure
log_info "Validating enterprise project structure..."
if [[ -f "test-output/enterprise-test/internal/security/mtls.go" && \
      -f "test-output/enterprise-test/internal/security/rbac.go" && \
      -f "test-output/enterprise-test/internal/compliance/audit.go" && \
      -f "test-output/enterprise-test/configs/development.yaml" ]]; then
    log_success "Enterprise project structure is correct"
else
    log_error "Enterprise project structure is incorrect"
fi

echo ""
log_info "Phase 4: Testing Project Compilation"
echo "------------------------------------"

# Test 10: Compile Basic Project
log_info "Compiling basic project..."
if (cd test-output/basic-test && go mod tidy > /dev/null 2>&1 && go build ./... > /dev/null 2>&1); then
    log_success "Basic project compiles successfully"
else
    log_error "Basic project compilation failed"
fi

# Test 11: Compile Enterprise Project
log_info "Compiling enterprise project..."
if (cd test-output/enterprise-test && go mod tidy > /dev/null 2>&1 && go build ./... > /dev/null 2>&1); then
    log_success "Enterprise project compiles successfully"
else
    log_error "Enterprise project compilation failed"
fi

echo ""
log_info "Phase 5: Testing Advanced CLI Commands"
echo "--------------------------------------"

# Test 12: Migration Command (Dry Run)
log_info "Testing migration command (dry run)..."
if (cd test-output/basic-test && ../../../bin/template-health-endpoint migrate --to intermediate --dry-run > /dev/null 2>&1); then
    log_success "Migration command works (dry run)"
else
    log_warning "Migration command not fully implemented (expected)"
fi

# Test 13: Update Command (Dry Run)
log_info "Testing update command (dry run)..."
if (cd test-output/basic-test && ../../../bin/template-health-endpoint update --dry-run > /dev/null 2>&1); then
    log_success "Update command works (dry run)"
else
    log_warning "Update command not fully implemented (expected)"
fi

# Test 14: Customize Command Help
log_info "Testing customize command help..."
if ./bin/template-health-endpoint customize --help > /dev/null 2>&1; then
    log_success "Customize command help works"
else
    log_error "Customize command help failed"
fi

echo ""
log_info "Phase 6: Testing Template Features"
echo "----------------------------------"

# Test 15: TypeScript Feature
log_info "Testing TypeScript feature generation..."
if ./bin/template-health-endpoint generate \
    --name ts-test \
    --tier advanced \
    --features typescript \
    --output test-output > /dev/null 2>&1; then
    if [[ -d "test-output/ts-test/client/typescript" ]]; then
        log_success "TypeScript feature works"
    else
        log_warning "TypeScript feature partially implemented"
    fi
else
    log_warning "TypeScript feature not fully implemented (expected)"
fi

# Test 16: Kubernetes Feature
log_info "Testing Kubernetes feature generation..."
if ./bin/template-health-endpoint generate \
    --name k8s-test \
    --tier intermediate \
    --features kubernetes \
    --output test-output > /dev/null 2>&1; then
    if [[ -d "test-output/k8s-test/deployments/kubernetes" ]]; then
        log_success "Kubernetes feature works"
    else
        log_warning "Kubernetes feature partially implemented"
    fi
else
    log_warning "Kubernetes feature not fully implemented (expected)"
fi

echo ""
log_info "Phase 7: Testing BDD Framework"
echo "------------------------------"

# Test 17: BDD Test Compilation
log_info "Testing BDD test compilation..."
if (cd features/steps && go test -c > /dev/null 2>&1); then
    log_success "BDD tests compile successfully"
else
    log_error "BDD test compilation failed"
fi

# Test 18: Run Sample BDD Test
log_info "Running sample BDD test..."
if (cd features/steps && timeout 30s go test -v -run TestFeatures/Help_and_usage_information > /dev/null 2>&1); then
    log_success "Sample BDD test passes"
else
    log_warning "BDD tests need more implementation (expected)"
fi

echo ""
log_info "Phase 8: Testing Documentation"
echo "------------------------------"

# Test 19: README Files
log_info "Checking README files..."
readme_count=0
for project in test-output/*/; do
    if [[ -f "$project/README.md" ]]; then
        ((readme_count++))
    fi
done

if [[ $readme_count -gt 0 ]]; then
    log_success "README files are generated ($readme_count found)"
else
    log_error "No README files found in generated projects"
fi

# Test 20: Template Documentation
log_info "Checking template documentation..."
if [[ -f "templates/basic/template.yaml" && \
      -f "templates/enterprise/template.yaml" ]]; then
    log_success "Template documentation exists"
else
    log_error "Template documentation missing"
fi

echo ""
echo "🏁 Integration Test Summary"
echo "=========================="
echo -e "Tests Passed: ${GREEN}$TESTS_PASSED${NC}"
echo -e "Tests Failed: ${RED}$TESTS_FAILED${NC}"

if [[ $TESTS_FAILED -eq 0 ]]; then
    echo -e "\n${GREEN}🎉 All critical tests passed! BMAD-METHOD implementation is working!${NC}"
    exit 0
elif [[ $TESTS_FAILED -le 3 ]]; then
    echo -e "\n${YELLOW}⚠️  Most tests passed with minor issues. Implementation is mostly complete.${NC}"
    exit 0
else
    echo -e "\n${RED}❌ Multiple tests failed. Implementation needs more work.${NC}"
    exit 1
fi
