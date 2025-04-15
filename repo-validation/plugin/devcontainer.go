package plugin

import (
	"github.com/LarsArtmann/templates/repo-validation/domain"
)

// DevContainerPlugin implements the DevContainer file group plugin
type DevContainerPlugin struct{}

// Name returns the name of the plugin
func (p *DevContainerPlugin) Name() string {
	return "devcontainer"
}

// Description returns the description of the plugin
func (p *DevContainerPlugin) Description() string {
	return "DevContainer related files"
}

// Files returns the files checked by this plugin
func (p *DevContainerPlugin) Files() []domain.File {
	return []domain.File{
		{
			Path:     ".devcontainer/devcontainer.json",
			Required: true,
			Category: domain.CategoryDevContainer,
			Priority: domain.PriorityMustHave,
			Template: "devcontainer.json.tmpl",
		},
		{
			Path:     ".devcontainer/Dockerfile",
			Required: false,
			Category: domain.CategoryDevContainer,
			Priority: domain.PriorityShouldHave,
			Template: "devcontainer.Dockerfile.tmpl",
		},
		{
			Path:     ".devcontainer/docker-compose.yml",
			Required: false,
			Category: domain.CategoryDevContainer,
			Priority: domain.PriorityShouldHave,
			Template: "devcontainer.docker-compose.yml.tmpl",
		},
	}
}
