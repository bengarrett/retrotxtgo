package create

import (
	"testing"

	"golang.org/x/text/encoding/charmap"
	"retrotxt.com/retrotxt/lib/filesystem"
)

func TestComment(t *testing.T) {
	const lf = "\x0A"
	r := []rune("hello world" + lf)
	args := Args{}
	args.Source.Name = "hi.txt"
	args.Source.Encoding = charmap.CodePage437

	const want = "encoding: IBM Code Page 437; line break: LF; length: 1; width: 11; name: hi.txt"

	t.Run("comment", func(t *testing.T) {
		if got := args.comment(filesystem.LF(), r...); got != want {
			t.Errorf("AutoFont() = %v, want %v", got, want)
		}
	})
}
