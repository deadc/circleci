package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Project struct {
	Reponame     string `json:"reponame"`
	BuildURL     string `json:"build_url"`
	Branch       string `json:"branch"`
	Username     string `json:"username"`
	VcsRev       string `json:"vcs_revision"`
	Status       string `json:"status"`
	CommiterName string `json:"commiter_name"`
	Subject      string `json:"subject"`
	Failed       bool   `json:"failed"`
	BuildNum     int    `json:"build_num"`

	Builder Builder `json:"user"`
}

type Builder struct {
	Name      string `json:"name"`
	Login     string `json:"login"`
	AvatarUrl string `json:"avatar_url"`
}

func main() {

	var projects []Project

	circle_token := flag.String("circle-token", "", "CircleCI API Token")
	circle_api := flag.String("circle-api", "https://circleci.com/api/v1.1", "CircleCI API Endpoint")
	project_name := flag.String("project", "", "Project name")
	vcs_type := flag.String("vcs", "github", "VCS: github or bitbucket")
	username := flag.String("user", "", "CircleCI User")
	build_num := flag.Int("build-num", 1, "Per project num builds")

	flag.Parse()

	circle_payload := fmt.Sprintf("%s/project/%s/%s/%s?circle-token=%s&limit=%d",
		*circle_api, *vcs_type, *username, *project_name, *circle_token, *build_num)

	req, err := http.NewRequest(http.MethodGet, circle_payload, nil)
	if err != nil {
		fmt.Printf("%#v\n", err)
		return
	}

	req.Header.Set("Accept", "*/*")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%#v\n", err)
		return
	}

	//fmt.Printf("req_url >> %s\n", circle_payload)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%#v\n", err)
		return
	}

	jsonErr := json.Unmarshal(body, &projects)
	if jsonErr != nil {
		fmt.Printf("%#v\n", jsonErr)
		return
	}

	for k := range projects {
		fmt.Printf("[%02d/%s] - Builder >> [%s](%s)\tStatus >> %-20s\tSubject >>%s\n",
			projects[k].BuildNum, projects[k].Branch, projects[k].Builder.Login, projects[k].Builder.Name, projects[k].Status, projects[k].Subject)
	}
}
