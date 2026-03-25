# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Alphabetical sorting of DNS profiles for consistent display
- Comprehensive documentation (ARCHITECTURE.md, CONTRIBUTING.md)
- GitHub Actions workflows for CI/CD
- Multi-platform binary releases (Linux and macOS, amd64 and arm64)

### Changed
- Complete rewrite in Go using Bubble Tea TUI framework
- Reorganized repository structure following Go best practices
- Improved status bar visibility with colored backgrounds
- Enhanced keyboard navigation (Esc key for back/quit)

## [2.0.0] - 2024-03-25

### Added
- Go implementation using Bubble Tea TUI framework
- Support for macOS (networksetup) and Linux (nmcli)
- Interactive DNS profile selection with keyboard shortcuts
- Network interface selection
- Real-time status updates with colored indicators
- Configuration loading from multiple locations
- Pre-configured popular DNS providers (Cloudflare, Google, Quad9, OpenDNS, AdGuard)
- DHCP/Automatic DNS option
- Refresh configuration without restarting (press 'r')
- Check current DNS configuration (press 'c')
- Task-based build system

### Changed
- Migrated from Python/Textual to Go/Bubble Tea
- Single binary distribution (no runtime dependencies)
- Faster startup and better performance
- Improved error handling and user feedback

### Removed
- Python implementation and dependencies
- pip/pipx installation method (now uses Go install or binary download)

## [1.0.0] - 2024-03-24

### Added
- Initial Python implementation with Textual TUI
- Basic DNS switching functionality for macOS and Linux
- YAML-based configuration
- Network interface selection
- Multiple DNS provider profiles

[Unreleased]: https://github.com/pinaka-io/dns-switch/compare/v2.0.0...HEAD
[2.0.0]: https://github.com/pinaka-io/dns-switch/compare/v1.0.0...v2.0.0
[1.0.0]: https://github.com/pinaka-io/dns-switch/releases/tag/v1.0.0
