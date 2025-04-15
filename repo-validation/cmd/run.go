package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/LarsArtmann/templates/repo-validation/internal/checker"
	"github.com/LarsArtmann/templates/repo-validation/internal/config"
	"github.com/LarsArtmann/templates/repo-validation/internal/errors"
	"github.com/LarsArtmann/templates/repo-validation/internal/reporter"
)

// Run executes the main application logic
func Run(opts ...config.ConfigOption) error {
	// Create default configuration
	cfg := &config.Config{
		RepoPath: ".",
	}

	// Apply options
	for _, opt := range opts {
		opt(cfg)
	}

	// Validate the configuration
	if err := cfg.Validate(); err != nil {
		// If interactive mode is enabled, prompt for missing parameters
		if cfg.Interactive {
			if err := PromptForMissingParameters(cfg); err != nil {
				return errors.NewInvalidConfigError(err.Error())
			}
		} else {
			return errors.NewInvalidConfigError(err.Error())
		}
	}

	// Resolve the repository path to an absolute path
	absPath, err := filepath.Abs(cfg.RepoPath)
	if err != nil {
		return errors.NewPathError(cfg.RepoPath, err)
	}

	// Check if the path exists and is a directory
	stat, err := os.Stat(absPath)
	if err != nil {
		return errors.NewFileAccessError(absPath, err)
	}
	if !stat.IsDir() {
		return errors.NewPathError(absPath, fmt.Errorf("path is not a directory"))
	}
	cfg.RepoPath = absPath

	// Create a checker
	chk := checker.NewChecker(cfg)

	// Check the repository
	results, err := chk.CheckRepository()
	if err != nil {
		return fmt.Errorf("error checking repository: %w", err)
	}

	// Create a reporter
	rep := reporter.NewReporter(cfg)

	// Report the results
	if err := rep.ReportResults(results); err != nil {
		return fmt.Errorf("error reporting results: %w", err)
	}

	// Fix missing files if requested
	if cfg.Fix {
		if err := chk.FixMissingFiles(results); err != nil {
			return fmt.Errorf("error fixing missing files: %w", err)
		}

		// Check the repository again after fixing
		results, err = chk.CheckRepository()
		if err != nil {
			return fmt.Errorf("error checking repository after fixing: %w", err)
		}

		// Report the results again
		if !cfg.JSONOutput {
			fmt.Println("\nAfter fixing:")
		}
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
