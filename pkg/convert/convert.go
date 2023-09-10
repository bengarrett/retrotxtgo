// Package convert is extends Go's x/text/encoding capability
// to convert legacy encoded text to a modern UTF-8 encoding.
package convert

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/bengarrett/retrotxtgo/pkg/byter"
	"github.com/bengarrett/retrotxtgo/pkg/fsys"
	"github.com/bengarrett/retrotxtgo/pkg/logs"
	"github.com/bengarrett/retrotxtgo/pkg/term"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
	"golang.org/x/text/transform"
)

var (
	ErrBytes     = errors.New("cannot transform an empty byte slice")
	ErrChainWrap = errors.New("wrapWidth() is a chain method that is to be" +
		" used in conjunction with swap: c.swap().wrapWidth()")
	ErrEncoding = errors.New("no encoding provided")
	ErrWidth    = errors.New("cannot determine the number columns from using line break")
)

// Convert legacy 8-bit codepage encode or Unicode byte array text to UTF-8 runes.
type Convert struct {
	Flags Flag // Flags are the cmd supplied flag values.
	Input struct {
		Encoding  encoding.Encoding // Encoding are the encoding of the input text.
		Bytes     []byte            // Bytes are the input text as bytes.
		LineBreak [2]rune           // Line break controls used by the text.
		Table     bool              // Table flags this text as a codepage table.
	}
	Output     []rune // Output are the transformed UTF-8 runes.
	Ignores    []rune // Ignores these runes.
	LineBreaks bool   // LineBreaks uses line break controls.
}

// Flag are the user supplied values.
type Flag struct {
	Controls  []string // Always use these control codes.
	SwapChars []string // Swap out these characters with UTF-8 alternatives.
	MaxWidth  int      // Maximum text width per-line.
}

// ANSI transforms legacy encoded ANSI into modern UTF-8 text.
// It displays ASCII control codes as characters.
// It obeys the DOS end of file marker.
func (c *Convert) ANSI(b ...byte) ([]rune, error) {
	c.LineBreaks = true
	c.Flags.SwapChars = nil
	c.Input.Bytes = byter.TrimEOF(b)
	if err := c.SkipCtrlCodes().Transform(); err != nil {
		return nil, fmt.Errorf("dump transform failed: %w", err)
	}
	c.Swap().ANSIControls().wrapWidth(c.Flags.MaxWidth)
	return c.Output, nil
}

// Chars transforms legacy encoded characters and text control codes into UTF-8 characters.
// It displays both ASCII and ANSI control codes as characters.
// It ignores the DOS end of file marker.
func (c *Convert) Chars(b ...byte) ([]rune, error) {
	c.Input.Table = true
	c.Input.Bytes = b
	if err := c.Transform(); err != nil {
		return nil, fmt.Errorf("chars transform failed: %w", err)
	}
	c.Swap().wrapWidth(c.Flags.MaxWidth)
	return c.Output, nil
}

// Dump transforms legacy encoded text or ANSI into modern UTF-8 text.
// It obeys common ASCII control codes.
// It ignores the DOS end of file marker.
func (c *Convert) Dump(b ...byte) ([]rune, error) {
	c.LineBreaks = true
	c.Input.Bytes = b
	if err := c.SkipCtrlCodes().Transform(); err != nil {
		return nil, fmt.Errorf("dump transform failed: %w", err)
	}
	c.Swap().ANSIControls().wrapWidth(c.Flags.MaxWidth)
	return c.Output, nil
}

// Text transforms legacy encoded text or ANSI into modern UTF-8 text.
// It obeys common ASCII control codes.
// It obeys the DOS end of file marker.
func (c *Convert) Text(b ...byte) ([]rune, error) {
	c.LineBreaks = true
	c.Input.Bytes = byter.TrimEOF(b)
	if err := c.SkipCtrlCodes().Transform(); err != nil {
		return nil, fmt.Errorf("text transform failed: %w", err)
	}
	c.Swap().ANSIControls().wrapWidth(c.Flags.MaxWidth)
	return c.Output, nil
}

// Transform byte data from named character map encoded text into UTF-8.
func (c *Convert) Transform() error {
	if c.Input.Encoding == nil {
		return ErrEncoding
	}
	if len(c.Input.Bytes) == 0 {
		return ErrBytes
	}
	// transform unicode encodings
	if r, err := unicodeDecoder(c.Input.Encoding, c.Input.Bytes); err != nil {
		return err
	} else if len(r) > 0 {
		c.Output = r
		return nil
	}
	// use the input bytes if they are already valid UTF-8 runes
	if utf8.Valid(c.Input.Bytes) {
		c.Output = bytes.Runes(c.Input.Bytes)
		return nil
	}
	// transform the input bytes into UTF-8 runes
	c.FixJISTable()
	b := bytes.Buffer{}
	t := transform.NewWriter(&b, c.Input.Encoding.NewDecoder())
	if _, err := t.Write(c.Input.Bytes); err != nil {
		return err
	}
	defer t.Close()
	c.Output = bytes.Runes(b.Bytes())
	return nil
}

