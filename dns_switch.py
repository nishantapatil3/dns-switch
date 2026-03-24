#!/usr/bin/env python3
"""
DNS Switch TUI - A user-friendly DNS configuration switcher
"""

import argparse
import subprocess
import sys
import yaml
from pathlib import Path
from typing import Dict, List, Optional

from textual.app import App, ComposeResult
from textual.containers import Container, Vertical, Horizontal
from textual.widgets import Header, Footer, Static, Button, ListView, ListItem, Label
from textual.binding import Binding
from textual.screen import Screen
from textual import on

__version__ = "1.0.0"


class DNSProfile:
    """Represents a DNS profile configuration"""

    def __init__(self, key: str, data: Dict):
        self.key = key
        self.name = data.get("name", key)
        self.description = data.get("description", "")
        self.primary = data.get("primary", "")
        self.secondary = data.get("secondary", "")

    def __str__(self):
        return f"{self.name} - {self.description}"


class InterfaceSelectionScreen(Screen):
    """Screen for selecting network interface"""

    BINDINGS = [
        Binding("q", "quit", "Quit"),
        Binding("escape", "app.pop_screen", "Back"),
    ]

    def __init__(self, interfaces: List[str], on_select):
        super().__init__()
        self.interfaces = interfaces
        self.on_select = on_select

    def compose(self) -> ComposeResult:
        yield Container(
            Static("Select Network Interface", classes="title"),
            ListView(
                *[ListItem(Label(iface)) for iface in self.interfaces],
                id="interface-list"
            ),
            id="interface-container"
        )
        yield Footer()

    @on(ListView.Selected, "#interface-list")
    def on_interface_selected(self, event: ListView.Selected) -> None:
        interface = self.interfaces[event.list_view.index]
        self.on_select(interface)
        self.app.pop_screen()


