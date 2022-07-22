package config_test

import (
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/set"
	"github.com/gookit/color"
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

func TestSortKeys(t *testing.T) {
	got, keys := config.SortKeys(), set.Keys()
	lenG, lenK := len(got), len(keys)
	if lenG == 0 {
		t.Errorf("SortKeys() is empty")
		return
	}
	if lenG != lenK {
		t.Errorf("SortKeys() length = %d, want %d", lenG, lenK)
	}
	const want = get.FontFamily
	if s := got[0]; s != want {
		t.Errorf("SortKeys()[0] = %s, want %s", s, want)
	}
}

func TestTip(t *testing.T) {
	color.Enable = false
	tip := config.Tip()
	i := 0
	for key := range get.Reset() {
		i++
		s := tip[key]
		if s == "" {
			t.Errorf("tip[%s] is missing", key)
		}
	}
}
