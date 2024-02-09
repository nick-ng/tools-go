package jira

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"time"
	"tools-go/utils"
)

// This function gets a the current sprint from the given boardId.
// If a sprintAdjustment is provided, instead gives the sprint that corresponds
// to the adjustment from the current sprint.
//
// Usage: GetJiraBoard("45", -2)
func GetJiraBoard(boardId string, sprintAdjustment int) {
	sprints := getJiraSprints(boardId)

	activeSprintIndex := 0

	for i, sprint := range sprints {
		if sprint.State == "active" {
			activeSprintIndex = i
			break
		}
	}

	sprintIndex := activeSprintIndex + sprintAdjustment

	if (sprintIndex < 0) || (sprintIndex > len(sprints)-1) {
		fmt.Printf("sprint adjustment out of range. must be between %d and %d\n", -activeSprintIndex, len(sprints)-activeSprintIndex-1)
		os.Exit(1)
	}

	sprint := sprints[sprintIndex]

	now := time.Now()

	startDate, err := time.Parse(time.RFC3339, sprint.StartDateString)
	includeStartEndDate := true

	if err != nil {
		includeStartEndDate = false
	}

	endDate, err := time.Parse(time.RFC3339, sprint.EndDateString)

	if err != nil {
		includeStartEndDate = false
	}

	fmt.Printf("\nToday: %s\n", utils.FormatDate(now))
	fmt.Println("")

	if includeStartEndDate {
		fmt.Printf("Sprint: %s (%s - %s)\n", sprint.Name, utils.FormatDate(startDate), utils.FormatDate(endDate))
	} else {
		fmt.Printf("Sprint: %s\n", sprint.Name)
	}

	if len(sprint.Goal) > 0 {
		fmt.Printf("Goal: %s\n", sprint.Goal)
	}

	if includeStartEndDate {
		diff := endDate.Sub(now)
		diffHours := diff.Hours()
		fmt.Printf("Days left: %0.1f\n", diffHours/24)
	}

	fmt.Println("")
	res, err := Get(fmt.Sprintf(
		"%s/rest/agile/1.0/board/%s/sprint/%d/issue",
		jiraUrl,
		boardId,
		sprint.Id,
	))

	if err != nil {
		fmt.Printf("error getting jira board: %s\n", err)
		os.Exit(1)
	}

	if res.StatusCode != 200 {
		fmt.Println("Unexpected status code while getting jira board", res.StatusCode)
		os.Exit(1)
	}

	resBody, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Println("error when reading body of jira board response", err)
		os.Exit(1)
	}

	utils.WriteBytesDebug("boardBody.json", resBody)

	var resBodyObj BoardResponseBody

	err = json.Unmarshal(resBody, &resBodyObj)

	if err != nil {
		fmt.Println("cannot unmarshal board resBody", err)
		os.Exit(1)
	}

	issuesByStatusName := map[string][]BoardIssue{}
	for _, issue := range resBodyObj.Issues {
		issueStatus := issue.Fields.Status.Name

		issuesByStatusName[issueStatus] = append(issuesByStatusName[issueStatus], issue)
	}

	issueStatuses := []string{}
	for status := range issuesByStatusName {
		issueStatuses = append(issueStatuses, status)
	}

	slices.SortFunc(issueStatuses, func(statusA, statusB string) int {
		a := getIssueSortValue(statusA)
		b := getIssueSortValue(statusB)

		return a - b
	})

	for _, statusName := range issueStatuses {
		issues := issuesByStatusName[statusName]
		fmt.Println(getIssueStatusWithColour(statusName))
		for _, issue := range issues {
			if issue.Fields.Assignee.DisplayName != "" {
				fmt.Printf("- %s: %s - %s\n", issue.Key, issue.Fields.Summary, issue.Fields.Assignee.DisplayName)

			} else {
				fmt.Printf("- %s: %s\n", issue.Key, issue.Fields.Summary)
			}
		}

	}
}

// @todo(nick-ng): cache sprints?
func getJiraSprints(boardId string) []Sprint {
	sprints := []Sprint{}

	isLast := false
	urlPattern := "%s/rest/agile/1.0/board/%s/sprint?startAt=%d&maxResults=%d"

	start := 0
	for !isLast {
		url := fmt.Sprintf(urlPattern, jiraUrl, boardId, start, limit)

		res, err := Get(url)

		if err != nil {
			fmt.Printf("error getting jira sprints: %s\n", err)
			os.Exit(1)
		}

		if res.StatusCode != 200 {
			fmt.Println(res)
			fmt.Println("unsuccessful jira sprint request")
			os.Exit(1)
		}

		resBody, err := io.ReadAll(res.Body)

		utils.WriteBytesDebug(fmt.Sprintf("sprintsBody-%d.json", start), resBody)

		if err != nil {
			fmt.Printf("error reading jira sprint body: %s\n", err)
			os.Exit(1)
		}

		var resBodyObj SprintResponseBody

		err = json.Unmarshal(resBody, &resBodyObj)

		if err != nil {
			fmt.Printf("couldn't parse sprints response: %s\n", err)
			os.Exit(1)
		}

		isLast = resBodyObj.IsLast

		start += len(resBodyObj.Values)

		sprints = append(sprints, resBodyObj.Values...)
	}

	return sprints
}

func GetJiraBoardOptionStructExample(options BoardOptions) {
	fmt.Println("Board ID:", options.BoardId)
	fmt.Println("Sprint Adjustment:", options.SprintAdjustment)
}

func getIssueSortValue(issueStatus string) int {
	lowercaseIssue := strings.ToLower(issueStatus)
	switch lowercaseIssue {
	case "in progress":
		{
			return 10
		}
	case "review":
		{
			return 20
		}
	case "done":
		{
			return 30
		}
	default:
		{
			return 9999
		}
	}
}
