package plugin

import "github.com/LarsArtmann/templates/repo-validation/domain"

// AugmentPlugin implements the Augment file group plugin
type AugmentPlugin struct{}

// Name returns the name of the plugin
func (p *AugmentPlugin) Name() string {
	return "augment"
}

// Description returns the description of the plugin
func (p *AugmentPlugin) Description() string {
	return "Augment AI related files"
}

// Files returns the files checked by this plugin
func (p *AugmentPlugin) Files() []domain.File {
	return []domain.File{
		{
			Path:     ".augment-guidelines",
			Required: true,
			Category: "Augment",
			Priority: "Must-have",
			Template: ".augment-guidelines.tmpl",
		},
		{
			Path:     ".augmentignore",
			Required: true,
			Category: "Augment",
			Priority: "Must-have",
			Template: ".augmentignore.tmpl",
		},
	}
}
