package plugin

import (
	"github.com/LarsArtmann/templates/repo-validation/domain"
)

// AugmentPlugin implements the Augment file group plugin
type AugmentPlugin struct{}

// The Augment plugin validates files used by AI code assistance tools.
// .augment-guidelines contains instructions for AI assistants
// .augmentignore specifies files to exclude from AI analysis

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
			Category: domain.CategoryAugment,
			Priority: domain.PriorityMustHave,
			Template: ".augment-guidelines.tmpl",
		},
		{
			Path:     ".augmentignore",
			Required: true,
			Category: domain.CategoryAugment,
			Priority: domain.PriorityMustHave,
			Template: ".augmentignore.tmpl",
		},
	}
}
