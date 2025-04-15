package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

	// If --all is specified, enable all optional file groups
	if *checkAll {
		*checkAugment = true
		*checkDocker = true
		*checkTypeScript = true
		*checkDevContainer = true
		*checkDevEnv = true
	}

	// Run the application
	if err := run(*dryRun, *fix, *jsonOutput, *repoPath, *interactive, *checkAugment, *checkDocker, *checkTypeScript, *checkDevContainer, *checkDevEnv); err != nil {
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

// promptForMissingParameters prompts the user for missing parameters
func promptForMissingParameters(cfg *config.Config) error {
	reader := bufio.NewReader(os.Stdin)

	// If DryRun and Fix are both true, prompt to resolve the conflict
	if cfg.DryRun && cfg.Fix {
		fmt.Println("--dry-run and --fix cannot be used together. Which one do you want to use?")
		fmt.Println("1. Dry Run (only report issues)")
		fmt.Println("2. Fix (generate missing files)")
		fmt.Print("Enter your choice (1 or 2): ")

		choice, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading input: %w", err)
		}

		choice = strings.TrimSpace(choice)
		switch choice {
		case "1":
			cfg.DryRun = true
			cfg.Fix = false
		case "2":
			cfg.DryRun = false
			cfg.Fix = true
		default:
			return fmt.Errorf("invalid choice: %s", choice)
		}
	} else if !cfg.Fix && !cfg.DryRun {
		// If neither Fix nor DryRun is set, prompt for Fix
		fmt.Print("Do you want to fix missing files? (y/n): ")

		choice, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading input: %w", err)
		}

		choice = strings.TrimSpace(strings.ToLower(choice))
		if choice == "y" || choice == "yes" {
			cfg.Fix = true
		}
	}

	// If no file groups are selected, prompt for which ones to check
	if !cfg.CheckAugment && !cfg.CheckDocker && !cfg.CheckTypeScript && !cfg.CheckDevContainer && !cfg.CheckDevEnv {
		fmt.Println("Which file groups do you want to check?")
		fmt.Println("1. Augment AI files (.augment-guidelines, .augmentignore)")
		fmt.Println("2. Docker files (Dockerfile, docker-compose.yaml, .dockerignore)")
		fmt.Println("3. TypeScript/JavaScript files (package.json, tsconfig.json)")
		fmt.Println("4. DevContainer files (.devcontainer.json)")
		fmt.Println("5. DevEnv files (devenv.nix)")
		fmt.Println("6. All file groups")
		fmt.Print("Enter your choices (comma-separated, e.g., 1,3,5): ")

		choice, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading input: %w", err)
		}

		choice = strings.TrimSpace(choice)
		choices := strings.Split(choice, ",")

		for _, c := range choices {
			c = strings.TrimSpace(c)
			switch c {
			case "1":
				cfg.CheckAugment = true
			case "2":
				cfg.CheckDocker = true
			case "3":
				cfg.CheckTypeScript = true
			case "4":
				cfg.CheckDevContainer = true
			case "5":
				cfg.CheckDevEnv = true
			case "6":
				cfg.CheckAugment = true
				cfg.CheckDocker = true
				cfg.CheckTypeScript = true
				cfg.CheckDevContainer = true
				cfg.CheckDevEnv = true
			}
		}
	}

	return nil
}

// run executes the main application logic
func run(dryRun, fix, jsonOutput bool, repoPath string, interactive bool, checkAugment, checkDocker, checkTypeScript, checkDevContainer, checkDevEnv bool) error {
	// Create the configuration
	cfg := &config.Config{
		DryRun:          dryRun,
		Fix:             fix,
		JSONOutput:      jsonOutput,
		RepoPath:        repoPath,
		Interactive:     interactive,
		CheckAugment:    checkAugment,
		CheckDocker:     checkDocker,
		CheckTypeScript: checkTypeScript,
		CheckDevContainer: checkDevContainer,
		CheckDevEnv:     checkDevEnv,
	}

	// Validate the configuration
	if err := cfg.Validate(); err != nil {
		// If interactive mode is enabled, prompt for missing parameters
		if interactive {
			if err := promptForMissingParameters(cfg); err != nil {
				return err
			}
		} else {
			return err
		}
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
