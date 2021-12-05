package color_test

import (
	"fmt"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/config/internal/color"
)

func ExampleChromaNamesMono_css() {
	c := color.ChromaNamesMono("css")
	s := strings.Split(c, "\n")
	fmt.Print(s[0])
	// Output:0 <abap="abap">
}

func ExampleChromaNamesMono_json() {
	c := color.ChromaNamesMono("json")
	s := strings.Split(c, "\n")
	fmt.Print(s[0])
	// Output:0 { "abap":"abap" }
}
