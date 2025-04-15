package domain

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
	exists, err := s.repo.CheckExists(file.Path)
	if err != nil {
		return err
	}
	
	file.Exists = exists
	
	// If the file is required and doesn't exist, return an error
	if file.Required && !file.Exists {
		return ErrRequiredFileMissing
	}
	
	return nil
}

// Generate generates a file from a template
func (s *FileService) Generate(file *File, data interface{}) error {
	if file.Template == "" {
		return ErrNoTemplate
	}
	
	return s.repo.Generate(file.Path, file.Template, data)
}
