package terraform

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Executor handles Terraform command execution
type Executor struct {
	workingDir string
	varsFile   string
}

// NewExecutor creates a new Terraform executor
func NewExecutor(workingDir, varsFile string) *Executor {
	return &Executor{
		workingDir: workingDir,
		varsFile:   varsFile,
	}
}

// Plan executes terraform plan with the configured variables
func (e *Executor) Plan(ctx context.Context) error {
	return e.executeCommand(ctx, "plan", []string{})
}

// Apply executes terraform apply with the configured variables
func (e *Executor) Apply(ctx context.Context) error {
	return e.executeCommand(ctx, "apply", []string{"-auto-approve"})
}

// executeCommand executes a terraform command with common options
func (e *Executor) executeCommand(ctx context.Context, command string, extraArgs []string) error {
	// Check if terraform binary exists
	terraformPath, err := exec.LookPath("terraform")
	if err != nil {
		return fmt.Errorf("terraform binary not found in PATH: %w", err)
	}

	// Build command arguments
	args := []string{command}
	
	// Add variable file if specified
	if e.varsFile != "" {
		args = append(args, fmt.Sprintf("-var-file=%s", e.varsFile))
	}
	
	// Add extra arguments
	args = append(args, extraArgs...)

	// Create command
	cmd := exec.CommandContext(ctx, terraformPath, args...)
	
	// Set working directory
	if e.workingDir != "" {
		cmd.Dir = e.workingDir
	}

	// Set up output capture
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Log command execution
	slog.Info("Executing Terraform command", 
		"command", strings.Join(append([]string{"terraform"}, args...), " "),
		"working_dir", e.workingDir)

	// Execute command
	err = cmd.Run()
	
	// Print outputs
	if stdout.Len() > 0 {
		fmt.Print(stdout.String())
	}
	
	if stderr.Len() > 0 {
		fmt.Fprint(os.Stderr, stderr.String())
	}

	if err != nil {
		return fmt.Errorf("terraform %s failed: %w", command, err)
	}

	slog.Info("Terraform command completed successfully", "command", command)
	return nil
}

// ValidateWorkingDirectory checks if the working directory contains Terraform files
func (e *Executor) ValidateWorkingDirectory() error {
	if e.workingDir == "" {
		e.workingDir = "."
	}

	// Check if directory exists
	if _, err := os.Stat(e.workingDir); os.IsNotExist(err) {
		return fmt.Errorf("terraform working directory does not exist: %s", e.workingDir)
	}

	// Look for Terraform files
	tfFiles, err := filepath.Glob(filepath.Join(e.workingDir, "*.tf"))
	if err != nil {
		return fmt.Errorf("failed to search for terraform files: %w", err)
	}

	if len(tfFiles) == 0 {
		slog.Warn("No .tf files found in working directory", "dir", e.workingDir)
		return fmt.Errorf("no .tf files found in working directory: %s", e.workingDir)
	}

	slog.Info("Terraform files found", "count", len(tfFiles), "dir", e.workingDir)
	return nil
}

// Init executes terraform init to initialize the working directory
func (e *Executor) Init(ctx context.Context) error {
	return e.executeCommand(ctx, "init", []string{})
}