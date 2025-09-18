# Troubleshooting Empty Output Issue

## Issue Description
The pongo2-prompt activity is returning empty output when logged in a Flogo application.

## Enhanced Logging (SOLUTION)

The activity has been updated with enhanced logging to help diagnose empty output issues. The enhanced version will now log detailed information about:

1. **Template Processing**: Original and trimmed output lengths
2. **Empty Output Detection**: Special warnings when output is empty
3. **Variable Information**: What variables were provided vs expected
4. **Output Setting**: Confirmation when output is set

## Diagnostic Steps

### 1. Check Activity Logs
Look for these specific log messages in your Flogo application logs:

```
INFO  - Template variables detected: [variable1, variable2, ...]
INFO  - Setting output 'renderedPrompt' with X characters
WARN  - Template expects variables that are not provided: [missing_vars]
WARN  - WARNING: Rendered output is empty!
```

### 2. Common Causes of Empty Output

#### A. Missing Template Variables
**Symptoms:**
- Log shows: `"Template expects variables that are not provided: [...]"`
- Output contains placeholder text like "Hello !"

**Solution:**
- Ensure all template variables (e.g., `{{ name }}`) have corresponding values in the `templateVariables` input
- Use the schema generation utility to identify required variables:
  ```bash
  cd utils
  ./generate_flogo_params.sh "Your template here"
  ```

#### B. Nil or Empty Variable Values
**Symptoms:**
- Variables are provided but output still contains empty placeholders
- Log shows: `"No variables provided to template"` even when variables are set

**Solution:**
- Check that variable values are not `null` or `nil`
- Ensure string values are not empty (`""`)
- Variables with `null` values are filtered out by the activity

#### C. Template Syntax Errors
**Symptoms:**
- Activity fails with parsing errors
- No output is generated

**Solution:**
- Verify pongo2 template syntax (not pure Jinja2)
- Use `{{ variable }}` for variables, not `{variable}`
- Use `{% for %}` loops, not `{%for%}` (spaces matter)
- Use `forloop.Counter` instead of `loop.index`

#### D. Input Field Mapping Issues
**Symptoms:**
- Variables appear to be set but aren't received by the activity
- Log shows: `"No variables provided to template"`

**Solution:**
- Verify the `templateVariables` input is properly mapped in Flogo Web UI
- Ensure you're using the `params` type correctly in the descriptor
- Check that parameter names in Flogo Web UI match template variable names exactly

### 3. Testing Your Template

Use these test patterns to isolate the issue:

#### Test 1: Plain Text (No Variables)
```
Template: "This is a plain text test."
Expected Output: "This is a plain text test."
```

#### Test 2: Simple Variable
```
Template: "Hello {{ name }}!"
Variables: {"name": "World"}
Expected Output: "Hello World!"
```

#### Test 3: Missing Variable
```
Template: "Hello {{ name }}!"
Variables: {} (empty)
Expected Output: "Hello !" (with warning logs)
```

### 4. Enhanced Logging Output Examples

#### Normal Operation:
```
INFO  - Template variables detected: [name, role]
INFO  - Setting output 'renderedPrompt' with 25 characters
DEBUG - Output 'renderedPrompt' has been set successfully
```

#### Empty Output Warning:
```
WARN  - WARNING: Rendered output is empty!
WARN  - Original template: 'Hello {{ name }}!'
WARN  - Template variables provided: map[]
WARN  - Raw rendered output: ''
```

#### Missing Variables:
```
INFO  - Template variables detected: [name, age]
WARN  - Template expects variables that are not provided: [age]
INFO  - Provide missing variables through the template variables input fields
```

### 5. Flogo Web UI Configuration

Ensure your activity is configured correctly:

1. **Template Input**: 
   - Use the large text editor (15 rows)
   - Enter your pongo2 template with `{{ variable }}` syntax

2. **Template Variables Input**:
   - Add a parameter for each variable in your template
   - Parameter names must match template variable names exactly
   - Map each parameter to appropriate data sources

3. **Output Mapping**:
   - Map the `renderedPrompt` output to your target field
   - Verify the output field name is spelled correctly

### 6. Quick Verification Commands

```bash
# Test template variable detection
cd /opt/tibco/flogo-extensions/pongo2-prompt/utils
./generate_flogo_params.sh "Your template content here"

# Run unit tests to verify activity functionality
cd /opt/tibco/flogo-extensions/pongo2-prompt
go test -v -run TestJinja2PromptActivity

# Test specific scenarios
go test -v -run TestEmptyOutputScenarios
```

### 7. Common Template Issues

#### Incorrect Pongo2 Syntax:
```
❌ WRONG: Hello {name}!
❌ WRONG: Hello {{name}}! (missing spaces)
❌ WRONG: {% for item in items%} (missing space before %)
✅ CORRECT: Hello {{ name }}!
✅ CORRECT: {% for item in items %}
```

#### Loop Variables:
```
❌ WRONG: {{ loop.index }} (Jinja2 syntax)
✅ CORRECT: {{ forloop.Counter }} (Pongo2 syntax)
```

### 8. If Issue Persists

If you're still getting empty output after checking all the above:

1. **Enable Debug Logging**: Set your Flogo application to DEBUG level
2. **Check All Logs**: Look for the enhanced logging messages added to the activity
3. **Test Incrementally**: Start with a plain text template, then add variables one by one
4. **Verify Input Data**: Log the input data before it reaches the activity
5. **Check Output Mapping**: Verify how the output is being used after the activity

The enhanced logging should provide clear indicators of where the issue lies in the template processing pipeline.