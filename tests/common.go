package tests

import (
	"os"
	"testing"
)

var HOME = os.Getenv("HOME")

func GetConfigDir(t *testing.T) string {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		t.Fatal("HOME environment variable is not set")
	}
	configDir := homeDir + "/.config/browsir"

	return configDir
}

func CleanUp(t *testing.T, configDir string) {
	// Clean up the mocked config after the test
	t.Cleanup(func() {
		err := os.RemoveAll(configDir)
		if err != nil {
			t.Fatalf("Error cleaning up config directory: %v", err)
		}
	})
}

func SetupEmptyEnvs() {
	os.Setenv("HOME", "/dev/null")
}

func CleanUpEnvs() {
	os.Setenv("HOME", HOME)
}
