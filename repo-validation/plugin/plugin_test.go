package plugin

import (
	"testing"

	"github.com/LarsArtmann/templates/repo-validation/domain"
)

// MockPlugin is a mock implementation of FileGroupPlugin for testing
type MockPlugin struct {
	name        string
	description string
	files       []domain.File
}

func (m *MockPlugin) Name() string {
	return m.name
}

func (m *MockPlugin) Description() string {
	return m.description
}

func (m *MockPlugin) Files() []domain.File {
	return m.files
}

func TestPluginRegistry_Register(t *testing.T) {
	registry := NewPluginRegistry()
	
	plugin := &MockPlugin{
		name:        "test",
		description: "Test plugin",
		files: []domain.File{
			{Path: "test.txt"},
		},
	}
	
	registry.Register(plugin)
	
	if len(registry.plugins) != 1 {
		t.Errorf("PluginRegistry.Register() registry.plugins = %v, want %v", len(registry.plugins), 1)
	}
	
	if registry.plugins["test"] != plugin {
		t.Errorf("PluginRegistry.Register() registry.plugins[\"test\"] = %v, want %v", registry.plugins["test"], plugin)
	}
}

func TestPluginRegistry_Get(t *testing.T) {
	registry := NewPluginRegistry()
	
	plugin := &MockPlugin{
		name:        "test",
		description: "Test plugin",
		files: []domain.File{
			{Path: "test.txt"},
		},
	}
	
	registry.Register(plugin)
	
	got, ok := registry.Get("test")
	
	if !ok {
		t.Errorf("PluginRegistry.Get() ok = %v, want %v", ok, true)
	}
	
	if got != plugin {
		t.Errorf("PluginRegistry.Get() got = %v, want %v", got, plugin)
	}
	
	_, ok = registry.Get("nonexistent")
	
	if ok {
		t.Errorf("PluginRegistry.Get() ok = %v, want %v", ok, false)
	}
}

func TestPluginRegistry_GetAll(t *testing.T) {
	registry := NewPluginRegistry()
	
	plugin1 := &MockPlugin{
		name:        "test1",
		description: "Test plugin 1",
		files: []domain.File{
			{Path: "test1.txt"},
		},
	}
	
	plugin2 := &MockPlugin{
		name:        "test2",
		description: "Test plugin 2",
		files: []domain.File{
			{Path: "test2.txt"},
		},
	}
	
	registry.Register(plugin1)
	registry.Register(plugin2)
	
	plugins := registry.GetAll()
	
	if len(plugins) != 2 {
		t.Errorf("PluginRegistry.GetAll() plugins = %v, want %v", len(plugins), 2)
	}
	
	if plugins["test1"] != plugin1 {
		t.Errorf("PluginRegistry.GetAll() plugins[\"test1\"] = %v, want %v", plugins["test1"], plugin1)
	}
	
	if plugins["test2"] != plugin2 {
		t.Errorf("PluginRegistry.GetAll() plugins[\"test2\"] = %v, want %v", plugins["test2"], plugin2)
	}
}

func TestPluginRegistry_GetFiles(t *testing.T) {
	registry := NewPluginRegistry()
	
	plugin1 := &MockPlugin{
		name:        "test1",
		description: "Test plugin 1",
		files: []domain.File{
			{Path: "test1.txt"},
		},
	}
	
	plugin2 := &MockPlugin{
		name:        "test2",
		description: "Test plugin 2",
		files: []domain.File{
			{Path: "test2.txt"},
		},
	}
	
	registry.Register(plugin1)
	registry.Register(plugin2)
	
	files := registry.GetFiles()
	
	if len(files) != 2 {
		t.Errorf("PluginRegistry.GetFiles() files = %v, want %v", len(files), 2)
	}
	
	if files[0].Path != "test1.txt" && files[1].Path != "test1.txt" {
		t.Errorf("PluginRegistry.GetFiles() files does not contain test1.txt")
	}
	
	if files[0].Path != "test2.txt" && files[1].Path != "test2.txt" {
		t.Errorf("PluginRegistry.GetFiles() files does not contain test2.txt")
	}
}
