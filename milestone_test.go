package zenhub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMilestones(t *testing.T) {
	org := "testorg"
	repo := "testrepo"

	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		expected := fmt.Sprintf("/v2/%s/%s/milestones", org, repo)
		if req.URL.Path != expected {
			t.Error(fmt.Sprintf("request url path should be %s but %s", expected, req.URL.Path))
		}

		if req.Method != "GET" {
			t.Error(fmt.Sprintf("request method should be %s but %s", "GET", req.Method))
		}

		respJson, _ := json.Marshal([]map[string]interface{}{
			map[string]interface{}{
				"_id":        "blahblahblah",
				"__v":        0,
				"repo_id":    9876,
				"start_date": "2015-01-01T00:00:00.000Z",
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJson))
	}))
	defer server.Close()

	client, _ := NewClientWithOptions("test-token", org, repo, server.URL, false)

	milestones, err := client.GetMilestones()

	if err != nil {
		t.Error("GetEstiamtes error should be nil but", err)
	}

	if len(milestones) != 1 {
		t.Error("number of milestones should be 1 but", len(milestones))
	}

	if milestones[0].RepositoryId != 9876 {
		t.Error("RepositoryId of first element of milestones should be 9876 but", milestones[0].RepositoryId)
	}
}
