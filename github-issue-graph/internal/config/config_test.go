package config

import (
	"testing"
)

func TestGetOwnerAndRepo(t *testing.T) {
	tests := []struct {
		name       string
		repoString string
		wantOwner  string
		wantRepo   string
	}{
		{
			name:       "valid repo format",
			repoString: "LarsArtmann/templates",
			wantOwner:  "LarsArtmann",
			wantRepo:   "templates",
		},
		{
			name:       "no slash",
			repoString: "LarsArtmann",
			wantOwner:  "LarsArtmann",
			wantRepo:   "",
		},
		{
			name:       "multiple slashes",
			repoString: "LarsArtmann/templates/extra",
			wantOwner:  "LarsArtmann",
			wantRepo:   "templates/extra",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Repo: tt.repoString,
			}
			gotOwner, gotRepo := cfg.GetOwnerAndRepo()
			if gotOwner != tt.wantOwner {
				t.Errorf("GetOwnerAndRepo() gotOwner = %v, want %v", gotOwner, tt.wantOwner)
			}
			if gotRepo != tt.wantRepo {
				t.Errorf("GetOwnerAndRepo() gotRepo = %v, want %v", gotRepo, tt.wantRepo)
			}
		})
	}
}

func TestIsRootIssueSpecified(t *testing.T) {
	tests := []struct {
		name        string
		issueNumber int
		want        bool
	}{
		{
			name:        "issue specified",
			issueNumber: 42,
			want:        true,
		},
		{
			name:        "no issue specified",
			issueNumber: 0,
			want:        false,
		},
		{
			name:        "negative issue number",
			issueNumber: -1,
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				IssueNumber: tt.issueNumber,
			}
			if got := cfg.IsRootIssueSpecified(); got != tt.want {
				t.Errorf("IsRootIssueSpecified() = %v, want %v", got, tt.want)
			}
		})
	}
}
