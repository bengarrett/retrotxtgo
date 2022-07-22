package colorise_test

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/config/internal/colorise"
)

func ExampleChromaNames() {
	w := new(bytes.Buffer)
	colorise.ChromaNames(w, "css")
	s := strings.Split(w.String(), ">")
	fmt.Print(s[0])
	// Output:0 <[38;5;16mabap[0m=[38;5;70m"abap"[0m
}

func ExampleChromaNamesMono_css() {
	w := new(bytes.Buffer)
	colorise.ChromaNamesMono(w, "css")
	s := strings.Split(w.String(), "\n")
	fmt.Print(s[0])
	// Output:0 <abap="abap">
}

func ExampleChromaNamesMono_json() {
	w := new(bytes.Buffer)
	colorise.ChromaNamesMono(w, "json")
	s := strings.Split(w.String(), "\n")
	fmt.Print(s[0])
	// Output:0 { "abap":"abap" }
}
