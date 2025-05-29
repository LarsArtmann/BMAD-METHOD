# Template Variable Processing and File Type Handling

## Prompt Name
**Template Variable Processing and File Type Handling**

## Context
When building template systems that generate code projects, proper template variable processing across all file types is crucial. Initial implementations often miss file types, leading to incomplete variable substitution and broken generated projects.

## Objective
Implement comprehensive template variable processing that handles all relevant file types and ensures complete variable substitution across the entire generated project.

## Key Requirements

### 1. Comprehensive File Type Coverage
Process template variables in ALL relevant file types:
- **Template files**: `.tmpl` extensions
- **Source code**: `.go`, `.js`, `.ts`, `.py`, `.java`, `.cs`
- **Configuration**: `.yaml`, `.yml`, `.json`, `.toml`, `.ini`
- **Scripts**: `.sh`, `.bat`, `.ps1`
- **Documentation**: `.md`, `.txt`, `.rst`
- **Build files**: `Makefile`, `Dockerfile`, `docker-compose.yml`
- **Package files**: `go.mod`, `package.json`, `requirements.txt`, `Cargo.toml`

### 2. Template Variable Standards
- **Naming Convention**: Use descriptive, hierarchical names
  ```
  {{.Config.Name}}
  {{.Config.GoModule}}
  {{.Config.Description}}
  {{.Features.Kubernetes}}
  {{.Version}}
  {{.Timestamp}}
  ```
- **Grouping**: Organize variables logically under namespaces
- **Defaults**: Provide sensible defaults for optional variables
- **Validation**: Ensure all required variables are provided

### 3. Processing Logic Implementation
```go
func needsTemplateProcessing(filePath string) bool {
    // Check file extension
    ext := filepath.Ext(filePath)
    if ext == ".tmpl" {
        return true
    }
    
    // Check specific file types
    processableExts := []string{
        ".go", ".js", ".ts", ".py", ".java", ".cs",
        ".yaml", ".yml", ".json", ".toml", ".ini",
        ".sh", ".bat", ".ps1",
        ".md", ".txt", ".rst",
    }
    
    for _, processableExt := range processableExts {
        if ext == processableExt {
            return true
        }
    }
    
    // Check specific filenames
    baseName := filepath.Base(filePath)
    processableFiles := []string{
        "go.mod", "package.json", "requirements.txt",
        "Makefile", "Dockerfile", "docker-compose.yml",
    }
    
    for _, processableFile := range processableFiles {
        if baseName == processableFile {
            return true
        }
    }
    
    return false
}
```

### 4. Template Processing Pipeline
1. **Load Template**: Read template file content
2. **Parse Variables**: Extract template variables and validate
3. **Apply Context**: Substitute variables with actual values
4. **Handle Extensions**: Remove `.tmpl` extension from output files
5. **Validate Output**: Ensure generated files are syntactically correct

### 5. Error Handling
- **Missing Variables**: Clear error messages for undefined variables
- **Syntax Errors**: Validate template syntax before processing
- **File Permissions**: Handle read/write permission issues gracefully
- **Path Issues**: Validate input and output paths

## Implementation Example

```go
func processTemplateFile(inputPath, outputPath string, context map[string]interface{}) error {
    // Read input file
    content, err := os.ReadFile(inputPath)
    if err != nil {
        return fmt.Errorf("failed to read template file %s: %w", inputPath, err)
    }
    
    // Create output directory if needed
    outputDir := filepath.Dir(outputPath)
    if err := os.MkdirAll(outputDir, 0755); err != nil {
        return fmt.Errorf("failed to create output directory: %w", err)
    }
    
    // Remove .tmpl extension from output path
    if filepath.Ext(outputPath) == ".tmpl" {
        outputPath = outputPath[:len(outputPath)-5]
    }
    
    // Check if file needs template processing
    if needsTemplateProcessing(inputPath) {
        // Process as template
        tmpl, err := template.New(filepath.Base(inputPath)).Parse(string(content))
        if err != nil {
            return fmt.Errorf("failed to parse template %s: %w", inputPath, err)
        }
        
        // Create output file
        outputFile, err := os.Create(outputPath)
        if err != nil {
            return fmt.Errorf("failed to create output file: %w", err)
        }
        defer outputFile.Close()
        
        // Execute template
        if err := tmpl.Execute(outputFile, context); err != nil {
            return fmt.Errorf("failed to execute template: %w", err)
        }
    } else {
        // Copy file as-is
        if err := os.WriteFile(outputPath, content, 0644); err != nil {
            return fmt.Errorf("failed to copy file: %w", err)
        }
    }
    
    return nil
}
```

## Testing Strategy

### 1. Variable Substitution Testing
- Test all file types with template variables
- Verify variables are substituted correctly
- Test edge cases (empty values, special characters)

### 2. File Type Coverage Testing
- Generate projects and verify all files are processed
- Check that no hardcoded values remain
- Validate generated projects compile/run correctly

### 3. Error Handling Testing
- Test with missing variables
- Test with invalid template syntax
- Test with permission issues

## Common Pitfalls

### 1. Incomplete File Type Coverage
**Problem**: Only processing `.tmpl` files and missing source code files
**Solution**: Comprehensive file type checking as shown above

### 2. Hardcoded Values in Generated Projects
**Problem**: Template variables not substituted in all files
**Solution**: Systematic testing of all generated files

### 3. Template Syntax Errors
**Problem**: Invalid Go template syntax causing processing failures
**Solution**: Template validation before processing

### 4. Path Handling Issues
**Problem**: Incorrect handling of file paths and extensions
**Solution**: Robust path manipulation and validation

## Success Criteria
- [ ] All relevant file types are processed for template variables
- [ ] Generated projects have no hardcoded template values
- [ ] Generated projects compile and run successfully
- [ ] Template processing handles errors gracefully
- [ ] Performance is acceptable for large projects

## Related Patterns
- Template engine design
- Code generation systems
- File processing pipelines
- Configuration management
- Build system integration
