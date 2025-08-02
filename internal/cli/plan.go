package cli

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/marcy326/tfve/internal/config"
	"github.com/marcy326/tfve/internal/terraform"
	"github.com/spf13/cobra"
)

var (
	planWorkingDir string
)

// NewPlanCmd creates the plan command.
func NewPlanCmd() *cobra.Command {
	planCmd := &cobra.Command{
		Use:   "plan [environment-name]",
		Short: "Execute Terraform plan for the specified environment",
		Long: `Loads configuration for the specified environment, prepares variable files,
and executes terraform plan.

Examples:
  tfve plan staging
  tfve plan production
  tfve plan staging --working-dir=./infrastructure`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			envName := args[0]
			return runPlan(envName, planWorkingDir)
		},
	}

	planCmd.Flags().StringVarP(&planWorkingDir, "working-dir", "w", ".", "Terraform working directory")

	return planCmd
}

// runPlan performs the actual processing of the plan command.
func runPlan(envName, workingDir string) error {
	slog.Info("Starting Terraform plan", "environment", envName, "working_dir", workingDir)

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
		"backend_type", getBackendType(env))

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
	tmpDir, err := os.MkdirTemp("", "tfve-*")
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

	// 5. Execute terraform plan
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

	// Execute terraform plan
	slog.Info("Executing Terraform plan")
	if err := executor.Plan(ctx); err != nil {
		return fmt.Errorf("terraform plan failed: %w", err)
	}

	fmt.Printf("✅ Terraform plan completed successfully for environment: %s\n", envName)
	fmt.Printf("📁 Variables file: %s\n", tmpVarsFile)
	fmt.Printf("📂 Working directory: %s\n", workingDir)

	return nil
}

// getBackendType safely retrieves the backend type.
func getBackendType(env *config.Environment) string {
	if env.Backend != nil {
		return env.Backend.Type
	}
	return "not-configured"
}
