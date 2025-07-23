# Creating Custom TIBCO Flogo Activities: A Go Learning Tutorial

This tutorial teaches you how to create custom TIBCO Flogo activities using Go (Golang), using our Prometheus Metrics activity as a real-world example. Perfect for developers learning Go or wanting to extend Flogo with custom functionality.

## 🎯 What You'll Learn

- **Go Language Fundamentals**: Structs, interfaces, methods, and packages
- **TIBCO Flogo Framework**: Activity lifecycle, inputs/outputs, and integration
- **JSON Processing**: Parsing, transforming, and validating JSON data
- **String Manipulation**: Building formatted output and text processing
- **Error Handling**: Go's error patterns and best practices
- **Testing**: Writing comprehensive unit tests for activities
- **Module Management**: Go modules, dependencies, and versioning

---

## 📚 Go Language Basics

### 1. Understanding Go Structs

In Go, a `struct` is a collection of fields that groups related data together. Our activity uses structs to define its structure and behavior.

```go
// A struct defines a custom type with named fields
type PrometheusActivity struct {
    // Fields can be basic types
    metricType string
    metricName string
    
    // Or more complex types
    settings map[string]interface{}
}

// You can also embed other types
type ActivitySettings struct {
    MetricType    string `json:"metricType"`
    MetricName    string `json:"metricName"`
    IncludeHelp   bool   `json:"includeHelp"`
    IncludeType   bool   `json:"includeType"`
    Timestamp     bool   `json:"timestamp"`
}
```

**Key Learning Points**:
- Structs group related data
- Field names starting with capital letters are **exported** (public)
- Field names starting with lowercase are **unexported** (private)
- Struct tags (like `json:"metricType"`) control JSON marshaling/unmarshaling

### 2. Interfaces in Go

Interfaces define behavior contracts. Flogo activities must implement the `activity.Activity` interface.

```go
// This is how Flogo defines what an activity should do
type Activity interface {
    Eval(ctx activity.Context) (done bool, err error)
}

// Our activity implements this interface
func (a *PrometheusActivity) Eval(ctx activity.Context) (done bool, err error) {
    // Activity logic goes here
    return true, nil
}
```

**Key Learning Points**:
- Interfaces define what methods a type must have
- Go uses "duck typing" - if it has the right methods, it implements the interface
- Small interfaces are preferred in Go

### 3. Methods and Receivers

Methods in Go are functions with a receiver - they belong to a specific type.

```go
// Method with pointer receiver (can modify the struct)
func (a *PrometheusActivity) processData(data map[string]interface{}) string {
    // Process the data and return result
    return "processed"
}

// Method with value receiver (read-only access)
func (s ActivitySettings) validate() error {
    if s.MetricName == "" {
        return errors.New("metric name cannot be empty")
    }
    return nil
}
```

**Key Learning Points**:
- Use pointer receivers (`*Type`) when you need to modify the struct
- Use value receivers (`Type`) when you only need to read data
- Method names starting with capital letters are exported

### 4. Error Handling in Go

Go uses explicit error handling - functions return errors as values.

```go
func processField(fieldName string, value interface{}) (float64, error) {
    switch v := value.(type) {
    case float64:
        return v, nil
    case int:
        return float64(v), nil
    case string:
        if parsed, err := strconv.ParseFloat(v, 64); err == nil {
            return parsed, nil
        }
        return 0, fmt.Errorf("cannot parse '%s' as number", v)
    default:
        return 0, fmt.Errorf("field '%s' is not a numeric type", fieldName)
    }
}

// Using the function
value, err := processField("temperature", "23.5")
if err != nil {
    log.Printf("Error processing field: %v", err)
    return err
}
// Use value here...
```

**Key Learning Points**:
- Always check errors: `if err != nil`
- Use `fmt.Errorf()` to create formatted error messages
- Return errors as the last return value

### 5. Type Assertions and Type Switches

Go provides ways to work with `interface{}` (empty interface) values.

