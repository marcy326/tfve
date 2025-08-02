# tfve Examples

This directory contains examples and test cases for demonstrating tfve functionality.

## 📁 Directory Structure

- `basic/` - Basic usage example showing environment inheritance and variable merging

## 🚀 Quick Start

1. Navigate to the basic example:
   ```bash
   cd examples/basic/
   ```

2. Test with development environment:
   ```bash
   tfve plan dev --working-dir=./terraform
   ```

3. Test with staging environment (shows inheritance):
   ```bash
   tfve plan staging --working-dir=./terraform
   ```

4. Apply to any environment (safe with null_resource):
   ```bash
   tfve apply dev --working-dir=./terraform
   ```

## 🎯 What You'll Learn

- How to structure `tfve.yaml` configuration files
- Environment inheritance patterns (common → dev → staging → production)  
- Variable file organization and merging behavior
- Integration with existing Terraform projects

## 📝 Notes

- These examples use `null_resource` for safe demonstration
- Variable files show realistic infrastructure configuration patterns
- All examples are self-contained and safe to run
- Generated Terraform state files are excluded from version control

## 🧪 Testing

These examples serve as both documentation and integration tests for tfve functionality.