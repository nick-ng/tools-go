package jira

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"slices"
	"strings"
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

type sprint struct {
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
	limit = 100
}

func request(method string, url string, body io.Reader) (*http.Response, error) {
	client := http.Client{}
	fmt.Println("User:", atlassianUser)

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

func GetJiraBoard(boardId string, sprintAdjustment int) {
	fmt.Println("Board ID:", boardId)
	fmt.Println("Sprint Adjustment:", sprintAdjustment)

	res, err := Get(fmt.Sprintf(
		"%s/rest/agile/1.0/board/%s/sprint/%s/issue",
		jiraUrl,
		boardId,
		"368",
	))

	if err != nil {
		fmt.Printf("error getting jira board: %s\n", err)
		os.Exit(1)
	}

	if res.StatusCode != 200 {
		fmt.Println("Unexpected status code while getting jira board", res.StatusCode)
		os.Exit(1)
	}

	fmt.Println("Content Type", res.Header.Get("Content-Type"))

	resBody, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Println("error when reading body of jira board response", err)
		os.Exit(1)
	}

	debugFilename := "resBody.json"
	debugFilepath := path.Join("debug", debugFilename)

	utils.MkDirIfNotExist("debug")

	err = os.WriteFile(debugFilepath, resBody, 0644)

	if err != nil {
		fmt.Println("cannot write resBody", err)
	}

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
			fmt.Printf("- %s: %s - %s\n", issue.Key, issue.Fields.Summary, issue.Fields.Assignee.DisplayName)
		}

	}

	// fmt.Println("Jira Board:", res)
}

func getJiraSprints() {
	sprints := []sprint{}
	for i := 0; i < limit; i++ {
		// res, err :=
		sprints = append(sprints, sprint{})
	}
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
