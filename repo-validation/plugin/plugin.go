package plugin

import "github.com/LarsArtmann/templates/repo-validation/domain"

// FileGroupPlugin defines the interface for file group plugins
type FileGroupPlugin interface {
	Name() string
	Description() string
	Files() []domain.File
}

// PluginRegistry manages file group plugins
type PluginRegistry struct {
	plugins map[string]FileGroupPlugin
}

// NewPluginRegistry creates a new PluginRegistry
func NewPluginRegistry() *PluginRegistry {
	return &PluginRegistry{
		plugins: make(map[string]FileGroupPlugin),
	}
}

// Register registers a plugin
func (r *PluginRegistry) Register(plugin FileGroupPlugin) {
	r.plugins[plugin.Name()] = plugin
}

// Get returns a plugin by name
func (r *PluginRegistry) Get(name string) (FileGroupPlugin, bool) {
	plugin, ok := r.plugins[name]
	return plugin, ok
}

// GetAll returns all registered plugins
func (r *PluginRegistry) GetAll() map[string]FileGroupPlugin {
	return r.plugins
}

// GetFiles returns all files from all registered plugins
func (r *PluginRegistry) GetFiles() []domain.File {
	var files []domain.File
	
	for _, plugin := range r.plugins {
		files = append(files, plugin.Files()...)
	}
	
	return files
}
