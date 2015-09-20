package zenhub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetEstimates(t *testing.T) {
	org := "testorg"
	repo := "testrepo"
	issue := 123

	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		expected := fmt.Sprintf("/v2/%s/%s/issues/%d/estimates", org, repo, issue)
		if req.URL.Path != expected {
			t.Error(fmt.Sprintf("request url path should be %s but %s", expected, req.URL.Path))
		}

		if req.Method != "GET" {
			t.Error(fmt.Sprintf("request method should be %s but %s", "GET", req.Method))
		}

		respJson, _ := json.Marshal([]map[string]interface{}{
			map[string]interface{}{
				"value":    1,
				"selected": true,
			},
			map[string]interface{}{
				"value":    2,
				"selected": false,
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJson))
	}))
	defer server.Close()

	client, _ := NewClientWithOptions("test-token", org, repo, server.URL, false)

	estimates, err := client.GetEstimates(issue)

	if err != nil {
		t.Error("GetEstiamtes error should be nil but", err)
	}

	if estimates[0].Value != 1 {
		t.Error("Value of first element of estimates should be 1 but", estimates[0].Value)
	}

	if estimates[1].Selected {
		t.Error("Selected of second element of estimates should be false")
	}
}

func TestUpdateEstimate(t *testing.T) {
	org := "testorg"
	repo := "testrepo"
	issue := 123
	estimateValue := 2

	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		expected := fmt.Sprintf("/v2/%s/%s/issues/%d/estimates", org, repo, issue)
		if req.URL.Path != expected {
			t.Error(fmt.Sprintf("request url path should be %s but %s", expected, req.URL.Path))
		}

		if req.Method != "PUT" {
			t.Error(fmt.Sprintf("request method should be %s but %s", "PUT", req.Method))
		}

		body, _ := ioutil.ReadAll(req.Body)
		if string(body) != fmt.Sprintf("estimate_value=%d", estimateValue) {
			t.Error("request body should be correct but", string(body))
		}

		respJson, _ := json.Marshal(map[string]interface{}{
			"_id":           "blahblahblah",
			"__v":           0,
			"estimate":      estimateValue,
			"issue_number":  issue,
			"last_modified": "2015-01-01T00:00:00.000Z",
			"repo_id":       12345,
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJson))
	}))
	defer server.Close()

	client, _ := NewClientWithOptions("test-token", org, repo, server.URL, false)

	res, err := client.UpdateEstimate(issue, estimateValue)

	if err != nil {
		t.Error("UpdateEstimate error should be nil but", err)
	}

	if res.Value != estimateValue {
		t.Error(fmt.Sprintf("Value should be %d but %d", estimateValue, res.Value))
	}

	if res.IssueNumber != issue {
		t.Error(fmt.Sprintf("IssueNumber should be %d but %d", issue, res.IssueNumber))
	}
}

func TestClearEstimate(t *testing.T) {
	org := "testorg"
	repo := "testrepo"
	issue := 123

	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		expected := fmt.Sprintf("/v2/%s/%s/issues/%d/estimates", org, repo, issue)
		if req.URL.Path != expected {
			t.Error(fmt.Sprintf("request url path should be %s but %s", expected, req.URL.Path))
		}

		if req.Method != "DELETE" {
			t.Error(fmt.Sprintf("request method should be %s but %s", "DELETE", req.Method))
		}
	}))
	defer server.Close()

	client, _ := NewClientWithOptions("test-token", org, repo, server.URL, false)

	err := client.ClearEstimate(issue)

	if err != nil {
		t.Error("UpdateEstimate error should be nil but", err)
	}
}
