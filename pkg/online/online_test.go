package online_test

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bengarrett/retrotxtgo/pkg/online"
)

func ExampleEndpoint() {
	_, p, _ := online.Endpoint("https://demozoo.org/api/v1/productions/126496/", `W/"0708012ac3fb439a46dd5156195901b4"`)
	fmt.Fprintln(os.Stdout, "id:", p["id"])
	// Output: id: 126496
}

func ExamplePing() {
	ok, _ := online.Ping("https://example.org")
	bad, _ := online.Ping("https://example.com/this/url/does/not/exist")
	fmt.Fprintln(os.Stdout, ok, bad)
	// Output: true false
}

func ExampleGet() {
	_, d, _ := online.Get("https://demozoo.org/api/v1/productions/126496/", "")
	_, g, _ := online.Get(online.ReleaseAPI, "")
	fmt.Fprintln(os.Stdout, "valid json api?", json.Valid(d))
	fmt.Fprintln(os.Stdout, "valid github api?", json.Valid(g))
	// Output: valid json api? true
	// valid github api? true
}
