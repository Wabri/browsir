package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
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

func findConfigFile() (string, error) {
  configHome := os.Getenv("XDG_CONFIG_HOME")
  if configHome == "" {
    configHome = os.Getenv("HOME") + "/.config"
  }
  configPath := configHome + "/browsir/config.yml"

  if _, err := os.Stat(configPath); err == nil {
    return configPath, nil
  }

  // Fallback to /etc
  configPath = "/etc/browsir/config.yml"
  if _, err := os.Stat(configPath); err == nil {
    return configPath, nil
  }

  return "", fmt.Errorf("config file not found")
}

func LoadConfig() Config {
    configPath, err := findConfigFile()
	if err != nil {
      fmt.Println("No config file found")
      return Config{}
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
