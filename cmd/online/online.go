// Package online is for simple HTTP interactions with the GitHub API.
// It is used to fetch the latest release information of the program.
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
	ErrJSON = errors.New("the response body syntax is not json")
	ErrMash = errors.New("cannot unmarshal the json response body")
)

const (
	// ReleaseAPI GitHub API v3 releases endpoint.
	// See: https://developer.github.com/v3/repos/releases/
	ReleaseAPI = "https://api.github.com/repos/bengarrett/retrotxtgo/releases/latest"
	timeout    = time.Second * 3
)

// API interface to store the JSON results from GitHub.
type API map[string]interface{}

// Endpoint requests an API endpoint from the URL.
// A HTTP ETag can be provided to validate local data cache against the server.
// It also reports whether the etag value matches the server ETag header.
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
		Timeout: timeout,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
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
	return resp, body, nil
}

// Ping requests a URL and reports whether if the status is successful.
// A server response status code between 200 and 299 is considered a success.
func Ping(url string) (bool, error) {
	client := &http.Client{
		Timeout: timeout,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	defer cancel()
	if err != nil {
		return false, fmt.Errorf("pinging a new request error: %w", err)
	}
	req.Header.Set("User-Agent", userAgent())
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("requesting to set the ping user-agent header: %w", err)
	}
	defer resp.Body.Close()
	const ok, maximum = http.StatusOK, 299
	success2xx := resp.StatusCode >= ok && resp.StatusCode <= maximum
	return success2xx, nil
}

func userAgent() string {
	return meta.Bin + " version ping"
}
