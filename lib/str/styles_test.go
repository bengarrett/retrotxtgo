package str

import (
	"fmt"
	"os"
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
			os.Setenv("TERM", tt.name)
			if gotTerm := Term(); gotTerm != tt.wantTerm {
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
			os.Setenv("COLORTERM", tt.name)
			if gotTerm := Term(); gotTerm != tt.wantTerm {
				t.Errorf("Term() = %v, want %v", gotTerm, tt.wantTerm)
			}
		})
	}
}

func Test_underlineRune(t *testing.T) {
	tests := []struct {
		name    string
		c       string
		wantS   string
		wantErr bool
	}{
		{"nil", "", "", false},
		{"ascii", "Z", "[0m[4mZ[0m", false},
		{"hex", "\x90", "", true},
		{"utf8", "\u005A", "[0m[4mZ[0m", false},
		{"░utf8", "\u2591", "[0m[4m░[0m", false},
		{"░hex", "\xe2\x96\x91", "[0m[4m░[0m", false},
		{"😀", "😀", "[0m[4m😀[0m", false},
		{"😀hex", "\xf0\x9f\x98\x80", "[0m[4m😀[0m", false},
		{"😀b", string([]byte{240, 159, 152, 128}), "[0m[4m😀[0m", false},
		{"string", "blahblah", "[0m[4mb[0m", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotS, err := underlineChar(tt.c)
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
		{"t", []string{"t"}, "[0m[4mt[0m"},
		{"test", []string{"test"}, "[0m[4mt[0mest"},
		{"tests", []string{"test1", "test2"}, "[0m[4mt[0mest1, [0m[4mt[0mest2"},
		{"😀", []string{"😀"}, "[0m[4m😀[0m"},
		{"test😀", []string{"test😀"}, "[0m[4mt[0mest😀"},
		{"😀test", []string{"😀test"}, "[0m[4m😀[0mtest"},
		{"file.min", []string{"file.min"}, "[0m[4mf[0mile.[0m[4mm[0min"},
		{"file.js.min", []string{"file.js.min"}, "[0m[4mf[0mile.js.[0m[4mm[0min"},
		{"📁.min", []string{"📁.min"}, "[0m[4m📁[0m.[0m[4mm[0min"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnderlineKeys(tt.keys); got != tt.want {
				t.Errorf("UnderlineKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}
