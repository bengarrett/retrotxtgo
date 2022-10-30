package config_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/gookit/color"
)

func TestInfo(t *testing.T) {
	color.Enable = false
	const success = "RetroTxtGo default settings in use"
	tests := []struct {
		name    string
		style   string
		wantW   string
		wantErr bool
	}{
		{"empty", "", success, false},
		{"none", "none", success, false},
		{"invalid style", "xxx", "unknown style \"xxx\", so using none", false},
		{"valid style", "vs", success, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := config.Info(w, tt.style); (err != nil) != tt.wantErr {
				t.Errorf("Info() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); !strings.Contains(gotW, tt.wantW) {
				t.Errorf("Info() does not contain %v", tt.wantW)
				fmt.Println(gotW)
			}
		})
	}
}

func TestAlert(t *testing.T) {
	color.Enable = false
	tests := []struct {
		name    string
		list    []string
		wantW   string
		wantErr bool
	}{
		{"empty", nil, "", false},
		{"one", []string{"x"}, "This setting is missing and should be configured", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := config.Alert(w, tt.list...); (err != nil) != tt.wantErr {
				t.Errorf("Alert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); !strings.Contains(gotW, tt.wantW) {
				t.Errorf("Alert() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
