# Basic tfve Example

This example demonstrates core tfve functionality using safe `null_resource` examples that simulate real infrastructure deployment.

## 📁 Files Overview

```
basic/
├── README.md          # This file
├── tfve.yaml         # Configuration with environment inheritance
├── variables/        # Variable files for each environment
│   ├── common.tfvars     # Shared across all environments  
│   ├── dev.tfvars        # Development-specific
│   ├── staging.tfvars    # Staging-specific (inherits from dev)
│   └── production.tfvars # Production-specific (inherits from staging)
└── terraform/        # Sample Terraform code
    ├── main.tf           # Main infrastructure definition
    ├── variables.tf      # Variable declarations
    └── outputs.tf        # Output definitions
```

## 🚀 Usage Examples

First, ensure you're in the `examples/basic/` directory:
```bash
cd examples/basic/
```

### Development Environment
```bash
# Plan for development (1 instance, no monitoring)  
../../bin/tfve plan dev --working-dir=./terraform

# Apply development configuration
../../bin/tfve apply dev --working-dir=./terraform
```

### Staging Environment  
```bash
# Plan for staging (2 instances, monitoring enabled)
../../bin/tfve plan staging --working-dir=./terraform

# Apply staging configuration
../../bin/tfve apply staging --working-dir=./terraform
```

### Production Environment
```bash
# Plan for production (5 instances, full monitoring, S3 backend)
# Note: Will fail because S3 backend is not implemented yet
../../bin/tfve plan production --working-dir=./terraform
```

## 🎓 Learning Points

### Environment Inheritance
- **dev**: Uses `common.tfvars` + `dev.tfvars`
- **staging**: Uses `common.tfvars` + `dev.tfvars` + `staging.tfvars`
- **production**: Uses all four files with later files overriding earlier ones

### Variable Merging Behavior
Variables with the same name are overridden by later files:
- `instance_count`: 1 (dev) → 2 (staging) → 5 (production)
- `enable_monitoring`: false (dev) → true (staging/production)
- `environment_tags`: Completely replaced per environment

### Safe Testing
- Uses `null_resource` with `local-exec` provisioners
- No real infrastructure is created
- Safe to run multiple times
- Demonstrates realistic variable patterns

## 🔍 Expected Output

When running `tfve plan staging`, you should see:
- 2 infrastructure instances
- 1 monitoring setup resource  
- Proper tag merging
- Environment-specific values applied

## 🛠️ Customization

To adapt this example for your own infrastructure:
1. Replace `null_resource` with your actual resources (aws_instance, etc.)
2. Modify variable files to match your infrastructure needs
3. Update the backend configuration for your remote state storage
4. Add environment-specific variables as needed