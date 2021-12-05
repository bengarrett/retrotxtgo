package set_test

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/lib/config/internal/set"
	"github.com/gookit/color"
)

func Example_recommend() {
	color.Enable = false
	fmt.Print(set.Recommend(""))
	// Output: (suggestion: do not use)
}
