package convert

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"retrotxt.com/retrotxt/lib/filesystem"
	"retrotxt.com/retrotxt/lib/logs"
)

// Convert for the transformation of legacy encoded text to UTF-8.
type Convert struct {
	Source   []byte            // Source legacy encoded text.
	Runes    []rune            // Runes with UTF-8 text.
	encode   encoding.Encoding // Source character set encoding
	len      int               // Runes count
	newline  bool              // use newline controls
	newlines [2]rune           // the newline controls rune values
	ignores  []rune            // these runes will not be transformed
}

// Args are user supplied flag values
type Args struct {
	Controls []string
	Encoding string
	Width    int
}

// Chars transforms legacy encoded characters and text control codes into UTF-8 characters.
func (a Args) Chars(b *[]byte) (utf8 []rune, err error) {
	var c = Convert{
		Source: *b,
	}
	if _, err = c.Transform(a.Encoding); err != nil {
		return nil, err
	}
	c.Swap()
	return c.Runes, nil
}

// Dump transforms legacy encoded text or ANSI into modern UTF-8 text
// including all text contained after any MS-DOS end-of-file markers.
func (a Args) Dump(b *[]byte) (utf8 []rune, err error) {
	var c = Convert{
		Source:  *b,
		newline: true,
	}
	c.controls(a)
	if _, err = c.Transform(a.Encoding); err != nil {
		return nil, err
	}
	c.Swap().ANSI()
	return c.Runes, nil
}

// Text transforms legacy encoded text or ANSI into modern UTF-8 text.
func (a Args) Text(b *[]byte) (utf8 []rune, err error) {
	var c = Convert{
		Source:  *b,
		newline: true,
	}
	c.controls(a)
	c.Source = EndOfFile(*b)
	if _, err = c.Transform(a.Encoding); err != nil {
		return nil, err
	}
	c.Swap().ANSI()
	if a.Width > 0 {
		err := c.width(a.Width)
		if err != nil {
			logs.Println("ignoring width argument", "", err)
		}
	}
	return c.Runes, nil
}

// Transform byte data from named character map encoded text into UTF-8.
func (c *Convert) Transform(name string) (*Convert, error) {
	if name == "" {
		name = "UTF-8"
	}
	var err error
	if c.encode, err = Encoding(name); err != nil {
		return c, err
	}
	if len(c.Source) == 0 {
		return c, nil
	}
	// only transform source if it is not already UTF-8
	if utf8.Valid(c.Source) {
		c.Runes = bytes.Runes(c.Source)
		c.len = len(c.Runes)
		return c, nil
	}
	if c.Source, err = c.encode.NewDecoder().Bytes(c.Source); err != nil {
		return c, err
	}
	c.Runes = bytes.Runes(c.Source)
	c.len = len(c.Runes)
	return c, nil
}

func (c *Convert) width(max int) error {
	cnt := len(c.Runes)
	cols, err := filesystem.Columns(bytes.NewReader(c.Source), c.newlines)
	if err != nil {
		return fmt.Errorf("width could not determine the columns: %s", err)
	}
	if cols <= max {
		return nil
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
	return nil
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
