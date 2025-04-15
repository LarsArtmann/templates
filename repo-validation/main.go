package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/LarsArtmann/templates/repo-validation/cmd"
	"github.com/LarsArtmann/templates/repo-validation/internal/config"
)

// Version is the current version of the repository validation script
const Version = "0.1.0"

func main() {
	// Parse command-line flags
	version := flag.Bool("version", false, "Show version information")
	dryRun := flag.Bool("dry-run", false, "Only report issues without making changes")
	fix := flag.Bool("fix", false, "Generate missing files")
	jsonOutput := flag.Bool("json", false, "Output results in JSON format")
	repoPath := flag.String("path", ".", "Path to the repository to validate")
	interactive := flag.Bool("interactive", false, "Prompt for missing parameters")

	// Optional file group flags
	checkAugment := flag.Bool("augment", false, "Check Augment AI related files (.augment-guidelines, .augmentignore)")
	checkDocker := flag.Bool("docker", false, "Check Docker related files (Dockerfile, docker-compose.yaml, .dockerignore)")
	checkTypeScript := flag.Bool("typescript", false, "Check TypeScript/JavaScript related files (package.json, tsconfig.json)")
	checkDevContainer := flag.Bool("devcontainer", false, "Check DevContainer related files (.devcontainer.json)")
	checkDevEnv := flag.Bool("devenv", false, "Check DevEnv related files (devenv.nix)")
	checkAll := flag.Bool("all", false, "Check all optional file groups")

	flag.Parse()

	// If --version is specified, print version information and exit
	if *version {
		fmt.Printf("Repository Validation Script v%s\n", Version)
		os.Exit(0)
	}

	// Prepare options for the run function
	options := []config.ConfigOption{
		config.WithDryRun(*dryRun),
		config.WithFix(*fix),
		config.WithJSONOutput(*jsonOutput),
		config.WithRepoPath(*repoPath),
		config.WithInteractive(*interactive),
	}

	// If --all is specified, add the all file group option
	if *checkAll {
		options = append(options, config.WithFileGroup("all", true))
	} else {
		// Otherwise, add individual file group options
		if *checkAugment {
			options = append(options, config.WithFileGroup("augment", true))
		}
		if *checkDocker {
			options = append(options, config.WithFileGroup("docker", true))
		}
		if *checkTypeScript {
			options = append(options, config.WithFileGroup("typescript", true))
		}
		if *checkDevContainer {
			options = append(options, config.WithFileGroup("devcontainer", true))
		}
		if *checkDevEnv {
			options = append(options, config.WithFileGroup("devenv", true))
		}
	}

	// Run the application with the options
	if err := cmd.Run(options...); err != nil {
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


