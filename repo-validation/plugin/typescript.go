package plugin

import (
	"github.com/LarsArtmann/templates/repo-validation/domain"
)

// TypeScriptPlugin implements the TypeScript file group plugin
type TypeScriptPlugin struct{}

// Name returns the name of the plugin
func (p *TypeScriptPlugin) Name() string {
	return "typescript"
}

// Description returns the description of the plugin
func (p *TypeScriptPlugin) Description() string {
	return "TypeScript/JavaScript related files"
}

// Files returns the files checked by this plugin
func (p *TypeScriptPlugin) Files() []domain.File {
	return []domain.File{
		{
			Path:     "package.json",
			Required: true,
			Category: domain.CategoryTypeScript,
			Priority: domain.PriorityMustHave,
			Template: "package.json.tmpl",
		},
		{
			Path:     "tsconfig.json",
			Required: true,
			Category: domain.CategoryTypeScript,
			Priority: domain.PriorityMustHave,
			Template: "tsconfig.json.tmpl",
		},
		{
			Path:     ".eslintrc.json",
			Required: false,
			Category: domain.CategoryTypeScript,
			Priority: domain.PriorityShouldHave,
			Template: ".eslintrc.json.tmpl",
		},
		{
			Path:     ".prettierrc",
			Required: false,
			Category: domain.CategoryTypeScript,
			Priority: domain.PriorityShouldHave,
			Template: ".prettierrc.tmpl",
		},
		{
			Path:     "jest.config.js",
			Required: false,
			Category: domain.CategoryTypeScript,
			Priority: domain.PriorityShouldHave,
			Template: "jest.config.js.tmpl",
		},
	}
}
