package local

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/marcy326/tfve/internal/backend"
)

// LocalBackend implements Backend interface for local filesystem access
type LocalBackend struct {
	basePath string
}

// New creates a new LocalBackend instance
func New(config backend.Config) (backend.Backend, error) {
	basePath := "."
	if path, ok := config["path"].(string); ok && path != "" {
		basePath = path
	}

	// Convert to absolute path for consistency
	absPath, err := filepath.Abs(basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve absolute path for %s: %w", basePath, err)
	}

	return &LocalBackend{
		basePath: absPath,
	}, nil
}

// GetVarsFile retrieves the content of a variable file from local filesystem
func (l *LocalBackend) GetVarsFile(ctx context.Context, path string) ([]byte, error) {
	// Handle both absolute and relative paths
	var fullPath string
	if filepath.IsAbs(path) {
		fullPath = path
	} else {
		fullPath = filepath.Join(l.basePath, path)
	}

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("variable file not found: %s", fullPath)
	}

	// Read file content
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read variable file %s: %w", fullPath, err)
	}

	return content, nil
}
