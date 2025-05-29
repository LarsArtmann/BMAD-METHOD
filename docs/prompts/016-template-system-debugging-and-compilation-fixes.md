# Template System Debugging and Compilation Fixes

## Prompt Name: Template System Debugging and Compilation Fixes

## Context
You need to debug and fix compilation issues in a Go-based template generation system. The system generates multi-tier health endpoint services, but some generated projects have compilation errors due to unused imports in template code.

## Problem Description
Template systems that generate Go code often have issues with:
1. **Unused imports** in generated code
2. **Template variable substitution** not covering all necessary files
3. **Integration tests** expecting files that aren't generated
4. **Compilation errors** in generated projects

## Debugging Approach

### Step 1: Identify Compilation Errors
```bash
# Generate a project and test compilation
./bin/tool generate --name test-project --tier enterprise --output test-output
cd test-output && go mod tidy && go build ./...
```

### Step 2: Locate Template Issues
1. **Find the template source**: Look for inline templates in generator code
2. **Check import statements**: Identify unused imports in template definitions
3. **Validate variable substitution**: Ensure all template variables are properly replaced

### Step 3: Fix Template Code
```go
// Example: Remove unused imports from template
"go-security-rbac": `package security

import (
    // Remove unused imports like "fmt" or "encoding/json" if not used
    "net/http"
    "strings"
)
```

### Step 4: Update Integration Tests
```bash
# Check what files are actually generated
find generated-project -type f | sort

# Update test expectations to match reality
if [[ -f "project/internal/security/rbac.go" && \
      -f "project/internal/security/mtls.go" ]]; then
    log_success "Enterprise structure correct"
fi
```

## Common Issues and Solutions

### Issue 1: Unused Import "encoding/json"
**Problem**: Template includes `"encoding/json"` import but only uses struct tags
**Solution**: Remove the import - struct tags don't require importing the package

### Issue 2: Unused Import "fmt"
**Problem**: Template includes `"fmt"` import but doesn't call fmt functions
**Solution**: Remove the import or add actual fmt usage

### Issue 3: Unused Import "context"
**Problem**: Template includes `"context"` import but only passes context, doesn't create it
**Solution**: Remove the import - `r.Context()` doesn't require importing context package

### Issue 4: Integration Test File Expectations
**Problem**: Test expects files that aren't generated (e.g., `configs/development.yaml`)
**Solution**: Update test to check for actually generated files

## Template Processing Best Practices

### 1. Import Management
```go
// Only include imports that are actually used in the template
import (
    "net/http"        // Used for http.Handler, http.Request
    "strings"         // Used for strings.HasPrefix()
    // Don't include "fmt" unless calling fmt.Printf, fmt.Sprintf, etc.
    // Don't include "encoding/json" unless calling json.Marshal, json.Unmarshal
)
```

### 2. Template Variable Scope
```go
// Process all relevant file types
func needsTemplateProcessing(filePath string) bool {
    ext := filepath.Ext(filePath)
    processableExts := []string{
        ".go", ".js", ".ts", ".py",           // Source code
        ".yaml", ".yml", ".json", ".toml",    // Configuration
        ".sh", ".bat", ".ps1",                // Scripts
        ".md", ".txt",                        // Documentation
    }
    return contains(processableExts, ext)
}
```

### 3. Compilation Validation
```bash
# Always test generated projects compile
generate_and_test() {
    local tier=$1
    ./bin/tool generate --name test-${tier} --tier ${tier} --output test-output/${tier}
    cd test-output/${tier}
    go mod tidy
    go build ./...
    cd ../..
}
```

## Testing Strategy

### 1. Unit Tests for Templates
- Test template parsing and variable substitution
- Validate generated code syntax
- Check import statements are necessary

### 2. Integration Tests
- Generate complete projects
- Compile generated projects
- Run generated applications
- Validate expected functionality

### 3. Regression Tests
- Test all tiers after changes
- Ensure no existing functionality breaks
- Validate performance doesn't degrade

## Debugging Tools

### 1. Go Build Analysis
```bash
# Check for unused imports
go build -v ./...

# Get detailed error information
go build -x ./...
```

### 2. Template Validation
```bash
# Validate template syntax
go run cmd/generator/main.go template validate

# Test template processing
go run cmd/generator/main.go generate --dry-run
```

### 3. Generated Code Analysis
```bash
# Check generated file structure
find generated-project -name "*.go" -exec go fmt -d {} \;

# Validate imports
find generated-project -name "*.go" -exec goimports -d {} \;
```

## Success Criteria
1. **All generated projects compile without warnings**
2. **No unused imports in generated code**
3. **Integration tests pass completely**
4. **Generated applications run correctly**

## Common Pitfalls
1. **Assuming struct tags require imports** - they don't
2. **Including imports "just in case"** - only include what's used
3. **Not testing all tiers** - changes can affect multiple tiers
4. **Hardcoded test expectations** - tests should match actual generation

This approach ensures clean, compilable generated code and reliable template systems.
