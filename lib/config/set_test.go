package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func Test_dirExpansion(t *testing.T) {
	h, err := os.UserHomeDir()
	hp := filepath.Dir(h)
	w, err := os.Getwd()
	wp := filepath.Dir(w)
	s := string(os.PathSeparator)
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name    string
		wantDir string
	}{
		{"~", h},
		{filepath.Join("~", "foo"), filepath.Join(h, "foo")},
		{".", w},
		{fmt.Sprintf(".%sfoo", s), filepath.Join(w, "foo")},
		{fmt.Sprintf("..%sfoo", s), filepath.Join(wp, "foo")},
		{fmt.Sprintf("~%s..%sfoo", s, s), filepath.Join(hp, "foo")},
		{fmt.Sprintf("%sroot%sfoo%s..%sblah", s, s, s, s), fmt.Sprintf("root%sblah", s)},
		{fmt.Sprintf("%sroot%sfoo%s.%sblah", s, s, s, s), fmt.Sprintf("root%sfoo%sblah", s, s)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDir := dirExpansion(tt.name); gotDir != tt.wantDir {
				t.Errorf("dirExpansion(%v) = %v, want %v", tt.name, gotDir, tt.wantDir)
			}
		})
	}
}
