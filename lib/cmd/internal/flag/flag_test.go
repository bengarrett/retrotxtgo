package flag_test

import (
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/flag"
)

func Test_configInfo(t *testing.T) {
	tests := []struct {
		name     string
		wantExit bool
	}{
		{"output", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotExit := flag.ConfigInfo(); gotExit != tt.wantExit {
				t.Errorf("ConfigInfo() = %v, want %v", gotExit, tt.wantExit)
			}
		})
	}
}
