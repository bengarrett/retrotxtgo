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
	"golang.org/x/text/encoding/unicode/utf32"
	"retrotxt.com/retrotxt/lib/filesystem"
	"retrotxt.com/retrotxt/lib/logs"
)

// ANSI transforms legacy encoded ANSI into modern UTF-8 text.
// It displays ASCII control codes as characters.
// It obeys the end of file marker.
func (c *Convert) ANSI(b *[]byte) (utf []rune, err error) {
	c.Output.lineBreaks = true
	c.Flags.SwapChars = nil
	c.Source.B = *b
	c.unicodeControls()
	c.Source.B = EndOfFile(*b...)
	if err = c.Transform(); err != nil {
		return nil, fmt.Errorf("dump transform failed: %w", err)
	}
	c.Swap().ANSIControls()
	c.width(c.Flags.Width)
	return c.Output.R, nil
}

// Chars transforms legacy encoded characters and text control codes into UTF-8 characters.
// It displays both ASCII and ANSI control codes as characters.
// It ignores the end of file marker.
func (c *Convert) Chars(b *[]byte) (utf []rune, err error) {
	c.Source.table = true
	c.Source.B = *b
	if err = c.Transform(); err != nil {
		return nil, fmt.Errorf("chars transform failed: %w", err)
	}
	c.Swap()
	c.width(c.Flags.Width)
	return c.Output.R, nil
}

// Dump transforms legacy encoded text or ANSI into modern UTF-8 text.
// It obeys common ASCII control codes.
// It ignores the end of file marker.
func (c *Convert) Dump(b *[]byte) (utf []rune, err error) {
	c.Output.lineBreaks = true
	c.Source.B = *b
	c.unicodeControls()
	if err = c.Transform(); err != nil {
		return nil, fmt.Errorf("dump transform failed: %w", err)
	}
	c.Swap().ANSIControls()
	c.width(c.Flags.Width)
	return c.Output.R, nil
}

// Text transforms legacy encoded text or ANSI into modern UTF-8 text.
// It obeys common ASCII control codes.
// It obeys the end of file marker.
func (c *Convert) Text(b *[]byte) (utf []rune, err error) {
	c.Output.lineBreaks = true
	c.Source.B = *b
	c.unicodeControls()
	c.Source.B = EndOfFile(*b...)
	if err = c.Transform(); err != nil {
		return nil, fmt.Errorf("text transform failed: %w", err)
	}
	c.Swap().ANSIControls()
	c.width(c.Flags.Width)
	return c.Output.R, nil
}

// Transform byte data from named character map encoded text into UTF-8.
func (c *Convert) Transform() error {
	if c.Source.E == nil {
		c.Source.E = unicode.UTF8
	}
	if len(c.Source.B) == 0 {
		return nil
	}
	if err := c.transformUnicode(); err != nil {
		return err
	} else if c.Output.len > 0 {
		return nil
	}
	c.transformFixJISTable()
	// return the source as runes if it is already in UTF-8 Unicode
	if utf8.Valid(c.Source.B) {
		c.Output.R = bytes.Runes(c.Source.B)
		c.Output.len = len(c.Output.R)
		return nil
	}
	var err error
	if c.Source.B, err = c.Source.E.NewDecoder().Bytes(c.Source.B); err != nil {
		return fmt.Errorf("transform new decoder error: %w", err)
	}
	c.Output.R = bytes.Runes(c.Source.B)
	c.Output.len = len(c.Output.R)
	return nil
}

// blank invalid shiftjis characters when printing 8-bit tables.
func (c *Convert) transformFixJISTable() {
	if c.Source.E == japanese.ShiftJIS && c.Source.table {
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
}

// transform Unicode-16 or Unicode-32 text into UTF-8 encoded Unicode.
func (c *Convert) transformUnicode() error {
	var decode = func(e encoding.Encoding) error {
		result, err := e.NewDecoder().Bytes(c.Source.B)
		if err != nil {
			return err
		}
		c.Output.R = bytes.Runes(result)
		c.Output.len = len(c.Output.R)
		return nil
	}
	var (
		u16be  = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
		u16beB = unicode.UTF16(unicode.BigEndian, unicode.ExpectBOM)
		u16le  = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
		u16leB = unicode.UTF16(unicode.LittleEndian, unicode.ExpectBOM)
		u32be  = utf32.UTF32(utf32.LittleEndian, utf32.IgnoreBOM)
		u32beB = utf32.UTF32(utf32.BigEndian, utf32.UseBOM)
		u32le  = utf32.UTF32(utf32.LittleEndian, utf32.IgnoreBOM)
		u32leB = utf32.UTF32(utf32.LittleEndian, utf32.UseBOM)
	)
	switch c.Source.E {
	case unicode.UTF8, unicode.UTF8BOM:
		c.Output.R = bytes.Runes(c.Source.B)
		c.Output.len = len(c.Output.R)
	case u16be:
		if err := decode(u16be); err != nil {
			return err
		}
	case u16le:
		if err := decode(u16le); err != nil {
			return err
		}
	case u16beB:
		if err := decode(u16beB); err != nil {
			return err
		}
	case u16leB:
		if err := decode(u16leB); err != nil {
			return err
		}
	case u32be:
		if err := decode(u32be); err != nil {
			return err
		}
	case u32beB:
		if err := decode(u32beB); err != nil {
			return err
		}
	case u32le:
		if err := decode(u32le); err != nil {
			return err
		}
	case u32leB:
		if err := decode(u32leB); err != nil {
			return err
		}
	}
	return nil
}

func (c *Convert) width(max int) {
	if max < 1 {
		return
	}
	cnt := len(c.Output.R)
	cols, err := filesystem.Columns(bytes.NewReader(c.Source.B), c.Source.lineBreak)
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
			fmt.Fprintf(&w, "%s\n", string(c.Output.R[0:max]))
		default:
			i := int(f)
			a, b := (i-1)*max, i*max
			if b >= cnt {
				fmt.Fprintf(&w, "%s\n", string(c.Output.R[a:cnt]))
			} else {
				fmt.Fprintf(&w, "%s\n", string(c.Output.R[a:b]))
			}
		}
	}
	c.Output.R = []rune(w.String())
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
	c.Output.ignores = append(c.Output.ignores, r)
}
