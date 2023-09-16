package info_test

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/pkg/info"
)

func ExampleLineBreak_Find() {
	lb := &info.LineBreak{}

	// A Windows example of a line break.
	windows := [2]rune{rune('\r'), rune('\n')}
	lb.Find(windows)
	fmt.Printf("%s %q %d %X\n", lb.Abbr, lb.Escape, lb.Decimal, lb.Decimal)

	// A macOS, Linux, or Unix example of a line break.
	linux := [2]rune{rune('\n')}
	lb.Find(linux)
	fmt.Printf("%s %q %d %X\n", lb.Abbr, lb.Escape, lb.Decimal, lb.Decimal)

	// Output: CRLF "\r\n" [13 10] [D A]
	// LF "\n" [10 0] [A 0]
}

func ExampleLineBreak_Total() {
	lb := &info.LineBreak{}
	linux := [2]rune{rune('\n')}
	lb.Find(linux)
	// Count the number of lines in a file.
	total, _ := lb.Total("testdata/textlf.txt")
	fmt.Println(total)
	// Output: 7
}
