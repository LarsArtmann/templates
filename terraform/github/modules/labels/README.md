# GitHub Standardized Labels Terraform Module

This Terraform module implements a consistent set of GitHub labels across repositories as described in [LarsArtmann/mono#38](https://github.com/LarsArtmann/mono/issues/38).

## Features

- Applies a standardized set of labels to multiple GitHub repositories
- Supports a three-tiered work structure (Epic → Workstream → SubIssue)
- Includes labels for priority, size, status, technology, area, maintenance, and release management
- Consistent colors and descriptions across all repositories
- Support for custom labels in addition to the standard set

## Usage

### Manual Usage

```hcl
module "github_labels" {
  source = "github.com/LarsArtmann/templates//terraform/github/modules/labels"

  repositories = [
    "repo1",
    "repo2",
    "repo3"
  ]

  github_token = var.github_token

  # Optional: Add custom labels specific to your organization
  custom_labels = {
    "custom/label" = {
      description = "A custom label for specific needs"
      color       = "ff0000"
    }
  }
}
```

### Automated Usage

This module can be automatically applied to repositories using the [GitHub Labels workflow](../../../../.github/workflows/github-labels.yml).

To run the workflow:
1. Go to the Actions tab in the repository
2. Select the "Apply Standardized GitHub Labels" workflow
3. Click "Run workflow"
4. Enter the repositories you want to update (or leave empty to use defaults)
5. Choose whether to include private repositories from secrets
6. Type "YES" to confirm
7. Optionally enable dry run mode to preview changes without applying them
8. Click "Run workflow"

## Requirements

| Name | Version |
|------|---------|
| terraform | >= 1.0.0 |
| github | >= 6.0.0, < 7.0.0 |

## Providers

| Name | Version |
|------|---------|
| github | >= 6.0.0, < 7.0.0 |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| repositories | List of repository names to apply the standardized labels to | `list(string)` | n/a | yes |
| github_token | GitHub personal access token with repo scope | `string` | n/a | yes |
| custom_labels | Map of custom labels to add to the standard set | `map(object({ description = string, color = string }))` | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| label_count | The number of labels created across all repositories |
| repositories_configured | List of repositories that have been configured with standardized labels |
| all_labels | Map of all labels that were created, including custom labels |

## Label Categories

The module includes the following label categories:

### Core Labels

| Label Name | Description | Color |
|------------|-------------|-------|
| `bug` | Something isn't working as expected | `#d73a4a` |
| `enhancement` | New feature or request | `#a2eeef` |
| `documentation` | Improvements or additions to documentation | `#0075ca` |
| `question` | Further information is requested | `#d876e3` |
| `security` | Security-related issues or improvements | `#b60205` |
| `duplicate` | This issue already exists | `#cfd3d7` |
| `wontfix` | This will not be worked on | `#ffffff` |
| `good-first-issue` | Good for newcomers | `#7057ff` |
| `help-wanted` | Extra attention is needed | `#008672` |
| `Epic` | A large, complex issue that contains multiple workstreams | `#9C27B0` |
| `Workstream` | A focused stream of related work within an Epic | `#00BCD4` |

### Priority Labels

| Label Name | Description | Color |
|------------|-------------|-------|
| `priority/critical` | Must be fixed ASAP | `#b60205` |
| `priority/high` | Important issue to address soon | `#d93f0b` |
| `priority/medium` | Normal priority | `#fbca04` |
| `priority/low` | Low priority, nice to have | `#0e8a16` |

### Size/Effort Labels

| Label Name | Description | Color | Line Count |
|------------|-------------|-------|------------|
| `size/XS` | Very small change | `#c2e0c6` | 0-9 lines |
| `size/S` | Small change | `#c2e0c6` | 10-29 lines |
| `size/M` | Medium change | `#fbca04` | 30-99 lines |
| `size/L` | Large change | `#eb6420` | 100-499 lines |
| `size/XL` | Very large change | `#b60205` | 500-999 lines |
| `size/XXL` | Extremely large change | `#000000` | 1000+ lines |

For the complete list of labels, see the `main.tf` file.

## Three-Tiered Work Structure

This label system supports a three-tiered work structure:

1. **Epic** (Top level) - Large initiatives that might span multiple quarters or significant project components
2. **Workstream** (Middle level) - Focused streams of related work within an Epic that can be managed as a cohesive unit
3. **SubIssue** (Bottom level) - Small, actionable tasks that can be completed in a single PR

This hierarchy is implemented as follows:
- Epics have the `Epic` label
- Workstreams have the `Workstream` label and an Epic as their parent
- SubIssues have a Workstream as their parent and appropriate size labels

## Recommended GitHub App

For automatic size labeling based on line count, we recommend installing the [pull-request-size](https://github.com/marketplace/pull-request-size) GitHub app.

## License

This module is licensed under the [European Union Public Licence (EUPL) v1.2](https://joinup.ec.europa.eu/software/page/eupl).
