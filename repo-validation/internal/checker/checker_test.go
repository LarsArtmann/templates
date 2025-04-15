package checker

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/LarsArtmann/templates/repo-validation/internal/config"
)

// setupTestDir creates a temporary directory with test files
func setupTestDir(t *testing.T) string {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "repo-validation-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	// Create some test files
	files := []string{
		"README.md",
		"LICENSE.md",
		".gitignore",
	}

	for _, file := range files {
		path := filepath.Join(tempDir, file)
		if err := os.WriteFile(path, []byte("test content"), 0644); err != nil {
			os.RemoveAll(tempDir)
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	return tempDir
}

// cleanupTestDir removes the temporary directory
func cleanupTestDir(dir string) {
	os.RemoveAll(dir)
}

func TestNewChecker(t *testing.T) {
	cfg := &config.Config{
		RepoPath: "/test/path",
	}

	chk := NewChecker(cfg)
	if chk == nil {
		t.Errorf("Expected non-nil Checker, got nil")
	}

	if chk.Config != cfg {
		t.Errorf("Expected Config to be %v, got %v", cfg, chk.Config)
	}
}

func TestCheckRepository(t *testing.T) {
	// Setup test directory
	tempDir := setupTestDir(t)
	defer cleanupTestDir(tempDir)

	// Create a checker with the test directory
	cfg := &config.Config{
		RepoPath: tempDir,
	}
	chk := NewChecker(cfg)

	// Check the repository
	results, err := chk.CheckRepository()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check that we have results
	if len(results) == 0 {
		t.Errorf("Expected non-empty results, got empty results")
	}

	// Check that existing files are marked as existing
	for _, result := range results {
		if result.Requirement.Path == "README.md" ||
		   result.Requirement.Path == "LICENSE.md" ||
		   result.Requirement.Path == ".gitignore" {
			if !result.Exists {
				t.Errorf("Expected file %s to exist, but it was marked as not existing", result.Requirement.Path)
			}
		}
	}
}

func TestFixMissingFiles(t *testing.T) {
	// This test is more of an integration test and requires actual templates
	// For unit testing, we'll mock the behavior instead

	// Setup test directory
	tempDir := setupTestDir(t)
	defer cleanupTestDir(tempDir)

	// Create a checker with the test directory
	cfg := &config.Config{
		RepoPath: tempDir,
		Fix:      true,
	}
	chk := NewChecker(cfg)

	// Create a test result with a missing file but no template
	// This should be skipped without error
	results := []ValidationResult{
		{
			Requirement: config.FileRequirement{
				Path:         "CONTRIBUTING.md",
				Priority:     config.PriorityMustHave,
				TemplatePath: "", // No template
			},
			Exists: false,
		},
	}

	// Fix the missing files
	err := chk.FixMissingFiles(results)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// The file should not be created since there's no template
	contributingPath := filepath.Join(tempDir, "CONTRIBUTING.md")
	if _, err := os.Stat(contributingPath); !os.IsNotExist(err) {
		t.Errorf("Expected file %s to not exist, but it does", contributingPath)
	}
}

func TestCheckFile(t *testing.T) {
	// Setup test directory
	tempDir := setupTestDir(t)
	defer cleanupTestDir(tempDir)

	// Create a checker with the test directory
	cfg := &config.Config{
		RepoPath: tempDir,
	}
	chk := NewChecker(cfg)

	// Test existing file
	t.Run("existing file", func(t *testing.T) {
		req := config.FileRequirement{
			Path:     "README.md",
			Priority: config.PriorityMustHave,
		}

		result := chk.checkFile(req)
		if !result.Exists {
			t.Errorf("Expected file to exist, got not existing")
		}
		if result.Error != nil {
			t.Errorf("Expected no error, got %v", result.Error)
		}
		if result.Requirement.Path != req.Path {
			t.Errorf("Expected path %s, got %s", req.Path, result.Requirement.Path)
		}
	})

	// Test non-existing file
	t.Run("non-existing file", func(t *testing.T) {
		req := config.FileRequirement{
			Path:     "nonexistent.txt",
			Priority: config.PriorityMustHave,
		}

		result := chk.checkFile(req)
		if result.Exists {
			t.Errorf("Expected file to not exist, got existing")
		}
		if result.Error != nil {
			t.Errorf("Expected no error, got %v", result.Error)
		}
		if result.Requirement.Path != req.Path {
			t.Errorf("Expected path %s, got %s", req.Path, result.Requirement.Path)
		}
	})

	// We can't easily create an error case in a platform-independent way
	// since the behavior of os.Stat on directories varies by OS.
	// Instead, we'll test the normal cases thoroughly.
}
