package plugin

import "github.com/LarsArtmann/templates/repo-validation/domain"

// DevEnvPlugin implements the DevEnv file group plugin
type DevEnvPlugin struct{}

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
			Category: "DevEnv",
			Priority: "Must-have",
			Template: "devenv.nix.tmpl",
		},
		{
			Path:     "devenv.yaml",
			Required: false,
			Category: "DevEnv",
			Priority: "Should-have",
			Template: "devenv.yaml.tmpl",
		},
		{
			Path:     "devenv.lock",
			Required: false,
			Category: "DevEnv",
			Priority: "Should-have",
			Template: "devenv.lock.tmpl",
		},
	}
}
