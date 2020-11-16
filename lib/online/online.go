package online

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	userAgent   = "retrotxt version ping"
	httpTimeout = time.Second * 3
)

var (
	// ErrJSON body is not valid json
	ErrJSON = errors.New("response body data is not valid json")
	// ErrMash body is not able to unmarshal
	ErrMash = errors.New("response body json could not unmarshal")
)

// API interface to store the JSON results from GitHub.
type API map[string]interface{}

// ReleaseAPI GitHub API v3 releases endpoint.
// See: https://developer.github.com/v3/repos/releases/
const ReleaseAPI = "https://api.github.com/repos/bengarrett/retrotxtgo/releases/latest"

// Endpoint request an API endpoint from the URL.
// A HTTP ETag can be provided to validate local data cache against the server.
// The useCache will return true with the etag value matches the server's ETag header.
func Endpoint(url, etag string) (useCache bool, data API, err error) {
	useCache = false
	resp, body, err := Get(url, etag)
	if err != nil {
		return useCache, data, fmt.Errorf("endpoint get failed: %w", err)
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
		return useCache, data, fmt.Errorf("endpoint %s: %w", url, ErrJSON)
	}
	if err = json.Unmarshal(body, &data); err != nil {
		return useCache, data, fmt.Errorf("endpoint %s: %w", url, ErrMash)
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
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, httpTimeout)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	defer cancel()
	if err != nil {
		return nil, body, fmt.Errorf("getting a new request error: %w", err)
	}
	if etag != "" {
		req.Header.Set("If-None-Match", etag)
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err = client.Do(req)
	if err != nil {
		return nil, body, fmt.Errorf("requesting to set the get user-agent header: %w", err)
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, body, fmt.Errorf("reading the response body failed: %w", err)
	}
	return resp, body, resp.Body.Close()
}

// Ping requests a URL and determines if the status is ok.
func Ping(url string) (ok bool, err error) {
	ok = false
	client := &http.Client{
		Timeout: httpTimeout,
	}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, httpTimeout)
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	defer cancel()
	if err != nil {
		return ok, fmt.Errorf("pinging a new request error: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return ok, fmt.Errorf("requesting to set the ping user-agent header: %w", err)
	}
	return (resp.StatusCode >= 200 && resp.StatusCode <= 299), resp.Body.Close()
}
