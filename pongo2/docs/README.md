# Pongo2 Prompt Template Activity

A Flogo activity for processing Pongo2/Django-style templates, specifically designed for AI prompt generation with dynamic variable mapping.

## Overview

This activity uses the **pongo2** library (a Go implementation of Django templates) to process templates with Jinja2-like syntax. While similar to Jinja2, pongo2 has some syntax differences that are important to understand for template compatibility.

## Why Pongo2 Instead of Jinja2?

- **Go Native**: Pongo2 is implemented in pure Go, providing better performance and integration
- **No External Dependencies**: No need for Python runtime or external processes
- **Flogo Optimized**: Designed specifically for Flogo's Go-based architecture
- **Django Heritage**: Based on Django templates with Jinja2-like syntax

## Key Syntax Differences from Jinja2

While pongo2 syntax is very similar to Jinja2, there are important differences:

### Loop Variables
- **Jinja2**: `{{ loop.index }}`
- **Pongo2**: `{{ forloop.Counter }}` (1-based) or `{{ forloop.Counter0 }}` (0-based)

### Complex Calculations
- **Jinja2**: `{% set result = (value / total * 100) | round(2) %}`
- **Pongo2**: Pre-calculate in Go and pass as template variable

### Custom Filters
- **Jinja2**: Supports extensive custom filters
- **Pongo2**: Limited built-in filters, use Go for complex formatting

## Key Features

✅ **Dynamic Variable Mapping**: Support for both individual inputs and complex object mapping  
✅ **Conditional Logic**: `{% if %}`, `{% else %}`, `{% endif %}` statements  
✅ **Loops**: `{% for %}` loops with proper counters (`forloop.Counter`, `forloop.First`, etc.)  
✅ **Built-in Filters**: `default`, `length`, `lower`, `upper`, `date`, etc.  
✅ **Complex Data Structures**: Nested objects and arrays  
✅ **AI Prompt Optimization**: Designed for dynamic prompt generation  

## Quick Start

### 1. Installation
```bash
# From local path
Copy the flogo plugin folder to the extensions folder and from vscode settings for flogo plugin mention this extension folder in extensions. Restart vscode

# Or from repository (when published)
flogo install github.com/your-org/pongo2-prompt
```

### 2. Basic Usage
```json
{
  "activity": {
    "ref": "#pongo2-prompt",
    "input": {
      "template": "Hello {{ name }}, you are {{ age }} years old and work as {{ role }}.",
      "templateVariables": {
        "name": "John",
        "age": 30,
        "role": "Developer"
      }
    }
  }
}
```

**Output**: `"Hello John, you are 30 years old and work as Developer."`

### 3. Complex Data Example
```json
{
  "activity": {
    "ref": "#pongo2-prompt",
    "input": {
      "template": "{% for task in tasks %}{{ forloop.Counter }}. {{ task.title }} ({{ task.status }}){% endfor %}",
      "templateVariables": {
        "tasks": [
          {"title": "Setup database", "status": "completed"},
          {"title": "Create API", "status": "in-progress"}
        ]
      }
    }
  }
}
```

**Output**: `"1. Setup database (completed)2. Create API (in-progress)"`

## Important Syntax Differences

### Loop Variables
- **Standard Jinja2**: `{{ loop.index }}`
- **Pongo2**: `{{ forloop.Counter }}`

### Mathematical Operations
- **Standard Jinja2**: `{% set result = (value1 + value2) / 2 %}`
- **Pongo2**: Pre-calculate in Go code and pass as variable

### Available Loop Variables
- `forloop.Counter` - 1-based index
- `forloop.Counter0` - 0-based index  
- `forloop.First` - boolean, true if first iteration
- `forloop.Last` - boolean, true if last iteration

## Configuration

### Inputs
1. **template** (string, required): Pongo2 template string with variables like `{{ variable_name }}`
   - Large text editor with 15 rows for complex templates
   - Variables are automatically detected and mapped to individual input fields
   - Example: `"Hello {{ name }}! Please analyze {{ data }} and provide results in {{ format }} format."`

2. **templateVariables** (object): Auto-generated mappable fields based on your template
   - Individual fields are created for each variable detected in your template
   - Each field can be mapped to different data sources in your flow
   - Schema-based approach ensures proper field generation in Flogo Web UI

### Outputs
- **renderedPrompt**: The processed template with variables substituted

## Advanced Examples

### AI Data Analysis Prompt
```pongo2
You are a data analyst working with {{ dataset_type }} data.

**Analysis Objective:** {{ objective }}

**Key Variables:**
{% for var in variables %}
- **{{ var.name }}** ({{ var.type }}){% if var.description %}: {{ var.description }}{% endif %}
{% endfor %}

**Tasks:**
{% for task in tasks %}
{{ forloop.Counter }}. {{ task }}
{% endfor %}
```

### Conditional Logic
```pongo2
{% if user.role == "admin" %}
You have administrative privileges.
{% elif user.role == "manager" %}  
You have management access.
{% else %}
You have standard user access.
{% endif %}
```



## Use Cases

🤖 **AI Prompt Templates**: Dynamic prompt generation for LLMs  
📊 **Report Generation**: Data-driven report templates  
📧 **Email Templates**: Personalized email content  
🔄 **Workflow Automation**: Dynamic content in automated processes  
📋 **Documentation**: Auto-generated documentation with data injection  

## Performance Notes

- Templates are parsed on each execution (future versions may include caching)
- Pre-calculate complex mathematical operations in Go code for better performance
- Use the `variables` input for complex data structures rather than many individual inputs

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass: `go test -v`
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For issues, questions, or contributions:
- Create an issue in the repository
- Check the documentation guides for detailed examples
- Review the test files for implementation patterns

---

## Documentation

- **[Pongo2 Examples](PONGO2_EXAMPLES.md)** - Comprehensive template examples with sample data and outputs
- **[JSON Schema Guide](JSON_SCHEMA_GUIDE.md)** - Generate schemas for template variables  
- **[Troubleshooting Guide](TROUBLESHOOTING.md)** - Diagnose and fix empty output issues  
- **[Build & Test Guide](BUILD_TEST_GUIDE.md)** - Development, testing, and deployment  
- **[Flogo Integration Guide](FLOGO_INTEGRATION_GUIDE.md)** - Step-by-step integration in Flogo flows  

**Note**: This activity uses pongo2 (Django template syntax) rather than pure Jinja2. While very similar, there are syntax differences as outlined in the "Key Syntax Differences" section above.