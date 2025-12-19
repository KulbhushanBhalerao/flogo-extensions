package prometheusmetrics

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/coerce"
)

// Constants for identifying settings and inputs
const (
	sMetricType  = "metricType"
	sMetricName  = "metricName"
	sIncludeHelp = "includeHelp"
	sIncludeType = "includeType"
	sTimestamp   = "timestamp"
	ivMetricData = "metricData"
)

// activityMd is the metadata for the activity.
var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

// Activity is the structure for your activity.
type Activity struct {
	metricType  string
	metricName  string
	includeHelp bool
	includeType bool
	timestamp   bool
}

func init() {
	_ = activity.Register(&Activity{}, New)
}

// Metadata returns the activity's metadata.
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// New creates a new instance of the Activity.
func New(ctx activity.InitContext) (activity.Activity, error) {
	s := &Settings{}
	err := s.FromMap(ctx.Settings())
	if err != nil {
		return nil, err
	}

	act := &Activity{
		metricType:  s.MetricType,
		metricName:  s.MetricName,
		includeHelp: s.IncludeHelp,
		includeType: s.IncludeType,
		timestamp:   s.Timestamp,
	}
	return act, nil
}

// Eval executes the main logic of the Activity.
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {
	logger := ctx.Logger()

	// --- 1. Get Inputs ---
	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return false, err
	}

	if input.MetricData == nil {
		logger.Warn("Input 'metricData' is empty. Nothing to convert.")
		return true, nil
	}

	// --- 2. Convert JSON to Prometheus Metric Format ---
	logger.Debugf("Input metric data: %+v", input.MetricData)
	logger.Debugf("Processing %d fields in metric data", len(input.MetricData))

	prometheusMetric, err := a.convertToPrometheusFormat(input.MetricData)
	if err != nil {
		logger.Errorf("Failed to convert JSON to Prometheus format: %v", err)
		return false, err
	}

	// Debug: Check the output before setting
	logger.Debugf("DEBUG: prometheusMetric before setting: '%s'", prometheusMetric)
	logger.Debugf("DEBUG: prometheusMetric contains newlines: %t", strings.Contains(prometheusMetric, "\n"))
	logger.Debugf("DEBUG: prometheusMetric line count: %d", strings.Count(prometheusMetric, "\n")+1)

	// Generate formatted version (multi-line) from single-line output
	prometheusMetricFormatted := a.formatForLogging(prometheusMetric)

	logger.Debugf("Generated prometheus metric output: %s", prometheusMetricFormatted)

	// --- 3. Set Output ---
	output := &Output{
		PrometheusMetric: prometheusMetricFormatted,
	}

	err = ctx.SetOutputObject(output)
	if err != nil {
		logger.Errorf("Error setting output object: %v", err)
		return false, err
	}

	logger.Debugf("Successfully generated Prometheus metrics. Output length: %d, lines: %d",
		len(prometheusMetric), strings.Count(prometheusMetric, "\n")+1)

	// Additional debug logging to verify the output contains all lines
	lines := strings.Split(prometheusMetric, "\n")
	logger.Debugf("Final output contains %d lines", len(lines))
	for i, line := range lines {
		if strings.TrimSpace(line) != "" {
			logger.Debugf("Output line %d: %s", i+1, line)
		}
	}

	logger.Debugf("Successfully converted JSON to Prometheus metric: %s", prometheusMetric)
	return true, nil
}

// convertToPrometheusFormat converts JSON data to Prometheus metric format
func (a *Activity) convertToPrometheusFormat(data map[string]interface{}) (string, error) {
	var parts []string

	// Add HELP comment if enabled
	if a.includeHelp {
		help := "Generated metric from JSON data"
		if helpValue, ok := data["help"]; ok {
			if helpStr, err := coerce.ToString(helpValue); err == nil {
				help = helpStr
			}
		}
		parts = append(parts, fmt.Sprintf("# HELP %s %s", a.metricName, help))
	}

	// Add TYPE comment if enabled
	if a.includeType {
		parts = append(parts, fmt.Sprintf("# TYPE %s %s", a.metricName, a.metricType))
	}

	// Check if data contains an array of metrics
	if metricsArray, ok := data["metrics"]; ok {
		// Handle array of metric objects
		if metrics, err := coerce.ToArray(metricsArray); err == nil {
			for _, metricItem := range metrics {
				if metricObj, err := coerce.ToObject(metricItem); err == nil {
					metricLine, err := a.processMetricObject(metricObj)
					if err != nil {
						continue // Skip invalid metric objects
					}
					parts = append(parts, metricLine)
				}
			}
		}
	} else {
		// Handle single metric object (backward compatibility)
		metricLine, err := a.processMetricObject(data)
		if err != nil {
			return "", err
		}
		parts = append(parts, metricLine)
	}

	// Join all parts with spaces to create a single line
	finalOutput := strings.Join(parts, " ")

	return finalOutput, nil
}

