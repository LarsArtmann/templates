package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

// Client is a wrapper around the GitHub API client
type Client struct {
	client *github.Client
	ctx    context.Context
}

// Issue represents a GitHub issue with its relationships
type Issue struct {
	Number     int
	Title      string
	URL        string
	State      string
	Labels     []Label
	ParentIDs  []int
	ChildIDs   []int
}

// Label represents a GitHub label
type Label struct {
	Name  string
	Color string
}

// NewClient creates a new GitHub client
func NewClient(token string) (*Client, error) {
	ctx := context.Background()
	var client *github.Client

	if token != "" {
		// Create authenticated client
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	} else {
		// Create unauthenticated client
		client = github.NewClient(nil)
	}

	return &Client{
		client: client,
		ctx:    ctx,
	}, nil
}

// GetIssue fetches a specific issue by number
func (c *Client) GetIssue(owner, repo string, number int) (*Issue, error) {
	issue, _, err := c.client.Issues.Get(c.ctx, owner, repo, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue #%d: %w", number, err)
	}

	return c.convertIssue(issue), nil
}

// GetAllIssues fetches all issues in a repository
func (c *Client) GetAllIssues(owner, repo string) ([]*Issue, error) {
	opts := &github.IssueListByRepoOptions{
		State:     "all",
		Sort:      "created",
		Direction: "desc",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var allIssues []*Issue
	for {
		issues, resp, err := c.client.Issues.ListByRepo(c.ctx, owner, repo, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list issues: %w", err)
		}

		for _, issue := range issues {
			// Skip pull requests
			if issue.IsPullRequest() {
				continue
			}
			allIssues = append(allIssues, c.convertIssue(issue))
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allIssues, nil
}

// FindIssueRelationships analyzes issue bodies to find parent-child relationships
// This is a simple implementation that looks for common patterns in issue descriptions
func (c *Client) FindIssueRelationships(owner, repo string, issues []*Issue) error {
	// Create a map for quick lookup
	issueMap := make(map[int]*Issue)
	for _, issue := range issues {
		issueMap[issue.Number] = issue
	}

	// For each issue, look for references to other issues
	for _, issue := range issues {
		// Get the full issue with body
		fullIssue, _, err := c.client.Issues.Get(c.ctx, owner, repo, issue.Number)
		if err != nil {
			return fmt.Errorf("failed to get issue #%d details: %w", issue.Number, err)
		}

		if fullIssue.Body == nil {
			continue
		}

		body := *fullIssue.Body
		
		// Look for parent references
		parentRefs := findParentReferences(body, owner, repo)
		for _, parentID := range parentRefs {
			if parent, ok := issueMap[parentID]; ok {
				issue.ParentIDs = append(issue.ParentIDs, parentID)
				parent.ChildIDs = append(parent.ChildIDs, issue.Number)
			}
		}

		// Look for child references
		childRefs := findChildReferences(body, owner, repo)
		for _, childID := range childRefs {
			if child, ok := issueMap[childID]; ok {
				issue.ChildIDs = append(issue.ChildIDs, childID)
				child.ParentIDs = append(child.ParentIDs, issue.Number)
			}
		}
	}

	return nil
}

// convertIssue converts a GitHub API issue to our internal Issue type
func (c *Client) convertIssue(issue *github.Issue) *Issue {
	result := &Issue{
		Number: issue.GetNumber(),
		Title:  issue.GetTitle(),
		URL:    issue.GetHTMLURL(),
		State:  issue.GetState(),
	}

	// Convert labels
	for _, label := range issue.Labels {
		result.Labels = append(result.Labels, Label{
			Name:  label.GetName(),
			Color: label.GetColor(),
		})
	}

	return result
}

// findParentReferences looks for common patterns that indicate parent issues
func findParentReferences(body, owner, repo string) []int {
	var parentIDs []int

	// Common patterns for parent references
	patterns := []string{
		"Parent issue: #",
		"Parent: #",
		"Part of #",
		"Related to #",
	}

	for _, pattern := range patterns {
		parentIDs = append(parentIDs, findIssueReferences(body, pattern)...)
	}

	return parentIDs
}

// findChildReferences looks for common patterns that indicate child issues
func findChildReferences(body, owner, repo string) []int {
	var childIDs []int

	// Common patterns for child references
	patterns := []string{
		"Child issue: #",
		"Child: #",
		"Sub-task: #",
		"Implements #",
	}

	for _, pattern := range patterns {
		childIDs = append(childIDs, findIssueReferences(body, pattern)...)
	}

	return childIDs
}

// findIssueReferences finds issue numbers after a specific pattern
func findIssueReferences(body, pattern string) []int {
	var issueIDs []int
	
	parts := strings.Split(body, pattern)
	for i := 1; i < len(parts); i++ { // Start from 1 to skip the first part (before the pattern)
		var issueID int
		_, err := fmt.Sscanf(parts[i], "%d", &issueID)
		if err == nil && issueID > 0 {
			issueIDs = append(issueIDs, issueID)
		}
	}
	
	return issueIDs
}
