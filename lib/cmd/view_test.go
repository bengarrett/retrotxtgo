package cmd

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/convert"
)

func Example_viewToFlag() {
	r := []rune("hello")
	viewFlag.to = "cp437"
	viewToFlag(r...)
	// Output: hello
}

func Example_endOfFile() {
	var f convert.Flags
	f.Controls = []string{eof}
	fmt.Print(endOfFile(f))
	// Output: true
}

func Test_viewEncode(t *testing.T) {
	const hi = "hello world"
	rn := []rune(hi)
	type args struct {
		name string
		r    []rune
	}
	tests := []struct {
		name    string
		args    args
		wantB   []byte
		wantErr bool
	}{
		{"empty", args{}, nil, true},
		{"invalid", args{"invalid", rn}, nil, true},
		{"okay", args{"cp437", rn}, []byte(hi), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotB, err := viewEncode(tt.args.name, tt.args.r...)
			if (err != nil) != tt.wantErr {
				t.Errorf("viewEncode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotB, tt.wantB) {
				t.Errorf("viewEncode() = %s, want %s", gotB, tt.wantB)
			}
		})
	}
}
