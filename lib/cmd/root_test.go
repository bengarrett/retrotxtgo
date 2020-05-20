package cmd

import (
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/spf13/viper"
)

func TestInitDefaults(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want string
	}{
		{"empty", "", ""},
		{"layout", "create.layout", "standard"},
		{"save dir", "create.save-directory", ""},
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
