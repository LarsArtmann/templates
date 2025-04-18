package config

import "fmt"

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

// WithFileGroup enables a specific file group
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

// Priority levels
const (
	PriorityMustHave   = "Must-have"
	PriorityShouldHave = "Should-have"
	PriorityNiceToHave = "Nice-to-have"
)

// Category types
const (
	CategoryGeneral    = "General"
	CategoryPublic     = "Public"
	CategoryDocker     = "Docker"
	CategoryJavaScript = "JavaScript"
	CategoryTypeScript = "TypeScript"
)

// Config represents the configuration for the repository validation script
type Config struct {
	// DryRun if true, only report issues without making changes
	DryRun bool
	// Fix if true, generate missing files
	Fix bool
	// JSONOutput if true, output results in JSON format
	JSONOutput bool
	// RepoPath path to the repository to validate
	RepoPath string
	// Interactive if true, prompt for missing parameters
	Interactive bool

	// File group flags
	CheckAll         bool // Check all file groups
	CheckAugment     bool // Check Augment AI related files (.augment-guidelines, .augmentignore)
	CheckDocker      bool // Check Docker related files (Dockerfile, docker-compose.yaml, .dockerignore)
	CheckTypeScript  bool // Check TypeScript/JavaScript related files (package.json, tsconfig.json)
	CheckDevContainer bool // Check DevContainer related files (.devcontainer.json)
	CheckDevEnv      bool // Check DevEnv related files (devenv.nix)
}

// ValidationOption is a function that performs additional validation on a Config
type ValidationOption func(*Config) error

// Validate checks if the configuration is valid and returns an error if not
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

// FileRequirement represents a file that should be present in a repository
type FileRequirement struct {
	// Path is the path to the file, relative to the repository root
	Path string
	// Category is the category of the file (General, Public, JavaScript, etc.)
	Category string
	// Priority is the priority of the file (Must-have, Should-have, Nice-to-have)
	Priority string
	// Description is a brief description of what the file is for
	Description string
	// TemplatePath is the path to the template file, if any
	TemplatePath string
}

// FileRequirementList represents a list of file requirements with helper methods
type FileRequirementList []FileRequirement

// FileGroup represents a group of file requirements
type FileGroup struct {
	Name         string
	Flag         *bool
	Requirements []FileRequirement
}

// GetFileGroups returns all file groups
func GetFileGroups(cfg *Config) []FileGroup {
	return []FileGroup{
		{
			Name:         "Core",
			Flag:         nil, // Always included
			Requirements: GetCoreFiles(),
		},
		{
			Name:         "Augment",
			Flag:         &cfg.CheckAugment,
			Requirements: GetAugmentFiles(),
		},
		{
			Name:         "Docker",
			Flag:         &cfg.CheckDocker,
			Requirements: GetDockerFiles(),
		},
		{
			Name:         "TypeScript",
			Flag:         &cfg.CheckTypeScript,
			Requirements: GetTypeScriptFiles(),
		},
		{
			Name:         "DevContainer",
			Flag:         &cfg.CheckDevContainer,
			Requirements: GetDevContainerFiles(),
		},
		{
			Name:         "DevEnv",
			Flag:         &cfg.CheckDevEnv,
			Requirements: GetDevEnvFiles(),
		},
	}
}

// GetAllFileRequirements returns all file requirements based on the configuration
func GetAllFileRequirements(cfg *Config) FileRequirementList {
	var allRequirements FileRequirementList

	// Get all file groups
	fileGroups := GetFileGroups(cfg)

	// Add requirements from each group based on flags
	for _, group := range fileGroups {
		// If the flag is nil or true, include the requirements
		if group.Flag == nil || *group.Flag {
			allRequirements = append(allRequirements, group.Requirements...)
		}
	}

	return allRequirements
}

// Filter returns file requirements that match the given filter function
func (list FileRequirementList) Filter(filterFn func(FileRequirement) bool) FileRequirementList {
	var filtered FileRequirementList
	for _, req := range list {
		if filterFn(req) {
			filtered = append(filtered, req)
		}
	}
	return filtered
}

// FilterByPriority returns file requirements with the specified priority
func (list FileRequirementList) FilterByPriority(priority string) FileRequirementList {
	return list.Filter(func(req FileRequirement) bool {
		return req.Priority == priority
	})
}

// FilterByCategory returns file requirements with the specified category
func (list FileRequirementList) FilterByCategory(category string) FileRequirementList {
	return list.Filter(func(req FileRequirement) bool {
		return req.Category == category
	})
}

// GetGeneralMustHaveFiles returns the list of general files that must be present in a repository
func GetGeneralMustHaveFiles() []FileRequirement {
	return []FileRequirement{
		{
			Path:        "README.md",
			Category:    CategoryGeneral,
			Priority:    PriorityMustHave,
			Description: "Primary documentation file that explains what the project does, how to install/use it, and other essential information",
			TemplatePath: "templates/README.md.tmpl",
		},
		{
			Path:        ".gitignore",
			Category:    CategoryGeneral,
			Priority:    PriorityMustHave,
			Description: "Specifies intentionally untracked files to ignore when using Git",
			TemplatePath: "templates/.gitignore.tmpl",
		},
		{
			Path:        "LICENSE.md",
			Category:    CategoryPublic,
			Priority:    PriorityMustHave,
			Description: "Defines the terms under which the software can be used, modified, and distributed",
			TemplatePath: "templates/LICENSE.md.tmpl",
		},
		{
			Path:        "SECURITY.md",
			Category:    CategoryPublic,
			Priority:    PriorityMustHave,
			Description: "Provides security policy and vulnerability reporting instructions",
			TemplatePath: "templates/SECURITY.md.tmpl",
		},
	}
}

