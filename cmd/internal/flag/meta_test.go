package flag_test

import (
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
)

func TestSort(t *testing.T) {
	help := flag.Init()
	t.Run("sorter", func(t *testing.T) {
		got := flag.Sort(help)
		if len(got) == 0 {
			t.Error("Sort() returned an empty array")
			return
		}
		if i := got[0]; i != 0 {
			t.Errorf("Sort() returned %v, expected 0", i)
		}
		const wantLen = 18
		if len(got) != wantLen {
			t.Errorf("Sort() returned %v items, expected %d", len(got), wantLen)
		}
	})
}
