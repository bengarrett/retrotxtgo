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
	encode    encoding.Encoding // Source character set encoding
	len       int               // Runes count
	newline   bool              // use newline controls
	table     bool
	newlines  [2]rune // the newline controls rune values
	Source    []byte  // Source legacy encoded text.
	ignores   []rune  // these runes will not be transformed
	Runes     []rune  // Runes with UTF-8 text.
	swapChars []int
}

// Args are user supplied flag values.
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
	if err = c.Transform(a.Encoding); err != nil {
		return nil, fmt.Errorf("chars transform failed: %w", err)
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
	if err = c.Transform(a.Encoding); err != nil {
		return nil, fmt.Errorf("dump transform failed: %w", err)
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
	c.Source = EndOfFile(*b...)
	if err = c.Transform(a.Encoding); err != nil {
		return nil, fmt.Errorf("text transform failed: %w", err)
	}
	c.Swap().ANSI()
	c.width(a.Width)
	return c.Runes, nil
}

// Transform byte data from named character map encoded text into UTF-8.
func (c *Convert) Transform(name string) error {
	if name == "" {
		name = "UTF-8"
	}
	var err error
	if c.encode, err = Encoding(name); err != nil {
		return fmt.Errorf("transform encoding error: %w", err)
	}
	if len(c.Source) == 0 {
		return nil
	}
	// don't transform unicode encoded strings
	switch c.encode {
	case unicode.UTF8, unicode.UTF8BOM:
		for _, b := range c.Source {
			c.Runes = append(c.Runes, rune(b))
		}
		c.len = len(c.Runes)
		return nil
	}
	// blank invalid shiftjis characters when printing 8-bit tables
	switch c.encode {
	case japanese.ShiftJIS:
		if !c.table {
			break
		}
		// this is only for the table command,
		// it will break normal shift-jis encode text
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
		return nil
	}
	if c.Source, err = c.encode.NewDecoder().Bytes(c.Source); err != nil {
		return fmt.Errorf("transform new decoder error: %w", err)
	}
	c.Runes = bytes.Runes(c.Source)
	c.len = len(c.Runes)
	return nil
}

func (c *Convert) width(max int) {
	if max < 1 {
		return
	}
	cnt := len(c.Runes)
	cols, err := filesystem.Columns(bytes.NewReader(c.Source), c.newlines)
	if err != nil {
		logs.Println("ignoring width argument", "",
			fmt.Errorf("width could not determine the columns: %w", err))
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
}

func (c *Convert) controls(a Args) {
	const (
		esc = 27
		del = 127
	)
	const (
		bell = iota + 7
		bs
		tab
		lf
		vt
		ff
		cr
	)
	for _, v := range a.Controls {
		switch strings.ToLower(v) {
		case "bell", "b":
			c.ignore(bell)
		case "backspace", "bs":
			c.ignore(bs)
		case "tab", "ht", "t":
			c.ignore(tab)
		case "lf", "l":
			c.ignore(lf)
		case "vtab", "vt", "v":
			c.ignore(vt)
		case "formfeed", "ff", "f":
			c.ignore(ff)
		case "cr", "c":
			c.ignore(cr)
		case "esc", "e":
			c.ignore(esc)
		case "del", "d":
			c.ignore(del)
		}
	}
}
