package reporter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/LarsArtmann/templates/repo-validation/internal/checker"
	"github.com/LarsArtmann/templates/repo-validation/internal/config"
	"github.com/LarsArtmann/templates/repo-validation/internal/exitcode"
	"github.com/charmbracelet/log"
)

// Reporter is responsible for reporting validation results
type Reporter struct {
	// Config is the configuration for the reporter
	Config *config.Config
}

// NewReporter creates a new Reporter
func NewReporter(cfg *config.Config) *Reporter {
	return &Reporter{
		Config: cfg,
	}
}

// JSONResult represents the JSON output of the validation
type JSONResult struct {
	// Success indicates whether all required files exist
	Success bool `json:"success"`
	// MissingMustHaveFiles is the list of must-have files that are missing
	MissingMustHaveFiles []string `json:"missingMustHaveFiles,omitempty"`
	// MissingShouldHaveFiles is the list of should-have files that are missing
	MissingShouldHaveFiles []string `json:"missingShouldHaveFiles,omitempty"`
	// Errors is the list of errors that occurred during validation
	Errors []string `json:"errors,omitempty"`
}

// ReportResults reports the validation results
func (r *Reporter) ReportResults(results []checker.ValidationResult) error {
	if r.Config.JSONOutput {
		return r.reportResultsJSON(results)
	}

	return r.reportResultsConsole(results)
}

// processResults extracts information about missing files and errors from validation results
func (r *Reporter) processResults(results []checker.ValidationResult) (missingMustHave, missingShouldHave, errors []string) {
	for _, result := range results {
		if result.Error != nil {
			errors = append(errors, fmt.Sprintf("%s: %s", result.Requirement.Path, result.Error))
			continue
		}

		if !result.Exists {
			if result.Requirement.Priority == config.PriorityMustHave {
				missingMustHave = append(missingMustHave, result.Requirement.Path)
			} else if result.Requirement.Priority == config.PriorityShouldHave {
				missingShouldHave = append(missingShouldHave, result.Requirement.Path)
			}
		}
	}
	return missingMustHave, missingShouldHave, errors
}

// reportResultsConsole reports the validation results to the console
func (r *Reporter) reportResultsConsole(results []checker.ValidationResult) error {
	missingMustHave, missingShouldHave, errors := r.processResults(results)

	// Print summary
	log.Info("Repository Validation Results")
	log.Info("===========================")

	if len(missingMustHave) == 0 && len(errors) == 0 {
		log.Info("✓ All must-have files are present", "status", "success")
	} else {
		log.Error("✗ Some must-have files are missing", "status", "failed")
	}

	// Print missing must-have files
	if len(missingMustHave) > 0 {
		log.Error("Missing must-have files:")
		for _, file := range missingMustHave {
			log.Error("  - " + file)
		}
	}

	// Print missing should-have files
	if len(missingShouldHave) > 0 {
		log.Warn("Missing should-have files:")
		for _, file := range missingShouldHave {
			log.Warn("  - " + file)
		}
	}

	// Print errors
	if len(errors) > 0 {
		log.Error("Errors:")
		for _, err := range errors {
			log.Error("  - " + err)
		}
	}

	// Print fix message
	if (len(missingMustHave) > 0 || len(missingShouldHave) > 0) && !r.Config.Fix {
		log.Info("Run with --fix to generate missing files")
	}

	return nil
}

// reportResultsJSON reports the validation results in JSON format
func (r *Reporter) reportResultsJSON(results []checker.ValidationResult) error {
	missingMustHave, missingShouldHave, errors := r.processResults(results)

	jsonResult := JSONResult{
		Success:              len(missingMustHave) == 0 && len(errors) == 0,
		MissingMustHaveFiles: missingMustHave,
		MissingShouldHaveFiles: missingShouldHave,
		Errors:               errors,
	}

	jsonData, err := json.MarshalIndent(jsonResult, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	fmt.Println(string(jsonData))

	return nil
}

// GetSummary returns a summary of the validation results
func (r *Reporter) GetSummary(results []checker.ValidationResult) string {
	// Reuse the processResults function for consistency
	missingMustHave, missingShouldHave, errors := r.processResults(results)

	// Extract just the paths from error messages for simpler output
	simplifiedErrors := make([]string, len(errors))
	for i, err := range errors {
		parts := strings.SplitN(err, ":", 2)
		simplifiedErrors[i] = parts[0]
	}
	errors = simplifiedErrors

	var summary strings.Builder

	if len(missingMustHave) == 0 && len(errors) == 0 {
		summary.WriteString("All must-have files are present")
	} else {
		summary.WriteString("Some must-have files are missing")
	}

	if len(missingMustHave) > 0 {
		summary.WriteString(fmt.Sprintf(". Missing must-have files: %s", strings.Join(missingMustHave, ", ")))
	}

	if len(missingShouldHave) > 0 {
		summary.WriteString(fmt.Sprintf(". Missing should-have files: %s", strings.Join(missingShouldHave, ", ")))
	}

	if len(errors) > 0 {
		summary.WriteString(fmt.Sprintf(". Errors: %s", strings.Join(errors, ", ")))
	}

	return summary.String()
}

// ShouldExitWithError returns true if there are missing must-have files or errors
// This helps the caller determine the correct exit code
func (r *Reporter) ShouldExitWithError(results []checker.ValidationResult) bool {
	missingMustHave, _, errors := r.processResults(results)
	return len(missingMustHave) > 0 || len(errors) > 0
}

// GetExitCode returns the appropriate exit code based on the validation results
func (r *Reporter) GetExitCode(results []checker.ValidationResult) int {
	missingMustHave, _, errors := r.processResults(results)

	if len(errors) > 0 {
		return exitcode.GeneralError
	}

	if len(missingMustHave) > 0 {
		return exitcode.MissingMustHaveFiles
	}

	return exitcode.Success
}
