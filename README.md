# LarsArtmann/templates

[![License: EUPL-1.2](https://img.shields.io/badge/License-EUPL--1.2-blue.svg)](https://joinup.ec.europa.eu/software/page/eupl)
[![Maintained by Lars Artmann](https://img.shields.io/badge/Maintained%20by-Lars%20Artmann-brightgreen)](https://lars.software)

This repository serves as a centralized collection of templates, best practices, and automation tools for software development projects.

## Purpose

This repository is part of the larger initiative described in [LarsArtmann/mono#42](https://github.com/LarsArtmann/mono/issues/42) to create a single source of truth for development standards and practices accumulated over 12+ years of software architecture experience.

The templates repository will include:
- Project templates for various languages and architectures
- GitHub automation tools and workflows
- Infrastructure as Code templates
- Code quality and security configurations
- Documentation templates and standards

## Current Status

This repository is in its initial setup phase. The structure and content will be expanded incrementally through future issues.

## Directory Structure

```
templates/
├── README.md (this file)
├── LICENSE (European Union Public Licence v1.2)
├── .gitignore
├── .github/
│   └── workflows/ (GitHub Actions workflows)
├── terraform/
│   ├── github/ (GitHub-related Terraform modules)
│   └── README.md (Terraform modules documentation)
└── github/
    └── README.md (GitHub-specific templates and automation)
```

## Available Tools

### Terraform Modules

- **[GitHub Labels](terraform/github/modules/labels)**: Terraform module for implementing a standardized GitHub label system across repositories

### GitHub Actions Workflows

- **[GitHub Labels](/.github/workflows/github-labels.yml)**: Workflow for automatically applying standardized labels to repositories

## Related Issues

- [LarsArtmann/mono#42](https://github.com/LarsArtmann/mono/issues/42) - LarsArtmann/templates Project: Centralized Repository for Best Practices and Automation
- [LarsArtmann/mono#38](https://github.com/LarsArtmann/mono/issues/38) - Implement standardized GitHub labels across all repositories with Terraform
- [LarsArtmann/mono#50](https://github.com/LarsArtmann/mono/issues/50) - Update templates repository to use European Union Public Licence (EUPL) v1.2
- [LarsArtmann/templates#4](https://github.com/LarsArtmann/templates/issues/4) - GitHub Actions Workflow for Automated Label Management

## Next Steps

Once this basic repository is created, additional issues will handle:
- Adding it as a submodule to the mono repository
- Implementing specific components described in issue #42
- Starting with the GitHub label standardization (issue #38)

## License

This repository is licensed under the [European Union Public Licence (EUPL) v1.2](https://joinup.ec.europa.eu/software/page/eupl). See the [LICENSE](LICENSE) file for details.
