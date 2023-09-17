// Package sample opens and encodes the example embedded text files.
// These files are used for demostrating the info and the view commands.
package sample

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/bengarrett/retrotxtgo/pkg/convert"
	"github.com/bengarrett/retrotxtgo/static"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
)

var (
	ErrEncode  = errors.New("no encoding provided")
	ErrConvert = errors.New("unknown convert method")
	ErrConvNil = errors.New("conv argument cannot be empty")
	ErrName    = errors.New("sample filename does not exist")
)

// Flags and configuration values by the user.
type Flags struct {
	Input  encoding.Encoding // Input encoding is set using the --from flag.
	Output encoding.Encoding // Output encoding is set using the --to flag.
}

// Sample textfile data.
type Sample struct {
	// the order of these fields must not be changed
	Convert     Output            // Convert text method.
	Encoding    encoding.Encoding // Encoding used by the sample.
	Name        string            // Name of the sample.
	Description string            // Description of the sample.
}

// Output method.
type Output int

const (
	Ansi Output = iota // Ansi method should only be use with ANSI text.
	Ctrl               // Ctrl method print the common text control codes as characters.
	Text               // Text obeys the common text controls.
	Dump               // Dump obeys the common text controls except for EOF, end-of-file.
)

// Map is the collection of sample text files.
// Each sample includes the output method, character encoding, the filename and a brief description.
func Map() map[string]Sample {
	var (
		cp037  = charmap.CodePage037
		cp437  = charmap.CodePage437
		cp865  = charmap.CodePage865 // ibm865
		cp1252 = charmap.Windows1252 // cp1252
		iso1   = charmap.ISO8859_1   // 1
		iso15  = charmap.ISO8859_15  // 15
		jis    = japanese.ShiftJIS   // shiftjis
		u8     = unicode.UTF8
		u8bom  = unicode.UTF8BOM
		u16    = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
		u16be  = unicode.UTF16(unicode.BigEndian, unicode.UseBOM)
		u16le  = unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)
		u32    = utf32.UTF32(utf32.LittleEndian, utf32.IgnoreBOM)
		u32be  = utf32.UTF32(utf32.BigEndian, utf32.UseBOM)
		u32le  = utf32.UTF32(utf32.LittleEndian, utf32.UseBOM)
	)
	m := map[string]Sample{
		"037":           {Text, cp037, "text/cp037.txt", "EBCDIC 037 IBM mainframe test"},
		"437":           {Dump, cp437, "text/cp437-crlf.txt", "CP-437 all characters test using Windows line breaks"},
		"437.cr":        {Dump, cp437, "text/cp437-cr.txt", "CP-437 all characters test using CR (carriage return)"},
		"437.lf":        {Dump, cp437, "text/cp437-lf.txt", "CP-437 all characters test using LF (line feed)"},
		"865":           {Text, cp865, "text/cp865.txt", "CP-865 and CP-860 Nordic test"},
		"1252":          {Text, cp1252, "text/cp1252.txt", "Windows-1252 English test"},
		"ascii":         {Text, cp437, "text/retrotxt.asc", fmt.Sprintf("%s ASCII logos", meta.Name)},
		"ansi":          {Ansi, cp437, "ansi/retrotxt.ans", fmt.Sprintf("%s 256 color ANSI logo", meta.Name)},
		"ansi.aix":      {Ansi, cp437, "ansi/ansi-aixterm.ans", "IBM AIX terminal colors"},
		"ansi.blank":    {Ansi, cp437, "ansi/ansi-blank.ans", "Empty file test"},
		"ansi.cp":       {Ansi, cp437, "ansi/ansi-cp.ans", "ANSI cursor position tests"},
		"ansi.cpf":      {Ansi, cp437, "ansi/ansi-cpf.ans", "ANSI cursor forward tests"},
		"ansi.hvp":      {Ansi, cp437, "ansi/ansi-hvp.ans", "ANSI horizontal and vertical cursor positioning"},
		"ansi.proof":    {Ansi, cp437, "text/ansi-proof.ans", "ANSI formatting proof sheet"},
		"ansi.rgb":      {Ansi, cp437, "ansi/ansi-rgb.ans", "ANSI RGB 24-bit color sheet"},
		"ansi.setmodes": {Ansi, cp437, "ansi/ansi-setmodes.ans", "MS-DOS ANSI.SYS Set Mode examples"},
		"iso-1":         {Text, iso1, "text/iso-8859-1.txt", "ISO 8859-1 select characters"},
		"iso-15":        {Text, iso15, "text/iso-8859-15.txt", "ISO 8859-15 select characters"},
		"sauce":         {Dump, cp437, "text/sauce.txt", "SAUCE metadata test"},
		"shiftjis":      {Text, jis, "text/shiftjis.txt", "Shift-JIS and Mona font test"},
		"us-ascii":      {Dump, u8, "text/us-ascii.txt", "US-ASCII controls test"},
		"utf8":          {Text, u8, "text/utf-8.txt", "UTF-8 test with no Byte Order Mark"},
		"utf8.bom":      {Text, u8bom, "text/utf-8-bom.txt", "UTF-8 test with a Byte Order Mark"},
		"utf16":         {Text, u16, "text/utf-16.txt", "UTF-16 test"},
		"utf16.be":      {Text, u16be, "text/utf-16-be.txt", "UTF-16 Big Endian test"},
		"utf16.le":      {Text, u16le, "text/utf-16-le.txt", "UTF-16 Little Endian test"},
		"utf32":         {Text, u32, "text/utf-32.txt", "UTF-32 test"},
		"utf32.be":      {Text, u32be, "text/utf-32-be.txt", "UTF-32 Big Endian test"},
		"utf32.le":      {Text, u32le, "text/utf-32-le.txt", "UTF-32 Little Endian test"},
	}
	return m
}

