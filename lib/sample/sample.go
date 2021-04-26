// Package sample opens and encodes the example textfiles embedded into the program.
package sample

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
	"retrotxt.com/retrotxt/lib/convert"
	"retrotxt.com/retrotxt/lib/create"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/static"
)

// Flags and configuration values by the user.
type Flags struct {
	From encoding.Encoding
	To   encoding.Encoding
}

// File details.
type File struct {
	Encoding encoding.Encoding
	Font     create.Font
	Runes    []rune
}

// Sample textfile data.
type Sample struct {
	// Convert text method.
	convert output
	// Font used the render the text.
	Font create.Font
	// Encoding used by the textfile.
	encoding encoding.Encoding
	// Name of the sample textfile.
	Name string
	// Description of the sample textfile.
	Description string
}

type output uint

const (
	Mona = create.Mona
	VGA  = create.VGA

	ansi output = iota // Only use with ANSI text.
	char               // Ignore and print the common text controls as characters.
	text               // Obey common text controls.
	dump               // Obey common text controls except end-of-file.
)

// Map the samples.
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
		"037":           {text, VGA, cp037, "text/cp037.txt", "EBCDIC 037 IBM mainframe test"},
		"437":           {dump, VGA, cp437, "text/cp437-crlf.txt", "CP-437 all characters test using Windows line breaks"},
		"437.cr":        {dump, VGA, cp437, "text/cp437-cr.txt", "CP-437 all characters test using CR (carriage return)"},
		"437.lf":        {dump, VGA, cp437, "text/cp437-lf.txt", "CP-437 all characters test using LF (line feed)"},
		"865":           {text, VGA, cp865, "text/cp865.txt", "CP-865 and CP-860 Nordic test"},
		"1252":          {text, VGA, cp1252, "text/cp1252.txt", "Windows-1252 English test"},
		"ascii":         {text, VGA, cp437, "text/retrotxt.asc", "RetroTxt ASCII logos"},
		"ansi":          {ansi, VGA, cp437, "text/retrotxt.ans", "RetroTxt 256 color ANSI logo"},
		"ansi.aix":      {ansi, VGA, cp437, "text/ansi-aixterm.ans", "IBM AIX terminal colors"},
		"ansi.blank":    {ansi, VGA, cp437, "text/ansi-blank.ans", "Empty file test"},
		"ansi.cp":       {ansi, VGA, cp437, "text/ansi-cp.ans", "ANSI cursor position tests"},
		"ansi.cpf":      {ansi, VGA, cp437, "text/ansi-cpf.ans", "ANSI cursor forward tests"},
		"ansi.hvp":      {ansi, VGA, cp437, "text/ansi-hvp.ans", "ANSI horizontal and vertical cursor positioning"},
		"ansi.proof":    {ansi, VGA, cp437, "text/ansi-proof.ans", "ANSI formatting proof sheet"},
		"ansi.rgb":      {ansi, VGA, cp437, "text/ansi-rgb.ans", "ANSI RGB 24-bit color sheet"},
		"ansi.setmodes": {ansi, VGA, cp437, "text/ansi-setmodes.ans", "MS-DOS ANSI.SYS Set Mode examples"},
		"iso-1":         {text, VGA, iso1, "text/iso-8859-1.txt", "ISO 8859-1 select characters"},
		"iso-15":        {text, VGA, iso15, "text/iso-8859-15.txt", "ISO 8859-15 select characters"},
		"sauce":         {text, VGA, cp437, "text/sauce.txt", "SAUCE metadata test"},
		"shiftjis":      {text, Mona, jis, "text/shiftjis.txt", "Shift-JIS and Mona font test"},
		"us-ascii":      {dump, VGA, u8, "text/us-ascii.txt", "US-ASCII controls test"},
		"utf8":          {text, VGA, u8, "text/utf-8.txt", "UTF-8 test with no Byte Order Mark"},
		"utf8.bom":      {text, VGA, u8bom, "text/utf-8-bom.txt", "UTF-8 test with a Byte Order Mark"},
		"utf16":         {text, VGA, u16, "text/utf-16.txt", "UTF-16 test"},
		"utf16.be":      {text, VGA, u16be, "text/utf-16-be.txt", "UTF-16 Big Endian test"},
		"utf16.le":      {text, VGA, u16le, "text/utf-16-le.txt", "UTF-16 Little Endian test"},
		"utf32":         {text, VGA, u32, "text/utf-32.txt", "UTF-32 test"},
		"utf32.be":      {text, VGA, u32be, "text/utf-32-be.txt", "UTF-32 Big Endian test"},
		"utf32.le":      {text, VGA, u32le, "text/utf-32-le.txt", "UTF-32 Little Endian test"},
	}
	return m
}

