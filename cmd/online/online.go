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
	ErrJSON   = errors.New("the response body syntax is not json")
	ErrMash   = errors.New("cannot unmarshal the json response body")
	ErrNoResp = errors.New("the response is nil and unusable")
)

const (
	// ReleaseAPI GitHub API v3 releases endpoint.
	// See: https://developer.github.com/v3/repos/releases/
	ReleaseAPI = "https://api.github.com/repos/bengarrett/retrotxtgo/releases/latest"
	timeout    = time.Second * 3
)

// API interface to store the JSON results from GitHub.
type API map[string]any

// Endpoint requests an API endpoint from the URL.
// A HTTP ETag can be provided to validate local data cache against the server.
// It also reports whether the etag value matches the server ETag header.
func Endpoint(ctx context.Context, url, etag string) (bool, API, error) {
	const format = "online api endpoint get %s: %w"
	resp, body, err := Get(ctx, url, etag)
	if err != nil {
		return false, API{}, fmt.Errorf(format, "failed", err)
	}
	if resp == nil {
		return false, API{}, fmt.Errorf(format, "no response", ErrNoResp)
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
		return false, API{}, fmt.Errorf(format, url, ErrJSON)
	}
	var data API
	if err := json.Unmarshal(body, &data); err != nil {
		return false, API{}, fmt.Errorf(format, url, ErrMash)
	}
	if data == nil {
		return false, API{}, fmt.Errorf(format, url, ErrMash)
	}
	val := resp.Header.Get("Etag")
	data["etag"] = val
	return false, data, nil
}

// Get fetches a URL and returns both its response and body.
// If an etag is provided a "If-None-Match" header request will be included.
func Get(ctx context.Context, url, etag string) (*http.Response, []byte, error) {
	client := &http.Client{
		Timeout: timeout,
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	defer cancel()
	if err != nil {
		const format = "getting a new request error: %w"
		return nil, nil, fmt.Errorf(format, err)
	}
	if etag != "" {
		req.Header.Set("If-None-Match", etag)
	}
	req.Header.Set("User-Agent", userAgent())
	resp, err := client.Do(req)
	if err != nil {
		const format = "requesting to set the get user-agent header: %w"
		return nil, nil, fmt.Errorf(format, err)
	}
	if resp == nil {
		const format = "getting response: %w"
		return nil, nil, fmt.Errorf(format, ErrNoResp)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		const format = "reading the response body failed: %w"
		return nil, nil, fmt.Errorf(format, err)
	}
	return resp, body, nil
}

// Ping requests a URL and reports whether if the status is successful.
// A server response status code between 200 and 299 is considered a success.
func Ping(ctx context.Context, url string) (bool, error) {
	client := &http.Client{
		Timeout: timeout,
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	defer cancel()
	if err != nil {
		const format = "pinging a new request error: %w"
		return false, fmt.Errorf(format, err)
	}
	req.Header.Set("User-Agent", userAgent())
	resp, err := client.Do(req)
	if err != nil {
		const format = "requesting to set the ping user-agent header: %w"
		return false, fmt.Errorf(format, err)
	}
	if resp == nil {
		const format = "ping response: %w"
		return false, fmt.Errorf(format, ErrNoResp)
	}
	defer resp.Body.Close()
	const ok, maximum = http.StatusOK, 299
	success2xx := resp.StatusCode >= ok && resp.StatusCode <= maximum
	return success2xx, nil
}

func userAgent() string {
	return meta.Bin + " version ping"
}
