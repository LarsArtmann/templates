package config

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

	// File group flags
	CheckAugment     bool // Check Augment AI related files (.augment-guidelines, .augmentignore)
	CheckDocker      bool // Check Docker related files (Dockerfile, docker-compose.yaml, .dockerignore)
	CheckTypeScript  bool // Check TypeScript/JavaScript related files (package.json, tsconfig.json)
	CheckDevContainer bool // Check DevContainer related files (.devcontainer.json)
	CheckDevEnv      bool // Check DevEnv related files (devenv.nix)
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

// GetAllFileRequirements returns all file requirements based on the configuration
func GetAllFileRequirements(cfg *Config) FileRequirementList {
	var allRequirements FileRequirementList

	// Always include general files
	allRequirements = append(allRequirements, GetMustHaveFiles()...)
	allRequirements = append(allRequirements, GetShouldHaveFiles()...)

	// Include optional file groups based on configuration
	if cfg.CheckAugment {
		allRequirements = append(allRequirements, GetAugmentFiles()...)
	}

	if cfg.CheckDocker {
		allRequirements = append(allRequirements, GetDockerFiles()...)
	}

	if cfg.CheckTypeScript {
		allRequirements = append(allRequirements, GetTypeScriptFiles()...)
	}

	if cfg.CheckDevContainer {
		allRequirements = append(allRequirements, GetDevContainerFiles()...)
	}

	if cfg.CheckDevEnv {
		allRequirements = append(allRequirements, GetDevEnvFiles()...)
	}

	return allRequirements
}

// FilterByPriority returns file requirements with the specified priority
func (list FileRequirementList) FilterByPriority(priority string) FileRequirementList {
	var filtered FileRequirementList
	for _, req := range list {
		if req.Priority == priority {
			filtered = append(filtered, req)
		}
	}
	return filtered
}

// FilterByCategory returns file requirements with the specified category
func (list FileRequirementList) FilterByCategory(category string) FileRequirementList {
	var filtered FileRequirementList
	for _, req := range list {
		if req.Category == category {
			filtered = append(filtered, req)
		}
	}
	return filtered
}

// GetGeneralMustHaveFiles returns the list of general files that must be present in a repository
func GetGeneralMustHaveFiles() []FileRequirement {
	return []FileRequirement{
		{
			Path:        "README.md",
			Category:    "General",
			Priority:    "Must-have",
			Description: "Primary documentation file that explains what the project does, how to install/use it, and other essential information",
			TemplatePath: "templates/README.md.tmpl",
		},
		{
			Path:        ".gitignore",
			Category:    "General",
			Priority:    "Must-have",
			Description: "Specifies intentionally untracked files to ignore when using Git",
			TemplatePath: "templates/.gitignore.tmpl",
		},
		{
			Path:        "LICENSE.md",
			Category:    "Public",
			Priority:    "Must-have",
			Description: "Defines the terms under which the software can be used, modified, and distributed",
			TemplatePath: "templates/LICENSE.md.tmpl",
		},
		{
			Path:        "SECURITY.md",
			Category:    "Public",
			Priority:    "Must-have",
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
			Category:    "Docker",
			Priority:    "Must-have",
			Description: "Instructions for building a Docker image for the application",
			TemplatePath: "", // No template for MVP
		},
		{
			Path:        ".dockerignore",
			Category:    "Docker",
			Priority:    "Should-have",
			Description: "Specifies files that should be excluded when building Docker images",
			TemplatePath: "", // No template for MVP
		},
		{
			Path:        "docker-compose.yaml",
			Category:    "Docker",
			Priority:    "Should-have",
			Description: "Defines and runs multi-container Docker applications",
			TemplatePath: "", // No template for MVP
		},
	}
}

// GetTypeScriptFiles returns the list of TypeScript/JavaScript-related files
func GetTypeScriptFiles() []FileRequirement {
	return []FileRequirement{
		{
			Path:        "package.json",
			Category:    "JavaScript",
			Priority:    "Must-have",
			Description: "Defines project metadata and dependencies for Node.js projects",
			TemplatePath: "", // No template for MVP
		},
		{
			Path:        "tsconfig.json",
			Category:    "TypeScript",
			Priority:    "Must-have",
			Description: "Configuration file for TypeScript compiler options",
			TemplatePath: "", // No template for MVP
		},
	}
}

// GetAugmentFiles returns the list of Augment AI-related files
func GetAugmentFiles() []FileRequirement {
	return []FileRequirement{
		{
			Path:        ".augment-guidelines",
			Category:    "General",
			Priority:    "Should-have",
			Description: "Provides guidelines for Augment AI to follow when working with the codebase",
			TemplatePath: "templates/.augment-guidelines.tmpl",
		},
		{
			Path:        ".augmentignore",
			Category:    "General",
			Priority:    "Should-have",
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
			Category:    "Public",
			Priority:    "Nice-to-have",
			Description: "Configuration for development in a containerized environment",
			TemplatePath: "", // No template for MVP
		},
	}
}

// GetDevEnvFiles returns the list of DevEnv-related files
func GetDevEnvFiles() []FileRequirement {
	return []FileRequirement{
		{
			Path:        "devenv.nix",
			Category:    "Public",
			Priority:    "Nice-to-have",
			Description: "Defines development environment using Nix for reproducible builds",
			TemplatePath: "", // No template for MVP
		},
	}
}

// GetMustHaveFiles returns the list of files that must be present in a repository
func GetMustHaveFiles() []FileRequirement {
	return GetGeneralMustHaveFiles()
}

// GetGeneralShouldHaveFiles returns the list of general files that should be present in a repository
func GetGeneralShouldHaveFiles() []FileRequirement {
	return []FileRequirement{
		{
			Path:        "AUTHORS",
			Category:    "General",
			Priority:    "Should-have",
			Description: "Lists all individuals who have contributed to the project",
			TemplatePath: "", // No template for MVP
		},
		{
			Path:        "MAINTAINERS.md",
			Category:    "General",
			Priority:    "Should-have",
			Description: "Identifies current maintainers and their responsibilities",
			TemplatePath: "", // No template for MVP
		},
		{
			Path:        ".editorconfig",
			Category:    "General",
			Priority:    "Should-have",
			Description: "Helps maintain consistent coding styles across various editors and IDEs",
			TemplatePath: "templates/.editorconfig.tmpl",
		},
		{
			Path:        "CONTRIBUTING.md",
			Category:    "Public",
			Priority:    "Should-have",
			Description: "Guidelines for how to contribute to the project",
			TemplatePath: "", // No template for MVP
		},
		{
			Path:        "CODE-OF-CONDUCT.md",
			Category:    "Public",
			Priority:    "Should-have",
			Description: "Establishes expectations for behavior within the project community",
			TemplatePath: "", // No template for MVP
		},
		{
			Path:        "CODEOWNERS",
			Category:    "Public",
			Priority:    "Should-have",
			Description: "Defines individuals or teams responsible for code in a repository",
			TemplatePath: "", // No template for MVP
		},
	}
}

// GetShouldHaveFiles returns the list of files that should be present in a repository
func GetShouldHaveFiles() []FileRequirement {
	return GetGeneralShouldHaveFiles()
}
