# GitHub Issue Graph

A tool for visualizing parent-child relationships between GitHub issues. This tool generates SVG diagrams that show the hierarchical structure of issues in a GitHub repository.

## Features

- Generate SVG visualizations of GitHub issue relationships
- Color-code nodes based on issue labels (Epic, Workstream, etc.)
- Support for different layout styles (tree, radial)
- Include issue numbers, titles, and status in the visualization
- Clickable links to the original GitHub issues
- Support for authentication via GitHub token for private repositories

## Installation

### Using Go Install

```bash
go install github.com/LarsArtmann/mono/public/templates/github-issue-graph@latest
```

### From Source

```bash
# Clone the repository
git clone https://github.com/LarsArtmann/mono.git
cd mono/public/templates/github-issue-graph

# Build the tool
go build -o github-issue-graph

# Optional: Install to your PATH
sudo mv github-issue-graph /usr/local/bin/
```

## Usage

```bash
# Set your GitHub token (optional, but recommended to avoid rate limits)
export GITHUB_TOKEN=your_github_token

# Generate SVG for all issues in a repository
github-issue-graph --repo="LarsArtmann/templates" --output="issues-graph.svg"

# Generate SVG for a specific epic and its children
github-issue-graph --repo="LarsArtmann/templates" --issue=42 --output="epic-42.svg"

# Limit depth of relationships
github-issue-graph --repo="LarsArtmann/templates" --issue=42 --depth=2 --output="epic-42-depth2.svg"

# Use a different layout style
github-issue-graph --repo="LarsArtmann/templates" --layout="radial" --output="radial-layout.svg"
```

### Options

- `--repo`: GitHub repository in owner/repo format (required)
- `--issue`: Root issue number to visualize (0 means all issues, default: 0)
- `--depth`: Maximum depth of relationships to traverse (default: 10)
- `--output`: Output file path (empty means stdout)
- `--layout`: Layout style (tree, radial, default: tree)
- `--version`: Show version information and exit

## How It Works

1. The tool connects to the GitHub API and fetches issue data
2. It analyzes the relationships between issues using timeline events and references
3. It builds a graph data structure representing the issue hierarchy
4. It generates an SVG visualization of the graph with appropriate styling

## Authentication

For private repositories or to avoid API rate limits, set your GitHub token as an environment variable:

```bash
export GITHUB_TOKEN=your_github_token
```

You can create a personal access token with the `repo` scope at [GitHub Developer Settings](https://github.com/settings/tokens).

## Integration with CI/CD

You can integrate this tool with GitHub Actions to automatically generate and publish relationship diagrams:

```yaml
name: Generate Issue Graph

on:
  issues:
    types: [opened, edited, closed, reopened]
  schedule:
    - cron: '0 0 * * *'  # Daily at midnight

jobs:
  generate-graph:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
          
      - name: Install github-issue-graph
        run: go install github.com/LarsArtmann/mono/public/templates/github-issue-graph@latest
        
      - name: Generate issue graph
        run: github-issue-graph --repo="${{ github.repository }}" --output="issue-graph.svg"
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          
      - name: Commit and push
        run: |
          git config --global user.name 'GitHub Actions'
          git config --global user.email 'actions@github.com'
          git add issue-graph.svg
          git commit -m "Update issue graph" || echo "No changes to commit"
          git push
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
