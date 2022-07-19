package cmd_test

import (
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd"
	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/spf13/viper"
)

func TestConfigs_Command(t *testing.T) {
	tests := []struct {
		name string
		c    cmd.Configs
		want string
	}{
		{"create", cmd.Create, "create"},
		{"del", cmd.Delete, "delete"},
		{"edit", cmd.Edit, "edit"},
		{"info", cmd.Info, "info"},
		{"set", cmd.Set, "set [setting names]"},
		{"setup", cmd.Setup, "setup"},
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
			cmd.Config.Configs = tt.flag
			err := cmd.ListAll(nil)
			if err != nil {
				t.Error(err)
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
	config.InitDefaults()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := viper.GetString(tt.key); got != tt.want {
				t.Errorf("config.InitDefaults() %v = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}
