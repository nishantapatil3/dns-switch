# GitHub Actions Workflows

This directory contains automated workflows for building, testing, and releasing dns-switch.

## Workflows

### build.yml
**Trigger:** Push to `main` branch or pull requests

**Purpose:** Continuous Integration (CI)

**Actions:**
- Runs on both Ubuntu and macOS
- Tests with Go 1.21 and 1.22
- Verifies code formatting (`gofmt`)
- Runs `go vet` for static analysis
- Executes tests with race detection
- Uploads coverage reports to Codecov

**Usage:**
Automatically runs on every push and PR. Ensures code quality before merging.

### release.yml
**Trigger:** Push of version tags (e.g., `v2.0.0`)

**Purpose:** Automated release creation and binary distribution

**Actions:**
- Builds binaries for multiple platforms:
  - Linux (amd64, arm64)
  - macOS (amd64, arm64)
- Generates SHA256 checksums
- Creates GitHub release with:
  - Compressed binaries (.tar.gz)
  - Checksums file
  - Installation instructions
  - Links to changelog

**Usage:**
To create a new release:
```bash
# Create and push a version tag
git tag -a v2.1.0 -m "Release v2.1.0"
git push origin v2.1.0
```

The workflow will automatically:
1. Build all platform binaries
2. Create a GitHub release
3. Upload all artifacts

## Requirements

- Repository must have `contents: write` permission (enabled by default)
- Tags must follow semantic versioning with `v` prefix (e.g., `v2.0.0`)

## Local Testing

You can test builds locally before pushing tags:

```bash
# Test Linux build
GOOS=linux GOARCH=amd64 go build -o dns-switch-linux-amd64 ./cmd/dns-switch

# Test macOS build
GOOS=darwin GOARCH=arm64 go build -o dns-switch-darwin-arm64 ./cmd/dns-switch

# Or use Task
task build
```

## Troubleshooting

### Build workflow fails
- Check Go syntax errors: `task fmt && task vet`
- Run tests locally: `task test`
- Verify Go version compatibility

### Release workflow fails
- Ensure tag follows `v*` pattern (e.g., `v2.0.0`)
- Check that GITHUB_TOKEN has write permissions
- Verify all build commands succeed locally

## References

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [setup-go action](https://github.com/actions/setup-go)
- [action-gh-release](https://github.com/softprops/action-gh-release)
