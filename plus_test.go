package zenhub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPluses(t *testing.T) {
	org := "testorg"
	repo := "testrepo"
	issue := 123
	repoId := 456

	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		expected_path := "/v1/pluses"
		if req.URL.Path != expected_path {
			t.Error(fmt.Sprintf("request url path should be %s but %s", expected_path, req.URL.Path))
		}

		expected_query := fmt.Sprintf("issue=%d&repoId=%d", issue, repoId)
		if req.URL.RawQuery != expected_query {
			t.Error(fmt.Sprintf("request url query should be %s but %s", expected_query, req.URL.RawQuery))
		}

		if req.Method != "GET" {
			t.Error(fmt.Sprintf("request method should be %s but %s", "GET", req.Method))
		}

		respJson, _ := json.Marshal([]map[string]interface{}{
			map[string]interface{}{
				"id":           "blahblahblah",
				"repository":   repo,
				"organization": org,
				"repoId":       repoId,
				"issue":        issue,
				"comment":      5678,
				"receiver":     9999,
				"createdAt":    "2015-01-01T00:00:00.000Z",
				"user": map[string]interface{}{
					"id": "blahblahblah",
					"github": map[string]interface{}{
						"id":        9877,
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

	pluses, err := client.GetPluses(issue, repoId)

	if err != nil {
		t.Error("GetPluses error should be nil but", err)
	}

	if len(pluses) != 1 {
		t.Error("number of pluses should be 1 but", len(pluses))
	}

	if pluses[0].RepositoryId != repoId {
		t.Error(fmt.Sprintf("Name of first element of pluses should be %d but %d", repoId, pluses[0].RepositoryId))
	}
}
