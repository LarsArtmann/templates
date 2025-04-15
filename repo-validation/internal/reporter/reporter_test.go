package reporter

import (
	"testing"

	"github.com/LarsArtmann/templates/repo-validation/internal/checker"
	"github.com/LarsArtmann/templates/repo-validation/internal/config"
	"github.com/LarsArtmann/templates/repo-validation/internal/exitcode"
)

func TestReporter_GetSummary(t *testing.T) {
	tests := []struct {
		name    string
		results []checker.ValidationResult
		want    string
	}{
		{
			name: "all files present",
			results: []checker.ValidationResult{
				{
					Requirement: config.FileRequirement{Path: "README.md", Priority: config.PriorityMustHave},
					Exists:      true,
				},
				{
					Requirement: config.FileRequirement{Path: "LICENSE.md", Priority: config.PriorityMustHave},
					Exists:      true,
				},
			},
			want: "All must-have files are present",
		},
		{
			name: "missing must-have files",
			results: []checker.ValidationResult{
				{
					Requirement: config.FileRequirement{Path: "README.md", Priority: config.PriorityMustHave},
					Exists:      false,
				},
				{
					Requirement: config.FileRequirement{Path: "LICENSE.md", Priority: config.PriorityMustHave},
					Exists:      true,
				},
			},
			want: "Some must-have files are missing. Missing must-have files: README.md",
		},
		{
			name: "missing should-have files",
			results: []checker.ValidationResult{
				{
					Requirement: config.FileRequirement{Path: "README.md", Priority: config.PriorityMustHave},
					Exists:      true,
				},
				{
					Requirement: config.FileRequirement{Path: "CONTRIBUTING.md", Priority: config.PriorityShouldHave},
					Exists:      false,
				},
			},
			want: "All must-have files are present. Missing should-have files: CONTRIBUTING.md",
		},
		{
			name: "errors",
			results: []checker.ValidationResult{
				{
					Requirement: config.FileRequirement{Path: "README.md", Priority: config.PriorityMustHave},
					Exists:      true,
					Error:       nil,
				},
				{
					Requirement: config.FileRequirement{Path: "LICENSE.md", Priority: config.PriorityMustHave},
					Exists:      false,
					Error:       nil,
				},
				{
					Requirement: config.FileRequirement{Path: "SECURITY.md", Priority: config.PriorityMustHave},
					Error:       &testError{message: "test error"},
				},
			},
			want: "Some must-have files are missing. Missing must-have files: LICENSE.md. Errors: SECURITY.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Reporter{
				Config: &config.Config{},
			}
			if got := r.GetSummary(tt.results); got != tt.want {
				t.Errorf("Reporter.GetSummary() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestReporter_ProcessResults(t *testing.T) {
	tests := []struct {
		name                string
		results             []checker.ValidationResult
		wantMissingMustHave []string
		wantMissingShouldHave []string
		wantErrors          []string
	}{
		{
			name: "all files present",
			results: []checker.ValidationResult{
				{
					Requirement: config.FileRequirement{Path: "README.md", Priority: config.PriorityMustHave},
					Exists:      true,
				},
				{
					Requirement: config.FileRequirement{Path: "LICENSE.md", Priority: config.PriorityMustHave},
					Exists:      true,
				},
			},
			wantMissingMustHave: []string{},
			wantMissingShouldHave: []string{},
			wantErrors:          []string{},
		},
		{
			name: "missing files and errors",
			results: []checker.ValidationResult{
				{
					Requirement: config.FileRequirement{Path: "README.md", Priority: config.PriorityMustHave},
					Exists:      false,
				},
				{
					Requirement: config.FileRequirement{Path: "CONTRIBUTING.md", Priority: config.PriorityShouldHave},
					Exists:      false,
				},
				{
					Requirement: config.FileRequirement{Path: "SECURITY.md", Priority: config.PriorityMustHave},
					Error:       &testError{message: "test error"},
				},
			},
			wantMissingMustHave: []string{"README.md"},
			wantMissingShouldHave: []string{"CONTRIBUTING.md"},
			wantErrors:          []string{"SECURITY.md: test error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Reporter{
				Config: &config.Config{},
			}
			gotMissingMustHave, gotMissingShouldHave, gotErrors := r.processResults(tt.results)

			// Check missing must-have files
			if !stringSlicesEqual(gotMissingMustHave, tt.wantMissingMustHave) {
				t.Errorf("Reporter.processResults() missingMustHave = %v, want %v", gotMissingMustHave, tt.wantMissingMustHave)
			}

			// Check missing should-have files
			if !stringSlicesEqual(gotMissingShouldHave, tt.wantMissingShouldHave) {
				t.Errorf("Reporter.processResults() missingShouldHave = %v, want %v", gotMissingShouldHave, tt.wantMissingShouldHave)
			}

			// Check errors
			if !stringSlicesEqual(gotErrors, tt.wantErrors) {
				t.Errorf("Reporter.processResults() errors = %v, want %v", gotErrors, tt.wantErrors)
			}
		})
	}
}

// Helper function to compare string slices
func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestReporter_ShouldExitWithError(t *testing.T) {
	tests := []struct {
		name    string
		results []checker.ValidationResult
		want    bool
	}{
		{
			name: "all files present",
			results: []checker.ValidationResult{
				{
					Requirement: config.FileRequirement{Path: "README.md", Priority: config.PriorityMustHave},
					Exists:      true,
				},
				{
					Requirement: config.FileRequirement{Path: "LICENSE.md", Priority: config.PriorityMustHave},
					Exists:      true,
				},
			},
			want: false,
		},
		{
			name: "missing must-have files",
			results: []checker.ValidationResult{
				{
					Requirement: config.FileRequirement{Path: "README.md", Priority: config.PriorityMustHave},
					Exists:      false,
				},
				{
					Requirement: config.FileRequirement{Path: "LICENSE.md", Priority: config.PriorityMustHave},
					Exists:      true,
				},
			},
			want: true,
		},
		{
			name: "errors",
			results: []checker.ValidationResult{
				{
					Requirement: config.FileRequirement{Path: "README.md", Priority: config.PriorityMustHave},
					Exists:      true,
					Error:       nil,
				},
				{
					Requirement: config.FileRequirement{Path: "SECURITY.md", Priority: config.PriorityMustHave},
					Error:       &testError{message: "test error"},
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Reporter{
				Config: &config.Config{},
			}
			if got := r.ShouldExitWithError(tt.results); got != tt.want {
				t.Errorf("Reporter.ShouldExitWithError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReporter_GetExitCode(t *testing.T) {
	tests := []struct {
		name    string
		results []checker.ValidationResult
		want    int
	}{
		{
			name: "all files present",
			results: []checker.ValidationResult{
				{
					Requirement: config.FileRequirement{Path: "README.md", Priority: config.PriorityMustHave},
					Exists:      true,
				},
				{
					Requirement: config.FileRequirement{Path: "LICENSE.md", Priority: config.PriorityMustHave},
					Exists:      true,
				},
			},
			want: exitcode.Success,
		},
		{
			name: "missing must-have files",
			results: []checker.ValidationResult{
				{
					Requirement: config.FileRequirement{Path: "README.md", Priority: config.PriorityMustHave},
					Exists:      false,
				},
				{
					Requirement: config.FileRequirement{Path: "LICENSE.md", Priority: config.PriorityMustHave},
					Exists:      true,
				},
			},
			want: exitcode.MissingMustHaveFiles,
		},
		{
			name: "errors",
			results: []checker.ValidationResult{
				{
					Requirement: config.FileRequirement{Path: "README.md", Priority: config.PriorityMustHave},
					Exists:      true,
					Error:       nil,
				},
				{
					Requirement: config.FileRequirement{Path: "SECURITY.md", Priority: config.PriorityMustHave},
					Error:       &testError{message: "test error"},
				},
			},
			want: exitcode.GeneralError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Reporter{
				Config: &config.Config{},
			}
			if got := r.GetExitCode(tt.results); got != tt.want {
				t.Errorf("Reporter.GetExitCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

// testError is a simple error implementation for testing
type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}
