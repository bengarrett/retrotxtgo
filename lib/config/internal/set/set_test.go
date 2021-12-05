package set_test

import (
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
