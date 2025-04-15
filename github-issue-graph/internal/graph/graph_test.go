package graph

import (
	"testing"

	"github.com/LarsArtmann/mono/public/templates/github-issue-graph/internal/github"
)

func TestDetermineNodeColor(t *testing.T) {
	tests := []struct {
		name   string
		labels []github.Label
		want   string
	}{
		{
			name: "epic label",
			labels: []github.Label{
				{Name: "epic", Color: "6f42c1"},
				{Name: "enhancement", Color: "a2eeef"},
			},
			want: "#6f42c1",
		},
		{
			name: "workstream label",
			labels: []github.Label{
				{Name: "workstream", Color: "1d76db"},
				{Name: "bug", Color: "d73a4a"},
			},
			want: "#1d76db",
		},
		{
			name: "epic in label name",
			labels: []github.Label{
				{Name: "type/epic", Color: "6f42c1"},
			},
			want: "#6f42c1",
		},
		{
			name: "no special labels",
			labels: []github.Label{
				{Name: "enhancement", Color: "a2eeef"},
				{Name: "bug", Color: "d73a4a"},
			},
			want: "#cccccc",
		},
		{
			name:   "no labels",
			labels: []github.Label{},
			want:   "#cccccc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := determineNodeColor(tt.labels); got != tt.want {
				t.Errorf("determineNodeColor() = %v, want %v", got, tt.want)
			}
		})
	}
}
