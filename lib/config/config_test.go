package config_test

import (
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/set"
)

func TestMissing(t *testing.T) {
	tests := []struct {
		name      string
		wantCount int
	}{
		{"want no values", len(set.Keys())},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotList := len(config.Missing()); gotList != tt.wantCount {
				t.Errorf("Missing() = %v, want %v", gotList, tt.wantCount)
			}
		})
	}
}

func TestKeySort(t *testing.T) {
	got, keys := config.KeySort(), set.Keys()
	lenG, lenK := len(got), len(keys)
	if lenG == 0 {
		t.Errorf("KeySort() is empty")
		return
	}
	if lenG != lenK {
		t.Errorf("KeySort() length = %d, want %d", lenG, lenK)
	}
	want := get.FontFamily
	if s := got[0]; s != want {
		t.Errorf("KeySort()[0] = %s, want %s", s, want)
	}
}
