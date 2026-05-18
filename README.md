# stackdiff

Compares two Terraform state files and outputs a human-readable drift report.

## Overview

`stackdiff` helps infrastructure engineers quickly identify differences between two Terraform state files — useful for detecting configuration drift, auditing environment parity, or reviewing changes between deployments.

## Installation

### From source

```bash
git clone https://github.com/yourorg/stackdiff.git
cd stackdiff
go build -o stackdiff ./...
```

### Using `go install`

```bash
go install github.com/yourorg/stackdiff@latest
```

## Usage

```bash
stackdiff <baseline-state-file> <target-state-file>
```

### Example

```bash
stackdiff terraform.tfstate.backup terraform.tfstate
```

### Sample Output

```
=== Terraform State Drift Report ===

[ADDED] aws_s3_bucket.my_bucket
[REMOVED] aws_instance.old_server
[CHANGED] aws_security_group.web
  ~ description: "old description" -> "new description"
  ~ ingress.0.from_port: "80" -> "443"

--- Summary ---
Added:   1
Removed: 1
Changed: 1
Total drift: 3 resource(s)
```

## How It Works

1. **Parse** — Both state files are read and decoded from JSON into an internal representation.
2. **Diff** — Resources are indexed by their address and compared attribute-by-attribute.
3. **Report** — A formatted drift report is written to stdout.

## Project Structure

```
stackdiff/
├── main.go                      # Entry point; wires parser, differ, reporter
├── internal/
│   ├── parser/
│   │   ├── state.go             # Parses Terraform state JSON files
│   │   └── state_test.go
│   ├── differ/
│   │   ├── diff.go              # Compares two parsed states, produces DriftResult
│   │   └── diff_test.go
│   └── reporter/
│       ├── report.go            # Formats and writes the drift report
│       └── report_test.go
└── README.md
```

## Supported Terraform State Versions

Currently supports Terraform state format **version 4**, which is produced by Terraform 0.13 and later.

## Development

### Running tests

```bash
go test ./...
```

### Running with sample data

```bash
go run . testdata/baseline.tfstate testdata/target.tfstate
```

## Contributing

Pull requests are welcome. For significant changes, please open an issue first to discuss what you'd like to change.

## License

MIT
