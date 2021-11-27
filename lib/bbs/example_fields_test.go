package bbs_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/bbs"
)

func ExampleFields() {
	r := strings.NewReader("@X03Hello world")
	s, err := bbs.Fields(r)
	if err != nil {
		log.Print(err)
	}
	fmt.Printf("Color sequences: %d", len(s))
	// Output: Color sequences: 1
}
