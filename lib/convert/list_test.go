package convert

import (
	"fmt"
	"testing"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

func TestEncodings(t *testing.T) {
	const totalEncodings = 53
	got, want := len(Encodings()), totalEncodings
	if got != want {
		t.Errorf("Encodings() count = %v, want %v", got, want)
	}
}

func TestList(t *testing.T) {
	if got := List(); got == nil {
		t.Errorf("List() do not want %v", got)
	}
}

func Test_cells(t *testing.T) {
	type args struct {
		e encoding.Encoding
	}
	tests := []struct {
		name string
		args args
		cell cell
	}{
		{"unknown", args{}, cell{"", "", "", ""}},
		{"cp437", args{charmap.CodePage437}, cell{"IBM Code Page 437", "cp437", "437", "msdos"}},
		{"latin6", args{charmap.ISO8859_10}, cell{"ISO 8859-10", "iso-8859-10", "10", "latin6"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := cells(tt.args.e)
			fmt.Printf("%#v", c)
			if c.name != tt.cell.name {
				t.Errorf("cells() gotN = %v, want %v", c.name, tt.cell.name)
			}
			if c.value != tt.cell.value {
				t.Errorf("cells() gotV = %v, want %v", c.value, tt.cell.value)
			}
			if c.numeric != tt.cell.numeric {
				t.Errorf("cells() gotV = %v, want %v", c.numeric, tt.cell.numeric)
			}
			if c.alias != tt.cell.alias {
				t.Errorf("cells() gotA = %v, want %v", c.alias, tt.cell.alias)
			}
		})
	}
}

func Test_alias(t *testing.T) {
	type args struct {
		s   string
		val string
		e   encoding.Encoding
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"err", args{"", "", nil}, ""},
		{"empty cp037", args{"", "", charmap.CodePage037}, "ibm037"},
		{"dupe cp037", args{"", "ibm037", charmap.CodePage037}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := alias(tt.args.s, tt.args.val, tt.args.e); got != tt.want {
				t.Errorf("alias() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uniform(t *testing.T) {
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
			if gotS := uniform(tt.args.mime); gotS != tt.wantS {
				t.Errorf("uniform() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}
