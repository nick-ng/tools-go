package jira

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"
	"tools-go/utils"
)

type BoardOptions struct {
	BoardId          string
	SprintAdjustment int
}

type BoardResponseBody struct {
	StartAt    int          `json:"startAt"`
	MaxResults int          `json:"maxResults"`
	Total      int          `json:"total"`
	Issues     []BoardIssue `json:"issues"`
}

type SprintResponseBody struct {
	MaxResults int      `json:"maxResults"`
	StartAt    int      `json:"startAt"`
	Total      int      `json:"total"`
	IsLast     bool     `json:"isLast"`
	Values     []Sprint `json:"values"`
}

type Sprint struct {
	Id                 int    `json:"id"`
	Self               string `json:"self"`
	State              string `json:"state"`
	Name               string `json:"name"`
	StartDateString    string `json:"startDate"`
	EndDateString      string `json:"endDate"`
	CompleteDateString string `json:"completeDate"`
	OriginBoardId      int    `json:"originBoardId"`
	Goal               string `json:"goal"`
}

type BoardIssue struct {
	Key    string           `json:"key"`
	Fields BoardIssueFields `json:"fields"`
}

type BoardIssueFields struct {
	IssueType BasicField `json:"issueType"`
	Status    BasicField `json:"status"`
	// Description string     `json: "description"` // It's in a weird format
	Summary  string    `json:"summary"`
	Assignee UserField `json:"assignee"`
}

type UserField struct {
	DisplayName string `json:"displayName"`
}

type BasicField struct {
	Name string `json:"name"`
}

var (
	jiraUrl           string
	atlassianUser     string
	atlassianApiToken string
	limit             int
)

func init() {
	jiraUrl = os.Getenv("JIRA_URL")
	atlassianUser = os.Getenv("ATLASSIAN_USER")
	atlassianApiToken = os.Getenv("ATLASSIAN_API_TOKEN")
	limit = 50
}

func request(method string, url string, body io.Reader) (*http.Response, error) {
	client := http.Client{}

	req, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, err
	}

	b64 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", atlassianUser, atlassianApiToken)))

	req.Header = http.Header{
		"Authorization": {fmt.Sprintf("Basic %s", b64)},
	}

	return client.Do(req)
}

func Get(url string) (*http.Response, error) {
	return request("GET", url, nil)
}

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
		fmt.Printf("sprint adjustment out of range. must be between %d and %d", -activeSprintIndex, len(sprints)-activeSprintIndex-1)
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
		fmt.Println("cannot unmarshal resBody", err)
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

func getIssueStatusWithColour(issueStatus string) string {
	lowercaseIssue := strings.ToLower(issueStatus)
	switch lowercaseIssue {
	case "in progress":
		{
			withColour := fmt.Sprintf("\x1b[1m\x1b[34m%s\x1b[0m", issueStatus)
			return withColour
		}
	case "review":
		{
			withColour := fmt.Sprintf("\x1b[1m\x1b[33m%s\x1b[0m", issueStatus)
			return withColour
		}
	case "done":
		{
			withColour := fmt.Sprintf("\x1b[32m%s\x1b[0m", issueStatus)
			return withColour
		}
	case "blocked":
		{
			withColour := fmt.Sprintf("\x1b[31m%s\x1b[0m", issueStatus)
			return withColour
		}
	default:
		{
			withColour := fmt.Sprintf("\x1b[1m\x1b[90m%s\x1b[0m", issueStatus)
			return withColour
		}
	}
}
