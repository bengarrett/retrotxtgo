// Package convert extends the interface for the character encodings
// that transform text to and from Unicode UTF-8.
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
	ErrANSI   = errors.New("ansi controls must be chained to c.swap")
	ErrBytes  = errors.New("cannot transform an empty byte slice")
	ErrEncode = errors.New("no input encoding provided")
	ErrName   = errors.New("unknown or unsupported code page name or alias")
	ErrOutput = errors.New("nothing to output")
	ErrWidth  = errors.New("cannot find the number columns from using line break")
	ErrWrap   = errors.New("wrap width must be chained to c.swap")
)

// Convert 8-bit code page text encodings or Unicode byte array text to UTF-8 runes.
type Convert struct {
	Args  Flag // Args are the cmd supplied flag arguments.
	Input struct {
		Encoding  encoding.Encoding // Encoding are the encoding of the input text.
		Input     []byte            // Bytes are the input text as bytes.
		Ignore    []rune            // Ignore these runes.
		LineBreak [2]rune           // Line break controls used by the text.
		UseBreaks bool              // UseBreaks uses the line break controls as new lines.
		Table     bool              // Table flags this text as a code page table.
	}
	Output []rune // Output are the transformed UTF-8 runes.
}

// Flag are the user supplied values.
type Flag struct {
	Controls  []string // Always use these control codes.
	SwapChars []string // Swap out these characters with common alternatives.
	MaxWidth  int      // Maximum text width per-line.
}

// ANSI transforms legacy encoded ANSI into modern UTF-8 text.
// It displays ASCII control codes as characters.
// It obeys the DOS end of file marker.
func (c *Convert) ANSI(b ...byte) ([]rune, error) {
	c.Input.UseBreaks = true
	c.Args.SwapChars = nil
	c.Input.Input = byter.TrimEOF(b)
	if err := c.SkipCode().Transform(); err != nil {
		return nil, fmt.Errorf("dump transform failed: %w", err)
	}
	c, err := c.Swap()
	if err != nil {
		return nil, err
	}
	c.ANSIControls().wrapWidth(c.Args.MaxWidth)
	return c.Output, nil
}

// Chars transforms legacy encoded characters and text control codes into UTF-8 characters.
// It displays both ASCII and ANSI control codes as characters.
// It ignores the DOS end of file marker.
func (c *Convert) Chars(b ...byte) ([]rune, error) {
	c.Input.Table = true
	c.Input.Input = b
	if err := c.Transform(); err != nil {
		return nil, fmt.Errorf("chars transform failed: %w", err)
	}
	c, err := c.Swap()
	if err != nil {
		return nil, err
	}
	c.wrapWidth(c.Args.MaxWidth)
	return c.Output, nil
}

// Dump transforms legacy encoded text or ANSI into modern UTF-8 text.
// It obeys common ASCII control codes.
// It ignores the DOS end of file marker.
func (c *Convert) Dump(b ...byte) ([]rune, error) {
	c.Input.UseBreaks = true
	c.Input.Input = b
	if err := c.SkipCode().Transform(); err != nil {
		return nil, fmt.Errorf("dump transform failed: %w", err)
	}
	c, err := c.Swap()
	if err != nil {
		return nil, err
	}
	c.ANSIControls().wrapWidth(c.Args.MaxWidth)
	return c.Output, nil
}

// Text transforms legacy encoded text or ANSI into modern UTF-8 text.
// It obeys common ASCII control codes.
// It obeys the DOS end of file marker.
func (c *Convert) Text(b ...byte) ([]rune, error) {
	c.Input.UseBreaks = true
	c.Input.Input = byter.TrimEOF(b)
	if err := c.SkipCode().Transform(); err != nil {
		return nil, fmt.Errorf("text transform failed: %w", err)
	}
	c, err := c.Swap()
	if err != nil {
		return nil, err
	}
	c.ANSIControls().wrapWidth(c.Args.MaxWidth)
	return c.Output, nil
}

