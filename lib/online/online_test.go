package online

import (
	"encoding/json"
	"fmt"
)

func ExampleEndpoint() {
	_, p, _ := Endpoint("https://demozoo.org/api/v1/productions/126496/", "")
	_, g, _ := Endpoint(ReleaseAPI, "")
	fmt.Println("id:", p["id"])
	fmt.Println("ver:", g["tag_name"])
	// Output: id: 126496
	// ver: 0.0.1
}

func ExamplePing() {
	pingOk, _ := Ping("https://example.org")
	pingBad, _ := Ping("https://example.com/this/url/does/not/exist")
	fmt.Println(pingOk, pingBad)
	// Output: true false
}

func ExampleGet() {
	_, d, _ := Get("https://demozoo.org/api/v1/productions/126496/", "")
	_, g, _ := Get(ReleaseAPI, "")
	fmt.Println("valid json api?", json.Valid(d))
	fmt.Println("valid github api?", json.Valid(g))
	// Output: valid json api? true
	// valid github api? true
}
