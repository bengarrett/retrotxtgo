package get_test

import (
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
	"github.com/bengarrett/retrotxtgo/meta"
)

func TestBool(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want bool
	}{
		{"bool", get.Genr, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := get.Bool(tt.key); got != tt.want {
				t.Errorf("Bool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want string
	}{
		{"string", get.LayoutTmpl, "standard"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := get.String(tt.key); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUInt(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want uint
	}{
		{"uint", "serve", meta.WebPort},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := get.UInt(tt.key); got != tt.want {
				t.Errorf("UInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
