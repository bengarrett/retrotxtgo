package create_test

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/lib/create"
	"github.com/gookit/color"
)

func ExampleColorScheme() {
	fmt.Print(create.ColorScheme()[0])
	// Output: normal
}

func ExampleReferrer() {
	fmt.Print(create.Referrer()[1])
	// Output: origin
}

func ExampleRobots() {
	fmt.Print(create.Robots()[2])
	// Output: follow
}

func ExampleStats() {
	// Disable ANSI color output
	color.Enable = false

	fmt.Println(create.Stats("filename.txt", 0))
	fmt.Println(create.Stats("filename.txt", 123))
	fmt.Println(create.Stats("filename.txt", 1234567890))
	// Output:saved to filename.txt (zero-byte file)
	// saved to filename.txt, 123B
	// saved to filename.txt, 1.23 GB (1234567890)
}
