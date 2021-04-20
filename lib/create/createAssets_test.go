package create

import (
	"fmt"

	"github.com/gookit/color"
)

func Example_bytesStats() {
	// Disable ANSI color output
	color.Enable = false

	fmt.Println(bytesStats("filename.txt", 0))
	fmt.Println(bytesStats("filename.txt", 123))
	fmt.Println(bytesStats("filename.txt", 1234567890))
	// Output:saved to filename.txt (zero-byte file)
	//saved to filename.txt, 123B
	//saved to filename.txt, 1.23 GB (1234567890)
}
