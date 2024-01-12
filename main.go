package main

import (
	"fmt"
	"os"
	"strconv"
	"tools-go/jira"
)

const (
	issue = iota
	board = iota
)

type flagOptions struct {
	resource  string
	id        string
	action    string
	subAction string
}

func main() {
	flags := readFlags(os.Args)

	// fmt.Println(os.Args)
	// fmt.Println(os.Args[1])
	// fmt.Println(flags)

	if flags.resource == "board" {
		sprintAdjustment, err := strconv.Atoi(flags.subAction)

		if err != nil {
			sprintAdjustment = 0
		}

		switch flags.action {
		default:
			{ // Get
				jira.GetJiraBoard(flags.id, sprintAdjustment)
			}
		}
	}
}

func readFlags(args []string) flagOptions {
	if len(args) > 1 {
		switch args[1] {
		case "board":
		case "b":
			{
				return readBoardFlags(args)
			}
		case "issue":
		case "i":
			{
				return readIssueFlags(args)
			}
		}
	}

	fmt.Println("invalid flags")
	fmt.Println("usage: jira resource-type resource-id [action]")
	fmt.Println("instead, got:")
	fmt.Println(args)
	os.Exit(1)

	return flagOptions{}
}

func readBoardFlags(args []string) flagOptions {
	// fmt.Println(args)
	if len(args) < 3 {
		return flagOptions{}
	}

	option := flagOptions{
		resource:  "board",
		id:        args[2],
		action:    "get",
		subAction: "0",
	}

	if len(args) >= 4 {
		// Sprint adjustment: -999 - 999/+999
		option.subAction = args[3]
	}

	return option
}

func readIssueFlags(args []string) flagOptions {
	return flagOptions{}
}
