package nl_test

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/bengarrett/retrotxtgo/pkg/nl"
)

func ExampleNewLine() {
	s := nl.NewLine(nl.Windows)
	fmt.Printf("%q\n", s)
	s = nl.NewLine(nl.Macintosh)
	fmt.Printf("%q\n", s)
	s = nl.NewLine(nl.Linux)
	fmt.Printf("%q\n", s)
	// Output: "\r\n"
	// "\r"
	// "\n"
}

func ExampleLineBreak_Find() {
	lb := &nl.LineBreak{}

	// A Windows example of a line break.
	windows := [2]rune{rune('\r'), rune('\n')}
	lb.Find(windows)
	fmt.Printf("%s %q %d %X\n", lb.Abbr, lb.Escape, lb.Decimal, lb.Decimal)

	// A macOS, Linux, or Unix example of a line break.
	linux := [2]rune{rune('\n')}
	lb.Find(linux)
	fmt.Printf("%s %q %d %X\n", lb.Abbr, lb.Escape, lb.Decimal, lb.Decimal)

	// Output: CRLF "\r\n" [13 10] [D A]
	// LF "\n" [10 0] [A 0]
}

func ExampleLineBreak_Total() {
	lb := &nl.LineBreak{}
	linux := [2]rune{rune('\n')}
	lb.Find(linux)
	// Count the number of lines in a file.
	total, _ := lb.Total("testdata/textlf.txt")
	fmt.Println(total)
	// Output: 7
}

func TestLines(t *testing.T) {
	lf := [2]rune{nl.LF}
	type args struct {
		r  io.Reader
		lb [2]rune
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		{"empty", args{strings.NewReader(""), lf}, 0, false},
		{"single line", args{strings.NewReader("hello world"), lf}, 1, false},
		{"multiple lines", args{strings.NewReader("hello\nworld\neof"), lf}, 3, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCount, err := nl.Lines(tt.args.r, tt.args.lb)
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