// FixJISTable blanks invalid ShiftJIS characters while printing 8-bit tables.
func (c *Convert) FixJISTable() {
	if c.Input.Encoding == japanese.ShiftJIS && c.Input.Table {
		// this is only for the table command,
		// it will break normal shift-jis encode text
		for i, b := range c.Input.Bytes {
			switch {
			case b > 0x7f && b <= 0xa0,
				b >= 0xe0 && b <= 0xff:
				c.Input.Bytes[i] = SP
			}
		}
	}
}

// decode transforms encoded input bytes into UTF-8 runes.
func decode(e encoding.Encoding, input []byte) ([]rune, error) {
	b, err := e.NewDecoder().Bytes(input)
	if err != nil {
		return nil, err
	}
	// c.Output
	return bytes.Runes(b), nil
}

// unicodeDecoder transforms UTF-8, UTF-16 or UTF-32 bytes into UTF-8 runes.
func unicodeDecoder(e encoding.Encoding, input []byte) ([]rune, error) {
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
	switch e {
	case unicode.UTF8, unicode.UTF8BOM:
		return bytes.Runes(input), nil
	case u16be:
		return decode(u16be, input)
	case u16le:
		return decode(u16le, input)
	case u16beB:
		return decode(u16beB, input)
	case u16leB:
		return decode(u16leB, input)
	case u32be:
		return decode(u32be, input)
	case u32beB:
		return decode(u32beB, input)
	case u32le:
		return decode(u32le, input)
	case u32leB:
		return decode(u32leB, input)
	}
	return nil, nil
}

func replaceNL(r []rune) []rune {
	re := regexp.MustCompile(`\r?\n`)
	s := re.ReplaceAllString(string(r), "")
	return []rune(s)
}

// wrapWidth enforces a row length by inserting newline characters.
func (c *Convert) wrapWidth(max int) {
	if max < 1 {
		return
	}
	// remove newlines
	c.Output = replaceNL(c.Output)
	cnt := len(c.Output)
	if cnt == 0 {
		log.Fatal(ErrChainWrap)
	}
	r := strings.NewReader(string(c.Output))
	cols, err := fsys.Columns(r, c.Input.LineBreak)
	if err != nil {
		logs.FatalMark(fmt.Sprint(c.Input.LineBreak), ErrWidth, err)
	}
	if cols <= max {
		return
	}
	limit := math.Ceil(float64(cnt) / float64(max))
	var w bytes.Buffer
	for f := float64(1); f <= limit; f++ {
		switch f {
		case 1:
			fmt.Fprintf(&w, "%s\n", string(c.Output[0:max]))
		default:
			i := int(f)
			a, b := (i-1)*max, i*max
			if b >= cnt {
				fmt.Fprintf(&w, "%s\n", string(c.Output[a:cnt]))
				continue
			}
			fmt.Fprintf(&w, "%s\n", string(c.Output[a:b]))
		}
	}
	c.Output = []rune(w.String())
}

// SkipCtrlCodes marks control characters to be ignored.
// It needs to be applied before Convert.transform().
func (c *Convert) SkipCtrlCodes() *Convert {
	unknown := []string{}
	for _, v := range c.Flags.Controls {
		switch strings.ToLower(v) {
		case "eof", "=":
			continue
		case "tab", "ht", "t":
			c.ignore(HT)
		case "bell", "bel", "b":
			c.ignore(BEL)
		case "cr", "c":
			c.ignore(CR)
		case "lf", "l":
			c.ignore(LF)
		case "backspace", "bs":
			c.ignore(BS)
		case "del", "d":
			c.ignore(DEL)
		case "esc", "e":
			c.ignore(ESC)
		case "formfeed", "ff", "f":
			c.ignore(FF)
		case "vtab", "vt", "v":
			c.ignore(VT)
		default:
			unknown = append(unknown, v)
		}
	}
	if len(unknown) > 0 {
		fmt.Fprintln(os.Stderr, term.Inform(), "unsupported --control values:", strings.Join(unknown, ","))
	}
	return c
}

// ignore adds the rune to an ignore runes list.
func (c *Convert) ignore(r rune) {
	c.Ignores = append(c.Ignores, r)
}
