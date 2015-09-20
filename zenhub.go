package zenhub

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

const (
	defaultBaseURL    = "https://api.zenhub.io/"
	defaultUserAgent  = "zenhub-client"
	apiRequestTimeout = 30 * time.Second
)

type Client struct {
	BaseURL      *url.URL
	AuthToken    string
	Verbose      bool
	UserAgent    string
	Organization string
	Repository   string
}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func NewClient(authToken, org, repo string) *Client {
	u, _ := url.Parse(defaultBaseURL)
	return &Client{
		BaseURL:      u,
		AuthToken:    authToken,
		Verbose:      false,
		UserAgent:    defaultUserAgent,
		Organization: org,
		Repository:   repo,
	}
}

func NewClientWithOptions(authToken string, org string, repo string, baseurl string, verbose bool) (*Client, error) {
	u, err := url.Parse(baseurl)
	if err != nil {
		return nil, err
	}
	return &Client{
		BaseURL:      u,
		AuthToken:    authToken,
		Verbose:      verbose,
		UserAgent:    defaultUserAgent,
		Organization: org,
		Repository:   repo,
	}, nil
}

func (c *Client) urlFor(params ...string) *url.URL {
	newUrl, err := url.Parse(c.BaseURL.String())
	if err != nil {
		panic("invalid BaseURL passed")
	}

	if len(params) <= 0 {
		panic("too few arguments")
	}

	newUrl.Path = params[0]
	if len(params) >= 2 && params[1] != "" {
		newUrl.RawQuery = params[1]
	}
	return newUrl
}

func (c *Client) buildReq(req *http.Request) *http.Request {
	req.Header.Set("x-authentication-token", c.AuthToken)
	req.Header.Set("User-Agent", c.UserAgent)
	if req.Method == "PUT" || req.Method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	return req
}

func (c *Client) Request(req *http.Request) (resp *http.Response, err error) {
	req = c.buildReq(req)

	if c.Verbose {
		dump, err := httputil.DumpRequest(req, true)
		if err == nil {
			log.Printf("%s", dump)
		}
	}

	client := &http.Client{}
	client.Timeout = apiRequestTimeout
	resp, err = client.Do(req)

	if err != nil {
		return nil, err
	}
	if c.Verbose {
		dump, err := httputil.DumpResponse(resp, true)
		if err == nil {
			log.Printf("%s", dump)
		}
	}

	if resp.StatusCode/100 != 2 {
		return nil, errors.New(fmt.Sprintf("request failed. status code is %d", resp.StatusCode))
	}

	return resp, nil
}
