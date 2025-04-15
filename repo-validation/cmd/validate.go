package cmd

import (
	"fmt"

	"github.com/LarsArtmann/templates/repo-validation/internal/checker"
	"github.com/LarsArtmann/templates/repo-validation/internal/config"
	"github.com/LarsArtmann/templates/repo-validation/internal/errors"
	"github.com/LarsArtmann/templates/repo-validation/internal/exitcode"
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
	if c.Config.Fix && !c.Config.DryRun {
		if err := chk.FixMissingFiles(results); err != nil {
			return fmt.Errorf("error fixing missing files: %w", err)
		}

		// Check the repository again after fixing
		results, err = chk.CheckRepository()
		if err != nil {
			return fmt.Errorf("error checking repository after fixing: %w", err)
		}

		// Report the results again
		if !c.Config.JSONOutput {
			fmt.Println("\nAfter fixing:")
		}
		if err := rep.ReportResults(results); err != nil {
			return fmt.Errorf("error reporting results after fixing: %w", err)
		}
	} else if c.Config.Fix && c.Config.DryRun {
		// Flag conflict - both dry-run and fix are set
		if !c.Config.JSONOutput {
			fmt.Println("Warning: Both --dry-run and --fix flags are set. No files will be modified due to dry-run mode.")
		}
	}

	// Use the reporter to determine if we should exit with an error
	exitCode := rep.GetExitCode(results)
	if exitCode != exitcode.Success {
		// Return an appropriate error based on the exit code
		switch exitCode {
		case exitcode.MissingMustHaveFiles:
			return errors.NewMissingMustHaveFilesError(rep.GetSummary(results))
		default:
			return fmt.Errorf("repository validation failed: %s", rep.GetSummary(results))
		}
	}

	return nil
}
