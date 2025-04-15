package svg

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/LarsArtmann/mono/public/templates/github-issue-graph/internal/config"
	"github.com/LarsArtmann/mono/public/templates/github-issue-graph/internal/graph"
)

// NodePosition represents the position of a node in the SVG
type NodePosition struct {
	X int
	Y int
}

// SVGNode represents a node in the SVG with its position
type SVGNode struct {
	ID       int
	Title    string
	URL      string
	State    string
	Color    string
	Position NodePosition
	Width    int
	Height   int
}

// SVGConnection represents a connection between nodes in the SVG
type SVGConnection struct {
	From *SVGNode
	To   *SVGNode
}

// SVGData contains all the data needed to render the SVG
type SVGData struct {
	Width       int
	Height      int
	Nodes       []*SVGNode
	Connections []*SVGConnection
}

// Generate generates an SVG representation of the graph
func Generate(g *graph.Graph, cfg *config.Config) (string, error) {
	// Calculate positions for nodes
	svgData, err := calculateLayout(g, cfg)
	if err != nil {
		return "", err
	}

	// Render the SVG template
	return renderSVG(svgData)
}

// calculateLayout calculates the positions of nodes in the SVG
func calculateLayout(g *graph.Graph, cfg *config.Config) (*SVGData, error) {
	// Choose the layout algorithm based on the configuration
	switch cfg.Layout {
	case "tree":
		return calculateTreeLayout(g, cfg)
	case "radial":
		return calculateRadialLayout(g, cfg)
	default:
		return calculateTreeLayout(g, cfg)
	}
}

// calculateTreeLayout calculates a tree layout for the graph
func calculateTreeLayout(g *graph.Graph, cfg *config.Config) (*SVGData, error) {
	// Create SVG data
	svgData := &SVGData{
		Width:  800,
		Height: 600,
	}

	// Create a map to track node positions
	nodePositions := make(map[int]*SVGNode)

	// Start with the root node
	if g.Root == nil {
		return nil, fmt.Errorf("graph has no root node")
	}

	// Calculate positions using a recursive approach
	calculateTreePositions(g.Root, nodePositions, 0, 0, 200, 100, svgData)

	// Adjust SVG dimensions based on node positions
	adjustSVGDimensions(svgData)

	return svgData, nil
}

// calculateTreePositions recursively calculates positions for a tree layout
func calculateTreePositions(node *graph.Node, nodePositions map[int]*SVGNode, level, startX, xSpacing, ySpacing int, svgData *SVGData) int {
	// Skip if this node has already been positioned (to handle cycles)
	if _, exists := nodePositions[node.ID]; exists {
		return startX
	}

	// Create SVG node
	svgNode := &SVGNode{
		ID:    node.ID,
		Title: node.Title,
		URL:   node.URL,
		State: node.State,
		Color: node.Color,
		Width: 180,
		Height: 70,
	}

	// If this is a leaf node or we've reached the maximum level
	if len(node.Children) == 0 {
		svgNode.Position = NodePosition{X: startX, Y: level * ySpacing + 50}
		nodePositions[node.ID] = svgNode
		svgData.Nodes = append(svgData.Nodes, svgNode)
		return startX + xSpacing
	}

	// Calculate positions for children first
	childStartX := startX
	for _, child := range node.Children {
		childStartX = calculateTreePositions(child, nodePositions, level+1, childStartX, xSpacing, ySpacing, svgData)
	}

	// Position this node centered above its children
	if len(node.Children) > 0 {
		firstChildX := nodePositions[node.Children[0].ID].Position.X
		lastChildX := nodePositions[node.Children[len(node.Children)-1].ID].Position.X
		centerX := (firstChildX + lastChildX) / 2
		svgNode.Position = NodePosition{X: centerX, Y: level * ySpacing + 50}
	} else {
		svgNode.Position = NodePosition{X: startX, Y: level * ySpacing + 50}
	}

	// Add node to the map and SVG data
	nodePositions[node.ID] = svgNode
	svgData.Nodes = append(svgData.Nodes, svgNode)

	// Add connections to children
	for _, child := range node.Children {
		childNode := nodePositions[child.ID]
		svgData.Connections = append(svgData.Connections, &SVGConnection{
			From: svgNode,
			To:   childNode,
		})
	}

	// Return the next X position
	return childStartX
}

// calculateRadialLayout calculates a radial layout for the graph
func calculateRadialLayout(g *graph.Graph, cfg *config.Config) (*SVGData, error) {
	// For simplicity, we'll just use the tree layout for now
	// A proper radial layout would require more complex calculations
	return calculateTreeLayout(g, cfg)
}

// adjustSVGDimensions adjusts the SVG dimensions based on node positions
func adjustSVGDimensions(svgData *SVGData) {
	// Find the maximum X and Y coordinates
	maxX := 0
	maxY := 0
	for _, node := range svgData.Nodes {
		rightEdge := node.Position.X + node.Width/2
		bottomEdge := node.Position.Y + node.Height/2
		if rightEdge > maxX {
			maxX = rightEdge
		}
		if bottomEdge > maxY {
			maxY = bottomEdge
		}
	}

	// Add some padding
	svgData.Width = maxX + 100
	svgData.Height = maxY + 100
}

