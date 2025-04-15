package plugin

import (
	"github.com/LarsArtmann/templates/repo-validation/domain"
)

// DockerPlugin implements the Docker file group plugin
type DockerPlugin struct{}

// Name returns the name of the plugin
func (p *DockerPlugin) Name() string {
	return "docker"
}

// Description returns the description of the plugin
func (p *DockerPlugin) Description() string {
	return "Docker related files"
}

// Files returns the files checked by this plugin
func (p *DockerPlugin) Files() []domain.File {
	return []domain.File{
		{
			Path:     "Dockerfile",
			Required: true,
			Category: domain.CategoryDocker,
			Priority: domain.PriorityMustHave,
			Template: "Dockerfile.tmpl",
		},
		{
			Path:     "docker-compose.yaml",
			Required: false,
			Category: domain.CategoryDocker,
			Priority: domain.PriorityShouldHave,
			Template: "docker-compose.yaml.tmpl",
		},
		{
			Path:     ".dockerignore",
			Required: false,
			Category: domain.CategoryDocker,
			Priority: domain.PriorityShouldHave,
			Template: ".dockerignore.tmpl",
		},
	}
}
