# TIBCO Flogo Prometheus Metrics Activity

A custom TIBCO Flogo activity that converts JSON messages into Prometheus metric format. This activity creates **multiple metrics** from a single JSON input, with each numeric field becoming a separate metric line with a `name` label.

## 🚀 Key Features

- **Multiple Metrics Generation**: Creates separate metrics for each numeric field in your JSON
- **Automatic Name Labels**: Each metric gets a `name` label with the JSON field name
- **Multiple Metric Types**: Support for gauge, counter, histogram, summary
- **Flexible Input**: Works with any JSON structure containing numeric values
- **Prometheus Compliant**: Properly formatted output ready for Prometheus ingestion
- **Label Sanitization**: Automatic sanitization of label names and values

### Installation

1. Install the activity using the Flogo CLI:
```bash
flogo install github.com/kulbhushanbhalerao/flogo-extensions/prometheus-metrics
```

2. Or add it to your `flogo.json` file:
```json
{
  "imports": [
    "github.com/kulbhushanbhalerao/flogo-extensions/prometheus-metrics"
  ]
}
```

## ⚙️ Configuration

### Settings

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| **Metric Type** | string | `gauge` | Type of Prometheus metric (gauge, counter, histogram, summary) |
| **Metric Name** | string | `flogo_metric` | Base name for all generated metrics |
| **Include HELP** | boolean | `true` | Include HELP comment in output |
| **Include TYPE** | boolean | `true` | Include TYPE comment in output |
| **Include Timestamp** | boolean | `false` | Include timestamp in metric output |

### Input

| Field | Type | Description |
|-------|------|-------------|
| **Metric Data** | object | JSON object containing numeric fields to convert to metrics |

## 💡 How It Works

The activity can handle two types of input:

### Single Metric Object
Pass a single JSON object where:
- **Numeric fields** become separate metrics with a `metric_name` label
- **String fields** become labels applied to all metrics
- **Reserved fields** (help, timestamp, type) are handled specially

### Multiple Metrics Array
Pass a JSON object with a `metrics` array where each array element is a metric object processed as above.

## 📋 Usage Examples

### Example 1: Single Metric Object with Multiple Values

**Input JSON:**
```json
{
  "cpu_usage": 75.2,
  "memory_usage": 68.5,
  "service": "user-service",
  "environment": "prod",
  "region": "us-east"
}
```

**Settings:**
- Metric Type: `gauge`
- Metric Name: `system_metrics`

**Output:**
```prometheus
# HELP system_metrics Generated metric from JSON data
# TYPE system_metrics gauge
system_metrics{metric_name="cpu_usage",environment="prod",region="us-east",service="user-service"} 75.2
system_metrics{metric_name="memory_usage",environment="prod",region="us-east",service="user-service"} 68.5
```

### Example 2: Multiple Metric Objects

**Input JSON:**
```json
{
  "metrics": [
    {
      "value": 45,
      "service": "my-service",
      "environment": "prod",
      "region": "us-east"
    },
    {
      "value": 67,
      "service": "my-service",
      "environment": "prod", 
      "region": "us-west"
    }
  ]
}
```

**Settings:**
- Metric Type: `counter`
- Metric Name: `request_count`

**Output:**
```prometheus
# HELP request_count Generated metric from JSON data
# TYPE request_count counter
request_count{metric_name="value",environment="prod",region="us-east",service="my-service"} 45
request_count{metric_name="value",environment="prod",region="us-west",service="my-service"} 67
```

### Example 3: Mixed Data Types

**Input JSON:**
```json
{
  "temperature": 23.5,
  "humidity": "65",
  "pressure": 1013,
  "location": "room_a",
  "sensor_id": "temp_001",
  "status": "online"
}
```

**Output:**
```prometheus
# HELP sensor_reading Generated metric from JSON data
# TYPE sensor_reading gauge
sensor_reading{metric_name="humidity",location="room_a",sensor_id="temp_001",status="online"} 65
sensor_reading{metric_name="pressure",location="room_a",sensor_id="temp_001",status="online"} 1013
sensor_reading{metric_name="temperature",location="room_a",sensor_id="temp_001",status="online"} 23.5
```

## Data Processing Rules

### Numeric Field Detection
The activity automatically detects numeric values from:
- Direct numeric types (int, float)
- String values that can be parsed as numbers (e.g., "45", "23.5")

### Label Creation  
- All non-numeric, non-reserved fields become labels
- Label names are sanitized for Prometheus compliance
- Label values are properly escaped
- Labels are sorted alphabetically for consistent output

