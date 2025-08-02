package cli

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

// NewSopsCmd creates the sops command.
func NewSopsCmd() *cobra.Command {
	sopsCmd := &cobra.Command{
		Use:   "sops <encrypt|decrypt> [file-path]",
		Short: "Encrypt and decrypt files using SOPS",
		Long: `Performs file encryption and decryption using the SOPS library.

Examples:
  tfve sops encrypt terraform/variables/secrets.tfvars
  tfve sops decrypt terraform/variables/secrets.tfvars.enc`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("please specify encrypt or decrypt")
			}

			action := args[0]
			var filePath string
			if len(args) > 1 {
				filePath = args[1]
			}

			return runSops(action, filePath)
		},
	}

	return sopsCmd
}

// runSops performs the actual processing of the sops command.
func runSops(action, filePath string) error {
	slog.Info("Starting SOPS command", "action", action, "file", filePath)

	switch action {
	case "encrypt":
		return runSopsEncrypt(filePath)
	case "decrypt":
		return runSopsDecrypt(filePath)
	default:
		return fmt.Errorf("invalid action: %s (please specify encrypt or decrypt)", action)
	}
}

// runSopsEncrypt performs file encryption.
func runSopsEncrypt(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("please specify file path to encrypt")
	}

	slog.Info("Preparing file encryption", "file", filePath)

	// TODO: Implement encryption processing using SOPS library
	// Currently only simple log output

	fmt.Printf("⚠️  The encryption feature is currently under development (file: %s)\n", filePath)
	fmt.Println("Integration with SOPS library is required.")

	return nil
}

// runSopsDecrypt performs file decryption.
func runSopsDecrypt(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("please specify file path to decrypt")
	}

	slog.Info("Preparing file decryption", "file", filePath)

	// TODO: Implement decryption processing using SOPS library
	// Currently only simple log output

	fmt.Printf("⚠️  The decryption feature is currently under development (file: %s)\n", filePath)
	fmt.Println("Integration with SOPS library is required.")

	return nil
}
