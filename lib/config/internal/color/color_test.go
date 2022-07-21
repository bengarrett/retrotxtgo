package color_test

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/config/internal/color"
)

func ExampleCSS() {
	//c.Enable = false
	if err := color.CSS(os.Stdout, "hello"); err != nil {
		log.Print(err)
	}
	// Output: [1m[38;5;254mhello[0m
}

func TestElm(t *testing.T) {
	tests := []struct {
		name string
		elm  string
		want string
	}{
		{"empty", "", ""},
		{"str", "hello", "\nhello\n\n"},
		{"basic", "<h1>hello</h1>", "\n<h1>hello</h1>\n\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := new(bytes.Buffer)
			err := color.Elm(b, tt.elm, "html", "bw", false)
			if b.String() != tt.want {
				t.Errorf("Elm() = %v, want %v", b.String(), tt.want)
			}
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestNames_string(t *testing.T) {
	tests := []struct {
		name  string
		n     color.Names
		theme bool
		want  string
	}{
		{"nil", nil, false, ""},
		{"empty", color.Names{""}, false, ""},
		{"one", color.Names{"okay"}, false, " 0 <okay=\"okay\">  \n"},
		{"two", color.Names{"hello", "world"}, false, " 0 <hello=\"hello\">  \n 1 <world=\"world\">\n\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := new(bytes.Buffer)
			tt.n.String(got, tt.theme, "")
			if got.String() != tt.want {
				t.Errorf("Names.string() = %q, want %q", got, tt.want)
			}
		})
	}
}
