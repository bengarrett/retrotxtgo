package input_test

import (
	"fmt"
	"log"

	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/input"
)

func ExampleDefaults() {
	s := input.Defaults(get.Author)
	fmt.Println(s)

	blank := input.Defaults("")
	fmt.Println(blank)
	// Output: Your name goes here
}

func ExamplePrintMeta() {
	s, err := input.PrintMeta("html.meta.author", "value")
	if err != nil {
		log.Print(err)
		return
	}
	fmt.Print(s)
	// Output:
	// <[1m[38;5;254mhead[0m>
	//     <[1m[38;5;254mmeta[0m [38;5;6mname[0m=[1m[38;5;51m"author"[0m [38;5;6mvalue[0m=[1m[38;5;51m"value"[0m>
	//   </[1m[38;5;254mhead[0m>
	//
	//   About this value: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/meta/name
	//   Defines the name of the page authors.
}
