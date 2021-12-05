package get_test

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
)

func ExampleTip() {
	tip := get.Tip()
	fmt.Print(tip[get.FontFamily])
	// Output: specifies the font to use with the HTML
}

func ExampleTextEditor() {
	ed := get.TextEditor()
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
