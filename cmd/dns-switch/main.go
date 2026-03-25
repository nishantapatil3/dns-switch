package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/pinaka-io/dns-switch/internal/tui"
)

var version = "dev" // Set by GoReleaser via ldflags

func main() {
	// Check for version flag
	for _, arg := range os.Args[1:] {
		if arg == "--version" || arg == "-v" {
			fmt.Printf("dns-switch %s\n", version)
			return
		}
		if arg == "--help" || arg == "-h" {
			printHelp()
			return
		}
	}

	m, err := tui.NewModel()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Check if interface needs to be selected
	m = m.SwitchToInterfaceSelectionIfNeeded()

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	help := `dns-switch - A user-friendly TUI for quickly switching between DNS configurations

Usage:
  dns-switch              Launch the TUI
  dns-switch --help       Show this help message
  dns-switch --version    Show version information

Configuration:
  Default config location: ~/.config/dns-switch/config.yaml

Keyboard Shortcuts:
  ↑/↓ or j/k             Navigate through options
  Enter                   Apply selected DNS profile / Select interface
  c                       Check current DNS configuration
  i                       Change network interface
  r                       Refresh configuration
  q                       Quit application
  Esc                     Back to interface selection / Quit

Note: Requires sudo to modify network settings

Examples:
  sudo dns-switch         Launch the application

For more information, visit: https://github.com/pinaka-io/dns-switch
`
	fmt.Print(help)
}
