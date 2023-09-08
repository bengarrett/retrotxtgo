package flag_test

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/pkg/convert"
)

func ExampleEndOfFile() {
	var f convert.Flag
	f.Controls = []string{"eof"}
	fmt.Print(flag.EndOfFile(f))
	// Output: true
}
