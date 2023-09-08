package xterm_test

import (
	"testing"

	"github.com/bengarrett/retrotxtgo/pkg/ansi/internal/bg"
	"github.com/bengarrett/retrotxtgo/pkg/ansi/internal/fg"
	"github.com/bengarrett/retrotxtgo/pkg/ansi/internal/xterm"
)

func TestForeground(t *testing.T) {
	tests := []struct {
		name string
		c    fg.Colors
		want xterm.Color
	}{
		{"black", fg.Black, 30},
		{"white", fg.White, 37},
		{"bcyan", fg.BrightCyan, 13},
		{"bwhite", fg.BrightWhite, 14},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := xterm.Foreground(tt.c); got != tt.want {
				t.Errorf("Foreground() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBackground(t *testing.T) {
	tests := []struct {
		name string
		c    bg.Colors
		want xterm.Color
	}{
		{"black", bg.Black, 40},
		{"white", bg.White, 47},
		{"bcyan", bg.BrightCyan, 13},
		{"bwhite", bg.BrightWhite, 14},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := xterm.Background(tt.c); got != tt.want {
				t.Errorf("Background() = %v, want %v", got, tt.want)
			}
		})
	}
}
