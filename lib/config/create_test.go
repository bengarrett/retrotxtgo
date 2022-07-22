package config_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd"
	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/gookit/color"
	"github.com/spf13/viper"
)

func TestCreate(t *testing.T) {
	if err := cmd.LoadTester(os.Stdout); err != nil {
		t.Error(err)
	}
	type args struct {
		name string
		ow   bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{}, true},
		{"overwrite", args{viper.ConfigFileUsed(), true}, false},
		{"no overwrite", args{viper.ConfigFileUsed(), false}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := config.Create(os.Stdout, tt.args.name, tt.args.ow); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNew(t *testing.T) {
	color.Enable = false
	if err := cmd.LoadTester(os.Stdout); err != nil {
		t.Error(err)
	}
	tests := []struct {
		name      string
		overwrite bool
		wantErr   bool
	}{
		{"ow false", false, true},
		{"ow", true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := config.New(w, tt.overwrite); (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDoesExist(t *testing.T) {
	color.Enable = false
	if err := cmd.LoadTester(os.Stdout); err != nil {
		t.Error(err)
	}
	type args struct {
		name   string
		suffix string
	}
	tests := []struct {
		name  string
		args  args
		wantW string
	}{
		{"empty", args{"", ""}, "Edit:\t edit"},
		{"empty", args{"abc-xyz", "-xyz"}, "Edit:\tabc edit"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			config.DoesExist(w, tt.args.name, tt.args.suffix)
			if gotW := w.String(); !strings.Contains(gotW, tt.wantW) {
				t.Errorf("DoesExist() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
