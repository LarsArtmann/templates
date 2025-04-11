# Define all labels with their properties
locals {
  # Core Labels
  core_labels = {
    "bug" = {
      description = "Something isn't working as expected"
      color       = "d73a4a"
    },
    "enhancement" = {
      description = "New feature or request"
      color       = "a2eeef"
    },
    "documentation" = {
      description = "Improvements or additions to documentation"
      color       = "0075ca"
    },
    "question" = {
      description = "Further information is requested"
      color       = "d876e3"
    },
    "security" = {
      description = "Security-related issues or improvements"
      color       = "b60205"
    },
    "duplicate" = {
      description = "This issue already exists"
      color       = "cfd3d7"
    },
    "wontfix" = {
      description = "This will not be worked on"
      color       = "ffffff"
    },
    "good-first-issue" = {
      description = "Good for newcomers"
      color       = "7057ff"
    },
    "help-wanted" = {
      description = "Extra attention is needed"
      color       = "008672"
    },
    "Epic" = {
      description = "A large, complex issue that contains multiple workstreams"
      color       = "9C27B0"
    },
    "Workstream" = {
      description = "A focused stream of related work within an Epic"
      color       = "00BCD4"
    }
  }

  # Priority Labels
  priority_labels = {
    "priority/critical" = {
      description = "Must be fixed ASAP"
      color       = "b60205"
    },
    "priority/high" = {
      description = "Important issue to address soon"
      color       = "d93f0b"
    },
    "priority/medium" = {
      description = "Normal priority"
      color       = "fbca04"
    },
    "priority/low" = {
      description = "Low priority, nice to have"
      color       = "0e8a16"
    }
  }

  # Size/Effort Labels
  size_labels = {
    "size/XS" = {
      description = "Very small change (0-9 lines)"
      color       = "c2e0c6"
    },
    "size/S" = {
      description = "Small change (10-29 lines)"
      color       = "c2e0c6"
    },
    "size/M" = {
      description = "Medium change (30-99 lines)"
      color       = "fbca04"
    },
    "size/L" = {
      description = "Large change (100-499 lines)"
      color       = "eb6420"
    },
    "size/XL" = {
      description = "Very large change (500-999 lines)"
      color       = "b60205"
    },
    "size/XXL" = {
      description = "Extremely large change (1000+ lines)"
      color       = "000000"
    }
  }

  # Status Labels
  status_labels = {
    "status/blocked" = {
      description = "Blocked by another issue or external factor"
      color       = "b60205"
    },
    "status/in-progress" = {
      description = "Currently being worked on"
      color       = "0e8a16"
    },
    "status/needs-review" = {
      description = "Ready for review"
      color       = "fbca04"
    },
    "status/needs-improvement" = {
      description = "Needs refinement before proceeding"
      color       = "c5def5"
    }
  }

  # Technology/Component Labels
  tech_labels = {
    "tech/terraform" = {
      description = "Related to Terraform"
      color       = "5319e7"
    },
    "tech/go" = {
      description = "Related to Go code"
      color       = "00ADD8"
    },
    "tech/javascript" = {
      description = "Related to JavaScript code"
      color       = "f1e05a"
    },
    "tech/typescript" = {
      description = "Related to TypeScript code"
      color       = "3178c6"
    },
    "tech/kotlin" = {
      description = "Related to Kotlin code"
      color       = "A97BFF"
    },
    "tech/python" = {
      description = "Related to Python code"
      color       = "3572A5"
    },
    "tech/rust" = {
      description = "Related to Rust code"
      color       = "DEA584"
    }
  }

  # Area Labels
  area_labels = {
    "area/frontend" = {
      description = "Frontend-related changes"
      color       = "1d76db"
    },
    "area/backend" = {
      description = "Backend-related changes"
      color       = "0052cc"
    },
    "area/devops" = {
      description = "DevOps-related changes"
      color       = "0e8a16"
    },
    "area/testing" = {
      description = "Testing-related changes"
      color       = "d4c5f9"
    },
    "area/automation" = {
      description = "Automation-related changes"
      color       = "fbca04"
    },
    "area/security" = {
      description = "Security-related changes"
      color       = "b60205"
    },
    "area/infrastructure" = {
      description = "Infrastructure-related changes"
      color       = "6b8e23"
    },
    "area/ci-cd" = {
      description = "CI/CD pipeline changes"
      color       = "4b0082"
    },
    "area/ui-ux" = {
      description = "User interface and experience changes"
      color       = "ff69b4"
    },
    "area/migration" = {
      description = "Changes related to migrating between systems or versions"
      color       = "5319e7"
    }
  }

  # Maintenance Labels
  maintenance_labels = {
    "maintenance/refactor" = {
      description = "Code refactoring without functional changes"
      color       = "c5def5"
    },
    "maintenance/dependencies" = {
      description = "Dependency updates"
      color       = "0366d6"
    },
    "maintenance/performance" = {
      description = "Performance improvements"
      color       = "0e8a16"
    },
    "maintenance/cleanup" = {
      description = "Code cleanup or removal"
      color       = "c5def5"
    },
    "maintenance/architecture" = {
      description = "Architecture improvements"
      color       = "1d76db"
    }
  }

  # Release Labels
  release_labels = {
    "release/next" = {
      description = "Targeted for the next release"
      color       = "0366d6"
    },
    "release/future" = {
      description = "Planned for a future release"
      color       = "5319e7"
    },
    "release/hotfix" = {
      description = "Needs immediate hotfix release"
      color       = "b60205"
    },
    "release/breaking-change" = {
      description = "Introduces breaking changes"
      color       = "d93f0b"
    }
  }

  # Combine all label categories
  all_labels = merge(
    local.core_labels,
    local.priority_labels,
    local.size_labels,
    local.status_labels,
    local.tech_labels,
    local.area_labels,
    local.maintenance_labels,
    local.release_labels,
    var.custom_labels
  )
}

# Create labels for each repository
resource "github_issue_label" "label" {
  for_each = {
    for pair in setproduct(var.repositories, keys(local.all_labels)) : "${pair[0]}_${pair[1]}" => {
      repository  = pair[0]
      name        = pair[1]
      description = local.all_labels[pair[1]].description
      color       = local.all_labels[pair[1]].color
    }
  }

  repository  = each.value.repository
  name        = each.value.name
  description = each.value.description
  color       = each.value.color
}
