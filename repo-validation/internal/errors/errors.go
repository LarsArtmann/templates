package errors

import (
	"fmt"
)

// PathError represents an error related to path resolution
type PathError struct {
	Path string
	Err  error
}

func (e *PathError) Error() string {
	return fmt.Sprintf("path error for %s: %v", e.Path, e.Err)
}

// NewPathError creates a new PathError
func NewPathError(path string, err error) *PathError {
	return &PathError{
		Path: path,
		Err:  err,
	}
}

// FileAccessError represents an error related to file access or permissions
type FileAccessError struct {
	Path string
	Err  error
}

func (e *FileAccessError) Error() string {
	return fmt.Sprintf("file access error for %s: %v", e.Path, e.Err)
}

// NewFileAccessError creates a new FileAccessError
func NewFileAccessError(path string, err error) *FileAccessError {
	return &FileAccessError{
		Path: path,
		Err:  err,
	}
}

// InvalidConfigError represents an error related to invalid configuration
type InvalidConfigError struct {
	Message string
}

func (e *InvalidConfigError) Error() string {
	return fmt.Sprintf("invalid configuration: %s", e.Message)
}

// NewInvalidConfigError creates a new InvalidConfigError
func NewInvalidConfigError(message string) *InvalidConfigError {
	return &InvalidConfigError{
		Message: message,
	}
}

// MissingMustHaveFilesError represents an error related to missing must-have files
type MissingMustHaveFilesError struct {
	Summary string
}

func (e *MissingMustHaveFilesError) Error() string {
	return fmt.Sprintf("repository validation failed: %s", e.Summary)
}

// NewMissingMustHaveFilesError creates a new MissingMustHaveFilesError
func NewMissingMustHaveFilesError(summary string) *MissingMustHaveFilesError {
	return &MissingMustHaveFilesError{
		Summary: summary,
	}
}
