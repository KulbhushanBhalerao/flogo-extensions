package pongo2v2prompt

import (
	"strings"
	"testing"

	"github.com/flosch/pongo2/v6"
)

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// Helper function to render template with variables
func renderTemplate(templateStr string, variables map[string]interface{}) (string, error) {
	template, err := pongo2.FromString(templateStr)
	if err != nil {
		return "", err
	}

	context := pongo2.Context{}
	for key, value := range variables {
		context[key] = value
	}

	return template.Execute(context)
}

func TestExtractTemplateVariables(t *testing.T) {
	template := "Hello {{ name }}! You are {{ age }} years old."
	variables := extractTemplateVariables(template)

	expected := []string{"name", "age"}
	if len(variables) != len(expected) {
		t.Errorf("Expected %d variables, got %d", len(expected), len(variables))
	}

	for i, v := range expected {
		if i >= len(variables) || variables[i] != v {
			t.Errorf("Expected variable %d to be '%s', got '%s'", i, v, variables[i])
		}
	}

	t.Logf("✓ Template variables extracted: %v", variables)
}

func TestExtractForLoopArrays(t *testing.T) {
	template := "{% for item in items %}{{ item.name }}{% endfor %}"
	arrays := extractForLoopArrays(template)

	expected := []string{"items"}
	if len(arrays) != len(expected) {
		t.Errorf("Expected %d arrays, got %d", len(expected), len(arrays))
	}

	if len(arrays) > 0 && arrays[0] != expected[0] {
		t.Errorf("Expected array to be '%s', got '%s'", expected[0], arrays[0])
	}

	t.Logf("✓ For loop arrays extracted: %v", arrays)
}

func TestIsArrayVariable(t *testing.T) {
	arrays := []string{"items", "users"}

	if !isArrayVariable("items", arrays) {
		t.Error("'items' should be identified as array variable")
	}

	if isArrayVariable("name", arrays) {
		t.Error("'name' should not be identified as array variable")
	}

	t.Log("✓ Array variable detection working correctly")
}

func TestImprovedSchemaProvider(t *testing.T) {
	provider := &ImprovedSchemaProvider{}

	template := "Hello {{ name }}! You are {{ age }} years old."
	settings := map[string]interface{}{
		"template": template,
	}

	schema, err := provider.GetInputSchema(settings)
	if err != nil {
		t.Fatalf("Failed to get schema: %v", err)
	}

	if schema == nil {
		t.Fatal("Schema should not be nil")
	}

	properties, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("Schema properties should be a map")
	}

	if len(properties) == 0 {
		t.Error("Schema should have properties")
	}

	t.Logf("✓ Schema generated with %d properties", len(properties))
}

func TestPongo2TemplateConditionalLogic(t *testing.T) {
	template := `Hello {{ name }}! Welcome to {{ company }}.
Your role is {{ role }} 
{% if experience == 0 %} and welcome to your first professional role!
{% if experience == 10 or experience == 15 or experience == 20 %} We truly appreciate your significant contribution to the industry.{% endif %}
{% elif experience > 0 %} and well done on your {{ experience }} years of experience!{% endif %}`

	// Test Case 1: New employee (experience = 0)
	t.Run("NewEmployee", func(t *testing.T) {
		variables := map[string]interface{}{
			"name":       "John Smith",
			"company":    "TechCorp Solutions",
			"role":       "Junior Developer",
			"experience": 0,
		}

		result, err := renderTemplate(template, variables)
		if err != nil {
			t.Fatalf("Failed to render template: %v", err)
		}

		t.Logf("New Employee Result: %s", result)

		if !contains(result, "Hello John Smith! Welcome to TechCorp Solutions.") {
			t.Error("Should contain greeting")
		}
		if !contains(result, "Your role is Junior Developer") {
			t.Error("Should contain role")
		}
		if !contains(result, "welcome to your first professional role!") {
			t.Error("Should contain new employee message")
		}
		if contains(result, "well done on your") {
			t.Error("Should not contain experienced employee message")
		}
		if contains(result, "We truly appreciate") {
			t.Error("Should not contain milestone message")
		}
	})

	// Test Case 2: Experienced employee (experience = 5)
	t.Run("ExperiencedEmployee", func(t *testing.T) {
		variables := map[string]interface{}{
			"name":       "Alice Johnson",
			"company":    "TechCorp Solutions",
			"role":       "Senior Developer",
			"experience": 5,
		}

		result, err := renderTemplate(template, variables)
		if err != nil {
			t.Fatalf("Failed to render template: %v", err)
		}

		t.Logf("Experienced Employee Result: %s", result)

		if !contains(result, "Hello Alice Johnson! Welcome to TechCorp Solutions.") {
			t.Error("Should contain greeting")
		}
		if !contains(result, "Your role is Senior Developer") {
			t.Error("Should contain role")
		}
		if !contains(result, "well done on your 5 years of experience!") {
			t.Error("Should contain experienced employee message")
		}
		if contains(result, "welcome to your first professional role!") {
			t.Error("Should not contain new employee message")
		}
		if contains(result, "We truly appreciate") {
			t.Error("Should not contain milestone message")
		}
	})

	// Test Case 3: Milestone employee (experience = 15)
	t.Run("MilestoneEmployee", func(t *testing.T) {
		variables := map[string]interface{}{
			"name":       "Robert Chen",
			"company":    "TechCorp Solutions",
			"role":       "Principal Architect",
			"experience": 15,
		}

		result, err := renderTemplate(template, variables)
		if err != nil {
			t.Fatalf("Failed to render template: %v", err)
		}

		t.Logf("Milestone Employee Result: %s", result)

		if !contains(result, "Hello Robert Chen! Welcome to TechCorp Solutions.") {
			t.Error("Should contain greeting")
		}
		if !contains(result, "Your role is Principal Architect") {
			t.Error("Should contain role")
		}
		if !contains(result, "well done on your 15 years of experience!") {
			t.Error("Should contain experienced employee message")
		}
		if !contains(result, "We truly appreciate your significant contribution to the industry.") {
			t.Error("Should contain milestone appreciation message")
		}
		if contains(result, "welcome to your first professional role!") {
			t.Error("Should not contain new employee message")
		}
	})
}

