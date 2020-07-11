package convert

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/unicode"
	"retrotxt.com/retrotxt/lib/filesystem"
	"retrotxt.com/retrotxt/lib/logs"
)

// Convert for the transformation of legacy encoded text to UTF-8.
type Convert struct {
	Source    []byte            // Source legacy encoded text.
	Runes     []rune            // Runes with UTF-8 text.
	encode    encoding.Encoding // Source character set encoding
	len       int               // Runes count
	newline   bool              // use newline controls
	newlines  [2]rune           // the newline controls rune values
	ignores   []rune            // these runes will not be transformed
	swapChars []int
	table     bool
}

// Args are user supplied flag values
type Args struct {
	Controls []string
	Encoding string
	Swap     []int
	Width    int
}

// Chars transforms legacy encoded characters and text control codes into UTF-8 characters.
func (a Args) Chars(b *[]byte) (utf8 []rune, err error) {
	var c = Convert{
		Source:    *b,
		swapChars: a.Swap,
		table:     true,
	}
	if _, err = c.Transform(a.Encoding); err != nil {
		return nil, fmt.Errorf("chars transform failed: %s", err)
	}
	c.Swap()
	c.width(a.Width)
	return c.Runes, nil
}

// Dump transforms legacy encoded text or ANSI into modern UTF-8 text
// including all text contained after any MS-DOS end-of-file markers.
func (a Args) Dump(b *[]byte) (utf8 []rune, err error) {
	var c = Convert{
		Source:    *b,
		newline:   true,
		swapChars: a.Swap,
	}
	c.controls(a)
	if _, err = c.Transform(a.Encoding); err != nil {
		return nil, fmt.Errorf("dump transform failed: %s", err)
	}
	c.Swap().ANSI()
	c.width(a.Width)
	return c.Runes, nil
}

// Text transforms legacy encoded text or ANSI into modern UTF-8 text.
func (a Args) Text(b *[]byte) (utf8 []rune, err error) {
	var c = Convert{
		Source:    *b,
		newline:   true,
		swapChars: a.Swap,
	}
	c.controls(a)
	c.Source = EndOfFile(*b)
	if _, err = c.Transform(a.Encoding); err != nil {
		return nil, fmt.Errorf("text transform failed: %s", err)
	}
	c.Swap().ANSI()
	c.width(a.Width)
	return c.Runes, nil
}

// Transform byte data from named character map encoded text into UTF-8.
func (c *Convert) Transform(name string) (*Convert, error) {
	if name == "" {
		name = "UTF-8"
	}
	var err error
	if c.encode, err = Encoding(name); err != nil {
		return c, fmt.Errorf("transform encoding error: %s", err)
	}
	if len(c.Source) == 0 {
		return c, nil
	}
	// don't transform unicode, japanese..
	switch c.encode {
	case unicode.UTF8, unicode.UTF8BOM:
		for _, b := range c.Source {
			c.Runes = append(c.Runes, rune(b))
		}
		c.len = len(c.Runes)
		return c, nil
	}
	// blank invalid shiftjis characters when printing 8-bit tables
	switch c.encode {
	case japanese.ShiftJIS:
		if !c.table {
			break
		}
		// this will break normal shift-jis encoding
		for i, b := range c.Source {
			switch {
			case b > 0x7f && b <= 0xa0,
				b >= 0xe0 && b <= 0xff:
				c.Source[i] = 32
			}
		}
	}
	// transform source if it is not already UTF-8
	if utf8.Valid(c.Source) {
		c.Runes = bytes.Runes(c.Source)
		c.len = len(c.Runes)
		return c, nil
	}
	if c.Source, err = c.encode.NewDecoder().Bytes(c.Source); err != nil {
		return c, fmt.Errorf("transform new decoder error: %s", err)
	}
	c.Runes = bytes.Runes(c.Source)
	c.len = len(c.Runes)
	return c, nil
}

func (c *Convert) width(max int) {
	if max < 1 {
		return
	}
	cnt := len(c.Runes)
	cols, err := filesystem.Columns(bytes.NewReader(c.Source), c.newlines)
	if err != nil {
		logs.Println("ignoring width argument", "",
			fmt.Errorf("width could not determine the columns: %s", err))
		return
	}
	if cols <= max {
		return
	}
	limit := math.Ceil(float64(cnt) / float64(max))
	var w bytes.Buffer
	for f := float64(1); f <= limit; f++ {
		switch f {
		case 1:
			fmt.Fprintf(&w, "%s\n", string(c.Runes[0:max]))
		default:
			i := int(f)
			a, b := (i-1)*max, i*max
			if b >= cnt {
				fmt.Fprintf(&w, "%s\n", string(c.Runes[a:cnt]))
			} else {
				fmt.Fprintf(&w, "%s\n", string(c.Runes[a:b]))
			}
		}
	}
	c.Runes = []rune(w.String())
	return
}

func (c *Convert) controls(a Args) {
	for _, v := range a.Controls {
		switch strings.ToLower(v) {
		case "bell", "b":
			c.ignore(7)
		case "backspace", "bs":
			c.ignore(8)
		case "tab", "ht", "t":
			c.ignore(9)
		case "lf", "l":
			c.ignore(10)
		case "vtab", "vt", "v":
			c.ignore(11)
		case "formfeed", "ff", "f":
			c.ignore(12)
		case "cr", "c":
			c.ignore(13)
		case "esc", "e":
			c.ignore(27)
		case "del", "d":
			c.ignore(127)
		}
	}
}
