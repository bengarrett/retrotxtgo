package convert

import (
	"testing"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

func TestEncodings(t *testing.T) {
	const totalEncodings = 50
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
		name  string
		args  args
		wantN string
		wantV string
		wantD string
		wantA string
	}{
		{"unknown", args{}, "", "", "", ""},
		{"cp437", args{charmap.CodePage437}, "IBM Code Page 437", "cp437", "437", "msdos"},
		{"latin6", args{charmap.ISO8859_10}, "ISO 8859-10", "iso-8859-10", "10", "latin6"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotN, gotV, gotD, gotA := cells(tt.args.e)
			if gotN != tt.wantN {
				t.Errorf("cells() gotN = %v, want %v", gotN, tt.wantN)
			}
			if gotV != tt.wantV {
				t.Errorf("cells() gotV = %v, want %v", gotV, tt.wantV)
			}
			if gotD != tt.wantD {
				t.Errorf("cells() gotV = %v, want %v", gotD, tt.wantD)
			}
			if gotA != tt.wantA {
				t.Errorf("cells() gotA = %v, want %v", gotA, tt.wantA)
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
