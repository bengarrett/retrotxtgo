package config

import (
	"testing"
)

func TestInfo(t *testing.T) {
	tests := []struct {
		name    string
		style   string
		wantErr string
	}{
		{"ok", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotErr := Info(tt.style); gotErr.Issue != tt.wantErr {
				t.Errorf("Info() = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}