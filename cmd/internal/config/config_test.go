package config_test

import (
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd/internal/config"
	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	cfg "github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/spf13/viper"
)

func TestConfigs_Command(t *testing.T) {
	tests := []struct {
		name string
		c    config.Configs
		want string
	}{
		{"create", config.Create, "create"},
		{"del", config.Delete, "delete"},
		{"edit", config.Edit, "edit"},
		{"info", config.Info, "info"},
		{"set", config.Set, "set [setting names]"},
		{"setup", config.Setup, "setup"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Command().Use; got != tt.want {
				t.Errorf("Configs.Command() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestListAll(t *testing.T) {
	tests := []struct {
		name string
		flag bool
		want bool
	}{
		{"dont list", false, false},
		{"list all", true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.Config.Configs = tt.flag
			if got := config.ListAll(); got != tt.want {
				t.Errorf("ListAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
	cfg.InitDefaults()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := viper.GetString(tt.key); got != tt.want {
				t.Errorf("config.InitDefaults() %v = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}
