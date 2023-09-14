package nl_test

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/pkg/fsys/nl"
)

func ExampleNewLine() {
	s := nl.NewLine(nl.Win)
	fmt.Printf("%q\n", s)
	s = nl.NewLine(nl.Mac)
	fmt.Printf("%q\n", s)
	s = nl.NewLine(nl.Linux)
	fmt.Printf("%q\n", s)
	// Output: "\r\n"
	// "\r"
	// "\n"
}
