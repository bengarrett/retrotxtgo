package color_test

import (
	"bytes"
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/config/internal/color"
)

// func TestCSS(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		elm  string
// 		want string
// 	}{
// 		{"empty", "", ""},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := color.CSS(tt.elm); got != tt.want {
// 				t.Errorf("CSS() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestElm(t *testing.T) {
// 	// set test mode for str.HighlightWriter()
// 	tests := []struct {
// 		name string
// 		elm  string
// 		want string
// 	}{
// 		{"empty", "", ""},
// 		{"str", "hello", "\nhello\n\n"},
// 		{"basic", "<h1>hello</h1>", "\n<h1>hello</h1>\n\n"},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := color.Elm(tt.elm, "html", "bw", false); got != tt.want {
// 				t.Errorf("Elm() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

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
