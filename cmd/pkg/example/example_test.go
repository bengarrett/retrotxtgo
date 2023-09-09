package example_test

import (
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd/pkg/example"
	"github.com/stretchr/testify/assert"
)

func TestExample(t *testing.T) {
	assert.Contains(t, example.Cmd.String(), "print text files partial info")
	assert.Contains(t, example.Info.String(), "info file.txt")
	assert.Contains(t, example.ListExamples.String(), "list the builtin examples")
	assert.Contains(t, example.ListTable.String(), "iso-8859-15")
	assert.Contains(t, example.List.String(), "list codepages")
	assert.Contains(t, example.View.String(), "view file.txt")
}
