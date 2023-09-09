package term_test

import (
	"fmt"
	"os"

	"github.com/bengarrett/retrotxtgo/pkg/term"
	"github.com/gookit/color"
)

func ExampleAlert() {
	color.Enable = false
	fmt.Fprint(os.Stdout, term.Alert())
	// Output:Problem:
}

func ExampleInfo() {
	color.Enable = false
	fmt.Fprint(os.Stdout, term.Info())
	// Output:Information:
}

func ExampleColSec() {
	color.Enable = false
	fmt.Fprint(os.Stdout, term.ColSec("Hi"))
	// Output:Hi
}

func ExampleColCmt() {
	color.Enable = false
	fmt.Fprint(os.Stdout, term.ColCmt("Hi"))
	// Output:Hi
}

func ExampleColFuz() {
	color.Enable = false
	fmt.Fprint(os.Stdout, term.ColFuz("Hi"))
	// Output:Hi
}

func ExampleColInf() {
	color.Enable = false
	fmt.Fprint(os.Stdout, term.ColInf("Hi"))
	// Output:Hi
}

func ExampleBool() {
	fmt.Fprint(os.Stdout, term.Bool(true))
	fmt.Fprint(os.Stdout, term.Bool(false))
	// Output:✓✗
}

func ExampleOptions() {
	fmt.Fprint(os.Stdout, term.Options("this is an example of a list of options",
		false, false, "option3", "option2", "option1"))
	// Output:this is an example of a list of options.
	//   Options: option1, option2, option3
}
