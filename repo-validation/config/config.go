package config

import (
	"fmt"
)

// Config represents the application configuration
type Config struct {
	// General options
	DryRun     bool
	Fix        bool
	JSONOutput bool
	RepoPath   string
	Interactive bool
	
	// File groups
	CheckAll         bool
	CheckAugment     bool
	CheckDocker      bool
	CheckTypeScript  bool
	CheckDevContainer bool
	CheckDevEnv      bool
}

// ValidationOption is a function that performs additional validation on a Config
type ValidationOption func(*Config) error

// Validate validates the configuration
func (c *Config) Validate(opts ...ValidationOption) error {
	// Check for incompatible parameters
	if c.DryRun && c.Fix {
		return fmt.Errorf("--dry-run and --fix cannot be used together")
	}

	// Check if the repository path exists and is a directory
	if c.RepoPath == "" {
		return fmt.Errorf("repository path cannot be empty")
	}

	// Check if JSON output is enabled with interactive mode
	if c.JSONOutput && c.Interactive {
		return fmt.Errorf("--json and --interactive cannot be used together")
	}

	// Validate file groups when --all is used
	if err := ValidateFileGroups(c); err != nil {
		return err
	}

	// Apply additional validation options
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return err
		}
	}

	return nil
}

// ValidateFileGroups checks if at least one file group is selected when the --all flag is used
func ValidateFileGroups(c *Config) error {
	// If the --all flag is not set, we don't need to validate file groups
	if !c.CheckAll {
		return nil
	}

	// If --all is set, at least one file group should be selected
	if c.CheckAugment || c.CheckDocker || c.CheckTypeScript || c.CheckDevContainer || c.CheckDevEnv {
		return nil
	}

	// No file groups selected with --all flag
	return fmt.Errorf("no file groups selected, use at least one of --augment, --docker, etc., or remove --all flag")
}

// ConfigOption is a function that configures a Config
type ConfigOption func(*Config)

// WithDryRun sets the DryRun option
func WithDryRun(dryRun bool) ConfigOption {
	return func(c *Config) {
		c.DryRun = dryRun
	}
}

// WithFix sets the Fix option
func WithFix(fix bool) ConfigOption {
	return func(c *Config) {
		c.Fix = fix
	}
}

// WithJSONOutput sets the JSONOutput option
func WithJSONOutput(jsonOutput bool) ConfigOption {
	return func(c *Config) {
		c.JSONOutput = jsonOutput
	}
}

// WithRepoPath sets the RepoPath option
func WithRepoPath(repoPath string) ConfigOption {
	return func(c *Config) {
		c.RepoPath = repoPath
	}
}

// WithInteractive sets the Interactive option
func WithInteractive(interactive bool) ConfigOption {
	return func(c *Config) {
		c.Interactive = interactive
	}
}

// WithFileGroup sets a file group option
func WithFileGroup(group string, enabled bool) ConfigOption {
	return func(c *Config) {
		switch group {
		case "augment":
			c.CheckAugment = enabled
		case "docker":
			c.CheckDocker = enabled
		case "typescript":
			c.CheckTypeScript = enabled
		case "devcontainer":
			c.CheckDevContainer = enabled
		case "devenv":
			c.CheckDevEnv = enabled
		case "all":
			c.CheckAll = enabled
			c.CheckAugment = enabled
			c.CheckDocker = enabled
			c.CheckTypeScript = enabled
			c.CheckDevContainer = enabled
			c.CheckDevEnv = enabled
		}
	}
}
