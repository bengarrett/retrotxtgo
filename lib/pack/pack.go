// Package pack ...
package pack

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
	"retrotxt.com/retrotxt/lib/convert"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/static"
)

// Flags and configuration values by the user.
type Flags struct {
	From encoding.Encoding
	To   encoding.Encoding
}

// Pack item details.
type Pack struct {
	Src   encoding.Encoding
	Font  Font
	Runes []rune
}

// Packs for items.
type Packs struct {
	// convert method
	convert output
	// font choice
	Font Font
	// default character encoding for the packed data
	encoding encoding.Encoding
	// package name used in internal/pack/blob.go
	Name string
	// package description
	Description string
}

type (
	// Font to use with pack item.
	Font   uint
	output uint
)

func (f Font) String() string {
	return [...]string{
		"", "vga", "mono",
	}[f]
}

const (
	ansi output = iota // Only use with ANSI text.
	char               // Ignore and print the common text controls as characters.
	text               // Obey common text controls.
	dump               // Obey common text controls except end-of-file.
)

const (
	_    Font = iota
	VGA       // VGA 8px font.
	Mona      // Mona font for shift-jis.
)

var (
	// ErrPackGet invalid pack name.
	ErrPackGet = errors.New("pack.get name is invalid")
	// ErrPackValue unknown pack value.
	ErrPackValue = errors.New("unknown package convert value")
)

// Map the pack items.
func Map() map[string]Packs {
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
	m := map[string]Packs{
		"037":           {text, VGA, cp037, "text/cp037.txt", "EBCDIC 037 IBM mainframe test"},
		"437.cr":        {dump, VGA, cp437, "text/cp437-cr.txt", "CP-437 all characters test using CR (carriage return)"},
		"437.crlf":      {dump, VGA, cp437, "text/cp437-crlf.txt", "CP-437 all characters test using Windows line breaks"},
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

// Open an internal packed example file.
func (f Flags) Open(conv *convert.Convert, name string) (p Pack, err error) {
	name = strings.ToLower(name)
	if _, err = os.Stat(name); !os.IsNotExist(err) {
		return p, nil
	}
	pkg, exist := Map()[name]
	if !exist {
		return p, nil
	}
	b, err := static.File.ReadFile(pkg.Name)
	if err != nil {
		return p, fmt.Errorf("view package %q: %w", pkg.Name, ErrPackGet)
	}
	p.Src = pkg.encoding
	p.Font = pkg.Font
	if f.From == nil {
		conv.Source.E = p.Src
	}
	if f.To != nil {
		// pack items that break the NewEncoder
		switch name {
		case "037", "shiftjis", "utf16.be", "utf16.le":
			return p, nil
		}
		transform(f.From, b...)
		return p, nil
	}
	// default overrides
	switch name {
	case "437.cr", "437.crlf", "437.lf":
		if conv.Flags.Controls != nil {
			fmt.Printf("\nTo correctly display all %q table cells the --controls flag is ignored\n", name)
			conv.Flags.Controls = nil
		}
	}
	// convert to runes and print
	return pkg.convPrint(&b, conv)
}

func (pkg Packs) convPrint(b *[]byte, conv *convert.Convert) (p Pack, err error) {
	switch pkg.convert {
	case ansi:
		if p.Runes, err = conv.ANSI(b); err != nil {
			return p, err
		}
	case char:
		if p.Runes, err = conv.Chars(b); err != nil {
			return p, err
		}
	case dump:
		if p.Runes, err = conv.Dump(b); err != nil {
			return p, err
		}
	case text:
		if p.Runes, err = conv.Text(b); err != nil {
			return p, err
		}
	default:
		return p, fmt.Errorf("view package %q: %w", pkg.convert, ErrPackValue)
	}
	return p, nil
}

// Valid package name?
func Valid(name string) bool {
	if _, exist := Map()[name]; !exist {
		return false
	}
	return true
}

func transform(e encoding.Encoding, b ...byte) {
	b, err := encode(e, b...)
	if err != nil {
		logs.Println("using the original encoding and not", fmt.Sprint(e), err)
	}
	fmt.Println(string(b))
}

func encode(e encoding.Encoding, b ...byte) ([]byte, error) {
	newer, err := e.NewEncoder().Bytes(b)
	if err != nil {
		if len(newer) == 0 {
			return b, fmt.Errorf("encoder could not convert bytes to %s: %w", e, err)
		}
		return newer, nil
	}
	return newer, nil
}
