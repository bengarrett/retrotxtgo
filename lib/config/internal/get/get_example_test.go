package get_test

import (
	"fmt"
	"os"

	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
)

func ExampleTextEditor() {
	ed := get.TextEditor(os.Stdout)
	found := len(ed) > 0
	fmt.Print("Text editor found? ", found)
	// Output: Text editor found? true
}

func ExampleDiscEditor() {
	ed := get.DiscEditor()
	found := len(ed) > 0
	fmt.Print("Text editor found? ", found)
	// Output: Text editor found? true
}
