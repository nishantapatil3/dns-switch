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
		return nil, fmt.Errorf("config file not found in any of the default locations")
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

// GetProfiles returns a slice of DNS profiles sorted alphabetically by name
func (c *Config) GetProfiles() []DNSProfile {
	profiles := make([]DNSProfile, 0, len(c.DNSProfiles))
	for _, profile := range c.DNSProfiles {
		profiles = append(profiles, profile)
	}

	// Sort profiles alphabetically by name for consistent display
	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].Name < profiles[j].Name
	})

	return profiles
}
