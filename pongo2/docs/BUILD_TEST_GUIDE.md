# Pongo2 Prompt Template Activity - Build & Test Guide

## Prerequisites

- Go 1.21 or later
- Flogo CLI (for integration testing)
- Git (for version control)

## Building the Activity

### 1. Clone/Navigate to the Activity Directory
```bash
cd /path/to/pongo2
```

### 2. Initialize Go Modules (if not already done)
```bash
go mod init pongo2
```

### 3. Install Dependencies
```bash
go mod tidy
```
This will download:
- `github.com/project-flogo/core` - Flogo core framework
- `github.com/flosch/pongo2/v6` - Jinja2-like template engine

### 4. Build the Activity
```bash
go build .
```

### 5. Build as Flogo Plugin (Optional)
For use as a plugin in Flogo applications:
```bash
go build -buildmode=plugin -o pongo2.so .
```

## Testing the Activity

### 1. Run All Tests
```bash
go test -v
```

### 2. Run Specific Test
```bash
go test -v -run TestJinja2PromptActivity
go test -v -run TestDataAnalysisTemplate
```

### 3. Run Tests with Coverage
```bash
go test -v -cover
go test -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### 4. Benchmark Tests (if added)
```bash
go test -bench=.
```

## Project Structure

The activity is organized into separate files for better maintainability:

- **`activity.go`** - Main activity logic and execution (`Eval` method)
- **`metadata.go`** - Input/Output struct definitions with `ToMap`/`FromMap` methods
- **`utils.go`** - Utility functions for template variable extraction and schema generation
- **`activity_test.go`** - Comprehensive test suite
- **`user_template_test.go`** - User-provided template tests
- **`schema_provider.go`** - Schema provider for dynamic input generation
- **`descriptor.json`** - Activity metadata and configuration

## Test Structure

The test suite includes:

1. **TestJinja2PromptActivity** - Basic template rendering
2. **TestJinja2PromptActivityWithVariablesMap** - Complex variable mapping
3. **TestJinja2PromptActivityComplexTemplate** - AI prompt templates
4. **TestDataAnalysisTemplate** - Advanced data analysis template
5. **TestTemplateVariablesInput** - Schema-based input testing
6. **TestCombinedInputApproaches** - Multiple input method testing
7. **TestJSONSchemaGeneration** - Utility function testing
8. **TestUserProvidedTemplate** - Real-world template testing

## Debugging

### Enable Verbose Logging
```bash
go test -v -args -test.v=true
```

### Debug Template Parsing
Add debug prints in the activity code:
```go
fmt.Printf("Template: %s\n", templateStr)
fmt.Printf("Variables: %+v\n", templateVars)
fmt.Printf("Detected variables: %v\n", extractTemplateVariables(templateStr))
```

### Check Template Syntax
Use pongo2 online validator or create a simple test:
```go
template, err := pongo2.FromString(templateString)
if err != nil {
    t.Fatalf("Template parsing error: %v", err)
}
```

## Common Build Issues

### Issue: Module not found
**Error**: `could not import github.com/project-flogo/core/activity`
**Solution**: Run `go mod tidy` to download dependencies

### Issue: Pongo2 import error  
**Error**: `could not import github.com/flosch/pongo2/v6`
**Solution**: Check internet connection and run `go mod download`

### Issue: Test failures
**Error**: Interface implementation errors
**Solution**: Ensure mock implementations have all required methods

### Issue: Metadata struct not found
**Error**: `undefined: Settings` or similar struct errors
**Solution**: Ensure `metadata.go` is in the same package and properly imports coerce

### Issue: Utility function not found
**Error**: `undefined: extractTemplateVariables`
**Solution**: Ensure `utils.go` is in the same package and functions are properly exported

## Development Workflow

1. **Make changes** to relevant files:
   - `activity.go` - Main execution logic
   - `metadata.go` - Input/Output definitions
   - `utils.go` - Utility functions
   - `descriptor.json` - Activity configuration
2. **Run tests** with `go test -v`
3. **Test specific functionality**:
   - Template parsing: `go test -v -run TestJSONSchemaGeneration`
   - Input handling: `go test -v -run TestTemplateVariablesInput`
   - Complex templates: `go test -v -run TestDataAnalysisTemplate`
4. **Fix any issues** revealed by tests
5. **Build** with `go build .`
6. **Integration test** in a Flogo application
7. **Commit changes** to version control

## Architecture & Design

The activity follows standard Flogo patterns with clean separation of concerns:

### File Organization
- **Metadata Layer** (`metadata.go`): Defines Input/Output structs with proper serialization
- **Business Logic** (`activity.go`): Core execution logic and Flogo integration
- **Utility Layer** (`utils.go`): Template parsing and schema generation functions
- **Schema Provider** (`schema_provider.go`): Dynamic input field generation
- **Tests** (`*_test.go`): Comprehensive test coverage (51.6% coverage)

### Benefits of This Architecture
- **Maintainability**: Easy to modify each layer independently
- **Testability**: Utility functions can be tested in isolation
- **Reusability**: Utility functions can be used by other components
- **Standards Compliance**: Follows Flogo activity best practices

## Continuous Integration

Example GitHub Actions workflow:
```yaml
name: Test Pongo2 Prompt Activity
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - uses: actions/setup-go@v2
      with:
        go-version: 1.21
    - name: Download dependencies
      run: go mod download
    - name: Run tests
      run: go test -v ./...
    - name: Test coverage
      run: go test -v -cover
    - name: Test schema utilities
      run: cd utils && ./generate_flogo_params.sh "Hello {{ name }}"
    - name: Build
      run: go build .
```

## Performance Testing

### Load Testing Template Rendering
```go
func BenchmarkTemplateRendering(b *testing.B) {
    // Setup activity and context
    for i := 0; i < b.N; i++ {
        act.Eval(ctx)
    }
}
```

### Memory Profiling
```bash
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof
```

## Schema Generation Utilities

The activity includes utilities for generating JSON schemas from templates:

### 1. Test Schema Generation
```bash
cd utils
go run schema_generator.go "Hello {{ name }}, you are {{ age }} years old"
```

### 2. Use the Shell Script (with auto Go installation)
```bash
cd utils
./generate_flogo_params.sh "Hello {{ name }}, you are {{ age }} years old"
```

### 3. Testing Utility Functions
```bash
# Test template variable extraction
go test -v -run TestJSONSchemaGeneration

# Test utility functions directly
go test -v utils.go activity_test.go
```

## Distribution

### Package for Distribution
```bash
# Create distribution package (includes all source files)
tar -czf pongo2-v1.0.0.tar.gz *.go descriptor.json go.mod go.sum *.md utils/

# Create zip for Windows
zip -r pongo2-v1.0.0.zip *.go descriptor.json go.mod go.sum *.md utils/
```

### Publish as Go Module (Optional)
1. Tag the release: `git tag v1.0.0`
2. Push tags: `git push --tags`
3. Publish to Go module registry