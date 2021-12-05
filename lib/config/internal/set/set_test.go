package set_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/config/internal/set"
)

func TestKeys(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"0 index", "editor"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := set.Keys()[0]; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirExpansion(t *testing.T) {
	u, err := os.UserHomeDir()
	if err != nil {
		t.Error(err)
	}
	w, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name    string
		wantDir string
	}{
		{"", ""},
		{"~", u},
		{".", w},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDir := set.DirExpansion(tt.name); gotDir != tt.wantDir {
				t.Errorf("DirExpansion() = %v, want %v", gotDir, tt.wantDir)
			}
		})
	}
}
