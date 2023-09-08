//nolint:gochecknoglobals
package cmd_test

import (
	"bytes"
	"strconv"
	"testing"
	"unicode/utf8"
)

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
	// utf32    = []string{"--encode", "utf32", "../static/text/utf-32.txt"} // TODO: failing
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

func Test_ViewCommand(t *testing.T) {
	utfResults := []rune{'☠', '☮', '♺'}
	tests := []struct {
		name       string
		args       []string
		checkRunes []rune
	}{
		{
			"default encoding", cp437,
			[]rune{'≡', '₧', 'Ç'},
		},
		{
			"cp037", cp037,
			[]rune{'¶', '¦', '÷'},
		},
		{
			"cp865", cp865,
			[]rune{'█', '▓', '▒'},
		},
		{
			"cp1252", cp1252,
			[]rune{'‘', '’', '“', '”', '…', '™'},
		},
		{
			"latin-1", latin1,
			[]rune{'¤', '¾', '½'},
		},
		{
			"latin-15", latin15,
			[]rune{'€', 'Ÿ', 'œ'},
		},
		{
			"japanese", shiftjis,
			[]rune{'□', '■', 'つ', '∪'},
		},
		{"utf-8", utf8bom, utfResults},
		{"utf-16", utf16, utfResults},
		//{"utf-32", utf32, utfResults}, // TODO: failing
		{
			"no controls", noCtrls,
			[]rune{'○', '•', '⌂', '←', '→'},
		},
		{
			"hide rune end-of-file", eof,
			[]rune{'→'},
		},
		{
			"hide rune tab", tab,
			[]rune{'○'},
		},
		{
			"hide rune bell", bell,
			[]rune{'•'},
		},
		{
			"hide rune backspace", bs,
			[]rune{'◘'},
		},
		{
			"hide rune delete", del,
			[]rune{'⌂'},
		},
		{
			"hide rune escape", esc,
			[]rune{'←'},
		},
		{
			"hide rune form feed", ff,
			[]rune{'♀'},
		},
		{
			"hide rune vertical tab", vtab,
			[]rune{'♂'},
		},
		{
			"all cp-437", allCP437,
			[]rune{'␀', '|', '⌂', '│', '√', ' '},
		},
		{
			"hide rune c null", null,
			[]rune{'␀'},
		},
		{
			"hide rune bar", bar,
			[]rune{'|'},
		},
		{
			"swap bar", bar,
			[]rune{'¦'},
		},
		{
			"hide rune house", house,
			[]rune{'⌂'},
		},
		{
			"swap house", house,
			[]rune{'Δ'},
		},
		{
			"hide rune pipe", pipe,
			[]rune{'│'},
		},
		{
			"swap pipe", pipe,
			[]rune{'⎮'},
		},
		{
			"hide rune square root", root,
			[]rune{'√'},
		},
		{
			"swap square root", root,
			[]rune{'✓'},
		},
		{
			"hide rune space", space,
			[]rune{' '},
		},
		{
			"swap space", space,
			[]rune{'␣'},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i, runes := range tt.checkRunes {
				gotB, err := viewT.tester(tt.args)
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
			gotB, err := viewT.tester([]string{"--width", w, "../static/text/cp437-crlf.txt"})
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
