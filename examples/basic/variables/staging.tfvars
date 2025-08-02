# Staging environment variables
environment = "staging"
instance_count = 2
instance_type = "t3.small"

# Staging specific settings
enable_monitoring = true
backup_retention_days = 14

# Environment-specific tags
environment_tags = {
  Environment = "staging"
  CostCenter = "qa"
}