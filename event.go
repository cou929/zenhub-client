package zenhub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Event struct {
	Id               string `json:"id"`
	Actor            Actor  `json:"actor"`
	Type             string `json:"type"`
	RepositoryId     int    `json:"repoId"`
	Organization     string `json:"organization"`
	Repository       string `json:"repository"`
	IssueNumber      int    `json:"issue"`
	CommentId        int    `json:"comment"`
	RecipientId      int    `json:"recipient"`
	CreatedAt        string `json:"createdAt"`
	SrcPipelineName  string `json:"srcPipelineName"`
	DestPipelineName string `json:"destPipelineName"`
}

type Actor struct {
	Id             string         `json:"id"`
	GithubUserInfo GithubUserInfo `json:"github"`
}

type GithubUserInfo struct {
	Id        int    `json:"id"`
	Name      string `json:"username"`
	AvatarUrl string `json:"avatarUrl"`
}

type Events []Event

func (c *Client) GetEvents(page int) (Events, error) {
	req, err := http.NewRequest("GET", c.urlFor("/v1/events", fmt.Sprintf("page=%d", page)).String(), nil)
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

	var data Events
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, err
}