// renderSVG renders the SVG using a template
func renderSVG(data *SVGData) (string, error) {
	// Define the SVG template
	const svgTemplate = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg width="{{.Width}}" height="{{.Height}}" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
  <title>GitHub Issue Relationships</title>
  <defs>
    <!-- Arrow marker definition -->
    <marker id="arrow" markerWidth="10" markerHeight="10" refX="9" refY="3" orient="auto" markerUnits="strokeWidth">
      <path d="M0,0 L0,6 L9,3 z" fill="#666666" />
    </marker>
  </defs>

  <!-- Nodes -->
  {{range .Nodes}}
  <a xlink:href="{{.URL}}" target="_blank">
    <rect x="{{nodeLeft .}}" y="{{nodeTop .}}" width="{{.Width}}" height="{{.Height}}" rx="5" ry="5" fill="{{.Color}}" stroke="#000000" stroke-width="1" />
    <text x="{{nodeCenter .}}" y="{{nodeTextY .}}" font-family="Arial" font-size="14" text-anchor="middle" fill="{{textColor .}}">
      <tspan x="{{nodeCenter .}}" dy="0">#{{.ID}} {{truncateTitle .Title 20}}</tspan>
      {{if gt (len .Title) 20}}
      <tspan x="{{nodeCenter .}}" dy="20">{{truncateTitle (titleRemainder .Title 20) 20}}</tspan>
      {{end}}
    </text>
  </a>
  {{end}}

  <!-- Connections -->
  {{range .Connections}}
  <line x1="{{connectionFromX .}}" y1="{{connectionFromY .}}" x2="{{connectionToX .}}" y2="{{connectionToY .}}" stroke="#666666" stroke-width="2" marker-end="url(#arrow)" />
  {{end}}

  <!-- Legend -->
  <rect x="{{legendX .}}" y="{{legendY .}}" width="20" height="20" fill="#6f42c1" stroke="#000000" stroke-width="1" />
  <text x="{{legendTextX .}}" y="{{legendTextY .}}" font-family="Arial" font-size="12" text-anchor="start">Epic</text>

  <rect x="{{legendX .}}" y="{{legendY2 .}}" width="20" height="20" fill="#1d76db" stroke="#000000" stroke-width="1" />
  <text x="{{legendTextX .}}" y="{{legendTextY2 .}}" font-family="Arial" font-size="12" text-anchor="start">Workstream</text>

  <rect x="{{legendX .}}" y="{{legendY3 .}}" width="20" height="20" fill="#cccccc" stroke="#000000" stroke-width="1" />
  <text x="{{legendTextX .}}" y="{{legendTextY3 .}}" font-family="Arial" font-size="12" text-anchor="start">Task</text>
</svg>`

	// Create template functions
	funcMap := template.FuncMap{
		"nodeLeft": func(node *SVGNode) int {
			return node.Position.X - node.Width/2
		},
		"nodeTop": func(node *SVGNode) int {
			return node.Position.Y - node.Height/2
		},
		"nodeCenter": func(node *SVGNode) int {
			return node.Position.X
		},
		"nodeTextY": func(node *SVGNode) int {
			return node.Position.Y + 5
		},
		"textColor": func(node *SVGNode) string {
			// Use white text for dark backgrounds, black for light backgrounds
			if node.Color == "#cccccc" {
				return "#000000"
			}
			return "#ffffff"
		},
		"truncateTitle": func(title string, length int) string {
			if len(title) <= length {
				return title
			}
			return title[:length]
		},
		"titleRemainder": func(title string, length int) string {
			if len(title) <= length {
				return ""
			}
			return title[length:]
		},
		"connectionFromX": func(conn *SVGConnection) int {
			return conn.From.Position.X
		},
		"connectionFromY": func(conn *SVGConnection) int {
			return conn.From.Position.Y + conn.From.Height/2
		},
		"connectionToX": func(conn *SVGConnection) int {
			return conn.To.Position.X
		},
		"connectionToY": func(conn *SVGConnection) int {
			return conn.To.Position.Y - conn.To.Height/2
		},
		"legendX": func(data *SVGData) int {
			return data.Width - 100
		},
		"legendY": func(data *SVGData) int {
			return data.Height - 80
		},
		"legendY2": func(data *SVGData) int {
			return data.Height - 50
		},
		"legendY3": func(data *SVGData) int {
			return data.Height - 20
		},
		"legendTextX": func(data *SVGData) int {
			return data.Width - 70
		},
		"legendTextY": func(data *SVGData) int {
			return data.Height - 65
		},
		"legendTextY2": func(data *SVGData) int {
			return data.Height - 35
		},
		"legendTextY3": func(data *SVGData) int {
			return data.Height - 5
		},
	}

	// Parse the template
	tmpl, err := template.New("svg").Funcs(funcMap).Parse(svgTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse SVG template: %w", err)
	}

	// Execute the template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute SVG template: %w", err)
	}

	return buf.String(), nil
}
