package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/LarsArtmann/templates/repo-validation/cmd"
	"github.com/LarsArtmann/templates/repo-validation/internal/config"
)

func main() {
	// Parse command-line flags
	dryRun := flag.Bool("dry-run", false, "Only report issues without making changes")
	fix := flag.Bool("fix", false, "Generate missing files")
	jsonOutput := flag.Bool("json", false, "Output results in JSON format")
	repoPath := flag.String("path", ".", "Path to the repository to validate")
	
	flag.Parse()
	
	// Create the configuration
	cfg := &config.Config{
		DryRun:     *dryRun,
		Fix:        *fix,
		JSONOutput: *jsonOutput,
		RepoPath:   *repoPath,
	}
	
	// Resolve the repository path to an absolute path
	absPath, err := filepath.Abs(cfg.RepoPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving repository path: %v\n", err)
		os.Exit(1)
	}
	// Check if the path exists and is a directory
	stat, err := os.Stat(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error accessing repository path: %v\n", err)
		os.Exit(1)
	}
	if !stat.IsDir() {
		fmt.Fprintf(os.Stderr, "Repository path is not a directory: %s\n", absPath)
		os.Exit(1)
	}
	cfg.RepoPath = absPath
	
	// Create the validate command
	validateCmd := cmd.NewValidateCommand(cfg)
	
	// Execute the command
	if err := validateCmd.Execute(); err != nil {
		if !cfg.JSONOutput {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		os.Exit(1)
	}
}
