package main

import (
	"fmt"
	"os"
	"strconv"
	"tools-go/jira"
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

	switch flags.resource {
	case "board":
		{
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
	case "issue":
		{
			jira.GetJiraIssue(flags.id)
		}
	}
}

func readFlags(args []string) flagOptions {
	if len(args) > 1 {
		switch args[1] {
		case "board":
			fallthrough
		case "b":
			{
				return readBoardFlags(args)
			}
		case "issue":
			fallthrough
		case "i":
			{
				return readIssueFlags(args)
			}
		default:
			{
				// noop
			}
		}

	}

	fmt.Println("invalid flags")
	fmt.Println("usage: jira resource-type resource-id [action]")
	fmt.Println("instead, got:")
	fmt.Println(args[1:])
	os.Exit(1)

	return flagOptions{}
}

/*
Example usage:
- jira board 45
- jira b 45
- jira b 45 1
- jira b 45 +1
- jira b 45 -1
*/
func readBoardFlags(args []string) flagOptions {
	if len(args) < 3 {
		// 0    1     2
		// jira board 45
		return flagOptions{}
	}

	options := flagOptions{
		resource:  "board",
		id:        args[2],
		action:    "get",
		subAction: "0",
	}

	if len(args) >= 4 {
		// Sprint adjustment: -999 - 999/+999
		options.subAction = args[3]
	}

	return options
}

// @todo(nick-ng): have some way of storing a "current issue" so you don't have to remember the issue id
func readIssueFlags(args []string) flagOptions {
	if len(args) < 3 {
		// 0    1     2
		// jira issue PLAT-100
		return flagOptions{}
	}

	options := flagOptions{
		resource:  "issue",
		id:        args[2],
		action:    "get",
		subAction: "",
	}

	if len(args) >= 4 {
		switch args[3] {
		case "c":
			fallthrough
		case "comment":
			{
				options.subAction = "comment"
			}
		}
	}

	fmt.Println(options)

	return options
}
