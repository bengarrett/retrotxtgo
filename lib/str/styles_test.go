package str

import (
	"fmt"
	"testing"
)

func ExampleBorder() {
	fmt.Printf("%s", Border("hi"))
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
			if gotTerm := Term("", tt.name); gotTerm != tt.wantTerm {
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
			if gotTerm := Term(tt.name, ""); gotTerm != tt.wantTerm {
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
			gotS, err := UnderlineChar(tt.c)
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
			if got := UnderlineKeys(tt.keys...); got != tt.want {
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
			if got := Center(tt.args.text, tt.args.width); got != tt.want {
				t.Errorf("Center() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNumberizeKeys(t *testing.T) {
	tests := []struct {
		name string
		keys []string
		want bool
	}{
		{"empty", nil, false},
		{"three", []string{"alpha", "bravo", "charlie"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NumberizeKeys(tt.keys...); bool(len(got) > 16) != tt.want {
				t.Errorf("NumberizeKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHighlight(t *testing.T) {
	type args struct {
		source string
		lexer  string
		style  string
		ansi   bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{"", "", "", false}, false},
		{"unknown", args{"abcde", "zxcvb", "sdfgh", false}, false},
		{"okay", args{"hello world", "json", "solarized-dark", true}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Highlight(tt.args.source, tt.args.lexer, tt.args.style, tt.args.ansi); (err != nil) != tt.wantErr {
				t.Errorf("Highlight() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValid(t *testing.T) {
	tests := []struct {
		name  string
		style string
		want  bool
	}{
		{"empty", "", false},
		{"err", "zxcv", false},
		{"okay", "rrt", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Valid(tt.style); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
