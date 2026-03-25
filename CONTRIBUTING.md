# Contributing to DNS Switch

Thank you for your interest in contributing to DNS Switch! This document provides guidelines and instructions for contributing.

## Development Setup

### Prerequisites

- Go 1.21 or later
- Task (optional, but recommended): `brew install go-task`
- Git

### Getting Started

1. Fork and clone the repository:
```bash
git clone https://github.com/yourusername/dns-switch.git
cd dns-switch
```

2. Download dependencies:
```bash
task deps
# or
go mod download
```

3. Build the project:
```bash
task build
# or
go build -o dns-switch ./cmd/dns-switch
```

4. Run tests:
```bash
task test
# or
go test ./...
```

## Project Structure

```
dns-switch/
├── cmd/dns-switch/      # Main application entry point
├── internal/            # Internal packages
│   ├── config/          # Configuration management
│   ├── dns/             # DNS operations
│   └── tui/             # Terminal UI
├── docs/                # Documentation
└── config.yaml          # Example configuration
```

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for detailed architecture information.

## Code Style

### Go Code Guidelines

- Follow standard Go conventions
- Run `gofmt` before committing
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions small and focused

### Formatting

Before submitting code, ensure it's properly formatted:

```bash
task fmt
task vet
# or
go fmt ./...
go vet ./...
```

## Making Changes

### Branching Strategy

- `main` branch is the stable branch
- Create feature branches from `main`
- Use descriptive branch names: `feature/add-ipv6-support`, `fix/interface-detection`

### Commit Messages

Write clear, concise commit messages:

```
Add IPv6 DNS support

- Add IPv6 detection in dns package
- Update config to support IPv6 profiles
- Add tests for IPv6 functionality
```

Format:
- First line: brief summary (50 chars or less)
- Blank line
- Detailed description (if needed)

### Pull Requests

1. Create a feature branch
2. Make your changes
3. Add/update tests
4. Update documentation
5. Run `task check` (fmt, vet, test)
6. Submit PR with clear description

PR Title Format:
- `feat: add IPv6 support`
- `fix: interface detection on Linux`
- `docs: update installation instructions`
- `refactor: simplify config loading`

## Testing

### Running Tests

```bash
task test
# or
go test ./...
```

### Writing Tests

- Place tests in the same package as the code
- Use table-driven tests when appropriate
- Mock external dependencies (commands, file system)

Example:
```go
func TestGetNetworkInterfaces(t *testing.T) {
    // Test implementation
}
```

## Adding Features

### New DNS Provider

To add a new DNS provider profile:

1. Edit `config.yaml`:
```yaml
dns_profiles:
  newprovider:
    name: "New Provider"
    description: "Description here"
    primary: "1.2.3.4"
    secondary: "5.6.7.8"
```

### New Platform Support

To add support for a new platform (e.g., Windows):

1. Update `internal/dns/dns.go`
2. Add platform-specific functions
3. Update runtime.GOOS checks
4. Update documentation

### UI Improvements

UI changes should be made in `internal/tui/tui.go`:

- Keep the TUI responsive
- Test with different terminal sizes
- Use consistent styling (Lipgloss)
- Update keyboard shortcuts in help text

## Documentation

### Code Documentation

- Add godoc comments for exported types and functions
- Use complete sentences
- Include examples when helpful

Example:
```go
// LoadConfig loads the DNS configuration from the specified file paths.
// It tries multiple locations in order and returns the first valid config found.
func LoadConfig() (*Config, error) {
    // ...
}
```

### User Documentation

Update README.md when adding:
- New features
- Configuration options
- Installation methods
- Usage examples

## Reporting Issues

### Bug Reports

Include:
- DNS Switch version (`dns-switch --version`)
- Operating system and version
- Steps to reproduce
- Expected vs actual behavior
- Error messages or screenshots

### Feature Requests

Include:
- Use case description
- Proposed solution
- Alternative approaches considered

## Code Review Process

1. Maintainers will review PRs
2. Address feedback and suggestions
3. Once approved, PR will be merged
4. Changes will be included in next release

## Release Process

Releases follow semantic versioning (MAJOR.MINOR.PATCH):

- MAJOR: Breaking changes
- MINOR: New features (backward compatible)
- PATCH: Bug fixes

## Getting Help

- GitHub Issues: Bug reports and feature requests
- Discussions: Questions and ideas
- Documentation: Check docs/ directory

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
