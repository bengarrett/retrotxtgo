package term_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/bengarrett/retrotxtgo/pkg/term"
	"github.com/gookit/color"
)

func ExampleAlert() {
	color.Enable = false
	fmt.Fprint(os.Stdout, term.Alert())
	// Output:Problem:
}

func ExampleInform() {
	color.Enable = false
	fmt.Fprint(os.Stdout, term.Inform())
	// Output:Information:
}

func ExampleSecondary() {
	color.Enable = false
	fmt.Fprint(os.Stdout, term.Secondary("Hi"))
	// Output:Hi
}

func ExampleComment() {
	color.Enable = false
	fmt.Fprint(os.Stdout, term.Comment("Hi"))
	// Output:Hi
}

func ExampleFuzzy() {
	color.Enable = false
	fmt.Fprint(os.Stdout, term.Fuzzy("Hi"))
	// Output:Hi
}

func ExampleInfo() {
	color.Enable = false
	fmt.Fprint(os.Stdout, term.Info("Hi"))
	// Output:Hi
}

func ExampleBool() {
	fmt.Fprint(os.Stdout, term.Bool(true))
	fmt.Fprint(os.Stdout, term.Bool(false))
	// Output:✓✗
}

func ExampleOptions() {
	_, _ = term.Options(os.Stdout, "this is an example of a list of options",
		false, false, "option3", "option2", "option1")
	// Output:this is an example of a list of options.
	//   Options: option1, option2, option3
}

func ExampleBorder() {
	term.Border(os.Stdout, "hi")
	// Output: ┌────┐
	// │ hi │
	// └────┘
}

func TestTerm(t *testing.T) {
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotTerm := term.Term("", tt.name); gotTerm != tt.wantTerm {
				t.Errorf("Term() = %v, want %v", gotTerm, tt.wantTerm)
			}
		})
	}
}

func TestTerm16M(t *testing.T) {
	tests := []struct {
		name     string
		wantTerm string
	}{
		{"24bit", "terminal16m"},
		{"truecolor", "terminal16m"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotTerm := term.Term(tt.name, ""); gotTerm != tt.wantTerm {
				t.Errorf("Term() = %v, want %v", gotTerm, tt.wantTerm)
			}
		})
	}
}

func TestUnderlineChar(t *testing.T) {
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotS, err := term.UnderlineChar(tt.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("underlineChar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotS != tt.wantS {
				t.Errorf("underlineChar() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}

func TestUnderlineKeys(t *testing.T) {
	tests := []struct {
		name string
		keys []string
		want string
	}{
		{"empty", []string{}, ""},
		{"t", []string{"t"}, "\x1b[0m\x1b[4mt\x1b[0m"},
		{"test", []string{"test"}, "\x1b[0m\x1b[4mt\x1b[0mest"},
		{"tests", []string{"test1", "test2"}, "\x1b[0m\x1b[4mt\x1b[0mest1, \x1b[0m\x1b[4mt\x1b[0mest2"},
		{"😀", []string{"😀"}, "\x1b[0m\x1b[4m😀\x1b[0m"},
		{"test😀", []string{"test😀"}, "\x1b[0m\x1b[4mt\x1b[0mest😀"},
		{"😀test", []string{"😀test"}, "\x1b[0m\x1b[4m😀\x1b[0mtest"},
		{"file.min", []string{"file.min"}, "\x1b[0m\x1b[4mf\x1b[0mile.\x1b[0m\x1b[4mm\x1b[0min"},
		{"file.js.min", []string{"file.js.min"}, "\x1b[0m\x1b[4mf\x1b[0mile.js.\x1b[0m\x1b[4mm\x1b[0min"},
		{"📁.min", []string{"📁.min"}, "\x1b[0m\x1b[4m📁\x1b[0m.\x1b[0m\x1b[4mm\x1b[0min"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := term.UnderlineKeys(tt.keys...); got != tt.want {
				t.Errorf("UnderlineKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCenter(t *testing.T) {
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := term.Center(tt.args.width, tt.args.text); got != tt.want {
				t.Errorf("Center() = %q, want %q", got, tt.want)
			}
		})
	}
}
