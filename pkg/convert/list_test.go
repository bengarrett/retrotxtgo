package convert_test

import (
	"fmt"
	"testing"

	"github.com/bengarrett/retrotxtgo/pkg/convert"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
)

func TestEncodings(t *testing.T) {
	const totalEncodings = 53
	got, want := len(convert.Encodings()), totalEncodings
	if got != want {
		t.Errorf("Encodings() count = %v, want %v", got, want)
	}
}

func TestList(t *testing.T) {
	if got := convert.List(); got == nil {
		t.Errorf("List() do not want %v", got)
	}
}

func TestCells(t *testing.T) {
	tests := []struct {
		name    string
		e       encoding.Encoding
		cell    convert.Cell
		wantErr bool
	}{
		{"unknown", nil, convert.Cell{"", "", "", ""}, true},
		{
			"cp437", charmap.CodePage437,
			convert.Cell{"IBM Code Page 437", "cp437", "437", "msdos"},
			false,
		},
		{
			"latin6", charmap.ISO8859_10,
			convert.Cell{"ISO 8859-10", "iso-8859-10", "10", "latin6"},
			false,
		},
		{
			"utf8", unicode.UTF8,
			convert.Cell{"UTF-8", "utf-8", "", "utf8"},
			false,
		},
		{
			"utf32", utf32.UTF32(utf32.BigEndian, utf32.UseBOM),
			convert.Cell{"UTF-32BE (Use BOM)", "utf-32", "", "utf32"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := convert.Cells(tt.e)
			if got := (err != nil); got != tt.wantErr {
				t.Errorf("Cells() error is %v, want %v", got, tt.wantErr)
			}
			fmt.Printf("%#v", c)
			if c.Name != tt.cell.Name {
				t.Errorf("Cells() gotN = %v, want %v", c.Name, tt.cell.Name)
			}
			if c.Value != tt.cell.Value {
				t.Errorf("Cells() gotV = %v, want %v", c.Value, tt.cell.Value)
			}
			if c.Numeric != tt.cell.Numeric {
				t.Errorf("Cells() gotV = %v, want %v", c.Numeric, tt.cell.Numeric)
			}
			if c.Alias != tt.cell.Alias {
				t.Errorf("Cells() gotA = %v, want %v", c.Alias, tt.cell.Alias)

				a, _ := ianaindex.MIB.Name(tt.e)
				t.Error(a)
			}
		})
	}
}

func TestAliasFmt(t *testing.T) {
	type args struct {
		s   string
		val string
		e   encoding.Encoding
	}
	tests := []struct {
		name    string
		args    args
		want    string
		WantErr bool
	}{
		{"err", args{"", "", nil}, "", true},
		{"empty cp037", args{"", "", charmap.CodePage037}, "ibm037", false},
		{"dupe cp037", args{"", "ibm037", charmap.CodePage037}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convert.AliasFmt(tt.args.s, tt.args.val, tt.args.e)
			if got != tt.want {
				t.Errorf("AliasFmt() = %v, want %v", got, tt.want)
			}
			if (err != nil) != tt.WantErr {
				t.Errorf("AliasFmt() error = %v, want %v", err, tt.WantErr)
			}
		})
	}
}

func TestUniform(t *testing.T) {
	type args struct {
		mime string
	}
	tests := []struct {
		name  string
		args  args
		wantS string
	}{
		{"IBM00", args{"IBM00858"}, "CP858"},
		{"IBM01", args{"IBM01140"}, "CP1140"},
		{"IBM", args{"IBM850"}, "CP850"},
		{"windows-1252", args{"windows-1252"}, "CP1252"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotS := convert.Uniform(tt.args.mime); gotS != tt.wantS {
				t.Errorf("uniform() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}

func TestNumeric(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{"", -1},
		{"a string with no numbers", -1},
		{"UTF-8", -1},
		{"ISO 8859", 8859},
		{"ISO 8859 1", 1},
		{"ISO 8859-1", 1},
		{"ISO 8859-16", 16},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convert.Numeric(tt.name); got != tt.want {
				t.Errorf("Numeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Encodings()
