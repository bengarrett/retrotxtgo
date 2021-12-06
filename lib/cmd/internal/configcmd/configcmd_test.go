package configcmd_test

import (
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/configcmd"
	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/flag"
)

func Test_configListAll(t *testing.T) {
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
			flag.Config.Configs = tt.flag
			if got := configcmd.ListAll(); got != tt.want {
				t.Errorf("configListAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
