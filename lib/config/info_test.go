package config_test

import (
	"errors"
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/config"
)

func TestInfo(t *testing.T) {
	tests := []struct {
		name    string
		style   string
		wantErr error
	}{
		{"ok", "", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, gotErr := config.Info(tt.style); !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("Info() = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}
