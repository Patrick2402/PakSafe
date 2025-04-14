# PakSafe - Technical Documentation

PakSafe is a command-line tool written in Go for scanning package dependencies for security vulnerabilities, with a focus on detecting dependency confusion attacks. It currently supports NPM packages.

## Architecture Overview

PakSafe follows a modular architecture with clear separation of concerns:

```
PakSafe/
├── cmd/             # Command-line interface implementation
├── format/          # Output formatting utilities
├── internal/        # Core functionality
│   ├── dependencies/ # Package dependency handling
│   ├── registry/    # Package registry checking
│   └── types/       # Data structures and types
├── output/          # Default location for scan results
├── tests/           # Test files and fixtures
└── utils/           # Utility functions
```

## Core Components

### Command Line Interface

PakSafe uses the [Cobra](https://github.com/spf13/cobra) library to implement its CLI. The main command structure is:

```
paksafe [flags]
paksafe npm [flags] [package.json path]
```

The application entry point is in `cmd/root.go`, which defines the root command and its subcommands.

### Dependency Analysis

Dependency analysis is handled by the `dependencies` package:

1. `GetDependencies()` - Parses a package.json file and extracts dependency information
2. `BuildJson()` - Formats dependency data and registry status into a standardized JSON structure
3. `SaveJsonToFile()` - Writes the results to a JSON file

### Registry Checking

Registry checking is implemented in `registry/registry.go`:

1. `CheckForDependencyConfusion()` - Main function that determines if dependencies are vulnerable to confusion attacks
2. `checkPackageInRegistry()` - Verifies if a package exists in a specific registry
3. `extractLatestVersion()` - Gets the latest version of a package from registry response
4. `compareVersions()` - Compares semantic versions to detect version bombing attacks

### Data Structures

Key data structures defined in `internal/types/types.go`:

```go
type Dependency struct {
    Name    string
    Version string
}

type DependencyStatus struct {
    IsVulnerable   bool
    Status         string
    Reason         string
    PrivateVersion string
    PublicVersion  string
}
```

### Output Formatting

Output formatting is handled in `format/table.go`, which uses the [tablewriter](https://github.com/olekukonko/tablewriter) library to display scan results in a formatted table.

## Dependency Confusion Detection

PakSafe implements a comprehensive approach to detect dependency confusion vulnerabilities:

### Vulnerability Types Detected

1. **Name Squatting** - When a package exists only in a private registry but not in the public one
2. **Version Bombing** - When a public version of a package has a higher version number than the private one
3. **Scope Issues** - Detecting issues with scoped packages (e.g., @organization/package)

### Detection Logic

The detection algorithm considers several factors:
- Whether the package is scoped with an organization prefix (@)
- Whether it belongs to a private scope
- Its existence in private and public registries
- Version differences between registries

Truth table for vulnerability detection:

| Scoped | Own Scope | In Private | In Public | Status | Reason |
|--------|-----------|------------|-----------|--------|--------|
| Yes | Yes | Yes | Yes | Secure | Own-scoped package exists in both registries (likely intentional) |
| Yes | Yes | Yes | No | Secure | Own-scoped package exists only in private registry (secure scenario) |
| Yes | Yes | No | Yes | Vulnerable | Own-scoped package exists only in public registry (high risk dependency confusion) |
| Yes | Yes | No | No | Not Found | Own-scoped package not found in any registry |
| Yes | No | Yes | Yes | Suspicious | External-scoped package exists in both registries |
| Yes | No | Yes | No | Suspicious | External-scoped package exists only in private registry (unusual scenario) |
| Yes | No | No | Yes | Secure | External-scoped public package (normal scenario) |
| Yes | No | No | No | Not Possible | External-scoped package not found anywhere |
| No | No | Yes | Yes | Suspicious | Unscoped package exists in both registries |
| No | No | Yes | No | Vulnerable | Unscoped package exists only in private registry (high risk for name squatting) |
| No | No | No | Yes | Secure | Public package exists only in public registry (normal scenario) |
| No | No | No | No | No Package | Package not found in any registry |

## Usage Examples

### Scanning with Default Settings

```bash
paksafe npm
```

This scans dependencies in the current directory's package.json and outputs the results as a table.

### Scanning a Specific Package File

```bash
paksafe npm ./tests/package.json
```

### Scan with Private Registry Configuration

```bash
paksafe npm --private-registry http://localhost:4873 --private-scope @your-scope
```

This helps detect dependency confusion vulnerabilities by checking both private and public registries.

### Output as JSON

```bash
paksafe npm --format json --output-file ./output/results.json
```

## Output Formats

### Table Output (Default)

The table output displays:
- Package name
- Version
- Status (color-coded):
  - Green: secure
  - Red: vulnerable
  - Yellow: suspicious
  - Cyan: not found
- Reason for the status
- Private latest version
- Public latest version

### JSON Output

JSON output is structured as follows:

```json
{
  "packages": [
    {
      "package": "react-scripts",
      "version": "5.0.1",
      "status": "secure",
      "reason": "Public package exists only in public registry (normal scenario)",
      "privateVersion": "",
      "publicVersion": "5.0.1"
    },
    {
      "package": "@your-scope/example",
      "version": "1.0.0",
      "status": "vulnerable",
      "reason": "Own-scoped package exists only in public registry (high risk dependency confusion)",
      "privateVersion": "",
      "publicVersion": "1.0.2"
    }
  ]
}
```

## Technical Implementation Details

### HTTP Requests

Package existence is checked by making HTTP GET requests to registry endpoints:

```go
url := fmt.Sprintf("%s/%s", strings.TrimRight(registryURL, "/"), packageName)
resp, err := http.Get(url)
```

### Version Comparison

PakSafe implements a simplified semantic version comparison to detect version bombing attacks:

```go
func compareVersions(version1, version2 string) int {
    if version1 == "" {
        return -1
    }
    if version2 == "" {
        return 1
    }
    
    if version1 > version2 {
        return 1
    } else if version1 < version2 {
        return -1
    }
    return 0
}
```

## Dependencies

- `github.com/spf13/cobra`: Command-line interface framework
- `github.com/olekukonko/tablewriter`: Table formatting for terminal output
- `github.com/fatih/color`: Terminal color output for better readability

## Building and Testing

### Building from Source

```bash
git clone https://github.com/your-username/paksafe.git
cd paksafe
go build -o paksafe
```

### Running Tests

```bash
go test ./... -v
```

## Future Development

Potential extensions for the project:

1. Support for additional package managers (Yarn, pip, Maven, etc.)
2. Integration with CVE databases for known vulnerability checking
3. License compliance scanning
4. CI/CD pipeline integration
5. Improved semantic version comparison
6. Automated vulnerability remediation recommendations