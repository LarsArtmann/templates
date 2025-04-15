package plugin

import (
	"github.com/LarsArtmann/templates/repo-validation/domain"
)

// CorePlugin implements the core file group plugin
type CorePlugin struct{}

// Name returns the name of the plugin
func (p *CorePlugin) Name() string {
	return "core"
}

// Description returns the description of the plugin
func (p *CorePlugin) Description() string {
	return "Core files that should be present in all repositories"
}

// Files returns the files checked by this plugin
func (p *CorePlugin) Files() []domain.File {
	return []domain.File{
		{
			Path:     "README.md",
			Required: true,
			Category: domain.CategoryDocumentation,
			Priority: domain.PriorityMustHave,
			Template: "README.md.tmpl",
		},
		{
			Path:     "LICENSE.md",
			Required: true,
			Category: domain.CategoryDocumentation,
			Priority: domain.PriorityMustHave,
			Template: "LICENSE.md.tmpl",
		},
		{
			Path:     ".gitignore",
			Required: true,
			Category: domain.CategoryGit,
			Priority: domain.PriorityMustHave,
			Template: ".gitignore.tmpl",
		},
		{
			Path:     "SECURITY.md",
			Required: true,
			Category: domain.CategoryDocumentation,
			Priority: domain.PriorityMustHave,
			Template: "SECURITY.md.tmpl",
		},
		{
			Path:     "AUTHORS",
			Required: false,
			Category: domain.CategoryDocumentation,
			Priority: domain.PriorityShouldHave,
			Template: "AUTHORS.tmpl",
		},
		{
			Path:     "MAINTAINERS.md",
			Required: false,
			Category: domain.CategoryDocumentation,
			Priority: domain.PriorityShouldHave,
			Template: "MAINTAINERS.md.tmpl",
		},
		{
			Path:     ".editorconfig",
			Required: false,
			Category: domain.CategoryDevelopment,
			Priority: domain.PriorityShouldHave,
			Template: ".editorconfig.tmpl",
		},
		{
			Path:     "CONTRIBUTING.md",
			Required: false,
			Category: domain.CategoryDocumentation,
			Priority: domain.PriorityShouldHave,
			Template: "CONTRIBUTING.md.tmpl",
		},
		{
			Path:     "CODE-OF-CONDUCT.md",
			Required: false,
			Category: domain.CategoryDocumentation,
			Priority: domain.PriorityShouldHave,
			Template: "CODE-OF-CONDUCT.md.tmpl",
		},
		{
			Path:     "CODEOWNERS",
			Required: false,
			Category: domain.CategoryGit,
			Priority: domain.PriorityShouldHave,
			Template: "CODEOWNERS.tmpl",
		},
	}
}
