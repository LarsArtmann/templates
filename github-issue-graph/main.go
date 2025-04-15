package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/LarsArtmann/mono/public/templates/github-issue-graph/internal/config"
	"github.com/LarsArtmann/mono/public/templates/github-issue-graph/internal/graph"
	"github.com/LarsArtmann/mono/public/templates/github-issue-graph/internal/github"
	"github.com/LarsArtmann/mono/public/templates/github-issue-graph/internal/svg"
)

// Version is the current version of the GitHub issue graph tool
const Version = "0.1.0"

func main() {
	// Parse command-line flags
	version := flag.Bool("version", false, "Show version information")
	repo := flag.String("repo", "", "GitHub repository in owner/repo format")
	issueNumber := flag.Int("issue", 0, "Root issue number (0 means all issues)")
	depth := flag.Int("depth", 10, "Maximum depth of relationships to traverse")
	output := flag.String("output", "", "Output file path (empty means stdout)")
	layout := flag.String("layout", "tree", "Layout style (tree, radial)")

	flag.Parse()

	// If --version is specified, print version information and exit
	if *version {
		fmt.Printf("GitHub Issue Graph v%s\n", Version)
		os.Exit(0)
	}

	// Validate required parameters
	if *repo == "" {
		fmt.Fprintln(os.Stderr, "Error: --repo parameter is required")
		flag.Usage()
		os.Exit(1)
	}

	// Get GitHub token from environment
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Fprintln(os.Stderr, "Warning: GITHUB_TOKEN environment variable not set. Using unauthenticated API (rate limited)")
	}

	// Create configuration
	cfg := &config.Config{
		Repo:        *repo,
		IssueNumber: *issueNumber,
		Depth:       *depth,
		OutputPath:  *output,
		Token:       token,
		Layout:      *layout,
	}

	// Create GitHub client
	client, err := github.NewClient(cfg.Token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating GitHub client: %v\n", err)
		os.Exit(1)
	}

	// Fetch issues and build graph
	issueGraph, err := graph.BuildGraph(client, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building graph: %v\n", err)
		os.Exit(1)
	}

	// Generate SVG
	svgContent, err := svg.Generate(issueGraph, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating SVG: %v\n", err)
		os.Exit(1)
	}

	// Output the SVG
	if cfg.OutputPath == "" {
		// Write to stdout
		fmt.Println(svgContent)
	} else {
		// Write to file
		err = os.WriteFile(cfg.OutputPath, []byte(svgContent), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("SVG graph written to %s\n", cfg.OutputPath)
	}
}
