package config_test

import (
	"os"
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd"
	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/gookit/color"
)

func TestSetup(t *testing.T) {
	color.Enable = false
	if err := cmd.LoadTester(os.Stdout); err != nil {
		t.Error(err)
	}
	tests := []struct {
		name    string
		start   int
		wantErr bool
	}{
		{"negative", -5, true},
		{"default", 0, false},
		{"skip items", 10, false},
		{"out of range", 100, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := config.Setup(os.Stdout, tt.start); (err != nil) != tt.wantErr {
				t.Errorf("Setup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
