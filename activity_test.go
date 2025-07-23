package prometheusmetrics

import (
	"testing"

	"github.com/project-flogo/core/support/test"
	"github.com/stretchr/testify/assert"
)

func TestActivity_TempHumidityExample(t *testing.T) {
	// Setup activity exactly as you described
	act := &Activity{
		metricType:  "gauge",
		metricName:  "flogo_metric",
		includeHelp: true,
		includeType: true,
		timestamp:   false,
	}

	// Create test context
	tc := test.NewActivityContext(act.Metadata())

	// Your exact example input
	input := &Input{
		MetricData: map[string]interface{}{
			"temp":     1,
			"humidity": 3,
		},
	}

	tc.SetInputObject(input)

	// Execute
	done, err := act.Eval(tc)

	// Assertions
	assert.True(t, done)
	assert.NoError(t, err)

	output := tc.GetOutput("prometheusMetric")
	assert.NotNil(t, output)

	outputStr := output.(string)
	t.Logf("Output for temp/humidity example:\n%s", outputStr)

	// Verify it contains the exact output you wanted
	assert.Contains(t, outputStr, `flogo_metric{name="temp"} 1`)
	assert.Contains(t, outputStr, `flogo_metric{name="humidity"} 3`)
}

func TestActivity_MultipleMetrics(t *testing.T) {
	// Setup activity
	act := &Activity{
		metricType:  "gauge",
		metricName:  "flogo_metric",
		includeHelp: true,
		includeType: true,
		timestamp:   false,
	}

	// Create test context
	tc := test.NewActivityContext(act.Metadata())

	// Test case: Multiple metrics from JSON
	input := &Input{
		MetricData: map[string]interface{}{
			"temp":     1,
			"humidity": 3,
		},
	}

	tc.SetInputObject(input)

	// Execute
	done, err := act.Eval(tc)

	// Assertions
	assert.True(t, done)
	assert.NoError(t, err)

	output := tc.GetOutput("prometheusMetric")
	assert.NotNil(t, output)

	// The output should contain both metrics with name labels
	outputStr := output.(string)
	assert.Contains(t, outputStr, `flogo_metric{name="temp"} 1`)
	assert.Contains(t, outputStr, `flogo_metric{name="humidity"} 3`)
	assert.Contains(t, outputStr, "# HELP flogo_metric Generated metric from JSON data")
	assert.Contains(t, outputStr, "# TYPE flogo_metric gauge")
}

func TestActivity_Eval(t *testing.T) {
	// Setup activity
	act := &Activity{
		metricType:  "gauge",
		metricName:  "test_metric",
		includeHelp: true,
		includeType: true,
		timestamp:   false,
	}

	// Create test context
	tc := test.NewActivityContext(act.Metadata())

	// Test case 1: Simple metric with value
	input := &Input{
		MetricData: map[string]interface{}{
			"value":   42.5,
			"service": "test-service",
			"env":     "prod",
		},
	}

	tc.SetInputObject(input)

	// Execute
	done, err := act.Eval(tc)

	// Assertions
	assert.True(t, done)
	assert.NoError(t, err)

	output := tc.GetOutput("prometheusMetric")
	assert.NotNil(t, output)

	// With the new implementation, each numeric field becomes a separate metric
	outputStr := output.(string)
	assert.Contains(t, outputStr, "# HELP test_metric Generated metric from JSON data")
	assert.Contains(t, outputStr, "# TYPE test_metric gauge")
	assert.Contains(t, outputStr, `test_metric{name="value",env="prod",service="test-service"} 42.5`)
}

func TestActivity_Eval_Counter(t *testing.T) {
	// Setup activity for counter
	act := &Activity{
		metricType:  "counter",
		metricName:  "http_requests_total",
		includeHelp: true,
		includeType: true,
		timestamp:   false,
	}

	// Create test context
	tc := test.NewActivityContext(act.Metadata())

	// Test case: Counter metric
	input := &Input{
		MetricData: map[string]interface{}{
			"count":    150,
			"endpoint": "/api/users",
			"method":   "GET",
			"status":   "200",
			"help":     "Total HTTP requests processed",
		},
	}

	tc.SetInputObject(input)

	// Execute
	done, err := act.Eval(tc)

	// Assertions
	assert.True(t, done)
	assert.NoError(t, err)

	output := tc.GetOutput("prometheusMetric")
	assert.NotNil(t, output)

	// With the new implementation, each numeric field becomes a separate metric
	outputStr := output.(string)
	assert.Contains(t, outputStr, "# HELP http_requests_total Total HTTP requests processed")
	assert.Contains(t, outputStr, "# TYPE http_requests_total counter")
	// Update to expect both metrics with their labels
	assert.Contains(t, outputStr, `http_requests_total{name="count",endpoint="/api/users",method="GET"} 150`)
	assert.Contains(t, outputStr, `http_requests_total{name="status",endpoint="/api/users",method="GET"} 200`)
}

func TestActivity_Eval_WithTimestamp(t *testing.T) {
	// Setup activity with timestamp
	act := &Activity{
		metricType:  "gauge",
		metricName:  "cpu_usage",
		includeHelp: false,
		includeType: false,
		timestamp:   true,
	}

	// Create test context
	tc := test.NewActivityContext(act.Metadata())

	// Test case: Metric with timestamp
	input := &Input{
		MetricData: map[string]interface{}{
			"metric_value": 75.2,
			"host":         "server-01",
			"timestamp":    int64(1705316200000),
		},
	}

	tc.SetInputObject(input)

	// Execute
	done, err := act.Eval(tc)

	// Assertions
	assert.True(t, done)
	assert.NoError(t, err)

	output := tc.GetOutput("prometheusMetric")
	assert.NotNil(t, output)
	outputStr := output.(string)
	assert.Contains(t, outputStr, `cpu_usage{name="metric_value",host="server-01"} 75.2 1705316200000`)
}

func TestActivity_Eval_NoValue(t *testing.T) {
	// Setup activity
	act := &Activity{
		metricType:  "gauge",
		metricName:  "test_metric",
		includeHelp: true,
		includeType: true,
		timestamp:   false,
	}

	// Create test context
	tc := test.NewActivityContext(act.Metadata())

	// Test case: No numeric value
	input := &Input{
		MetricData: map[string]interface{}{
			"service": "test-service",
			"message": "no numeric data",
		},
	}

	tc.SetInputObject(input)

	// Execute
	done, err := act.Eval(tc)

	// Assertions - Since there are no numeric values, it should fail with an error
	assert.False(t, done)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no numeric value found in metric object")
}

func TestSanitizeLabelName(t *testing.T) {
	act := &Activity{}

	tests := []struct {
		input    string
		expected string
	}{
		{"valid_name", "valid_name"},
		{"123invalid", "_23invalid"},
		{"with-dashes", "with_dashes"},
		{"with.dots", "with_dots"},
		{"with spaces", "with_spaces"},
		{"UPPERCASE", "UPPERCASE"},
		{"mixed-Case.123", "mixed_Case_123"},
	}

	for _, test := range tests {
		result := act.sanitizeLabelName(test.input)
		assert.Equal(t, test.expected, result, "Input: %s", test.input)
	}
}
