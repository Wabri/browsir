package browsir

import (
	"fmt"

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
	return fmt.Errorf("preview not implemented")
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
