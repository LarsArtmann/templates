package reporter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/LarsArtmann/templates/repo-validation/internal/checker"
	"github.com/LarsArtmann/templates/repo-validation/internal/config"
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

// reportResultsConsole reports the validation results to the console
func (r *Reporter) reportResultsConsole(results []checker.ValidationResult) error {
	var missingMustHave, missingShouldHave []string
	var errors []string

	for _, result := range results {
		if result.Error != nil {
			errors = append(errors, fmt.Sprintf("%s: %s", result.Requirement.Path, result.Error))
			continue
		}

		if !result.Exists {
			if result.Requirement.Priority == "Must-have" {
				missingMustHave = append(missingMustHave, result.Requirement.Path)
			} else if result.Requirement.Priority == "Should-have" {
				missingShouldHave = append(missingShouldHave, result.Requirement.Path)
			}
		}
	}

	// Print summary
	fmt.Println("Repository Validation Results:")
	fmt.Println("=============================")

	if len(missingMustHave) == 0 && len(errors) == 0 {
		fmt.Println("\033[32m✓ All must-have files are present\033[0m")
	} else {
		fmt.Println("\033[31m✗ Some must-have files are missing\033[0m")
	}

	// Print missing must-have files
	if len(missingMustHave) > 0 {
		fmt.Println("\n\033[31mMissing must-have files:\033[0m")
		for _, file := range missingMustHave {
			fmt.Printf("  - %s\n", file)
		}
	}

	// Print missing should-have files
	if len(missingShouldHave) > 0 {
		fmt.Println("\n\033[33mMissing should-have files:\033[0m")
		for _, file := range missingShouldHave {
			fmt.Printf("  - %s\n", file)
		}
	}

	// Print errors
	if len(errors) > 0 {
		fmt.Println("\n\033[31mErrors:\033[0m")
		for _, err := range errors {
			fmt.Printf("  - %s\n", err)
		}
	}

	// Print fix message
	if (len(missingMustHave) > 0 || len(missingShouldHave) > 0) && !r.Config.Fix {
		fmt.Println("\nRun with --fix to generate missing files")
	}

	return nil
}

// reportResultsJSON reports the validation results in JSON format
func (r *Reporter) reportResultsJSON(results []checker.ValidationResult) error {
	var missingMustHave, missingShouldHave []string
	var errors []string

	for _, result := range results {
		if result.Error != nil {
			errors = append(errors, fmt.Sprintf("%s: %s", result.Requirement.Path, result.Error))
			continue
		}

		if !result.Exists {
			if result.Requirement.Priority == "Must-have" {
				missingMustHave = append(missingMustHave, result.Requirement.Path)
			} else if result.Requirement.Priority == "Should-have" {
				missingShouldHave = append(missingShouldHave, result.Requirement.Path)
			}
		}
	}

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
	var missingMustHave, missingShouldHave []string
	var errors []string

	for _, result := range results {
		if result.Error != nil {
			errors = append(errors, result.Requirement.Path)
			continue
		}

		if !result.Exists {
			if result.Requirement.Priority == "Must-have" {
				missingMustHave = append(missingMustHave, result.Requirement.Path)
			} else if result.Requirement.Priority == "Should-have" {
				missingShouldHave = append(missingShouldHave, result.Requirement.Path)
			}
		}
	}

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