```go
func examineValue(value interface{}) {
    // Type assertion - check if value is a specific type
    if str, ok := value.(string); ok {
        fmt.Printf("It's a string: %s\n", str)
    }
    
    // Type switch - handle multiple types
    switch v := value.(type) {
    case string:
        fmt.Printf("String: %s\n", v)
    case int:
        fmt.Printf("Integer: %d\n", v)
    case float64:
        fmt.Printf("Float: %.2f\n", v)
    default:
        fmt.Printf("Unknown type: %T\n", v)
    }
}
```

**Key Learning Points**:
- Type assertions return value and boolean: `value, ok := x.(Type)`
- Type switches are clean ways to handle multiple types
- `%T` in `fmt.Printf` prints the type name

---

## 🏗️ Building a Custom Flogo Activity

### Step 1: Project Structure

Every Flogo activity needs these essential files:

```
my-custom-activity/
├── activity.go          # Main activity implementation
├── descriptor.json      # Activity metadata and schema
├── go.mod              # Go module definition
├── go.sum              # Dependency checksums (auto-generated)
├── activity_test.go    # Unit tests
└── README.md           # Documentation
```

### Step 2: Define Your Activity Struct

```go
package main

import (
    "github.com/project-flogo/core/activity"
    "github.com/project-flogo/core/data/metadata"
)

// Settings defines configuration options for your activity
type Settings struct {
    OutputFormat string `md:"outputFormat,required"` // Flogo metadata tag
    Prefix       string `md:"prefix"`                // Optional setting
}

// Input defines what data the activity receives
type Input struct {
    Data map[string]interface{} `md:"data,required"`
}

// Output defines what data the activity returns
type Output struct {
    Result string `md:"result"`
    Count  int    `md:"count"`
}

// MyActivity is the main activity struct
type MyActivity struct {
    settings *Settings
}
```

**Key Learning Points**:
- Use `md:` tags for Flogo metadata
- `required` makes fields mandatory
- Separate structs for Settings, Input, and Output keeps code organized

### Step 3: Implement Required Methods

```go
func init() {
    // Register your activity with Flogo
    _ = activity.Register(&MyActivity{}, New)
}

// New creates a new instance of your activity
func New(ctx activity.InitContext) (activity.Activity, error) {
    settings := &Settings{}
    
    // Get settings from Flogo configuration
    err := metadata.MapToStruct(ctx.Settings(), settings, true)
    if err != nil {
        return nil, err
    }
    
    return &MyActivity{settings: settings}, nil
}

// Metadata returns schema information about your activity
func (a *MyActivity) Metadata() *activity.Metadata {
    return activityMd // This comes from generated metadata
}

// Eval performs the main activity logic
func (a *MyActivity) Eval(ctx activity.Context) (done bool, err error) {
    // Get input data
    input := &Input{}
    err = ctx.GetInputObject(input)
    if err != nil {
        return false, err
    }
    
    // Process the data
    result, count, err := a.processData(input.Data)
    if err != nil {
        return false, err
    }
    
    // Set output
    output := &Output{
        Result: result,
        Count:  count,
    }
    
    err = ctx.SetOutputObject(output)
    if err != nil {
        return false, err
    }
    
    return true, nil
}
```

### Step 4: Implement Your Business Logic

```go
func (a *MyActivity) processData(data map[string]interface{}) (string, int, error) {
    var result strings.Builder
    count := 0
    
    // Add prefix if configured
    if a.settings.Prefix != "" {
        result.WriteString(a.settings.Prefix)
        result.WriteString(": ")
    }
    
    // Process each field in the input data
    for key, value := range data {
        processed, err := a.processField(key, value)
        if err != nil {
            // Log error but continue processing other fields
            ctx.Logger().Warnf("Error processing field %s: %v", key, err)
            continue
        }
        
        result.WriteString(processed)
        result.WriteString("\n")
        count++
    }
    
    return result.String(), count, nil
}

func (a *MyActivity) processField(key string, value interface{}) (string, error) {
    switch a.settings.OutputFormat {
    case "json":
        return a.formatAsJSON(key, value)
    case "text":
        return a.formatAsText(key, value)
    default:
        return "", fmt.Errorf("unsupported output format: %s", a.settings.OutputFormat)
    }
}

func (a *MyActivity) formatAsJSON(key string, value interface{}) (string, error) {
    data := map[string]interface{}{
        "field": key,
        "value": value,
    }
    
    jsonData, err := json.Marshal(data)
    if err != nil {
        return "", err
    }
    
    return string(jsonData), nil
}

func (a *MyActivity) formatAsText(key string, value interface{}) (string, error) {
    return fmt.Sprintf("%s = %v", key, value), nil
}
```