// Open a sample text file.
func Open(name string) ([]byte, error) {
	name = strings.ToLower(name)
	samp, exist := Map()[name]
	if !exist {
		return nil, logs.ErrSampFile
	}
	b, err := static.File.ReadFile(samp.Name)
	if err != nil {
		return nil, fmt.Errorf("open sample %q: %w", samp.Name, logs.ErrSampFile)
	}
	return b, nil
}

// Open and convert a sample textfile.
func (f Flags) Open(name string, conv *convert.Convert) (s File, err error) {
	name = strings.ToLower(name)
	if _, err = os.Stat(name); !os.IsNotExist(err) {
		return s, nil
	}
	samp, exist := Map()[name]
	if !exist {
		return s, logs.ErrSampFile
	}
	b, err := static.File.ReadFile(samp.Name)
	if err != nil {
		return s, fmt.Errorf("open sample %q: %w", samp.Name, logs.ErrSampFile)
	}
	if conv == nil {
		return s, ErrConvNil
	}
	s.Encoding = samp.encoding
	s.Font = samp.Font
	if f.From == nil {
		conv.Source.E = s.Encoding
	}
	if f.To != nil {
		// sample items that break the NewEncoder
		switch name {
		case "037", "shiftjis", "utf16.be", "utf16.le":
			return s, nil
		}
		err = printT(f.From, b...)
		if err != nil {
			return s, err
		}
		return s, nil
	}
	// default overrides
	switch name {
	case "437", "437.cr", "437.lf":
		if conv.Flags.Controls != nil {
			fmt.Printf("\nTo correctly display all %q table cells the --controls flag is ignored\n", name)
			conv.Flags.Controls = nil
		}
	}
	return samp.transform(&b, conv)
}

// Transform converts the raw byte data of the textfile into UTF8 runes.
func (samp Sample) transform(b *[]byte, conv *convert.Convert) (s File, err error) {
	switch samp.convert {
	case ansi:
		if s.Runes, err = conv.ANSI(b); err != nil {
			return s, err
		}
	case char:
		if s.Runes, err = conv.Chars(b); err != nil {
			return s, err
		}
	case dump:
		if s.Runes, err = conv.Dump(b); err != nil {
			return s, err
		}
	case text:
		if s.Runes, err = conv.Text(b); err != nil {
			return s, err
		}
	default:
		return s, fmt.Errorf("transform sample %q: %w", samp.convert, ErrConvert)
	}
	return s, nil
}

// Valid textfile name.
func Valid(name string) bool {
	if _, exist := Map()[name]; !exist {
		return false
	}
	return true
}

// PrintT encodes and prints the bytes as a string.
func printT(e encoding.Encoding, b ...byte) error {
	if e == nil {
		fmt.Println(string(b))
		return nil
	}
	b, err := encode(e, b...)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func encode(e encoding.Encoding, b ...byte) ([]byte, error) {
	nb, err := e.NewEncoder().Bytes(b)
	if err != nil {
		if len(nb) == 0 {
			return b, fmt.Errorf("encoder could not convert bytes to %s: %w", e, err)
		}
		return nb, nil
	}
	return nb, nil
}
