package config

import (
	"os"
	"testing"
)

type TableTest struct {
	name string
	want Config
	got  Config
}

func TestConfig(t *testing.T) {

	tt := []TableTest{
		{
			name: "Test config loading",
			want: Config{
				AppName:     "browsir",
				BrowserName: "chrome",
				Profiles: []Profile{
					{Name: "personal", ProfileDir: "Default", Description: "Default profile"},
					{Name: "work", ProfileDir: "Profile 1", Description: "Work profile"},
				},
				Shortcuts: map[string]string{
					"cal": "calendar.google.com",
				},
			},
			got: Config{},
		},
		{
			name: "Test config loading with empty profile",
			want: Config{
				AppName:     "browsir",
				BrowserName: "chrome",
				Profiles: []Profile{
					{Name: "personal", ProfileDir: "Default", Description: "Default profile"},
					{Name: "work", ProfileDir: "Profile 1", Description: "Work profile"},
				},
				Shortcuts: map[string]string{
					"cal": "calendar.google.com",
				},
			},
			got: Config{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			TestSetupConfigDir(t)
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
		})
	}

}

func TestSetupConfigDir(t *testing.T) {
	// Set up config files as would be done by installing browsir
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		t.Fatal("HOME environment variable is not set")
	}
	configDir := homeDir + "/.config/browsir"
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

	// Clean up the mocked config after the test
	t.Cleanup(func() {
		err := os.RemoveAll(configDir)
		if err != nil {
			t.Fatalf("Error cleaning up config directory: %v", err)
		}
	})
}
