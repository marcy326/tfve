package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewInitCmd creates the init command.
func NewInitCmd() *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Generate a tfve.yaml template",
		Long:  "An interactive command to generate a tfve.yaml template.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit()
		},
	}

	return initCmd
}

// runInit performs the actual processing of the init command.
func runInit() error {
	configPath := "tfve.yaml"

	// Check if file already exists
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("File %s already exists. Overwrite? (y/N): ", configPath)
		var response string
		_, err := fmt.Scanln(&response)
		if err != nil {
			// If there's an input error, default to "no"
			response = "n"
		}
		if response != "y" && response != "Y" {
			fmt.Println("Operation cancelled.")
			return nil
		}
	}

	// Template content
	template := `# tfve.yaml - Terraform Variable Environment Configuration File

# Configuration file version
version: "1.0"

# Default settings for the entire project
defaults:
  # Path to files defining variables common across environments
  vars_files:
    - "terraform/variables/common.tfvars"

# Secret management configuration
secrets:
  # Secret engine to use (v1.0 supports sops only)
  engine: sops
  # Path to SOPS configuration file
  sops_config_path: ".sops.yaml"

# List of environment definitions
environments:
  # --- Development Environment ---
  - name: dev
    vars_files:
      - "terraform/variables/dev.tfvars"
    backend:
      type: local

  # --- Staging Environment ---
  - name: staging
    # Inherit settings from dev environment
    inherits: dev
    vars_files:
      - "terraform/variables/staging.tfvars"
    backend:
      type: local

  # --- Production Environment ---
  - name: production
    # Inherit settings from staging environment
    inherits: staging
    vars_files:
      - "terraform/variables/production.tfvars"
    # Example of using remote backend for production
    backend:
      type: s3
      config:
        bucket: "your-tfstate-bucket"
        key: "terraform.tfstate"
        region: "ap-northeast-1"
        dynamodb_table: "your-terraform-lock-table"
`

	// Write to file
	if err := os.WriteFile(configPath, []byte(template), 0644); err != nil {
		return fmt.Errorf("failed to create configuration file: %w", err)
	}

	fmt.Printf("✅ Created %s.\n", configPath)
	fmt.Println("Please edit the configuration as needed.")

	return nil
}
