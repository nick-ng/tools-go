package jira

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
)

type BoardOptions struct {
	BoardId          string
	SprintAdjustment int
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

	fmt.Println("resBody", string(resBody))

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