### Reserved Fields
The following fields are treated specially:
- `help` - Used for HELP comment
- `timestamp` - Used for timestamp  
- `type` - Reserved
- Built-in metric fields are processed as metrics, not labels

## Integration Patterns

### Single Metrics
Use when you have one measurement with multiple dimensions:
```json
{
  "response_time": 150,
  "service": "api",
  "endpoint": "/users"
}
```

### Multiple Metrics
Use when you have multiple related measurements:
```json
{
  "metrics": [
    {"cpu": 75, "host": "server1"},
    {"cpu": 68, "host": "server2"}
  ]
}
```

**Output:**
```prometheus
## Error Handling

- **No numeric fields found**: Activity returns error with list of available fields
- **Invalid JSON input**: Gracefully handled with descriptive error messages  
- **Missing settings**: Use sensible defaults (gauge metrics, standard naming)
- **Invalid array format**: Individual invalid metrics are skipped, valid ones processed

## Output Format

The activity generates standard Prometheus exposition format:
```
# HELP <metric_name> <help_text>
# TYPE <metric_name> <metric_type>  
<metric_name>{<labels>} <value> [<timestamp>]
```

## Integration with Prometheus

The output can be used with:
- **Prometheus Pushgateway**: Push metrics directly
- **HTTP endpoint**: Serve metrics for scraping
- **File export**: Write to files for collection
- **Remote write**: Send to Prometheus-compatible systems

## Dependencies

- TIBCO Flogo Core v1.2.0+
- Go 1.19+
```

### Example 4: With Timestamps

**Input JSON:**
```json
{
    "response_time": 125.5,
    "requests_per_sec": 45,
    "timestamp": 1705316200000
}
```

**Settings:**
- Metric Type: `gauge`
- Metric Name: `web_metrics`
- Include Timestamp: `true`

**Output:**
```prometheus
# HELP web_metrics Generated metric from JSON data
# TYPE web_metrics gauge
web_metrics{name="requests_per_sec"} 45 1705316200000
web_metrics{name="response_time"} 125.5 1705316200000
```

## 🔧 Field Processing Rules

### Numeric Fields → Metrics
- **Integers**: `1`, `100`, `-50` → Become metric values
- **Floats**: `23.5`, `3.14159`, `0.001` → Become metric values
- **String Numbers**: `"42"`, `"3.14"` → Automatically converted to numeric values

### Non-Numeric Fields → Ignored
- **Strings**: `"hello"`, `"production"` → Skipped (not converted to metrics)
- **Booleans**: `true`, `false` → Skipped
- **Objects/Arrays**: `{}`, `[]` → Skipped

### Reserved Fields
These fields are treated specially and not converted to metrics:
- `help` → Used for HELP comment text
- `timestamp` → Used for timestamp value
- `type` → Reserved field

## 🏷️ Label Sanitization

Field names are automatically sanitized to comply with Prometheus label naming rules:
- Must start with letter or underscore: `123field` → `_23field`
- Only letters, numbers, underscores allowed: `field-name` → `field_name`
- Spaces replaced with underscores: `field name` → `field_name`

## 🔗 Integration Examples

### With Prometheus Pushgateway
```go
// Use the activity output to push to Pushgateway
prometheusOutput := activityOutput.PrometheusMetric
http.Post("http://pushgateway:9091/metrics/job/flogo-app", "text/plain", strings.NewReader(prometheusOutput))
```

### With File Export
```go
// Write metrics to file for Prometheus file_sd_config
os.WriteFile("/var/lib/prometheus/flogo-metrics.prom", []byte(activityOutput.PrometheusMetric), 0644)
```

### With HTTP Endpoint
```go
// Serve metrics via HTTP endpoint
http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    w.Write([]byte(activityOutput.PrometheusMetric))
})
```

## 🛠️ Building and Deployment

### Quick Deployment (Recommended)
```bash
# Complete build and deployment guidance
./deploy.sh
```

This script will:
1. Build the activity
2. Clean Go module cache
3. Check for Flogo applications that need rebuilding
4. Provide specific rebuild instructions

### Quick Build
```bash
# Build just the activity
./build.sh
```

### Manual Build Steps
```bash
# Download dependencies
go mod tidy

# Build the activity
go build -v .

# Run tests (optional)
go test -v
```

### Flogo Application Integration

#### Step 1: Update Your Activity Code
If you made changes to the activity, ensure they're saved and built:

```bash
cd /path/to/prometheus-metrics
go build
```

#### Step 2: Clear Go Module Cache
To ensure Flogo picks up the latest activity version:

```bash
go clean -modcache
```

#### Step 3: Rebuild Flogo Application
Depending on your setup:

