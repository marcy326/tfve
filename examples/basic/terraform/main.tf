# Main Terraform configuration
terraform {
  required_version = ">= 1.0"
  required_providers {
    null = {
      source  = "hashicorp/null"
      version = "~> 3.0"
    }
  }
}

# Local values for computed configurations
locals {
  merged_tags = merge(
    var.common_tags,
    var.environment_tags,
    {
      Terraform = "true"
      Environment = var.environment
    }
  )
  
  resource_name_prefix = "${var.project_name}-${var.environment}"
}

# Null resource to simulate infrastructure deployment
resource "null_resource" "demo_infrastructure" {
  count = var.instance_count

  triggers = {
    project_name = var.project_name
    environment = var.environment
    instance_type = var.instance_type
    monitoring = var.enable_monitoring
    backup_days = var.backup_retention_days
    timestamp = timestamp()
  }

  provisioner "local-exec" {
    command = "echo 'Simulating deployment of instance ${count.index + 1} for ${var.environment} environment'"
  }
}

# Conditional resource based on monitoring setting
resource "null_resource" "monitoring_setup" {
  count = var.enable_monitoring ? 1 : 0

  triggers = {
    environment = var.environment
    monitoring_enabled = var.enable_monitoring
  }

  provisioner "local-exec" {
    command = "echo 'Setting up monitoring for ${var.environment} environment'"
  }
}