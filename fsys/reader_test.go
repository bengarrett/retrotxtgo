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

	"github.com/bengarrett/retrotxtgo/fsys"
	"github.com/bengarrett/retrotxtgo/sample"
)

func ExampleWordsEBCDIC() {
	b, err := sample.File.ReadFile("plaintext/cp037.txt")
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
	t.Parallel()
	tests := []struct {
		name string
		text []rune
		want [2]rune
	}{
		{"unix", []rune("hello\x0aworld\x0a"), fsys.LF()},
		{"win", []rune("hello\x0d\x0aworld\x0d\x0a\x1a"), fsys.CRLF()},
		{"c64", []rune("hello\x0dworld\x0d"), fsys.CR()},
		{"ibm", []rune("hello\x15world\x15"), fsys.NL()},
		{"mix", []rune("\x15Windows line break: \x0d\x0a\x15Unix line break: \x0a\x15"), fsys.NL()},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			if got := fsys.LineBreaks(false, tt.text...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Data.LineBreaks() = %v, want %v", got, tt.want)
			}
		}
	})
}

func TestLineBreak(t *testing.T) {
	t.Parallel()
	type args struct {
		r         [2]rune
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
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			if got := fsys.LineBreak(tt.args.r, tt.args.extraInfo); got != tt.want {
				t.Errorf("LineBreak() = %v, want %v", got, tt.want)
			}
		}
	})
}

func TestColumns(t *testing.T) {
	t.Parallel()
	type args struct {
		r  io.Reader
		lb [2]rune
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
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			gotWidth, err := fsys.Columns(tt.args.r, tt.args.lb)
			if (err != nil) != tt.wantErr {
				t.Errorf("Columns() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWidth != tt.wantWidth {
				t.Errorf("Columns() = %v, want %v", gotWidth, tt.wantWidth)
			}
		}
	})
}
