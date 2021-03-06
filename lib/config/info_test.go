package config

import (
	"errors"
	"testing"
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
			if gotErr := Info(tt.style); !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("Info() = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}
