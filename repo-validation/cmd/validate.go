package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/LarsArtmann/templates/repo-validation/config"
	"github.com/LarsArtmann/templates/repo-validation/infrastructure"
	"github.com/LarsArtmann/templates/repo-validation/internal/errors"
	"github.com/LarsArtmann/templates/repo-validation/internal/exitcode"
	"github.com/LarsArtmann/templates/repo-validation/plugin"
	"github.com/LarsArtmann/templates/repo-validation/validation"
)

// runValidate executes the validation command
func runValidate(opts ...config.ConfigOption) error {
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
			if err := promptForMissingParameters(cfg); err != nil {
				return fmt.Errorf("error in interactive mode: %w", err)
			}
		} else {
			return fmt.Errorf("invalid configuration: %w", err)
		}
	}

	// Resolve the repository path to an absolute path
	absPath, err := filepath.Abs(cfg.RepoPath)
	if err != nil {
		return fmt.Errorf("failed to resolve repository path: %w", err)
	}

	// Check if the path exists and is a directory
	stat, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("failed to access repository path: %w", err)
	}
	if !stat.IsDir() {
		return fmt.Errorf("repository path is not a directory: %s", absPath)
	}
	cfg.RepoPath = absPath

	// Create the file repository
	fileRepo := infrastructure.NewFileRepository(cfg.RepoPath)

	// Create the plugin registry
	registry := plugin.NewPluginRegistry()

	// Register plugins
	registry.Register(&plugin.CorePlugin{})

	// Register optional plugins based on configuration
	if cfg.CheckAugment {
		registry.Register(&plugin.AugmentPlugin{})
	}
	if cfg.CheckDocker {
		registry.Register(&plugin.DockerPlugin{})
	}
	if cfg.CheckTypeScript {
		registry.Register(&plugin.TypeScriptPlugin{})
	}
	if cfg.CheckDevContainer {
		registry.Register(&plugin.DevContainerPlugin{})
	}
	if cfg.CheckDevEnv {
		registry.Register(&plugin.DevEnvPlugin{})
	}

	// Get all files to check
	files := registry.GetFiles()

	// Create the validation service
	service := createValidationService(cfg, fileRepo)

	// Validate the files
	results := service.ValidateFiles(files)

	// Create a reporter
	reporter := createReporter(cfg)

	// Report the results
	if err := reporter.ReportResults(results); err != nil {
		return fmt.Errorf("error reporting results: %w", err)
	}

	// Fix missing files if requested
	if cfg.Fix && !cfg.DryRun {
		if err := service.GenerateMissingFiles(results); err != nil {
			return fmt.Errorf("error fixing missing files: %w", err)
		}

		// Check the repository again after fixing
		results = service.ValidateFiles(files)

		// Report the results again
		if !cfg.JSONOutput {
			fmt.Println("\nAfter fixing:")
		}
		if err := reporter.ReportResults(results); err != nil {
			return fmt.Errorf("error reporting results after fixing: %w", err)
		}
	} else if cfg.Fix && cfg.DryRun {
		// Flag conflict - both dry-run and fix are set
		if !cfg.JSONOutput {
			fmt.Println("Warning: Both --dry-run and --fix flags are set. No files will be modified due to dry-run mode.")
		}
	}

	// Get the exit code
	exitCode := reporter.GetExitCode(results)

	// Return an error if the exit code is not success
	if exitCode != exitcode.Success {
		// Return an appropriate error based on the exit code
		switch exitCode {
		case exitcode.MissingMustHaveFiles:
			return errors.NewMissingMustHaveFilesError(reporter.GetSummary(results))
		default:
			return fmt.Errorf("repository validation failed: %s", reporter.GetSummary(results))
		}
	}

	return nil
}

// createValidationService creates a validation service
func createValidationService(cfg *config.Config, fileRepo *infrastructure.FileRepository) *validation.ValidationService {
	// Create the checker adapter
	checker := infrastructure.NewCheckerAdapter(cfg)

	// Create the generator adapter
	generator := infrastructure.NewGeneratorAdapter(fileRepo)

	// Create the reporter adapter
	reporter := createReporter(cfg)

	// Create the validation service
	return validation.NewValidationService(checker, reporter, generator)
}

// createReporter creates a reporter
func createReporter(cfg *config.Config) validation.Reporter {
	// Create the reporter adapter
	return infrastructure.NewReporterAdapter(cfg)
}
