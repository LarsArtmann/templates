package infrastructure

import (
	"fmt"
	"github.com/LarsArtmann/templates/repo-validation/config"
	"github.com/LarsArtmann/templates/repo-validation/domain"
	"github.com/LarsArtmann/templates/repo-validation/internal/checker"
	internalConfig "github.com/LarsArtmann/templates/repo-validation/internal/config"
	"github.com/LarsArtmann/templates/repo-validation/internal/reporter"
	"github.com/LarsArtmann/templates/repo-validation/validation"
)

// CheckerAdapter adapts the internal checker to the validation.Checker interface
type CheckerAdapter struct {
	checker *checker.Checker
}

// convertConfig converts the external config to internal config format
func convertConfig(cfg *config.Config) *internalConfig.Config {
	return &internalConfig.Config{
		DryRun:            cfg.DryRun,
		Fix:               cfg.Fix,
		JSONOutput:        cfg.JSONOutput,
		RepoPath:          cfg.RepoPath,
		Interactive:       cfg.Interactive,
		CheckAugment:      cfg.CheckAugment,
		CheckDocker:       cfg.CheckDocker,
		CheckTypeScript:   cfg.CheckTypeScript,
		CheckDevContainer: cfg.CheckDevContainer,
		CheckDevEnv:       cfg.CheckDevEnv,
	}
}

// NewCheckerAdapter creates a new CheckerAdapter
func NewCheckerAdapter(cfg *config.Config) *CheckerAdapter {
	internalCfg := convertConfig(cfg)

	return &CheckerAdapter{
		checker: checker.NewChecker(internalCfg),
	}
}

// CheckFile checks if a file exists
func (a *CheckerAdapter) CheckFile(path string) (bool, error) {
	// Use the internal checker to check if the file exists
	result := a.checker.CheckFile(internalConfig.FileRequirement{
		Path: path,
	})

	if result.Error != nil {
		return result.Exists, fmt.Errorf("failed to check file %s: %w", path, result.Error)
	}
	return result.Exists, nil
}

// ReporterAdapter adapts the internal reporter to the validation.Reporter interface
type ReporterAdapter struct {
	reporter *reporter.Reporter
}

// NewReporterAdapter creates a new ReporterAdapter
func NewReporterAdapter(cfg *config.Config) *ReporterAdapter {
	internalCfg := convertConfig(cfg)

	return &ReporterAdapter{
		reporter: reporter.NewReporter(internalCfg),
	}
}

// convertResults converts validation results to internal results
func convertResults(results []validation.Result) []checker.ValidationResult {
	internalResults := make([]checker.ValidationResult, 0, len(results))

	for _, result := range results {
		internalResults = append(internalResults, checker.ValidationResult{
			Requirement: internalConfig.FileRequirement{
				Path:     result.File.Path,
				Priority: result.File.Priority,
			},
			Exists: result.File.Exists,
			Error:  result.Error,
		})
	}

	return internalResults
}

// ReportResults reports the validation results
func (a *ReporterAdapter) ReportResults(results []validation.Result) error {
	// Convert the validation results to internal results
	internalResults := convertResults(results)

	// Use the internal reporter to report the results
	return a.reporter.ReportResults(internalResults)
}

// GetSummary returns a summary of the validation results
func (a *ReporterAdapter) GetSummary(results []validation.Result) string {
	// Convert the validation results to internal results
	internalResults := convertResults(results)

	// Use the internal reporter to get the summary
	return a.reporter.GetSummary(internalResults)
}

// GetExitCode returns the exit code based on the validation results
func (a *ReporterAdapter) GetExitCode(results []validation.Result) int {
	// Convert the validation results to internal results
	internalResults := convertResults(results)

	// Use the internal reporter to get the exit code
	return a.reporter.GetExitCode(internalResults)
}

// GeneratorAdapter adapts the domain.FileRepository to the validation.Generator interface
type GeneratorAdapter struct {
	repository domain.FileRepository
}

// NewGeneratorAdapter creates a new GeneratorAdapter
func NewGeneratorAdapter(repository domain.FileRepository) *GeneratorAdapter {
	return &GeneratorAdapter{
		repository: repository,
	}
}

// GenerateFile generates a file from a template
func (a *GeneratorAdapter) GenerateFile(path, template string, data interface{}) error {
	err := a.repository.Generate(path, template, data)
	if err != nil {
		return fmt.Errorf("failed to generate file %s from template %s: %w", path, template, err)
	}
	return nil
}
