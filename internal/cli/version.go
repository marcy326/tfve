package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	// Version is the version of tfve
	Version = "0.1.0"
)

// NewVersionCmd creates the version command.
func NewVersionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Display version information for tfve",
		Long:  "Display version information for tfve.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("tfve version %s\n", Version)
		},
	}

	return versionCmd
}
