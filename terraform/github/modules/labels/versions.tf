terraform {
  required_version = ">= 1.0.0"
  
  required_providers {
    github = {
      source  = "integrations/github"
      version = ">= 6.0.0, < 7.0.0"
    }
  }
}
