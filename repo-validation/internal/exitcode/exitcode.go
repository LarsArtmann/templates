package exitcode

// Exit codes for the repository validation script
const (
	// Success indicates that all must-have files are present
	Success = 0
	
	// GeneralError indicates a general error or failure
	GeneralError = 1
	
	// PathError indicates a path resolution error
	PathError = 2
	
	// MissingMustHaveFiles indicates that some must-have files are missing
	MissingMustHaveFiles = 3
	
	// InvalidConfig indicates invalid configuration options
	InvalidConfig = 4
	
	// FileAccessError indicates a file access or permission error
	FileAccessError = 5
)
