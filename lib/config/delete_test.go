package config_test

import (
	"os"
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd"
	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/gookit/color"
)

func TestDelete(t *testing.T) {
	color.Enable = false
	tests := []struct {
		name    string
		ask     bool
		wantErr bool
	}{
		{"no prompt", false, false},
		{"prompt", true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cmd.LoadTester(os.Stdout); err != nil {
				t.Error(err)
			}
			if gotErr := config.Delete(os.Stdout, tt.ask); (gotErr != nil) != tt.wantErr {
				t.Errorf("Delete() = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}
