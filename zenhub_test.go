package zenhub

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Header.Get("x-authentication-token") != "test-token" {
			t.Error("x-authentication-token header should contain passed token")
		}
	}))
	defer server.Close()

	client, err := NewClientWithOptions("test-token", "org", "repo", server.URL, false)
	if err != nil {
		t.Error("error should be nil but", err)
	}

	req, _ := http.NewRequest("GET", client.urlFor("/").String(), nil)

	client.Request(req)
}

func TestBuildReq(t *testing.T) {
	client := NewClient("test-token", "org", "repo")
	client.UserAgent = "test-ua"

	req, _ := http.NewRequest("GET", client.urlFor("/").String(), nil)
	req = client.buildReq(req)

	if req.Header.Get("x-authentication-token") != "test-token" {
		t.Error("x-authentication-token header should contain passed token")
	}

	if req.Header.Get("User-Agent") != "test-ua" {
		t.Error("User-Agent header should contain passed user agent string")
	}

	req, _ = http.NewRequest("PUT", client.urlFor("/").String(), strings.NewReader("test=1"))
	req = client.buildReq(req)

	if req.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		t.Error("Content-Type header should x-www-form-urlencoded when method is POST/PUT")
	}
}

func TestUrlFor(t *testing.T) {
	client := NewClient("test-token", "org", "repo")

	path := "test/path"
	query := "testquery=1&foo=bar"
	url := client.urlFor(path, query)

	expected := fmt.Sprintf("%s%s?%s", defaultBaseURL, path, query)

	if url.String() != expected {
		t.Error("url should be built with passed path and query")
	}
}
