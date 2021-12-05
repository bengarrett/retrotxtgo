package color_test

import (
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/config/internal/color"
)

func TestColorCSS(t *testing.T) {
	tests := []struct {
		name string
		elm  string
		want string
	}{
		{"empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := color.ColorCSS(tt.elm); got != tt.want {
				t.Errorf("ColorCSS() = %v, want %v", got, tt.want)
			}
		})
	}
}
