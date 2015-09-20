package zenhub

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Pipeline struct {
	Id   string `json:"_id"`
	Name string `json:"name"`
	IsIn bool   `json:"isIn"`
}

type Pipelines []Pipeline

func (c *Client) GetPipelines(issueNumber int) (Pipelines, error) {
	req, err := http.NewRequest("GET", c.urlFor(fmt.Sprintf("/v2/%s/%s/issues/%d/pipelines", c.Organization, c.Repository, issueNumber)).String(), nil)
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

	var data Pipelines
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, err
}

func (c *Client) UpdatePipeline(issueNumber int, pipelineName string) (Pipelines, error) {
	pipelines, err := c.GetPipelines(issueNumber)
	if err != nil {
		return nil, err
	}

	var pipelineId string
	for _, p := range pipelines {
		if p.Name == pipelineName {
			pipelineId = p.Id
		}
	}
	if pipelineId == "" {
		return nil, errors.New("no such pipeline")
	}

	return c.updatePipelineWithPipelineId(issueNumber, pipelineId)
}

func (c *Client) updatePipelineWithPipelineId(issueNumber int, pipelineId string) (Pipelines, error) {
	req, err := http.NewRequest("PUT", c.urlFor(fmt.Sprintf("/v1/github/%s/%s/issues/%d/pipelines", c.Organization, c.Repository, issueNumber)).String(), strings.NewReader(fmt.Sprintf("pipeline_id=%s", pipelineId)))
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

	var data Pipelines
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, err
}
