# DNS Switch TUI

A fast, user-friendly Terminal User Interface (TUI) for quickly switching between different DNS configurations. Written in Go with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Features

- 🎨 Beautiful, intuitive TUI built with Bubble Tea
- 🚀 Quick DNS profile switching with real-time status updates
- ⚙️ YAML-based configuration with automatic alphabetical sorting
- 🔄 Support for multiple network interfaces
- 🍎 macOS support (using `networksetup`)
- 🐧 Linux support (using `nmcli`)
- 📋 Pre-configured popular DNS providers (Cloudflare, Google, Quad9, OpenDNS, AdGuard)
- ⚡ Fast and lightweight (single binary, no runtime dependencies)
- 🎯 Consistent display - profiles always shown in alphabetical order
- 🔍 Check current DNS configuration with one keystroke

## Installation

### Using Go Install

```bash
go install github.com/pinaka-io/dns-switch@latest

# Run
sudo dns-switch
```

### Build from Source

```bash
git clone https://github.com/pinaka-io/dns-switch.git
cd dns-switch

# Using Task (recommended)
task build
task install

# Or build directly with Go
go build -o dns-switch .
sudo mv dns-switch /usr/local/bin/

# Run directly
task run
```

### Download Binary

Download the latest binary from the [releases page](https://github.com/pinaka-io/dns-switch/releases):

**Linux (amd64):**
```bash
curl -LO https://github.com/pinaka-io/dns-switch/releases/latest/download/dns-switch-linux-amd64.tar.gz
tar xzf dns-switch-linux-amd64.tar.gz
sudo mv dns-switch-linux-amd64 /usr/local/bin/dns-switch
sudo chmod +x /usr/local/bin/dns-switch
```

**macOS (Apple Silicon):**
```bash
curl -LO https://github.com/pinaka-io/dns-switch/releases/latest/download/dns-switch-darwin-arm64.tar.gz
tar xzf dns-switch-darwin-arm64.tar.gz
sudo mv dns-switch-darwin-arm64 /usr/local/bin/dns-switch
sudo chmod +x /usr/local/bin/dns-switch
```

**macOS (Intel):**
```bash
curl -LO https://github.com/pinaka-io/dns-switch/releases/latest/download/dns-switch-darwin-amd64.tar.gz
tar xzf dns-switch-darwin-amd64.tar.gz
sudo mv dns-switch-darwin-amd64 /usr/local/bin/dns-switch
sudo chmod +x /usr/local/bin/dns-switch
```

See the [releases page](https://github.com/pinaka-io/dns-switch/releases) for other platforms (Linux arm64) and checksums.

### Development

```bash
git clone https://github.com/pinaka-io/dns-switch.git
cd dns-switch
go mod download
go run .
```

## Configuration

Create or edit `~/.config/dns-switch/config.yaml` (or use `config.yaml` in the current directory):

```yaml
dns_profiles:
  cloudflare:
    name: "Cloudflare (1.1.1.1)"
    description: "Fast and privacy-focused DNS"
    primary: "1.1.1.1"
    secondary: "1.0.0.1"

  # Add your custom profiles here
  custom:
    name: "My Custom DNS"
    description: "My preferred DNS servers"
    primary: "192.168.1.1"
    secondary: "192.168.1.2"
```

### Configuration Options

- `name`: Display name for the profile (profiles are sorted alphabetically by name)
- `description`: Brief description of the DNS provider
- `primary`: Primary DNS server IP
- `secondary`: Secondary DNS server IP (optional)
- Use `"auto"` for both primary and secondary to use DHCP
- `network_interface`: Pre-select a network interface (optional)

**Note:** Profiles are automatically sorted alphabetically by name for consistent display.

## Usage

Run the application with sudo:
```bash
sudo dns-switch
```

### Keyboard Shortcuts

- **Arrow Keys / j/k**: Navigate through options
- **Enter**: Apply selected DNS profile / Select interface
- **c**: Check current DNS configuration
- **i**: Change network interface
- **r**: Refresh configuration
- **q / Ctrl+C**: Quit application
- **Esc**: Back to interface selection (from profile list) / Quit (from interface selection)

### Steps to Switch DNS

1. Launch the application: `sudo dns-switch`
2. Select your network interface (if not already configured)
3. Use arrow keys (`↑/↓` or `j/k`) to select a DNS profile
4. Press `Enter` to apply the DNS profile
5. Check the status bar for confirmation message (green = success, red = error)

## Permissions

The application requires administrator privileges to modify DNS settings:

### macOS
Uses `networksetup` command (requires sudo):
```bash
sudo dns-switch
```

### Linux
Uses `nmcli` command (may require sudo):
```bash
sudo dns-switch
```

**Optional:** Configure sudo to avoid password prompt:
```bash
sudo visudo
# Add this line (replace USERNAME):
USERNAME ALL=(ALL) NOPASSWD: /usr/bin/nmcli
```

## Troubleshooting

### "No network interfaces found"
- **macOS**: Ensure you have permission to run `networksetup`
- **Linux**: Install NetworkManager and nmcli: `sudo apt install network-manager`

### "Permission denied"
- Run the application with `sudo`

### Changes don't take effect
- Try disabling and re-enabling your network interface
- Check if you selected the correct interface

## Adding Custom DNS Profiles

Edit your configuration file (`~/.config/dns-switch/config.yaml` or `config.yaml`) and add your profile:

```yaml
dns_profiles:
  my_custom:
    name: "My ISP DNS"
    description: "Optimized for my network"
    primary: "10.0.0.1"
    secondary: "10.0.0.2"
```

Then restart the application or press `r` to refresh. Your new profile will appear in alphabetical order by name.

## Project Structure

```
dns-switch/
├── cmd/dns-switch/      # Main application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── dns/             # DNS operations
│   └── tui/             # Terminal UI (Bubble Tea)
├── config.yaml          # Example configuration
└── Taskfile.yaml        # Build tasks
```

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for detailed architecture documentation.

## Development

### Building

```bash
# Download dependencies
task deps

# Build the binary
task build

# Run directly
task run
```

### Available Tasks

```bash
task              # Show all available tasks
task build        # Build the binary
task install      # Install to /usr/local/bin
task run          # Build and run
task clean        # Clean build artifacts
task test         # Run tests
task fmt          # Format code
task vet          # Run go vet
task deps         # Update dependencies
task check        # Run all checks (fmt, vet, test)
```

## Contributing

Contributions welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Quick Start for Contributors

1. Fork and clone the repository
2. Install Task: `brew install go-task`
3. Build: `task build`
4. Test: `task test`
5. Format: `task fmt`
6. Submit PR

## Releases

Releases are automatically built and published via GitHub Actions when a new tag is pushed:

```bash
git tag -a v2.1.0 -m "Release v2.1.0"
git push origin v2.1.0
```

This will trigger the release workflow which builds binaries for:
- Linux (amd64, arm64)
- macOS (amd64, arm64)

See [CHANGELOG.md](CHANGELOG.md) for release history.

## License

MIT License - feel free to use and modify as needed!
