package infrastructure

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/LarsArtmann/templates/repo-validation/domain"
	"github.com/LarsArtmann/templates/repo-validation/internal/templates"
)

// FileRepository implements the domain.FileRepository interface
type FileRepository struct {
	basePath string
}

// NewFileRepository creates a new FileRepository
func NewFileRepository(basePath string) *FileRepository {
	return &FileRepository{
		basePath: basePath,
	}
}

// CheckExists checks if a file exists
func (r *FileRepository) CheckExists(path string) (bool, error) {
	fullPath := filepath.Join(r.basePath, path)
	
	// Check if the file exists
	_, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("error checking if file exists: %w", err)
	}
	
	return true, nil
}

// Generate generates a file from a template
func (r *FileRepository) Generate(path, templatePath string, data interface{}) error {
	// Skip if no template path is provided
	if templatePath == "" {
		return domain.ErrNoTemplate
	}
	
	// Extract the template filename from the path
	templateFilename := filepath.Base(templatePath)
	
	// Read the template from the embedded filesystem
	templateContent, err := templates.TemplateFS.ReadFile(templateFilename)
	if err != nil {
		// Try to read from the filesystem as a fallback
		templatePath := filepath.Join(filepath.Dir(os.Args[0]), templatePath)
		templateContent, err = os.ReadFile(templatePath)
		if err != nil {
			return fmt.Errorf("error reading template %s: %w", templatePath, err)
		}
	}
	
	// Parse template
	tmpl, err := template.New(path).Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("error parsing template %s: %w", templatePath, err)
	}
	
	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("error executing template %s: %w", templatePath, err)
	}
	
	// Create the file
	fullPath := filepath.Join(r.basePath, path)
	
	// Create the directory if it doesn't exist
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating directory %s: %w", dir, err)
	}
	
	// Write the file
	if err := os.WriteFile(fullPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("error writing file %s: %w", fullPath, err)
	}
	
	return nil
}
