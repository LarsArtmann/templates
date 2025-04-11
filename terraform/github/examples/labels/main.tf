terraform {
  required_providers {
    github = {
      source  = "integrations/github"
      version = "~> 6.0"
    }
  }
  required_version = ">= 1.0.0"
}

provider "github" {
  token = var.github_token
  owner = var.github_owner
}

# Use repository lists from variables
locals {
  # Combine public and private repositories if provided
  repositories = concat(var.public_repositories, var.private_repositories)
}

module "github_labels" {
  source = "../../modules/labels"

  repositories = local.repositories
  github_token = var.github_token
  
  # Optional: Add custom labels specific to your organization
  custom_labels = var.custom_labels
}
