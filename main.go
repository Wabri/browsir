package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

type Profile struct {
	Name        string `yaml:"name"`
	ProfileDir  string `yaml:"profile_dir"`
	Description string `yaml:"description"`
}

type Config struct {
	AppName     string            `yaml:"app_name"`
	BrowserName string            `yaml:"browser_name"`
	Profiles    []Profile         `yaml:"profiles"`
	Shortcuts   map[string]string `yaml:"shortcuts"`
}

func loadConfig() Config {
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

func loadLocalShortcuts() map[string]string {
	shortcuts := make(map[string]string)

	home, err := os.UserHomeDir()
	if err != nil {
		return shortcuts
	}

	// First check system config directory
	shortcutsPath := "/etc/browsir/shortcuts"

	// Then check XDG config directory
	if _, err := os.Stat(shortcutsPath); os.IsNotExist(err) {
		xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfigHome == "" {
			xdgConfigHome = filepath.Join(home, ".config")
		}
		shortcutsPath = filepath.Join(xdgConfigHome, "browsir", "shortcuts")
	}

	// Then check home directory
	if _, err := os.Stat(shortcutsPath); os.IsNotExist(err) {
		shortcutsPath = filepath.Join(home, ".browsir_shortcuts")
	}

	// Finally check current directory
	if _, err := os.Stat(shortcutsPath); os.IsNotExist(err) {
		shortcutsPath = "shortcuts"
	}

	data, err := os.ReadFile(shortcutsPath)
	if err != nil {
		return shortcuts
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			shortcuts[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return shortcuts
}

func saveLocalShortcut(shortcut, url string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// First try system config directory
	shortcutsPath := "/etc/browsir/shortcuts"
	f, err := os.OpenFile(shortcutsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		defer f.Close()
		_, err = fmt.Fprintf(f, "%s=%s\n", shortcut, url)
		return err
	}

	// Then try XDG config directory
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		xdgConfigHome = filepath.Join(home, ".config")
	}
	configDir := filepath.Join(xdgConfigHome, "browsir")
	if err := os.MkdirAll(configDir, 0755); err == nil {
		shortcutsPath := filepath.Join(configDir, "shortcuts")
		f, err := os.OpenFile(shortcutsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			defer f.Close()
			_, err = fmt.Fprintf(f, "%s=%s\n", shortcut, url)
			return err
		}
	}

	// Then try home directory
	shortcutsPath = filepath.Join(home, ".browsir_shortcuts")
	f, err = os.OpenFile(shortcutsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// Finally try current directory
		f, err = os.OpenFile("shortcuts", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "%s=%s\n", shortcut, url)
	return err
}

func getBrowserPath(browserName string) (string, error) {
	switch runtime.GOOS {
	case "darwin":
		switch browserName {
		case "chrome":
			return "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", nil
		case "brave":
			return "/Applications/Brave Browser.app/Contents/MacOS/Brave Browser", nil
		case "arc":
			return "/Applications/Arc.app/Contents/MacOS/Arc", nil
		}
	case "linux":
		switch browserName {
		case "chrome":
			return "google-chrome", nil
		case "brave":
			return "brave-browser", nil
		}
	case "windows":
		switch browserName {
		case "chrome":
			return filepath.Join(os.Getenv("ProgramFiles"), "Google", "Chrome", "Application", "chrome.exe"), nil
		case "brave":
			return filepath.Join(os.Getenv("ProgramFiles"), "BraveSoftware", "Brave-Browser", "Application", "brave.exe"), nil
		}
	}
	return "", fmt.Errorf("unsupported browser %s on %s", browserName, runtime.GOOS)
}

func openBrowser(browserName string, profile Profile, url string) error {
	browserPath, err := getBrowserPath(browserName)
	if err != nil {
		return err
	}

	args := []string{"--profile-directory=" + profile.ProfileDir}
	if url != "" {
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "https://" + url
		}
		args = append(args, url)
	}

	cmd := exec.Command(browserPath, args...)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start browser: %v", err)
	}

	return nil
}

func printUsage(profiles []Profile, shortcuts map[string]string, localShortcuts map[string]string) {
	fmt.Printf("browsir v1.0.0\n\n")
	fmt.Println("Usage: browsir [profile] [url|shortcut]")
	fmt.Println("\nProfiles:")
	for _, p := range profiles {
		fmt.Printf("  %-12s - %s\n", p.Name, p.Description)
	}
	fmt.Println("\nShortcuts:")
	for shortcut, url := range shortcuts {
		fmt.Printf("  %-12s -> %s\n", shortcut, url)
	}
	if len(localShortcuts) > 0 {
		fmt.Println("\nLocal Shortcuts:")
		for shortcut, url := range localShortcuts {
			fmt.Printf("  %-12s -> %s\n", shortcut, url)
		}
	}
	fmt.Println("\nExamples:")
	fmt.Println("  browsir work                    # Open browser with work profile")
	fmt.Println("  browsir personal gmail.com      # Open Gmail with personal profile")
	fmt.Println("  browsir default mail           # Open mail shortcut with default profile")
}

func findSimilarShortcuts(input string, shortcuts, localShortcuts map[string]string) []string {
	var similar []string
	for shortcut := range shortcuts {
		if strings.Contains(shortcut, input) || strings.Contains(input, shortcut) {
			similar = append(similar, shortcut)
		}
	}
	for shortcut := range localShortcuts {
		if strings.Contains(shortcut, input) || strings.Contains(input, shortcut) {
			similar = append(similar, shortcut)
		}
	}
	return similar
}

func promptYesNo(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s (y/n): ", prompt)
		response, _ := reader.ReadString('\n')
		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" || response == "yes" {
			return true
		}
		if response == "n" || response == "no" {
			return false
		}
	}
}

