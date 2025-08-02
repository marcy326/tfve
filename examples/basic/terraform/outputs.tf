# Output values
output "environment_summary" {
  description = "Summary of the environment configuration"
  value = {
    project_name = var.project_name
    environment = var.environment
    region = var.region
    instance_count = var.instance_count
    instance_type = var.instance_type
    monitoring_enabled = var.enable_monitoring
    backup_retention_days = var.backup_retention_days
  }
}

output "resource_name_prefix" {
  description = "Prefix used for resource naming"
  value = local.resource_name_prefix
}

output "merged_tags" {
  description = "All tags merged together"
  value = local.merged_tags
}

output "deployment_info" {
  description = "Information about the simulated deployment"
  value = {
    total_instances = var.instance_count
    monitoring_resources = var.enable_monitoring ? 1 : 0
    deployment_timestamp = timestamp()
  }
}