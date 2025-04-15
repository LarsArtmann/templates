package cmd

import (
	"fmt"

	"github.com/LarsArtmann/templates/repo-validation/internal/checker"
	"github.com/LarsArtmann/templates/repo-validation/internal/config"
	"github.com/LarsArtmann/templates/repo-validation/internal/errors"
	"github.com/LarsArtmann/templates/repo-validation/internal/reporter"
)

// ValidateCommand represents the validate command
type ValidateCommand struct {
	// Config is the configuration for the command
	Config *config.Config
}

// NewValidateCommand creates a new ValidateCommand
func NewValidateCommand(cfg *config.Config) *ValidateCommand {
	return &ValidateCommand{
		Config: cfg,
	}
}

// Execute executes the validate command
func (c *ValidateCommand) Execute() error {
	// Create a checker
	chk := checker.NewChecker(c.Config)

	// Check the repository
	results, err := chk.CheckRepository()
	if err != nil {
		return fmt.Errorf("error checking repository: %w", err)
	}

	// Create a reporter
	rep := reporter.NewReporter(c.Config)

	// Report the results
	if err := rep.ReportResults(results); err != nil {
		return fmt.Errorf("error reporting results: %w", err)
	}

	// Fix missing files if requested
	if c.Config.Fix {
		if err := chk.FixMissingFiles(results); err != nil {
			return fmt.Errorf("error fixing missing files: %w", err)
		}

		// Check the repository again after fixing
		results, err = chk.CheckRepository()
		if err != nil {
			return fmt.Errorf("error checking repository after fixing: %w", err)
		}

		// Report the results again
		fmt.Println("\nAfter fixing:")
		if err := rep.ReportResults(results); err != nil {
			return fmt.Errorf("error reporting results after fixing: %w", err)
		}
	}

	// Check for errors and missing must-have files
	var hasErrors bool
	var firstError error
	for _, result := range results {
		if result.Error != nil {
			hasErrors = true
			if firstError == nil {
				firstError = result.Error
			}
			continue
		}

		if !result.Exists && result.Requirement.Priority == config.PriorityMustHave {
			hasErrors = true
			if firstError == nil {
				firstError = fmt.Errorf("missing must-have file: %s", result.Requirement.Path)
			}
		}
	}

	// Return an error if there are any issues
	if hasErrors {
		return errors.NewMissingMustHaveFilesError(rep.GetSummary(results))
	}

	return nil
}
