// Package str for strings and styles.
// nolint:gocritic,gochecknoinits
package str_test

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/gookit/color"
)

func init() {
	color.Enable = false
}

func ExampleAlert() {
	fmt.Print(str.Alert())
	// Output:Problem:
}
func ExampleInfo() {
	fmt.Print(str.Info())
	// Output:Information:
}

func ExampleColSec() {
	fmt.Print(str.ColSec("Hi"))
	// Output:Hi
}
func ExampleColCmt() {
	fmt.Print(str.ColCmt("Hi"))
	// Output:Hi
}
func ExampleColFuz() {
	fmt.Print(str.ColFuz("Hi"))
	// Output:Hi
}
func ExampleItalic() {
	fmt.Print(str.Italic("Hi"))
	// Output:Hi
}
func ExampleColInf() {
	fmt.Print(str.ColInf("Hi"))
	// Output:Hi
}
func ExampleColPri() {
	fmt.Print(str.ColPri("Hi"))
	// Output:Hi
}
func ExampleColSuc() {
	fmt.Print(str.ColSuc("Hi"))
	// Output:Hi
}

func ExampleBool() {
	fmt.Print(str.Bool(true))
	fmt.Print(str.Bool(false))
	// Output:✓✗
}

func ExampleDefault() {
	fmt.Print(str.Default("hi, bye", "hi"))
	// Output:hi, bye (default "hi")
}

func ExampleOptions() {
	fmt.Print(str.Options("this is an example of a list of options",
		false, false, "option3", "option2", "option1"))
	// Output:this is an example of a list of options.
	//   Options: option1, option2, option3
}

func ExampleRequired() {
	fmt.Print(str.Required("hi"))
	// Output:hi (required)
}
