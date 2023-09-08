package filesystem_test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"
	"testing"

	"github.com/bengarrett/retrotxtgo/pkg/filesystem"
	"github.com/bengarrett/retrotxtgo/static"
)

func ExampleWordsEBCDIC() {
	b, err := static.File.ReadFile("text/cp037.txt")
	if err != nil {
		log.Fatal(err)
	}
	nr := bytes.NewReader(b)
	words, err := filesystem.WordsEBCDIC(nr)
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
		want filesystem.LB
	}{
		{"unix", []rune("hello\x0aworld\x0a"), filesystem.LF()},
		{"win", []rune("hello\x0d\x0aworld\x0d\x0a\x1a"), filesystem.CRLF()},
		{"c64", []rune("hello\x0dworld\x0d"), filesystem.CR()},
		{"ibm", []rune("hello\x15world\x15"), filesystem.NL()},
		{"mix", []rune("\x15Windows line break: \x0d\x0a\x15Unix line break: \x0a\x15"), filesystem.NL()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filesystem.LineBreaks(false, tt.text...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Data.LineBreaks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLineBreak(t *testing.T) {
	type args struct {
		r         filesystem.LB
		extraInfo bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{}, "??"},
		{"nl", args{filesystem.NL(), false}, "NL"},
		{"nl", args{filesystem.NL(), true}, "NL (IBM EBCDIC)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filesystem.LineBreak(tt.args.r, tt.args.extraInfo); got != tt.want {
				t.Errorf("LineBreak() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLines(t *testing.T) {
	type args struct {
		r  io.Reader
		lb filesystem.LB
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		{"empty", args{strings.NewReader(""), filesystem.LF()}, 0, false},
		{"single line", args{strings.NewReader("hello world"), filesystem.LF()}, 1, false},
		{"multiple lines", args{strings.NewReader("hello\nworld\neof"), filesystem.LF()}, 3, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCount, err := filesystem.Lines(tt.args.r, tt.args.lb)
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

func TestColumns(t *testing.T) {
	type args struct {
		r  io.Reader
		lb filesystem.LB
	}
	tests := []struct {
		name      string
		args      args
		wantWidth int
		wantErr   bool
	}{
		{"empty", args{}, 0, true},
		{"4 chars", args{strings.NewReader("abcd\n"), filesystem.LF()}, 4, false},
		{"4 runes", args{bytes.NewReader([]byte("😁😋😃🤫\n")), filesystem.LF()}, 16, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWidth, err := filesystem.Columns(tt.args.r, tt.args.lb)
			if (err != nil) != tt.wantErr {
				t.Errorf("Columns() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWidth != tt.wantWidth {
				t.Errorf("Columns() = %v, want %v", gotWidth, tt.wantWidth)
			}
		})
	}
}