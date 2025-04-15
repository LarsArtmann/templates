package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestPathError(t *testing.T) {
	// Create a test error
	testErr := errors.New("test error")
	path := "/test/path"
	
	// Create a PathError
	pathErr := NewPathError(path, testErr)
	
	// Check that the error message is formatted correctly
	expected := fmt.Sprintf("path error for %s: %v", path, testErr)
	if pathErr.Error() != expected {
		t.Errorf("Expected error message %q, got %q", expected, pathErr.Error())
	}
	
	// Check that the path and error are stored correctly
	if pathErr.Path != path {
		t.Errorf("Expected path %q, got %q", path, pathErr.Path)
	}
	
	if pathErr.Err != testErr {
		t.Errorf("Expected error %v, got %v", testErr, pathErr.Err)
	}
}

func TestFileAccessError(t *testing.T) {
	// Create a test error
	testErr := errors.New("test error")
	path := "/test/path"
	
	// Create a FileAccessError
	fileErr := NewFileAccessError(path, testErr)
	
	// Check that the error message is formatted correctly
	expected := fmt.Sprintf("file access error for %s: %v", path, testErr)
	if fileErr.Error() != expected {
		t.Errorf("Expected error message %q, got %q", expected, fileErr.Error())
	}
	
	// Check that the path and error are stored correctly
	if fileErr.Path != path {
		t.Errorf("Expected path %q, got %q", path, fileErr.Path)
	}
	
	if fileErr.Err != testErr {
		t.Errorf("Expected error %v, got %v", testErr, fileErr.Err)
	}
}

func TestInvalidConfigError(t *testing.T) {
	// Create a test message
	message := "invalid configuration"
	
	// Create an InvalidConfigError
	configErr := NewInvalidConfigError(message)
	
	// Check that the error message is formatted correctly
	expected := fmt.Sprintf("invalid configuration: %s", message)
	if configErr.Error() != expected {
		t.Errorf("Expected error message %q, got %q", expected, configErr.Error())
	}
	
	// Check that the message is stored correctly
	if configErr.Message != message {
		t.Errorf("Expected message %q, got %q", message, configErr.Message)
	}
}

func TestMissingMustHaveFilesError(t *testing.T) {
	// Create a test summary
	summary := "missing files: file1.txt, file2.txt"
	
	// Create a MissingMustHaveFilesError
	missingErr := NewMissingMustHaveFilesError(summary)
	
	// Check that the error message is formatted correctly
	expected := fmt.Sprintf("repository validation failed: %s", summary)
	if missingErr.Error() != expected {
		t.Errorf("Expected error message %q, got %q", expected, missingErr.Error())
	}
	
	// Check that the summary is stored correctly
	if missingErr.Summary != summary {
		t.Errorf("Expected summary %q, got %q", summary, missingErr.Summary)
	}
}
