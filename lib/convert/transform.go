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

// Convert 8-bit legacy or other Unicode text to UTF-8.
type Convert struct {
	// Source text for conversion.
	Source struct {
		B []byte              // Text as bytes
		E []encoding.Encoding // Encoding
	}
	// User supplied values.
	Flags Flags
	Runes []rune // Runes with UTF-8 text
	//Source    []byte            // Source legacy encoded text
	encode    encoding.Encoding // Source character set encoding
	ignores   []rune            // these runes will not be transformed
	len       int               // Runes count
	lineBreak [2]rune           // line break controls
	swapChars []int
	table     bool
	useBreaks bool // use line break controls

}

// Flags are the user supplied values.
type Flags struct {
	Controls  []string
	Encoding  encoding.Encoding
	SwapChars []int
	Width     int
}

// ANSI transforms legacy encoded ANSI into modern UTF-8 text.
// It displays ASCII control codes as characters.
// It obeys the end of file marker.
func (c Convert) ANSI(b *[]byte) (utf []rune, err error) {
	c.useBreaks = true
	c.swapChars = nil
	c.Source.B = *b
	c.unicodeControls()
	c.Source.B = EndOfFile(*b...)
	if err = c.Transform(c.Flags.Encoding); err != nil {
		return nil, fmt.Errorf("dump transform failed: %w", err)
	}
	c.Swap().ANSIControls()
	c.width(c.Flags.Width)
	return c.Runes, nil
}

// Chars transforms legacy encoded characters and text control codes into UTF-8 characters.
// It displays both ASCII and ANSI control codes as characters.
// It ignores the end of file marker.
func (c Convert) Chars(b *[]byte) (utf []rune, err error) {
	c.swapChars = c.Flags.SwapChars
	c.table = true
	c.Source.B = *b
	if err = c.Transform(c.Flags.Encoding); err != nil {
		return nil, fmt.Errorf("chars transform failed: %w", err)
	}
	c.Swap()
	c.width(c.Flags.Width)
	return c.Runes, nil
}

// Dump transforms legacy encoded text or ANSI into modern UTF-8 text.
// It obeys common ASCII control codes.
// It ignores the end of file marker.
func (c Convert) Dump(b *[]byte) (utf []rune, err error) {
	c.useBreaks = true
	c.swapChars = c.Flags.SwapChars
	c.Source.B = *b
	c.unicodeControls()
	if err = c.Transform(c.Flags.Encoding); err != nil {
		return nil, fmt.Errorf("dump transform failed: %w", err)
	}
	c.Swap().ANSIControls()
	c.width(c.Flags.Width)
	return c.Runes, nil
}

// Text transforms legacy encoded text or ANSI into modern UTF-8 text.
// It obeys common ASCII control codes.
// It obeys the end of file marker.
func (c Convert) Text(b *[]byte) (utf []rune, err error) {
	c.useBreaks = true
	c.swapChars = c.Flags.SwapChars //todo: redundant?
	c.Source.B = *b
	c.unicodeControls()
	c.Source.B = EndOfFile(*b...)
	if err = c.Transform(c.Flags.Encoding); err != nil {
		return nil, fmt.Errorf("text transform failed: %w", err)
	}
	c.Swap().ANSIControls()
	c.width(c.Flags.Width)
	return c.Runes, nil
}

// Transform byte data from named character map encoded text into UTF-8.
func (c *Convert) Transform(from encoding.Encoding) error {
	if from == nil {
		from = unicode.UTF8
	}
	c.encode = from // TODO: check if needed
	var err error
	if len(c.Source.B) == 0 {
		return nil
	}
	// don't transform, instead copy unicode encoded strings
	switch c.encode {
	case unicode.UTF8, unicode.UTF8BOM:
		c.Runes = []rune(string(c.Source.B))
		c.len = len(c.Runes)
		return nil
	}
	// blank invalid shiftjis characters when printing 8-bit tables
	if c.encode == japanese.ShiftJIS && c.table {
		// this is only for the table command,
		// it will break normal shift-jis encode text
		for i, b := range c.Source.B {
			switch {
			case b > 0x7f && b <= 0xa0,
				b >= 0xe0 && b <= 0xff:
				c.Source.B[i] = 32
			}
		}
	}
	// transform source if it is not already UTF-8
	if utf8.Valid(c.Source.B) {
		c.Runes = bytes.Runes(c.Source.B)
		c.len = len(c.Runes)
		return nil
	}
	if c.Source.B, err = c.encode.NewDecoder().Bytes(c.Source.B); err != nil {
		return fmt.Errorf("transform new decoder error: %w", err)
	}
	c.Runes = bytes.Runes(c.Source.B)
	c.len = len(c.Runes)
	return nil
}

func (c *Convert) width(max int) {
	if max < 1 {
		return
	}
	cnt := len(c.Runes)
	cols, err := filesystem.Columns(bytes.NewReader(c.Source.B), c.lineBreak)
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

func (c *Convert) unicodeControls() {
	const (
		bell = iota + 7 // BEL = x07
		bs
		tab
		lf
		vt
		ff
		cr

		esc = 27
		del = 127
	)
	for _, v := range c.Flags.Controls {
		switch strings.ToLower(v) {
		case "bell", "bel", "b":
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

func (c *Convert) ignore(r rune) {
	c.ignores = append(c.ignores, r)
}
