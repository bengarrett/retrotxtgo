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

func ExampleCb() {
	fmt.Print(Cb("Hi"))
	// Output:Hi
}
func ExampleCc() {
	fmt.Print(Cc("Hi"))
	// Output:Hi
}
func ExampleCe() {
	fmt.Print(Ce("Hi"))
	// Output:Hi
}
func ExampleCf() {
	fmt.Print(Cf("Hi"))
	// Output:Hi
}
func ExampleCi() {
	fmt.Print(Ci("Hi"))
	// Output:Hi
}
func ExampleCinf() {
	fmt.Print(Cinf("Hi"))
	// Output:Hi
}
func ExampleCp() {
	fmt.Print(Cp("Hi"))
	// Output:Hi
}
func ExampleCs() {
	fmt.Print(Cs("Hi"))
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
