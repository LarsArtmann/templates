package cmd

import (
	"fmt"
	"os"

	"github.com/LarsArtmann/templates/repo-validation/config"
	"github.com/LarsArtmann/templates/repo-validation/logger"
	"github.com/spf13/cobra"
)

var (
	// Configuration options
	dryRun          bool
	fix             bool
	jsonOutput      bool
	repoPath        string
	interactive     bool
	checkAll        bool
	checkAugment    bool
	checkDocker     bool
	checkTypeScript bool
	checkDevContainer bool
	checkDevEnv     bool
	
	// Version information
	version = "0.1.0"
	
	// Logger
	log = logger.NewCharmLogger()
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "repo-validate",
	Short: "Validates repository files",
	Long:  `A tool for validating the presence of common project files in a repository.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create configuration options
		options := []config.ConfigOption{
			config.WithDryRun(dryRun),
			config.WithFix(fix),
			config.WithJSONOutput(jsonOutput),
			config.WithRepoPath(repoPath),
			config.WithInteractive(interactive),
		}
		
		// Add file group options
		if checkAll {
			options = append(options, config.WithFileGroup("all", true))
		} else {
			if checkAugment {
				options = append(options, config.WithFileGroup("augment", true))
			}
			if checkDocker {
				options = append(options, config.WithFileGroup("docker", true))
			}
			if checkTypeScript {
				options = append(options, config.WithFileGroup("typescript", true))
			}
			if checkDevContainer {
				options = append(options, config.WithFileGroup("devcontainer", true))
			}
			if checkDevEnv {
				options = append(options, config.WithFileGroup("devenv", true))
			}
		}
		
		// Run the validation command
		return runValidate(options...)
	},
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("repo-validate version %s\n", version)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Add version command
	rootCmd.AddCommand(versionCmd)
	
	// Add flags
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Only report issues without making changes")
	rootCmd.Flags().BoolVar(&fix, "fix", false, "Generate missing files")
	rootCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output results in JSON format")
	rootCmd.Flags().StringVar(&repoPath, "path", ".", "Path to the repository")
	rootCmd.Flags().BoolVar(&interactive, "interactive", false, "Run in interactive mode")
	
	// File group flags
	rootCmd.Flags().BoolVar(&checkAll, "all", false, "Check all file groups")
	rootCmd.Flags().BoolVar(&checkAugment, "augment", false, "Check Augment AI related files")
	rootCmd.Flags().BoolVar(&checkDocker, "docker", false, "Check Docker related files")
	rootCmd.Flags().BoolVar(&checkTypeScript, "typescript", false, "Check TypeScript/JavaScript related files")
	rootCmd.Flags().BoolVar(&checkDevContainer, "devcontainer", false, "Check DevContainer related files")
	rootCmd.Flags().BoolVar(&checkDevEnv, "devenv", false, "Check DevEnv related files")
}
