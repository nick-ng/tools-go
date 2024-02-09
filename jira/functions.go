package jira

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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
