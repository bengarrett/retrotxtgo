package bbs_test

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/bbs"
)

func ExampleFields() {
	r := strings.NewReader("@X03Hello @XF0world")
	s, b, err := bbs.Fields(r)
	if err != nil {
		log.Print(err)
	}
	fmt.Printf("Found %d %s sequences.", len(s), b)
	// Output: Found 2 PCBoard @X sequences.
}

func ExampleFields_none() {
	r := strings.NewReader("Hello world.")
	_, _, err := bbs.Fields(r)
	if errors.Is(err, bbs.ErrColorCodes) {
		fmt.Print(err)
	}
	// Output: no bbs color codes found
}
