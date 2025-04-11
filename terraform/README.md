# Terraform Modules and Examples

This directory contains reusable Terraform modules and examples for various infrastructure components.

## Directory Structure

```
terraform/
├── github/                  # GitHub-related Terraform modules
│   ├── modules/             # Reusable modules
│   │   └── labels/          # GitHub standardized labels module
│   └── examples/            # Usage examples
│       └── labels/          # Example for GitHub labels module
└── README.md                # This file
```

## Available Modules

### GitHub

- **[github/modules/labels](github/modules/labels)**: Terraform module for implementing a standardized GitHub label system across repositories (implements [LarsArtmann/mono#38](https://github.com/LarsArtmann/mono/issues/38))

## Examples

- **[github/examples/labels](github/examples/labels)**: Example usage of the GitHub labels module

## Usage

Each module has its own README.md with specific usage instructions. Generally, you can use these modules in your Terraform configurations like this:

```hcl
module "module_name" {
  source = "github.com/LarsArtmann/templates//terraform/path/to/module"
  
  # Module-specific variables
  var1 = "value1"
  var2 = "value2"
}
```

Note the double slash (`//`) in the source URL, which is required to specify a subdirectory within a Git repository.

## License

All modules in this directory are licensed under the [European Union Public Licence (EUPL) v1.2](https://joinup.ec.europa.eu/software/page/eupl).
