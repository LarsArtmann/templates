package validation

import (
	"errors"
	"testing"

	"github.com/LarsArtmann/templates/repo-validation/domain"
)

// MockChecker is a mock implementation of Checker for testing
type MockChecker struct {
	checkFileFunc func(path string) (bool, error)
}

func (m *MockChecker) CheckFile(path string) (bool, error) {
	return m.checkFileFunc(path)
}

// MockReporter is a mock implementation of Reporter for testing
type MockReporter struct {
	reportResultsFunc func(results []Result) error
	getSummaryFunc    func(results []Result) string
	getExitCodeFunc   func(results []Result) int
}

func (m *MockReporter) ReportResults(results []Result) error {
	return m.reportResultsFunc(results)
}

func (m *MockReporter) GetSummary(results []Result) string {
	return m.getSummaryFunc(results)
}

func (m *MockReporter) GetExitCode(results []Result) int {
	return m.getExitCodeFunc(results)
}

// MockGenerator is a mock implementation of Generator for testing
type MockGenerator struct {
	generateFileFunc func(path, template string, data interface{}) error
}

func (m *MockGenerator) GenerateFile(path, template string, data interface{}) error {
	return m.generateFileFunc(path, template, data)
}

func TestValidationService_ValidateFiles(t *testing.T) {
	tests := []struct {
		name         string
		files        []domain.File
		checkFileFunc func(path string) (bool, error)
		wantResults  int
		wantErrors   int
	}{
		{
			name: "all files exist",
			files: []domain.File{
				{Path: "file1.txt"},
				{Path: "file2.txt"},
			},
			checkFileFunc: func(path string) (bool, error) {
				return true, nil
			},
			wantResults: 2,
			wantErrors:  0,
		},
		{
			name: "some files don't exist",
			files: []domain.File{
				{Path: "file1.txt"},
				{Path: "file2.txt"},
			},
			checkFileFunc: func(path string) (bool, error) {
				return path == "file1.txt", nil
			},
			wantResults: 2,
			wantErrors:  0,
		},
		{
			name: "error checking files",
			files: []domain.File{
				{Path: "file1.txt"},
				{Path: "file2.txt"},
			},
			checkFileFunc: func(path string) (bool, error) {
				if path == "file1.txt" {
					return false, errors.New("test error")
				}
				return true, nil
			},
			wantResults: 2,
			wantErrors:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checker := &MockChecker{
				checkFileFunc: tt.checkFileFunc,
			}
			
			reporter := &MockReporter{}
			generator := &MockGenerator{}
			
			service := NewValidationService(checker, reporter, generator)
			
			results := service.ValidateFiles(tt.files)
			
			if len(results) != tt.wantResults {
				t.Errorf("ValidationService.ValidateFiles() results = %v, want %v", len(results), tt.wantResults)
			}
			
			errorCount := 0
			for _, result := range results {
				if result.Error != nil {
					errorCount++
				}
			}
			
			if errorCount != tt.wantErrors {
				t.Errorf("ValidationService.ValidateFiles() errors = %v, want %v", errorCount, tt.wantErrors)
			}
		})
	}
}

func TestValidationService_GenerateMissingFiles(t *testing.T) {
	tests := []struct {
		name            string
		results         []Result
		generateFileFunc func(path, template string, data interface{}) error
		wantErr         bool
	}{
		{
			name: "generate all files successfully",
			results: []Result{
				{
					File: domain.File{
						Path:     "file1.txt",
						Exists:   false,
						Template: "template1.tmpl",
					},
				},
				{
					File: domain.File{
						Path:     "file2.txt",
						Exists:   false,
						Template: "template2.tmpl",
					},
				},
			},
			generateFileFunc: func(path, template string, data interface{}) error {
				return nil
			},
			wantErr: false,
		},
		{
			name: "skip existing files",
			results: []Result{
				{
					File: domain.File{
						Path:     "file1.txt",
						Exists:   true,
						Template: "template1.tmpl",
					},
				},
				{
					File: domain.File{
						Path:     "file2.txt",
						Exists:   false,
						Template: "template2.tmpl",
					},
				},
			},
			generateFileFunc: func(path, template string, data interface{}) error {
				return nil
			},
			wantErr: false,
		},
		{
			name: "skip files without templates",
			results: []Result{
				{
					File: domain.File{
						Path:     "file1.txt",
						Exists:   false,
						Template: "",
					},
				},
				{
					File: domain.File{
						Path:     "file2.txt",
						Exists:   false,
						Template: "template2.tmpl",
					},
				},
			},
			generateFileFunc: func(path, template string, data interface{}) error {
				return nil
			},
			wantErr: false,
		},
		{
			name: "error generating file",
			results: []Result{
				{
					File: domain.File{
						Path:     "file1.txt",
						Exists:   false,
						Template: "template1.tmpl",
					},
				},
			},
			generateFileFunc: func(path, template string, data interface{}) error {
				return errors.New("test error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checker := &MockChecker{}
			reporter := &MockReporter{}
			generator := &MockGenerator{
				generateFileFunc: tt.generateFileFunc,
			}
			
			service := NewValidationService(checker, reporter, generator)
			
			err := service.GenerateMissingFiles(tt.results)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidationService.GenerateMissingFiles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
