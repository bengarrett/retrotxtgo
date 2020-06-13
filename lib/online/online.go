package online

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

const userAgent = "retrotxt version ping"
const httpTimeout = time.Second * 3

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
		return useCache, data, err
	}
	if etag != "" {
		s := resp.StatusCode
		if s == 304 || (s == 200 && data == nil) {
			// Not Modified
			return true, data, nil
		}
	}
	if ok := json.Valid(body); !ok {
		return useCache, data, errors.New("the response from is not in json syntax: " + url)
	}
	err = json.Unmarshal(body, &data)
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
		return nil, body, err
	}
	if etag != "" {
		req.Header.Set("If-None-Match", etag)
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err = client.Do(req)
	if err != nil {
		return nil, body, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	return resp, body, err
}

// Ping requests a URL and determines if the status is ok.
func Ping(url string) (ok bool, err error) {
	ok = false
	client := &http.Client{
		Timeout: httpTimeout,
	}
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return ok, err
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return ok, err
	}
	defer resp.Body.Close()
	return (resp.StatusCode >= 200 && resp.StatusCode <= 299), nil
}
