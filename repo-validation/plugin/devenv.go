package plugin

import (
	"github.com/LarsArtmann/templates/repo-validation/domain"
)

// DevEnvPlugin implements the DevEnv file group plugin
type DevEnvPlugin struct{}

// DevEnv is a Nix-based developer environment manager (https://devenv.sh/)
// that provides reproducible and declarative development environments.
// The devenv.nix file contains the environment definition,
// devenv.yaml contains project-specific configuration,
// and devenv.lock is the dependency lockfile.

// Name returns the name of the plugin
func (p *DevEnvPlugin) Name() string {
	return "devenv"
}

// Description returns the description of the plugin
func (p *DevEnvPlugin) Description() string {
	return "DevEnv related files"
}

// Files returns the files checked by this plugin
func (p *DevEnvPlugin) Files() []domain.File {
	return []domain.File{
		{
			Path:     "devenv.nix",
			Required: true,
			Category: domain.CategoryDevEnv,
			Priority: domain.PriorityMustHave,
			Template: "devenv.nix.tmpl",
		},
		{
			Path:     "devenv.yaml",
			Required: false,
			Category: domain.CategoryDevEnv,
			Priority: domain.PriorityShouldHave,
			Template: "devenv.yaml.tmpl",
		},
		{
			Path:     "devenv.lock",
			Required: false,
			Category: domain.CategoryDevEnv,
			Priority: domain.PriorityShouldHave,
			Template: "devenv.lock.tmpl",
		},
	}
}
