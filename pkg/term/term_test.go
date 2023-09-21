package term_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/bengarrett/retrotxtgo/pkg/term"
	"github.com/gookit/color"
)

func init() {
	color.Enable = false
}

func ExampleAlert() {
	fmt.Print(term.Alert(), "something went wrong")
	// Output:Problem:
	// something went wrong
}

func ExampleComment() {
	fmt.Print(term.Comment("Hi"))
	// Output:Hi
}

func ExampleBool() {
	fmt.Println(term.Bool(true), "yes")
	fmt.Println(term.Bool(false), "no")
	// Output:✓ yes
	// ✗ no
}

func ExampleOptions() {
	term.Options(os.Stdout, "this is a list of options",
		false, false, "option3", "option2", "option1")
	// Output:this is a list of options.
	//   Options: option1, option2, option3
}

func ExampleBorder() {
	term.Border(os.Stdout, "hi")
	// Output: ┌────┐
	// │ hi │
	// └────┘
}

func ExampleCenter() {
	fmt.Print("[" + term.Center(10, "hi") + "]")
	// Output:[    hi]
}

func ExampleHR() {
	term.HR(os.Stdout, 8)
	// Output:────────
}

func ExampleHead() {
	term.Head(os.Stdout, 10, "heading")
	// Output:──────────
	//  heading
}

func ExampleUnderlineChar() {
	color.Enable = true
	defer func() { color.Enable = false }()
	s, _ := term.UnderlineChar("Z")
	fmt.Print(s)
	// Output:[0m[4mZ[0m
}

func TestTerm(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		wantTerm string
	}{
		{"", "terminal256"},
		{"xterm", "terminal"},
		{"xterm-color", "terminal"},
		{"xterm-mono", "none"},
		{"rxvt-unicode-256color", "terminal256"},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			if gotTerm := term.Term("", tt.name); gotTerm != tt.wantTerm {
				t.Errorf("Term() = %v, want %v", gotTerm, tt.wantTerm)
			}
		}
	})
}

func TestTerm16M(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		wantTerm string
	}{
		{"24bit", "terminal16m"},
		{"truecolor", "terminal16m"},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			if gotTerm := term.Term(tt.name, ""); gotTerm != tt.wantTerm {
				t.Errorf("Term() = %v, want %v", gotTerm, tt.wantTerm)
			}
		}
	})
}

func TestUnderlineChar(t *testing.T) {
	color.Enable = true
	t.Cleanup(func() { color.Enable = false })
	t.Parallel()
	tests := []struct {
		name    string
		c       string
		wantS   string
		wantErr bool
	}{
		{"nil", "", "", false},
		{"ascii", "Z", "\x1b[0m\x1b[4mZ\x1b[0m", false},
		{"hex", "\x90", "", true},
		{"utf8", "\u005A", "\x1b[0m\x1b[4mZ\x1b[0m", false},
		{"░utf8", "\u2591", "\x1b[0m\x1b[4m░\x1b[0m", false},
		{"░hex", "\xe2\x96\x91", "\x1b[0m\x1b[4m░\x1b[0m", false},
		{"😀", "😀", "\x1b[0m\x1b[4m😀\x1b[0m", false},
		{"😀hex", "\xf0\x9f\x98\x80", "\x1b[0m\x1b[4m😀\x1b[0m", false},
		{"😀b", string([]byte{240, 159, 152, 128}), "\x1b[0m\x1b[4m😀\x1b[0m", false},
		{"string", "blahblah", "\x1b[0m\x1b[4mb\x1b[0m", false},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			gotS, err := term.UnderlineChar(tt.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("underlineChar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotS != tt.wantS {
				t.Errorf("underlineChar() = %v, want %v", gotS, tt.wantS)
			}
		}
	})
}

func TestCenter(t *testing.T) {
	t.Parallel()
	type args struct {
		text  string
		width int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{"", 0}, ""},
		{"even", args{"hi", 10}, "    hi"},
		{"odd", args{"hello", 10}, "  hello"},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			if got := term.Center(tt.args.width, tt.args.text); got != tt.want {
				t.Errorf("Center() = %q, want %q", got, tt.want)
			}
		}
	})
}
