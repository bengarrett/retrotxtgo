package online_test

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/pkg/online"
)

func ExamplePing() {
	ok, _ := online.Ping("https://example.org")
	fmt.Print(ok)
	// Output: true
}
