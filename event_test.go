package zenhub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetEvents(t *testing.T) {
	org := "testorg"
	repo := "testrepo"
	page := 2

	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/v1/events" {
			t.Error("request url path should be /v1/events but", req.URL.Path)
		}

		expected_query := fmt.Sprintf("page=%d", page)
		if req.URL.RawQuery != expected_query {
			t.Error(fmt.Sprintf("request url query should be %s but %s", expected_query, req.URL.RawQuery))
		}

		if req.Method != "GET" {
			t.Error(fmt.Sprintf("request method should be %s but %s", "GET", req.Method))
		}

		respJson, _ := json.Marshal([]map[string]interface{}{
			map[string]interface{}{
				"id":               "blahblahblah",
				"type":             "transferIssue",
				"repoId":           123456,
				"organization":     org,
				"repository":       repo,
				"issue":            1234,
				"srcPipelineName":  "New Issues",
				"destPipelineName": "Todo",
				"createdAt":        "2015-01-01T00:00:00.000Z",
				"actor": map[string]interface{}{
					"id": "blahblahblah",
					"github": map[string]interface{}{
						"id":        9876,
						"username":  "john",
						"avatarUrl": "https://avatars.githubusercontent.com/u/000?v=3",
					},
				},
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJson))
	}))
	defer server.Close()

	client, _ := NewClientWithOptions("test-token", org, repo, server.URL, false)

	events, err := client.GetEvents(page)

	if err != nil {
		t.Error("GetEvents error should be nil but", err)
	}

	if len(events) != 1 {
		t.Error("number of events should be 1 but", len(events))
	}

	if events[0].Id != "blahblahblah" {
		t.Error("Id of first element of events should be blahblahblah but", events[0].Id)
	}
}
