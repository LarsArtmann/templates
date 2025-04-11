variable "github_token" {
  description = "GitHub personal access token with repo scope"
  type        = string
  sensitive   = true
}

variable "github_owner" {
  description = "GitHub owner (username or organization)"
  type        = string
  default     = "LarsArtmann"
}

variable "public_repositories" {
  description = "List of public repository names to apply the standardized labels to"
  type        = list(string)
  default     = [
    "mono",
    "LarsArtmann",
    "Setup-Mac",
    "templates"
  ]
}

variable "private_repositories" {
  description = "List of private repository names to apply the standardized labels to"
  type        = list(string)
  default     = []
  # Note: This variable cannot be marked as sensitive because it's used in for_each
}

variable "custom_labels" {
  description = "Map of custom labels to add to the standard set"
  type = map(object({
    description = string
    color       = string
  }))
  default = {}
}
