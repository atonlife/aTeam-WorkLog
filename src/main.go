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
	cfgPath    = "./etc"
	cfgConnect = "secret.connect.json"

	urlRest             = "/rest/api/2"
	urlMethodCreatemeta = "/issue/createmeta"
)


type CfgConnectT struct {
	Server   string `json:"server"`
	User     string `json:"user"`
	Password string `json:"password"`
}


func main() {
	connect := getCfgConnect(cfgConnect)
	response := getResponse(connect, urlMethodCreatemeta)
	readResponse(response)
}

func getCfgConnect(filename string) CfgConnectT {
	var connect CfgConnectT

	filesystem := os.DirFS(cfgPath)
	connectJSON, err := fs.ReadFile(filesystem, filename)
	if err != nil {
		log.Fatalf("Error in JSON: %s", err.Error())
	}

	json.Unmarshal(connectJSON, &connect)

	return connect
}

func getResponse(connect CfgConnectT, method string) *http.Response {
	url := connect.Server + urlRest + method
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.SetBasicAuth(connect.User, connect.Password)

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	return response
}

func readResponse(response *http.Response) {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	//TODO: DEBUG
	fmt.Println(response.StatusCode == http.StatusOK)
	fmt.Println(body != nil)
}
