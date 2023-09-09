package ansi_test

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/bengarrett/retrotxtgo/pkg/ansi"
)

func ExampleHTMLReader() {
	file, err := os.Open("../../static/ansi/ansi-aixterm.ans")
	if err != nil {
		log.Print(err)
	}
	defer file.Close()

	b := bytes.Buffer{}
	if err := ansi.HTMLReader(&b, file); err != nil {
		log.Print(err)
	}
	fmt.Fprintln(os.Stdout, b.String())
	// Output:
}

func ExampleHTMLReader_ansi() {
	file, err := os.Open("../../static/ansi/preview_01.ans")
	if err != nil {
		log.Print(err)
	}
	defer file.Close()

	b := bytes.Buffer{}
	if err := ansi.HTMLReader(&b, file); err != nil {
		log.Print(err)
	}
	fmt.Fprintln(os.Stdout, b.String())
	// Output:
}
