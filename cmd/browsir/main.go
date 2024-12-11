package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	cnf "github.com/404answernotfound/browsir/config"
	"github.com/404answernotfound/browsir/utils"
)

func main() {
	config := cnf.LoadConfig()
	localShortcuts := utils.LoadLocalShortcuts()

	flags := utils.GetFlags(os.Args[1:])

	if flags["-h"] != "" || flags["--help"] != "" {
		utils.PrintUsage(config.Profiles, config.Shortcuts, localShortcuts)
		os.Exit(0)
	}

	if flags["-v"] != "" || flags["--version"] != "" {
		if version := os.Getenv("BROWSIR_VERSION"); version != "" {
			fmt.Println("browsir v" + version)
		} else {
			fmt.Println("browsir version not set")
		}
		os.Exit(0)
	}

	if flags["-ls"] != "" || flags["--list-shortcuts"] != "" {
		utils.PrintLocalShortcuts(config.Shortcuts)
		os.Exit(0)
	}

	var profileName string
	var url string
	args := os.Args[1:]

	var selectedProfile cnf.Profile
	var found bool
	profileName = args[0]

	for _, p := range config.Profiles {
		if p.Name == profileName {
			selectedProfile = p
			found = true
			break
		}
	}

	if flags["-q"] != "" {
		fmt.Println("Searching...")
		cleanQuery := strings.ReplaceAll(flags["-q"], " ", "+")
		cleanQuery = strings.ReplaceAll(cleanQuery, "\"", "")

		if flags["-se"] != "" || flags["--search-engine"] != "" {
			searchEngine := flags["-se"]
			if searchEngine == "" {
				searchEngine = flags["--search-engine"]
			}
			utils.Search(config.BrowserName, selectedProfile, searchEngine, cleanQuery)
			os.Exit(0)
		}

		searchEngine := "google"
		utils.Search(config.BrowserName, selectedProfile, searchEngine, cleanQuery)

		os.Exit(0)
	}

	if len(args) > 1 {
		url = args[1]
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
