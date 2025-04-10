# PakSafe

PakSafe is a command-line tool for scanning package dependencies for vulnerabilities and other security issues. It currently supports NPM packages.

## Installation

```bash
go get github.com/your-username/paksafe
```

## Building from Source

Clone the repository and build using Go:

```bash
git clone https://github.com/your-username/paksafe.git
cd paksafe
go build -o paksafe
```

## Usage

### Basic Commands

```bash
# Scan dependencies in the default package.json
paksafe npm

# Scan dependencies in a specific package.json
paksafe npm path/to/package.json

# Show help information
paksafe --help
paksafe npm --help
```

### Options

- `--only-list`: Only list dependencies without scanning
- `--scan`: Explicitly scan dependencies for vulnerabilities
- `--format`: Output format (json, stdout, terminal, console)
- `--output-path`: Specify output file path for JSON format

### Examples

List dependencies only:
```bash
paksafe npm --only-list
```

Scan and output as a formatted table:
```bash
paksafe npm --format stdout
```

Scan and save results to a JSON file:
```bash
paksafe npm --format json --output-path ./results/scan-results.json
```

## Output Formats

### Table Format (stdout/terminal/console)

The table output displays the package name, version, and status:

```
PACKAGE                       VERSION   STATUS      
-----------------------------+---------+-------------
react-scripts                 1.0.0     available   
react-webcam                  1.0.0     available   
socket.io-client              1.0.0     vulnerable  
```

### JSON Format

The JSON output is structured as follows:

```json
{
  "packages": [
    {
      "package": "react-scripts",
      "version": "1.0.0",
      "status": "available"
    },
    {
      "package": "react-webcam",
      "version": "1.0.0",
      "status": "available"
    },
    {
      "package": "socket.io-client",
      "version": "1.0.0",
      "status": "vulnerable"
    }
  ]
}
```

## Technical Details

### Dependencies

- `github.com/spf13/cobra`: Command-line interface
- `github.com/olekukonko/tablewriter`: Table formatting
- `github.com/fatih/color`: Terminal color output

### Package Structure

- `cmd/`: Command-line interface implementation
- `internal/dependencies/`: Package dependency handling
- `internal/registry/`: Package registry checking
- `format/`: Output formatting
- `utils/`: Utility functions

### Registry Checking

PakSafe checks each dependency against the NPM registry to determine if it exists and is safe to use. It performs HTTP requests to the registry API and analyzes the responses.

## Future Plans

- Support for additional package managers (Yarn, pip, Maven, etc.)
- Vulnerability scanning against known CVE databases
- License compliance checking
- Integration with CI/CD pipelines

## License

[MIT License](LICENSE)