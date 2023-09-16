package nl_test

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/pkg/nl"
)

func ExampleNewLine() {
	s := nl.NewLine(nl.Windows)
	fmt.Printf("%q\n", s)
	s = nl.NewLine(nl.Macintosh)
	fmt.Printf("%q\n", s)
	s = nl.NewLine(nl.Linux)
	fmt.Printf("%q\n", s)
	// Output: "\r\n"
	// "\r"
	// "\n"
}
