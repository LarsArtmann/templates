# Description

A script for validating the presence of common project files in a repository. 
It can be used to ensure that all required files are present and to generate missing files if desired.

## Architecture

```mermaid
sequenceDiagram
    participant User
    participant main.go
    participant ValidateCommand
    participant Checker
    participant Reporter

    User->>main.go: Run CLI with flags
    main.go->>ValidateCommand: NewValidateCommand(config)
    main.go->>ValidateCommand: Execute()
    ValidateCommand->>Checker: CheckRepository()
    Checker-->>ValidateCommand: ValidationResults
    ValidateCommand->>Reporter: ReportResults(results)
    alt Fix requested and missing files detected
        ValidateCommand->>Checker: FixMissingFiles(results)
        Checker-->>ValidateCommand: (fixes files)
        ValidateCommand->>Checker: CheckRepository()
        Checker-->>ValidateCommand: ValidationResults (updated)
        ValidateCommand->>Reporter: ReportResults(results)
    end
    ValidateCommand-->>main.go: error (if any)
    main.go->>User: Output results / error

```