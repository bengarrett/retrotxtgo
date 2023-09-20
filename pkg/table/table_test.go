package table_test

import (
	"reflect"
	"testing"

	"github.com/bengarrett/retrotxtgo/pkg/table"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	uni "golang.org/x/text/encoding/unicode"
)

func TestTable(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"IBM437", false},
		{"cp437", false},
		{"win", false},
		{"xxx", true},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			err := table.Table(nil, tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Table() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		}
	})
}

func Test_CodePage(t *testing.T) {
	t.Parallel()
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    encoding.Encoding
		wantErr bool
	}{
		{"empty", args{}, nil, true},
		{"none", args{"helloworld"}, nil, true},
		{"437", args{"cp437"}, charmap.CodePage437, false},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			got, err := table.CodePage(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("CodePage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CodePage() = %v, want %v", got, tt.want)
			}
		}
	})
}

func Test_character(t *testing.T) {
	t.Parallel()
	const RLM = 8207
	iso11 := charmap.XUserDefined
	type args struct {
		pos int
		r   rune
		cp  encoding.Encoding
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{}, "\x00"},
		{"iso8859_11", args{128, 0, iso11}, " "},
		{"utf8 space", args{128, 0, uni.UTF8}, " "},
		{"utf8 space", args{160, 0, uni.UTF8}, "\u00a0"},
		{"utf8 space", args{RLM, 0, uni.UTF8}, "\u200f"},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			if got := table.Character(tt.args.cp, tt.args.pos, tt.args.r); got != tt.want {
				t.Errorf("Character() = %q, want %q", got, tt.want)
			}
		}
	})
}

func Test_CharmapAlias(t *testing.T) {
	t.Parallel()
	type args struct {
		cp encoding.Encoding
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{}, ""},
		{"win874", args{charmap.Windows874}, " (Thai)"},
		{"shiftjis", args{japanese.ShiftJIS}, " (Japanese)"},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			if got := table.CharmapAlias(tt.args.cp); got != tt.want {
				t.Errorf("CharmapAlias() = %v, want %v", got, tt.want)
			}
		}
	})
}