func main() {
	config := loadConfig()
	localShortcuts := loadLocalShortcuts()

	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		printUsage(config.Profiles, config.Shortcuts, localShortcuts)
		os.Exit(0)
	}

	if os.Args[1] == "-v" || os.Args[1] == "--version" {
		fmt.Println("browsir v1.0.0")
		os.Exit(0)
	}

	var profileName string
	var url string
	args := os.Args[1:]

	// Check if first arg contains a dot or is a shortcut
	if strings.Contains(args[0], ".") || strings.HasPrefix(args[0], "http") ||
		localShortcuts[args[0]] != "" || config.Shortcuts[args[0]] != "" {
		profileName = "default"
		url = args[0]
	} else {
		profileName = args[0]
		if len(args) > 1 {
			url = args[1]
		}
	}

	var selectedProfile Profile
	var found bool
	for _, p := range config.Profiles {
		if p.Name == profileName {
			selectedProfile = p
			found = true
			break
		}
	}

	if !found {
		fmt.Fprintf(os.Stderr, "Error: unknown profile: %s\n", profileName)
		printUsage(config.Profiles, config.Shortcuts, localShortcuts)
		os.Exit(1)
	}

	if url != "" {
		// If the argument contains a dot or protocol, treat it as a direct URL
		if !strings.Contains(url, ".") && !strings.HasPrefix(url, "http") {
			// Check local shortcuts first
			if localURL, exists := localShortcuts[url]; exists {
				url = localURL
			} else if configURL, exists := config.Shortcuts[url]; exists {
				url = configURL
			} else {
				// Check for similar shortcuts
				similar := findSimilarShortcuts(url, config.Shortcuts, localShortcuts)
				if len(similar) > 0 {
					fmt.Printf("\033[33mDid you mean one of these shortcuts?\033[0m\n")
					for _, s := range similar {
						if u, exists := localShortcuts[s]; exists {
							fmt.Printf("\033[36m  %s\033[0m -> %s (local)\n", s, u)
						} else {
							fmt.Printf("\033[36m  %s\033[0m -> %s\n", s, config.Shortcuts[s])
						}
					}
					os.Exit(1)
				}

				// If no similar shortcuts found, ask if they want to save it
				if promptYesNo("Would you like to save this as a shortcut?") {
					fmt.Print("Enter the website URL: ")
					reader := bufio.NewReader(os.Stdin)
					websiteURL, _ := reader.ReadString('\n')
					websiteURL = strings.TrimSpace(websiteURL)

					if err := saveLocalShortcut(url, websiteURL); err != nil {
						fmt.Fprintf(os.Stderr, "Error saving shortcut: %v\n", err)
					} else {
						fmt.Printf("\033[32mShortcut saved: %s -> %s\033[0m\n", url, websiteURL)
					}
					os.Exit(0)
				} else {
					fmt.Println("\033[32mTip: You can add shortcuts in your .browsir.yml config file or local shortcuts file\033[0m")
					os.Exit(1)
				}
			}
		}
	}

	if err := openBrowser(config.BrowserName, selectedProfile, url); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
