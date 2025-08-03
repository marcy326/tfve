# tivor - Terraform Infrastructure Variable Orchestrator

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.22-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

**tivor** is an orchestration tool for intuitive, safe, and declarative lifecycle management across multiple Terraform environments. It adopts a GitOps-First approach and manages all environment configurations from a single `tivor.yaml` configuration file.

## ✨ Key Features

- 🏗️ **Declarative Configuration** - Define all environments in a single `tivor.yaml` file
- 🔄 **Environment Inheritance** - Share and override configurations across environments  
- 🔀 **Smart Variable Merging** - Automatically merge and override variable files
- 🎯 **GitOps-First** - Git as the single source of truth for all configurations
- 🔐 **Secure by Default** - Built-in SOPS integration for secret management
- 🚀 **Multiple Backends** - Support for local, S3, and other storage backends
- 🛡️ **Safe Execution** - Validation and temporary file management

## 🚀 Quick Start

### Installation

```bash
# Build from source
git clone https://github.com/marcy326/tivor
cd tivor
go build -o bin/tivor ./cmd/tivor
```

### Basic Usage

1. **Initialize a configuration file:**
   ```bash
   tivor init
   ```

2. **Plan for development environment:**
   ```bash
   tivor plan dev --working-dir=./infrastructure
   ```

3. **Apply to staging environment:**
   ```bash
   tivor apply staging --working-dir=./infrastructure
   ```

## 📁 Project Structure

```
tivor/
├── cmd/                     # CLI entry point
├── internal/
│   ├── cli/                # Cobra commands (plan, apply, init, etc.)
│   ├── config/            # Configuration file parsing and validation
│   ├── backend/           # Storage backends (local, s3, etc.)
│   ├── terraform/         # Terraform execution engine
│   └── tfvars/           # Variable file parsing and merging
├── examples/              # Usage examples and integration tests
│   └── basic/            # Basic usage example with inheritance
└── configs/              # Configuration file templates
```

## ⚙️ Configuration

### Sample `tivor.yaml`

```yaml
# Configuration file version
version: "1.0"

# Default settings for all environments
defaults:
  vars_files:
    - "variables/common.tfvars"

# Secret management
secrets:
  engine: sops
  sops_config_path: ".sops.yaml"

# Environment definitions
environments:
  # Development environment
  - name: dev
    vars_files:
      - "variables/dev.tfvars"
    backend:
      type: local

  # Staging inherits from dev
  - name: staging
    inherits: dev
    vars_files:
      - "variables/staging.tfvars"
    backend:
      type: local

  # Production inherits from staging
  - name: production
    inherits: staging
    vars_files:
      - "variables/production.tfvars"
    backend:
      type: s3
      config:
        bucket: "terraform-state-bucket"
        key: "terraform.tfstate"
        region: "us-west-2"
```

### Environment Inheritance

tivor supports powerful environment inheritance patterns:

```
common.tfvars → dev.tfvars → staging.tfvars → production.tfvars
```

Later files override earlier ones, allowing you to:
- 📄 Define common settings once
- 🔧 Override specific values per environment  
- 🎯 Minimize configuration duplication

## 🛠️ CLI Commands

### Core Commands

```bash
# Initialize configuration
tivor init

# Plan infrastructure changes
tivor plan <environment> [--working-dir=<path>]

# Apply infrastructure changes  
tivor apply <environment> [--working-dir=<path>]

# Manage encrypted secrets
tivor sops encrypt <file>
tivor sops decrypt <file>

# Show version
tivor version
```

### Global Flags

```bash
-c, --config string      Path to configuration file (default "tivor.yaml")
    --log-level string   Log level: debug, info, warn, error (default "info")
```

## 🎯 Examples

Explore the `examples/` directory for complete working examples:

### Basic Example

```bash
cd examples/basic/

# Test development environment (1 instance, no monitoring)
../../bin/tivor plan dev --working-dir=./terraform

# Test staging environment (2 instances, monitoring enabled)  
../../bin/tivor plan staging --working-dir=./terraform
```

The basic example demonstrates:
- ✅ Environment inheritance (dev → staging → production)
- ✅ Variable file merging and override behavior
- ✅ Safe null_resource examples for testing
- ✅ Realistic infrastructure configuration patterns

## 🔧 Supported Backends

### Local Backend ✅
- Read variable files from local filesystem
- Relative and absolute path support
- Automatic path resolution

### S3 Backend 🚧
- Remote variable file storage (planned)
- Versioning and encryption support (planned)
- Cross-region replication (planned)

## 🔐 Secret Management

tivor integrates with SOPS for secure secret management:

```yaml
secrets:
  engine: sops
  sops_config_path: ".sops.yaml"
```

**Status**: 🚧 Integration planned for next release

## 🏗️ Architecture

### Variable Processing Pipeline

1. **Load Configuration** - Parse `tivor.yaml` and resolve inheritance
2. **Read Variable Files** - Load from configured backend  
3. **Merge Variables** - Smart merging with override precedence
4. **Decrypt Secrets** - SOPS integration for encrypted values
5. **Generate Combined File** - Create temporary `.tfvars` file
6. **Execute Terraform** - Run `terraform plan/apply` with merged variables

### Smart Variable Merging

tivor's revolutionary variable merging handles duplicate variables intelligently:

```hcl
# common.tfvars
instance_count = 1
environment = "dev"

# staging.tfvars  
instance_count = 2  # Overrides common.tfvars
monitoring = true   # New variable

# Result: instance_count = 2, environment = "dev", monitoring = true
```

## 🧪 Testing

Run the integration tests:

```bash
# Test with examples
cd examples/basic/
../../bin/tivor plan dev --working-dir=./terraform
../../bin/tivor apply staging --working-dir=./terraform
```

## 📋 Requirements

- **Go**: 1.22 or later
- **Terraform**: 1.0 or later
- **SOPS**: For secret management (optional)

## 🗺️ Roadmap

### v1.0 (Current) ✅
- [x] Basic CLI framework
- [x] Configuration file management
- [x] Environment inheritance  
- [x] Local backend support
- [x] Variable file merging
- [x] Terraform execution engine

### v1.1 (Next) 🚧
- [ ] SOPS integration for secret management
- [ ] S3 backend implementation
- [ ] Enhanced error handling and validation
- [ ] Configuration templates

### v1.2 (Future) 📋
- [ ] GCS backend support
- [ ] Terraform workspace management
- [ ] Plugin system for custom backends
- [ ] CI/CD integration helpers

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Terraform](https://terraform.io) - Infrastructure as Code
- [Cobra](https://github.com/spf13/cobra) - CLI framework  
- [SOPS](https://github.com/mozilla/sops) - Secret management
- The Go community for excellent tooling

---

**Ready to simplify your Terraform multi-environment workflows?** 🚀

Start with `tivor init` and explore the examples!