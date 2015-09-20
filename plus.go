package zenhub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Plus struct {
	Id           string `json:"id"`
	Repository   string `json:"repository"`
	Organization string `json:"organization"`
	RepositoryId int    `json:"repoId"`
	IssueNumber  int    `json:"issue"`
	Comment      int    `json:"comment"`
	User         User   `json:"user"`
	Receiver     int    `json"receiver"`
	CreatedAt    string `json:"createdAt"`
}

type User struct {
	Id             string         `json:"id"`
	GithubUserInfo GithubUserInfo `json:"github"`
}

type Pluses []Plus

func (c *Client) GetPluses(issueNumber, repositoryId int) (Pluses, error) {
	req, err := http.NewRequest("GET", c.urlFor("/v1/pluses", fmt.Sprintf("issue=%d&repoId=%d", issueNumber, repositoryId)).String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data Pluses
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, err
}
