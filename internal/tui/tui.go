package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/pinaka-io/dns-switch/internal/config"
	"github.com/pinaka-io/dns-switch/internal/dns"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	statusStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#5A5A5A")).
			Padding(0, 2).
			MarginTop(1)

	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#04B575")).
			Padding(0, 2).
			MarginTop(1)

	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#FF5F87")).
			Padding(0, 2).
			MarginTop(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Padding(1, 0)

	docStyle = lipgloss.NewStyle().
			Padding(1, 2)
)

type viewMode int

const (
	profileListView viewMode = iota
	interfaceSelectionView
)

type profileItem struct {
	profile config.DNSProfile
}

func (i profileItem) Title() string {
	return i.profile.Name
}

func (i profileItem) Description() string {
	dnsInfo := i.profile.Primary
	if i.profile.Secondary != "" && i.profile.Secondary != i.profile.Primary {
		dnsInfo += ", " + i.profile.Secondary
	}
	return fmt.Sprintf("%s (%s)", i.profile.Description, dnsInfo)
}

func (i profileItem) FilterValue() string {
	return i.profile.Name
}

type interfaceItem struct {
	name string
}

func (i interfaceItem) Title() string       { return i.name }
func (i interfaceItem) Description() string { return "" }
func (i interfaceItem) FilterValue() string { return i.name }

// Model represents the TUI application model
type Model struct {
	config     *config.Config
	list       list.Model
	mode       viewMode
	iface      string
	status     string
	statusType string // "", "success", "error"
	width      int
	height     int
	interfaces []string
	quitting   bool
}

// NewModel creates a new TUI model
func NewModel() (Model, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return Model{}, fmt.Errorf("failed to load config: %w", err)
	}

	profiles := cfg.GetProfiles()
	items := make([]list.Item, len(profiles))
	for i, p := range profiles {
		items[i] = profileItem{profile: p}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true).
		Padding(0, 0, 0, 1)

	delegate.Styles.SelectedDesc = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Foreground(lipgloss.Color("#626262")).
		Padding(0, 0, 0, 1)

	l := list.New(items, delegate, 0, 0)
	l.Title = "DNS Switch"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle

	m := Model{
		config:     cfg,
		list:       l,
		mode:       profileListView,
		iface:      cfg.NetworkInterface,
		status:     "Select profile and press Enter to apply",
		statusType: "",
	}

	return m, nil
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-5)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "r":
			if m.mode == profileListView {
				return m.handleRefresh()
			}

		case "i":
			if m.mode == profileListView {
				return m.switchToInterfaceSelection()
			}

		case "c":
			if m.mode == profileListView {
				return m.checkCurrentDNS()
			}

		case "enter":
			return m.handleEnter()

		case "esc":
			if m.mode == interfaceSelectionView {
				m.quitting = true
				return m, tea.Quit
			} else if m.mode == profileListView {
				return m.switchToInterfaceSelection()
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	if m.mode == interfaceSelectionView {
		if item, ok := m.list.SelectedItem().(interfaceItem); ok {
			m.iface = item.name
			m.status = fmt.Sprintf("Interface: %s", m.iface)
			m.statusType = "success"
			return m.switchToProfileList()
		}
	} else if m.mode == profileListView {
		if m.iface == "" {
			m.status = "Please select a network interface first! Press 'i'"
			m.statusType = "error"
			return m, nil
		}

		if item, ok := m.list.SelectedItem().(profileItem); ok {
			return m.applyDNSProfile(item.profile)
		}
	}
	return m, nil
}

func (m Model) applyDNSProfile(profile config.DNSProfile) (tea.Model, tea.Cmd) {
	m.status = fmt.Sprintf("Applying %s...", profile.Name)
	m.statusType = ""

	err := dns.ApplyDNS(m.iface, profile)
	if err != nil {
		m.status = fmt.Sprintf("✗ Failed to apply %s: %s", profile.Name, err.Error())
		m.statusType = "error"
	} else {
		m.status = fmt.Sprintf("✓ Successfully applied: %s", profile.Name)
		m.statusType = "success"
	}

	return m, nil
}

func (m Model) checkCurrentDNS() (tea.Model, tea.Cmd) {
	if m.iface == "" {
		m.status = "Please select a network interface first! Press 'i'"
		m.statusType = "error"
		return m, nil
	}

	dnsServers, err := dns.GetCurrentDNS(m.iface)
	if err != nil {
		m.status = fmt.Sprintf("Unable to retrieve DNS: %s", err.Error())
		m.statusType = "error"
	} else {
		m.status = fmt.Sprintf("Current DNS: %s", dnsServers)
		m.statusType = "success"
	}

	return m, nil
}

func (m Model) handleRefresh() (tea.Model, tea.Cmd) {
	cfg, err := config.LoadConfig()
	if err != nil {
		m.status = fmt.Sprintf("Failed to reload config: %s", err.Error())
		m.statusType = "error"
		return m, nil
	}

	m.config = cfg
	profiles := cfg.GetProfiles()
	items := make([]list.Item, len(profiles))
	for i, p := range profiles {
		items[i] = profileItem{profile: p}
	}
	m.list.SetItems(items)

	m.status = "Configuration refreshed"
	m.statusType = "success"
	return m, nil
}

func (m Model) switchToInterfaceSelection() (tea.Model, tea.Cmd) {
	interfaces, err := dns.GetNetworkInterfaces()
	if err != nil {
		m.status = fmt.Sprintf("Failed to get interfaces: %s", err.Error())
		m.statusType = "error"
		return m, nil
	}

	if len(interfaces) == 0 {
		m.status = "No network interfaces found!"
		m.statusType = "error"
		return m, nil
	}

	m.interfaces = interfaces
	items := make([]list.Item, len(interfaces))
	for i, iface := range interfaces {
		items[i] = interfaceItem{name: iface}
	}

	m.list.SetItems(items)
	m.list.Title = "Select Network Interface"
	m.mode = interfaceSelectionView
	m.status = "Select interface and press Enter"
	m.statusType = ""

	return m, nil
}

func (m Model) switchToProfileList() (tea.Model, tea.Cmd) {
	profiles := m.config.GetProfiles()
	items := make([]list.Item, len(profiles))
	for i, p := range profiles {
		items[i] = profileItem{profile: p}
	}

	m.list.SetItems(items)
	title := "DNS Switch"
	if m.iface != "" {
		title += fmt.Sprintf(" [%s]", m.iface)
	}
	m.list.Title = title
	m.mode = profileListView

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var statusBar string
	switch m.statusType {
	case "success":
		statusBar = successStyle.Render(m.status)
	case "error":
		statusBar = errorStyle.Render(m.status)
	default:
		statusBar = statusStyle.Render(m.status)
	}

	var helpText string
	if m.mode == interfaceSelectionView {
		helpText = helpStyle.Render("↑/↓: navigate • enter: select • esc/q: quit")
	} else {
		helpText = helpStyle.Render("↑/↓: navigate • enter: apply • c: check • i: interface • r: refresh • esc: back • q: quit")
	}

	return docStyle.Render(
		m.list.View() + "\n" +
			statusBar + "\n" +
			helpText,
	)
}

// SwitchToInterfaceSelectionIfNeeded switches to interface selection if no interface is set
func (m Model) SwitchToInterfaceSelectionIfNeeded() Model {
	if m.iface == "" {
		updatedModel, _ := m.switchToInterfaceSelection()
		return updatedModel.(Model)
	}
	return m
}
