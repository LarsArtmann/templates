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
	force    bool
}

// FileRepositoryOption is a function that configures a FileRepository
type FileRepositoryOption func(*FileRepository)

// WithForce configures the FileRepository to force overwrite existing files
func WithForce(force bool) FileRepositoryOption {
	return func(r *FileRepository) {
		r.force = force
	}
}

// NewFileRepository creates a new FileRepository
func NewFileRepository(basePath string, opts ...FileRepositoryOption) *FileRepository {
	repo := &FileRepository{
		basePath: basePath,
		force:    false, // Default to not overwriting existing files
	}

	// Apply options
	for _, opt := range opts {
		opt(repo)
	}

	return repo
}

// CheckExists checks if a file exists
func (r *FileRepository) CheckExists(path string) (bool, error) {
	if path == "" {
		return false, fmt.Errorf("path cannot be empty")
	}

	fullPath := filepath.Join(r.basePath, path)

	// Check if the file exists
	_, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("error checking if file exists at %s: %w", fullPath, err)
	}

	return true, nil
}

// Generate generates a file from a template
func (r *FileRepository) Generate(path, templatePath string, data interface{}) error {
	// Validate inputs
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

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

	// Check if the file already exists
	if _, err := os.Stat(fullPath); err == nil {
		// File exists, check if we should overwrite it
		if !r.force {
			return fmt.Errorf("file %s already exists, not overwriting (use force option to override)", fullPath)
		}
		// If force is true, we'll overwrite the file
	} else if !os.IsNotExist(err) {
		// Some other error occurred
		return fmt.Errorf("error checking if file %s exists: %w", fullPath, err)
	}

	// Write the file
	if err := os.WriteFile(fullPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("error writing file %s: %w", fullPath, err)
	}

	return nil
}
