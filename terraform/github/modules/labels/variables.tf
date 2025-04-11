variable "repositories" {
  description = "List of repository names to apply the standardized labels to"
  type        = list(string)
  # Note: This variable cannot be marked as sensitive because it's used in for_each
}

variable "github_token" {
  description = "GitHub personal access token with repo scope"
  type        = string
  sensitive   = true
}

variable "custom_labels" {
  description = "Map of custom labels to add to the standard set. Each label should have a description and color."
  type = map(object({
    description = string
    color       = string
  }))
  default = {}
}
