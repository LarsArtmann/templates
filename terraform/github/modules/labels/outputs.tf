output "label_count" {
  description = "The number of labels created across all repositories"
  value       = length(github_issue_label.label)
}

output "repositories_configured" {
  description = "List of repositories that have been configured with standardized labels"
  value       = var.repositories
  sensitive   = true
}

output "all_labels" {
  description = "Map of all labels that were created, including custom labels"
  value       = local.all_labels
}
