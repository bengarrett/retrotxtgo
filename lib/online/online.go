package online

import (
	"io/ioutil"
	"net/http"
	"time"
)

const userAgent = "retrotxtgo ping"
const httpTimeout = time.Second * 3

// ReleaseAPI GitHub API v3 releases endpoint
// See: https://developer.github.com/v3/repos/releases/
const ReleaseAPI = "https://api.github.com/repos/bengarrett/RetroTxt/releases"

// Get ...
func Get(url string) (body []byte, err error) {
	client := &http.Client{
		Timeout: httpTimeout,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return body, err
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return body, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// Ping ..
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
