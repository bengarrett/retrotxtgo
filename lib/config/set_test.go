package config_test

import (
	"io"
	"os"
	"strconv"
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd"
	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
	"github.com/gookit/color"
)

func TestList(t *testing.T) {
	tests := []struct {
		name    string
		w       io.Writer
		wantErr bool
	}{
		{"nil output", nil, true},
		{"list okay", os.Stdout, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := config.List(tt.w); (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSet(t *testing.T) {
	color.Enable = false
	if err := cmd.LoadTester(os.Stdout); err != nil {
		t.Error(err)
	}
	tests := []struct {
		name    string
		w       io.Writer
		wantErr bool
	}{
		{"", nil, true},
		{"", os.Stdout, true},
		{"3", os.Stdout, false},
		{"html.layout", os.Stdout, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := config.Set(tt.w, tt.name); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSets(t *testing.T) {
	color.Enable = false
	if err := cmd.LoadTester(os.Stdout); err != nil {
		t.Error(err)
	}
	for key := range get.Reset() {
		t.Run(key, func(t *testing.T) {
			if err := config.Set(os.Stdout, key); err != nil {
				t.Errorf("Set(%s) error = %v", key, err)
			}
		})
	}
	sum := -1
	for key := range get.Reset() {
		sum++
		item := strconv.Itoa(sum)
		t.Run(item, func(t *testing.T) {
			if err := config.Set(os.Stdout, item); err != nil {
				t.Errorf("Set(%s) [%s] error = %v", item, key, err)
			}
		})
	}
}