// Transform byte data from named character map encoded text into UTF-8.
func (c *Convert) Transform() error {
	if c.Input.Encoding == nil {
		return ErrEncode
	}
	if len(c.Input.Input) == 0 {
		return nil
	}
	// transform unicode encodings
	if r, err := unicodeDecoder(c.Input.Encoding, c.Input.Input...); err != nil {
		return err
	} else if len(r) > 0 {
		c.Output = r
		return nil
	}
	// use the input bytes if they are already valid UTF-8 runes
	if utf8.Valid(c.Input.Input) {
		c.Output = bytes.Runes(c.Input.Input)
		return nil
	}
	// transform the input bytes into UTF-8 runes
	c.FixJISTable()
	b := &bytes.Buffer{}
	t := transform.NewWriter(b, c.Input.Encoding.NewDecoder())
	defer t.Close()
	if _, err := t.Write(c.Input.Input); err != nil {
		return err
	}
	c.Output = bytes.Runes(b.Bytes())
	return nil
}

// FixJISTable blanks invalid ShiftJIS characters while printing 8-bit tables.
func (c *Convert) FixJISTable() {
	if c.Input.Encoding == japanese.ShiftJIS && c.Input.Table {
		// this is only for the table command,
		// it will break normal shift-jis encode text
		for i, b := range c.Input.Input {
			switch {
			case b > 0x7f && b <= 0xa0,
				b >= 0xe0 && b <= 0xff:
				c.Input.Input[i] = SP
			}
		}
	}
}

// decode transforms encoded bytes into UTF-8 runes.
func decode(e encoding.Encoding, b ...byte) ([]rune, error) {
	p, err := e.NewDecoder().Bytes(b)
	if err != nil {
		return nil, err
	}
	return bytes.Runes(p), nil
}

// unicodeDecoder transforms UTF-8, UTF-16 or UTF-32 bytes into UTF-8 runes.
func unicodeDecoder(e encoding.Encoding, b ...byte) ([]rune, error) {
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
		return bytes.Runes(b), nil
	case u16be:
		return decode(u16be, b...)
	case u16le:
		return decode(u16le, b...)
	case u16beB:
		return decode(u16beB, b...)
	case u16leB:
		return decode(u16leB, b...)
	case u32be:
		return decode(u32be, b...)
	case u32beB:
		return decode(u32beB, b...)
	case u32le:
		return decode(u32le, b...)
	case u32leB:
		return decode(u32leB, b...)
	}
	return nil, nil
}

// replaceNL replaces newlines with single spaces.
func replaceNL(r ...rune) []rune {
	re := regexp.MustCompile(`\r?\n`)
	const space = " "
	s := re.ReplaceAllString(string(r), space)
	return []rune(s)
}

// wrapWidth enforces a row length by inserting newline characters.
// Any tab characters are replaced with three spaces.
func (c *Convert) wrapWidth(max int) {
	if max < 1 {
		return
	}
	// remove newlines
	c.Output = replaceNL(c.Output...)
	// replace tabs with three spaces
	s := string(c.Output)
	const threeSpaces = "   "
	s = strings.ReplaceAll(s, "\t", threeSpaces)
	c.Output = []rune(s)
	cnt := len(c.Output)
	if cnt == 0 {
		log.Fatal(ErrWrap)
	}
	r := strings.NewReader(string(c.Output))
	cols, err := fsys.Columns(r, c.Input.LineBreak)
	if err != nil {
		logs.FatalS(ErrWidth, err, fmt.Sprint(c.Input.LineBreak))
	}
	if cols <= max {
		return
	}
	limit := math.Ceil(float64(cnt) / float64(max))
	b := &bytes.Buffer{}
	for f := float64(1); f <= limit; f++ {
		switch f {
		case 1:
			fmt.Fprintf(b, "%s\n", string(c.Output[0:max]))
		default:
			i := int(f)
			x, y := (i-1)*max, i*max
			if y >= cnt {
				fmt.Fprintf(b, "%s\n", string(c.Output[x:cnt]))
				continue
			}
			fmt.Fprintf(b, "%s\n", string(c.Output[x:y]))
		}
	}
	c.Output = []rune(b.String())
}

// SkipCode marks control characters to be ignored.
// It needs to be applied before Convert.transform().
func (c *Convert) SkipCode() *Convert {
	unknown := []string{}
	for _, v := range c.Args.Controls {
		v = strings.TrimSpace(v)
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
	if len(unknown) > 1 {
		fmt.Fprintln(os.Stderr, term.Inform(),
			"unsupported control values:", strings.Join(unknown, ","))
	}
	return c
}

// ignore adds the rune to an ignore runes list.
func (c *Convert) ignore(r rune) {
	c.Input.Ignore = append(c.Input.Ignore, r)
}
