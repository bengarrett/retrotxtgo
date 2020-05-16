package cmd

import (
	"testing"
)

func Test_versionPrint(t *testing.T) {
	tests := []struct {
		name   string
		format string
		wantOk bool
	}{
		{"empty", "", true},
		{"invalid", "abcde", false},
		{"j", "j", true},
		{"jm", "jm", true},
		{"t", "t", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOk := versionPrint(tt.format); gotOk != tt.wantOk {
				t.Errorf("versionPrint() = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
