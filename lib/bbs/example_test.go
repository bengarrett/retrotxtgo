package bbs_test

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/bbs"
)

func ExampleBBS_String() {
	fmt.Print(bbs.PCBoard)
	// Output: PCBoard @X
}

func ExampleBBS_Name() {
	fmt.Print(bbs.PCBoard.Name())
	// Output: PCBoard
}

func ExampleBBS_Bytes() {
	b := bbs.PCBoard.Bytes()
	fmt.Printf("%s %v", b, b)
	// Output: @X [64 88]
}

func ExampleBBS_HTML() {
	var out bytes.Buffer
	src := []byte("@X03Hello world")
	if err := bbs.PCBoard.HTML(&out, src); err != nil {
		log.Print(err)
	}
	fmt.Print(out.String())
	// Output: <i class="PB0 PF3">Hello world</i>
}

func ExampleBBS_CSS() {
	css, err := bbs.PCBoard.CSS()
	if err != nil {
		fmt.Print(err)
	}
	fmt.Print(css)
	// Output:
}

func ExampleHTML() {
	var out bytes.Buffer
	src := strings.NewReader("@X03Hello world")
	if err := bbs.HTML(&out, src); err != nil {
		fmt.Print(err)
	}
	fmt.Print(out.String())
	// Output: <i class="PB0 PF3">Hello world</i>
}

func ExampleFind() {
	r := strings.NewReader("@X03Hello world")
	f := bbs.Find(r)
	fmt.Printf("Reader is in a %s BBS format", f.Name())
	// Output: Reader is in a PCBoard BBS format
}

func ExampleIsPCBoard() {
	b := []byte("@X03Hello world")
	fmt.Printf("Is PCBoard BBS text: %v", bbs.IsPCBoard(b))
	// Output: Is PCBoard BBS text: true
}

func ExampleFieldsBars() {
	s := "|03Hello |07|19world"
	fmt.Printf("Color sequences: %d", len(bbs.FieldsBars(s)))
	// Output: Color sequences: 3
}

func ExampleHTMLRenegade() {
	var out bytes.Buffer
	src := "|03Hello |07|19world"
	if err := bbs.HTMLRenegade(&out, src); err != nil {
		log.Print(err)
	}
	fmt.Print(out.String())
	// Output: <i class="P0 P3">Hello </i><i class="P0 P7"></i><i class="P19 P7">world</i>
}

func ExampleFieldsCelerity() {
	s := "|cHello |C|S|wworld"
	fmt.Printf("Color sequences: %d", len(bbs.FieldsCelerity(s)))
	// Output: Color sequences: 4
}

func ExampleHTMLCelerity() {
	var out bytes.Buffer
	src := "|cHello |C|S|wworld"
	if err := bbs.HTMLCelerity(&out, src); err != nil {
		fmt.Print(err)
	}
	fmt.Print(out.String())
	// Output: <i class="PBk PFc">Hello </i><i class="PBk PFC"></i><i class="PBw PFC">world</i>
}

func ExampleFieldsPCBoard() {
	s := "@X03Hello world"
	fmt.Printf("Color sequences: %d", len(bbs.FieldsPCBoard(s)))
	// Output: Color sequences: 1
}

func ExampleHTMLPCBoard() {
	var out bytes.Buffer
	src := "@X03Hello world"
	if err := bbs.HTMLPCBoard(&out, src); err != nil {
		log.Print(err)
	}
	fmt.Print(out.String())
	// Output: <i class="PB0 PF3">Hello world</i>
}
