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
	Fields IssueFields `json:"fields"`
}

type IssueFields struct {
	Description IssueDescription `json:"description"`
}

type IssueDescription struct {
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

	utils.WriteBytesDebug(fmt.Sprintf("issueBody-%s.json", issueId), resBody)

	var issue Issue

	err = json.Unmarshal(resBody, &issue)

	if err != nil {
		fmt.Println("cannot unmarshal issue resBody", err)
	}

	output := ContentToString(issue.Fields.Description.Content, "")

	fmt.Println(output)
}
