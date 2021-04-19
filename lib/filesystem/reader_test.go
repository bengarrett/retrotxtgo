package filesystem

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"
	"testing"

	"retrotxt.com/retrotxt/static"
)

func ExampleWordsEBCDIC() {
	b, err := static.File.ReadFile("text/cp037.txt")
	if err != nil {
		log.Fatal(err)
	}
	nr := bytes.NewReader(b)
	words, err := WordsEBCDIC(nr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d words", words)
	// Output: 16 words
}

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

func TestLines(t *testing.T) {
	type args struct {
		r  io.Reader
		lb LB
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		{"empty", args{strings.NewReader(""), LF()}, 0, false},
		{"single line", args{strings.NewReader("hello world"), LF()}, 1, false},
		{"multiple lines", args{strings.NewReader("hello\nworld\neof"), LF()}, 3, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCount, err := Lines(tt.args.r, tt.args.lb)
			if (err != nil) != tt.wantErr {
				t.Errorf("Lines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCount != tt.wantCount {
				t.Errorf("Lines() = %v, want %v", gotCount, tt.wantCount)
			}
		})
	}
}
