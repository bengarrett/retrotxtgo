// Package str for strings and styles.
// nolint:gocritic,gochecknoinits
package str

import (
	"fmt"

	"github.com/gookit/color"
)

func init() {
	color.Enable = false
}

func ExampleAlert() {
	fmt.Print(Alert())
	// Output:Problem:
}
func ExampleInfo() {
	fmt.Print(Info())
	// Output:Information:
}

func ExampleColSec() {
	fmt.Print(ColSec("Hi"))
	// Output:Hi
}
func ExampleColCmt() {
	fmt.Print(ColCmt("Hi"))
	// Output:Hi
}
func ExampleColFuz() {
	fmt.Print(ColFuz("Hi"))
	// Output:Hi
}
func ExampleItalic() {
	fmt.Print(Italic("Hi"))
	// Output:Hi
}
func ExampleColInf() {
	fmt.Print(ColInf("Hi"))
	// Output:Hi
}
func ExampleColPri() {
	fmt.Print(ColPri("Hi"))
	// Output:Hi
}
func ExampleColSuc() {
	fmt.Print(ColSuc("Hi"))
	// Output:Hi
}

func ExampleBool() {
	fmt.Print(Bool(true))
	fmt.Print(Bool(false))
	// Output:✓✗
}

func ExampleDefault() {
	fmt.Print(Default("hi, bye", "hi"))
	// Output:hi, bye (default "hi")
}

func ExampleOptions() {
	fmt.Print(Options("this is an example of a list of options",
		false, false, "option3", "option2", "option1"))
	// Output:this is an example of a list of options.
	//   Options: option1, option2, option3
}

func ExampleRequired() {
	fmt.Print(Required("hi"))
	// Output:hi (required)
}
