package main

import (
	"log/slog"
	"os"

	"github.com/marcy326/tfve/internal/cli"
)

func main() {
	rootCmd := cli.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		slog.Error("Failed to execute command", "error", err)
		os.Exit(1)
	}
}
