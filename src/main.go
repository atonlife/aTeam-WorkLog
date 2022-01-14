package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
)

const (
	configPath    = "./etc"
	configConnect = "secret.config.json"

	urlRest             = "/rest/api/2"
	urlMethodCreatemeta = "/issue/createmeta"
	urlMethodSearch     = "/search"
)

type ConfigT struct {
	Hostname string `json:"hostname,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	JQL      string `json:"jql,omitempty"`
}

type JiraSearchIssueT struct {
	MaxResults int           `json:"maxResults,omitempty"`
	Total      int           `json:"total,omitempty"`
	Issues     []JiraIssuesT `json:"issues,omitempty"`
}

type JiraIssuesT struct {
	Key string `json:"key,omitempty"`
}

// MAIN ------------------------------------------------------------------------

func main() {
	config := getconfigConnect(configConnect)

	issues := jiraSearchIssue(config)

	var totalWorklog int
	for _, issue := range issues {
		totalWorklog += jiraIssueGetWorklog(issue.Key)
	}

	fmt.Printf("Total Worklog: %d\n", totalWorklog)
}

// CONFIG ----------------------------------------------------------------------

func getconfigConnect(filename string) ConfigT {
	filesystem := os.DirFS(configPath)
	file, err := fs.ReadFile(filesystem, filename)
	if err != nil {
		log.Fatalf("Error in file: '%s'", err.Error())
	}

	var connect ConfigT
	err = json.Unmarshal(file, &connect)
	if err != nil {
		log.Fatalf("Error in unmarshal '%s'", err.Error())
	}

	return connect
}

// JIRA ------------------------------------------------------------------------

// JIRA ISSUE

func jiraSearchIssue(config ConfigT) []JiraIssuesT {
	url := urlMethodSearch + "?jql=" + config.JQL
	response := jiraRestGet(config, url)

	var search JiraSearchIssueT
	err := json.Unmarshal(response, &search)
	if err != nil {
		log.Fatalf("Error in unmarshal '%s'", err.Error())
	}

	// Max Result by default is 50 items
	if search.Total > search.MaxResults {
		log.Fatalf("Total items '%d' more max result '%d'", search.Total, search.MaxResults)
	}

	return search.Issues
}

// JIRA ISSUE WORKLOG

func jiraIssueGetWorklog(issue string) int {
	//TODO: Get sum issue worklog by period and author (get_issue_worklog)
	var totalWorklog int

	return totalWorklog
}

// JIRA REST -------------------------------------------------------------------

func jiraRestGet(config ConfigT, urlRequest string) []byte {
	url := config.Hostname + urlRest + urlRequest
	log.Printf("Request URL '%s'", url)
	response := sendGet(url, config.Username, config.Password)

	log.Printf("Response Status '%d'", response.StatusCode)
	if response.StatusCode != http.StatusOK {
		log.Fatalf("Jira Response Status isn't OK")
	}

	return readBody(response)
}

// REST ------------------------------------------------------------------------

func sendGet(url, user, password string) *http.Response {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.SetBasicAuth(user, password)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	return response
}

func readBody(response *http.Response) []byte {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	return body
}
