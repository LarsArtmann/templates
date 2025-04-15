package domain

import "fmt"

// File represents a file in a repository
type File struct {
	Path     string
	Exists   bool
	Required bool
	Category string
	Priority string
	Template string
}

// FileRepository defines the interface for file operations
type FileRepository interface {
	CheckExists(path string) (bool, error)
	Generate(path, templatePath string, data interface{}) error
}

// FileService contains business logic for file operations
type FileService struct {
	repo FileRepository
}

// NewFileService creates a new FileService
func NewFileService(repo FileRepository) *FileService {
	return &FileService{
		repo: repo,
	}
}

// Validate checks if a file exists and meets requirements
func (s *FileService) Validate(file *File) error {
	if file == nil {
		return fmt.Errorf("file cannot be nil")
	}

	if file.Path == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	exists, err := s.repo.CheckExists(file.Path)
	if err != nil {
		return fmt.Errorf("failed to check if file %s exists: %w", file.Path, err)
	}

	file.Exists = exists

	// If the file is required and doesn't exist, return an error
	if file.Required && !file.Exists {
		return fmt.Errorf("%w: %s", ErrRequiredFileMissing, file.Path)
	}

	return nil
}

// Generate generates a file from a template
func (s *FileService) Generate(file *File, data interface{}) error {
	if file == nil {
		return fmt.Errorf("file cannot be nil")
	}

	if file.Path == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	if file.Template == "" {
		return fmt.Errorf("%w: %s", ErrNoTemplate, file.Path)
	}

	err := s.repo.Generate(file.Path, file.Template, data)
	if err != nil {
		return fmt.Errorf("failed to generate file %s from template %s: %w", file.Path, file.Template, err)
	}

	return nil
}
