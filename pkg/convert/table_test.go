package convert_test

import (
	"reflect"
	"testing"

	"github.com/bengarrett/retrotxtgo/pkg/convert"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	uni "golang.org/x/text/encoding/unicode"
)

func TestTable(t *testing.T) {
	tests := []struct {
		name    string
		wantNil bool
		wantErr bool
	}{
		{"IBM437", false, false},
		{"cp437", false, false},
		{"win", false, false},
		{"xxx", true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convert.Table(tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Table() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil != tt.wantNil {
				t.Errorf("Table() = %v, want %v", got, tt.wantNil)
			}
		})
	}
}

func Test_DefaultCP(t *testing.T) {
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convert.DefaultCP(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("DefaultCP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultCP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_character(t *testing.T) {
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convert.Character(tt.args.pos, tt.args.r, tt.args.cp); got != tt.want {
				t.Errorf("Character() = %q, want %q", got, tt.want)
			}
		})
	}
}

func Test_CharmapAlias(t *testing.T) {
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convert.CharmapAlias(tt.args.cp); got != tt.want {
				t.Errorf("CharmapAlias() = %v, want %v", got, tt.want)
			}
		})
	}
}
