package input_test

import (
	"bytes"
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

func ExamplePortInfo() {
	color.Enable = false
	fmt.Print(input.PortInfo())
	// Output: 1-65535 (suggestion: 8086)
}

func ExamplePrintMeta() {
	color.Enable = false
	b := new(bytes.Buffer)
	if err := input.PrintMeta(b, "html.meta.author", "value"); err != nil {
		log.Print(err)
		return
	}
	spl := strings.Split(b.String(), "\n")
	fmt.Print(strings.Join(spl[5:], "\n"))
	// Output:
	//   About this value: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/meta/name
	//   Defines the name of the page authors.
}