// Open the named sample text file.
func Open(name string) ([]byte, error) {
	name = strings.ToLower(name)
	samp, exist := Map()[name]
	if !exist {
		return nil, fmt.Errorf("%s: %w", name, ErrName)
	}
	b, err := static.File.ReadFile(samp.Name)
	if err != nil {
		return nil, fmt.Errorf("open sample %q: %w", samp.Name, err)
	}
	return b, nil
}

// Transform the byte array to use the supplied character encoding.
func Transform(e encoding.Encoding, b ...byte) ([]byte, error) {
	if e == nil {
		return nil, ErrEncode
	}
	p, err := e.NewEncoder().Bytes(b)
	if err != nil {
		if len(p) == 0 {
			return b, fmt.Errorf("encoder could not convert bytes to %s: %w", e, err)
		}
		return nil, fmt.Errorf("%s: %w", e, err)
	}
	return p, nil
}

// Open and convert the named sample text file into Unicode runes.
func (flag Flags) Open(conv *convert.Convert, name string) ([]rune, error) {
	name = strings.ToLower(name)
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		return nil, nil
	}
	samp, exist := Map()[name]
	if !exist {
		return nil, fmt.Errorf("%s: %w", name, ErrName)
	}
	b, err := static.File.ReadFile(samp.Name)
	if err != nil {
		return nil, fmt.Errorf("open sample %q: %w", samp.Name, err)
	}
	if conv == nil {
		return nil, ErrConvNil
	}
	conv.Input.Encoding = flag.Input
	if conv.Input.Encoding == nil {
		conv.Input.Encoding = samp.Encoding
	}
	// override --control flag values for codepage table samples
	ignoreCtrls := false
	switch name {
	case "437", "437.cr", "437.lf":
		ignoreCtrls = true
	default:
	}
	if ignoreCtrls {
		conv.Args.Controls = []string{}
	}
	r, err := samp.transform(conv, b...)
	if err != nil {
		return nil, err
	}
	return r, err
}

func (samp *Sample) transform(conv *convert.Convert, b ...byte) ([]rune, error) {
	switch samp.Convert {
	case Ansi:
		return conv.ANSI(b...)
	case Ctrl:
		return conv.Chars(b...)
	case Dump:
		return conv.Dump(b...)
	case Text:
		return conv.Text(b...)
	default:
		return nil, fmt.Errorf("transform sample %q: %w", samp.Convert, ErrConvert)
	}
}

// Valid conforms the named sample text file exists.
func Valid(name string) bool {
	if _, exist := Map()[name]; !exist {
		return false
	}
	return true
}
