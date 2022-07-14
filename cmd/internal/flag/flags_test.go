package flag

import (
	"testing"
)

func TestCreate(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		got := Create()
		if ctrl := got.Controls; ctrl[0] != "eof" {
			t.Errorf("Create() returned %v for the Controls first time, expected eof", ctrl[0])
		}
	})
}
