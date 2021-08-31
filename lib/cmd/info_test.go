package cmd

import (
	"testing"
)

func Test_infoSample(t *testing.T) {
	tests := []struct {
		m            string
		name         string
		wantFilename bool
		wantErr      bool
	}{
		{"empty", "", false, false},
		{"invalid", "text/retrotxt.asc", false, false},
		{"empty", "ascii", true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFilename, err := infoSample(tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("infoSample() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if bool(len(gotFilename) > 0) != tt.wantFilename {
				t.Errorf("infoSample() = %v, want %v", gotFilename, tt.wantFilename)
			}
		})
	}
}
