package exitcode

import (
	"testing"
)

func TestExitCodes(t *testing.T) {
	// Test that exit codes are unique
	exitCodes := map[int]string{
		Success:           "Success",
		GeneralError:      "GeneralError",
		PathError:         "PathError",
		MissingMustHaveFiles: "MissingMustHaveFiles",
		InvalidConfig:     "InvalidConfig",
		FileAccessError:   "FileAccessError",
	}

	// Check for uniqueness
	if len(exitCodes) != 6 {
		t.Errorf("Expected 6 unique exit codes, got %d", len(exitCodes))
	}

	// Check specific values
	if Success != 0 {
		t.Errorf("Expected Success to be 0, got %d", Success)
	}

	if GeneralError != 1 {
		t.Errorf("Expected GeneralError to be 1, got %d", GeneralError)
	}

	// Check that exit codes are in ascending order
	if !(Success < GeneralError && 
		GeneralError < PathError && 
		PathError < MissingMustHaveFiles && 
		MissingMustHaveFiles < InvalidConfig && 
		InvalidConfig < FileAccessError) {
		t.Errorf("Exit codes are not in ascending order")
	}
}
