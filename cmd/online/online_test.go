package online_test

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/cmd/online"
)

func ExamplePing() {
	ok, _ := online.Ping("https://example.org")
	fmt.Print(ok)
	// Output: true
}
