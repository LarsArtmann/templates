package domain

import "fmt"

// Priority represents the importance of a file
type Priority string

const (
	// PriorityMustHave indicates that the file must be present
	PriorityMustHave Priority = "Must-have"
	// PriorityShouldHave indicates that the file should be present
	PriorityShouldHave Priority = "Should-have"
	// PriorityNiceToHave indicates that the file is nice to have
	PriorityNiceToHave Priority = "Nice-to-have"
)

// Category represents the type of file
type Category string

const (
	// CategoryDocumentation indicates that the file is documentation
	CategoryDocumentation Category = "Documentation"
	// CategoryGit indicates that the file is related to Git
	CategoryGit Category = "Git"
	// CategoryDevelopment indicates that the file is related to development
	CategoryDevelopment Category = "Development"
	// CategoryDocker indicates that the file is related to Docker
	CategoryDocker Category = "Docker"
	// CategoryTypeScript indicates that the file is related to TypeScript
	CategoryTypeScript Category = "TypeScript"
	// CategoryDevContainer indicates that the file is related to DevContainer
	CategoryDevContainer Category = "DevContainer"
	// CategoryDevEnv indicates that the file is related to DevEnv
	CategoryDevEnv Category = "DevEnv"
	// CategoryAugment indicates that the file is related to Augment
	CategoryAugment Category = "Augment"
)

// File represents a file in a repository
type File struct {
	Path     string
	Exists   bool
	Required bool
	Category Category
	Priority Priority
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
