package bbs_test

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/bbs"
)

func ExampleFieldsBars() {
	s := "|03Hello |07|19world"
	l := len(bbs.FieldsBars(s))
	fmt.Printf("Color sequences: %d", l)
	// Output: Color sequences: 3
}

func ExampleFieldsCelerity() {
	s := "|cHello |C|S|wworld"
	l := len(bbs.FieldsCelerity(s))
	fmt.Printf("Color sequences: %d", l)
	// Output: Color sequences: 4
}

func ExampleFieldsPCBoard() {
	s := "@X03Hello world"
	l := len(bbs.FieldsPCBoard(s))
	fmt.Printf("Color sequences: %d", l)
	// Output: Color sequences: 1
}

func ExampleHTML() {
	var out bytes.Buffer
	src := strings.NewReader("@X03Hello world")
	if _, err := bbs.HTML(&out, src); err != nil {
		fmt.Print(err)
	}
	fmt.Print(out.String())
	// Output: <i class="PB0 PF3">Hello world</i>
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

func ExampleHTMLCelerity() {
	var out bytes.Buffer
	src := "|cHello |C|S|wworld"
	if err := bbs.HTMLCelerity(&out, src); err != nil {
		fmt.Print(err)
	}
	fmt.Print(out.String())
	// Output: <i class="PBk PFc">Hello </i><i class="PBk PFC"></i><i class="PBw PFC">world</i>
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

func ExampleIsCelerity() {
	b := []byte("|cHello |C|S|wworld")
	fmt.Printf("Is b Celerity BBS text? %v", bbs.IsCelerity(b))
	// Output: Is b Celerity BBS text? true
}

func ExampleIsPCBoard() {
	b := []byte("@X03Hello world")
	fmt.Printf("Is b PCBoard BBS text? %v", bbs.IsPCBoard(b))
	// Output: Is b PCBoard BBS text? true
}

func ExampleIsRenegade() {
	b := []byte("|03Hello |07|19world")
	fmt.Printf("Is b Renegade BBS text? %v", bbs.IsRenegade(b))
	// Output: Is b Renegade BBS text? true
}

func ExampleIsTelegard() {
	b := []byte("`07Hello world")
	fmt.Printf("Is b Telegard BBS text? %v", bbs.IsTelegard(b))
	// Output: Is b Telegard BBS text? true
}

func ExampleIsWHash() {
	b := []byte("|#7Hello world")
	fmt.Printf("Is b WVIV BBS # text? %v", bbs.IsWHash(b))
	// Output: Is b WVIV BBS # text? true
}
func ExampleIsWHeart() {
	b := []byte("\x037Hello world")
	fmt.Printf("Is b WWIV BBS ♥ text? %v", bbs.IsWHeart(b))
	// Output: Is b WWIV BBS ♥ text? true
}
func ExampleIsWildcat() {
	b := []byte("@0F@Hello world")
	fmt.Printf("Is b Wildcat! BBS text? %v", bbs.IsWildcat(b))
	// Output: Is b Wildcat! BBS text? true
}

func ExampleFind() {
	r := strings.NewReader("@X03Hello world")
	f := bbs.Find(r)
	fmt.Printf("Reader is in a %s BBS format", f.Name())
	// Output: Reader is in a PCBoard BBS format
}

func ExampleFind_none() {
	r := strings.NewReader("Hello world")
	f := bbs.Find(r)
	if !f.Valid() {
		fmt.Print("reader is plain text")
	}
	// Output: reader is plain text
}

func ExampleBBS_Bytes() {
	b := bbs.PCBoard.Bytes()
	fmt.Printf("%s %v", b, b)
	// Output: @X [64 88]
}

func ExampleBBS_CSS() {
	var buf bytes.Buffer
	if err := bbs.PCBoard.CSS(&buf); err != nil {
		fmt.Print(err)
	}

	f, err := os.OpenFile("pcboard.css", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(f.Name()) // clean up

	if _, err := buf.WriteTo(f); err != nil {
		log.Fatal(err)
	}
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

func ExampleBBS_Name() {
	fmt.Print(bbs.PCBoard.Name())
	// Output: PCBoard
}

func ExampleBBS_Remove() {
	var out bytes.Buffer
	src := []byte("@X03Hello @X07world")
	if err := bbs.PCBoard.Remove(&out, src); err != nil {
		log.Print(err)
	}
	fmt.Print(out.String())
	// Output: Hello world
}

func ExampleBBS_Remove_find() {
	var out bytes.Buffer
	src := []byte("@X03Hello @X07world")
	r := bytes.NewReader(src)
	b := bbs.Find(r)
	if err := b.Remove(&out, src); err != nil {
		log.Print(err)
	}
	fmt.Print(out.String())
	// Output: Hello world
}

func ExampleBBS_String() {
	fmt.Print(bbs.PCBoard)
	// Output: PCBoard @X
}

func ExampleBBS_Valid() {
	r := strings.NewReader("Hello world")
	f := bbs.Find(r)
	fmt.Print("reader is plain text? ", f.Valid())
	// Output: reader is plain text? false
}