**Key Learning Points**:
- Use `strings.Builder` for efficient string concatenation
- Range over maps with `for key, value := range mapVar`
- Handle errors gracefully - log warnings for non-critical errors
- Use switch statements for multiple format options

### Step 5: Create the Descriptor

The `descriptor.json` file tells Flogo about your activity:

```json
{
  "name": "my-custom-activity",
  "type": "flogo:activity",
  "version": "1.0.0",
  "title": "My Custom Activity",
  "description": "An example custom activity for learning",
  "author": "Your Name <your.email@company.com>",
  "settings": [
    {
      "name": "outputFormat",
      "type": "string",
      "required": true,
      "allowed": ["json", "text"],
      "value": "text"
    },
    {
      "name": "prefix",
      "type": "string",
      "required": false
    }
  ],
  "input": [
    {
      "name": "data",
      "type": "object",
      "required": true
    }
  ],
  "output": [
    {
      "name": "result",
      "type": "string"
    },
    {
      "name": "count",
      "type": "int"
    }
  ]
}
```

### Step 6: Write Tests

Testing is crucial in Go. Here's how to test your activity:

```go
package main

import (
    "testing"
    "github.com/project-flogo/core/activity"
    "github.com/project-flogo/core/support/test"
    "github.com/stretchr/testify/assert"
)

func TestActivity_Eval(t *testing.T) {
    // Create activity settings
    settings := map[string]interface{}{
        "outputFormat": "text",
        "prefix":       "DATA",
    }
    
    // Create mock activity context
    mf := mapper.NewFactory(resolve.GetBasicResolver())
    ctx := test.NewActivityInitContext(settings, mf)
    
    // Create activity instance
    act, err := New(ctx)
    assert.NoError(t, err)
    
    // Create test input
    tc := test.NewActivityContext(act.Metadata())
    input := map[string]interface{}{
        "name":  "John",
        "age":   30,
        "city":  "New York",
    }
    tc.SetInput("data", input)
    
    // Execute activity
    done, err := act.Eval(tc)
    
    // Verify results
    assert.True(t, done)
    assert.NoError(t, err)
    
    result := tc.GetOutput("result").(string)
    count := tc.GetOutput("count").(int)
    
    assert.Contains(t, result, "DATA:")
    assert.Contains(t, result, "name = John")
    assert.Equal(t, 3, count)
}

func TestActivity_InvalidFormat(t *testing.T) {
    settings := map[string]interface{}{
        "outputFormat": "invalid",
    }
    
    mf := mapper.NewFactory(resolve.GetBasicResolver())
    ctx := test.NewActivityInitContext(settings, mf)
    act, err := New(ctx)
    assert.NoError(t, err)
    
    tc := test.NewActivityContext(act.Metadata())
    tc.SetInput("data", map[string]interface{}{"test": "value"})
    
    done, err := act.Eval(tc)
    
    assert.False(t, done)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "unsupported output format")
}
```

**Key Learning Points**:
- Use table-driven tests for multiple test cases
- Test both success and error scenarios
- Use `assert` package for clean test assertions
- Mock external dependencies in tests

### Step 7: Go Module Setup

Create your `go.mod` file:

```bash
# Initialize a new Go module
go mod init github.com/yourusername/my-custom-activity

# Add required dependencies
go get github.com/project-flogo/core/activity
go get github.com/project-flogo/core/data/metadata
go get github.com/stretchr/testify/assert

# Clean up dependencies
go mod tidy
```

Your `go.mod` should look like:

```go
module github.com/yourusername/my-custom-activity

go 1.19

require (
    github.com/project-flogo/core v1.6.0
    github.com/stretchr/testify v1.8.0
)

require (
    // Indirect dependencies will be listed here automatically
)
```

---

## 🧪 Advanced Go Concepts in Our Prometheus Activity

### 1. Working with JSON

Our activity heavily uses JSON processing:

```go
// Unmarshal JSON into a map
var data map[string]interface{}
err := json.Unmarshal(jsonBytes, &data)
if err != nil {
    return fmt.Errorf("invalid JSON: %v", err)
}

// Check if a field exists in JSON
if value, exists := data["metrics"]; exists {
    // Field exists, process it
    if metricsArray, ok := value.([]interface{}); ok {
        // It's an array, process each element
        for _, item := range metricsArray {
            // Process each metric
        }
    }
}

// Marshal Go data to JSON
result := map[string]interface{}{
    "status": "success",
    "count":  len(data),
}
jsonBytes, err := json.Marshal(result)
```

### 2. String Building and Formatting

Efficient string operations are important for performance:

```go
func buildPrometheusMetric(metricName string, labels map[string]string, value float64) string {
    var builder strings.Builder
    
    // Pre-allocate space if you know approximate size
    builder.Grow(100)
    
    builder.WriteString(metricName)
    
    if len(labels) > 0 {
        builder.WriteString("{")
        first := true
        for key, val := range labels {
            if !first {
                builder.WriteString(",")
            }
            // Escape special characters in label values
            escapedValue := strings.ReplaceAll(val, `"`, `\"`)
            builder.WriteString(fmt.Sprintf(`%s="%s"`, key, escapedValue))
            first = false
        }
        builder.WriteString("}")
    }
    
    builder.WriteString(fmt.Sprintf(" %g", value))
    
    return builder.String()
}
```

### 3. Regular Expressions for Validation

Prometheus has strict naming rules that we enforce with regex:

```go
import "regexp"

var (
    // Compile regex patterns once, reuse many times
    metricNamePattern = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
    labelNamePattern  = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
)

func sanitizeMetricName(name string) string {
    // Replace invalid characters with underscores
    sanitized := regexp.MustCompile(`[^a-zA-Z0-9_]`).ReplaceAllString(name, "_")
    
    // Ensure it starts with letter or underscore
    if len(sanitized) > 0 && !regexp.MustCompile(`^[a-zA-Z_]`).MatchString(sanitized) {
        sanitized = "_" + sanitized
    }
    
    return sanitized
}

func isValidLabelName(name string) bool {
    return labelNamePattern.MatchString(name)
}
```

### 4. Logging in Flogo Activities

Proper logging helps with debugging:

```go
func (a *PrometheusActivity) Eval(ctx activity.Context) (done bool, err error) {
    logger := ctx.Logger()
    
    // Different log levels
    logger.Debug("Starting prometheus activity evaluation")
    logger.Info("Processing metric data")
    logger.Warn("Field type not supported, skipping")
    logger.Error("Failed to process metric data")
    
    // Structured logging with fields
    logger.Debugf("Processing field '%s' with value '%v'", fieldName, value)
    
    return true, nil
}
```

---

## 🛠️ Development Workflow

### 1. Setting Up Your Development Environment

```bash
# Create project directory
mkdir my-custom-activity
cd my-custom-activity

# Initialize Go module
go mod init github.com/yourusername/my-custom-activity

# Create basic structure
touch activity.go descriptor.json activity_test.go README.md

# Install dependencies
go get github.com/project-flogo/core/activity
go get github.com/stretchr/testify/assert
```

### 2. Development Loop

```bash
# 1. Write/modify code
vim activity.go

# 2. Run tests
go test -v

# 3. Build to check for compile errors
go build

# 4. Run tests with coverage
go test -cover

# 5. Format code (Go standard)
go fmt ./...

# 6. Check for common issues
go vet ./...

# 7. Update dependencies
go mod tidy
```

### 3. Building and Testing

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v

# Run tests with coverage report
go test -cover -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run specific test
go test -run TestActivity_Eval

# Run benchmarks
go test -bench=.

# Build for different platforms
GOOS=linux GOARCH=amd64 go build
GOOS=windows GOARCH=amd64 go build
```

