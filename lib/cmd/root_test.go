package cmd

import (
	"strings"
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/gookit/color"
	"github.com/spf13/viper"
)

func TestInitDefaults(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want string
	}{
		{"empty", "", ""},
		{"layout", "html.layout", "standard"},
		{"save dir", "save-directory", ""},
		{"style", "style.html", "lovelace"},
	}
	config.InitDefaults()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := viper.GetString(tt.key); got != tt.want {
				t.Errorf("config.InitDefaults() %v = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}

func Test_exampleCmd(t *testing.T) {
	color.Enable = false
	tests := []struct {
		name string
		tmpl string
		want string
	}{
		{"empty", "", ""},
		{"word", "Hello", "Hello"},
		{"words", "Hello world", "Hello world"},
		{"comment", "Hello # world", "Hello # world"},
		{"comments", "Hello # hash # world", "Hello # hash # world"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := exampleCmd(tt.tmpl); strings.TrimSpace(got) != tt.want {
				t.Errorf("exampleCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}
