package main

import (
	"bufio"
	"fmt"
	cnf "github.com/404answernotfound/browsir/config"
	"github.com/404answernotfound/browsir/utils"
	"os"
	"strings"
)

func main() {
	config := cnf.LoadConfig()
	localShortcuts := utils.LoadLocalShortcuts()

	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		utils.PrintUsage(config.Profiles, config.Shortcuts, localShortcuts)
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

	var selectedProfile cnf.Profile
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
		utils.PrintUsage(config.Profiles, config.Shortcuts, localShortcuts)
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
				similar := utils.FindSimilarShortcuts(url, config.Shortcuts, localShortcuts)
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
				if utils.PromptYesNo("Would you like to save this as a shortcut?") {
					fmt.Print("Enter the website URL: ")
					reader := bufio.NewReader(os.Stdin)
					websiteURL, _ := reader.ReadString('\n')
					websiteURL = strings.TrimSpace(websiteURL)

					if err := utils.SaveLocalShortcut(url, websiteURL); err != nil {
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

	if err := utils.OpenBrowser(config.BrowserName, selectedProfile, url); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
