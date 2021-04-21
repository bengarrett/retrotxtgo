package cmd

import (
	"testing"
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
			if gotExit := configInfo(); gotExit != tt.wantExit {
				t.Errorf("configInfo() = %v, want %v", gotExit, tt.wantExit)
			}
		})
	}
}

func Test_configSet(t *testing.T) {
	tests := []struct {
		name string
		flag bool
		want bool
	}{
		{"dont list", false, false},
		{"list all", true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configFlag.configs = tt.flag
			if got := configSet(); got != tt.want {
				t.Errorf("configSet() = %v, want %v", got, tt.want)
			}
		})
	}
}