// GetDockerFiles returns the list of Docker-related files
func GetDockerFiles() []FileRequirement {
	return []FileRequirement{
		{
			Path:        "Dockerfile",
			Category:    CategoryDocker,
			Priority:    PriorityMustHave,
			Description: "Instructions for building a Docker image for the application",
			TemplatePath: "", // No template for MVP - TODO(https://github.com/LarsArtmann/mono/issues/66)
		},
		{
			Path:        ".dockerignore",
			Category:    CategoryDocker,
			Priority:    PriorityShouldHave,
			Description: "Specifies files that should be excluded when building Docker images",
			TemplatePath: "", // No template for MVP - TODO(https://github.com/LarsArtmann/mono/issues/66)
		},
		{
			Path:        "docker-compose.yaml",
			Category:    CategoryDocker,
			Priority:    PriorityShouldHave,
			Description: "Defines and runs multi-container Docker applications",
			TemplatePath: "", // No template for MVP - TODO(https://github.com/LarsArtmann/mono/issues/66)
		},
	}
}

// GetTypeScriptFiles returns the list of TypeScript/JavaScript-related files
func GetTypeScriptFiles() []FileRequirement {
	return []FileRequirement{
		{
			Path:        "package.json",
			Category:    CategoryJavaScript,
			Priority:    PriorityMustHave,
			Description: "Defines project metadata and dependencies for Node.js projects",
			TemplatePath: "", // No template for MVP - TODO(https://github.com/LarsArtmann/mono/issues/66)
		},
		{
			Path:        "tsconfig.json",
			Category:    CategoryTypeScript,
			Priority:    PriorityMustHave,
			Description: "Configuration file for TypeScript compiler options",
			TemplatePath: "", // No template for MVP - TODO(https://github.com/LarsArtmann/mono/issues/66)
		},
	}
}

// GetAugmentFiles returns the list of Augment AI-related files
func GetAugmentFiles() []FileRequirement {
	return []FileRequirement{
		{
			Path:        ".augment-guidelines",
			Category:    CategoryGeneral,
			Priority:    PriorityShouldHave,
			Description: "Provides guidelines for Augment AI to follow when working with the codebase",
			TemplatePath: "templates/.augment-guidelines.tmpl",
		},
		{
			Path:        ".augmentignore",
			Category:    CategoryGeneral,
			Priority:    PriorityShouldHave,
			Description: "Controls what files Augment AI indexes in the workspace",
			TemplatePath: "templates/.augmentignore.tmpl",
		},
	}
}

// GetDevContainerFiles returns the list of DevContainer-related files
func GetDevContainerFiles() []FileRequirement {
	return []FileRequirement{
		{
			Path:        ".devcontainer.json",
			Category:    CategoryPublic,
			Priority:    PriorityNiceToHave,
			Description: "Configuration for development in a containerized environment",
			TemplatePath: "", // No template for MVP - TODO(https://github.com/LarsArtmann/mono/issues/66)
		},
	}
}

// GetDevEnvFiles returns the list of DevEnv-related files
func GetDevEnvFiles() []FileRequirement {
	return []FileRequirement{
		{
			Path:        "devenv.nix",
			Category:    CategoryPublic,
			Priority:    PriorityNiceToHave,
			Description: "Defines development environment using Nix for reproducible builds",
			TemplatePath: "", // No template for MVP - TODO(https://github.com/LarsArtmann/mono/issues/66)
		},
	}
}

// GetCoreFiles returns all core files (must-have and should-have)
func GetCoreFiles() []FileRequirement {
	return append(GetGeneralMustHaveFiles(), GetGeneralShouldHaveFiles()...)
}

// GetGeneralShouldHaveFiles returns the list of general files that should be present in a repository
func GetGeneralShouldHaveFiles() []FileRequirement {
	return []FileRequirement{
		{
			Path:        "AUTHORS",
			Category:    CategoryGeneral,
			Priority:    PriorityShouldHave,
			Description: "Lists all individuals who have contributed to the project",
			TemplatePath: "", // No template for MVP - TODO(https://github.com/LarsArtmann/mono/issues/66)
		},
		{
			Path:        "MAINTAINERS.md",
			Category:    CategoryGeneral,
			Priority:    PriorityShouldHave,
			Description: "Identifies current maintainers and their responsibilities",
			TemplatePath: "", // No template for MVP - TODO(https://github.com/LarsArtmann/mono/issues/66)
		},
		{
			Path:        ".editorconfig",
			Category:    CategoryGeneral,
			Priority:    PriorityShouldHave,
			Description: "Helps maintain consistent coding styles across various editors and IDEs",
			TemplatePath: "templates/.editorconfig.tmpl",
		},
		{
			Path:        "CONTRIBUTING.md",
			Category:    CategoryPublic,
			Priority:    PriorityShouldHave,
			Description: "Guidelines for how to contribute to the project",
			TemplatePath: "", // No template for MVP - TODO(https://github.com/LarsArtmann/mono/issues/66)
		},
		{
			Path:        "CODE-OF-CONDUCT.md",
			Category:    CategoryPublic,
			Priority:    PriorityShouldHave,
			Description: "Establishes expectations for behavior within the project community",
			TemplatePath: "", // No template for MVP - TODO(https://github.com/LarsArtmann/mono/issues/66)
		},
		{
			Path:        "CODEOWNERS",
			Category:    CategoryPublic,
			Priority:    PriorityShouldHave,
			Description: "Defines individuals or teams responsible for code in a repository",
			TemplatePath: "", // No template for MVP - TODO(https://github.com/LarsArtmann/mono/issues/66)
		},
	}
}


