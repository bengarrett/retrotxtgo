package convert

import (
	"bytes"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding"
)

// Data for the transformation of legacy encoded text to UTF-8.
type Data struct {
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
	var d = Data{
		Source: *b,
	}
	if _, err = d.Transform(a.Encoding); err != nil {
		return nil, err
	}
	d.Swap()
	return d.Runes, nil
}

// Dump transforms legacy encoded text or ANSI into modern UTF-8 text
// including all text contained after any MS-DOS end-of-file markers.
func (a Args) Dump(b *[]byte) (utf8 []rune, err error) {
	var d = Data{
		Source:  *b,
		newline: true,
	}
	d.controls(a)
	if _, err = d.Transform(a.Encoding); err != nil {
		return nil, err
	}
	d.Swap().ANSI()
	return d.Runes, nil
}

// Text transforms legacy encoded text or ANSI into modern UTF-8 text.
func (a Args) Text(b *[]byte) (utf8 []rune, err error) {
	var d = Data{
		Source:  *b,
		newline: true,
	}
	d.controls(a)
	d.Source = EndOfFile(*b)
	if _, err = d.Transform(a.Encoding); err != nil {
		return nil, err
	}
	d.Swap().ANSI()
	return d.Runes, nil
}

// Transform byte data from named character map encoded text into UTF-8.
func (d *Data) Transform(name string) (*Data, error) {
	if name == "" {
		name = "UTF-8"
	}
	var err error
	if d.encode, err = Encoding(name); err != nil {
		return d, err
	}
	if len(d.Source) == 0 {
		return d, nil
	}
	// only transform source if it is not already UTF-8
	if utf8.Valid(d.Source) {
		d.Runes = bytes.Runes(d.Source)
		d.len = len(d.Runes)
		return d, nil
	}
	if d.Source, err = d.encode.NewDecoder().Bytes(d.Source); err != nil {
		return d, err
	}
	d.Runes = bytes.Runes(d.Source)
	d.len = len(d.Runes)
	return d, nil
}

func (d *Data) controls(a Args) {
	for _, v := range a.Controls {
		switch strings.ToLower(v) {
		case "bell", "b":
			d.ignore(7)
		case "backspace", "bs":
			d.ignore(8)
		case "tab", "ht", "t":
			d.ignore(9)
		case "lf", "l":
			d.ignore(10)
		case "vtab", "vt", "v":
			d.ignore(11)
		case "formfeed", "ff", "f":
			d.ignore(12)
		case "cr", "c":
			d.ignore(13)
		case "esc", "e":
			d.ignore(27)
		case "del", "d":
			d.ignore(127)
		}
	}
}
