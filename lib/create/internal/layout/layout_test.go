package layout_test

import (
	"reflect"
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/create/internal/layout"
)

func TestTemplates(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want string
	}{
		{"empty", "", "unknown"},
		{"none", "none", "none"},
		{"standard", "standard", "standard"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, _ := layout.ParseLayout(tt.key)
			if got := l.Pack(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseLayout() = %v, want %v", got, tt.want)
			}
		})
	}
}
