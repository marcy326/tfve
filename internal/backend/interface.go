package backend

import "context"

// Config is generic configuration passed to each backend implementation.
type Config map[string]interface{}

// Backend is an interface for retrieving assets such as variable files.
// In v1.0, the primary responsibility is retrieving variable files.
type Backend interface {
	// GetVarsFile retrieves the content of a variable file at the specified path.
	// For remote backends (like S3), this performs download processing.
	// For local backends (like git), this returns the local path.
	GetVarsFile(ctx context.Context, path string) (content []byte, err error)
}

