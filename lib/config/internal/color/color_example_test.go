package color_test

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/config/internal/color"
)

func ExampleChromaNamesMono_css() {
	w := new(bytes.Buffer)
	color.ChromaNamesMono(w, "css")
	s := strings.Split(w.String(), "\n")
	fmt.Print(s[0])
	// Output:0 <abap="abap">
}

func ExampleChromaNamesMono_json() {
	w := new(bytes.Buffer)
	color.ChromaNamesMono(w, "json")
	s := strings.Split(w.String(), "\n")
	fmt.Print(s[0])
	// Output:0 { "abap":"abap" }
}
