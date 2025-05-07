package config

import (
	"fmt"
	"os"
	"testing"
)

// /Users/404answernotfound
var HOME = os.Getenv("HOME")

func TestConfig(t *testing.T) {

	tc := struct {
		got  Config
		want Config
	}{
		got: Config{},
		want: Config{
			AppName:     "browsir",
			BrowserName: "chrome",
			Profiles: []Profile{
				{Name: "personal", ProfileDir: "Default", Description: "Default profile"},
				{Name: "work", ProfileDir: "Profile 1", Description: "Work profile"},
			},
			Shortcuts: map[string]string{
				"google": "google.com",
				"github": "github.com",
				"mail":   "gmail.com",
			},
		},
	}

	t.Run("Test config loading", func(t *testing.T) {
		setupConfigDir(t)
		config, err := LoadConfig()
		if err != nil {
			t.Fatalf("Error loading config: %v", err)
		}

		tc.got = config
		if tc.got.AppName != tc.want.AppName {
			t.Errorf("got %v, want %v", tc.got.AppName, tc.want.AppName)
		}
		if tc.got.BrowserName != tc.want.BrowserName {
			t.Errorf("got %v, want %v", tc.got.BrowserName, tc.want.BrowserName)
		}
		if len(tc.got.Profiles) != len(tc.want.Profiles) {
			t.Errorf("got %v, want %v", len(tc.got.Profiles), len(tc.want.Profiles))
		}

		configDir := getConfigDir(t)
		cleanUp(t, configDir)
	})

	t.Run("Test empty configuration should return error", func(t *testing.T) {
		setupEmptyEnvs()
		_, err := LoadConfig()
		if err != nil {
			fmt.Println(err)
		}

		cleanUpEnvs()
	})

}

func setupConfigDir(t *testing.T) {
	// Set up config files as would be done by installing browsir
	configDir := getConfigDir(t)
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("Error creating config directory: %v", err)
	}
	configFile := configDir + "/config.yml"
	exampleFile, err := os.ReadFile("../config.example.yml")
	if err != nil {
		t.Fatalf("Error reading example config file: %v", err)
	}
	err = os.WriteFile(configFile, exampleFile, 0644)
	if err != nil {
		t.Fatalf("Error writing config file: %v", err)
	}
}

func getConfigDir(t *testing.T) string {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		t.Fatal("HOME environment variable is not set")
	}
	configDir := homeDir + "/.config/browsir"

	return configDir
}

func cleanUp(t *testing.T, configDir string) {
	// Clean up the mocked config after the test
	t.Cleanup(func() {
		err := os.RemoveAll(configDir)
		if err != nil {
			t.Fatalf("Error cleaning up config directory: %v", err)
		}
	})
}

func setupEmptyEnvs() {
	os.Setenv("HOME", "/dev/null")
}

func cleanUpEnvs() {
	os.Setenv("HOME", HOME)
}
