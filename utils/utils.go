package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/404answernotfound/browsir/config"
)

func LoadLocalShortcuts() map[string]string {
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

func SaveLocalShortcut(shortcut, url string) error {
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

func GetBrowserPath(browserName string) (string, error) {
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
		case "firefox":
			return "firefox", nil
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
			return filepath.Join(
				os.Getenv("ProgramFiles"),
				"BraveSoftware",
				"Brave-Browser",
				"Application",
				"brave.exe",
			), nil
		}
	}
	return "", fmt.Errorf("unsupported browser %s on %s", browserName, runtime.GOOS)
}

func OpenBrowser(browserName string, profile config.Profile, url string) error {
	browserPath, err := GetBrowserPath(browserName)
	if err != nil {
		return err
	}

	var args []string
	if browserName == "firefox" {
		args = []string{"-P", profile.ProfileDir}
	} else {
		args = []string{"--profile-directory=" + profile.ProfileDir}
	}

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

func PrintUsage(profiles []config.Profile, shortcuts map[string]string, localShortcuts map[string]string) {
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

func PrintProfiles(profiles []config.Profile) {
	fmt.Println("\nProfiles:")
	for _, p := range profiles {
		fmt.Printf("  %-12s - %s\n", p.Name, p.Description)
	}
}

func FindSimilarShortcuts(input string, shortcuts, localShortcuts map[string]string) []string {
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

func PromptYesNo(prompt string) bool {
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

func PrintLocalShortcuts(shortcuts map[string]string) {
	for shortcut, url := range shortcuts {
		fmt.Printf("  %-12s -> %s\n", shortcut, url)
	}
}

func Search(browserName string, profile config.Profile, searchEngine string, searchTerm string) error {

	var url string
	if searchTerm != "" {
		switch searchEngine {
		case "google":
			url = fmt.Sprintf("google.com/search?q=%s", searchTerm)
		case "duckduckgo":
			url = fmt.Sprintf("duckduckgo.com//?q=%s", searchTerm)
		case "brave":
			url = fmt.Sprintf("search.brave.com/search?q=%s", searchTerm)
		default:
			url = fmt.Sprintf("google.com/search?q=%s", searchTerm)
		}
		url = "https://" + url
	}

	OpenBrowser(browserName, profile, url)

	return nil
}

func GetFlags(args []string) map[string]string {
	flags := make(map[string]string)
	if args == nil {
		return flags
	}
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") || strings.HasPrefix(arg, "--") {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				flags[parts[0]] = parts[1]
			} else {
				flags[arg] = ""
			}
		}
	}
	return flags
}

func Contains(args []string, value string) bool {
	for _, arg := range args {
		if arg == value {
			return true
		}
	}
	return false
}
