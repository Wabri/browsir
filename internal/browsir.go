package browsir

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/404answernotfound/browsir/utils"
)

type ICommand interface {
	add(args []string)
	rm(args []string)
	last(args []string)
	preview(args []string)
}

type Command struct{}

// add, rm, list, preview
func (c Command) add(args []string) error {
	switch args[0] {
	case "link":
		if len(args) < 2 {
			utils.ExitLog("Please, see --help flag to check usage for add command")
		}
		flags := utils.GetFlags(args[2:])
		link := args[1]

		categories, ok := flags["-c"]
		if !ok {
			return fmt.Errorf("not a good flag")
		}

		err := utils.SaveLink(link, categories)
		if err != nil {
			return err
		}
		return nil
	case "shortcut":
		if len(args) < 3 {
			utils.ExitLog("Please, see --help flag to check usage for add command")
		}
		shortcut := args[1]
		url := args[2]
		err := utils.SaveLocalShortcut(shortcut, url)
		return err
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func (c Command) rm(args []string) error {
	return fmt.Errorf("rm not implemented")
}

func (c Command) list(args []string) error {
	links := utils.LoadLinks()
	for link, categories := range links {
		fmt.Printf("Link: %s - Categories: %s\n", link, categories)
	}
	return nil
}
func (c Command) preview(args []string) error {
	ctx := context.Background()
	deadline := time.Now().Add(3000 * time.Millisecond)
	ctx, cancel := context.WithDeadline(ctx, deadline)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://404answernotfound.eu", nil)
	if err != nil {
		return fmt.Errorf("error creating request: %s", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %s", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %s", err)
	}

	fmt.Println(string(body))
	return nil
}

func RunCommand(mainCmd string, otherArgs []string) error {
	command := &Command{}
	var err error

	switch mainCmd {
	case "add":
		err = command.add(otherArgs)
		return err
	case "rm":
		err = command.rm(otherArgs)
		return err
	case "list":
		err = command.list(otherArgs)
		return err
	case "preview":
		err = command.preview(otherArgs)
		return err
	default:
		return fmt.Errorf("not implemented")
	}
}