class DNSSwitchApp(App):
    """Main DNS Switch TUI Application"""

    CSS = """
    Screen {
        background: $surface;
        align: center middle;
    }

    #main-container {
        width: 90;
        height: auto;
        max-height: 25;
        background: $panel;
        border: solid $primary;
    }

    .title {
        width: 100%;
        height: 1;
        content-align: center middle;
        text-style: bold;
        color: $accent;
        background: $surface;
        padding: 0 1;
    }

    #profile-list {
        width: 100%;
        height: auto;
        max-height: 12;
        border: none;
    }

    ListItem {
        height: 1;
        padding: 0 1;
    }

    ListItem:hover {
        background: $primary-darken-2;
    }

    ListItem > Horizontal {
        width: 100%;
        height: 1;
    }

    .profile-name {
        text-style: bold;
        color: $accent;
        width: 25;
    }

    .profile-dns {
        color: $text-muted;
        width: 1fr;
    }

    #button-container {
        width: 100%;
        height: auto;
        align: center middle;
        padding: 0;
        background: $surface;
    }

    #status {
        width: 100%;
        height: 1;
        content-align: center middle;
        background: $surface;
        color: $text-muted;
        padding: 0 1;
    }

    .error {
        color: $error;
    }

    .success {
        color: $success;
    }

    Button {
        margin: 0 1;
        min-width: 14;
        height: 3;
    }

    #interface-container {
        width: 60;
        height: auto;
        max-height: 20;
        background: $panel;
        border: solid $primary;
        align: center middle;
    }

    #interface-list {
        width: 100%;
        height: auto;
        max-height: 15;
        border: none;
    }

    #interface-list ListItem {
        height: 1;
        padding: 0 1;
    }
    """

    BINDINGS = [
        Binding("q", "quit", "Quit"),
        Binding("r", "refresh", "Refresh"),
        Binding("i", "change_interface", "Change Interface"),
        Binding("c", "check_dns", "Check DNS"),
    ]

    def __init__(self):
        super().__init__()
        # Check multiple config locations
        config_locations = [
            Path.home() / ".config" / "dns-switch" / "config.yaml",
            Path("config.yaml"),
            Path(__file__).parent / "config.yaml",
        ]

        self.config_path = None
        for config_path in config_locations:
            if config_path.exists():
                self.config_path = config_path
                break

        if not self.config_path:
            self.config_path = Path("config.yaml")

        self.profiles: Dict[str, DNSProfile] = {}
        self.current_interface: Optional[str] = None
        self.load_config()

    def compose(self) -> ComposeResult:
        yield Container(
            Static(f"DNS Switch [{self.current_interface or 'No Interface'}]",
                   classes="title", id="title"),
            ListView(id="profile-list"),
            Horizontal(
                Button("Apply", variant="primary", id="apply-btn"),
                Button("Check", variant="default", id="check-btn"),
                Button("Interface", variant="default", id="interface-btn"),
                Button("Refresh", variant="default", id="refresh-btn"),
                id="button-container"
            ),
            Static("Select profile and press Apply", id="status"),
            id="main-container"
        )
        yield Footer()

    def on_mount(self) -> None:
        """Initialize the app when mounted"""
        self.refresh_profile_list()
        if not self.current_interface:
            self.select_interface()

    def load_config(self) -> None:
        """Load DNS profiles from config.yaml"""
        try:
            with open(self.config_path, 'r') as f:
                config = yaml.safe_load(f)

            profiles_data = config.get('dns_profiles', {})
            self.profiles = {
                key: DNSProfile(key, data)
                for key, data in profiles_data.items()
            }

            # Get network interface from config
            interface = config.get('network_interface', '')
            if interface:
                self.current_interface = interface

        except FileNotFoundError:
            self.notify("Config file not found!", severity="error")
            sys.exit(1)
        except yaml.YAMLError as e:
            self.notify(f"Error parsing config: {e}", severity="error")
            sys.exit(1)

    def refresh_profile_list(self) -> None:
        """Refresh the DNS profile list"""
        profile_list = self.query_one("#profile-list", ListView)
        profile_list.clear()

        for profile in self.profiles.values():
            dns_info = f"{profile.primary}"
            if profile.secondary and profile.secondary != profile.primary:
                dns_info += f", {profile.secondary}"

            item = ListItem(
                Horizontal(
                    Label(profile.name, classes="profile-name"),
                    Label(dns_info, classes="profile-dns"),
                )
            )
            profile_list.append(item)

    def get_current_dns(self) -> Optional[str]:
        """Get current DNS configuration for the selected interface"""
        if not self.current_interface:
            return None

        try:
            if sys.platform == "darwin":
                # macOS using networksetup
                result = subprocess.run(
                    ["networksetup", "-getdnsservers", self.current_interface],
                    capture_output=True,
                    text=True,
                    check=True
                )
                dns_output = result.stdout.strip()

                if "There aren't any DNS Servers set" in dns_output or not dns_output:
                    return "DHCP (Automatic)"

                # Parse DNS servers
                dns_servers = [line.strip() for line in dns_output.split('\n') if line.strip()]
                return ", ".join(dns_servers)

            elif sys.platform.startswith("linux"):
                # Linux using nmcli
                result = subprocess.run(
                    ["nmcli", "-t", "-f", "IP4.DNS", "con", "show", self.current_interface],
                    capture_output=True,
                    text=True,
                    check=True
                )

                dns_servers = []
                for line in result.stdout.split('\n'):
                    if line.startswith('IP4.DNS'):
                        dns = line.split(':')[1].strip()
                        if dns:
                            dns_servers.append(dns)

                if not dns_servers:
                    return "DHCP (Automatic)"

                return ", ".join(dns_servers)

            return None

        except subprocess.CalledProcessError:
            return None
        except Exception:
            return None

    def get_network_interfaces(self) -> List[str]:
        """Get list of available network interfaces"""
        try:
            if sys.platform == "darwin":
                # macOS
                result = subprocess.run(
                    ["networksetup", "-listallnetworkservices"],
                    capture_output=True,
                    text=True
                )
                interfaces = [
                    line.strip()
                    for line in result.stdout.split('\n')[1:]  # Skip first line
                    if line.strip() and not line.startswith('*')
                ]
            elif sys.platform.startswith("linux"):
                # Linux
                result = subprocess.run(
                    ["nmcli", "-t", "-f", "NAME", "connection", "show"],
                    capture_output=True,
                    text=True
                )
                interfaces = [
                    line.strip()
                    for line in result.stdout.split('\n')
                    if line.strip()
                ]
            else:
                interfaces = []

            return interfaces
        except Exception as e:
            self.notify(f"Error getting interfaces: {e}", severity="error")
            return []

    def select_interface(self) -> None:
        """Show interface selection screen"""
        interfaces = self.get_network_interfaces()
        if not interfaces:
            self.notify("No network interfaces found!", severity="error")
            return

        def on_select(interface: str):
            self.current_interface = interface
            title = self.query_one("#title", Static)
            title.update(f"DNS Switch [{self.current_interface}]")
            self.notify(f"Interface: {interface}", severity="information")

        self.push_screen(InterfaceSelectionScreen(interfaces, on_select))

    def apply_dns(self, profile: DNSProfile) -> bool:
        """Apply DNS settings to the network interface"""
        if not self.current_interface:
            self.notify("Please select a network interface first!", severity="error")
            return False

        try:
            if sys.platform == "darwin":
                # macOS using networksetup
                if profile.primary == "auto":
                    # Set to DHCP
                    subprocess.run(
                        ["networksetup", "-setdnsservers", self.current_interface, "Empty"],
                        check=True
                    )
                else:
                    # Set custom DNS
                    dns_servers = [profile.primary]
                    if profile.secondary:
                        dns_servers.append(profile.secondary)

                    subprocess.run(
                        ["networksetup", "-setdnsservers", self.current_interface] + dns_servers,
                        check=True
                    )

            elif sys.platform.startswith("linux"):
                # Linux using nmcli
                if profile.primary == "auto":
                    subprocess.run(
                        ["nmcli", "con", "mod", self.current_interface, "ipv4.dns", ""],
                        check=True
                    )
                    subprocess.run(
                        ["nmcli", "con", "mod", self.current_interface, "ipv4.ignore-auto-dns", "no"],
                        check=True
                    )
                else:
                    dns_servers = profile.primary
                    if profile.secondary:
                        dns_servers += f" {profile.secondary}"

                    subprocess.run(
                        ["nmcli", "con", "mod", self.current_interface, "ipv4.dns", dns_servers],
                        check=True
                    )
                    subprocess.run(
                        ["nmcli", "con", "mod", self.current_interface, "ipv4.ignore-auto-dns", "yes"],
                        check=True
                    )

                # Restart connection
                subprocess.run(
                    ["nmcli", "con", "down", self.current_interface],
                    check=False  # Don't fail if already down
                )
                subprocess.run(
                    ["nmcli", "con", "up", self.current_interface],
                    check=True
                )

            return True

        except subprocess.CalledProcessError as e:
            self.notify(f"Error applying DNS: {e}", severity="error")
            return False
        except Exception as e:
            self.notify(f"Unexpected error: {e}", severity="error")
            return False

    @on(Button.Pressed, "#apply-btn")
    def on_apply_button(self) -> None:
        """Handle apply button press"""
        profile_list = self.query_one("#profile-list", ListView)

        if profile_list.index is None:
            self.notify("Please select a DNS profile first!", severity="warning")
            return

        profile_key = list(self.profiles.keys())[profile_list.index]
        profile = self.profiles[profile_key]

        status = self.query_one("#status", Static)
        status.update(f"Applying {profile.name}...")

        if self.apply_dns(profile):
            status.update(f"✓ Successfully applied: {profile.name}")
            status.remove_class("error")
            status.add_class("success")
            self.notify(f"DNS switched to {profile.name}", severity="information")
        else:
            status.update(f"✗ Failed to apply: {profile.name}")
            status.remove_class("success")
            status.add_class("error")

    @on(Button.Pressed, "#interface-btn")
    def on_interface_button(self) -> None:
        """Handle interface change button"""
        self.select_interface()

    @on(Button.Pressed, "#check-btn")
    def on_check_button(self) -> None:
        """Handle check DNS button"""
        self.action_check_dns()

    @on(Button.Pressed, "#refresh-btn")
    def action_refresh(self) -> None:
        """Refresh the configuration"""
        self.load_config()
        self.refresh_profile_list()
        self.notify("Configuration refreshed", severity="information")

    def action_check_dns(self) -> None:
        """Check current DNS configuration"""
        if not self.current_interface:
            self.notify("Please select a network interface first!", severity="warning")
            return

        dns = self.get_current_dns()
        status = self.query_one("#status", Static)

        if dns:
            status.update(f"Current DNS: {dns}")
            status.remove_class("error")
            status.add_class("success")
            self.notify(f"Current DNS: {dns}", severity="information", timeout=5)
        else:
            status.update("Unable to retrieve DNS configuration")
            status.remove_class("success")
            status.add_class("error")
            self.notify("Unable to retrieve DNS configuration", severity="error")

    def action_change_interface(self) -> None:
        """Change network interface"""
        self.select_interface()


def main():
    """Main entry point"""
    parser = argparse.ArgumentParser(
        prog="dns-switch",
        description="A user-friendly TUI for quickly switching between DNS configurations",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  dns-switch              Launch the TUI
  dns-switch --help       Show this help message
  dns-switch --version    Show version information

Configuration:
  Default config location: ~/.config/dns-switch/config.yaml

Keyboard Shortcuts:
  Arrow Keys / Mouse      Navigate through DNS profiles
  Enter / Apply Button    Apply selected DNS profile
  c                       Check current DNS configuration
  i                       Change network interface
  r                       Refresh configuration
  q                       Quit application

Note: Requires sudo to modify network settings
      """
    )
    parser.add_argument(
        "--version",
        action="version",
        version=f"%(prog)s {__version__}"
    )

    args = parser.parse_args()

    # If we get here, no flags were provided, so launch the TUI
    app = DNSSwitchApp()
    app.run()


if __name__ == "__main__":
    main()
