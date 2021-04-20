package config

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gookit/color"
)

func Example_updateBool() {
	color.Enable = false
	updateBool(false, "example")
	// Output: example is currently not in use
}

func Example_updateString() {
	color.Enable = false
	updateString("", "example", "")
	updateString("x", "save-directory", "")
	// Output: example is currently not in use
	//
	//   save-directory is set to "" ✗
}

func Example_recommend() {
	color.Enable = false
	fmt.Print(recommend(""))
	// Output: (suggestion: do not use)
}

func TestList(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"list", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := List(); (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSet(t *testing.T) {
	tests := []struct {
		name string
		args string
	}{
		{"empty", ""},
		{"0", "0"},
		{"valid", "editor"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Set(tt.name)
		})
	}
}

func Test_names_string(t *testing.T) {
	tests := []struct {
		name  string
		n     names
		theme bool
		want  string
	}{
		{"nil", nil, false, ""},
		{"empty", names{""}, false, ""},
		{"ok", names{"okay"}, false, "okay"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.string(tt.theme, ""); got != tt.want {
				t.Errorf("names.string() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dirExpansion(t *testing.T) {
	u, err := os.UserHomeDir()
	if err != nil {
		t.Error(err)
	}
	w, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name    string
		wantDir string
	}{
		{"", ""},
		{"~", u},
		{".", w},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDir := dirExpansion(tt.name); gotDir != tt.wantDir {
				t.Errorf("dirExpansion() = %v, want %v", gotDir, tt.wantDir)
			}
		})
	}
}

func Test_colorElm(t *testing.T) {
	// set test mode for str.HighlightWriter()
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
			if got := colorElm(tt.elm, "html", "bw", false); got != tt.want {
				t.Errorf("colorhtml() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColorCSS(t *testing.T) {
	tests := []struct {
		name string
		elm  string
		want string
	}{
		{"empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ColorCSS(tt.elm); got != tt.want {
				t.Errorf("ColorCSS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_previewPrompt(t *testing.T) {
	type args struct {
		name  string
		value string
	}
	tests := []struct {
		name  string
		args  args
		wantP string
	}{
		{"empty", args{}, "Set"},
		{"key", args{"html.meta.keywords", "ooooh"}, "Replace"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotP := previewPrompt(tt.args.name, tt.args.value)
			firstWord := strings.Split(strings.TrimSpace(gotP), " ")[0]
			if firstWord != tt.wantP {
				t.Errorf("previewPrompt() = %v, want %v", firstWord, tt.wantP)
			}
		})
	}
}
