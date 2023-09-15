package online_test

import (
	"encoding/json"
	"fmt"

	"github.com/bengarrett/retrotxtgo/pkg/online"
)

func ExampleEndpoint() {
	etag := ""
	cached, p, _ := online.Endpoint("https://api.github.com/repos/bengarrett/retrotxtgo/releases/121077170", etag)
	fmt.Print(cached, p["name"])
	// Output: v0.4.0 false
}

func ExampleGet() {
	etag := ""
	resp, body, _ := online.Get("https://api.github.com/repos/bengarrett/retrotxtgo/releases/121077170", etag)
	fmt.Print(resp.StatusCode, json.Valid(body))
	// Output: 200 true
}

func ExamplePing() {
	ok, _ := online.Ping("https://example.org")
	fmt.Print(ok)
	// Output: true
}
