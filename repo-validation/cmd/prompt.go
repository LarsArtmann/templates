package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/LarsArtmann/templates/repo-validation/internal/config"
)

// PromptForMissingParameters prompts the user for missing parameters
func PromptForMissingParameters(cfg *config.Config) error {
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
		if choice == "" {
			return fmt.Errorf("no choice provided, please enter 1 or 2")
		}

		switch choice {
		case "1":
			cfg.DryRun = true
			cfg.Fix = false
		case "2":
			cfg.DryRun = false
			cfg.Fix = true
		default:
			return fmt.Errorf("invalid choice: %s (must be 1 or 2)", choice)
		}
	} else if !cfg.Fix && !cfg.DryRun {
		// If neither Fix nor DryRun is set, prompt for Fix
		fmt.Print("Do you want to fix missing files? (y/n): ")

		choice, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading input: %w", err)
		}

		choice = strings.TrimSpace(strings.ToLower(choice))
		if choice == "" {
			// Default to not fixing if no input is provided
			fmt.Println("No input provided, defaulting to no fix")
			cfg.Fix = false
		} else if choice == "y" || choice == "yes" {
			cfg.Fix = true
		} else if choice == "n" || choice == "no" {
			cfg.Fix = false
		} else {
			return fmt.Errorf("invalid choice: %s (must be y/yes or n/no)", choice)
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
		fmt.Println("7. None (only check core files)")
		fmt.Print("Enter your choices (comma-separated, e.g., 1,3,5): ")

		choice, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading input: %w", err)
		}

		choice = strings.TrimSpace(choice)
		if choice == "" {
			// Default to core files only if no input is provided
			fmt.Println("No input provided, defaulting to core files only")
			return nil
		}

		choices := strings.Split(choice, ",")
		validChoices := map[string]bool{"1": true, "2": true, "3": true, "4": true, "5": true, "6": true, "7": true}
		hasInvalidChoice := false
		invalidChoices := []string{}

		for _, c := range choices {
			c = strings.TrimSpace(c)
			if c == "" {
				continue
			}

			if !validChoices[c] {
				hasInvalidChoice = true
				invalidChoices = append(invalidChoices, c)
				continue
			}

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
				// All file groups
				cfg.CheckAugment = true
				cfg.CheckDocker = true
				cfg.CheckTypeScript = true
				cfg.CheckDevContainer = true
				cfg.CheckDevEnv = true
			case "7":
				// None (only check core files)
				// This is already the default, so we don't need to do anything
			}
		}

		if hasInvalidChoice {
			return fmt.Errorf("invalid choices: %s (must be numbers between 1-7)", strings.Join(invalidChoices, ", "))
		}
	}

	return nil
}
