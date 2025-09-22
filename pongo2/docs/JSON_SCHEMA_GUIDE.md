---
layout: default
---

# JSON Schema Generation from Pongo2 Templates

## 🎯 Quick Start

Generate JSON schema for your Pongo2 template variables to use in Flogo Web UI.

### Method 1: Simplified Script (Recommended)

```bash
cd /opt/tibco/flogo-extensions/pongo2-prompt/utils
./generate_flogo_params.sh "Your template here"
```

**Output:**
```
🎯 JSON Schema (copy this to Flogo Web UI):
{"$schema":"http://json-schema.org/draft-04/schema#","properties":{...}}

📋 Variables detected: role, domain, objective
💡 For params type, add these parameters in Flogo Web UI:
  - role
  - domain  
  - objective
```

### Method 2: Go Direct

```bash
cd /opt/tibco/flogo-extensions/pongo2-prompt/utils  
go run schema_generator.go "Your template here"
```

### Method 3: Web Interface

Open `utils/schema_generator.html` in your browser.

## 📋 Example Usage

### Input Template:
```
You are a {{ "{{" }} role {{ "}}" }} working with {{ "{{" }} domain {{ "}}" }} data.

**Analysis Objective:** {{ "{{" }} objective {{ "}}" }}

{{ "{%" }} for item in data_list {{ "%}" }}
- {{ "{{" }} item.name {{ "}}" }}: {{ "{{" }} item.value {{ "}}" }}
{{ "{%" }} endfor {{ "%}" }}
```

### Output from Simplified Script:
```bash
🎯 JSON Schema (copy this to Flogo Web UI):
{"$schema":"http://json-schema.org/draft-04/schema#","properties":{"role":{"type":"string"},"domain":{"type":"string"},"objective":{"type":"string"},"data_list":{"type":"array","items":{"type":"object","additionalProperties":true}}}}

📋 Variables detected: role, domain, objective, data_list
💡 For params type, add these parameters in Flogo Web UI:
  - role
  - domain
  - objective
  - data_list
```

### Formatted JSON Schema (for reference):
```json
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "role": {"type": "string"},
    "domain": {"type": "string"}, 
    "objective": {"type": "string"},
    "data_list": {
      "type": "array",
      "items": {
        "type": "object",
        "additionalProperties": true
      }
    }
  }
}
```

## 🔧 How to Use in Flogo Web UI

### Option A: Object Type (Legacy)
1. Copy the JSON schema above
2. Paste into Template Variables field value in Flogo Web UI
3. Individual mappable fields will appear

### Option B: Params Type (Recommended) 
1. Use the detected variables list: `[role, domain, objective, data_list]`
2. In Flogo Web UI, click "Add Parameter" for each variable
3. Set parameter names exactly as shown

## 🚀 Features

- **Detects simple variables**: `{{ "{{" }} variable_name {{ "}}" }}`
- **Detects arrays in loops**: `{{ "{%" }} for item in array_name {{ "%}" }}`
- **Filters loop iterators**: Ignores `{{ "{{" }} item.property {{ "}}" }}` inside loops
- **Handles complex objects**: `{{ "{{" }} user.profile.name {{ "}}" }}` → base object `user`
- **JSON Schema compliant**: Works with Flogo Web UI schema system

## ⚠️ Requirements

- Go 1.16+ (auto-installed by script if missing)
- Access to `/opt/tibco/flogo-extensions/pongo2-prompt/utils/` directory