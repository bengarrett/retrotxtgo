package online

import "fmt"

func ExamplePing() {
	pingOk, _ := Ping("https://demozoo.org/api/v1/productions/126496/")
	pingBad, _ := Ping("https://example.com/this/url/does/not/exist")
	fmt.Println(pingOk, pingBad)
	// Output: true false
}

func ExampleRequest() {
	body, _ := Get("https://demozoo.org/api/v1/productions/126496/")
	fmt.Println(string(body))
	// Output: >
}
