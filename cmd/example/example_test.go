package example_test

import (
	"strings"
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd/example"
	"github.com/nalgeon/be"
)

func TestExample(t *testing.T) {
	t.Parallel()
	s := &strings.Builder{}
	example.Cmd.String(s)
	find := strings.Contains(s.String(), "retrotxt info")
	be.True(t, find)
	example.Info.String(s)
	find = strings.Contains(s.String(), "info file.txt")
	be.True(t, find)
	example.Examples.String(s)
	find = strings.Contains(s.String(), "list the builtin examples")
	be.True(t, find)
	example.Table.String(s)
	find = strings.Contains(s.String(), "iso-8859-15")
	be.True(t, find)
	example.View.String(s)
	find = strings.Contains(s.String(), "view file.txt")
	be.True(t, find)
	s.Reset()
}
