package checker

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/LarsArtmann/templates/repo-validation/internal/config"
)

// ValidationResult represents the result of validating a file requirement
type ValidationResult struct {
	// Requirement is the file requirement that was validated
	Requirement config.FileRequirement
	// Exists indicates whether the file exists
	Exists bool
	// Error is any error that occurred during validation
	Error error
}

// Checker is responsible for checking if files exist in a repository
type Checker struct {
	// Config is the configuration for the checker
	Config *config.Config
}

// NewChecker creates a new Checker
func NewChecker(cfg *config.Config) *Checker {
	return &Checker{
		Config: cfg,
	}
}

// CheckRepository checks if all required files exist in the repository
func (c *Checker) CheckRepository() ([]ValidationResult, error) {
	var results []ValidationResult

	// Check all files using the consolidated list
	allRequirements := config.GetAllFileRequirements()
	for _, req := range allRequirements {
		result := c.checkFile(req)
		results = append(results, result)
	}

	return results, nil
}

// checkFile checks if a file exists in the repository
func (c *Checker) checkFile(req config.FileRequirement) ValidationResult {
	filePath := filepath.Join(c.Config.RepoPath, req.Path)

	_, err := os.Stat(filePath)
	exists := !os.IsNotExist(err)

	if err != nil && !os.IsNotExist(err) {
		return ValidationResult{
			Requirement: req,
			Exists:      false,
			Error:       fmt.Errorf("error checking file %s: %w", req.Path, err),
		}
	}

	return ValidationResult{
		Requirement: req,
		Exists:      exists,
		Error:       nil,
	}
}

// FixMissingFiles generates missing files based on templates
func (c *Checker) FixMissingFiles(results []ValidationResult) error {
	if c.Config.DryRun {
		return nil
	}

	for _, result := range results {
		if !result.Exists && result.Error == nil && result.Requirement.TemplatePath != "" {
			if err := c.generateFile(result.Requirement); err != nil {
				return fmt.Errorf("error generating file %s: %w", result.Requirement.Path, err)
			}
		}
	}

	return nil
}

// generateFile generates a file from a template
func (c *Checker) generateFile(req config.FileRequirement) error {
	// Skip if no template path is provided
	if req.TemplatePath == "" {
		return nil
	}

	// Get the absolute path to the template
	templatePath := filepath.Join(filepath.Dir(os.Args[0]), req.TemplatePath)

	// Read the template file
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("error reading template %s: %w", req.TemplatePath, err)
	}

	// Parse template
	tmpl, err := template.New(req.Path).Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("error parsing template %s: %w", req.TemplatePath, err)
	}

	// Prepare template data
	data := map[string]interface{}{
		"RepoName": filepath.Base(c.Config.RepoPath),
		// Add more variables as needed
	}

	// Render template to buffer
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("error executing template %s: %w", req.TemplatePath, err)
	}

	// Create the output file
	outputPath := filepath.Join(c.Config.RepoPath, req.Path)

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("error creating directory for %s: %w", req.Path, err)
	}

	// Write the file
	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("error writing file %s: %w", req.Path, err)
	}

	return nil
}
