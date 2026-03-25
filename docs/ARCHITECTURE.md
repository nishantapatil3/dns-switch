# Architecture

This document describes the architecture and organization of the dns-switch project.

## Project Structure

```
dns-switch/
├── cmd/
│   └── dns-switch/          # Main application entry point
│       └── main.go          # CLI argument parsing, version info, help text
├── internal/                # Internal packages (not importable by other projects)
│   ├── config/              # Configuration management
│   │   └── config.go        # Load and parse YAML configuration
│   ├── dns/                 # DNS operations
│   │   └── dns.go           # Network interface detection, DNS get/set operations
│   └── tui/                 # Terminal User Interface
│       └── tui.go           # Bubble Tea TUI implementation
├── config.yaml              # Example configuration file
├── go.mod                   # Go module definition
├── go.sum                   # Dependency checksums
├── README.md                # Project documentation
├── LICENSE                  # MIT License
└── Taskfile.yaml            # Build and development tasks
```

## Package Organization

### cmd/dns-switch
The main entry point for the application. Handles:
- Command-line argument parsing (--version, --help)
- Initializing the TUI
- Error handling and exit codes

### internal/config
Configuration management package. Responsibilities:
- Loading YAML configuration from multiple locations (`~/.config/dns-switch/config.yaml`, `./config.yaml`)
- Parsing DNS profile definitions
- Providing configuration data structures
- Sorting DNS profiles alphabetically by name for consistent display

### internal/dns
DNS operations package. Handles:
- Detecting available network interfaces (macOS and Linux)
- Getting current DNS configuration
- Applying DNS settings to network interfaces
- Platform-specific operations (networksetup for macOS, nmcli for Linux)

### internal/tui
Terminal User Interface package using Bubble Tea. Features:
- Interactive DNS profile selection
- Network interface selection
- Real-time status updates
- Keyboard navigation
- Visual styling with Lipgloss

## Design Principles

### Standard Go Layout
The project follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout):
- `cmd/` for command-line applications
- `internal/` for private application and library code
- Root level for configuration files and documentation

### Separation of Concerns
- **UI logic** is isolated in `internal/tui`
- **DNS operations** are isolated in `internal/dns`
- **Configuration** is managed separately in `internal/config`
- **Main function** only handles CLI parsing and initialization

### Platform Support
The application supports both macOS and Linux through runtime detection:
- macOS uses `networksetup` command
- Linux uses `nmcli` (NetworkManager CLI)
- Platform-specific code is clearly marked

## Dependencies

### Direct Dependencies
- **bubbletea**: TUI framework (MIT License)
- **bubbles**: TUI components for lists (MIT License)
- **lipgloss**: Terminal styling (MIT License)
- **yaml.v3**: YAML parsing (MIT License)

### Why These Dependencies?
- **Bubble Tea**: Modern, well-maintained TUI framework with excellent documentation
- **Bubbles**: Pre-built components that handle common UI patterns
- **Lipgloss**: Powerful styling without complex configuration
- **yaml.v3**: Standard YAML library for Go

## Build System

Uses [Task](https://taskfile.dev/) for build automation:
- `task build` - Build the binary
- `task run` - Build and run with sudo
- `task install` - Install to system
- `task clean` - Remove build artifacts
- `task test` - Run tests
- `task fmt` - Format code
- `task vet` - Run static analysis

## Configuration

Configuration is loaded from (in order):
1. `~/.config/dns-switch/config.yaml` (preferred)
2. `./config.yaml` (fallback)

The YAML format supports:
- Multiple DNS profiles (sorted alphabetically by name in the UI)
- Primary and secondary DNS servers
- Descriptive names and descriptions
- DHCP/automatic configuration option (`"auto"`)
- Optional pre-selected network interface

### Profile Sorting
DNS profiles are automatically sorted alphabetically by their `name` field when displayed. This ensures consistent ordering regardless of how they're defined in the YAML file or how many times the configuration is refreshed.

## Future Improvements

Potential enhancements:
- Configuration file validation
- Unit tests for DNS operations
- Integration tests with mocked commands
- Windows support
- IPv6 DNS support
- DNS-over-HTTPS (DoH) profile support
- Export/import profile functionality