// formatForLogging converts single-line format back to multi-line for better log readability
func (a *Activity) formatForLogging(singleLineOutput string) string {
	// Split the single line back into logical components
	// Look for patterns: "# HELP", "# TYPE", and metric lines

	formatted := singleLineOutput

	// Add newlines after HELP comment
	formatted = strings.ReplaceAll(formatted, "# HELP ", "\n# HELP ")

	// Add newlines after TYPE comment
	formatted = strings.ReplaceAll(formatted, "# TYPE ", "\n# TYPE ")

	// Add newlines before each metric line (looking for metric_name pattern)
	// This regex finds the start of each new metric line
	parts := strings.Split(formatted, a.metricName+"{")
	if len(parts) > 1 {
		// Rejoin with newlines between metric lines
		result := parts[0] + a.metricName + "{"
		for i := 1; i < len(parts); i++ {
			if i == 1 {
				result += parts[i]
			} else {
				result += "\n" + a.metricName + "{" + parts[i]
			}
		}
		formatted = result
	}

	// Clean up any leading newlines
	formatted = strings.TrimPrefix(formatted, "\n")

	return formatted
}

// processMetricObject processes a single metric object and returns a Prometheus metric line
func (a *Activity) processMetricObject(metricObj map[string]interface{}) (string, error) {
	// Get timestamp once if enabled
	var timestamp int64
	if a.timestamp {
		timestamp = time.Now().UnixMilli()
		if timestampValue, ok := metricObj["timestamp"]; ok {
			if ts, err := coerce.ToInt64(timestampValue); err == nil {
				timestamp = ts
			} else if tsStr, err := coerce.ToString(timestampValue); err == nil {
				if t, err := time.Parse(time.RFC3339, tsStr); err == nil {
					timestamp = t.UnixMilli()
				}
			}
		}
	}

	// Process each numeric field as a separate metric
	var metricLines []string

	// Get all keys and sort them for consistent output
	keys := make([]string, 0, len(metricObj))
	for k := range metricObj {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Extract labels (all non-numeric, non-reserved fields)
	labels := a.extractLabelsFromObject(metricObj)

	// Process each numeric field
	numericFieldsFound := 0
	for _, key := range keys {
		val := metricObj[key]

		// Skip reserved fields
		if a.isReservedField(key) {
			continue
		}

		// Check if this field is numeric
		var value string
		if floatVal, err := coerce.ToFloat64(val); err == nil {
			value = strconv.FormatFloat(floatVal, 'f', -1, 64)
		} else if intVal, err := coerce.ToInt64(val); err == nil {
			value = strconv.FormatInt(intVal, 10)
		} else if strVal, err := coerce.ToString(val); err == nil {
			// Try to parse string as number
			if floatVal, err := strconv.ParseFloat(strVal, 64); err == nil {
				value = strconv.FormatFloat(floatVal, 'f', -1, 64)
			} else if intVal, err := strconv.ParseInt(strVal, 10, 64); err == nil {
				value = strconv.FormatInt(intVal, 10)
			} else {
				// Not a numeric value, skip
				continue
			}
		} else {
			// Not a numeric value, skip
			continue
		}

		numericFieldsFound++

		// Build metric line
		metricLine := a.metricName

		// Add name label to distinguish different metrics (use "name" instead of "metric_name")
		allLabels := fmt.Sprintf(`name="%s"`, a.sanitizeLabelValue(key))
		if len(labels) > 0 {
			allLabels += "," + labels
		}

		if len(allLabels) > 0 {
			metricLine += "{" + allLabels + "}"
		}
		metricLine += " " + value // Add timestamp if enabled
		if a.timestamp {
			metricLine += " " + strconv.FormatInt(timestamp, 10)
		}

		metricLines = append(metricLines, metricLine)
	}

	if len(metricLines) == 0 {
		// Create a list of available keys for debugging
		var availableKeys []string
		for k := range metricObj {
			availableKeys = append(availableKeys, k)
		}
		return "", fmt.Errorf("no numeric value found in metric object. Available fields: %v", availableKeys)
	}

	// Join all metric lines into a single line separated by spaces instead of newlines
	finalResult := strings.Join(metricLines, " ")

	// Ensure we return properly formatted string as single line
	return finalResult, nil
}

// extractLabelsFromObject extracts and formats labels from a metric object
func (a *Activity) extractLabelsFromObject(metricObj map[string]interface{}) string {
	var labelPairs []string

	// Get all keys and sort them for consistent output
	keys := make([]string, 0, len(metricObj))
	for k := range metricObj {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		val := metricObj[key]

		// Skip reserved fields
		if a.isReservedField(key) {
			continue
		}

		// Only include non-numeric string values as labels
		if strVal, err := coerce.ToString(val); err == nil {
			// Check if this value can be parsed as a number
			isNumeric := false
			if _, err := strconv.ParseFloat(strVal, 64); err == nil {
				isNumeric = true
			} else if _, err := strconv.ParseInt(strVal, 10, 64); err == nil {
				isNumeric = true
			} else if _, err := coerce.ToFloat64(val); err == nil {
				isNumeric = true
			} else if _, err := coerce.ToInt64(val); err == nil {
				isNumeric = true
			}

			// Only add as label if it's not numeric
			if !isNumeric {
				labelPairs = append(labelPairs, fmt.Sprintf("%s=\"%s\"", a.sanitizeLabelName(key), a.sanitizeLabelValue(strVal)))
			}
		}
	}

	return strings.Join(labelPairs, ",")
}

// sanitizeLabelValue escapes special characters in label values
func (a *Activity) sanitizeLabelValue(value string) string {
	// Escape quotes and backslashes in label values
	escapedVal := strings.ReplaceAll(value, "\\", "\\\\")
	escapedVal = strings.ReplaceAll(escapedVal, "\"", "\\\"")
	return escapedVal
}

// isReservedField checks if a field is reserved and should not be used as a label
func (a *Activity) isReservedField(key string) bool {
	reservedFields := map[string]bool{
		"help":      true,
		"timestamp": true,
		"type":      true,
	}
	return reservedFields[strings.ToLower(key)]
}

// sanitizeLabelName ensures label names conform to Prometheus requirements
func (a *Activity) sanitizeLabelName(name string) string {
	// Prometheus label names must match [a-zA-Z_][a-zA-Z0-9_]*
	var result strings.Builder

	for i, r := range name {
		if i == 0 {
			// First character must be letter or underscore
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_' {
				result.WriteRune(r)
			} else {
				result.WriteRune('_')
			}
		} else {
			// Subsequent characters can be letters, digits, or underscores
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
				result.WriteRune(r)
			} else {
				result.WriteRune('_')
			}
		}
	}

	return result.String()
}

