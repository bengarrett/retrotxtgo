package cmd

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/lib/convert"
)

func Example_endOfFile() {
	var f convert.Flag
	f.Controls = []string{eof}
	fmt.Print(endOfFile(f))
	// Output: true
}
