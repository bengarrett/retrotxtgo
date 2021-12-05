package config_test

import (
	"fmt"
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/gookit/color"
)

func Example_recommend() {
	color.Enable = false
	fmt.Print(config.Recommend(""))
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
			if err := config.List(); (err != nil) != tt.wantErr {
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
			config.Set(tt.name)
		})
	}
}

func Test_names_string(t *testing.T) {
	tests := []struct {
		name  string
		n     config.Names
		theme bool
		want  string
	}{
		{"nil", nil, false, ""},
		{"empty", config.Names{""}, false, ""},
		{"ok", config.Names{"okay"}, false, "okay"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.String(tt.theme, ""); got != tt.want {
				t.Errorf("Names.string() = %v, want %v", got, tt.want)
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
			if got := config.ColorElm(tt.elm, "html", "bw", false); got != tt.want {
				t.Errorf("ColorElm() = %v, want %v", got, tt.want)
			}
		})
	}
}
