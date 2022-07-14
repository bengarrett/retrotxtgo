package create_test

import (
	"reflect"
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd/internal/create"
	"github.com/bengarrett/retrotxtgo/lib/convert"
)

const (
	row1 = "☺☻♥♦♣♠•◘◙♂♀♪♫☼"
	want = "␀" + row1 + "\n"
)

var cp437 = []byte("\x00" + row1 + "\n")

func TestRunes(t *testing.T) {
	type args struct {
		encode string
		flags  convert.Flag
		src    *[]byte
	}
	tests := []struct {
		name    string
		args    args
		want    []rune
		wantErr bool
	}{
		{"empty", args{}, nil, true},
		{"no src", args{"invalid", convert.Flag{}, nil}, nil, true},
		{"no encoding", args{"", convert.Flag{}, &cp437}, nil, true},
		{"invalid encoding", args{"invalid", convert.Flag{}, &cp437}, nil, true},
		{"CP-437 encoding", args{"cp437", convert.Flag{}, &cp437}, []rune(string(want)), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := create.Runes(tt.args.encode, tt.args.flags, tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("HTML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HTML() = %v, want %v", got, tt.want)
			}
		})
	}
}
