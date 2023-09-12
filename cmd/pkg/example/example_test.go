package example_test

import (
	"strings"
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd/pkg/example"
	"github.com/stretchr/testify/assert"
)

func TestExample(t *testing.T) {
	s := &strings.Builder{}
	example.Cmd.String(s)
	assert.Contains(t, s.String(), "print text files partial info")

	example.Info.String(s)
	assert.Contains(t, s.String(), "info file.txt")

	example.ListExamples.String(s)
	assert.Contains(t, s.String(), "list the builtin examples")

	example.ListTable.String(s)
	assert.Contains(t, s.String(), "iso-8859-15")

	example.List.String(s)
	assert.Contains(t, s.String(), "list codepages")

	example.View.String(s)
	assert.Contains(t, s.String(), "view file.txt")
	s.Reset()
}
