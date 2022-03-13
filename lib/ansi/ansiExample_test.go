package ansi_test

import (
	"log"
	"os"

	"github.com/bengarrett/retrotxtgo/lib/ansi"
)

func ExampleHTMLReader() {
	file, err := os.Open("../../static/ansi/ansi-aixterm.ans")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err := ansi.HTMLReader(os.Stdout, file); err != nil {
		log.Fatal(err)
	}
	// Output:
}

func ExampleHTMLReader_() {
	file, err := os.Open("../../static/ansi/preview_01.ans")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err := ansi.HTMLReader(os.Stdout, file); err != nil {
		log.Fatal(err)
	}
	// Output:
}
