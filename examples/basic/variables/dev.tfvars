# Development environment variables
environment = "dev"
instance_count = 1
instance_type = "t3.micro"

# Development specific settings
enable_monitoring = false
backup_retention_days = 7

# Environment-specific tags
environment_tags = {
  Environment = "development"
  CostCenter = "engineering"
}