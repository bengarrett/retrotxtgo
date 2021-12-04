package create_test

import (
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/create"
	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"golang.org/x/text/encoding/charmap"
)

func TestComment(t *testing.T) {
	const lf = "\x0A"
	r := []rune("hello world" + lf)
	args := create.Args{}
	args.Source.Name = "hi.txt"
	args.Source.Encoding = charmap.CodePage437

	const want = "encoding: IBM Code Page 437; line break: LF; length: 1; width: 11; name: hi.txt"

	t.Run("Comment", func(t *testing.T) {
		if got := args.Comment(filesystem.LF(), r...); got != want {
			t.Errorf("Comment() = %v, want %v", got, want)
		}
	})
}
