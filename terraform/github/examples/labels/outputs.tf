output "label_count" {
  description = "The number of labels created across all repositories"
  value       = module.github_labels.label_count
}

output "repositories_configured" {
  description = "List of repositories that have been configured with standardized labels"
  value       = module.github_labels.repositories_configured
  sensitive   = true
}
