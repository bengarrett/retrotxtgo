package online

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

const userAgent = "retrotxt ping"
const httpTimeout = time.Second * 3

// API blah
type API map[string]interface{}

// ReleaseAPI GitHub API v3 releases endpoint
// See: https://developer.github.com/v3/repos/releases/
const ReleaseAPI = "https://api.github.com/repos/bengarrett/retrotxtgo/releases/latest"

/*
access-control-expose-headers: ETag
cache-control: public, max-age=60, s-maxage=60
last-modified: Fri, 12 Jun 2020 01:04:23 GMT
etag: W/"3715383704fac6f3568e9039b347937a"
*/
/*
cache-control: no-cache
If-Modified-Since: Fri, 12 Jun 2020 01:04:23 GMT
If-None-Match: W/"3715383704fac6f3568e9039b347937a"
*/

// Endpoint ..
func Endpoint(etag, url string) (useCache bool, data API, err error) {
	useCache = false
	header, body, err := Get(url, etag)
	if err != nil {
		return useCache, data, err
	}
	if etag != "" {
		s := header.Get("Status")
		if s == "304 Not Modified" || (s == "200 OK" && data == nil) {
			// Not Modified
			return true, data, nil
		}
	}
	if ok := json.Valid(body); !ok {
		return useCache, data, errors.New("the response from is not in json syntax: " + url)
	}
	err = json.Unmarshal(body, &data)
	data["etag"] = header.Get("Etag")
	return useCache, data, nil
}

// Get ...
func Get(url, etag string) (header http.Header, body []byte, err error) {
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
	resp, err := client.Do(req)
	if err != nil {
		return nil, body, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	return resp.Header, body, err
}

// Ping requests a URL and determines if its status is ok.
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
