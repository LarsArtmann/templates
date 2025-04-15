package graph

import (
	"fmt"
	"strings"

	"github.com/LarsArtmann/mono/public/templates/github-issue-graph/internal/config"
	"github.com/LarsArtmann/mono/public/templates/github-issue-graph/internal/github"
)

// Node represents a node in the issue graph
type Node struct {
	ID       int
	Title    string
	URL      string
	State    string
	Color    string
	Parents  []*Node
	Children []*Node
}

// Graph represents the entire issue relationship graph
type Graph struct {
	Nodes map[int]*Node
	Root  *Node
}

// BuildGraph fetches issues from GitHub and builds the relationship graph
func BuildGraph(client *github.Client, cfg *config.Config) (*Graph, error) {
	owner, repo := cfg.GetOwnerAndRepo()
	if repo == "" {
		return nil, fmt.Errorf("invalid repository format: %s (expected owner/repo)", cfg.Repo)
	}

	// Create a new graph
	graph := &Graph{
		Nodes: make(map[int]*Node),
	}

	var issues []*github.Issue
	var err error

	// Fetch issues
	if cfg.IsRootIssueSpecified() {
		// Fetch specific issue
		issue, err := client.GetIssue(owner, repo, cfg.IssueNumber)
		if err != nil {
			return nil, err
		}
		issues = []*github.Issue{issue}
	} else {
		// Fetch all issues
		issues, err = client.GetAllIssues(owner, repo)
		if err != nil {
			return nil, err
		}
	}

	// Find relationships between issues
	err = client.FindIssueRelationships(owner, repo, issues)
	if err != nil {
		return nil, err
	}

	// Convert issues to nodes
	for _, issue := range issues {
		node := &Node{
			ID:    issue.Number,
			Title: issue.Title,
			URL:   issue.URL,
			State: issue.State,
			Color: determineNodeColor(issue.Labels),
		}
		graph.Nodes[issue.Number] = node
	}

	// Connect nodes based on relationships
	for _, issue := range issues {
		node := graph.Nodes[issue.Number]
		
		// Connect parents
		for _, parentID := range issue.ParentIDs {
			if parent, ok := graph.Nodes[parentID]; ok {
				node.Parents = append(node.Parents, parent)
			}
		}
		
		// Connect children
		for _, childID := range issue.ChildIDs {
			if child, ok := graph.Nodes[childID]; ok {
				node.Children = append(node.Children, child)
			}
		}
	}

	// Set the root node
	if cfg.IsRootIssueSpecified() {
		graph.Root = graph.Nodes[cfg.IssueNumber]
		if graph.Root == nil {
			return nil, fmt.Errorf("root issue #%d not found", cfg.IssueNumber)
		}
	} else {
		// If no root is specified, find the top-level nodes (nodes without parents)
		var topLevelNodes []*Node
		for _, node := range graph.Nodes {
			if len(node.Parents) == 0 {
				topLevelNodes = append(topLevelNodes, node)
			}
		}
		
		// Create a virtual root node that connects to all top-level nodes
		graph.Root = &Node{
			ID:    0,
			Title: "All Issues",
			Children: topLevelNodes,
		}
	}

	return graph, nil
}

// determineNodeColor determines the color of a node based on its labels
func determineNodeColor(labels []github.Label) string {
	// Check for special labels
	for _, label := range labels {
		labelLower := strings.ToLower(label.Name)
		if strings.Contains(labelLower, "epic") {
			return "#" + label.Color // Use the label's color
		}
		if strings.Contains(labelLower, "workstream") {
			return "#" + label.Color // Use the label's color
		}
	}
	
	// Default color for regular issues
	return "#cccccc" // Gray
}
