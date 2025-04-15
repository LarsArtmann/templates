package config

// Config represents the configuration for the GitHub issue graph tool
type Config struct {
	// Repository owner and name (e.g., "LarsArtmann/templates")
	Repo string
	// Specific issue number to use as root (0 means all issues)
	IssueNumber int
	// Maximum depth of relationships to traverse
	Depth int
	// Output file path (empty means stdout)
	OutputPath string
	// GitHub token for authentication
	Token string
	// Layout style (tree, radial, etc.)
	Layout string
}

// IsRootIssueSpecified returns true if a specific root issue is specified
func (c *Config) IsRootIssueSpecified() bool {
	return c.IssueNumber > 0
}

// GetOwnerAndRepo splits the repo string into owner and repo parts
func (c *Config) GetOwnerAndRepo() (string, string) {
	// Find the first slash
	for i := 0; i < len(c.Repo); i++ {
		if c.Repo[i] == '/' {
			return c.Repo[:i], c.Repo[i+1:]
		}
	}
	return c.Repo, ""
}
