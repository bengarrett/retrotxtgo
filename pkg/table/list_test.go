package table_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bengarrett/retrotxtgo/pkg/table"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
)

func ExampleAlias() {
	s, _ := table.Alias("", "cp437", nil)
	fmt.Println(s)
	// Output:msdos
}

func TestCharmaps(t *testing.T) {
	t.Parallel()
	const totalCharmaps = 53
	got, want := len(table.Charmaps()), totalCharmaps
	if got != want {
		t.Errorf("Charmaps() count = %v, want %v", got, want)
	}
}

func TestList(t *testing.T) {
	t.Parallel()
	w := &strings.Builder{}
	_ = table.List(w)
	if w.String() == "" {
		t.Errorf("List() do not want %v", w)
	}
}

func TestRows(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		e       encoding.Encoding
		cell    table.Row
		wantErr bool
	}{
		{"unknown", nil, table.Row{"", "", "", ""}, true},
		{
			"cp437", charmap.CodePage437,
			table.Row{"IBM Code Page 437", "cp437", "437", "msdos"},
			false,
		},
		{
			"latin6", charmap.ISO8859_10,
			table.Row{"ISO 8859-10", "iso-8859-10", "10", "latin6"},
			false,
		},
		{
			"utf8", unicode.UTF8,
			table.Row{"UTF-8", "utf-8", "", "utf8"},
			false,
		},
		{
			"utf32", utf32.UTF32(utf32.BigEndian, utf32.UseBOM),
			table.Row{"UTF-32BE (Use BOM)", "utf-32", "", "utf32"},
			false,
		},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			c, err := table.Rows(tt.e)
			if got := (err != nil); got != tt.wantErr {
				t.Errorf("Rows() error is %v, want %v", got, tt.wantErr)
			}
			if c.Name != tt.cell.Name {
				t.Errorf("Rows() gotN = %v, want %v", c.Name, tt.cell.Name)
			}
			if c.Value != tt.cell.Value {
				t.Errorf("Rows() gotV = %v, want %v", c.Value, tt.cell.Value)
			}
			if c.Numeric != tt.cell.Numeric {
				t.Errorf("Rows() gotV = %v, want %v", c.Numeric, tt.cell.Numeric)
			}
			if c.Alias != tt.cell.Alias {
				t.Errorf("Rows() gotA = %v, want %v", c.Alias, tt.cell.Alias)

				a, _ := ianaindex.MIB.Name(tt.e)
				t.Error(a)
			}
		}
	})
}

func TestAliasFmt(t *testing.T) {
	t.Parallel()
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
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			got, err := table.Alias(tt.args.s, tt.args.val, tt.args.e)
			if got != tt.want {
				t.Errorf("Alias() = %v, want %v", got, tt.want)
			}
			if (err != nil) != tt.WantErr {
				t.Errorf("Alias() error = %v, want %v", err, tt.WantErr)
			}
		}
	})
}

func TestUniform(t *testing.T) {
	t.Parallel()
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
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			if gotS := table.Uniform(tt.args.mime); gotS != tt.wantS {
				t.Errorf("uniform() = %v, want %v", gotS, tt.wantS)
			}
		}
	})
}

func TestNumeric(t *testing.T) {
	t.Parallel()
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
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			if got := table.Numeric(tt.name); got != tt.want {
				t.Errorf("Numeric() = %v, want %v", got, tt.want)
			}
		}
	})
}
