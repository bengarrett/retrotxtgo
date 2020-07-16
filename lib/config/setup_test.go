package config

import (
	"testing"

	"github.com/gookit/color"
)

func Test_hr(t *testing.T) {
	tests := []struct {
		name  string
		width uint
		want  string
	}{
		{"empty", 0, ""},
		{"5", 5, "-----"},
	}
	color.Disable()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hr(tt.width); got != tt.want {
				t.Errorf("hr() = %q, want %q", got, tt.want)
			}
		})
	}
}
