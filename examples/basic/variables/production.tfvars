# Production environment variables
environment = "production"
instance_count = 5
instance_type = "t3.medium"

# Production specific settings
enable_monitoring = true
backup_retention_days = 30

# Production security settings
enable_encryption = true
enable_backup = true

# Environment-specific tags
environment_tags = {
  Environment = "production"
  CostCenter = "operations"
}