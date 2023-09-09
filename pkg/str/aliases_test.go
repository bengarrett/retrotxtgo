// Package str for strings and styles.
package str_test

import (
	"fmt"
	"os"

	"github.com/bengarrett/retrotxtgo/pkg/str"
	"github.com/gookit/color"
)

func ExampleAlert() {
	color.Enable = false
	fmt.Fprint(os.Stdout, str.Alert())
	// Output:Problem:
}

func ExampleInfo() {
	color.Enable = false
	fmt.Fprint(os.Stdout, str.Info())
	// Output:Information:
}

func ExampleColSec() {
	color.Enable = false
	fmt.Fprint(os.Stdout, str.ColSec("Hi"))
	// Output:Hi
}

func ExampleColCmt() {
	color.Enable = false
	fmt.Fprint(os.Stdout, str.ColCmt("Hi"))
	// Output:Hi
}

func ExampleColFuz() {
	color.Enable = false
	fmt.Fprint(os.Stdout, str.ColFuz("Hi"))
	// Output:Hi
}

func ExampleItalic() {
	color.Enable = false
	fmt.Fprint(os.Stdout, str.Italic("Hi"))
	// Output:Hi
}

func ExampleColInf() {
	color.Enable = false
	fmt.Fprint(os.Stdout, str.ColInf("Hi"))
	// Output:Hi
}

func ExampleColPri() {
	color.Enable = false
	fmt.Fprint(os.Stdout, str.ColPri("Hi"))
	// Output:Hi
}

func ExampleColSuc() {
	color.Enable = false
	fmt.Fprint(os.Stdout, str.ColSuc("Hi"))
	// Output:Hi
}

func ExampleBool() {
	fmt.Fprint(os.Stdout, str.Bool(true))
	fmt.Fprint(os.Stdout, str.Bool(false))
	// Output:✓✗
}

func ExampleDefault() {
	fmt.Fprint(os.Stdout, str.Default("hi, bye", "hi"))
	// Output:hi, bye (default "hi")
}

func ExampleOptions() {
	fmt.Fprint(os.Stdout, str.Options("this is an example of a list of options",
		false, false, "option3", "option2", "option1"))
	// Output:this is an example of a list of options.
	//   Options: option1, option2, option3
}

func ExampleRequired() {
	fmt.Fprint(os.Stdout, str.Required("hi"))
	// Output:hi (required)
}
