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

	// Run the application
	if err := run(*dryRun, *fix, *jsonOutput, *repoPath); err != nil {
		if *jsonOutput {
			// Output error in JSON format
			fmt.Printf("{\"error\": \"%s\"}\n", err.Error())
		} else {
			// Output error in human-readable format
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		os.Exit(1)
	}
}

// run executes the main application logic
func run(dryRun, fix, jsonOutput bool, repoPath string) error {
	// Create the configuration
	cfg := &config.Config{
		DryRun:     dryRun,
		Fix:        fix,
		JSONOutput: jsonOutput,
		RepoPath:   repoPath,
	}

	// Resolve the repository path to an absolute path
	absPath, err := filepath.Abs(cfg.RepoPath)
	if err != nil {
		return fmt.Errorf("failed to resolve repository path: %w", err)
	}

	// Check if the path exists and is a directory
	stat, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("failed to access repository path: %w", err)
	}
	if !stat.IsDir() {
		return fmt.Errorf("repository path is not a directory: %s", absPath)
	}
	cfg.RepoPath = absPath

	// Create the validate command
	validateCmd := cmd.NewValidateCommand(cfg)

	// Execute the command
	return validateCmd.Execute()
}
