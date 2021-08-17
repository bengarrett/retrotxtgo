package cmd

import (
	"testing"

	"github.com/gookit/color"
)

func Test_serveBytes(t *testing.T) {
	color.Enable = false
	html.Test = true
	b := []byte("hello world")
	type args struct {
		i       int
		changed bool
		b       *[]byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"too many", args{2, true, &b}, false},
		{"okay", args{0, true, &b}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := serveBytes(tt.args.i, tt.args.changed, tt.args.b); got != tt.want {
				t.Errorf("serveBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
