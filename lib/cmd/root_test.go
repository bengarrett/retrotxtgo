package cmd

import (
	"testing"

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
	logs.InitDefaults()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := viper.GetString(tt.key); got != tt.want {
				t.Errorf("logs.InitDefaults() %v = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}
