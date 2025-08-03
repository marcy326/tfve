package cli

import (
	"log/slog"
	"os"

	"github.com/marcy326/tivor/internal/config"

	"github.com/spf13/cobra"
)

var (
	// Global flags
	configPath string
	logLevel   string

	// Global variables
	globalConfig *config.Config
)

// NewRootCmd creates the root command for tivor.
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "tivor",
		Short: "Terraform Infrastructure Variable Orchestrator",
		Long: `tivor is an orchestration tool for intuitive, safe, and declarative
lifecycle management across multiple Terraform environments.

It adopts a GitOps-First approach and manages all environment configurations
from a single tivor.yaml configuration file.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			setupLogging()
			// Skip config loading for commands that don't need it
			skipConfigCommands := []string{"version", "init", "sops"}
			shouldSkip := false
			for _, cmdName := range skipConfigCommands {
				if cmd.Name() == cmdName {
					shouldSkip = true
					break
				}
			}
			if !shouldSkip {
				loadConfig()
			}
		},
	}

	// Define global flags
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "tivor.yaml", "Path to configuration file")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Log level (debug, info, warn, error)")

	// Add subcommands
	rootCmd.AddCommand(NewVersionCmd())
	rootCmd.AddCommand(NewInitCmd())
	rootCmd.AddCommand(NewPlanCmd())
	rootCmd.AddCommand(NewApplyCmd())
	rootCmd.AddCommand(NewSopsCmd())

	return rootCmd
}

// setupLogging configures the log level.
func setupLogging() {
	var level slog.Level
	switch logLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}
	handler := slog.NewTextHandler(os.Stderr, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

// loadConfig loads the configuration file.
func loadConfig() {
	var err error
	globalConfig, err = config.LoadConfig(configPath)
	if err != nil {
		slog.Error("Failed to load configuration file", "error", err)
		os.Exit(1)
	}

	slog.Info("Configuration file loaded", "path", configPath, "version", globalConfig.Version)
}

// GetConfig returns the loaded configuration.
func GetConfig() *config.Config {
	return globalConfig
}
