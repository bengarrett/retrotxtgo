package convert_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/bengarrett/retrotxtgo/pkg/convert"
	"github.com/bengarrett/retrotxtgo/pkg/fsys"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
)

const (
	ascii    = "Hello world!"
	abc      = "\x00 \x01 \x02\x0D\x0A\x1b[0mA B C"
	eof      = "Hello\x1Aworld!"
	wantAbc  = "␀ ☺ ☻♪◙←[0mA B C"
	wantEOF1 = "Hello"
	wantEOF0 = "Hello→world!"

	filename = "testdata/convert_test.ans"
)

func Test_SkipCtrlCodes(t *testing.T) {
	tests := []struct {
		name string
		ctrl []string
		want []rune
	}{
		{"nil", []string{}, []rune{}},
		{"bs", []string{"bs"}, []rune{convert.BS}},
		{"v,del", []string{"v", "del"}, []rune{convert.VT, convert.DEL}},
		{"invalid", []string{"xxx"}, []rune{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := convert.Convert{}
			c.LineBreaks = true
			c.Flags.Controls = tt.ctrl
			c.SkipCtrlCodes()
			if got := c.Ignores; string(got) != string(tt.want) {
				t.Errorf("Convert.SkipCtrlCodes() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvert_ANSI(t *testing.T) {
	const wantHi = "\x1b[0m\x0D\x0A" +
		"\x1b[1;33m  ╓──────────────╖\x1b[0m\x0D\x0A" +
		"\x1b[1;33m  ║\x1b[0m  Hello world \x1b[1;33m║\x1b[0m\x0D\x0A" +
		"\x1b[1;33m  ╙──────────────╜\x1b[0m\x0D\x0A"
	ba := []byte(ascii)
	be := []byte(eof)
	hi, err := fsys.Read(filename)
	if err != nil {
		t.Error(err)
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []rune
		wantErr bool
	}{
		{"empty", args{}, nil, true},
		{"ascii", args{ba}, bytes.Runes(ba), false},
		{"eof", args{be}, bytes.Runes([]byte(wantEOF1)), false},
		{"ansi", args{hi}, bytes.Runes([]byte(wantHi)), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := convert.Convert{}
			c.Input.Encoding = charmap.CodePage437
			got, err := c.ANSI(tt.args.b...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert.ANSI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				g, w := string(got), string(tt.want)
				t.Errorf("Convert.ANSI() = %s (%d), want %s (%d)",
					g, len(g), w, len(w))
			}
		})
	}
}

func TestConvert_Chars(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []rune
		wantErr bool
	}{
		{"empty", args{}, nil, true},
		{"ascii", args{[]byte(ascii)}, bytes.Runes([]byte(ascii)), false},
		{"eof", args{[]byte(eof)}, bytes.Runes([]byte(wantEOF0)), false},
		{"ansi", args{[]byte(abc)}, bytes.Runes([]byte(wantAbc)), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := convert.Convert{}
			c.Input.Encoding = charmap.CodePage437
			got, err := c.Chars(tt.args.b...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert.Chars() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				g, w := string(got), string(tt.want)
				t.Errorf("Convert.Chars() = %s (%d), want %s (%d)",
					g, len(g), w, len(w))
			}
		})
	}
}

func TestConvert_Dump(t *testing.T) {
	bhi := []byte("hello\nworld")
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []rune
		wantErr bool
	}{
		{"empty", args{}, nil, true},
		{"nl", args{bhi}, bytes.Runes(bhi), false},
		{"ascii", args{[]byte(ascii)}, bytes.Runes([]byte(ascii)), false},
		{"eof", args{[]byte(eof)}, bytes.Runes([]byte(wantEOF0)), false},
		{"abc", args{[]byte(abc)}, []rune("␀ ☺ ☻\r\n\x1b[0mA B C"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := convert.Convert{}
			c.Input.Encoding = charmap.CodePage437
			got, err := c.Dump(tt.args.b...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert.Dump() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				g, w := string(got), string(tt.want)
				t.Errorf("Convert.Dump() = %q (%dB), want %q (%dB)",
					g, len(g), w, len(w))
			}
		})
	}
}

func TestConvert_Text(t *testing.T) {
	const wantHi = "\x1b[0m\x0D\x0A" +
		"\x1b[1;33m  ╓──────────────╖\x1b[0m\x0D\x0A" +
		"\x1b[1;33m  ║\x1b[0m  Hello world \x1b[1;33m║\x1b[0m\x0D\x0A" +
		"\x1b[1;33m  ╙──────────────╜\x1b[0m\x0D\x0A"
	ba := []byte(ascii)
	be := []byte(eof)
	hi, err := fsys.Read(filename)
	if err != nil {
		t.Error(err)
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []rune
		wantErr bool
	}{
		{"empty", args{}, nil, true},
		{"ascii", args{ba}, bytes.Runes(ba), false},
		{"eof", args{be}, bytes.Runes([]byte(wantEOF1)), false},
		{"ansi", args{hi}, bytes.Runes([]byte(wantHi)), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := convert.Convert{}
			c.Flags.SwapChars = []string{"null", "bar"}
			c.Input.Encoding = charmap.CodePage437
			got, err := c.Text(tt.args.b...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert.Text() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				g, w := string(got), string(tt.want)
				t.Errorf("Convert.Text() = %s (%d), want %s (%d)",
					g, len(g), w, len(w))
			}
		})
	}
}

func TestConvert_FixJISTable(t *testing.T) {
	fix := []byte("\x7f\xa0\xe0\xff")
	c := convert.Convert{}
	c.Input.Bytes = fix
	c.Input.Table = true
	c.FixJISTable()
	if !reflect.DeepEqual(c.Input.Bytes, fix) {
		t.Errorf("Convert.FixJISTable() = %s, want %s", c.Input.Bytes, fix)
	}
	c.Input.Encoding = japanese.ShiftJIS
	c.Input.Bytes = fix
	c.FixJISTable()
	if want := []byte("\u007f   "); !reflect.DeepEqual(c.Input.Bytes, want) {
		t.Errorf("Convert.FixJISTable() = %q, want %q", c.Input.Bytes, want)
	}
}

func TestConvert_wrapWidth(t *testing.T) {
	type args struct {
		max   int
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"empty", args{}, "", true},
		{"no string", args{80, ""}, "", true},
		{"string", args{80, "abcdefghi"}, "abcdefghi", false},
		{"3 chrs", args{3, "abcdefghi"}, "abc\ndef\nghi\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := convert.Convert{}
			c.Input.Encoding = charmap.CodePage437
			c.Flags.MaxWidth = tt.args.max
			c.Input.LineBreak = [2]rune{13, 0}
			r, err := c.Chars([]byte(tt.args.input)...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert.wrapWidth() error = %v, want %v", err, tt.wantErr)
			}
			if string(r) != tt.want {
				t.Errorf("Convert.wrapWidth(%d) = %q, want %q", c.Flags.MaxWidth, string(r), tt.want)
			}
		})
	}
}
