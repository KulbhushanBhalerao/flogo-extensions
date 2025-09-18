# Using Pongo2 Prompt Activity in Flogo Flows

## Overview
This guide shows how to integrate the pongo2-prompt activity into your Flogo flows for dynamic AI prompt generation using pongo2 templates (Django/Jinja2-like syntax).

## Installation in Flogo Application

### 1. Add Activity to Your Flogo App

#### Option A: Using Flogo CLI
```bash
# Navigate to your Flogo app directory
cd my-flogo-app

# Install the activity (if published as module)
```bash
flogo install github.com/your-org/pongo2-prompt 
```

# Or install from local path
flogo install /opt/tibco/flogo-extensions/pongo2-prompt
```

#### Option B: Manual Installation
1. Copy the activity files to your app's `src` directory
2. Update your `flogo.json` imports:
```json
{
  "imports": [
    "github.com/project-flogo/contrib/activity/log",
    "path/to/pongo2-prompt"
  ]
}
```

### 2. Register Activity in Application
Add to your app's `activity` imports:
```json
{
  "name": "my-app",
  "type": "flogo:app", 
  "version": "1.0.0",
  "imports": [
    "github.com/project-flogo/contrib/activity/log",
    "#pongo2-prompt"
  ]
}
```

## Basic Flow Configuration

### Simple AI Prompt Generation Flow

```json
{
  "name": "SimplePromptFlow",
  "description": "Generate AI prompts using templates",
  "tasks": [
    {
      "id": "generate_prompt",
      "name": "Generate AI Prompt",
      "activity": {
        "ref": "#pongo2-prompt",
        "settings": {
          "template": "You are a {{ role }}. Please {{ task }} about {{ topic }}."
        },
        "input": {
          "role": "helpful assistant",
          "task": "provide information", 
          "topic": "machine learning"
        }
      }
    },
    {
      "id": "log_prompt",
      "name": "Log Generated Prompt",
      "activity": {
        "ref": "github.com/project-flogo/contrib/activity/log",
        "input": {
          "message": "=$.generate_prompt.renderedPrompt"
        }
      }
    }
  ],
  "links": [
    {"from": "generate_prompt", "to": "log_prompt"}
  ]
}
```

## Advanced Flow Patterns

### 1. Dynamic Prompt with REST Input

```json
{
  "triggers": [
    {
      "id": "rest_trigger",
      "ref": "github.com/project-flogo/contrib/trigger/rest",
      "settings": {
        "port": 8080
      },
      "handlers": [
        {
          "settings": {
            "method": "POST",
            "path": "/generate-prompt"
          },
          "action": {
            "ref": "github.com/project-flogo/flow",
            "settings": {
              "flowURI": "res://flow:dynamic_prompt"
            }
          }
        }
      ]
    }
  ],
  "resources": [
    {
      "id": "flow:dynamic_prompt",
      "data": {
        "name": "DynamicPromptFlow",
        "tasks": [
          {
            "id": "transform_input",
            "name": "Transform Request Data",
            "activity": {
              "ref": "github.com/project-flogo/contrib/activity/mapper",
              "input": {
                "mapping": {
                  "user_role": "=$trigger.content.role",
                  "analysis_type": "=$trigger.content.analysis",
                  "dataset_info": "=$trigger.content.data"
                }
              }
            }
          },
          {
            "id": "generate_analysis_prompt",
            "name": "Generate Analysis Prompt",
            "activity": {
              "ref": "#pongo2-prompt",
              "settings": {
                "template": "You are a {{ user_role }}. Analyze the {{ dataset_info.type }} dataset with {{ dataset_info.records }} records. Focus on {{ analysis_type }} analysis."
              },
              "input": {
                "variables": "=$.transform_input.mapping"
              }
            }
          },
          {
            "id": "return_response",
            "name": "Return Prompt",
            "activity": {
              "ref": "github.com/project-flogo/flow/activity/reply",
              "input": {
                "code": 200,
                "data": {
                  "generated_prompt": "=$.generate_analysis_prompt.renderedPrompt",
                  "timestamp": "=$activity[generate_analysis_prompt].timestamp"
                }
              }
            }
          }
        ],
        "links": [
          {"from": "transform_input", "to": "generate_analysis_prompt"},
          {"from": "generate_analysis_prompt", "to": "return_response"}
        ]
      }
    }
  ]
}
```

### 2. Multi-Stage AI Pipeline

```json
{
  "name": "AIProcessingPipeline",
  "description": "Multi-stage AI processing with dynamic prompts",
  "tasks": [
    {
      "id": "preprocessing_prompt",
      "name": "Generate Preprocessing Prompt", 
      "activity": {
        "ref": "#pongo2-prompt",
        "settings": {
          "template": "Preprocess the following {{ data_type }} data: {{ raw_data }}. Focus on {{ preprocessing_steps }}."
        },
        "input": {
          "data_type": "=$.input.dataType",
          "raw_data": "=$.input.data",
          "preprocessing_steps": "cleaning and normalization"
        }
      }
    },
    {
      "id": "analysis_prompt", 
      "name": "Generate Analysis Prompt",
      "activity": {
        "ref": "#pongo2-prompt",
        "settings": {
          "template": "Analyze the preprocessed data using {{ analysis_method }}. Generate insights about {{ focus_areas }}."
        },
        "input": {
          "analysis_method": "=$.input.method",
          "focus_areas": "=$.input.objectives"
        }
      }
    },
    {
      "id": "reporting_prompt",
      "name": "Generate Reporting Prompt",
      "activity": {
        "ref": "#pongo2-prompt", 
        "settings": {
          "template": "Create a {{ report_type }} report summarizing: {{ analysis_results }}. Target audience: {{ audience }}."
        },
        "input": {
          "report_type": "executive",
          "analysis_results": "=$.analysis_prompt.renderedPrompt", 
          "audience": "=$.input.audience"
        }
      }
    }
  ],
  "links": [
    {"from": "preprocessing_prompt", "to": "analysis_prompt"},
    {"from": "analysis_prompt", "to": "reporting_prompt"}
  ]
}
```

