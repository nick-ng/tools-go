package jira

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"tools-go/utils"
)

type Issue struct {
	Key    string      `json:"key"`
	Fields IssueFields `json:"fields"`
}

type IssueFields struct {
	Description IssueDescription `json:"description"`
	Assignee    Assignee         `json:"assignee"`
	Status      Status           `json:"status"`
	Comment     CommentField     `json:"comment"`
}

type IssueDescription struct {
	Content []Content `json:"content"`
}

type Assignee struct {
	DisplayName string `json:"displayName"`
}

type Status struct {
	Name string `json:"name"`
}

type CommentField struct {
	Comments []Comment `json:"comments"`
	Total    int       `json:"total"`
}

type Comment struct {
	Author CommentAuthor `json:"author"`
	Body   CommentBody   `json:"body"`
}

type CommentAuthor struct {
	DisplayName string `json:"displayName"`
}

type CommentBody struct {
	Content []Content `json:"content"`
}

func GetJiraIssue(issueId string) {
	urlPattern := "%s/rest/api/3/issue/%s"

	url := fmt.Sprintf(urlPattern, jiraUrl, strings.ToUpper(issueId))

	res, err := Get(url)

	if err != nil {
		fmt.Printf("error getting jira issue: %s\n", err)
		os.Exit(1)
	}

	resBody, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Printf("error reading jira issue body: %s\n", err)
		os.Exit(1)
	}

	filename := fmt.Sprintf("issueBody-%s.json", strings.ToUpper(issueId))

	utils.WriteBytesDebug(filename, resBody)

	var issue Issue

	err = json.Unmarshal(resBody, &issue)

	if err != nil {
		fmt.Println("cannot unmarshal issue resBody", err)
	}

	content := ContentToString(issue.Fields.Description.Content, "")

	fmt.Println(content)
}
