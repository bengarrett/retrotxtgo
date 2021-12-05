package input_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/input"
	"github.com/gookit/color"
)

func ExampleDefaults() {
	s := input.Defaults(get.Author)
	fmt.Println(s)

	blank := input.Defaults("")
	fmt.Println(blank)
	// Output: Your name goes here
}

func ExamplePrintMeta() {
	color.Enable = false
	s, err := input.PrintMeta("html.meta.author", "value")
	if err != nil {
		log.Print(err)
		return
	}
	spl := strings.Split(s, "\n")
	fmt.Print(strings.Join(spl[5:], "\n"))
	// Output:
	//   About this value: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/meta/name
	//   Defines the name of the page authors.
}
