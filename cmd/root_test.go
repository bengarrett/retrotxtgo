package cmd_test

import (
	"bytes"
	"io/ioutil"

	"github.com/bengarrett/retrotxtgo/cmd"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

type cmdT int

const (
	listT cmdT = iota
	viewT
)

// tester initialises, runs and returns the results of the a cmd command.
// args are the command line arguments to test with the command.
func (t cmdT) tester(args []string) ([]byte, error) {
	color.Enable = false
	var c *cobra.Command
	b := bytes.NewBufferString("")
	switch t {
	case listT:
		c = cmd.ListInit()
	case viewT:
		c = cmd.ViewInit()
	default:
	}
	c.SetOut(b)
	c.SetArgs(args)
	c.Execute()
	out, err := ioutil.ReadAll(b)
	if err != nil {
		return nil, err
	}
	return out, nil
}
