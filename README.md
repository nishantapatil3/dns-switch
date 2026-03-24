# DNS Switch TUI

A user-friendly Terminal User Interface (TUI) application for quickly switching between different DNS configurations.

## Features

- 🎨 Beautiful, intuitive TUI built with Textual
- 🚀 Quick DNS profile switching
- ⚙️ YAML-based configuration
- 🔄 Support for multiple network interfaces
- 🍎 macOS support (using `networksetup`)
- 🐧 Linux support (using `nmcli`)
- 📋 Pre-configured popular DNS providers (Cloudflare, Google, Quad9, OpenDNS)

## Installation

### Recommended: pipx (Isolated Installation)

[pipx](https://pipx.pypa.io/) installs the application in an isolated environment:

```bash
# Install pipx if you haven't already
brew install pipx  # macOS
# or: python3 -m pip install --user pipx

# Install dns-switch
pipx install dns-switch

# Run
sudo dns-switch
```

### Alternative: pip

```bash
pip install dns-switch
sudo dns-switch
```

### From Source

```bash
git clone https://github.com/pinaka-io/dns-switch.git
cd dns-switch
pip install -e .
sudo dns-switch
```

## Configuration

Edit `config.yaml` to add or modify DNS profiles:

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

- `name`: Display name for the profile
- `description`: Brief description of the DNS provider
- `primary`: Primary DNS server IP
- `secondary`: Secondary DNS server IP (optional)
- Use `"auto"` for both primary and secondary to use DHCP

## Usage

Run the application with sudo:
```bash
sudo dns-switch
```

### Keyboard Shortcuts

- **Arrow Keys / Mouse**: Navigate through DNS profiles
- **Enter / Apply Button**: Apply selected DNS profile
- **i**: Change network interface
- **r**: Refresh configuration
- **q**: Quit application

### Steps to Switch DNS

1. Launch the application
2. Select your network interface (if not already configured)
3. Use arrow keys or mouse to select a DNS profile
4. Press Enter or click "Apply DNS" button
5. Wait for confirmation message

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

Edit `config.yaml` and add your profile:

```yaml
dns_profiles:
  my_custom:
    name: "My ISP DNS"
    description: "Optimized for my network"
    primary: "10.0.0.1"
    secondary: "10.0.0.2"
```

Then restart the application or press `r` to refresh.

## License

MIT License - feel free to use and modify as needed!

## Contributing

Contributions welcome! Feel free to submit issues or pull requests.
