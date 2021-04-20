package cmd

import (
	"testing"
)

func Test_examples(t *testing.T) {
	t.Run("example", func(t *testing.T) {
		if got := examples(); got == nil {
			t.Errorf("examples() failed to return anything")
		}
	})
}
