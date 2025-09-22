---
layout: default
title: TIBCO Flogo Extensions Collection
---

# TIBCO Flogo Extensions Collection

A comprehensive collection of custom TIBCO Flogo activities for enterprise integration and AI-powered workflows.

## 🚀 Available Extensions

### 1. [Pongo2 AI Prompt Engine](./pongo2/docs/)
**The world's first enterprise workflow platform with native AI prompt engineering capabilities**

- 🤖 **AI-Powered Templates**: Jinja2-like templating with advanced logic
- 📊 **JSON Schema Generation**: Auto-generate Flogo Web UI parameters
- 🔄 **Dynamic Content**: Variables, loops, conditionals, and filters
- 🎯 **Enterprise Ready**: Built for production workflows

**Key Features:**
- Template variable extraction and JSON schema generation
- Support for complex data structures and arrays
- Built-in template examples and troubleshooting guides
- Visual prompt engineering for business users

[📖 Documentation](./pongo2/docs/) | [🧪 Examples](./pongo2/docs/PONGO2_EXAMPLES.md) | [🔧 Build Guide](./pongo2/docs/BUILD_TEST_GUIDE.md)

---

### 2. [Pongo2 V2 (Enhanced)](./pongo2v2/)
**Next-generation AI prompt engine with advanced features**

- ⚡ **Performance Optimized**: Faster template processing
- 🛡️ **Enhanced Security**: Improved input validation
- 📈 **Scalability**: Better memory management for large templates
- 🔌 **Extended API**: Additional template functions and filters

*Coming Soon - Enhanced version with improved performance and features*

---

### 3. [Prime Numbers Generator](./prime-numbers/)
**Mathematical utility for generating prime numbers**

- 🔢 **Efficient Algorithms**: Optimized prime number generation
- 📊 **Flexible Output**: Various output formats supported
- 🎯 **Range Support**: Generate primes within specified ranges
- 🧮 **Mathematical Operations**: Prime validation and testing

*Perfect for cryptographic applications and mathematical workflows*

---

### 4. [Prometheus Metrics Activity](./prometheus-metrics/)
**Convert JSON data into Prometheus metrics format**

- 📈 **Auto-Conversion**: JSON to Prometheus format transformation
- 🏷️ **Smart Labeling**: Automatic label generation from data
- 📊 **Multiple Metrics**: Support for gauge, counter, histogram, summary
- 🔍 **Observability**: Real-time monitoring and alerting ready

**Key Capabilities:**
- Converts any JSON input to Prometheus metrics
- Kubernetes ServiceMonitor compatible
- Built-in `/metrics` endpoint support
- Business KPI monitoring from data streams

[🐳 Docker Image](https://hub.docker.com/r/kulbhushanbhalerao/flogo-prometheus-metrics-activity) | [📋 Usage Examples](./prometheus-metrics/samples/)

## 🛠️ Quick Start

### Installation
Each extension can be imported into your Flogo application using the following reference format:

```json
{
  "ref": "github.com/KulbhushanBhalerao/flogo-extensions/[extension-name]",
  "version": "v1.0.0"
}
```

### Build from Source
```bash
# Clone the repository
git clone https://github.com/KulbhushanBhalerao/flogo-extensions.git

# Navigate to specific extension
cd flogo-extensions/[extension-name]

# Build the extension
go build .

# Run tests
go test -v
```

## 📚 Documentation

- **[Pongo2 Documentation](./pongo2/docs/)** - Comprehensive guide for AI prompt engineering
- **[JSON Schema Guide](./pongo2/docs/JSON_SCHEMA_GUIDE.md)** - Generate Flogo Web UI parameters
- **[Template Examples](./pongo2/docs/PONGO2_EXAMPLES.md)** - Real-world template samples
- **[Troubleshooting](./pongo2/docs/TROUBLESHOOTING.md)** - Common issues and solutions
- **[Build & Test Guide](./pongo2/docs/BUILD_TEST_GUIDE.md)** - Development workflow

## 🎯 Use Cases

### Enterprise AI Workflows
- **Dynamic Prompt Generation**: Create AI prompts from business data
- **Template-Driven Reports**: Generate consistent business reports
- **Data Analysis Queries**: Build analytical prompts from datasets
- **Customer Communication**: Personalized messaging at scale

### Monitoring & Observability
- **Business Metrics**: Convert business KPIs to Prometheus format
- **Real-time Dashboards**: Feed data to Grafana and monitoring tools
- **Alert Management**: Custom alerting based on business logic
- **Performance Tracking**: Monitor application and business metrics

### Mathematical Operations
- **Cryptographic Applications**: Prime number generation for security
- **Data Science**: Mathematical operations in data pipelines
- **Algorithm Testing**: Prime number validation and testing

## 🤝 Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the individual extension directories for specific license files.

## 🔗 Links

- **GitHub Repository**: [flogo-extensions](https://github.com/KulbhushanBhalerao/flogo-extensions)
- **Docker Hub**: [Prometheus Metrics Activity](https://hub.docker.com/r/kulbhushanbhalerao/flogo-prometheus-metrics-activity)
- **TIBCO Flogo**: [Official Documentation](https://docs.tibco.com/products/tibco-flogo-enterprise)

## 📞 Support

For questions, issues, or feature requests:
- Create an issue on GitHub
- Check the troubleshooting guides in each extension
- Review the comprehensive documentation

---

*Built with ❤️ for the TIBCO Flogo community*