---

## 📖 Go Learning Resources

### Essential Go Concepts for Flogo Development

1. **Basic Syntax**: Variables, functions, control structures
2. **Types**: Basic types, structs, interfaces, pointers
3. **Error Handling**: Error interface, error wrapping, best practices
4. **Packages**: Importing, exporting, module system
5. **JSON**: Marshaling, unmarshaling, struct tags
6. **Testing**: Unit tests, table-driven tests, mocking
7. **Concurrency**: Goroutines and channels (for advanced activities)

### Recommended Learning Path

1. **Start Here**: [Tour of Go](https://tour.golang.org/)
2. **Go by Example**: [gobyexample.com](https://gobyexample.com/)
3. **Effective Go**: [golang.org/doc/effective_go](https://golang.org/doc/effective_go)
4. **Flogo Documentation**: [project-flogo.github.io](https://project-flogo.github.io/)

### Common Go Patterns in Flogo Activities

```go
// 1. Option pattern for flexible configuration
type ActivityOption func(*MyActivity)

func WithTimeout(timeout time.Duration) ActivityOption {
    return func(a *MyActivity) {
        a.timeout = timeout
    }
}

func NewActivity(opts ...ActivityOption) *MyActivity {
    a := &MyActivity{
        timeout: 30 * time.Second, // default
    }
    for _, opt := range opts {
        opt(a)
    }
    return a
}

// 2. Builder pattern for complex objects
type MetricBuilder struct {
    name   string
    labels map[string]string
    value  float64
}

func NewMetricBuilder(name string) *MetricBuilder {
    return &MetricBuilder{
        name:   name,
        labels: make(map[string]string),
    }
}

func (b *MetricBuilder) WithLabel(key, value string) *MetricBuilder {
    b.labels[key] = value
    return b
}

func (b *MetricBuilder) WithValue(value float64) *MetricBuilder {
    b.value = value
    return b
}

func (b *MetricBuilder) Build() string {
    // Build the metric string
    return "metric_line"
}

// Usage:
metric := NewMetricBuilder("http_requests").
    WithLabel("method", "GET").
    WithLabel("status", "200").
    WithValue(42).
    Build()
```

---

## 🚀 Next Steps

After completing this tutorial, you'll be able to:

1. **Create Custom Activities**: Build activities for your specific use cases
2. **Understand Go**: Apply Go language concepts in real projects  
3. **Debug Issues**: Use logging and testing to troubleshoot problems
4. **Extend Flogo**: Add new capabilities to your Flogo applications
5. **Contribute**: Understand how to contribute to open-source Flogo projects

### Ideas for Your Next Activity

- **Database Connector**: Connect to specific databases
- **Message Transformer**: Convert between different message formats
- **Authentication Handler**: Implement custom auth mechanisms
- **Metrics Collector**: Gather metrics from various sources
- **File Processor**: Handle specific file formats
- **API Client**: Integrate with specific REST APIs

Remember: Start simple, test thoroughly, and gradually add complexity. The Prometheus Metrics activity started as a simple JSON-to-text converter and evolved into a comprehensive metrics processor!

---

## 📋 Quick Reference

### Essential Go Commands
```bash
go mod init module-name    # Initialize new module
go mod tidy               # Clean up dependencies
go build                  # Compile package
go test                   # Run tests
go fmt                    # Format code
go vet                    # Check for issues
go get package-name       # Add dependency
```

### Flogo Activity Checklist
- [ ] Implement `activity.Activity` interface
- [ ] Create `descriptor.json` with schema
- [ ] Handle input validation and errors
- [ ] Write comprehensive tests
- [ ] Document your activity
- [ ] Follow Go naming conventions
- [ ] Use proper logging levels

### Common Gotchas
- Always check for `nil` pointers
- Use pointer receivers for methods that modify structs
- Handle JSON unmarshaling errors
- Don't forget to export struct fields (capital letters)
- Test edge cases and error conditions
- Use `go mod tidy` to keep dependencies clean

Happy coding! 🎉
