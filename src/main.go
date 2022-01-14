package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	configPath    = "./etc"
	configConnect = "secret.config.json"

	urlRest             = "/rest/api/2"
	urlMethodCreatemeta = "/issue/createmeta"
	urlMethodSearch     = "/search"
	urlMethodIssue      = "/issue/"
	urlMethodWorklog    = "/worklog"
)

// CONFIG

type ConfigT struct {
	Hostname string         `json:"hostname,omitempty"`
	Username string         `json:"username,omitempty"`
	Password string         `json:"password,omitempty"`
	JQL      string         `json:"jql,omitempty"`
	Worklog  ConfigWorklogT `json:"worklog,omitempty"`
}

type ConfigWorklogT struct {
	Author string `json:"author,omitempty"`
	Begin  string `json:"begin,omitempty"`
	End    string `json:"end,omitempty"`
}

// JIRA

type JiraSearchIssueT struct {
	MaxResults int           `json:"maxResults,omitempty"`
	Total      int           `json:"total,omitempty"`
	Issues     []JiraIssuesT `json:"issues,omitempty"`
}

type JiraIssuesT struct {
	Key string `json:"key,omitempty"`
}

type JiraIssueWorklogT struct {
	MaxResults int            `json:"maxResults,omitempty"`
	Total      int            `json:"total,omitempty"`
	Worklogs   []JiraWorklogT `json:"worklogs,omitempty"`
}

type JiraWorklogT struct {
	Started          string      `json:"started,omitempty"`
	TimeSpentSeconds int         `json:"timeSpentSeconds,omitempty"`
	Author           JiraAuthorT `json:"author,omitempty"`
}

type JiraAuthorT struct {
	Name string `json:"name,omitempty"`
}

// MAIN ------------------------------------------------------------------------

func main() {
	config := getconfigConnect(configConnect)

	issues := jiraSearchIssue(config)

	var totalWorklog int
	for _, issue := range issues {
		totalWorklog += jiraIssueGetWorklog(config, issue.Key)
	}

	// TODO: Convert seconds to hours
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
	// TODO: generate JQL from config's values
	url := urlMethodSearch + "?jql=" + config.JQL
	response := jiraRestGet(config, url)

	var found JiraSearchIssueT
	err := json.Unmarshal(response, &found)
	if err != nil {
		log.Fatalf("Error in unmarshal '%s'", err.Error())
	}

	// TODO: merge duplicate comparison "Total" and "MaxResults"
	// Max Result by default is 50 items
	if found.Total > found.MaxResults {
		log.Fatalf("Total items '%d' more max result '%d'", found.Total, found.MaxResults)
	}

	return found.Issues
}

// JIRA ISSUE WORKLOG

func jiraIssueGetWorklog(config ConfigT, issue string) int {
	var totalWorklog int

	worklogs := jiraGetWorklogs(config, issue)
	for _, worklog := range worklogs {
		// TODO: filter by Author and search period
		if worklog.Author.Name != config.Worklog.Author {
			continue
		}

		// Use field 'started', but not 'created', exemple:
		//* '2021-01-26T07:20:00.000+0400'
		started := strings.SplitN(worklog.Started, "T", 2)[0]
		if started < config.Worklog.Begin || config.Worklog.End < started {
			continue
		}

		totalWorklog += worklog.TimeSpentSeconds
	}

	return totalWorklog
}

func jiraGetWorklogs(config ConfigT, issue string) []JiraWorklogT {
	url := urlMethodIssue + issue + urlMethodWorklog
	response := jiraRestGet(config, url)

	var found JiraIssueWorklogT
	err := json.Unmarshal(response, &found)
	if err != nil {
		log.Fatalf("Error in unmarshal '%s'", err.Error())
	}

	// TODO: merge duplicate comparison "Total" and "MaxResults"
	// Max Result by default is 50 items
	if found.Total > found.MaxResults {
		log.Fatalf("Total items '%d' more max result '%d'", found.Total, found.MaxResults)
	}

	return found.Worklogs
}

// JIRA REST -------------------------------------------------------------------

// TODO: 'config' is static data for all REST methods
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
