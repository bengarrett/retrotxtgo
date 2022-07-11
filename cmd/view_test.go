package cmd_test

import (
	"bytes"
	"io/ioutil"
	"strconv"
	"testing"
	"unicode/utf8"

	"github.com/bengarrett/retrotxtgo/cmd"
	"github.com/gookit/color"
)

//  Formal name             | Named value   | Numeric value  | Alias value   |
//  *IBM Code Page 037      | cp037         | 37             | ibm037        |
//  IBM Code Page 437       | cp437         | 437            | msdos         |
//  IBM Code Page 850       | cp850         | 850            | latinI        |
//  IBM Code Page 852       | cp852         | 852            | latinII       |
//  IBM Code Page 855       | cp855         | 855            | ibm855        |
//  Windows Code Page 858   | cp858         | 858            | ibm00858      |
//  IBM Code Page 860       | cp860         | 860            | ibm860        |
//  IBM Code Page 862       | cp862         | 862            |               |
//  IBM Code Page 863       | cp863         | 863            | ibm863        |
//  IBM Code Page 865       | cp865         | 865            | ibm865        |
//  IBM Code Page 866       | ibm866        | 866            |               |
//  *IBM Code Page 1047     | cp1047        | 1047           | ibm1047       |
//  *IBM Code Page 1140     | cp1140        | 1140           | ibm01140      |
//  ISO 8859-1              | iso-8859-1    | 1              | latin1        |
//  ISO 8859-2              | iso-8859-2    | 2              | latin2        |
//  ISO 8859-3              | iso-8859-3    | 3              | latin3        |
//  ISO 8859-4              | iso-8859-4    | 4              | latin4        |
//  ISO 8859-5              | iso-8859-5    | 5              | cyrillic      |
//  ISO 8859-6              | iso-8859-6    | 6              | arabic        |
//  ISO-8859-6E             | iso-8859-6-e  |                | iso88596e     |
//  ISO-8859-6I             | iso-8859-6-i  |                | iso88596i     |
//  ISO 8859-7              | iso-8859-7    | 7              | greek         |
//  ISO 8859-8              | iso-8859-8    | 8              | hebrew        |
//  ISO-8859-8E             | iso-8859-8-e  |                | iso88598e     |
//  ISO-8859-8I             | iso-8859-8-i  |                | iso88598i     |
//  ISO 8859-9              | iso-8859-9    | 9              | latin5        |
//  ISO 8859-10             | iso-8859-10   | 10             | latin6        |
//  ISO 8895-11             | iso-8895-11   | 11             | iso889511     |
//  ISO 8859-13             | iso-8859-13   | 13             | iso885913     |
//  ISO 8859-14             | iso-8859-14   | 14             | iso885914     |
//  ISO 8859-15             | iso-8859-15   | 15             | iso885915     |
//  ISO 8859-16             | iso-8859-16   | 16             | iso885916     |
//  KOI8-R                  | koi8-r        |                | koi8r         |
//  KOI8-U                  | koi8-u        |                | koi8u         |
//  Macintosh               | macintosh     |                | mac           |
//  Windows 874             | cp874         | 874            | windows-874   |
//  Windows 1250            | cp1250        | 1250           | windows-1250  |
//  Windows 1251            | cp1251        | 1251           | windows-1251  |
//  Windows 1252            | cp1252        | 1252           | windows-1252  |
//  Windows 1253            | cp1253        | 1253           | windows-1253  |
//  Windows 1254            | cp1254        | 1254           | windows-1254  |
//  Windows 1255            | cp1255        | 1255           | windows-1255  |
//  Windows 1256            | cp1256        | 1256           | windows-1256  |
//  Windows 1257            | cp1257        | 1257           | windows-1257  |
//  Windows 1258            | cp1258        | 1258           | windows-1258  |
//  Shift JIS               | shift_jis     |                | shiftjis      |
//  UTF-8                   | utf-8         | 8              | utf8          |
//  UTF-16BE (Use BOM)      | utf-16        |                | utf16         |
//  UTF-16BE (Ignore BOM)   | utf-16be      |                | utf16be       |
//  UTF-16LE (Ignore BOM)   | utf-16le      |                | utf16le       |
//  UTF-32BE (Use BOM)      | utf-32        |                | utf32         |
//  UTF-32BE (Ignore BOM)   | utf-32be      |                | utf32be       |
//  UTF-32LE (Ignore BOM)   | utf-32le      |                | utf32le       |
//  **ASA X3.4 1963         | ascii-63      |                |               |
//  **ASA X3.4 1965         | ascii-65      |                |               |
//  **ANSI X3.4 1967/77/86  | ascii-67      |                |               |

