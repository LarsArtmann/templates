package domain

import "errors"

var (
	// ErrRequiredFileMissing is returned when a required file is missing
	ErrRequiredFileMissing = errors.New("required file is missing")
	
	// ErrNoTemplate is returned when a file has no template
	ErrNoTemplate = errors.New("no template available for file")
)