// --- Supporting Structs ---

type Settings struct {
	MetricType  string `md:"metricType"`
	MetricName  string `md:"metricName"`
	IncludeHelp bool   `md:"includeHelp"`
	IncludeType bool   `md:"includeType"`
	Timestamp   bool   `md:"timestamp"`
}

// FromMap populates the struct from a map.
func (s *Settings) FromMap(values map[string]interface{}) error {
	if values == nil {
		s.MetricType = "gauge"
		s.MetricName = "flogo_metric"
		s.IncludeHelp = true
		s.IncludeType = true
		s.Timestamp = false
		return nil
	}

	var err error

	if val, ok := values[sMetricType]; ok && val != nil {
		s.MetricType, err = coerce.ToString(val)
		if err != nil {
			return err
		}
	}
	if s.MetricType == "" {
		s.MetricType = "gauge"
	}

	if val, ok := values[sMetricName]; ok && val != nil {
		s.MetricName, err = coerce.ToString(val)
		if err != nil {
			return err
		}
	}
	if s.MetricName == "" {
		s.MetricName = "flogo_metric"
	}

	if val, ok := values[sIncludeHelp]; ok && val != nil {
		s.IncludeHelp, err = coerce.ToBool(val)
		if err != nil {
			s.IncludeHelp = true // Default on error
		}
	} else {
		s.IncludeHelp = true // Default if not present
	}

	if val, ok := values[sIncludeType]; ok && val != nil {
		s.IncludeType, err = coerce.ToBool(val)
		if err != nil {
			s.IncludeType = true // Default on error
		}
	} else {
		s.IncludeType = true // Default if not present
	}

	if val, ok := values[sTimestamp]; ok && val != nil {
		s.Timestamp, err = coerce.ToBool(val)
		if err != nil {
			s.Timestamp = false // Default on error
		}
	} else {
		s.Timestamp = false // Default if not present
	}

	return nil
}

type Input struct {
	MetricData map[string]interface{} `md:"metricData"`
}

// FromMap populates the struct from the activity's inputs.
func (i *Input) FromMap(values map[string]interface{}) error {
	if values == nil {
		return nil
	}

	var err error
	i.MetricData, err = coerce.ToObject(values[ivMetricData])
	if err != nil {
		return err
	}
	return nil
}

// ToMap converts the struct to a map.
func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		ivMetricData: i.MetricData,
	}
}

type Output struct {
	PrometheusMetric string `md:"prometheusMetric"`
}

// ToMap converts the struct to a map.
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"prometheusMetric": o.PrometheusMetric,
	}
}

// FromMap populates the struct from a map.
func (o *Output) FromMap(values map[string]interface{}) error {
	if values == nil {
		return nil
	}

	var err error
	if val, ok := values["prometheusMetric"]; ok && val != nil {
		o.PrometheusMetric, err = coerce.ToString(val)
		if err != nil {
			return err
		}
	}
	return nil
}
