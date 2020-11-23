package filesystem

import (
	"reflect"
	"testing"
)

func TestLineBreaks(t *testing.T) {
	tests := []struct {
		name string
		text []rune
		want LB
	}{
		{"unix", []rune("hello\x0aworld\x0a"), LF()},
		{"win", []rune("hello\x0d\x0aworld\x0d\x0a\x1a"), CRLF()},
		{"c64", []rune("hello\x0dworld\x0d"), CR()},
		{"ibm", []rune("hello\x15world\x15"), NL()},
		{"mix", []rune("\x15Windows line break: \x0d\x0a\x15Unix line break: \x0a\x15"), NL()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LineBreaks(false, tt.text...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Data.LineBreaks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLineBreak(t *testing.T) {
	type args struct {
		r         LB
		extraInfo bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{}, "??"},
		{"nl", args{NL(), false}, "NL"},
		{"nl", args{NL(), true}, "NL (IBM EBCDIC)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LineBreak(tt.args.r, tt.args.extraInfo); got != tt.want {
				t.Errorf("LineBreak() = %v, want %v", got, tt.want)
			}
		})
	}
}
