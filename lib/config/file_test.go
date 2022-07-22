package config_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd"
	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/gookit/color"
)

func TestPath(t *testing.T) {
	if err := cmd.LoadTester(os.Stdout); err != nil {
		t.Error(err)
	}
	t.Run("path test", func(t *testing.T) {
		if gotDir := config.Path(); gotDir == "tt.wantDir" {
			t.Errorf("Path() = \"\"")
		}
	})
}

func TestSetConfig(t *testing.T) {
	color.Enable = false
	if err := cmd.LoadTester(os.Stdout); err != nil {
		t.Error(err)
	}
	tests := []struct {
		name    string
		flag    string
		wantErr bool
	}{
		{"default", "", false},
		{"invalid", "this-file-doesnt-exist", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := config.SetConfig(os.Stdout, tt.flag); (err != nil) != tt.wantErr {
				t.Errorf("SetConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigMissing(t *testing.T) {
	color.Enable = false
	if err := cmd.LoadTester(os.Stdout); err != nil {
		t.Error(err)
	}
	t.Run("config missing", func(t *testing.T) {
		w := new(bytes.Buffer)
		config.ConfigMissing(w, "aaaxxx", "xxx")
		const want = "create it: aaa create --config="
		if gotW := w.String(); !strings.Contains(gotW, want) {
			t.Errorf("ConfigMissing() = %v, want %v", gotW, want)
		}
	})
}