// tester initialises, runs and returns the results of the view command.
// args are the command line arguments to test with the command.
func tester(args []string) ([]byte, error) {
	color.Enable = false
	b := bytes.NewBufferString("")
	cmd := cmd.ViewInit()
	cmd.SetOut(b)
	cmd.SetArgs(args)
	cmd.Execute()
	out, err := ioutil.ReadAll(b)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var (
	// arguments to test the --encoding flag
	cp037    = []string{"--encode", "cp037", "../static/text/cp037.txt"}
	cp437    = []string{"--controls", "", "../static/text/cp437-crlf.txt"} // overwrite default controls
	cp865    = []string{"../static/text/cp865.txt"}
	cp1252   = []string{"-e", "1252", "../static/text/cp1252.txt"}
	latin1   = []string{"--encode", "iso-8859-1", "../static/text/iso-8859-1.txt"}
	latin15  = []string{"--encode", "iso-8859-15", "../static/text/iso-8859-15.txt"}
	shiftjis = []string{"--encode", "shiftjis", "../static/text/shiftjis.txt"}
	utf8bom  = []string{"../static/text/utf-8-bom.txt"}
	utf16    = []string{"--encode", "utf16", "../static/text/utf-16.txt"}
	//utf32    = []string{"--encode", "utf32", "../static/text/utf-32.txt"} // TODO: failing
	// arguments to test the --controls flag
	noCtrls = []string{"--controls", "", "../static/text/us-ascii.txt"}
	eof     = []string{"--controls", "eof", "../static/text/us-ascii.txt"}
	tab     = []string{"--controls", "tab", "../static/text/us-ascii.txt"}
	bell    = []string{"--controls", "bell", "../static/text/us-ascii.txt"}
	bs      = []string{"--controls", "bs", "../static/text/us-ascii.txt"}
	del     = []string{"--controls", "del", "../static/text/us-ascii.txt"}
	esc     = []string{"--controls", "esc", "../static/text/us-ascii.txt"}
	ff      = []string{"--controls", "ff", "../static/text/us-ascii.txt"}
	vtab    = []string{"--controls", "vtab", "../static/text/us-ascii.txt"}
	// arguments to test the --swap-chars flag
	allCP437 = []string{"-c", "", "-x", "", "../static/text/cp437-crlf.txt"}
	null     = []string{"-c", "", "--swap-chars", "null", "../static/text/cp437-crlf.txt"}
	bar      = []string{"-c", "", "-x", "bar", "../static/text/cp437-crlf.txt"}
	house    = []string{"-c", "", "-x", "house", "../static/text/cp437-crlf.txt"}
	pipe     = []string{"-c", "", "-x", "pipe", "../static/text/cp437-crlf.txt"}
	root     = []string{"-c", "", "-x", "root", "../static/text/cp437-crlf.txt"}
	space    = []string{"-c", "", "-x", "space", "../static/text/cp437-crlf.txt"}
)

// todo, test for BOM with utf-8-bom.txt

func Test_ExecuteCommand(t *testing.T) {
	utfResults := []rune{'☠', '☮', '♺'}
	tests := []struct {
		name       string
		args       []string
		checkRunes []rune
	}{
		{"default encoding", cp437,
			[]rune{'≡', '₧', 'Ç'}},
		{"cp037", cp037,
			[]rune{'¶', '¦', '÷'}},
		{"cp865", cp865,
			[]rune{'█', '▓', '▒'}},
		{"cp1252", cp1252,
			[]rune{'‘', '’', '“', '”', '…', '™'}},
		{"latin-1", latin1,
			[]rune{'¤', '¾', '½'}},
		{"latin-15", latin15,
			[]rune{'€', 'Ÿ', 'œ'}},
		{"japanese", shiftjis,
			[]rune{'□', '■', 'つ', '∪'}},
		{"utf-8", utf8bom, utfResults},
		{"utf-16", utf16, utfResults},
		//{"utf-32", utf32, utfResults}, // TODO: failing
		{"no controls", noCtrls,
			[]rune{'○', '•', '⌂', '←', '→'}},
		{"hide rune end-of-file", eof,
			[]rune{'→'}},
		{"hide rune tab", tab,
			[]rune{'○'}},
		{"hide rune bell", bell,
			[]rune{'•'}},
		{"hide rune backspace", bs,
			[]rune{'◘'}},
		{"hide rune delete", del,
			[]rune{'⌂'}},
		{"hide rune escape", esc,
			[]rune{'←'}},
		{"hide rune form feed", ff,
			[]rune{'♀'}},
		{"hide rune vertical tab", vtab,
			[]rune{'♂'}},
		{"all cp-437", allCP437,
			[]rune{'␀', '|', '⌂', '│', '√', ' '}},
		{"hide rune c null", null,
			[]rune{'␀'}},
		{"hide rune bar", bar,
			[]rune{'|'}},
		{"swap bar", bar,
			[]rune{'¦'}},
		{"hide rune house", house,
			[]rune{'⌂'}},
		{"swap house", house,
			[]rune{'Δ'}},
		{"hide rune pipe", pipe,
			[]rune{'│'}},
		{"swap pipe", pipe,
			[]rune{'⎮'}},
		{"hide rune square root", root,
			[]rune{'√'}},
		{"swap square root", root,
			[]rune{'✓'}},
		{"hide rune space", space,
			[]rune{' '}},
		{"swap space", space,
			[]rune{'␣'}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i, runes := range tt.checkRunes {
				gotB, err := tester(tt.args)
				if err != nil {
					t.Error(err)
					return
				}
				findUnexpectedRune := len(tt.name) > 9 && tt.name[0:9] == "hide rune"
				if findUnexpectedRune {
					if bytes.ContainsRune(gotB, runes) {
						t.Errorf("%d. result contains a rune that it shouldn't, %q got:\n%s", i, runes, gotB)
						return
					}
					continue
				}
				if !bytes.ContainsRune(gotB, runes) {
					t.Errorf("%d. result doesn't include the expected rune, %q got:\n%s", i, runes, gotB)
					return
				}
			}
		})
	}
}

func Test_ExecuteCommand_width(t *testing.T) {
	tests := []struct {
		name  string
		width int
		wantW int
	}{
		{"80 char width", 80, 80},
		{"40 char width", 40, 40},
		{"1 char width", 1, 1},
		{"0 char width is ignored", 0, 80},
		{"-100 char width is ignored", -100, 80},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := strconv.Itoa(tt.width)
			gotB, err := tester([]string{"--width", w, "../static/text/cp437-crlf.txt"})
			if err != nil {
				t.Error(err)
				return
			}
			rows := bytes.Split(gotB, []byte("\n"))
			for i, row := range rows {
				if utf8.RuneCount(row) > tt.wantW {
					t.Errorf("row %d contains %d runes, wanted %d: %q",
						i, utf8.RuneCount(row), tt.wantW, row)
				}
			}
		})
	}
}
