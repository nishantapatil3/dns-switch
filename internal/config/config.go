package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"
)

// DNSProfile represents a DNS configuration profile
type DNSProfile struct {
	Key         string
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Primary     string `yaml:"primary"`
	Secondary   string `yaml:"secondary"`
}

// Config represents the application configuration
type Config struct {
	DNSProfiles      map[string]DNSProfile `yaml:"dns_profiles"`
	NetworkInterface string                `yaml:"network_interface"`
}

// defaultConfigYAML is the default configuration content
const defaultConfigYAML = `# DNS Switch Configuration
# Add your DNS profiles here

dns_profiles:
  # Adguard DNS
  adguard:
    name: "AdGuard DNS"
    description: "AdGuard DNS for blocking ads and trackers"
    primary: "192.168.4.69"
    secondary: ""

  # Cloudflare DNS
  cloudflare:
    name: "Cloudflare (1.1.1.1)"
    description: "Fast and privacy-focused DNS"
    primary: "1.1.1.1"
    secondary: "1.0.0.1"

  # Google DNS
  google:
    name: "Google Public DNS"
    description: "Reliable DNS by Google"
    primary: "8.8.8.8"
    secondary: "8.8.4.4"

  # Quad9 DNS
  quad9:
    name: "Quad9 (Security)"
    description: "DNS with malware blocking"
    primary: "9.9.9.9"
    secondary: "149.112.112.112"

  # OpenDNS
  opendns:
    name: "OpenDNS"
    description: "Fast and reliable DNS"
    primary: "208.67.222.222"
    secondary: "208.67.220.220"

  # DHCP (Automatic)
  dhcp:
    name: "DHCP (Automatic)"
    description: "Use DNS provided by your network"
    primary: "auto"
    secondary: "auto"

# Network interface to configure (e.g., "Wi-Fi", "Ethernet", "en0")
# Leave empty to show a list of available interfaces
network_interface: ""
`

// LoadConfig loads the configuration from config.yaml
func LoadConfig() (*Config, error) {
	// Try multiple config locations
	configPaths := []string{
		filepath.Join(os.Getenv("HOME"), ".config", "dns-switch", "config.yaml"),
		"config.yaml",
	}

	var configPath string
	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			configPath = path
			break
		}
	}

	if configPath == "" {
		// Create default config in user's home directory
		configPath = filepath.Join(os.Getenv("HOME"), ".config", "dns-switch", "config.yaml")
		if err := createDefaultConfig(configPath); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Add keys to profiles
	for key, profile := range config.DNSProfiles {
		profile.Key = key
		config.DNSProfiles[key] = profile
	}

	return &config, nil
}

// createDefaultConfig creates the default configuration file
func createDefaultConfig(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(path, []byte(defaultConfigYAML), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetProfiles returns a slice of DNS profiles with DHCP first, then the rest sorted alphabetically by name
func (c *Config) GetProfiles() []DNSProfile {
	var dhcpProfile DNSProfile
	profiles := make([]DNSProfile, 0, len(c.DNSProfiles)-1)

	for _, profile := range c.DNSProfiles {
		if profile.Key == "dhcp" {
			dhcpProfile = profile
		} else {
			profiles = append(profiles, profile)
		}
	}

	// Sort non-DHCP profiles alphabetically by name for consistent display
	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].Name < profiles[j].Name
	})

	// Put DHCP at the top
	if dhcpProfile.Key != "" {
		profiles = append([]DNSProfile{dhcpProfile}, profiles...)
	}

	return profiles
}
