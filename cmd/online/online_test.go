package online_test

import (
	"net/http"
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd/online"
	"github.com/nalgeon/be"
)

func TestErrors(t *testing.T) {
	t.Parallel()

	be.Equal(t, online.ErrJSON.Error(), "the response body syntax is not json")
	be.Equal(t, online.ErrMash.Error(), "cannot unmarshal the json response body")
	be.Equal(t, online.ErrNoResp.Error(), "the response is nil and unusable")
}

func TestReleaseAPI(t *testing.T) {
	t.Parallel()
	be.Equal(t, online.ReleaseAPI, "https://api.github.com/repos/bengarrett/retrotxtgo/releases/latest")
}

func TestPing(t *testing.T) {
	t.Parallel()

	// Test with invalid URL
	ok, err := online.Ping(t.Context(), "http://this-url-definitely-does-not-exist-anywhere.com/invalid")
	be.Equal(t, ok, false)
	be.True(t, err != nil)
	ok, err = online.Ping(t.Context(), online.ReleaseAPI)
	if err != nil {
		t.Logf("Network error (expected in some environments): %v", err)
		return
	}
	be.True(t, ok)
}

func TestGet(t *testing.T) {
	t.Parallel()
	resp, body, err := online.Get(t.Context(), "http://this-url-definitely-does-not-exist-anywhere.com/invalid", "")
	if resp != nil {
		_ = resp.Body.Close()
	}
	be.True(t, err != nil)
	be.Equal(t, resp, (*http.Response)(nil))
	be.Equal(t, body, ([]byte)(nil))

	// Test with GitHub API (real network call - skipped)
	t.Skip("Skipping real network test for GitHub API")
	resp, body, err = online.Get(t.Context(), online.ReleaseAPI, "")
	if err != nil {
		t.Logf("Network error (expected in some environments): %v", err)
		return
	}
	be.True(t, resp != nil)
	be.True(t, len(body) > 0)
}

func TestEndpoint(t *testing.T) {
	t.Parallel()

	// Test with invalid URL
	matched, api, err := online.Endpoint(t.Context(), "http://this-url-definitely-does-not-exist-anywhere.com/invalid", "")
	be.Equal(t, matched, false)
	be.True(t, err != nil)
	be.Equal(t, api, online.API{})

	// Test with GitHub API (real network call - skipped)
	t.Skip("Skipping real network test for GitHub API")
	matched, api, err = online.Endpoint(t.Context(), online.ReleaseAPI, "")
	if err != nil {
		t.Logf("Network error (expected in some environments): %v", err)
		return
	}
	be.Equal(t, matched, false)
	be.True(t, len(api) > 0)
}
