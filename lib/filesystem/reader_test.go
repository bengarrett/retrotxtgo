package filesystem

import (
	"reflect"
	"testing"
)

func TestNewlines(t *testing.T) {
	tests := []struct {
		name string
		text []rune
		want [2]rune
	}{
		{"unix", []rune("hello\x0aworld\x0a"), [2]rune{10}},
		{"win", []rune("hello\x0d\x0aworld\x0d\x0a\x1a"), [2]rune{13, 10}},
		{"c64", []rune("hello\x0dworld\x0d"), [2]rune{13}},
		{"ibm", []rune("hello\x15world\x15"), [2]rune{21}},
		{"mix", []rune("\x15Windows newline: \x0d\x0a\x15Unix newline: \x0a\x15"), [2]rune{21}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Newlines(tt.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Data.Newlines() = %v, want %v", got, tt.want)
			}
		})
	}
}