## Integration with External Services

### 1. OpenAI Integration Flow

```json
{
  "name": "OpenAIIntegrationFlow",
  "tasks": [
    {
      "id": "generate_openai_prompt",
      "name": "Generate OpenAI Prompt",
      "activity": {
        "ref": "#pongo2-prompt",
        "settings": {
          "template": "{{ system_message }}\n\nUser: {{ user_query }}\n\nContext: {{ context_info }}\n\nPlease respond in {{ response_format }} format."
        },
        "input": {
          "system_message": "You are a helpful AI assistant specialized in data analysis.",
          "user_query": "=$.input.question", 
          "context_info": "=$.input.context",
          "response_format": "JSON"
        }
      }
    },
    {
      "id": "call_openai",
      "name": "Call OpenAI API",
      "activity": {
        "ref": "github.com/project-flogo/contrib/activity/rest",
        "input": {
          "method": "POST",
          "uri": "https://api.openai.com/v1/chat/completions",
          "headers": {
            "Authorization": "Bearer YOUR_API_KEY",
            "Content-Type": "application/json"
          },
          "content": {
            "model": "gpt-4",
            "messages": [
              {
                "role": "user",
                "content": "=$.generate_openai_prompt.renderedPrompt"
              }
            ],
            "max_tokens": 1000
          }
        }
      }
    }
  ],
  "links": [
    {"from": "generate_openai_prompt", "to": "call_openai"}
  ]
}
```

## Best Practices

### 1. Template Organization
- **Keep templates modular** - Use separate activities for different prompt types
- **Version templates** - Store templates in configuration or external files
- **Test templates** - Validate template syntax before deployment

### 2. Error Handling
```json
{
  "id": "generate_prompt_with_error_handling",
  "name": "Generate Prompt (Safe)",
  "activity": {
    "ref": "#pongo2-prompt",
    "settings": {
      "template": "{{ prompt_template | default:'Default prompt when template is missing' }}"
    }
  },
  "errorHandler": {
    "tasks": [
      {
        "id": "log_error",
        "activity": {
          "ref": "github.com/project-flogo/contrib/activity/log",
          "input": {
            "message": "Template rendering failed: =$.error.message"
          }
        }
      }
    ]
  }
}
```

### 3. Performance Optimization
- **Cache templates** - Reuse activity instances when possible
- **Minimize variable complexity** - Pre-process complex data structures
- **Use async patterns** - Don't block on template generation

### 4. Security Considerations
- **Sanitize inputs** - Validate user inputs before template processing
- **Template injection** - Use allowlists for dynamic template parts
- **Sensitive data** - Avoid logging rendered prompts containing PII

## Deployment Patterns

### 1. Microservice Pattern
Deploy as dedicated prompt generation service:
```json
{
  "triggers": [
    {
      "ref": "github.com/project-flogo/contrib/trigger/rest",
      "settings": {"port": 8080}
    }
  ],
  "resources": [
    {
      "id": "flow:prompt_service",
      "data": {
        "name": "PromptGenerationService",
        "description": "Centralized prompt generation microservice"
      }
    }
  ]
}
```

### 2. Event-Driven Pattern  
Use with message queues:
```json
{
  "triggers": [
    {
      "ref": "github.com/project-flogo/contrib/trigger/kafka",
      "settings": {
        "brokers": ["localhost:9092"],
        "topic": "prompt_requests"
      }
    }
  ]
}
```

### 3. Batch Processing Pattern
Process multiple templates:
```json
{
  "tasks": [
    {
      "id": "iterate_requests",
      "activity": {
        "ref": "github.com/project-flogo/contrib/activity/iterator",
        "input": {
          "collection": "=$.input.prompt_requests"
        },
        "iterator": {
          "tasks": [
            {
              "id": "generate_individual_prompt",
              "activity": {
                "ref": "#pongo2-prompt"
              }
            }
          ]
        }
      }
    }
  ]
}
```

## Monitoring and Debugging

### 1. Add Logging
```json
{
  "tasks": [
    {
      "id": "log_template_input",
      "activity": {
        "ref": "github.com/project-flogo/contrib/activity/log",
        "input": {
          "message": "Template variables: =$.template_variables"
        }
      }
    },
    {
      "id": "generate_prompt",
      "activity": {"ref": "#pongo2-prompt"}
    },
    {
      "id": "log_template_output", 
      "activity": {
        "ref": "github.com/project-flogo/contrib/activity/log",
        "input": {
          "message": "Generated prompt: =$.generate_prompt.renderedPrompt"
        }
      }
    }
  ]
}
```

### 2. Add Metrics
```json
{
  "tasks": [
    {
      "id": "record_template_metric",
      "activity": {
        "ref": "github.com/project-flogo/contrib/activity/counter",
        "input": {
          "counterName": "prompts_generated",
          "increment": 1
        }
      }
    }
  ]
}
```

This guide provides comprehensive coverage for integrating the pongo2-prompt activity into various Flogo flow patterns and deployment scenarios.