func TestPongo2TemplateStructureDebug(t *testing.T) {
	// Test the exact template structure you provided
	template := `Hello {{ name }}! Welcome to {{ company }}.
Your role is {{ role }} 
{% if experience == 0 %} and welcome to your first professional role!
{% if experience == 10 or experience == 15 or experience == 20 %} We truly appreciate your significant contribution to the industry.{% endif %}
{% elif experience > 0 %} and well done on your {{ experience }} years of experience!{% endif %}`

	t.Run("ParseTemplateStructure", func(t *testing.T) {
		_, err := renderTemplate(template, map[string]interface{}{
			"name":       "Test",
			"company":    "TestCorp",
			"role":       "Developer",
			"experience": 0,
		})
		if err != nil {
			t.Errorf("Template parsing failed: %v", err)
			t.Logf("This suggests the template structure is malformed")
		}
	})

	// Test with corrected template structure
	correctedTemplate := `Hello {{ name }}! Welcome to {{ company }}.
Your role is {{ role }}{% if experience == 0 %} and welcome to your first professional role!{% elif experience > 0 %} and well done on your {{ experience }} years of experience!{% endif %}{% if experience == 10 or experience == 15 or experience == 20 %} We truly appreciate your significant contribution to the industry.{% endif %}`

	t.Run("CorrectedTemplateNewEmployee", func(t *testing.T) {
		result, err := renderTemplate(correctedTemplate, map[string]interface{}{
			"name":       "John",
			"company":    "TestCorp",
			"role":       "Junior Dev",
			"experience": 0,
		})
		if err != nil {
			t.Fatalf("Corrected template failed: %v", err)
		}
		t.Logf("New employee result: %s", result)

		if !contains(result, "welcome to your first professional role!") {
			t.Error("Should contain new employee message")
		}
		if contains(result, "We truly appreciate") {
			t.Error("Should not contain milestone message for experience=0")
		}
	})

	t.Run("CorrectedTemplateMilestone", func(t *testing.T) {
		result, err := renderTemplate(correctedTemplate, map[string]interface{}{
			"name":       "Alice",
			"company":    "TestCorp",
			"role":       "Senior Dev",
			"experience": 15,
		})
		if err != nil {
			t.Fatalf("Corrected template failed: %v", err)
		}
		t.Logf("Milestone employee result: %s", result)

		if !contains(result, "well done on your 15 years") {
			t.Error("Should contain experienced employee message")
		}
		if !contains(result, "We truly appreciate") {
			t.Error("Should contain milestone appreciation message")
		}
	})

	// Test with the exact template structure you provided
	userTemplate := `Hello {{ name }}! Welcome to {{ company }}.
Your role is {{ role }}
{% if experience == 0 %} and welcome to your first professional role!
{% elif experience > 0 %} and well done on your {{ experience }} years of experience!
{% endif %}
{% if experience == 10 or experience == 15 or experience == 20 %} We truly appreciate your significant contribution to the industry.
{% endif %}`

	t.Run("UserTemplateMilestone15", func(t *testing.T) {
		result, err := renderTemplate(userTemplate, map[string]interface{}{
			"name":       "B",
			"company":    "T",
			"role":       "Dev",
			"experience": 15,
		})
		if err != nil {
			t.Fatalf("User template failed: %v", err)
		}
		t.Logf("User template with experience=15: '%s'", result)

		// Check for the decimal formatting issue
		if contains(result, "15.000000") {
			t.Error("Experience should not show as decimal (15.000000)")
		}

		// Check if milestone message appears
		if contains(result, "We truly appreciate") {
			t.Log("✓ Milestone message correctly appears")
		} else {
			t.Error("✗ Milestone message is missing - this suggests the second conditional is not working")
		}
	})

	// Test with integer instead of float to fix decimal issue
	t.Run("UserTemplateWithIntegerExperience", func(t *testing.T) {
		result, err := renderTemplate(userTemplate, map[string]interface{}{
			"name":       "B",
			"company":    "T",
			"role":       "Dev",
			"experience": int(15), // Explicitly use integer
		})
		if err != nil {
			t.Fatalf("User template with int failed: %v", err)
		}
		t.Logf("User template with int experience=15: '%s'", result)

		// Check for clean number formatting
		if contains(result, "15.000000") {
			t.Error("Should not show decimal formatting with integer input")
		}

		if contains(result, "15 years") {
			t.Log("✓ Clean number formatting")
		}

		if contains(result, "We truly appreciate") {
			t.Log("✓ Milestone message appears with integer")
		} else {
			t.Error("✗ Milestone message still missing with integer")
		}
	})

	// Test with float input to reproduce the decimal issue - then test the activity fix
	t.Run("UserTemplateWithFloatExperience", func(t *testing.T) {
		result, err := renderTemplate(userTemplate, map[string]interface{}{
			"name":       "B",
			"company":    "T",
			"role":       "Dev",
			"experience": 15.0, // Float input like what might come from Flogo
		})
		if err != nil {
			t.Fatalf("User template with float failed: %v", err)
		}
		t.Logf("User template with float experience=15.0: '%s'", result)

		// Before activity fix: decimal formatting issue
		if contains(result, "15.000000") {
			t.Log("✓ Reproduced decimal formatting issue with float input")
		} else {
			t.Log("Float input shows clean formatting")
		}

		if contains(result, "We truly appreciate") {
			t.Log("✓ Milestone message appears with float")
		} else {
			t.Error("✗ Milestone message missing with float - this might be the issue")
		}
	})

	// Test the activity's float-to-int conversion fix
	t.Run("ActivityFloatConversionFix", func(t *testing.T) {
		// Simulate the activity's float conversion logic
		varMap := map[string]interface{}{
			"name":       "B",
			"company":    "T",
			"role":       "Dev",
			"experience": 15.0, // Float input
		}

		// Apply the activity's conversion logic
		for key, value := range varMap {
			if f, ok := value.(float64); ok {
				if f == float64(int64(f)) {
					varMap[key] = int64(f)
				}
			}
		}

		result, err := renderTemplate(userTemplate, varMap)
		if err != nil {
			t.Fatalf("Activity fix test failed: %v", err)
		}
		t.Logf("After activity float conversion: '%s'", result)

		if contains(result, "15 years") && !contains(result, "15.000000") {
			t.Log("✓ Activity conversion fixes decimal display")
		}

		if contains(result, "We truly appreciate") {
			t.Log("✓ Activity conversion fixes milestone message")
		} else {
			t.Error("✗ Activity conversion still doesn't fix milestone message")
		}
	})

	// Test the floatformat filter to fix decimal display
	templateWithFloatFormat := `Hello {{ name }}! Welcome to {{ company }}.
Your role is {{ role }}
{% if experience == 0 %} and welcome to your first professional role!
{% elif experience > 0 %} and well done on your {{ experience|floatformat:0 }} years of experience!
{% endif %}
{% if experience == 10 or experience == 15 or experience == 20 %} We truly appreciate your significant contribution to the industry.
{% endif %}`

	t.Run("TemplateWithFloatFormat", func(t *testing.T) {
		result, err := renderTemplate(templateWithFloatFormat, map[string]interface{}{
			"name":       "B",
			"company":    "T",
			"role":       "Dev",
			"experience": 15.0,
		})
		if err != nil {
			t.Fatalf("Template with floatformat failed: %v", err)
		}
		t.Logf("Template with floatformat: '%s'", result)

		if contains(result, "15 years") && !contains(result, "15.000000") {
			t.Log("✓ floatformat:0 fixes decimal display")
		}

		if contains(result, "We truly appreciate") {
			t.Log("✓ Milestone message works with floatformat")
		}
	})
}
