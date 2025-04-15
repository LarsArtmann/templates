package validation

import "github.com/LarsArtmann/templates/repo-validation/domain"

// Result represents the result of a validation
type Result struct {
	File  domain.File
	Error error
}

// Checker defines the interface for checking files
type Checker interface {
	CheckFile(path string) (bool, error)
}

// Reporter defines the interface for reporting results
type Reporter interface {
	ReportResults(results []Result) error
	GetSummary(results []Result) string
	GetExitCode(results []Result) int
}

// Generator defines the interface for generating files
type Generator interface {
	GenerateFile(path, template string, data interface{}) error
}

// ValidationService orchestrates the validation process
type ValidationService struct {
	checker   Checker
	reporter  Reporter
	generator Generator
}

// NewValidationService creates a new ValidationService
func NewValidationService(checker Checker, reporter Reporter, generator Generator) *ValidationService {
	return &ValidationService{
		checker:   checker,
		reporter:  reporter,
		generator: generator,
	}
}

// ValidateFiles validates a list of files
func (s *ValidationService) ValidateFiles(files []domain.File) []Result {
	results := make([]Result, 0, len(files))
	
	for _, file := range files {
		exists, err := s.checker.CheckFile(file.Path)
		
		result := Result{
			File: domain.File{
				Path:     file.Path,
				Exists:   exists,
				Required: file.Required,
				Category: file.Category,
				Priority: file.Priority,
				Template: file.Template,
			},
		}
		
		if err != nil {
			result.Error = err
		}
		
		results = append(results, result)
	}
	
	return results
}

// GenerateMissingFiles generates missing files from templates
func (s *ValidationService) GenerateMissingFiles(results []Result) error {
	for _, result := range results {
		if !result.File.Exists && result.File.Template != "" {
			if err := s.generator.GenerateFile(result.File.Path, result.File.Template, nil); err != nil {
				return err
			}
		}
	}
	
	return nil
}
