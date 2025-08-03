package config

import (
	"context"
	"fmt"
	"os"

	"github.com/marcy326/tivor/internal/backend"
	"github.com/marcy326/tivor/internal/backend/local"
	"github.com/marcy326/tivor/internal/tfvars"
	"gopkg.in/yaml.v3"
)

// LoadConfig loads and parses tivor.yaml from the specified path into Config.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file (%s): %w", path, err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file (%s): %w", path, err)
	}

	// Validation
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid config file (%s): %w", path, err)
	}

	return &config, nil
}

// validateConfig performs basic validation of the configuration file.
func validateConfig(config *Config) error {
	if config.Version == "" {
		return fmt.Errorf("version field is required")
	}

	if len(config.Environments) == 0 {
		return fmt.Errorf("at least one environment is required")
	}

	// Check for duplicate environment names
	envNames := make(map[string]bool)
	for _, env := range config.Environments {
		if env.Name == "" {
			return fmt.Errorf("environment name is required")
		}
		if envNames[env.Name] {
			return fmt.Errorf("duplicate environment name: %s", env.Name)
		}
		envNames[env.Name] = true
	}

	// Check inheritance relationships
	for _, env := range config.Environments {
		if env.Inherits != "" {
			if !envNames[env.Inherits] {
				return fmt.Errorf("environment %s inherits from non-existent environment %s", env.Name, env.Inherits)
			}
		}
	}

	return nil
}

// GetEnvironment retrieves environment configuration by name.
func (c *Config) GetEnvironment(name string) (*Environment, error) {
	for i := range c.Environments {
		if c.Environments[i].Name == name {
			return &c.Environments[i], nil
		}
	}
	return nil, fmt.Errorf("environment %s not found", name)
}

// ResolveEnvironment returns environment configuration with inheritance resolved.
func (c *Config) ResolveEnvironment(name string) (*Environment, error) {
	env, err := c.GetEnvironment(name)
	if err != nil {
		return nil, err
	}

	resolved := *env

	// Merge parent settings if inheritance is defined
	if env.Inherits != "" {
		parentEnv, err := c.ResolveEnvironment(env.Inherits)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve parent environment (%s): %w", env.Inherits, err)
		}

		// Copy parent settings first, then override with child settings
		if resolved.Backend == nil {
			resolved.Backend = parentEnv.Backend
		}

		// Merge VarsFiles from parent and child with deduplication
		if len(parentEnv.VarsFiles) > 0 {
			if len(resolved.VarsFiles) == 0 {
				resolved.VarsFiles = parentEnv.VarsFiles
			} else {
				// Add parent settings first, then child settings
				mergedVarsFiles := make([]string, 0, len(parentEnv.VarsFiles)+len(resolved.VarsFiles))
				mergedVarsFiles = append(mergedVarsFiles, parentEnv.VarsFiles...)
				mergedVarsFiles = append(mergedVarsFiles, resolved.VarsFiles...)
				resolved.VarsFiles = deduplicateSlice(mergedVarsFiles)
			}
		}
	}

	// Merge with default settings with deduplication
	if c.Defaults != nil && len(c.Defaults.VarsFiles) > 0 {
		if len(resolved.VarsFiles) == 0 {
			resolved.VarsFiles = c.Defaults.VarsFiles
		} else {
			// Add default settings first, then environment-specific settings
			mergedVarsFiles := make([]string, 0, len(c.Defaults.VarsFiles)+len(resolved.VarsFiles))
			mergedVarsFiles = append(mergedVarsFiles, c.Defaults.VarsFiles...)
			mergedVarsFiles = append(mergedVarsFiles, resolved.VarsFiles...)
			resolved.VarsFiles = deduplicateSlice(mergedVarsFiles)
		}
	}

	return &resolved, nil
}

// LoadVarsFiles loads and combines variable files for the specified environment
func (c *Config) LoadVarsFiles(ctx context.Context, envName string) ([]byte, error) {
	// Resolve environment configuration
	env, err := c.ResolveEnvironment(envName)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve environment: %w", err)
	}

	// Create backend
	var backendInstance backend.Backend
	if env.Backend != nil {
		switch env.Backend.Type {
		case "local":
			backendInstance, err = local.New(env.Backend.Config)
			if err != nil {
				return nil, fmt.Errorf("failed to create local backend: %w", err)
			}
		case "s3":
			return nil, fmt.Errorf("s3 backend not yet implemented")
		default:
			return nil, fmt.Errorf("unknown backend type: %s", env.Backend.Type)
		}
	} else {
		// Default to local backend if no backend is specified
		backendInstance, err = local.New(backend.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to create default local backend: %w", err)
		}
	}

	// Load and parse all variable files
	var allVariableSets [][]tfvars.Variable

	for _, varsFile := range env.VarsFiles {
		content, err := backendInstance.GetVarsFile(ctx, varsFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load vars file %s: %w", varsFile, err)
		}

		variables, err := tfvars.ParseTfvars(content)
		if err != nil {
			return nil, fmt.Errorf("failed to parse vars file %s: %w", varsFile, err)
		}

		allVariableSets = append(allVariableSets, variables)
	}

	// Merge all variables (later definitions override earlier ones)
	mergedVariables := tfvars.MergeVariables(allVariableSets...)

	// Generate final tfvars content
	finalContent := tfvars.GenerateTfvars(mergedVariables, envName)

	return []byte(finalContent), nil
}

// deduplicateSlice removes duplicate strings from a slice while preserving order
func deduplicateSlice(slice []string) []string {
	if len(slice) == 0 {
		return slice
	}

	seen := make(map[string]bool)
	result := make([]string, 0, len(slice))

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}
