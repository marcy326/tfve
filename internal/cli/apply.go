package cli

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/marcy326/tivor/internal/config"
	"github.com/marcy326/tivor/internal/terraform"
	"github.com/spf13/cobra"
)

var (
	applyWorkingDir string
)

// NewApplyCmd creates the apply command.
func NewApplyCmd() *cobra.Command {
	applyCmd := &cobra.Command{
		Use:   "apply [environment-name]",
		Short: "Execute Terraform apply for the specified environment",
		Long: `Loads configuration for the specified environment, prepares variable files,
and executes terraform apply.

Examples:
  tivor apply staging
  tivor apply production
  tivor apply staging --working-dir=./infrastructure`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			envName := args[0]
			return runApply(envName, applyWorkingDir)
		},
	}

	applyCmd.Flags().StringVarP(&applyWorkingDir, "working-dir", "w", ".", "Terraform working directory")

	return applyCmd
}

// runApply performs the actual processing of the apply command.
func runApply(envName, workingDir string) error {
	slog.Info("Starting Terraform apply", "environment", envName, "working_dir", workingDir)

	config := GetConfig()
	if config == nil {
		return fmt.Errorf("configuration file not loaded")
	}

	// 1. Resolve environment configuration (including inheritance)
	env, err := config.ResolveEnvironment(envName)
	if err != nil {
		return fmt.Errorf("failed to resolve environment configuration: %w", err)
	}

	slog.Info("Environment configuration loaded",
		"environment", env.Name,
		"vars_files", env.VarsFiles,
		"backend_type", getBackendTypeForApply(env))

	// 2. Load variable files
	slog.Info("Loading variable files", "files", env.VarsFiles)
	ctx := context.Background()
	combinedVars, err := config.LoadVarsFiles(ctx, envName)
	if err != nil {
		return fmt.Errorf("failed to load variable files: %w", err)
	}
	slog.Info("Variable files loaded successfully", "total_size", len(combinedVars))

	// 3. Decrypt secrets (not implemented yet)
	if config.Secrets != nil && config.Secrets.Engine == "sops" {
		slog.Info("SOPS secret decryption will be implemented in future", "engine", config.Secrets.Engine)
		// TODO: Implement SOPS decryption processing
	}

	// 4. Create temporary variable file
	slog.Info("Creating temporary variable file")
	tmpDir, err := os.MkdirTemp("", "tivor-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			slog.Warn("Failed to cleanup temporary directory", "dir", tmpDir, "error", err)
		}
	}()

	tmpVarsFile := filepath.Join(tmpDir, fmt.Sprintf("%s.tfvars", envName))
	if err := os.WriteFile(tmpVarsFile, combinedVars, 0600); err != nil {
		return fmt.Errorf("failed to write temporary vars file: %w", err)
	}
	slog.Info("Temporary variable file created", "path", tmpVarsFile)

	// 5. Execute terraform apply
	executor := terraform.NewExecutor(workingDir, tmpVarsFile)

	// Validate working directory
	if err := executor.ValidateWorkingDirectory(); err != nil {
		return fmt.Errorf("terraform working directory validation failed: %w", err)
	}

	// Initialize terraform if needed
	slog.Info("Initializing Terraform")
	if err := executor.Init(ctx); err != nil {
		return fmt.Errorf("terraform init failed: %w", err)
	}

	// Execute terraform apply
	slog.Info("Executing Terraform apply")
	if err := executor.Apply(ctx); err != nil {
		return fmt.Errorf("terraform apply failed: %w", err)
	}

	fmt.Printf("✅ Terraform apply completed successfully for environment: %s\n", envName)
	fmt.Printf("📁 Variables file: %s\n", tmpVarsFile)
	fmt.Printf("📂 Working directory: %s\n", workingDir)

	return nil
}

// getBackendTypeForApply safely retrieves the backend type.
func getBackendTypeForApply(env *config.Environment) string {
	if env.Backend != nil {
		return env.Backend.Type
	}
	return "not-configured"
}
