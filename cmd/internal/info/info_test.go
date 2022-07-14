package info_test

import (
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd/internal/info"
)

func TestSample(t *testing.T) {
	tests := []struct {
		m            string
		name         string
		wantFilename bool
		wantErr      bool
	}{
		{"empty", "", false, true},
		{"invalid", "text/retrotxt.asc", false, true},
		{"logos", "ascii", true, false},
		{"ansi", "ansi", true, false},
		{"utf8", "utf8", true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFilename, err := info.Sample(tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sample() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if bool(len(gotFilename) > 0) != tt.wantFilename {
				t.Errorf("Sample() = %v, want %v", gotFilename, tt.wantFilename)
			}
		})
	}
}
