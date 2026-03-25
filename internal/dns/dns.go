package dns

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/pinaka-io/dns-switch/internal/config"
)

// GetNetworkInterfaces returns a list of available network interfaces
func GetNetworkInterfaces() ([]string, error) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		// macOS
		cmd = exec.Command("networksetup", "-listallnetworkservices")
	case "linux":
		// Linux
		cmd = exec.Command("nmcli", "-t", "-f", "NAME", "connection", "show")
	default:
		return nil, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var interfaces []string

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Skip first line on macOS (header)
		if runtime.GOOS == "darwin" && i == 0 {
			continue
		}

		// Skip disabled interfaces on macOS
		if runtime.GOOS == "darwin" && strings.HasPrefix(line, "*") {
			continue
		}

		interfaces = append(interfaces, line)
	}

	return interfaces, nil
}

// GetCurrentDNS retrieves the current DNS configuration for the interface
func GetCurrentDNS(iface string) (string, error) {
	if iface == "" {
		return "", fmt.Errorf("no interface specified")
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("networksetup", "-getdnsservers", iface)
	case "linux":
		cmd = exec.Command("nmcli", "-t", "-f", "IP4.DNS", "con", "show", iface)
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get DNS servers: %w", err)
	}

	result := strings.TrimSpace(string(output))

	if runtime.GOOS == "darwin" {
		if strings.Contains(result, "There aren't any DNS Servers set") || result == "" {
			return "DHCP (Automatic)", nil
		}
		return strings.ReplaceAll(result, "\n", ", "), nil
	}

	// Linux
	var dnsServers []string
	for _, line := range strings.Split(result, "\n") {
		if strings.HasPrefix(line, "IP4.DNS") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				dns := strings.TrimSpace(parts[1])
				if dns != "" {
					dnsServers = append(dnsServers, dns)
				}
			}
		}
	}

	if len(dnsServers) == 0 {
		return "DHCP (Automatic)", nil
	}

	return strings.Join(dnsServers, ", "), nil
}

// ApplyDNS applies the DNS profile to the specified interface
func ApplyDNS(iface string, profile config.DNSProfile) error {
	if iface == "" {
		return fmt.Errorf("no interface specified")
	}

	switch runtime.GOOS {
	case "darwin":
		return applyDNSMacOS(iface, profile)
	case "linux":
		return applyDNSLinux(iface, profile)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func applyDNSMacOS(iface string, profile config.DNSProfile) error {
	var cmd *exec.Cmd

	if profile.Primary == "auto" {
		// Set to DHCP
		cmd = exec.Command("networksetup", "-setdnsservers", iface, "Empty")
	} else {
		// Set custom DNS
		args := []string{"-setdnsservers", iface, profile.Primary}
		if profile.Secondary != "" && profile.Secondary != profile.Primary {
			args = append(args, profile.Secondary)
		}
		cmd = exec.Command("networksetup", args...)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to apply DNS settings: %w", err)
	}

	return nil
}

func applyDNSLinux(iface string, profile config.DNSProfile) error {
	if profile.Primary == "auto" {
		// Set to DHCP
		cmd := exec.Command("nmcli", "con", "mod", iface, "ipv4.dns", "")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to clear DNS: %w", err)
		}

		cmd = exec.Command("nmcli", "con", "mod", iface, "ipv4.ignore-auto-dns", "no")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to enable auto DNS: %w", err)
		}
	} else {
		// Set custom DNS
		dnsServers := profile.Primary
		if profile.Secondary != "" && profile.Secondary != profile.Primary {
			dnsServers += " " + profile.Secondary
		}

		cmd := exec.Command("nmcli", "con", "mod", iface, "ipv4.dns", dnsServers)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to set DNS: %w", err)
		}

		cmd = exec.Command("nmcli", "con", "mod", iface, "ipv4.ignore-auto-dns", "yes")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to disable auto DNS: %w", err)
		}
	}

	// Restart connection
	_ = exec.Command("nmcli", "con", "down", iface).Run() // Ignore error if already down

	cmd := exec.Command("nmcli", "con", "up", iface)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to restart connection: %w", err)
	}

	return nil
}
