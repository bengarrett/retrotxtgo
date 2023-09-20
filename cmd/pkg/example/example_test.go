package example_test

import (
	"strings"
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd/pkg/example"
	"github.com/stretchr/testify/assert"
)

func TestExample(t *testing.T) {
	t.Parallel()
	s := &strings.Builder{}
	example.Cmd.String(s)
	assert.Contains(t, s.String(), "print text files partial info")

	example.Info.String(s)
	assert.Contains(t, s.String(), "info file.txt")

	example.Examples.String(s)
	assert.Contains(t, s.String(), "list the builtin examples")

	example.Table.String(s)
	assert.Contains(t, s.String(), "iso-8859-15")

	example.View.String(s)
	assert.Contains(t, s.String(), "view file.txt")
	s.Reset()
}
