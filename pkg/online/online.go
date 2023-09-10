// Package online is for HTTP interactions.
package online

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bengarrett/retrotxtgo/meta"
)

var (
	ErrJSON = errors.New("cannot understand the response body as the syntax is not json")
	ErrMash = errors.New("cannot unmarshal the json response body")
)

const (
	httpTimeout = time.Second * 3
	// ReleaseAPI GitHub API v3 releases endpoint.
	// See: https://developer.github.com/v3/repos/releases/
	ReleaseAPI = "https://api.github.com/repos/bengarrett/retrotxtgo/releases/latest"
)

// API interface to store the JSON results from GitHub.
type API map[string]interface{}

// Endpoint request an API endpoint from the URL.
// A HTTP ETag can be provided to validate local data cache against the server.
// The return is true when the etag value matches the server's ETag header.
func Endpoint(url, etag string) (bool, API, error) {
	resp, body, err := Get(url, etag)
	if err != nil {
		return false, API{}, fmt.Errorf("endpoint get failed: %w", err)
	}
	defer resp.Body.Close()
	if etag != "" {
		s := resp.StatusCode
		if s == 304 || (s == 200 && body == nil) {
			// Not Modified
			return true, API{}, nil
		}
	}
	if ok := json.Valid(body); !ok {
		return false, API{}, fmt.Errorf("endpoint %s: %w", url, ErrJSON)
	}
	var data API
	if err := json.Unmarshal(body, &data); err != nil {
		return false, API{}, fmt.Errorf("endpoint %s: %w", url, ErrMash)
	}
	data["etag"] = resp.Header.Get("Etag")
	return false, data, nil
}

// Get fetches a URL and returns both its response and body.
// If an etag is provided a "If-None-Match" header request will be included.
func Get(url, etag string) (*http.Response, []byte, error) {
	client := &http.Client{
		Timeout: httpTimeout,
	}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, httpTimeout)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	defer cancel()
	if err != nil {
		return nil, nil, fmt.Errorf("getting a new request error: %w", err)
	}
	if etag != "" {
		req.Header.Set("If-None-Match", etag)
	}
	req.Header.Set("User-Agent", userAgent())
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("requesting to set the get user-agent header: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("reading the response body failed: %w", err)
	}
	return resp, body, resp.Body.Close()
}

// Ping requests a URL and determines if the status is ok.
func Ping(url string) (bool, error) {
	ok := false
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
	req.Header.Set("User-Agent", userAgent())
	resp, err := client.Do(req)
	if err != nil {
		return ok, fmt.Errorf("requesting to set the ping user-agent header: %w", err)
	}
	return (resp.StatusCode >= 200 && resp.StatusCode <= 299), resp.Body.Close()
}

func userAgent() string {
	return fmt.Sprintf("%s version ping", meta.Bin)
}