**If using Flogo CLI:**
```bash
# Navigate to your app directory
cd /path/to/your/flogo/app

# Rebuild the application
flogo build -f your-app.flogo
```

**If using development environment:**
- Stop the current Flogo application completely
- Clear any build cache in your development environment
- Restart/rebuild the application

**If using a compiled binary:**
```bash
# Remove the old binary
rm bin/your-flogo-app

# Rebuild using your development environment or CI/CD pipeline
```

#### Step 4: Restart Application
```bash
# If using the compiled binary
./bin/your-flogo-app

# If using flogo run
flogo run
```

### Troubleshooting Build Issues

#### Activity Not Recognized
**Problem**: `failed to resolve activity attr: 'prometheusMetric', not found`

**Solution**:
1. Ensure the activity is properly built: `go build`
2. Clear module cache: `go clean -modcache`
3. Rebuild the entire Flogo application
4. Restart the application completely

#### Module Import Issues
**Problem**: Import path not found

**Solution**:
1. Verify the module name in `go.mod` matches the import in your `.flogo` file
2. Ensure the activity directory is accessible to the Flogo build process
3. Check that all dependencies are available: `go mod tidy`

#### Build Cache Issues
**Problem**: Changes not reflected in running application

**Solution**:
```bash
# Clean everything
go clean -cache -modcache
cd prometheus-metrics && go build
# Then rebuild your Flogo application
```

### Development Workflow

1. **Make Changes**: Edit `activity.go` or other source files
2. **Test**: Run `go test -v` to verify changes
3. **Build Activity**: Run `./build.sh` or `go build`
4. **Clean Cache**: Run `go clean -modcache`
5. **Rebuild App**: Rebuild your Flogo application
6. **Restart**: Restart the Flogo application to pick up changes

### Verification

After rebuilding and restarting, verify the activity is working:

1. **Check Logs**: Look for DEBUG statements from the activity
2. **Validate Output**: Ensure metrics are generated correctly
3. **Test Flow**: Run your flow and check the output format

Example expected log output:
```
DEBUG: Processing single metric object
DEBUG: Processing field 'cpu_usage' with value '75'
DEBUG: Generated labels: 'metric_name="cpu_usage",environment="prod",service="web-server"'
DEBUG: Final result from processMetricObject: '...'
```

## 🔧 Activity Configuration

### In Flogo Flow JSON
```json
{
  "activity": {
    "ref": "github.com/kulbhushanbhalerao/flogo-activities/prometheus-metrics",
    "settings": {
      "metricType": "gauge",
      "metricName": "my_metric",
      "includeHelp": true,
      "includeType": true,
      "timestamp": false
    },
    "input": {
      "metricData": "=$flow.body"
    }
  }
}
```

## 📊 Real-World Use Cases

1. **IoT Sensor Data**: Convert sensor readings (temperature, humidity, pressure) into individual metrics
2. **Application Performance**: Monitor different performance counters (CPU, memory, network) as separate metrics
3. **Business Metrics**: Track various KPIs (sales, users, transactions) from a single JSON payload
4. **System Monitoring**: Convert system stats into multiple Prometheus metrics for dashboarding

## 🚨 Error Handling

- **No Numeric Fields**: Activity succeeds but produces only HELP/TYPE comments (no metric lines)
- **Invalid JSON**: Activity returns error with descriptive message
- **Invalid Settings**: Uses sensible defaults

## 📋 Requirements

- TIBCO Flogo Core v1.2.0+
- Go 1.19+

## � Quick Reference

### Build and Deploy Activity
```bash
./deploy.sh          # Complete deployment with guidance
./build.sh           # Build activity only
```

### Module Information
```
Module Path: github.com/kulbhushanbhalerao/flogo-extensions/prometheus-metrics
Activity ID: prometheus-metrics
Reference:   #prometheus-metrics
```

### Common Commands
```bash
# Clean build everything
go clean -cache -modcache && go build

# Test the activity
go test -v

# Remove old Flogo binary (replace with your app name)
rm ../../bin/your-flogo-app

# Check activity is working (look for these in logs)
"DEBUG: Processing field 'field_name' with value 'value'"
"DEBUG: Number of metric lines generated: X"
```

### Expected Output Format
```prometheus
# HELP metric_name Generated metric from JSON data
# TYPE metric_name gauge
metric_name{metric_name="field1",label1="value1"} 123
metric_name{metric_name="field2",label1="value1"} 456
```

### Troubleshooting
- **Field not found**: Rebuild Flogo application completely
- **No DEBUG logs**: Activity not updated, clean cache and rebuild
- **Single metric only**: Check input JSON has multiple numeric fields

## �📄 License

This activity follows the same license as your TIBCO Flogo installation.
