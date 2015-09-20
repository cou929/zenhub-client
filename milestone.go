package zenhub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Milestone struct {
	Id           string `json:"_id"`
	V            int    `json:"__v"`
	RepositoryId int    `json:"repo_id"`
	StartDate    string `json:"start_date"`
}

type Milestones []Milestone

func (c *Client) GetMilestones() (Milestones, error) {
	req, err := http.NewRequest("GET", c.urlFor(fmt.Sprintf("/v2/%s/%s/milestones", c.Organization, c.Repository)).String(), nil)
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

	var data Milestones
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, err
}
