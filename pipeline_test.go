package zenhub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPipelines(t *testing.T) {
	org := "testorg"
	repo := "testrepo"
	issue := 123

	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		expected := fmt.Sprintf("/v2/%s/%s/issues/%d/pipelines", org, repo, issue)
		if req.URL.Path != expected {
			t.Error(fmt.Sprintf("request url path should be %s but %s", expected, req.URL.Path))
		}

		if req.Method != "GET" {
			t.Error(fmt.Sprintf("request method should be %s but %s", "GET", req.Method))
		}

		respJson, _ := json.Marshal([]map[string]interface{}{
			map[string]interface{}{
				"_id":  "abc",
				"name": "New Issue",
			},
			map[string]interface{}{
				"_id":  "efg",
				"name": "Todo",
				"isIn": true,
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJson))
	}))
	defer server.Close()

	client, _ := NewClientWithOptions("test-token", org, repo, server.URL, false)

	pipelines, err := client.GetPipelines(issue)

	if err != nil {
		t.Error("GetPipelines error should be nil but", err)
	}

	if len(pipelines) != 2 {
		t.Error("number of pipelines should be 2 but", len(pipelines))
	}

	if pipelines[0].Name != "New Issue" {
		t.Error("Name of first element of pipelines should be New Issue but", pipelines[0].Name)
	}

	if pipelines[1].Name != "Todo" {
		t.Error("Name of second element of pipelines should be Todo but", pipelines[1].Name)
	}

	if pipelines[0].IsIn != false {
		t.Error("IsIn of first element of pipelines should be false but", pipelines[0].IsIn)
	}

	if pipelines[1].IsIn != true {
		t.Error("IsIn of second element of pipelines should be true but", pipelines[1].IsIn)
	}
}

func TestUpdatePipelineWithPipelineId(t *testing.T) {
	org := "testorg"
	repo := "testrepo"
	issue := 123
	pipelineId := "efg"

	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		expected := fmt.Sprintf("/v1/github/%s/%s/issues/%d/pipelines", org, repo, issue)
		if req.URL.Path != expected {
			t.Error(fmt.Sprintf("request url path should be %s but %s", expected, req.URL.Path))
		}

		if req.Method != "PUT" {
			t.Error(fmt.Sprintf("request method should be %s but %s", "PUT", req.Method))
		}

		body, _ := ioutil.ReadAll(req.Body)
		if string(body) != fmt.Sprintf("pipeline_id=%s", pipelineId) {
			t.Error("request body should be correct but", string(body))
		}

		respJson, _ := json.Marshal([]map[string]interface{}{
			map[string]interface{}{
				"_id":  "abc",
				"name": "New Issue",
			},
			map[string]interface{}{
				"_id":  "efg",
				"name": "Todo",
				"isIn": true,
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJson))
	}))
	defer server.Close()

	client, _ := NewClientWithOptions("test-token", org, repo, server.URL, false)

	pipelines, err := client.updatePipelineWithPipelineId(issue, pipelineId)

	if err != nil {
		t.Error("UpdatePipeline error should be nil but", err)
	}

	if len(pipelines) != 2 {
		t.Error("number of pipelines should be 2 but", len(pipelines))
	}

	if pipelines[1].IsIn != true {
		t.Error("IsIn of second element of pipelines should be true but", pipelines[1].IsIn)
	}
}

func TestUpdatePipeline(t *testing.T) {
	org := "testorg"
	repo := "testrepo"
	issue := 123
	pipelineName := "unexistName"

	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		respJson, _ := json.Marshal([]map[string]interface{}{
			map[string]interface{}{
				"_id":  "abc",
				"name": "New Issue",
			},
			map[string]interface{}{
				"_id":  "efg",
				"name": "Todo",
				"isIn": true,
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJson))
	}))
	defer server.Close()

	client, _ := NewClientWithOptions("test-token", org, repo, server.URL, false)

	_, err := client.UpdatePipeline(issue, pipelineName)

	if err == nil {
		t.Error("UpdatePipeline should be error")
	}

	if err.Error() != "no such pipeline" {
		t.Error("UpdatePipeline should be error due to pipelineName argument is invalid but", err)
	}
}
