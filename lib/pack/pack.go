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
	"retrotxt.com/retrotxt/internal/pack"
	"retrotxt.com/retrotxt/lib/convert"
	"retrotxt.com/retrotxt/lib/logs"
)

// Flags and configuration values by the user.
type Flags struct {
	Encode  encoding.Encoding
	xEncode string // todo enocoding type
	To      string
}

type internalPack struct {
	// convert method
	convert output
	// font choice
	font font
	// default character encoding for the packed data
	encoding encoding.Encoding
	// package name used in internal/pack/blob.go
	name string
	// package description
	description string
}

type (
	font   uint
	output uint
)

const (
	text output = iota // Obey common text controls.
	dump               // Ignore and print the common text controls as characters.

	vga  font = iota // VGA 8px font.
	mona             // Mona font for shift-jis.
)

var (
	cp037  = charmap.CodePage037
	cp367  = charmap.XUserDefined // IBM367 = us-ascii
	cp437  = charmap.CodePage437
	cp865  = charmap.CodePage865 // ibm865
	cp1252 = charmap.Windows1252 // cp1252
	iso1   = charmap.ISO8859_1   // 1
	iso15  = charmap.ISO8859_15  // 15
	jis    = japanese.ShiftJIS   // shiftjis
	u8     = unicode.UTF8
	u8bom  = unicode.UTF8BOM
	u16be  = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	u16le  = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
)

var (
	// ErrPackGet invalid pack name.
	ErrPackGet = errors.New("pack.get name is invalid")
	// ErrPackValue unknown pack value.
	ErrPackValue = errors.New("unknown package convert value")
)

var internalPacks = map[string]internalPack{
	"037":           {text, vga, cp037, "text/cp037.txt", "EBCDIC 037 IBM mainframe test"},
	"437.cr":        {dump, vga, cp437, "text/cp437-cr.txt", "CP-437 all characters test using CR (carriage return)"},  //
	"437.crlf":      {dump, vga, cp437, "text/cp437-crlf.txt", "CP-437 all characters test using Windows line breaks"}, //
	"437.lf":        {dump, vga, cp437, "text/cp437-lf.txt", "CP-437 all characters test using LF (line feed)"},        //
	"865":           {text, vga, cp865, "text/cp865.txt", "CP-865 and CP-860 Nordic test"},                             //
	"1252":          {text, vga, cp1252, "text/cp1252.txt", "Windows-1252 English test"},                               //
	"ascii":         {text, vga, cp437, "text/retrotxt.asc", "RetroTxt ASCII logos"},                                   //
	"ansi":          {text, vga, cp437, "text/retrotxt.ans", "RetroTxt 256 color ANSI logo"},                           //
	"ansi.aix":      {text, vga, cp437, "text/ansi-aixterm.ans", "IBM AIX terminal colours"},                           //
	"ansi.blank":    {text, vga, cp437, "text/ansi-blank.ans", "Empty file test"},                                      //
	"ansi.cp":       {text, vga, cp437, "text/ansi-cp.ans", "ANSI cursor position tests"},                              //
	"ansi.cpf":      {text, vga, cp437, "text/ansi-cpf.ans", "ANSI cursor forward tests"},                              //
	"ansi.hvp":      {text, vga, cp437, "text/ansi-hvp.ans", "ANSI horizontal and vertical cursor positioning"},        //
	"ansi.proof":    {text, vga, cp437, "text/ansi-proof.ans", "ANSI formatting proof sheet"},                          //
	"ansi.rgb":      {text, vga, cp437, "text/ansi-rgb.ans", "ANSI RGB 24-bit color sheet"},                            //
	"ansi.setmodes": {text, vga, cp437, "text/ansi-setmodes.ans", "MS-DOS ANSI.SYS Set Mode examples"},                 //
	"iso-1":         {text, vga, iso1, "text/iso-8859-1.txt", "ISO 8859-1 select characters"},                          //
	"iso-15":        {text, vga, iso15, "text/iso-8859-15.txt", "ISO 8859-15 select characters"},
	"sauce":         {text, vga, cp437, "text/sauce.txt", "SAUCE metadata test"},                   // todo: check the charmap is okay
	"shiftjis":      {dump, mona, jis, "text/shiftjis.txt", "Shift-JIS and Mona font test"},        // outputs to utf8?
	"us-ascii":      {dump, vga, cp367, "text/us-ascii.txt", "US-ASCII controls test"},             // us-ascii outputs to UTF8 because the control codes fail with the 8-bit codepages
	"utf8":          {text, vga, u8, "text/utf-8.txt", "UTF-8 test with no Byte Order Mark"},       //
	"utf8.bom":      {text, vga, u8bom, "text/utf-8-bom.txt", "UTF-8 test with a Byte Order Mark"}, //
	"utf16.be":      {text, vga, u16be, "text/utf-16-be.txt", "UTF-16 Big Endian test"},            //
	"utf16.le":      {text, vga, u16le, "text/utf-16-le.txt", "UTF-16 Little Endian test"},         //
}

// Open an internal packed example file.
func (f Flags) Open(conv convert.Args, name string) (r []rune, e encoding.Encoding, err error) {
	var s = strings.ToLower(name)
	if _, err = os.Stat(s); !os.IsNotExist(err) {
		return r, e, nil
	}
	pkg, exist := internalPacks[s]
	if !exist {
		return r, e, nil
	}
	b := pack.Get(pkg.name)
	if b == nil {
		return r, e, fmt.Errorf("view package %q: %w", pkg.name, ErrPackGet)
	}
	e = pkg.encoding
	if f.Encode == nil {
		conv.Encoding = fmt.Sprint(e)
	}
	if f.To != "" {
		// example exceptions that break the NewEncoder
		switch s {
		case "037", "shiftjis", "utf16.be", "utf16.le":
			return r, e, nil
		}
		transform(f.Encode, b...)
		return r, e, nil
	}
	// convert to runes and print
	switch pkg.convert {
	case dump:
		fmt.Println("dumping pack")
		if r, err = conv.Dump(&b); err != nil {
			return r, e, err
		}
	case text:
		fmt.Println("texting pack")
		if r, err = conv.Text(&b); err != nil {
			return r, e, err
		}
	default:
		return r, e, fmt.Errorf("view package %q: %w", pkg.convert, ErrPackValue)
	}
	return r, e, nil
}

// Valid package name?
func Valid(name string) bool {
	if _, exist := internalPacks[name]; !exist {
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
