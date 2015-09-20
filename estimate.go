package zenhub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Estimation struct {
	Value    int  `json:"value"`
	Selected bool `json:"selected"`
}

type Estimates []Estimation

type UpdateResponse struct {
	Id           string `json:"_id"`
	V            int    `json:"__v"`
	Value        int    `json:"estimate"`
	IssueNumber  int    `json:"issue_number"`
	LastModified string `json:"last_modified"`
	RepositoryId int    `json:"repo_id"`
}

func (c *Client) GetEstimates(issueNumber int) (Estimates, error) {
	req, err := http.NewRequest("GET", c.urlFor(fmt.Sprintf("/v2/%s/%s/issues/%d/estimates", c.Organization, c.Repository, issueNumber)).String(), nil)
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

	var data Estimates
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (c *Client) UpdateEstimate(issueNumber, estimateValue int) (*UpdateResponse, error) {
	req, err := http.NewRequest("PUT", c.urlFor(fmt.Sprintf("/v2/%s/%s/issues/%d/estimates", c.Organization, c.Repository, issueNumber)).String(), strings.NewReader(fmt.Sprintf("estimate_value=%d", estimateValue)))
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

	var data UpdateResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, err
}

func (c *Client) ClearEstimate(issueNumber int) error {
	req, err := http.NewRequest("DELETE", c.urlFor(fmt.Sprintf("/v2/%s/%s/issues/%d/estimates", c.Organization, c.Repository, issueNumber)).String(), nil)
	if err != nil {
		return err
	}

	resp, err := c.Request(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
