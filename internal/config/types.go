package config

// Config represents the overall structure of tivor.yaml
type Config struct {
	Version      string        `yaml:"version"`
	Defaults     *Defaults     `yaml:"defaults,omitempty"`
	Secrets      *Secrets      `yaml:"secrets,omitempty"`
	Environments []Environment `yaml:"environments"`
}

// Defaults represents common default settings across environments
type Defaults struct {
	VarsFiles []string `yaml:"vars_files,omitempty"`
}

// Secrets represents secret management configuration
type Secrets struct {
	Engine         string `yaml:"engine"`
	SopsConfigPath string `yaml:"sops_config_path,omitempty"`
}

// Environment represents configuration for individual environments
type Environment struct {
	Name      string   `yaml:"name"`
	Inherits  string   `yaml:"inherits,omitempty"`
	VarsFiles []string `yaml:"vars_files,omitempty"`
	Backend   *Backend `yaml:"backend,omitempty"`
}

// Backend represents storage backend configuration
type Backend struct {
	Type   string                 `yaml:"type"`
	Config map[string]interface{} `yaml:"config,omitempty"`
}
