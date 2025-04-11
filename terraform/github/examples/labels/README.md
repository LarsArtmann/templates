# GitHub Standardized Labels Example

This example demonstrates how to use the `labels` module to apply standardized labels to multiple GitHub repositories.

## Usage

1. Create a `terraform.tfvars` file with your GitHub token and optional configuration:

```hcl
github_token = "your-github-token"

# Optional: Override the GitHub owner (default: LarsArtmann)
github_owner = "your-github-username-or-org"

# Optional: Override the default public repositories list
public_repositories = [
  "repo1",
  "repo2"
]

# Optional: Add private repositories
private_repositories = [
  "private-repo1",
  "private-repo2"
]

# Optional: Add custom labels
custom_labels = {
  "team/frontend" = {
    description = "Issues related to the frontend team"
    color       = "0075ca"
  },
  "team/backend" = {
    description = "Issues related to the backend team"
    color       = "fbca04"
  }
}
```

2. Initialize Terraform:

```bash
terraform init
```

3. Apply the configuration:

```bash
terraform apply
```

## What This Example Does

This example applies the standardized GitHub label system to all repositories listed in the `public_repositories` and `private_repositories` variables. The label system includes categories for:

- Core labels (bug, enhancement, documentation, etc.)
- Priority labels (critical, high, medium, low)
- Size/effort labels (XS, S, M, L, XL, XXL)
- Status labels (blocked, in-progress, needs-review, needs-improvement)
- Technology/component labels (terraform, go, javascript, etc.)
- Area labels (frontend, backend, devops, etc.)
- Maintenance labels (refactor, dependencies, performance, etc.)
- Release labels (next, future, hotfix, breaking-change)

## Security Notes

- The `private_repositories` variable is marked as sensitive to prevent accidental exposure in logs
- You need a GitHub personal access token with the `repo` scope to apply these changes
- The token should be kept secret and not committed to version control
- This will overwrite any existing labels with the same names in the repositories

## Customization

You can customize the labels by adding your own custom labels through the `custom_labels` variable. This allows you to add organization-specific or team-specific labels while maintaining the standard set.
