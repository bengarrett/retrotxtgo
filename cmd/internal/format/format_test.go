package format_test

import (
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd/internal/format"
	"github.com/nalgeon/be"
)

func TestFormat(t *testing.T) {
	t.Parallel()

	s := format.Format()
	be.Equal(t, len(s.Info), 5)
	be.Equal(t, s.Info[0], "color")
	be.Equal(t, s.Info[1], "json")
	be.Equal(t, s.Info[2], "json.min")
	be.Equal(t, s.Info[3], "text")
	be.Equal(t, s.Info[4], "xml")
}
