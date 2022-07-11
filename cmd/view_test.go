package cmd_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd"
	"github.com/gookit/color"
)

func tester(args []string) ([]byte, error) {
	color.Enable = false
	b := bytes.NewBufferString("")
	cmd := cmd.ViewInit()
	cmd.SetOut(b)
	cmd.SetArgs(args)
	cmd.Execute()
	out, err := ioutil.ReadAll(b)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var (
	cp1252 = []string{"-e", "1252", "../static/text/cp1252.txt"}
)

func Test_ExecuteCommand(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		wants []rune
	}{
		{"cp1252", cp1252, []rune{'‘', '’', '“', '”', '…', '™'}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i, want := range tt.wants {
				gotB, err := tester([]string{"-e", "1252", "../static/text/cp1252.txt"})
				if err != nil {
					t.Error(err)
					return
				}
				if !bytes.ContainsRune(gotB, want) {
					t.Errorf("%d. returned result doesn't include expected rune, %q got:\n%s", i, want, gotB)
					return
				}
			}
		})
	}
}
