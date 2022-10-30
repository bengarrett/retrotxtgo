package config_test

import (
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
	"github.com/gookit/color"
)

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
