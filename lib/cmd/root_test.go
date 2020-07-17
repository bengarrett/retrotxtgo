package cmd

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"retrotxt.com/retrotxt/lib/config"
)

func TestInitDefaults(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name string
		key  string
		want string
	}{
		{"empty", "", ""},
		{"layout", "html.layout", "standard"},
		{"save dir", "save-directory", home},
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

func Test_examples(t *testing.T) {
	t.Run("example", func(t *testing.T) {
		if got := examples(); got == nil {
			t.Errorf("examples() failed to return anything")
		}
	})
}
