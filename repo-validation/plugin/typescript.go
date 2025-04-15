package plugin

import "github.com/LarsArtmann/templates/repo-validation/domain"

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
			Category: "TypeScript",
			Priority: "Must-have",
			Template: "package.json.tmpl",
		},
		{
			Path:     "tsconfig.json",
			Required: true,
			Category: "TypeScript",
			Priority: "Must-have",
			Template: "tsconfig.json.tmpl",
		},
		{
			Path:     ".eslintrc.json",
			Required: false,
			Category: "TypeScript",
			Priority: "Should-have",
			Template: ".eslintrc.json.tmpl",
		},
		{
			Path:     ".prettierrc",
			Required: false,
			Category: "TypeScript",
			Priority: "Should-have",
			Template: ".prettierrc.tmpl",
		},
		{
			Path:     "jest.config.js",
			Required: false,
			Category: "TypeScript",
			Priority: "Should-have",
			Template: "jest.config.js.tmpl",
		},
	}
}
