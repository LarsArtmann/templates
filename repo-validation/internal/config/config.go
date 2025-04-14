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

// GetMustHaveFiles returns the list of files that must be present in a repository
func GetMustHaveFiles() []FileRequirement {
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

// GetShouldHaveFiles returns the list of files that should be present in a repository
func GetShouldHaveFiles() []FileRequirement {
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
