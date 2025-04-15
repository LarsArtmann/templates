package config

import (
	"fmt"
	"testing"
)

func TestConfigOptions(t *testing.T) {
	// Test WithDryRun
	t.Run("WithDryRun", func(t *testing.T) {
		cfg := &Config{}
		opt := WithDryRun(true)
		opt(cfg)
		if !cfg.DryRun {
			t.Errorf("Expected DryRun to be true, got false")
		}
	})

	// Test WithFix
	t.Run("WithFix", func(t *testing.T) {
		cfg := &Config{}
		opt := WithFix(true)
		opt(cfg)
		if !cfg.Fix {
			t.Errorf("Expected Fix to be true, got false")
		}
	})

	// Test WithJSONOutput
	t.Run("WithJSONOutput", func(t *testing.T) {
		cfg := &Config{}
		opt := WithJSONOutput(true)
		opt(cfg)
		if !cfg.JSONOutput {
			t.Errorf("Expected JSONOutput to be true, got false")
		}
	})

	// Test WithRepoPath
	t.Run("WithRepoPath", func(t *testing.T) {
		cfg := &Config{}
		path := "/test/path"
		opt := WithRepoPath(path)
		opt(cfg)
		if cfg.RepoPath != path {
			t.Errorf("Expected RepoPath to be %q, got %q", path, cfg.RepoPath)
		}
	})

	// Test WithInteractive
	t.Run("WithInteractive", func(t *testing.T) {
		cfg := &Config{}
		opt := WithInteractive(true)
		opt(cfg)
		if !cfg.Interactive {
			t.Errorf("Expected Interactive to be true, got false")
		}
	})

	// Test WithFileGroup
	t.Run("WithFileGroup", func(t *testing.T) {
		tests := []struct {
			name    string
			group   string
			enabled bool
			check   func(*Config) bool
		}{
			{
				name:    "augment enabled",
				group:   "augment",
				enabled: true,
				check:   func(c *Config) bool { return c.CheckAugment },
			},
			{
				name:    "docker enabled",
				group:   "docker",
				enabled: true,
				check:   func(c *Config) bool { return c.CheckDocker },
			},
			{
				name:    "typescript enabled",
				group:   "typescript",
				enabled: true,
				check:   func(c *Config) bool { return c.CheckTypeScript },
			},
			{
				name:    "devcontainer enabled",
				group:   "devcontainer",
				enabled: true,
				check:   func(c *Config) bool { return c.CheckDevContainer },
			},
			{
				name:    "devenv enabled",
				group:   "devenv",
				enabled: true,
				check:   func(c *Config) bool { return c.CheckDevEnv },
			},
			{
				name:    "all enabled",
				group:   "all",
				enabled: true,
				check: func(c *Config) bool {
					return c.CheckAugment && c.CheckDocker && c.CheckTypeScript && c.CheckDevContainer && c.CheckDevEnv
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				cfg := &Config{}
				opt := WithFileGroup(tt.group, tt.enabled)
				opt(cfg)
				if !tt.check(cfg) {
					t.Errorf("Expected %s to be enabled", tt.group)
				}
			})
		}
	})
}

func TestConfigValidate(t *testing.T) {
	// Test valid configuration
	t.Run("valid config", func(t *testing.T) {
		cfg := &Config{
			RepoPath: "/test/path",
		}
		if err := cfg.Validate(); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	// Test missing repo path
	t.Run("missing repo path", func(t *testing.T) {
		cfg := &Config{
			RepoPath: "",
		}
		if err := cfg.Validate(); err == nil {
			t.Errorf("Expected error for missing repo path, got nil")
		}
	})

	// Test conflicting flags
	t.Run("conflicting flags", func(t *testing.T) {
		cfg := &Config{
			RepoPath: "/test/path",
			DryRun:   true,
			Fix:      true,
		}
		if err := cfg.Validate(); err == nil {
			t.Errorf("Expected error for conflicting flags, got nil")
		}
	})

	// Test JSON output with interactive mode
	t.Run("json with interactive", func(t *testing.T) {
		cfg := &Config{
			RepoPath:    "/test/path",
			JSONOutput:  true,
			Interactive: true,
		}
		if err := cfg.Validate(); err == nil {
			t.Errorf("Expected error for JSON output with interactive mode, got nil")
		}
	})

	// Test validation options
	t.Run("validation options", func(t *testing.T) {
		cfg := &Config{
			RepoPath: "/test/path",
		}

		// Create a validation option that always returns an error
		opt := func(c *Config) error {
			return fmt.Errorf("test error")
		}

		if err := cfg.Validate(opt); err == nil {
			t.Errorf("Expected error from validation option, got nil")
		}
	})

	// Test ValidateFileGroups
	t.Run("validate file groups", func(t *testing.T) {
		cfg := &Config{
			RepoPath: "/test/path",
		}

		if err := ValidateFileGroups(cfg); err != nil {
			t.Errorf("Expected no error for empty file groups, got %v", err)
		}

		cfg.CheckAugment = true
		if err := ValidateFileGroups(cfg); err != nil {
			t.Errorf("Expected no error for non-empty file groups, got %v", err)
		}
	})
}

func TestFileRequirementList(t *testing.T) {
	// Create a test list
	list := FileRequirementList{
		{Path: "file1.txt", Priority: PriorityMustHave},
		{Path: "file2.txt", Priority: PriorityShouldHave},
		{Path: "file3.txt", Priority: PriorityMustHave},
	}

	// Test FilterByPriority
	t.Run("FilterByPriority", func(t *testing.T) {
		filtered := list.FilterByPriority(PriorityMustHave)
		if len(filtered) != 2 {
			t.Errorf("Expected 2 must-have files, got %d", len(filtered))
		}
		for _, req := range filtered {
			if req.Priority != PriorityMustHave {
				t.Errorf("Expected priority %s, got %s", PriorityMustHave, req.Priority)
			}
		}
	})
}

func TestGetGeneralMustHaveFiles(t *testing.T) {
	// Test that GetGeneralMustHaveFiles returns a non-empty list
	files := GetGeneralMustHaveFiles()
	if len(files) == 0 {
		t.Errorf("Expected non-empty list of must-have files, got empty list")
	}

	// Check that all files have the correct priority
	for _, file := range files {
		if file.Priority != PriorityMustHave {
			t.Errorf("Expected priority %s, got %s", PriorityMustHave, file.Priority)
		}
	}
}

func TestGetGeneralShouldHaveFiles(t *testing.T) {
	// Test that GetGeneralShouldHaveFiles returns a non-empty list
	files := GetGeneralShouldHaveFiles()
	if len(files) == 0 {
		t.Errorf("Expected non-empty list of should-have files, got empty list")
	}

	// Check that all files have the correct priority
	for _, file := range files {
		if file.Priority != PriorityShouldHave && file.Priority != PriorityNiceToHave {
			t.Errorf("Expected priority %s or %s, got %s", PriorityShouldHave, PriorityNiceToHave, file.Priority)
		}
	}
}
