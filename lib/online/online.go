package online

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	userAgent   = "retrotxt version ping"
	httpTimeout = time.Second * 3
)

// API interface to store the JSON results from GitHub
type API map[string]interface{}

// ReleaseAPI GitHub API v3 releases endpoint
// See: https://developer.github.com/v3/repos/releases/
const ReleaseAPI = "https://api.github.com/repos/bengarrett/retrotxtgo/releases/latest"

// Endpoint request an API endpoint from the URL.
// A HTTP ETag can be provided to validate local data cache against the server.
// The useCache will return true with the etag value matches the server's ETag header.
func Endpoint(url, etag string) (useCache bool, data API, err error) {
	useCache = false
	resp, body, err := Get(url, etag)
	if err != nil {
		return useCache, data, fmt.Errorf("endpoint get failed: %s", err)
	}
	defer resp.Body.Close()
	if etag != "" {
		s := resp.StatusCode
		if s == 304 || (s == 200 && body == nil) {
			// Not Modified
			return true, data, nil
		}
	}
	if ok := json.Valid(body); !ok {
		return useCache, data, fmt.Errorf("the endpoint response is not valid json: %s", url)
	}
	if err = json.Unmarshal(body, &data); err != nil {
		return useCache, data, fmt.Errorf("could not unmarshal the endpoint: %s", url)
	}
	data["etag"] = resp.Header.Get("Etag")
	return useCache, data, nil
}

// Get fetches a URL and returns both its response and body.
// If an etag is provided a "If-None-Match" header request will be included.
func Get(url, etag string) (resp *http.Response, body []byte, err error) {
	client := &http.Client{
		Timeout: httpTimeout,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, body, fmt.Errorf("getting a new request error: %s", err)
	}
	if etag != "" {
		req.Header.Set("If-None-Match", etag)
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err = client.Do(req)
	if err != nil {
		return nil, body, fmt.Errorf("requesting to set the get user-agent header: %s", err)
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, body, fmt.Errorf("reading the response body failed: %s", err)
	}
	return resp, body, resp.Body.Close()
}

// Ping requests a URL and determines if the status is ok.
func Ping(url string) (ok bool, err error) {
	ok = false
	client := &http.Client{
		Timeout: httpTimeout,
	}
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return ok, fmt.Errorf("pinging a new request error: %s", err)
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return ok, fmt.Errorf("requesting to set the ping user-agent header: %s", err)
	}
	return (resp.StatusCode >= 200 && resp.StatusCode <= 299), resp.Body.Close()
}
