package fsys_test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/bengarrett/retrotxtgo/pkg/fsys"
	"github.com/bengarrett/retrotxtgo/static"
)

func ExampleWordsEBCDIC() {
	b, err := static.File.ReadFile("text/cp037.txt")
	if err != nil {
		log.Fatal(err)
	}
	nr := bytes.NewReader(b)
	words, err := fsys.WordsEBCDIC(nr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(os.Stdout, "%d words", words)
	// Output: 16 words
}

func TestLineBreaks(t *testing.T) {
	tests := []struct {
		name string
		text []rune
		want fsys.LB
	}{
		{"unix", []rune("hello\x0aworld\x0a"), fsys.LF()},
		{"win", []rune("hello\x0d\x0aworld\x0d\x0a\x1a"), fsys.CRLF()},
		{"c64", []rune("hello\x0dworld\x0d"), fsys.CR()},
		{"ibm", []rune("hello\x15world\x15"), fsys.NL()},
		{"mix", []rune("\x15Windows line break: \x0d\x0a\x15Unix line break: \x0a\x15"), fsys.NL()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fsys.LineBreaks(false, tt.text...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Data.LineBreaks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLineBreak(t *testing.T) {
	type args struct {
		r         fsys.LB
		extraInfo bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{}, "??"},
		{"nl", args{fsys.NL(), false}, "NL"},
		{"nl", args{fsys.NL(), true}, "NL (IBM EBCDIC)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fsys.LineBreak(tt.args.r, tt.args.extraInfo); got != tt.want {
				t.Errorf("LineBreak() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLines(t *testing.T) {
	type args struct {
		r  io.Reader
		lb fsys.LB
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		{"empty", args{strings.NewReader(""), fsys.LF()}, 0, false},
		{"single line", args{strings.NewReader("hello world"), fsys.LF()}, 1, false},
		{"multiple lines", args{strings.NewReader("hello\nworld\neof"), fsys.LF()}, 3, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCount, err := fsys.Lines(tt.args.r, tt.args.lb)
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
		lb fsys.LB
	}
	tests := []struct {
		name      string
		args      args
		wantWidth int
		wantErr   bool
	}{
		{"empty", args{}, 0, true},
		{"4 chars", args{strings.NewReader("abcd\n"), fsys.LF()}, 4, false},
		{"4 runes", args{bytes.NewReader([]byte("😁😋😃🤫\n")), fsys.LF()}, 16, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWidth, err := fsys.Columns(tt.args.r, tt.args.lb)
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
