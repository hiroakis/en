package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const CIRCLECI_API_V1 = "https://circleci.com/api/v1.1"

func getProjectEnvironment(project Project) (ProjectEnvironment, error) {
	// https://circleci.com/api/v1.1/project/:vcs-type/:username/:project/envvar?circle-token=:token
	resourcePath := fmt.Sprintf("/project/%s/%s/%s/envvar", project.VCS, project.User, project.Repository)
	body, err := callAPI("GET", resourcePath, nil)
	if err != nil {
		return ProjectEnvironment{}, err
	}

	var Envs []Environment
	json.Unmarshal(body, &Envs)

	return ProjectEnvironment{
		Project:      project,
		Environments: Envs,
	}, nil
}

func getProjects() ([]Project, error) {

	var projects []Project

	body, err := callAPI("GET", "/me", nil)
	if err != nil {
		return nil, err
	}

	var projectResp interface{}
	json.Unmarshal(body, &projectResp)
	for k := range projectResp.(map[string]interface{})["projects"].(map[string]interface{}) {
		// https://github.com/username/repository
		// [github.com", "username", "repository"]
		urlComponents := strings.Split(strings.TrimLeft(k, "https://"), "/")
		project := Project{
			VCS:        strings.Split(urlComponents[0], ".")[0],
			User:       urlComponents[1],
			Repository: urlComponents[2],
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func addEnvironment(project Project, environment Environment) error {

	resourcePath := fmt.Sprintf("/project/%s/%s/%s/envvar",
		project.VCS, project.User, project.Repository)
	envVar, _ := json.Marshal(environment)
	_, err := callAPI("POST", resourcePath, bytes.NewBuffer(envVar))
	if err != nil {
		return err
	}
	return nil
}

func updateEnvironment(project Project, environment Environment) error {
	err := addEnvironment(project, environment)
	if err != nil {
		return err
	}
	return nil
}

func deleteEnvironment(project Project, envName string) error {

	resourcePath := fmt.Sprintf("/project/%s/%s/%s/envvar/%s",
		project.VCS, project.User, project.Repository, envName)

	_, err := callAPI("DELETE", resourcePath, nil)
	if err != nil {
		return err
	}
	return nil
}

func callAPI(method, resourcePath string, payload io.Reader) ([]byte, error) {
	var (
		body []byte
		err  error
	)

	endpoint := fmt.Sprintf("%s%s", CIRCLECI_API_V1, resourcePath)

	req, err := http.NewRequest(method, endpoint, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	values := url.Values{}
	values.Add("circle-token", token)
	req.URL.RawQuery = values.Encode()

	if err != nil {
		return nil, err
	}

	client := http.Client{Timeout: time.Duration(10 * time.Second)}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		err = fmt.Errorf("%d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
		return nil, err
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return body, nil
}
