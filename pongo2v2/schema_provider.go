package pongo2v2prompt

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
)

// extractTemplateVariables extracts variable names from pongo2 template using regex
func extractTemplateVariables(template string) []string {
	re := regexp.MustCompile(`\{\{\s*([^}|]+)(?:\|[^}]*)?\s*\}\}`)
	matches := re.FindAllStringSubmatch(template, -1)

	var variables []string
	seen := make(map[string]bool)
	for _, match := range matches {
		if len(match) > 1 {
			varName := strings.TrimSpace(match[1])
			if varName != "" && !seen[varName] {
				variables = append(variables, varName)
				seen[varName] = true
			}
		}
	}
	return variables
}

// extractForLoopArrays finds arrays used in for loops
func extractForLoopArrays(template string) []string {
	re := regexp.MustCompile(`\{\%\s*for\s+\w+\s+in\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\%\}`)
	matches := re.FindAllStringSubmatch(template, -1)

	var arrays []string
	for _, match := range matches {
		if len(match) > 1 {
			arrays = append(arrays, match[1])
		}
	}
	return arrays
}

// isArrayVariable checks if a variable is an array used in for loops
func isArrayVariable(variable string, forLoopArrays []string) bool {
	for _, array := range forLoopArrays {
		if variable == array {
			return true
		}
	}
	return false
}

// ImprovedSchemaProvider generates better schema with smart type detection
type ImprovedSchemaProvider struct{}

// detectVariableType attempts to determine the likely data type based on variable name and context
func detectVariableType(variableName string, template string) string {
	// Convert to lowercase for pattern matching
	lowerName := strings.ToLower(variableName)

	// Numeric indicators in variable names
	numericPatterns := []string{
		"age", "count", "number", "num", "amount", "price", "cost", "year", "month", "day",
		"experience", "years", "months", "days", "quantity", "size", "length", "width", "height",
		"score", "rating", "level", "rank", "percentage", "percent", "total", "sum", "average",
		"min", "max", "id", "index", "position", "order", "sequence", "version",
	}

	// Boolean indicators in variable names
	booleanPatterns := []string{
		"is_", "has_", "can_", "should_", "will_", "was_", "are_", "enabled", "disabled",
		"active", "inactive", "valid", "invalid", "visible", "hidden", "public", "private",
		"available", "unavailable", "allowed", "denied", "confirmed", "verified",
	}

	// Check for numeric patterns
	for _, pattern := range numericPatterns {
		if strings.Contains(lowerName, pattern) {
			return "number"
		}
	}

	// Check for boolean patterns
	for _, pattern := range booleanPatterns {
		if strings.Contains(lowerName, pattern) {
			return "boolean"
		}
	}

	// Look for conditional usage patterns in the template to infer numeric types
	// Check for mathematical comparisons
	mathPatterns := []regexp.Regexp{
		*regexp.MustCompile(`\{\%\s*if\s+` + regexp.QuoteMeta(variableName) + `\s*[><=!]+\s*\d+`),
		*regexp.MustCompile(`\{\%\s*elif\s+` + regexp.QuoteMeta(variableName) + `\s*[><=!]+\s*\d+`),
		*regexp.MustCompile(`\{\{\s*` + regexp.QuoteMeta(variableName) + `\s*\|\s*floatformat`),
	}

	for _, pattern := range mathPatterns {
		if pattern.MatchString(template) {
			return "number"
		}
	}

	// Check for equality comparisons with numbers
	eqPattern := regexp.MustCompile(`\{\%\s*if\s+` + regexp.QuoteMeta(variableName) + `\s*==\s*\d+`)
	if eqPattern.MatchString(template) {
		return "number"
	}

	// Default to string
	return "string"
}

// GetInputSchema generates schema with intelligent type detection
func (isp *ImprovedSchemaProvider) GetInputSchema(settings map[string]interface{}) (map[string]interface{}, error) {
	template, ok := settings["template"].(string)
	if !ok || template == "" {
		return nil, nil
	}

	variables := extractTemplateVariables(template)
	forLoopArrays := extractForLoopArrays(template)

	properties := make(map[string]interface{})
	for _, variable := range variables {
		if isArrayVariable(variable, forLoopArrays) {
			properties[variable] = map[string]interface{}{
				"type":        "array",
				"description": "Array variable used in for loop: {{ " + variable + " }}",
				"items": map[string]interface{}{
					"type":                 "object",
					"additionalProperties": true,
				},
			}
		} else {
			// Use intelligent type detection
			detectedType := detectVariableType(variable, template)
			properties[variable] = map[string]interface{}{
				"type":        detectedType,
				"description": "Template variable: {{ " + variable + " }} (detected as " + detectedType + ")",
			}
		}
	}

	schema := map[string]interface{}{
		"$schema":    "http://json-schema.org/draft-04/schema#",
		"type":       "object",
		"properties": properties,
	}

	return schema, nil
}

// normalizeVariableValue converts input values to appropriate types for template processing
func normalizeVariableValue(value interface{}) interface{} {
	if value == nil {
		return nil
	}

	// If it's already the right type, return as-is
	switch v := value.(type) {
	case int, int32, int64, float32, float64:
		return v
	case bool:
		return v
	case string:
		// Try to convert strings that look like numbers
		if v == "" {
			return v
		}

		// Try boolean conversion
		if strings.ToLower(v) == "true" {
			return true
		}
		if strings.ToLower(v) == "false" {
			return false
		}

		// Try integer conversion
		if intVal, err := strconv.Atoi(v); err == nil {
			return intVal
		}

		// Try float conversion
		if floatVal, err := strconv.ParseFloat(v, 64); err == nil {
			return floatVal
		}

		// Return as string if no conversion possible
		return v
	default:
		return v
	}
}

// NormalizeTemplateVariables processes all template variables for type consistency
func NormalizeTemplateVariables(templateVars map[string]interface{}) map[string]interface{} {
	normalized := make(map[string]interface{})

	for key, value := range templateVars {
		normalized[key] = normalizeVariableValue(value)
	}

	return normalized
}

// GetImprovedTemplateSchemaAsJSON returns improved JSON schema with type detection
func GetImprovedTemplateSchemaAsJSON(template string) string {
	provider := &ImprovedSchemaProvider{}
	settings := map[string]interface{}{
		"template": template,
	}

	schema, err := provider.GetInputSchema(settings)
	if err != nil || schema == nil {
		return "{}"
	}

	schemaBytes, err := json.Marshal(schema)
	if err != nil {
		return "{}"
	}

	return string(schemaBytes)
}
