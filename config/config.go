package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type Config struct {
	AppName     string            `yaml:"app_name"`
	BrowserName string            `yaml:"browser_name"`
	Profiles    []Profile         `yaml:"profiles"`
	Shortcuts   map[string]string `yaml:"shortcuts"`
}

type Profile struct {
	Name        string `yaml:"name"`
	ProfileDir  string `yaml:"profile_dir"`
	Description string `yaml:"description"`
}

func LoadConfig() Config {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	// First check system config directory
	configPath := "/etc/browsir/config.yml"

	// Then check XDG config directory
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfigHome == "" {
			xdgConfigHome = filepath.Join(home, ".config")
		}
		configPath = filepath.Join(xdgConfigHome, "browsir", "config.yml")
	}

	// Then check home directory
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = filepath.Join(home, ".browsir.yml")
	}

	// Finally check current directory
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = ".browsir.yml"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		// Return default config if file doesn't exist
		return Config{
			AppName:     "browsir",
			BrowserName: "chrome",
			Profiles: []Profile{
				{Name: "default", ProfileDir: "Default", Description: "Default profile"},
			},
			Shortcuts: map[string]string{
				"cal": "calendar.google.com",
			},
		}
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing config file: %v\n", err)
		os.Exit(1)
	}

	if config.AppName == "" {
		config.AppName = "browsir"
	}
	if config.BrowserName == "" {
		config.BrowserName = "chrome"
	}
	if len(config.Profiles) == 0 {
		config.Profiles = []Profile{
			{Name: "default", ProfileDir: "Default", Description: "Default profile"},
		}
	}
	if config.Shortcuts == nil {
		config.Shortcuts = map[string]string{
			"cal": "calendar.google.com",
		}
	}

	return config
}
