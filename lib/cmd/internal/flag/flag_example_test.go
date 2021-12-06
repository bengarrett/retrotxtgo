package flag_test

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/lib/convert"
)

func ExampleEndOfFile() {
	var f convert.Flag
	f.Controls = []string{"eof"}
	fmt.Print(flag.EndOfFile(f))
	// Output: true
}
